package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	_ sdk.Msg = &MsgSetOrchestratorAddress{}
	_ sdk.Msg = &MsgValsetConfirm{}
	_ sdk.Msg = &MsgSendToEth{}
	_ sdk.Msg = &MsgRequestBatch{}
	_ sdk.Msg = &MsgConfirmBatch{}
	_ sdk.Msg = &MsgDepositClaim{}
	_ sdk.Msg = &MsgWithdrawClaim{}
	_ sdk.Msg = &MsgFxOriginatedTokenClaim{}
)

// NewMsgSetOrchestratorAddress returns a new msgSetOrchestratorAddress
func NewMsgSetOrchestratorAddress(val sdk.ValAddress, oper sdk.AccAddress, eth string) *MsgSetOrchestratorAddress {
	return &MsgSetOrchestratorAddress{
		Validator:    val.String(),
		Orchestrator: oper.String(),
		EthAddress:   eth,
	}
}

// Route should return the name of the module
func (m *MsgSetOrchestratorAddress) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSetOrchestratorAddress) Type() string { return "set_operator_address" }

// ValidateBasic performs stateless checks
func (m *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(m.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Validator)
	}
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSetOrchestratorAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSetOrchestratorAddress) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(m.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// NewMsgValsetConfirm returns a new msgValsetConfirm
func NewMsgValsetConfirm(nonce uint64, ethAddress string, validator sdk.AccAddress, signature string) *MsgValsetConfirm {
	return &MsgValsetConfirm{
		Nonce:        nonce,
		Orchestrator: validator.String(),
		EthAddress:   ethAddress,
		Signature:    signature,
	}
}

// Route should return the name of the module
func (m *MsgValsetConfirm) Route() string { return RouterKey }

// Type should return the action
func (m *MsgValsetConfirm) Type() string { return "valset_confirm" }

// ValidateBasic performs stateless checks
func (m *MsgValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthAddress); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgSendToEth returns a new msgSendToEth
func NewMsgSendToEth(sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEth {
	return &MsgSendToEth{
		Sender:    sender.String(),
		EthDest:   destAddress,
		Amount:    send,
		BridgeFee: bridgeFee,
	}
}

// Route should return the name of the module
func (m MsgSendToEth) Route() string { return RouterKey }

// Type should return the action
func (m MsgSendToEth) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m MsgSendToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !m.BridgeFee.IsValid() || m.BridgeFee.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if err := ValidateEthAddress(m.EthDest); err != nil {
		return sdkerrors.Wrap(err, "ethereum dest address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatch returns a new msgRequestBatch
func NewMsgRequestBatch(orchestrator sdk.AccAddress, denom string, minimumFee sdk.Int, feeReceive string, baseFee sdk.Int) *MsgRequestBatch {
	return &MsgRequestBatch{
		Sender:     orchestrator.String(),
		Denom:      denom,
		MinimumFee: minimumFee,
		FeeReceive: feeReceive,
		BaseFee:    baseFee,
	}
}

// Route should return the name of the module
func (m MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (m MsgRequestBatch) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (m MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(ErrEmpty, fmt.Sprintf("denom is empty:%s", m.Denom))
	}
	if !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(ErrEmpty, "minimum fee is lg zero")
	}
	if err := ValidateEthAddressAndValidateChecksum(m.FeeReceive); err != nil {
		return sdkerrors.Wrap(ErrInvalid, fmt.Sprintf("err feeReceive address:%s", m.FeeReceive))
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgConfirmBatch returns a new msgConfirmBatch
func NewMsgConfirmBatch(nonce uint64, tokenContract, bscAddress, signature string, orchestrator sdk.AccAddress) *MsgConfirmBatch {
	return &MsgConfirmBatch{
		Nonce:         nonce,
		TokenContract: tokenContract,
		EthSigner:     bscAddress,
		Signature:     signature,
		Orchestrator:  orchestrator.String(),
	}
}

// Route should return the name of the module
func (m MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (m MsgConfirmBatch) Type() string { return "confirm_batch" }

// ValidateBasic performs stateless checks
func (m MsgConfirmBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthSigner); err != nil {
		return sdkerrors.Wrap(err, "eth signer")
	}
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "token contract")
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "signature is empty")
	}
	_, err := hex.DecodeString(m.Signature)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Could not decode hex string %s", m.Signature)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// the delegate address of the claimer, for MsgDepositClaim and MsgWithdrawClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ EthereumClaim = &MsgDepositClaim{}
	_ EthereumClaim = &MsgWithdrawClaim{}
	_ EthereumClaim = &MsgFxOriginatedTokenClaim{}
)

func NewMsgDepositClaim(eventNonce, blockHeight uint64, tokenContract string, amount sdk.Int, ethSender, fxReceiver,
	targetIbc, orchestrator string) *MsgDepositClaim {
	return &MsgDepositClaim{
		EventNonce:    eventNonce,
		BlockHeight:   blockHeight,
		TokenContract: tokenContract,
		Amount:        amount,
		EthSender:     ethSender,
		FxReceiver:    fxReceiver,
		TargetIbc:     targetIbc,
		Orchestrator:  orchestrator,
	}
}

// GetType returns the type of the claim
func (m *MsgDepositClaim) GetType() ClaimType {
	return CLAIM_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (m *MsgDepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FxReceiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.FxReceiver)
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthSender); err != nil {
		return sdkerrors.Wrap(err, "eth sender")
	}
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalid, "amount cannot be negative")
	}
	if _, err := hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not decode hex targetIbc string: %s", m.TargetIbc)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgDepositClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDepositClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgDepositClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(m.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (m MsgDepositClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (m MsgDepositClaim) Type() string { return "deposit_claim" }

