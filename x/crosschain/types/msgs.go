package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// cross chain message types
const (
	TypeMsgBondedOracle   = "bonded_oracle"
	TypeMsgAddDelegate    = "add_delegate"
	TypeMsgReDelegate     = "re_delegate"
	TypeMsgEditBridger    = "edit_bridger"
	TypeMsgWithdrawReward = "withdraw_reward"
	TypeMsgUnbondedOracle = "unbonded_oracle"

	TypeMsgOracleSetConfirm      = "valset_confirm"
	TypeMsgOracleSetUpdatedClaim = "valset_updated_claim"

	TypeMsgBridgeTokenClaim = "bridge_token_claim"

	TypeMsgSendToFxClaim = "send_to_fx_claim"

	TypeMsgSendToExternal        = "send_to_external"
	TypeMsgCancelSendToExternal  = "cancel_send_to_external"
	TypeMsgIncreaseBridgeFee     = "increase_bridge_fee"
	TypeMsgAddPendingPoolRewards = "add_pending_pool_rewards"
	TypeMsgSendToExternalClaim   = "send_to_external_claim"

	TypeMsgBridgeCallClaim = "bridge_call_claim"

	TypeMsgBridgeCall            = "bridge_call"
	TypeMsgBridgeCallConfirm     = "bridge_call_confirm"
	TypeMsgBridgeCallResultClaim = "bridge_call_result_claim"

	TypeMsgRequestBatch = "request_batch"
	TypeMsgConfirmBatch = "confirm_batch"

	TypeMsgUpdateParams = "update_params"

	TypeMsgUpdateChainOracles = "update_chain_oracles"
)

type (
	// CrossChainMsg cross msg must implement GetChainName interface.. using in router
	CrossChainMsg interface {
		GetChainName() string
	}
)

var (
	_ sdk.Msg       = &MsgBondedOracle{}
	_ CrossChainMsg = &MsgBondedOracle{}
	_ sdk.Msg       = &MsgAddDelegate{}
	_ CrossChainMsg = &MsgAddDelegate{}
	_ sdk.Msg       = &MsgReDelegate{}
	_ CrossChainMsg = &MsgReDelegate{}
	_ sdk.Msg       = &MsgEditBridger{}
	_ CrossChainMsg = &MsgEditBridger{}
	_ sdk.Msg       = &MsgWithdrawReward{}
	_ CrossChainMsg = &MsgWithdrawReward{}
	_ sdk.Msg       = &MsgUnbondedOracle{}
	_ CrossChainMsg = &MsgUnbondedOracle{}

	_ sdk.Msg       = &MsgOracleSetConfirm{}
	_ CrossChainMsg = &MsgOracleSetConfirm{}
	_ sdk.Msg       = &MsgOracleSetUpdatedClaim{}
	_ CrossChainMsg = &MsgOracleSetUpdatedClaim{}

	_ sdk.Msg       = &MsgBridgeTokenClaim{}
	_ CrossChainMsg = &MsgBridgeTokenClaim{}

	_ sdk.Msg       = &MsgSendToFxClaim{}
	_ CrossChainMsg = &MsgSendToFxClaim{}

	_ sdk.Msg       = &MsgSendToExternal{}
	_ CrossChainMsg = &MsgSendToExternal{}
	_ sdk.Msg       = &MsgCancelSendToExternal{}
	_ CrossChainMsg = &MsgCancelSendToExternal{}
	_ sdk.Msg       = &MsgIncreaseBridgeFee{}
	_ CrossChainMsg = &MsgIncreaseBridgeFee{}
	_ sdk.Msg       = &MsgSendToExternalClaim{}
	_ CrossChainMsg = &MsgSendToExternalClaim{}
	_ sdk.Msg       = &MsgAddPendingPoolRewards{}
	_ CrossChainMsg = &MsgAddPendingPoolRewards{}

	_ sdk.Msg       = &MsgRequestBatch{}
	_ CrossChainMsg = &MsgRequestBatch{}
	_ sdk.Msg       = &MsgConfirmBatch{}
	_ CrossChainMsg = &MsgConfirmBatch{}

	_ sdk.Msg       = &MsgBridgeCallClaim{}
	_ CrossChainMsg = &MsgBridgeCallClaim{}

	_ sdk.Msg       = &MsgBridgeCall{}
	_ CrossChainMsg = &MsgBridgeCall{}
	_ sdk.Msg       = &MsgBridgeCallConfirm{}
	_ CrossChainMsg = &MsgBridgeCallConfirm{}
	_ sdk.Msg       = &MsgBridgeCallResultClaim{}
	_ CrossChainMsg = &MsgBridgeCallResultClaim{}

	_ sdk.Msg       = &MsgUpdateParams{}
	_ CrossChainMsg = &MsgUpdateParams{}

	_ sdk.Msg       = &MsgUpdateChainOracles{}
	_ CrossChainMsg = &MsgUpdateChainOracles{}
)

