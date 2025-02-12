package types

import (
	"math"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
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
		return sdkerrors.ErrInvalidCoins.Wrapf("denom: %s", err.Error())
	}
	if m.Coin.Amount.IsNil() || !m.Coin.Amount.IsPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap("amount: must be positive")
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

func (m *MsgToggleTokenConversion) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority address: %s", err.Error())
	}
	if err := contract.ValidateEthereumAddress(m.Token); err != nil {
		if err = sdk.ValidateDenom(m.Token); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("token must be a valid denom or erc20 address")
		}
	}
	return nil
}

func (m *MsgRegisterNativeCoin) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority address: %s", err.Error())
	}

	if strings.TrimSpace(m.Name) == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("name: cannot be blank")
	}

	if err := sdk.ValidateDenom(strings.ToLower(m.Symbol)); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("symbol: %s", err.Error())
	}

	if m.Decimals > math.MaxUint8 {
		return sdkerrors.ErrInvalidRequest.Wrapf("overflow decimals: %d", m.Decimals)
	}

	return nil
}

func (m *MsgRegisterNativeERC20) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority address: %s", err.Error())
	}
	if err := contract.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("contract address: %s", err.Error())
	}
	return nil
}

func (m *MsgRegisterBridgeToken) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("authority address: %s", err.Error())
	}

	if err := sdk.ValidateDenom(m.BaseDenom); err != nil {
		return sdkerrors.ErrInvalidCoins.Wrapf("denom: %s", err.Error())
	}

	if len(m.Channel) > 0 || len(m.IbcDenom) > 0 {
		if !ibcchanneltypes.IsValidChannelID(m.Channel) {
			return sdkerrors.ErrInvalidRequest.Wrap("channel id")
		}
		if err := sdk.ValidateDenom(m.IbcDenom); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrapf("ibc denom: %s", err.Error())
		}
		return nil
	}

	if err := fxtypes.ValidateExternalAddr(m.ChainName, m.ContractAddress); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("contract address: %s", m.ContractAddress)
	}
	return nil
}
