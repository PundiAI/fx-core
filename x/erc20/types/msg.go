package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
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
		return sdkerrors.ErrInvalidAddress.Wrapf("sender address: %s", err.Error())
	}
	if err = contract.ValidateEthereumAddress(m.Receiver); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("receiver address: %s", err.Error())
	}
	if err = ibctransfertypes.ValidateIBCDenom(m.Coin.Denom); err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("coin denom: %s", err.Error())
	}
	if m.Coin.Amount.IsNil() || !m.Coin.Amount.IsPositive() {
		return sdkerrors.ErrInvalidRequest.Wrap("amount")
	}
	return nil
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority address: %s", err.Error())
	}
	if err := m.Params.Validate(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("params: %s", err.Error())
	}
	return nil
}

func (m *MsgRegisterCoin) ValidateBasic() error {
	return nil
}

func (m *MsgRegisterERC20) ValidateBasic() error {
	return nil
}

func (m *MsgToggleTokenConversion) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority address: %s", err.Error())
	}
	if err := contract.ValidateEthereumAddress(m.Token); err != nil {
		if err = sdk.ValidateDenom(m.Token); err != nil {
			return sdkerrors.ErrInvalidCoins.Wrapf("token denom: %s", err.Error())
		}
	}
	return nil
}
