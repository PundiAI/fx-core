package middleware

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/ibc/middleware/keeper"
)

var _ porttypes.Middleware = &IBCMiddleware{}

// IBCMiddleware implements the ICS26 interface for transfer given the transfer keeper.
type IBCMiddleware struct {
	porttypes.IBCModule
	porttypes.ICS4Wrapper
	Keeper keeper.Keeper
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application
func NewIBCMiddleware(k keeper.Keeper, ics porttypes.ICS4Wrapper, ibcModule porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{
		IBCModule:   ibcModule,
		ICS4Wrapper: ics,
		Keeper:      k,
	}
}

// OnRecvPacket implements the IBCModule interface
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		err = sdkerrors.ErrInvalidType.Wrap("cannot unmarshal ICS-20 transfer packet data")
		return channeltypes.NewErrorAcknowledgement(err)
	}

	// parse receive address, compatible with evm addresses
	receiver, _, err := fxtypes.ParseAddress(data.Receiver)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}
	newPacketData := data
	newPacketData.Receiver = receiver.String()

	newPacket := packet
	newPacket.Data = newPacketData.GetBytes()

	ack := im.IBCModule.OnRecvPacket(ctx, newPacket, relayer)

	// return if the acknowledgement is an error ACK
	if !ack.Success() {
		return ack
	}

	if err = im.Keeper.OnRecvPacket(zeroGasConfigCtx(ctx), packet, data); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	if err := im.IBCModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer); err != nil {
		return err
	}
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if err := im.Keeper.OnAcknowledgementPacket(zeroGasConfigCtx(ctx), packet, data, ack); err != nil {
		return err
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	if err := im.IBCModule.OnTimeoutPacket(ctx, packet, relayer); err != nil {
		return err
	}
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if err := im.Keeper.OnTimeoutPacket(zeroGasConfigCtx(ctx), packet, data); err != nil {
		return err
	}

	return nil
}

// zeroGasConfigCtx returns a context with a zero gas meter
// use a zero gas config to avoid extra costs for the relayers
func zeroGasConfigCtx(ctx sdk.Context) sdk.Context {
	return ctx.
		WithKVGasConfig(storetypes.GasConfig{}).
		WithTransientKVGasConfig(storetypes.GasConfig{})
}