// MsgBondedOracle

func (m *MsgBondedOracle) Route() string { return RouterKey }

func (m *MsgBondedOracle) Type() string { return TypeMsgBondedOracle }

func (m *MsgBondedOracle) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if !m.DelegateAmount.IsValid() || m.DelegateAmount.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid delegation amount")
	}
	if m.OracleAddress == m.BridgerAddress {
		return errortypes.ErrInvalidRequest.Wrap("same address")
	}
	return nil
}

func (m *MsgBondedOracle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBondedOracle) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// MsgAddDelegate

func (m *MsgAddDelegate) Route() string { return RouterKey }

func (m *MsgAddDelegate) Type() string {
	return TypeMsgAddDelegate
}

func (m *MsgAddDelegate) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

func (m *MsgAddDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// MsgReDelegate

func (m *MsgReDelegate) Route() string { return RouterKey }

func (m *MsgReDelegate) Type() string {
	return TypeMsgReDelegate
}

func (m *MsgReDelegate) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if _, err = sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	return nil
}

func (m *MsgReDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgReDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// MsgEditBridger

func (m *MsgEditBridger) Route() string { return RouterKey }

func (m *MsgEditBridger) Type() string { return TypeMsgEditBridger }

func (m *MsgEditBridger) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	if _, err = sdk.ValAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if m.OracleAddress == m.BridgerAddress {
		return errortypes.ErrInvalidRequest.Wrap("same address")
	}
	return nil
}

func (m *MsgEditBridger) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgEditBridger) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// MsgWithdrawReward

func (m *MsgWithdrawReward) Route() string { return RouterKey }

func (m *MsgWithdrawReward) Type() string { return TypeMsgWithdrawReward }

func (m *MsgWithdrawReward) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	return nil
}

func (m *MsgWithdrawReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// MsgUnbondedOracle

func (m *MsgUnbondedOracle) Route() string { return RouterKey }

func (m *MsgUnbondedOracle) Type() string { return TypeMsgUnbondedOracle }

func (m *MsgUnbondedOracle) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid oracle address: %s", err)
	}
	return nil
}

func (m *MsgUnbondedOracle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUnbondedOracle) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// MsgOracleSetConfirm

// Route should return the name of the module
func (m *MsgOracleSetConfirm) Route() string { return RouterKey }

// Type should return the action
func (m *MsgOracleSetConfirm) Type() string { return TypeMsgOracleSetConfirm }

// ValidateBasic performs stateless checks
func (m *MsgOracleSetConfirm) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
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
func (m *MsgOracleSetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgOracleSetConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// MsgSendToExternal

// Route should return the name of the module
func (m *MsgSendToExternal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendToExternal) Type() string { return TypeMsgSendToExternal }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m *MsgSendToExternal) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Dest); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid dest address: %s", err)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return errortypes.ErrInvalidRequest.Wrap("bridge fee denom not equal amount denom")
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSendToExternal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// MsgRequestBatch