// Route should return the name of the module
func (m MsgDepositClaim) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (m *MsgDepositClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/", m.TokenContract, m.EthSender, m.FxReceiver)
	return tmhash.Sum([]byte(path))
}

// GetType returns the claim type
func (m *MsgWithdrawClaim) GetType() ClaimType {
	return CLAIM_TYPE_WITHDRAW
}

// ValidateBasic performs stateless checks
func (m *MsgWithdrawClaim) ValidateBasic() error {
	if m.EventNonce == 0 {
		return fmt.Errorf("event_nonce == 0")
	}
	if m.BatchNonce == 0 {
		return fmt.Errorf("batch_nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	return nil
}

// Hash implements WithdrawBatch.Hash
func (m *MsgWithdrawClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", m.TokenContract, m.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (m MsgWithdrawClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgWithdrawClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgWithdrawClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(m.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (m MsgWithdrawClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (m MsgWithdrawClaim) Route() string { return RouterKey }

// Type should return the action
func (m MsgWithdrawClaim) Type() string { return "withdraw_claim" }

// NewMsgCancelSendToEth returns a new MsgCancelSendToEth
func NewMsgCancelSendToEth(sender sdk.AccAddress, id uint64) *MsgCancelSendToEth {
	return &MsgCancelSendToEth{
		Sender:        sender.String(),
		TransactionId: id,
	}
}

// Route should return the name of the module
func (m *MsgCancelSendToEth) Route() string { return RouterKey }

// Type should return the action
func (m *MsgCancelSendToEth) Type() string { return "cancel_send_to_eth" }

// ValidateBasic performs stateless checks
func (m *MsgCancelSendToEth) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(ErrInvalid, m.Sender)
	}
	if m.TransactionId == 0 {
		return sdkerrors.Wrap(ErrInvalid, "Transaction == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgCancelSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgCancelSendToEth) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

func (m *MsgFxOriginatedTokenClaim) Route() string {
	return RouterKey
}

func (m *MsgFxOriginatedTokenClaim) Type() string {
	return "fx_originated_token_claim"
}

func (m *MsgFxOriginatedTokenClaim) ValidateBasic() error {
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrap(err, "erc20 token")
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Orchestrator)
	}
	if m.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "token name is empty")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "token symbol is empty")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	return nil
}

func (m *MsgFxOriginatedTokenClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgFxOriginatedTokenClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgFxOriginatedTokenClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgFxOriginatedTokenClaim failed ValidateBasic! Should have been handled earlier")
	}
	val, _ := sdk.AccAddressFromBech32(m.Orchestrator)
	return val
}

func (m *MsgFxOriginatedTokenClaim) GetType() ClaimType {
	return CLAIM_TYPE_ORIGINATED_TOKEN
}

func (m *MsgFxOriginatedTokenClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%d/", m.TokenContract, m.Name, m.Symbol, m.Decimals)
	return tmhash.Sum([]byte(path))
}

// EthereumClaim implementation for MsgValsetUpdatedClaim
// ======================================================

// GetType returns the type of the claim
func (e *MsgValsetUpdatedClaim) GetType() ClaimType {
	return CLAIM_TYPE_VALSET_UPDATED
}

// ValidateBasic performs stateless checks
func (e *MsgValsetUpdatedClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(e.Orchestrator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, e.Orchestrator)
	}
	if len(e.Members) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "members len == 0")
	}
	for _, member := range e.Members {
		if err := ValidateEthAddress(member.EthAddress); err != nil {
			return sdkerrors.Wrap(ErrInvalid, err.Error())
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(ErrInvalid, "member power == 0")
		}
	}
	if e.EventNonce == 0 {
		return fmt.Errorf("nonce == 0")
	}
	if e.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgValsetUpdatedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgValsetUpdatedClaim) GetClaimer() sdk.AccAddress {
	err := msg.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, _ := sdk.AccAddressFromBech32(msg.Orchestrator)
	return val
}

// GetSigners defines whose signature is required
func (msg MsgValsetUpdatedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (msg MsgValsetUpdatedClaim) Type() string { return "Valset_Updated_Claim" }

// Route should return the name of the module
func (msg MsgValsetUpdatedClaim) Route() string { return RouterKey }

// Hash implements BridgeDeposit.Hash
func (b *MsgValsetUpdatedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%s/", b.ValsetNonce, b.EventNonce, b.BlockHeight, b.Members)
	return tmhash.Sum([]byte(path))
}
