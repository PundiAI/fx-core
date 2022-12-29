package types

import (
	"fmt"
	"regexp"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	TypeMsgSendToExternal       = "send_to_external"
	TypeMsgCancelSendToExternal = "cancel_send_to_external"
	TypeMsgSendToExternalClaim  = "send_to_external_claim"

	TypeMsgRequestBatch = "request_batch"
	TypeMsgConfirmBatch = "confirm_batch"
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
	_ sdk.Msg       = &MsgSendToExternalClaim{}
	_ CrossChainMsg = &MsgSendToExternalClaim{}

	_ sdk.Msg       = &MsgRequestBatch{}
	_ CrossChainMsg = &MsgRequestBatch{}
	_ sdk.Msg       = &MsgConfirmBatch{}
	_ CrossChainMsg = &MsgConfirmBatch{}
)

type MsgValidateBasic interface {
	MsgBondedOracleValidate(m *MsgBondedOracle) (err error)
	MsgAddDelegateValidate(m *MsgAddDelegate) (err error)
	MsgReDelegateValidate(m *MsgReDelegate) (err error)
	MsgEditBridgerValidate(m *MsgEditBridger) (err error)
	MsgWithdrawRewardValidate(m *MsgWithdrawReward) (err error)
	MsgUnbondedOracleValidate(m *MsgUnbondedOracle) (err error)

	MsgOracleSetConfirmValidate(m *MsgOracleSetConfirm) (err error)
	MsgOracleSetUpdatedClaimValidate(m *MsgOracleSetUpdatedClaim) (err error)
	MsgBridgeTokenClaimValidate(m *MsgBridgeTokenClaim) (err error)
	MsgSendToExternalClaimValidate(m *MsgSendToExternalClaim) (err error)

	MsgSendToFxClaimValidate(m *MsgSendToFxClaim) (err error)
	MsgSendToExternalValidate(m *MsgSendToExternal) (err error)

	MsgCancelSendToExternalValidate(m *MsgCancelSendToExternal) (err error)
	MsgRequestBatchValidate(m *MsgRequestBatch) (err error)
	MsgConfirmBatchValidate(m *MsgConfirmBatch) (err error)
}

var reModuleName *regexp.Regexp

func init() {
	reModuleNameString := `[a-zA-Z][a-zA-Z0-9/]{1,32}`
	reModuleName = regexp.MustCompile(fmt.Sprintf(`^%s$`, reModuleNameString))
}

// ValidateModuleName is the default validation function for crosschain moduleName.
func ValidateModuleName(moduleName string) error {
	if !reModuleName.MatchString(moduleName) {
		return fmt.Errorf("invalid module name: %s", moduleName)
	}
	return nil
}

var msgValidateBasicRouter = make(map[string]MsgValidateBasic)

func RegisterValidateBasic(chainName string, validate MsgValidateBasic) {
	if err := ValidateModuleName(chainName); err != nil {
		panic(sdkerrors.Wrap(ErrInvalid, "chain name"))
	}
	if _, ok := msgValidateBasicRouter[chainName]; ok {
		panic(fmt.Sprintf("duplicate registry msg validateBasic! chainName: %s", chainName))
	}
	msgValidateBasicRouter[chainName] = validate
}

// MsgBondedOracle

func (m *MsgBondedOracle) Route() string { return RouterKey }

func (m *MsgBondedOracle) Type() string { return TypeMsgBondedOracle }

func (m *MsgBondedOracle) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgBondedOracleValidate(m)
	}
}

func (m *MsgBondedOracle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBondedOracle) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgAddDelegate

func (m *MsgAddDelegate) Route() string { return RouterKey }

func (m *MsgAddDelegate) Type() string {
	return TypeMsgAddDelegate
}

func (m *MsgAddDelegate) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgAddDelegateValidate(m)
	}
}

func (m *MsgAddDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgAddDelegate) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgReDelegate

func (m *MsgReDelegate) Route() string { return RouterKey }

func (m *MsgReDelegate) Type() string {
	return TypeMsgReDelegate
}

func (m *MsgReDelegate) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgReDelegateValidate(m)
	}
}

func (m *MsgReDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgReDelegate) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgEditBridger

func (m *MsgEditBridger) Route() string { return RouterKey }

func (m *MsgEditBridger) Type() string { return TypeMsgEditBridger }

func (m *MsgEditBridger) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgEditBridgerValidate(m)
	}
}

func (m *MsgEditBridger) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgEditBridger) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgWithdrawReward

func (m *MsgWithdrawReward) Route() string { return RouterKey }

func (m *MsgWithdrawReward) Type() string { return TypeMsgWithdrawReward }