// Route should return the name of the module
func (m *MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRequestBatch) Type() string { return TypeMsgRequestBatch }

// ValidateBasic performs stateless checks
func (m *MsgRequestBatch) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.Denom) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty denom")
	}
	if m.MinimumFee.IsNil() || !m.MinimumFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid minimum fee")
	}
	if err = ValidateExternalAddr(m.ChainName, m.FeeReceive); err != nil {
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

// MsgConfirmBatch

// Route should return the name of the module
func (m *MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConfirmBatch) Type() string { return TypeMsgConfirmBatch }

// ValidateBasic performs stateless checks
func (m *MsgConfirmBatch) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
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
func (m *MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// MsgBridgeCallConfirm

func (m *MsgBridgeCallConfirm) Route() string { return RouterKey }

func (m *MsgBridgeCallConfirm) Type() string { return TypeMsgBridgeCallConfirm }

func (m *MsgBridgeCallConfirm) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.ExternalAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
	}
	if len(m.Signature) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	return nil
}

func (m *MsgBridgeCallConfirm) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

func (m *MsgBridgeCallConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// MsgCancelSendToExternal

// Route should return the name of the module
func (m *MsgCancelSendToExternal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgCancelSendToExternal) Type() string { return TypeMsgCancelSendToExternal }

// ValidateBasic performs stateless checks
func (m *MsgCancelSendToExternal) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero transaction id")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgCancelSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgCancelSendToExternal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// MsgIncreaseBridgeFee

// Route should return the name of the module
func (m *MsgIncreaseBridgeFee) Route() string { return RouterKey }

// Type should return the action
func (m *MsgIncreaseBridgeFee) Type() string { return TypeMsgIncreaseBridgeFee }

// ValidateBasic performs stateless checks
func (m *MsgIncreaseBridgeFee) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero transaction id")
	}
	if !m.AddBridgeFee.IsValid() || !m.AddBridgeFee.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid bridge fee")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgIncreaseBridgeFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgIncreaseBridgeFee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// MsgAddPendingPoolRewards

// Route should return the name of the module
func (m *MsgAddPendingPoolRewards) Route() string { return RouterKey }

// Type should return the action
func (m *MsgAddPendingPoolRewards) Type() string { return TypeMsgAddPendingPoolRewards }

// ValidateBasic performs stateless checks
func (m *MsgAddPendingPoolRewards) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if m.TransactionId == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero transaction id")
	}
	rewards := sdk.NewCoins(m.Rewards...)
	if rewards.Empty() || !rewards.IsValid() || !rewards.IsAllPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid rewards")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgAddPendingPoolRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgAddPendingPoolRewards) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// ExternalClaim represents a claim on ethereum state
type ExternalClaim interface {
	proto.Message
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
	// GetClaimer the delegate address of the claimer, for MsgSendToExternalClaim and MsgSendToFxClaim
	// this is sent in as the sdk.AccAddress of the delegated key. it is up to the user
	// to disambiguate this into a sdk.ValAddress
	GetClaimer() sdk.AccAddress
	// GetType Which type of claim this is
	GetType() ClaimType
	ValidateBasic() error
	ClaimHash() []byte
}

var (
	_ ExternalClaim = &MsgSendToFxClaim{}
	_ ExternalClaim = &MsgBridgeCallClaim{}
	_ ExternalClaim = &MsgBridgeTokenClaim{}
	_ ExternalClaim = &MsgSendToExternalClaim{}
	_ ExternalClaim = &MsgOracleSetUpdatedClaim{}
	_ ExternalClaim = &MsgBridgeCallResultClaim{}
)

func UnpackAttestationClaim(cdc codectypes.AnyUnpacker, att *Attestation) (ExternalClaim, error) {
	var msg ExternalClaim
	err := cdc.UnpackAny(att.Claim, &msg)
	return msg, err
}

