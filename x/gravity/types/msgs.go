package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	TypeMsgSetOrchestratorAddress = "set_operator_address"
	TypeMsgValsetConfirm          = "valset_confirm"
	TypeMsgSendToEth              = "send_to_eth"
	TypeMsgRequestBatch           = "request_batch"
	TypeMsgConfirmBatch           = "confirm_batch"
	TypeMsgDepositClaim           = "deposit_claim"
	TypeMsgWithdrawClaim          = "withdraw_claim"
	TypeMsgFxOriginatedTokenClaim = "fx_originated_token_claim"
	TypeMsgCancelSendToEth        = "cancel_send_to_eth"
	TypeMsgValsetUpdatedClaim     = "Valset_Updated_Claim"
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
	_ sdk.Msg = &MsgCancelSendToEth{}
	_ sdk.Msg = &MsgValsetUpdatedClaim{}
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
func (m *MsgSetOrchestratorAddress) Type() string { return TypeMsgSetOrchestratorAddress }

// ValidateBasic performs stateless checks
func (m *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(m.Validator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "validator: %s, err: %s", m.Validator, err.Error())
	}
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "ethereum address: %s, err: %s", m.EthAddress, err.Error())
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
func (m *MsgValsetConfirm) Type() string { return TypeMsgValsetConfirm }

// ValidateBasic performs stateless checks
func (m *MsgValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "ethereum address: %s, err: %s", m.EthAddress, err.Error())
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(err, "could not hex decode signature: %s", m.Signature)
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
func (m MsgSendToEth) Type() string { return TypeMsgSendToEth }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m MsgSendToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s, err: %s", m.Sender, err.Error())
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "amount: %s", m.Amount)
	}
	if !m.BridgeFee.IsValid() || m.BridgeFee.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "bridge fee: %s", m.BridgeFee)
	}
	if err := ValidateEthAddress(m.EthDest); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "ethereum dest address: %s, err: %s", m.EthDest, err.Error())
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
		BaseFee:    &baseFee,
	}
}

// Route should return the name of the module
func (m MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (m MsgRequestBatch) Type() string { return TypeMsgRequestBatch }

// ValidateBasic performs stateless checks
func (m MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s, err: %s", m.Sender, err.Error())
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(ErrEmpty, "denom is empty")
	}
	if !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(ErrEmpty, "minimum fee is lg zero")
	}
	if err := ValidateEthAddressAndValidateChecksum(m.FeeReceive); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "fee receive: %s, err: %s", m.FeeReceive, err.Error())
	}
	if m.BaseFee == nil || m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidRequestBatchBaseFee, m.BaseFee.String())
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
func (m MsgConfirmBatch) Type() string { return TypeMsgConfirmBatch }

// ValidateBasic performs stateless checks
func (m MsgConfirmBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthSigner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "eth signer: %s, err: %s", m.EthSigner, err.Error())
	}
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "token: %s, err: %s", m.TokenContract, err.Error())
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "signature is empty")
	}
	_, err := hex.DecodeString(m.Signature)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not decode hex string %s, err: %s", m.Signature, err.Error())
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

// Route should return the name of the module
func (m *MsgCancelSendToEth) Route() string { return RouterKey }

// Type should return the action
func (m *MsgCancelSendToEth) Type() string { return TypeMsgCancelSendToEth }

// ValidateBasic performs stateless checks
func (m *MsgCancelSendToEth) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %s, err: %s", m.Sender, err.Error())
	}
	if m.TransactionId == 0 {
		return sdkerrors.Wrap(ErrEmpty, "transaction == 0")
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
	return []sdk.AccAddress{acc}
}

