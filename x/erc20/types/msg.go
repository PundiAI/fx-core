package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

var (
	_ sdk.Msg = &MsgConvertCoin{}
	_ sdk.Msg = &MsgConvertERC20{}

	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgRegisterCoin{}
	_ sdk.Msg = &MsgRegisterERC20{}
	_ sdk.Msg = &MsgToggleTokenConversion{}
	_ sdk.Msg = &MsgUpdateDenomAlias{}
)

const (
	TypeMsgConvertCoin  = "convert_coin"
	TypeMsgConvertERC20 = "convert_ERC20"
	TypeMsgConvertDenom = "convert_denom"
	TypeMsgUpdateParams = "update_params"

	TypeMsgRegisterCoin          = "register_coin"
	TypeMsgRegisterERC20         = "register_erc20"
	TypeMsgToggleTokenConversion = "toggle_token_conversion" // #nosec G101
	TypeMsgUpdateDenomAlias      = "update_denom_alias"
)

// NewMsgConvertCoin creates a new instance of MsgConvertCoin
func NewMsgConvertCoin(coin sdk.Coin, receiver common.Address, sender sdk.AccAddress) *MsgConvertCoin {
	return &MsgConvertCoin{
		Coin:     coin,
		Receiver: receiver.Hex(),
		Sender:   sender.String(),
	}
}

// Route should return the name of the module
func (m *MsgConvertCoin) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConvertCoin) Type() string { return TypeMsgConvertCoin }

// ValidateBasic runs stateless checks on the message
func (m *MsgConvertCoin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if err = fxtypes.ValidateEthereumAddress(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err.Error())
	}
	if err = ibctransfertypes.ValidateIBCDenom(m.Coin.Denom); err != nil {
		return errortypes.ErrInvalidCoins.Wrapf("invalid coin denom %s", err.Error())
	}
	if m.Coin.Amount.IsNil() || !m.Coin.Amount.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgConvertCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConvertCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// NewMsgConvertERC20 creates a new instance of MsgConvertERC20
func NewMsgConvertERC20(amount sdkmath.Int, receiver sdk.AccAddress, contract, sender common.Address) *MsgConvertERC20 {
	return &MsgConvertERC20{
		ContractAddress: contract.String(),
		Amount:          amount,
		Receiver:        receiver.String(),
		Sender:          sender.Hex(),
	}
}

// Route should return the name of the module
func (m *MsgConvertERC20) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConvertERC20) Type() string { return TypeMsgConvertERC20 }

// ValidateBasic runs stateless checks on the message
func (m *MsgConvertERC20) ValidateBasic() error {
	if err := fxtypes.ValidateEthereumAddress(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err.Error())
	}
	if err := fxtypes.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid contract address: %s", err.Error())
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgConvertERC20) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConvertERC20) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{common.HexToAddress(m.Sender).Bytes()}
}

func NewMsgConvertDenom(sender, receiver sdk.AccAddress, coin sdk.Coin, target string) *MsgConvertDenom {
	return &MsgConvertDenom{
		Sender:   sender.String(),
		Receiver: receiver.String(),
		Coin:     coin,
		Target:   target,
	}
}

// Route should return the name of the module
func (m *MsgConvertDenom) Route() string { return RouterKey }

// Type should return the action
func (m *MsgConvertDenom) Type() string { return TypeMsgConvertDenom }

// ValidateBasic runs stateless checks on the message
func (m *MsgConvertDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err.Error())
	}
	if !m.Coin.IsValid() || !m.Coin.IsPositive() {
		return errortypes.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgConvertDenom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgConvertDenom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
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
	if err := m.Params.Validate(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	return nil
}

// Route returns the MsgRegisterCoin message route.
func (m *MsgRegisterCoin) Route() string { return ModuleName }

// Type returns the MsgRegisterCoin message type.
func (m *MsgRegisterCoin) Type() string { return TypeMsgRegisterCoin }

// GetSignBytes returns the raw bytes for a MsgRegisterCoin message that
// the expected signer needs to sign.
func (m *MsgRegisterCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRegisterCoin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := m.Metadata.Validate(); err != nil {
		return errorsmod.Wrap(err, "metadata")
	}
	if err := fxtypes.ValidateMetadata(m.Metadata); err != nil {
		return errorsmod.Wrap(err, "metadata")
	}
	if err := ibctransfertypes.ValidateIBCDenom(m.Metadata.Base); err != nil {
		return errorsmod.Wrap(err, "metadata base")
	}
	return nil
}

func (m *MsgRegisterCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// Route returns the MsgRegisterERC20 message route.
func (m *MsgRegisterERC20) Route() string { return ModuleName }

// Type returns the MsgRegisterERC20 message type.
func (m *MsgRegisterERC20) Type() string { return TypeMsgRegisterERC20 }

// GetSignBytes returns the raw bytes for a MsgRegisterERC20 message that
// the expected signer needs to sign.
func (m *MsgRegisterERC20) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRegisterERC20) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := fxtypes.ValidateEthereumAddress(m.Erc20Address); err != nil {
		return errorsmod.Wrap(err, "ERC20 address")
	}
	seenAliases := make(map[string]bool)
	for _, alias := range m.Aliases {
		if seenAliases[alias] {
			return fmt.Errorf("duplicate denomination unit alias %s", alias)
		}
		if strings.TrimSpace(alias) == "" {
			return fmt.Errorf("alias for denom unit %s cannot be blank", alias)
		}
		if err := sdk.ValidateDenom(alias); err != nil {
			return errorsmod.Wrap(err, "alias")
		}
		seenAliases[alias] = true
	}
	return nil
}

func (m *MsgRegisterERC20) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// Route returns the MsgToggleTokenConversion message route.
func (m *MsgToggleTokenConversion) Route() string { return ModuleName }

// Type returns the MsgToggleTokenConversion message type.
func (m *MsgToggleTokenConversion) Type() string { return TypeMsgToggleTokenConversion }

// GetSignBytes returns the raw bytes for a MsgToggleTokenConversion message that
// the expected signer needs to sign.
func (m *MsgToggleTokenConversion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgToggleTokenConversion) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := fxtypes.ValidateEthereumAddress(m.Token); err != nil {
		if err = sdk.ValidateDenom(m.Token); err != nil {
			return errorsmod.Wrap(err, "token")
		}
	}
	return nil
}

func (m *MsgToggleTokenConversion) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// Route returns the MsgUpdateDenomAlias message route.
func (m *MsgUpdateDenomAlias) Route() string { return ModuleName }

// Type returns the MsgUpdateDenomAlias message type.
func (m *MsgUpdateDenomAlias) Type() string { return TypeMsgUpdateDenomAlias }

// GetSignBytes returns the raw bytes for a MsgUpdateDenomAlias message that
// the expected signer needs to sign.
func (m *MsgUpdateDenomAlias) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateDenomAlias) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return errorsmod.Wrap(err, "denom")
	}
	if err := sdk.ValidateDenom(m.Alias); err != nil {
		return errorsmod.Wrap(err, "alias")
	}
	return nil
}

func (m *MsgUpdateDenomAlias) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}