// MsgSendToFxClaim

// GetType returns the type of the claim
func (m *MsgSendToFxClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_FX
}

// ValidateBasic performs stateless checks
func (m *MsgSendToFxClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
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
func (m *MsgSendToFxClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSendToFxClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

// GetSigners defines whose signature is required
func (m *MsgSendToFxClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// Type should return the action
func (m *MsgSendToFxClaim) Type() string { return TypeMsgSendToFxClaim }

// Route should return the name of the module
func (m *MsgSendToFxClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m *MsgSendToFxClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d%s/%s/%s/%s/%s", m.BlockHeight, m.EventNonce, m.TokenContract, m.Sender, m.Amount.String(), m.Receiver, m.TargetIbc)
	return tmhash.Sum([]byte(path))
}

// MsgBridgeCallClaim

// GetType returns the type of the claim
func (m *MsgBridgeCallClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_CALL
}

// ValidateBasic performs stateless checks
func (m *MsgBridgeCallClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	return msgBridgeCallClaimValidateBasic(m)
}

func msgBridgeCallClaimValidateBasic(m *MsgBridgeCallClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.To) > 0 {
		if err = ValidateExternalAddr(m.ChainName, m.To); err != nil {
			return errortypes.ErrInvalidAddress.Wrapf("invalid to contract: %s", err)
		}
	}
	if err = ValidateExternalAddr(m.ChainName, m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if m.Value.IsNil() || m.Value.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid value")
	}
	if len(m.Data) > 0 {
		if _, err = hex.DecodeString(m.Data); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid data")
		}
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	if err = ValidateExternalAddr(m.ChainName, m.TxOrigin); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid tx origin: %s", err)
	}
	if len(m.Memo) > 0 {
		if _, err = hex.DecodeString(m.Memo); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid memo")
		}
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgBridgeCallClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBridgeCallClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

// GetSigners defines whose signature is required
func (m *MsgBridgeCallClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// Type should return the action
func (m *MsgBridgeCallClaim) Type() string { return TypeMsgBridgeCallClaim }

// Route should return the name of the module
func (m *MsgBridgeCallClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m *MsgBridgeCallClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%s/%v/%v/%s", m.BlockHeight, m.EventNonce, m.Sender, m.Receiver, m.To, m.TokenContracts, m.Amounts, m.Data, m.Value.String())
	return tmhash.Sum([]byte(path))
}

func (m *MsgBridgeCallClaim) GetSenderAddr() common.Address {
	return common.BytesToAddress(ExternalAddrToAccAddr(m.ChainName, m.Sender).Bytes())
}

func (m *MsgBridgeCallClaim) GetReceiverAddr() sdk.AccAddress {
	return ExternalAddrToAccAddr(m.ChainName, m.Receiver)
}

func (m *MsgBridgeCallClaim) GetToAddr() *common.Address {
	if len(m.To) == 0 {
		return nil
	}
	addr := common.BytesToAddress(ExternalAddrToAccAddr(m.ChainName, m.To).Bytes())
	return &addr
}

func (m *MsgBridgeCallClaim) MustData() []byte {
	if len(m.Data) == 0 {
		return []byte{}
	}
	bz, err := hex.DecodeString(m.Data)
	if err != nil {
		panic(err)
	}
	return bz
}

func (m *MsgBridgeCallClaim) GetTokensAddr() []common.Address {
	addrs := make([]common.Address, 0, len(m.TokenContracts))
	for _, token := range m.TokenContracts {
		addr := ExternalAddrToAccAddr(m.ChainName, token)
		addrs = append(addrs, common.BytesToAddress(addr))
	}
	return addrs
}

func (m *MsgBridgeCallClaim) GetAmounts() []*big.Int {
	amts := make([]*big.Int, 0, len(m.Amounts))
	for _, a := range m.Amounts {
		amts = append(amts, a.BigInt())
	}
	return amts
}

// MsgBridgeCallResultClaim

// GetType returns the type of the claim
func (m *MsgBridgeCallResultClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_CALL_RESULT
}

// ValidateBasic performs stateless checks
func (m *MsgBridgeCallResultClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.To) > 0 {
		if err = ValidateExternalAddr(m.ChainName, m.To); err != nil {
			return errortypes.ErrInvalidAddress.Wrapf("invalid to contract: %s", err)
		}
	}
	if err = ValidateExternalAddr(m.ChainName, m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if m.Nonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero nonce")
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	if err = ValidateExternalAddr(m.ChainName, m.TxOrigin); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid tx origin: %s", err)
	}
	if len(m.Cause) > 0 {
		if _, err = hex.DecodeString(m.Cause); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid cause")
		}
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgBridgeCallResultClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBridgeCallResultClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

// GetSigners defines whose signature is required
func (m *MsgBridgeCallResultClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// Type should return the action
func (m *MsgBridgeCallResultClaim) Type() string { return TypeMsgBridgeCallResultClaim }

// Route should return the name of the module
func (m *MsgBridgeCallResultClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m *MsgBridgeCallResultClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%s/%s/%d/%t/%s", m.BlockHeight, m.EventNonce, m.Sender, m.Receiver, m.To, m.Nonce, m.Success, m.Cause)
	return tmhash.Sum([]byte(path))
}

// MsgSendToExternalClaim

// GetType returns the claim type
func (m *MsgSendToExternalClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_EXTERNAL
}

// ValidateBasic performs stateless checks
func (m *MsgSendToExternalClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	if m.BatchNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero batch nonce")
	}
	return nil
}

// ClaimHash Hash implements SendToFxBatch.Hash
func (m *MsgSendToExternalClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%s/%d/", m.BlockHeight, m.EventNonce, m.TokenContract, m.BatchNonce)
	return tmhash.Sum([]byte(path))
}

// GetSignBytes encodes the message for signing
func (m *MsgSendToExternalClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSendToExternalClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

// GetSigners defines whose signature is required
func (m *MsgSendToExternalClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// Route should return the name of the module
func (m *MsgSendToExternalClaim) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendToExternalClaim) Type() string { return TypeMsgSendToExternalClaim }

// MsgBridgeTokenClaim

func (m *MsgBridgeTokenClaim) Route() string { return RouterKey }

func (m *MsgBridgeTokenClaim) Type() string { return TypeMsgBridgeTokenClaim }

func (m *MsgBridgeTokenClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.TokenContract); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid token contract: %s", err)
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return errortypes.ErrInvalidRequest.Wrap("could not decode hex channelIbc string")
	}
	if len(m.Name) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty token name")
	}
	if len(m.Symbol) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty token symbol")
	}
	if m.EventNonce == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero event nonce")
	}
	if m.BlockHeight == 0 {
		return errortypes.ErrInvalidRequest.Wrap("zero block height")
	}
	return nil
}