// EthereumClaim represents a claim on ethereum state
type EthereumClaim interface {
	// GetEventNonce All Ethereum claims that we relay from the Gravity contract and into the module
	// have a nonce that is monotonically increasing and unique, since this nonce is
	// issued by the Ethereum contract it is immutable and must be agreed on by all validators
	// any disagreement on what claim goes to what nonce means someone is lying.
	GetEventNonce() uint64
	// GetBlockHeight The block height that the claimed event occurred on. This EventNonce provides sufficient
	// ordering for the execution of all claims. The block height is used only for batchTimeouts + logicTimeouts
	// when we go to create a new batch we set the timeout some number of batches out from the last
	// known height plus projected block progress since then.
	GetBlockHeight() uint64
	// GetClaimer the delegate address of the claimer, for MsgDepositClaim and MsgWithdrawClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// GetType Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ EthereumClaim = &MsgDepositClaim{}
	_ EthereumClaim = &MsgWithdrawClaim{}
	_ EthereumClaim = &MsgFxOriginatedTokenClaim{}
	_ EthereumClaim = &MsgValsetUpdatedClaim{}
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "fx receiver: %s, err: %s", m.FxReceiver, err.Error())
	}
	if err := ValidateEthAddressAndValidateChecksum(m.EthSender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "eth sender: %s, err: %s", m.EthSender, err.Error())
	}
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "token: %s, err: %s", m.TokenContract, err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
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
	if err := m.ValidateBasic(); err != nil {
		panic("MsgDepositClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
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
func (m MsgDepositClaim) Type() string { return TypeMsgDepositClaim }

// Route should return the name of the module
func (m MsgDepositClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeDeposit.Hash
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
		return sdkerrors.Wrap(ErrEmpty, "event nonce == 0")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.Wrap(ErrEmpty, "batch_nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrEmpty, "block height == 0")
	}
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "token: %s, err: %s", m.TokenContract, err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
	}
	return nil
}

// ClaimHash Hash implements WithdrawBatch.Hash
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
	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
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
func (m MsgWithdrawClaim) Type() string { return TypeMsgWithdrawClaim }

// NewMsgCancelSendToEth returns a new MsgCancelSendToEth
func NewMsgCancelSendToEth(sender sdk.AccAddress, id uint64) *MsgCancelSendToEth {
	return &MsgCancelSendToEth{
		Sender:        sender.String(),
		TransactionId: id,
	}
}

func (m *MsgFxOriginatedTokenClaim) Route() string {
	return RouterKey
}

func (m *MsgFxOriginatedTokenClaim) Type() string { return TypeMsgFxOriginatedTokenClaim }

func (m *MsgFxOriginatedTokenClaim) ValidateBasic() error {
	if err := ValidateEthAddressAndValidateChecksum(m.TokenContract); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "token: %s, err: %s", m.TokenContract, err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrEmpty, "event nonce == 0")
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "token name is empty")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "token symbol is empty")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrEmpty, "block height == 0")
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
	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
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
func (m *MsgValsetUpdatedClaim) GetType() ClaimType {
	return CLAIM_TYPE_VALSET_UPDATED
}

// ValidateBasic performs stateless checks
func (m *MsgValsetUpdatedClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "orchestrator: %s, err: %s", m.Orchestrator, err.Error())
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "members len == 0")
	}
	for _, member := range m.Members {
		if err := ValidateEthAddress(member.EthAddress); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "eth address: %s, err: %s", member.EthAddress, err.Error())
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(ErrEmpty, "member power == 0")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrEmpty, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrEmpty, "block height == 0")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgValsetUpdatedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgValsetUpdatedClaim) GetClaimer() sdk.AccAddress {
	err := m.ValidateBasic()
	if err != nil {
		panic("MsgERC20DeployedClaim failed ValidateBasic! Should have been handled earlier")
	}

	val, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}
	return val
}

// GetSigners defines whose signature is required
func (m MsgValsetUpdatedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Orchestrator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// Type should return the action
func (m MsgValsetUpdatedClaim) Type() string { return TypeMsgValsetUpdatedClaim }

// Route should return the name of the module
func (m MsgValsetUpdatedClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeDeposit.Hash
func (m *MsgValsetUpdatedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%s/", m.ValsetNonce, m.EventNonce, m.BlockHeight, m.Members)
	return tmhash.Sum([]byte(path))
}
