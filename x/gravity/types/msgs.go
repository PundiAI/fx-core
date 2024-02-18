package types

import (
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

const (
	ModuleName = "gravity"
	RouterKey  = ModuleName

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

// Route should return the name of the module
func (m *MsgSetOrchestratorAddress) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSetOrchestratorAddress) Type() string { return TypeMsgSetOrchestratorAddress }

// ValidateBasic performs stateless checks
func (m *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(m.Validator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	if err = fxtypes.ValidateEthereumAddress(m.EthAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid ethereum address: %s", err)
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

// Route should return the name of the module
func (m *MsgValsetConfirm) Route() string { return RouterKey }

// Type should return the action
func (m *MsgValsetConfirm) Type() string { return TypeMsgValsetConfirm }

// ValidateBasic performs stateless checks
func (m *MsgValsetConfirm) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	if err = fxtypes.ValidateEthereumAddress(m.EthAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid ethereum address: %s", err)
	}
	if len(m.Signature) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgValsetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgValsetConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Orchestrator)}
}

// Route should return the name of the module
func (m *MsgSendToEth) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendToEth) Type() string { return TypeMsgSendToEth }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m *MsgSendToEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return errortypes.ErrInvalidRequest.Wrap("bridge fee denom not equal amount denom")
	}
	if !m.Amount.IsValid() || m.Amount.IsZero() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	if !m.BridgeFee.IsValid() || m.BridgeFee.IsZero() {
		return errortypes.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	if err := fxtypes.ValidateEthereumAddress(m.EthDest); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid dest address: %s", err)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSendToEth) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// Route should return the name of the module
func (m *MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRequestBatch) Type() string { return TypeMsgRequestBatch }

// ValidateBasic performs stateless checks
func (m *MsgRequestBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.Denom) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty denom")
	}
	if m.MinimumFee.IsNil() || !m.MinimumFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid minimum fee")
	}
	if err := fxtypes.ValidateEthereumAddress(m.FeeReceive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid fee receive address: %s", err)
	}
	if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid base fee")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgRequestBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// Route should return the name of the module
func (m *MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConfirmBatch) Type() string { return TypeMsgConfirmBatch }

// ValidateBasic performs stateless checks
func (m *MsgConfirmBatch) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	if err := fxtypes.ValidateEthereumAddress(m.EthSigner); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid eth signer address: %s", err)
	}
	if err := fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract address: %s", err)
	}
	if len(m.Signature) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err := hex.DecodeString(m.Signature); err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Orchestrator)}
}

// Route should return the name of the module
func (m *MsgCancelSendToEth) Route() string { return RouterKey }

// Type should return the action
func (m *MsgCancelSendToEth) Type() string { return TypeMsgCancelSendToEth }

// ValidateBasic performs stateless checks
func (m *MsgCancelSendToEth) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero transaction id")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgCancelSendToEth) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgCancelSendToEth) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// Deprecated: EthereumClaim represents a claim on ethereum state
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

func UnpackAttestationClaim(cdc codectypes.AnyUnpacker, att *Attestation) (EthereumClaim, error) {
	var msg EthereumClaim
	err := cdc.UnpackAny(att.Claim, &msg)
	return msg, err
}

// GetType returns the type of the claim
func (m *MsgDepositClaim) GetType() ClaimType {
	return CLAIM_TYPE_DEPOSIT
}

// ValidateBasic performs stateless checks
func (m *MsgDepositClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FxReceiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid fx receiver address: %s", err)
	}
	if err := fxtypes.ValidateEthereumAddress(m.EthSender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid eth sender address: %s", err)
	}
	if err := fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	if _, err := hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not decode hex targetIbc")
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgDepositClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgDepositClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.Orchestrator)
}

// GetSigners defines whose signature is required
func (m *MsgDepositClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Orchestrator)}
}

// Type should return the action
func (m *MsgDepositClaim) Type() string { return TypeMsgDepositClaim }

// Route should return the name of the module
func (m *MsgDepositClaim) Route() string { return RouterKey }

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
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BatchNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero batch nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	if err := fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	return nil
}

// ClaimHash Hash implements WithdrawBatch.Hash
func (m *MsgWithdrawClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%d/", m.TokenContract, m.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (m *MsgWithdrawClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.Orchestrator)
}

// GetSigners defines whose signature is required
func (m *MsgWithdrawClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Orchestrator)}
}

// Route should return the name of the module
func (m *MsgWithdrawClaim) Route() string { return RouterKey }

// Type should return the action
func (m *MsgWithdrawClaim) Type() string { return TypeMsgWithdrawClaim }

func (m *MsgFxOriginatedTokenClaim) Route() string {
	return RouterKey
}

func (m *MsgFxOriginatedTokenClaim) Type() string { return TypeMsgFxOriginatedTokenClaim }

func (m *MsgFxOriginatedTokenClaim) ValidateBasic() error {
	if err := fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if len(m.Name) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty token name")
	}
	if len(m.Symbol) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty token symbol")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m *MsgFxOriginatedTokenClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgFxOriginatedTokenClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Orchestrator)}
}

func (m *MsgFxOriginatedTokenClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.Orchestrator)
}

func (m *MsgFxOriginatedTokenClaim) GetType() ClaimType {
	return CLAIM_TYPE_ORIGINATED_TOKEN
}

func (m *MsgFxOriginatedTokenClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%s/%s/%s/%d/", m.TokenContract, m.Name, m.Symbol, m.Decimals)
	return tmhash.Sum([]byte(path))
}

// GetType returns the type of the claim
func (m *MsgValsetUpdatedClaim) GetType() ClaimType {
	return CLAIM_TYPE_VALSET_UPDATED
}

// ValidateBasic performs stateless checks
func (m *MsgValsetUpdatedClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid orchestrator address: %s", err)
	}
	if len(m.Members) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty members")
	}
	for _, member := range m.Members {
		if err := fxtypes.ValidateEthereumAddress(member.EthAddress); err != nil {
			return errortypes.ErrInvalidAddress.Wrapf("invalid eth address: %s", err)
		}
		if member.Power == 0 {
			return errortypes.ErrInvalidRequest.Wrap("zero power")
		}
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgValsetUpdatedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgValsetUpdatedClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.Orchestrator)
}

// GetSigners defines whose signature is required
func (m *MsgValsetUpdatedClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Orchestrator)}
}

// Type should return the action
func (m *MsgValsetUpdatedClaim) Type() string { return TypeMsgValsetUpdatedClaim }

// Route should return the name of the module
func (m *MsgValsetUpdatedClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeDeposit.Hash
func (m *MsgValsetUpdatedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%s/", m.ValsetNonce, m.EventNonce, m.BlockHeight, m.Members)
	return tmhash.Sum([]byte(path))
}
