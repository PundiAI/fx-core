package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func NewDelegateAmount(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(fxtypes.DefaultDenom, amount)
}

// --- ERC20Token --- //

func NewERC20Token(amount sdkmath.Int, contract string) ERC20Token {
	return ERC20Token{Amount: amount, Contract: contract}
}

// ValidateBasic permforms stateless validation
func (m *ERC20Token) ValidateBasic() error {
	if err := contract.ValidateEthereumAddress(m.Contract); err != nil {
		return ErrInvalid.Wrap("contract address")
	}
	if !m.Amount.IsPositive() {
		return ErrInvalid.Wrapf("amount")
	}
	return nil
}

type ERC20Tokens []ERC20Token

func (e ERC20Tokens) GetContracts() []gethcommon.Address {
	contracts := make([]gethcommon.Address, 0, len(e))
	for _, token := range e {
		contracts = append(contracts, gethcommon.HexToAddress(token.Contract))
	}
	return contracts
}

func (e ERC20Tokens) GetAmounts() []sdkmath.Int {
	amounts := make([]sdkmath.Int, 0, len(e))
	for _, token := range e {
		amounts = append(amounts, token.Amount)
	}
	return amounts
}

// --- BRIDGE VALIDATOR(S) --- //

// ValidateBasic performs stateless checks on validity
func (m *BridgeValidator) ValidateBasic() error {
	if m.Power == 0 {
		return ErrInvalid.Wrapf("power")
	}
	if err := contract.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return ErrInvalid.Wrapf("external address")
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
		return ErrInvalid.Wrapf("no members")
	}
	for i := range b {
		if err := b[i].ValidateBasic(); err != nil {
			return err
		}
	}
	if b.HasDuplicates() {
		return ErrInvalid.Wrapf("duplicate members")
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
		return nil, fmt.Errorf("parse gravity id: %w", err)
	}
	checkpoint, err := fxtypes.StrToByte32("checkpoint")
	if err != nil {
		return nil, err
	}

	memberAddresses := make([]gethcommon.Address, len(m.Members))
	convertedPowers := make([]*big.Int, len(m.Members))
	for i, m := range m.Members {
		memberAddresses[i] = toHexAddr(gravityIDStr, m.ExternalAddress)
		convertedPowers[i] = big.NewInt(int64(m.Power))
	}
	// the word 'checkpoint' needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	packBytes, err := contract.PackOracleSetCheckpoint(gravityID, checkpoint, big.NewInt(int64(m.Nonce)), memberAddresses, convertedPowers)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, fmt.Errorf("encode oracle set checkpoint: %w", err)
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	hash := crypto.Keccak256Hash(packBytes[4:])
	return hash.Bytes(), nil
}

