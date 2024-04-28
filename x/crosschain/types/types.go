package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func NewDelegateAmount(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(OracleDelegateDenom, amount)
}

// --- ERC20Token --- //

func NewERC20Token(amount sdkmath.Int, contract string) ERC20Token {
	return ERC20Token{Amount: amount, Contract: contract}
}

// ValidateBasic permforms stateless validation
func (m *ERC20Token) ValidateBasic() error {
	if err := contract.ValidateEthereumAddress(m.Contract); err != nil {
		return errorsmod.Wrap(err, "invalid contract address")
	}
	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(ErrInvalid, "amount")
	}
	return nil
}

// --- BRIDGE VALIDATOR(S) --- //

// ValidateBasic performs stateless checks on validity
func (m *BridgeValidator) ValidateBasic() error {
	if m.Power == 0 {
		return errorsmod.Wrap(ErrEmpty, "power")
	}
	if err := contract.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return errorsmod.Wrap(ErrInvalid, "external address")
	}
	return nil
}

// BridgeValidators is the sorted set of validator data for Ethereum bridge MultiSig set
type BridgeValidators []BridgeValidator

func (b BridgeValidators) Len() int {
	return len(b)
}

func (b BridgeValidators) Less(i, j int) bool {
	if b[i].Power == b[j].Power {
		// Secondary sort on eth address in case powers are equal
		return bytes.Compare([]byte(b[i].ExternalAddress), []byte(b[j].ExternalAddress)) == -1
	}
	return b[i].Power > b[j].Power
}

func (b BridgeValidators) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
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
			return errorsmod.Wrapf(err, "member %d", i)
		}
	}
	if b.HasDuplicates() {
		return errorsmod.Wrap(ErrDuplicate, "address")
	}
	return nil
}

func (b BridgeValidators) Equal(o BridgeValidators) bool {
	if len(b) != len(o) {
		return false
	}

	for i, bv := range b {
		ov := o[i]
		if bv.Power != ov.Power || bv.ExternalAddress != ov.ExternalAddress {
			return false
		}
	}

	return true
}

// --- OracleSet(S) --- //

// NewOracleSet returns a new OracleSet
func NewOracleSet(nonce, height uint64, members BridgeValidators) *OracleSet {
	sort.Sort(members)
	return &OracleSet{
		Nonce:   nonce,
		Members: members,
		Height:  height,
	}
}

