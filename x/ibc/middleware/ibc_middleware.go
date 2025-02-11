package middleware

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/keeper"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
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
func (im IBCMiddleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	data, err := UnMarshalPacket(ctx.ChainID(), packet)
	if err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}

	var ack exported.Acknowledgement
	if data.Router != "" {
		ack = im.Keeper.OnRecvPacketWithRouter(ctx, im.IBCModule, packet, data, relayer)
	} else {
		ack = im.Keeper.OnRecvPacketWithoutRouter(ctx, im.IBCModule, packet, data, relayer)
	}

	return ack
}

func UnMarshalPacket(chainID string, packet channeltypes.Packet) (types.FungibleTokenPacketData, error) {
	var data types.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return data, err
	}

	if len(data.GetFee()) == 0 {
		data.Fee = sdkmath.ZeroInt().String()
	}

	needWrap, wrapDenom, packetDenom := fxtypes.OnRecvDenomNeedWrap(chainID, packet.SourcePort, packet.SourceChannel, data.Denom)
	if needWrap {
		data.Denom = packetDenom
		newAmount, err := fxtypes.OnRecvAmountCovert(wrapDenom, data.Amount)
		if err != nil {
			return data, err
		}
		data.Amount = newAmount
	}

	return data, data.ValidateBasic()
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.Keeper.OnAcknowledgementPacket(ctx, im.IBCModule, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.Keeper.OnTimeoutPacket(ctx, packet, im.IBCModule, relayer)
}
