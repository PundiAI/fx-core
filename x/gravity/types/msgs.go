package types

import (
	"fmt"

	"github.com/functionx/fx-core/x/crosschain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgSendToEth = "send_to_eth"
)

var (
	_ sdk.Msg = &MsgSendToEth{}
)

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
	if err := types.ValidateEthereumAddress(m.EthDest); err != nil {
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