// GetCheckpoint returns the checkpoint
func (m *OracleSet) GetCheckpoint(gravityIDStr string) ([]byte, error) {
	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := fxtypes.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	checkpointBytes := []uint8("checkpoint")
	var checkpoint [32]uint8
	copy(checkpoint[:], checkpointBytes)

	memberAddresses := make([]gethcommon.Address, len(m.Members))
	convertedPowers := make([]*big.Int, len(m.Members))
	for i, m := range m.Members {
		memberAddresses[i] = gethcommon.HexToAddress(m.ExternalAddress)
		convertedPowers[i] = big.NewInt(int64(m.Power))
	}
	// the word 'checkpoint' needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	packBytes, packErr := contract.GetFxBridgeABI().Pack("oracleSetCheckpoint", gravityID, checkpoint, big.NewInt(int64(m.Nonce)), memberAddresses, convertedPowers)

	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if packErr != nil {
		return nil, errorsmod.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	hash := crypto.Keccak256Hash(packBytes[4:])
	return hash.Bytes(), nil
}

func (m *OracleSet) Equal(o *OracleSet) (bool, error) {
	if m.Height != o.Height {
		return false, errorsmod.Wrap(ErrInvalid, "oracle set heights mismatch")
	}

	if m.Nonce != o.Nonce {
		return false, errorsmod.Wrap(ErrInvalid, "oracle set nonce mismatch")
	}

	if !BridgeValidators(m.Members).Equal(o.Members) {
		return false, errorsmod.Wrap(ErrInvalid, "oracle set members mismatch")
	}

	return true, nil
}

func (m *OracleSet) GetTotalPower() uint64 {
	if m == nil {
		return 0
	}
	totalPower := uint64(0)
	for _, member := range m.Members {
		totalPower += member.Power
	}
	return totalPower
}

func (m *OracleSet) GetBridgePower(externalAddress string) (uint64, bool) {
	if m == nil {
		return 0, false
	}
	for _, member := range m.Members {
		if externalAddress == member.ExternalAddress {
			return member.Power, true
		}
	}
	return 0, false
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

// --- OutgoingTxBatch --- //

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
func (m *OutgoingTxBatch) GetFees() sdkmath.Int {
	sum := sdkmath.ZeroInt()
	for _, t := range m.Transactions {
		sum = sum.Add(t.Fee.Amount)
	}
	return sum
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (m *OutgoingTxBatch) GetCheckpoint(gravityIDString string) ([]byte, error) {
	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := fxtypes.StrToByte32(gravityIDString)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("transactionBatch")
	var batchMethodName [32]uint8
	copy(batchMethodName[:], methodNameBytes)

	// Run through the elements of the batch and serialize them
	txAmounts := make([]*big.Int, len(m.Transactions))
	txDestinations := make([]gethcommon.Address, len(m.Transactions))
	txFees := make([]*big.Int, len(m.Transactions))
	for i, tx := range m.Transactions {
		txAmounts[i] = tx.Token.Amount.BigInt()
		txDestinations[i] = gethcommon.HexToAddress(tx.DestAddress)
		txFees[i] = tx.Fee.Amount.BigInt()
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := contract.GetFxBridgeABI().Pack("submitBatchCheckpoint",
		gravityID,
		batchMethodName,
		txAmounts,
		txDestinations,
		txFees,
		big.NewInt(int64(m.BatchNonce)),
		gethcommon.HexToAddress(m.TokenContract),
		big.NewInt(int64(m.BatchTimeout)),
		gethcommon.HexToAddress(m.FeeReceive),
	)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, errorsmod.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes(), nil
}

// --- Oracle(S) --- //

func (m *Oracle) GetOracle() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.OracleAddress)
}

func (m *Oracle) GetBridger() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *Oracle) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(m.DelegateValidator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *Oracle) GetSlashAmount(slashFraction sdk.Dec) sdkmath.Int {
	slashAmount := sdk.NewDecFromInt(m.DelegateAmount).Mul(slashFraction).MulInt64(m.SlashTimes).TruncateInt()
	slashAmount = sdkmath.MinInt(slashAmount, m.DelegateAmount)
	slashAmount = sdkmath.MaxInt(slashAmount, sdkmath.ZeroInt())
	return slashAmount
}

func (m *Oracle) GetPower() sdkmath.Int {
	return m.DelegateAmount.Quo(sdk.DefaultPowerReduction)
}

func (m *Oracle) GetDelegateAddress(moduleName string) sdk.AccAddress {
	data := append(m.GetOracle(), []byte(moduleName)...)
	return crypto.Keccak256(data)[12:]
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

func MinBatchFeeToBaseFees(ms []MinBatchFee) map[string]sdkmath.Int {
	kv := make(map[string]sdkmath.Int, len(ms))
	for _, m := range ms {
		if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
			continue
		}
		kv[m.TokenContract] = m.BaseFee
	}
	return kv
}

type OutgoingTransferTxs []*OutgoingTransferTx

func (bs OutgoingTransferTxs) TotalFee() sdkmath.Int {
	totalFee := sdkmath.NewInt(0)
	for _, tx := range bs {
		totalFee = totalFee.Add(tx.Fee.Amount)
	}
	return totalFee
}

// GetCheckpoint gets the checkpoint signature from the given outgoing bridge call
func (m *OutgoingBridgeCall) GetCheckpoint(gravityIDString, chainName string) ([]byte, error) {
	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := fxtypes.StrToByte32(gravityIDString)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("bridgeCallCheckpoint")
	var batchMethodName [32]uint8
	copy(batchMethodName[:], methodNameBytes)

	messagesBytes, err := hex.DecodeString(m.Message)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse message")
	}
	contracts := make([]gethcommon.Address, 0, len(m.Tokens))
	amounts := make([]*big.Int, 0, len(m.Tokens))
	for _, token := range m.Tokens {
		contracts = append(contracts, gethcommon.HexToAddress(token.Contract))
		amounts = append(amounts, token.Amount.BigInt())
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := contract.GetFxBridgeABI().Pack("bridgeCallCheckpoint",
		gravityID,
		batchMethodName,
		gethcommon.HexToAddress(m.Sender),
		gethcommon.HexToAddress(m.To),
		gethcommon.HexToAddress(m.Receiver),
		m.Value.BigInt(),
		big.NewInt(int64(m.Nonce)),
		big.NewInt(int64(m.GasLimit)),
		big.NewInt(int64(m.Timeout)),
		messagesBytes,
		contracts,
		amounts,
	)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, errorsmod.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes(), nil
}

func (m *SnapshotOracle) HasExternalAddress(address string) bool {
	if m == nil {
		return false
	}
	for _, member := range m.Members {
		if address == member.ExternalAddress {
			return true
		}
	}
	return false
}

func (m *SnapshotOracle) GetExternalAddressPower(address string) uint64 {
	if m == nil {
		return 0
	}
	for _, member := range m.Members {
		if address == member.ExternalAddress {
			return member.Power
		}
	}
	return 0
}

func (m *SnapshotOracle) GetTotalPower() uint64 {
	if m == nil {
		return 0
	}
	totalPower := uint64(0)
	for _, member := range m.Members {
		totalPower = totalPower + member.Power
	}
	return totalPower
}

func NewPendingOutgoingTx(txID uint64, sender sdk.AccAddress, receiver string, tokenContract string, amount, fee sdk.Coin, rewawrds sdk.Coins) PendingOutgoingTransferTx {
	return PendingOutgoingTransferTx{
		Id:            txID,
		Sender:        sender.String(),
		DestAddress:   receiver,
		TokenContract: tokenContract,
		Token:         amount,
		Fee:           fee,
		Rewards:       rewawrds,
	}
}

func NewSnapshotOracle(oracleSet *OracleSet, nonce uint64) *SnapshotOracle {
	return &SnapshotOracle{
		Nonces:         []uint64{nonce},
		OracleSetNonce: oracleSet.Nonce,
		Members:        oracleSet.Members,
	}
}

func (m *SnapshotOracle) AppendNonce(nonce uint64) *SnapshotOracle {
	m.Nonces = append(m.Nonces, nonce)
	return m
}

func (m *SnapshotOracle) RemoveNonce(nonce uint64) *SnapshotOracle {
	for i, n := range m.Nonces {
		if n == nonce {
			m.Nonces = append(m.Nonces[:i], m.Nonces[i+1:]...)
			break
		}
	}
	return m
}

func NewBridgeDenom(moduleName string, token string) string {
	return fmt.Sprintf("%s%s", moduleName, token)
}

func ExternalAddressToAccAddress(chainName, addr string) sdk.AccAddress {
	router, ok := msgValidateBasicRouter[chainName]
	if !ok {
		panic("unrecognized cross chain name")
	}
	accAddr, err := router.ExternalAddressToAccAddress(addr)
	if err != nil {
		panic(err)
	}
	return accAddr
}

func NewERC20Tokens(module string, tokenAddrs []gethcommon.Address, tokenAmounts []*big.Int) ([]ERC20Token, error) {
	if len(tokenAddrs) != len(tokenAmounts) {
		return nil, fmt.Errorf("invalid length")
	}
	tokens := make([]ERC20Token, 0)
	for i := 0; i < len(tokenAddrs); i++ {
		contractAddr := fxtypes.AddressToStr(tokenAddrs[i].Bytes(), module)
		amount := sdkmath.NewIntFromBigInt(tokenAmounts[i])
		found := false
		for j := 0; j < len(tokens); j++ {
			if contractAddr == tokens[j].Contract {
				tokens[j].Amount = tokens[j].Amount.Add(amount)
				found = true
				break
			}
		}
		if !found {
			tokens = append(tokens, ERC20Token{
				Contract: contractAddr,
				Amount:   amount,
			})
		}
	}
	return tokens, nil
}