func (m *MsgBridgeTokenClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBridgeTokenClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

func (m *MsgBridgeTokenClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

func (m *MsgBridgeTokenClaim) GetType() ClaimType {
	return CLAIM_TYPE_BRIDGE_TOKEN
}

func (m *MsgBridgeTokenClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d%s/%s/%s/%d/%s/", m.BlockHeight, m.EventNonce, m.TokenContract, m.Name, m.Symbol, m.Decimals, m.ChannelIbc)
	return tmhash.Sum([]byte(path))
}

// MsgOracleSetUpdatedClaim

// GetType returns the type of the claim
func (m *MsgOracleSetUpdatedClaim) GetType() ClaimType {
	return CLAIM_TYPE_ORACLE_SET_UPDATED
}

// ValidateBasic performs stateless checks
func (m *MsgOracleSetUpdatedClaim) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid bridger address: %s", err)
	}
	if len(m.Members) == 0 {
		return errortypes.ErrInvalidRequest.Wrap("empty members")
	}
	for _, member := range m.Members {
		if err = ValidateExternalAddr(m.ChainName, member.ExternalAddress); err != nil {
			return errortypes.ErrInvalidAddress.Wrapf("invalid external address: %s", err)
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
func (m *MsgOracleSetUpdatedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgOracleSetUpdatedClaim) GetClaimer() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.BridgerAddress)
}

// GetSigners defines whose signature is required
func (m *MsgOracleSetUpdatedClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

// Type should return the action
func (m *MsgOracleSetUpdatedClaim) Type() string { return TypeMsgOracleSetUpdatedClaim }

// Route should return the name of the module
func (m *MsgOracleSetUpdatedClaim) Route() string { return RouterKey }

// ClaimHash Hash implements BridgeSendToExternal.Hash
func (m *MsgOracleSetUpdatedClaim) ClaimHash() []byte {
	path := fmt.Sprintf("%d/%d/%d/%v/", m.BlockHeight, m.OracleSetNonce, m.EventNonce, m.Members)
	return tmhash.Sum([]byte(path))
}

func (m *MsgSetOrchestratorAddress) ValidateBasic() (err error) {
	return nil
}

func (m *MsgSetOrchestratorAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.BridgerAddress)}
}