func (m *MsgWithdrawReward) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgWithdrawRewardValidate(m)
	}
}

func (m *MsgWithdrawReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgWithdrawReward) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgUnbondedOracle

func (m *MsgUnbondedOracle) Route() string { return RouterKey }

func (m *MsgUnbondedOracle) Type() string { return TypeMsgUnbondedOracle }

func (m *MsgUnbondedOracle) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgUnbondedOracleValidate(m)
	}
}

func (m *MsgUnbondedOracle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUnbondedOracle) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgOracleSetConfirm

// Route should return the name of the module
func (m *MsgOracleSetConfirm) Route() string { return RouterKey }

// Type should return the action
func (m *MsgOracleSetConfirm) Type() string { return TypeMsgOracleSetConfirm }

// ValidateBasic performs stateless checks
func (m *MsgOracleSetConfirm) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgOracleSetConfirmValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgOracleSetConfirm) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgOracleSetConfirm) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgSendToExternal

// Route should return the name of the module
func (m *MsgSendToExternal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendToExternal) Type() string { return TypeMsgSendToExternal }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (m *MsgSendToExternal) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgSendToExternalValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSendToExternal) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgRequestBatch

// Route should return the name of the module
func (m *MsgRequestBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRequestBatch) Type() string { return TypeMsgRequestBatch }

// ValidateBasic performs stateless checks
func (m *MsgRequestBatch) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgRequestBatchValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgRequestBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgRequestBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgConfirmBatch

// Route should return the name of the module
func (m *MsgConfirmBatch) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConfirmBatch) Type() string { return TypeMsgConfirmBatch }

// ValidateBasic performs stateless checks
func (m *MsgConfirmBatch) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgConfirmBatchValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgConfirmBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConfirmBatch) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// MsgCancelSendToExternal

// Route should return the name of the module
func (m *MsgCancelSendToExternal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgCancelSendToExternal) Type() string { return TypeMsgCancelSendToExternal }

// ValidateBasic performs stateless checks
func (m *MsgCancelSendToExternal) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgCancelSendToExternalValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgCancelSendToExternal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgCancelSendToExternal) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
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
	_ ExternalClaim = &MsgBridgeTokenClaim{}
	_ ExternalClaim = &MsgSendToExternalClaim{}
	_ ExternalClaim = &MsgOracleSetUpdatedClaim{}
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
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgSendToFxClaimValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgSendToFxClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgSendToFxClaim) GetClaimer() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return val
}

// GetSigners defines whose signature is required
func (m *MsgSendToFxClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
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

// MsgSendToExternalClaim

// GetType returns the claim type
func (m *MsgSendToExternalClaim) GetType() ClaimType {
	return CLAIM_TYPE_SEND_TO_EXTERNAL
}

// ValidateBasic performs stateless checks
func (m *MsgSendToExternalClaim) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgSendToExternalClaimValidate(m)
	}
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
	val, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(fmt.Sprintf("invalid address %s", m.BridgerAddress))
	}
	return val
}

// GetSigners defines whose signature is required
func (m *MsgSendToExternalClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// Route should return the name of the module
func (m *MsgSendToExternalClaim) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendToExternalClaim) Type() string { return TypeMsgSendToExternalClaim }

// MsgBridgeTokenClaim

func (m *MsgBridgeTokenClaim) Route() string { return RouterKey }

func (m *MsgBridgeTokenClaim) Type() string { return TypeMsgBridgeTokenClaim }

func (m *MsgBridgeTokenClaim) ValidateBasic() (err error) {
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgBridgeTokenClaimValidate(m)
	}
}

func (m *MsgBridgeTokenClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgBridgeTokenClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgBridgeTokenClaim) GetClaimer() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(fmt.Sprintf("invalid address %s", m.BridgerAddress))
	}
	return val
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
	if err = ValidateModuleName(m.ChainName); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "chain name")
	}
	if router, ok := msgValidateBasicRouter[m.ChainName]; !ok {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", m.ChainName))
	} else {
		return router.MsgOracleSetUpdatedClaimValidate(m)
	}
}

// GetSignBytes encodes the message for signing
func (m *MsgOracleSetUpdatedClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgOracleSetUpdatedClaim) GetClaimer() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(fmt.Sprintf("invalid address %s", m.BridgerAddress))
	}
	return val
}

// GetSigners defines whose signature is required
func (m *MsgOracleSetUpdatedClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
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
	acc, err := sdk.AccAddressFromBech32(m.BridgerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgAddOracleDeposit) ValidateBasic() (err error) {
	return nil
}

func (m *MsgAddOracleDeposit) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OracleAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}
