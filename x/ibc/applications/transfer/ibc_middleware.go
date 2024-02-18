package transfer

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"

	"github.com/functionx/fx-core/v7/x/ibc"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/keeper"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

var _ porttypes.Middleware = &IBCMiddleware{}

// IBCMiddleware implements the ICS26 interface for transfer given the transfer keeper.
type IBCMiddleware struct {
	*ibc.Module
	keeper keeper.Keeper
}

func (im IBCMiddleware) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	return im.keeper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

func (im IBCMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return im.keeper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (im IBCMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return im.keeper.GetAppVersion(ctx, portID, channelID)
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application
func NewIBCMiddleware(k keeper.Keeper, app porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{
		Module: ibc.NewModule(app),
		keeper: k,
	}
}

// OnRecvPacket implements the IBCModule interface
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) exported.Acknowledgement {
	ack := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	var data types.FungibleTokenPacketData
	var ackErr error
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		ackErr = sdkerrors.ErrInvalidType.Wrap("cannot unmarshal ICS-20 transfer packet data")
		ack = channeltypes.NewErrorAcknowledgement(ackErr)
	}

	// only attempt the application logic if the packet data
	// was successfully decoded
	if ack.Success() {
		if len(data.GetFee()) == 0 {
			data.Fee = sdkmath.ZeroInt().String()
		}
		if err := im.keeper.OnRecvPacket(ctx, packet, data); err != nil {
			ack = channeltypes.NewErrorAcknowledgement(err)
			ackErr = err
		}
	}

	event := sdk.NewEvent(
		transfertypes.EventTypePacket,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, data.Sender),
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
		sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, data.Amount),
		sdk.NewAttribute(transfertypes.AttributeKeyMemo, data.Memo),
		sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", ack.Success())),
	)

	if ackErr != nil {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyRecvError, ackErr.Error()))
	}
	ctx.EventManager().EmitEvent(event)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	_ sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data types.FungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	if err := im.keeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		transfertypes.EventTypePacket,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, data.Sender),
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
		sdk.NewAttribute(transfertypes.AttributeKeyDenom, data.Denom),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, data.Amount),
		sdk.NewAttribute(transfertypes.AttributeKeyMemo, data.Memo),
		sdk.NewAttribute(transfertypes.AttributeKeyAck, ack.String()),
	))

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			transfertypes.EventTypePacket,
			sdk.NewAttribute(transfertypes.AttributeKeyAckSuccess, string(resp.Result)),
		))
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			transfertypes.EventTypePacket,
			sdk.NewAttribute(transfertypes.AttributeKeyAckError, resp.Error),
		))
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) error {
	var data types.FungibleTokenPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	// refund tokens
	if err := im.keeper.OnTimeoutPacket(ctx, packet, data); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		transfertypes.EventTypeTimeout,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(transfertypes.AttributeKeyRefundReceiver, data.Sender),
		sdk.NewAttribute(transfertypes.AttributeKeyRefundDenom, data.Denom),
		sdk.NewAttribute(transfertypes.AttributeKeyRefundAmount, data.Amount),
		sdk.NewAttribute(transfertypes.AttributeKeyMemo, data.Memo),
	))

	return nil
}
