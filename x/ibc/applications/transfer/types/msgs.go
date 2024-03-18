package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

const (
	TypeMsgTransfer = "transfer"
)

// NewMsgTransfer creates a new MsgTransfer instance
func NewMsgTransfer(
	sourcePort, sourceChannel string,
	token sdk.Coin, sender string, receiver string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin,
) *MsgTransfer {
	return &MsgTransfer{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Token:            token,
		Sender:           sender,
		Receiver:         receiver,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Router:           router,
		Fee:              fee,
	}
}

// Route implements sdk.Msg
func (MsgTransfer) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgTransfer) Type() string {
	return TypeMsgTransfer
}

// ValidateBasic performs a basic check of the MsgTransfer fields.
// NOTE: timeout height or timestamp values can be 0 to disable the timeout.
// NOTE: The recipient addresses format is not validated as the format defined by
// the chain is not known to IBC.
func (msg MsgTransfer) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return errorsmod.Wrap(err, "invalid source port ID")
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return errorsmod.Wrap(err, "invalid source channel ID")
	}
	if !msg.Token.IsValid() {
		return errorsmod.Wrap(errortypes.ErrInvalidCoins, msg.Token.String())
	}
	if !msg.Token.IsPositive() {
		return errorsmod.Wrap(errortypes.ErrInsufficientFunds, msg.Token.String())
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(errortypes.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	if strings.TrimSpace(msg.Receiver) == "" {
		return errorsmod.Wrap(errortypes.ErrInvalidAddress, "missing recipient address")
	}
	if msg.Fee.Amount.IsNil() || !msg.Fee.IsValid() {
		return errorsmod.Wrap(errortypes.ErrInvalidCoins, "fees")
	}
	if msg.Fee.Denom != msg.Token.Denom {
		return errorsmod.Wrap(ErrFeeDenomNotMatchTokenDenom, fmt.Sprintf("token denom:%s, fee denom:%s", msg.Token.Denom, msg.Fee.Denom))
	}
	return transfertypes.ValidateIBCDenom(msg.Token.Denom)
}

// GetSignBytes implements sdk.Msg.
func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Sender)}
}
