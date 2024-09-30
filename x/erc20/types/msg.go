package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
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

func NewMsgConvertCoin(coin sdk.Coin, receiver common.Address, sender sdk.AccAddress) *MsgConvertCoin {
	return &MsgConvertCoin{
		Coin:     coin,
		Receiver: receiver.Hex(),
		Sender:   sender.String(),
	}
}

func (m *MsgConvertCoin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if err = contract.ValidateEthereumAddress(m.Receiver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err.Error())
	}
	if err = ibctransfertypes.ValidateIBCDenom(m.Coin.Denom); err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid coin denom %s", err.Error())
	}
	if m.Coin.Amount.IsNil() || !m.Coin.Amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

func NewMsgConvertERC20(amount sdkmath.Int, receiver sdk.AccAddress, contract, sender common.Address) *MsgConvertERC20 {
	return &MsgConvertERC20{
		ContractAddress: contract.String(),
		Amount:          amount,
		Receiver:        receiver.String(),
		Sender:          sender.Hex(),
	}
}

func (m *MsgConvertERC20) ValidateBasic() error {
	if err := contract.ValidateEthereumAddress(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err.Error())
	}
	if err := contract.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid contract address: %s", err.Error())
	}
	if m.Amount.IsNil() || !m.Amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
}

func NewMsgConvertDenom(sender, receiver sdk.AccAddress, coin sdk.Coin, target string) *MsgConvertDenom {
	return &MsgConvertDenom{
		Sender:   sender.String(),
		Receiver: receiver.String(),
		Coin:     coin,
		Target:   target,
	}
}

func (m *MsgConvertDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err.Error())
	}
	if !m.Coin.IsValid() || !m.Coin.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid amount")
	}
	return nil
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

func (m *MsgRegisterERC20) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := contract.ValidateEthereumAddress(m.Erc20Address); err != nil {
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

func (m *MsgToggleTokenConversion) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := contract.ValidateEthereumAddress(m.Token); err != nil {
		if err = sdk.ValidateDenom(m.Token); err != nil {
			return errorsmod.Wrap(err, "token")
		}
	}
	return nil
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