func (m *MsgAddOracleDeposit) ValidateBasic() (err error) {
	return nil
}

func (m *MsgAddOracleDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.OracleAddress)}
}

// Route returns the MsgUpdateParams message route.
func (m *MsgUpdateParams) Route() string { return ModuleName }

// Type returns the MsgUpdateParams message type.
func (m *MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	if len(m.Params.Oracles) > 0 {
		return errors.New("deprecated oracles")
	}
	return nil
}

// Route returns the MsgUpdateChainOracles message route.
func (m *MsgUpdateChainOracles) Route() string { return ModuleName }

// Type returns the MsgUpdateChainOracles message type.
func (m *MsgUpdateChainOracles) Type() string { return TypeMsgUpdateChainOracles }

// GetSignBytes returns the raw bytes for a MsgUpdateChainOracles message that
// the expected signer needs to sign.
func (m *MsgUpdateChainOracles) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateChainOracles) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if len(m.Oracles) == 0 {
		return errors.New("empty oracles")
	}
	oraclesMap := make(map[string]bool)
	for _, addr := range m.Oracles {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrap(err, "oracle address")
		}
		if oraclesMap[addr] {
			return errors.New("duplicate oracle address")
		}
		oraclesMap[addr] = true
	}
	return nil
}

func (m *MsgUpdateChainOracles) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m *MsgBridgeCall) Route() string { return RouterKey }
func (m *MsgBridgeCall) Type() string  { return TypeMsgBridgeCall }
func (m *MsgBridgeCall) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgBridgeCall) ValidateBasic() (err error) {
	if _, ok := externalAddressRouter[m.ChainName]; !ok {
		return errortypes.ErrInvalidRequest.Wrap("unrecognized cross chain name")
	}
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if err = ValidateExternalAddr(m.ChainName, m.To); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}
	if m.Value.IsNil() || m.Value.IsNegative() {
		return errortypes.ErrInvalidRequest.Wrap("invalid value")
	}
	if err = m.Coins.Validate(); err != nil {
		return errortypes.ErrInvalidCoins.Wrap(err.Error())
	}
	if len(m.Data) > 0 {
		if _, err = hex.DecodeString(m.Data); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid data")
		}
	}
	if len(m.Memo) > 0 {
		if _, err = hex.DecodeString(m.Memo); err != nil {
			return errortypes.ErrInvalidRequest.Wrap("invalid memo")
		}
	}
	return nil
}

func (m *MsgBridgeCall) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}
