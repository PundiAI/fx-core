package types

import (
	math "math"
	"math/big"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//////////////////////////////////////
//      BRIDGE VALIDATOR(S)         //
//////////////////////////////////////

// ValidateBasic performs stateless checks on validity
func (m BridgeValidator) ValidateBasic() error {
	if m.Power == 0 {
		return sdkerrors.Wrap(ErrEmpty, "power")
	}
	if err := ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "external address")
	}
	return nil
}

// BridgeValidators is the sorted set of validator data for Ethereum bridge MultiSig set
type BridgeValidators []BridgeValidator

// Sort sorts the validators by power
func (b BridgeValidators) Sort() {
	sort.Slice(b, func(i, j int) bool {
		if b[i].Power == b[j].Power {
			// Secondary sort on eth address in case powers are equal
			return ethereumAddrLessThan(b[i].ExternalAddress, b[j].ExternalAddress)
		}
		return b[i].Power > b[j].Power
	})
}

// PowerDiff returns the difference in power between two bridge validator sets
// note this is Gravity bridge power *not* Cosmos voting power. Cosmos voting
// power is based on the absolute number of tokens in the staking pool at any given
// time Gravity bridge power is normalized using the equation.
//
// validators cosmos voting power / total cosmos voting power in this block = gravity bridge power / u32_max
//
// As an example if someone has 52% of the Cosmos voting power when a validator set is created their Gravity
// bridge voting power is u32_max * .52
//
// Normalized voting power dramatically reduces how often we have to produce new validator set updates. For example
// if the total on chain voting power increases by 1% due to inflation, we shouldn't have to generate a new validator
// set, after all the validators retained their relative percentages during inflation and normalized Gravity bridge power
// shows no difference.
func (b BridgeValidators) PowerDiff(c BridgeValidators) float64 {
	powers := map[string]int64{}
	// loop over b and initialize the map with their powers
	for _, bv := range b {
		powers[bv.ExternalAddress] = int64(bv.Power)
	}

	// subtract c powers from powers in the map, initializing
	// uninitialized keys with negative numbers
	for _, bv := range c {
		if val, ok := powers[bv.ExternalAddress]; ok {
			powers[bv.ExternalAddress] = val - int64(bv.Power)
		} else {
			powers[bv.ExternalAddress] = -int64(bv.Power)
		}
	}

	var delta float64
	for _, v := range powers {
		// NOTE: we care about the absolute value of the changes
		delta += math.Abs(float64(v))
	}

	return math.Abs(delta / float64(math.MaxUint32))
}

// TotalPower returns the total power in the bridge validator set
func (b BridgeValidators) TotalPower() (out uint64) {
	for _, v := range b {
		out += v.Power
	}
	return
}

// HasDuplicates returns true if there are duplicates in the set
func (b BridgeValidators) HasDuplicates() bool {
	m := make(map[string]struct{}, len(b))
	for i := range b {
		m[b[i].ExternalAddress] = struct{}{}
	}
	return len(m) != len(b)
}

// GetPowers returns only the power values for all members
func (b BridgeValidators) GetPowers() []uint64 {
	r := make([]uint64, len(b))
	for i := range b {
		r[i] = b[i].Power
	}
	return r
}

// ValidateBasic performs stateless checks
func (b BridgeValidators) ValidateBasic() error {
	if len(b) == 0 {
		return ErrEmpty
	}
	for i := range b {
		if err := b[i].ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "member %d", i)
		}
	}
	if b.HasDuplicates() {
		return sdkerrors.Wrap(ErrDuplicate, "address")
	}
	return nil
}

//////////////////////////////////////
//          OracleSet(S)            //
//////////////////////////////////////

// NewOracleSet returns a new OracleSet
func NewOracleSet(nonce, height uint64, members BridgeValidators) *OracleSet {
	members.Sort()
	var mem []BridgeValidator
	for _, val := range members {
		mem = append(mem, val)
	}
	return &OracleSet{
		Nonce:   nonce,
		Members: mem,
		Height:  height,
	}
}

// GetCheckpoint returns the checkpoint
func (m OracleSet) GetCheckpoint(gravityIDStr string) ([]byte, error) {
	// error case here should not occur outside of testing since the above is a constant
	contractAbi, err := abi.JSON(strings.NewReader(OracleSetCheckpointABIJSON))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}

	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := StrToFixByteArray(gravityIDStr)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "parse gravity id")
	}
	checkpointBytes := []uint8("checkpoint")
	var checkpoint [32]uint8
	copy(checkpoint[:], checkpointBytes[:])

	memberAddresses := make([]gethcommon.Address, len(m.Members))
	convertedPowers := make([]*big.Int, len(m.Members))
	for i, m := range m.Members {
		memberAddresses[i] = gethcommon.HexToAddress(m.ExternalAddress)
		convertedPowers[i] = big.NewInt(int64(m.Power))
	}
	// the word 'checkpoint' needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	bytes, packErr := contractAbi.Pack("checkpoint", gravityID, checkpoint, big.NewInt(int64(m.Nonce)), memberAddresses, convertedPowers)

	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if packErr != nil {
		return nil, sdkerrors.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	hash := crypto.Keccak256Hash(bytes[4:])
	return hash.Bytes(), nil
}

type OracleSets []*OracleSet

func (v OracleSets) Len() int {
	return len(v)
}

func (v OracleSets) Less(i, j int) bool {
	return v[i].Nonce > v[j].Nonce
}

func (v OracleSets) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

//////////////////////////////////////
//         OutgoingTxBatch          //
//////////////////////////////////////

type OutgoingTxBatches []*OutgoingTxBatch

func (v OutgoingTxBatches) Len() int {
	return len(v)
}

func (v OutgoingTxBatches) Less(i, j int) bool {
	return v[i].BatchNonce > v[j].BatchNonce
}

func (v OutgoingTxBatches) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// GetFees returns the total fees contained within a given batch
func (m OutgoingTxBatch) GetFees() sdk.Int {
	sum := sdk.ZeroInt()
	for _, t := range m.Transactions {
		sum = sum.Add(t.Fee.Amount)
	}
	return sum
}

//////////////////////////////////////
//            Oracle(S)             //
//////////////////////////////////////

func (m Oracle) GetOracle() sdk.AccAddress {
	// oracle address can't be empty
	addr, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (m Oracle) GetPower() sdk.Int {
	if m.IsValidator {
		return m.DelegateAmount
	}
	return m.DelegateAmount.Quo(sdk.DefaultPowerReduction)
}

type Oracles []Oracle

func (v Oracles) Len() int {
	return len(v)
}

func (v Oracles) Less(i, j int) bool {
	return v[i].DelegateAmount.Sub(v[j].DelegateAmount).IsPositive()
}

func (v Oracles) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func MinBatchFeeToBaseFees(ms []MinBatchFee) map[string]sdk.Int {
	kv := make(map[string]sdk.Int, len(ms))
	for _, m := range ms {
		if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
			continue
		}
		kv[m.TokenContract] = m.BaseFee
	}
	return kv
}

func CovertIbcPacketReceiveAddressByPrefix(targetIbcPrefix string, receiver sdk.AccAddress) (ibcReceiveAddr string, err error) {
	if strings.ToLower(targetIbcPrefix) == "0x" {
		return gethcommon.BytesToAddress(receiver.Bytes()).String(), nil
	}
	return bech32.ConvertAndEncode(targetIbcPrefix, receiver)
}

func GetOracleDelegateAddress(moduleName string, oracleAddr sdk.AccAddress) sdk.AccAddress {
	data := append(oracleAddr, []byte(moduleName)...)
	return crypto.Keccak256(data)[12:]
}