func (m *OracleSet) Equal(o *OracleSet) (bool, error) {
	if m.Height != o.Height {
		return false, ErrInvalid.Wrapf("oracle set heights mismatch")
	}

	if m.Nonce != o.Nonce {
		return false, ErrInvalid.Wrapf("oracle set nonce mismatch")
	}

	if !BridgeValidators(m.Members).Equal(o.Members) {
		return false, ErrInvalid.Wrapf("oracle set members mismatch")
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

func (s OracleSets) Len() int {
	return len(s)
}

func (s OracleSets) Less(i, j int) bool {
	return s[i].Nonce > s[j].Nonce
}

func (s OracleSets) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// --- OutgoingTxBatch --- //

type OutgoingTxBatches []*OutgoingTxBatch

func (b OutgoingTxBatches) Len() int {
	return len(b)
}

func (b OutgoingTxBatches) Less(i, j int) bool {
	return b[i].BatchNonce > b[j].BatchNonce
}

func (b OutgoingTxBatches) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
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
		return nil, fmt.Errorf("parse gravity id: %w", err)
	}

	// Create the methodName argument which salts the signature
	batchMethodName, err := fxtypes.StrToByte32("transactionBatch")
	if err != nil {
		return nil, err
	}

	// Run through the elements of the batch and serialize them
	txAmounts := make([]*big.Int, len(m.Transactions))
	txDestinations := make([]gethcommon.Address, len(m.Transactions))
	txFees := make([]*big.Int, len(m.Transactions))
	for i, tx := range m.Transactions {
		txAmounts[i] = tx.Token.Amount.BigInt()
		txDestinations[i] = toHexAddr(gravityIDString, tx.DestAddress)
		txFees[i] = tx.Fee.Amount.BigInt()
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := contract.PackSubmitBatchCheckpoint(
		gravityID,
		batchMethodName,
		txAmounts,
		txDestinations,
		txFees,
		big.NewInt(int64(m.BatchNonce)),
		toHexAddr(gravityIDString, m.TokenContract),
		big.NewInt(int64(m.BatchTimeout)),
		toHexAddr(gravityIDString, m.FeeReceive),
	)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, fmt.Errorf("encode batch checkpoint: %w", err)
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

func (m *Oracle) GetSlashAmount(slashFraction sdkmath.LegacyDec) sdkmath.Int {
	slashAmount := sdkmath.LegacyNewDecFromInt(m.DelegateAmount).Mul(slashFraction).MulInt64(m.SlashTimes).TruncateInt()
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

func (o Oracles) Len() int {
	return len(o)
}

func (o Oracles) Less(i, j int) bool {
	return o[i].DelegateAmount.Sub(o[j].DelegateAmount).IsPositive()
}

func (o Oracles) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type OutgoingTransferTxs []*OutgoingTransferTx

func (txs OutgoingTransferTxs) TotalFee() sdkmath.Int {
	totalFee := sdkmath.NewInt(0)
	for _, tx := range txs {
		totalFee = totalFee.Add(tx.Fee.Amount)
	}
	return totalFee
}

// GetCheckpoint gets the checkpoint signature from the given outgoing bridge call
func (m *OutgoingBridgeCall) GetCheckpoint(gravityIDString string) ([]byte, error) {
	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := fxtypes.StrToByte32(gravityIDString)
	if err != nil {
		return nil, fmt.Errorf("parse gravity id: %w", err)
	}

	// Create the methodName argument which salts the signature
	bridgeCallMethodName, err := fxtypes.StrToByte32("bridgeCall")
	if err != nil {
		return nil, err
	}

	dataBytes, err := hex.DecodeString(m.Data)
	if err != nil {
		return nil, fmt.Errorf("parse data: %w", err)
	}
	memoBytes, err := hex.DecodeString(m.Memo)
	if err != nil {
		return nil, fmt.Errorf("parse memo: %w", err)
	}
	contracts := make([]gethcommon.Address, 0, len(m.Tokens))
	amounts := make([]*big.Int, 0, len(m.Tokens))
	for _, token := range m.Tokens {
		contracts = append(contracts, toHexAddr(gravityIDString, token.Contract))
		amounts = append(amounts, token.Amount.BigInt())
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := contract.PackBridgeCallCheckpoint(
		gravityID,
		bridgeCallMethodName,
		big.NewInt(int64(m.Nonce)),
		&contract.FxBridgeBaseBridgeCallData{
			Sender:     toHexAddr(gravityIDString, m.Sender),
			Refund:     toHexAddr(gravityIDString, m.Refund),
			Tokens:     contracts,
			Amounts:    amounts,
			To:         toHexAddr(gravityIDString, m.To),
			Data:       dataBytes,
			Memo:       memoBytes,
			Timeout:    big.NewInt(int64(m.Timeout)),
			GasLimit:   big.NewInt(int64(m.GasLimit)),
			EventNonce: big.NewInt(int64(m.EventNonce)),
		},
	)
	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, fmt.Errorf("encode bridge call checkpoint: %w", err)
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes(), nil
}

func NewBridgeDenom(moduleName, token string) string {
	return erc20types.NewBridgeDenom(moduleName, token)
}

func (m *MsgBridgeCallClaim) GetERC20Tokens() []ERC20Token {
	erc20Tokens := make([]ERC20Token, 0, len(m.TokenContracts))
	for i, tokenContract := range m.TokenContracts {
		erc20Tokens = append(erc20Tokens, ERC20Token{
			Contract: tokenContract,
			Amount:   m.Amounts[i],
		})
	}
	return erc20Tokens
}

func NewOriginTokenKey(moduleName string, id uint64) string {
	return fmt.Sprintf("%s/%d", moduleName, id)
}

func NewIBCTransferKey(ibcChannel string, ibcSequence uint64) string {
	return fmt.Sprintf("%s/%d", ibcChannel, ibcSequence)
}

func NewQuoteInfo(quote contract.IBridgeFeeQuoteQuoteInfo) QuoteInfo {
	return QuoteInfo{
		Id:       quote.Id.Uint64(),
		Token:    quote.TokenName,
		Fee:      sdkmath.NewIntFromBigInt(quote.Fee),
		Oracle:   quote.Oracle.Hex(),
		GasLimit: quote.GasLimit.Uint64(),
		Expiry:   quote.Expiry.Uint64(),
	}
}

func (q QuoteInfo) OracleAddress() gethcommon.Address {
	return gethcommon.HexToAddress(q.Oracle)
}

func toHexAddr(gravityId, addr string) gethcommon.Address {
	if gravityId == "fx-tron-bridge" || gravityId == "fx-tron-bridge-testnet" {
		tronAddr, err := address.Base58ToAddress(addr)
		if err != nil {
			panic(err)
		}
		return gethcommon.BytesToAddress(tronAddr.Bytes()[1:])
	}
	return gethcommon.HexToAddress(addr)
}
