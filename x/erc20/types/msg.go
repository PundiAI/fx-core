package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

var (
	_ sdk.Msg = &MsgConvertCoin{}
	_ sdk.Msg = &MsgConvertERC20{}
)

const (
	TypeMsgConvertCoin  = "convert_coin"
	TypeMsgConvertERC20 = "convert_ERC20"
	TypeMsgConvertDenom = "convert_denom"
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
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if err = fxtypes.ValidateEthereumAddress(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
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
	addr := sdk.MustAccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}

// NewMsgConvertERC20 creates a new instance of MsgConvertERC20
func NewMsgConvertERC20(amount sdk.Int, receiver sdk.AccAddress, contract, sender common.Address) *MsgConvertERC20 {
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
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
	}
	if err := fxtypes.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid contract address: %s", err)
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
	addr := common.HexToAddress(m.Sender)
	return []sdk.AccAddress{addr.Bytes()}
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
		return errortypes.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receiver address: %s", err)
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
	addr := sdk.MustAccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}
