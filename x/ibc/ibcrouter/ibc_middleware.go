package ibcrouter

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/tendermint/tendermint/libs/log"

	fxtransfertypes "github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
	"github.com/functionx/fx-core/v7/x/ibc/ibcrouter/parser"
	"github.com/functionx/fx-core/v7/x/ibc/ibcrouter/types"
)

const (
	ForwardPacketTimeHour time.Duration = 12
)

var _ porttypes.Middleware = &IBCMiddleware{}

// IBCMiddleware implements the ICS26 interface for transfer given the transfer keeper.
type IBCMiddleware struct {
	app            porttypes.IBCModule
	transferKeeper types.TransferKeeper
	ics4Wrapper    porttypes.ICS4Wrapper
}

// ICS 30 callbacks

// OnChanOpenInit implements the IBCModule interface
func (im IBCMiddleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, chanCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version)
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCMiddleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, chanCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string,
) (version string, err error) {
	return im.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID,
		chanCap, counterparty, counterpartyVersion)
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// NewIBCMiddleware creates a new IBCMiddleware given the keeper and underlying application
func NewIBCMiddleware(app porttypes.IBCModule, ics4Wrapper porttypes.ICS4Wrapper, transferKeeper transferkeeper.Keeper) IBCMiddleware {
	return IBCMiddleware{
		app:            app,
		transferKeeper: transferKeeper,
		ics4Wrapper:    ics4Wrapper,
	}
}

// Logger returns a module-specific logger.
func (im IBCMiddleware) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+"ibcroutermiddleware")
}

// OnRecvPacket implements the IBCModule interface
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	var data fxtransfertypes.FungibleTokenPacketData
	if err := fxtransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return channeltypes.NewErrorAcknowledgement(fmt.Errorf("cannot unmarshal ICS-20 transfer packet data"))
	}
	// check the packet has router
	if len(data.Router) > 0 {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	ack, err := handlerForwardTransferPacket(ctx, im, packet, data.ToIBCPacketData(), relayer)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(err)
	}

	if err != nil {
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRouter,
			sdk.NewAttribute(types.AttributeKeyRouteError, err.Error()),
		))
	}
	return ack
}

func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return im.app.OnTimeoutPacket(ctx, packet, relayer)
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
	return im.ics4Wrapper.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

func (im IBCMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	return im.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

func (im IBCMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return im.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
}

func handlerForwardTransferPacket(ctx sdk.Context, im IBCMiddleware, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData, relayer sdk.AccAddress) (exported.Acknowledgement, error) {
	// parse out any forwarding info
	parsedReceiver, err := parser.ParseReceiverData(data.Receiver)
	if err != nil {
		return nil, err
	}

	if !parsedReceiver.ShouldForward {
		return im.app.OnRecvPacket(ctx, packet, relayer), nil
	}

	newData := data
	newData.Receiver = parsedReceiver.HostAccAddr.String()

	newPacket := packet
	newPacket.Data = newData.GetBytes()

	ack := im.app.OnRecvPacket(ctx, newPacket, relayer)
	if ack.Success() {
		// recalculate denom, skip checks that were already done in app.OnRecvPacket
		denom := GetDenomByIBCPacket(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetDestPort(), packet.GetDestChannel(), newData.GetDenom())
		// parse the transfer amount
		transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
		if !ok {
			return nil, errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse forward transfer amount (%s) into sdkmath.Int", data.Amount)
		}

		token := sdk.NewCoin(denom, transferAmount)
		msgTransfer := transfertypes.NewMsgTransfer(
			parsedReceiver.Port,
			parsedReceiver.Channel,
			token,
			parsedReceiver.HostAccAddr.String(),
			parsedReceiver.Destination,
			clienttypes.Height{},
			uint64(ctx.BlockTime().Add(ForwardPacketTimeHour*time.Hour).UnixNano()),
			"",
		)
		msgTransfer.Memo = newData.Memo

		// send tokens to destination
		if _, err = im.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), msgTransfer); err != nil {
			return nil, errorsmod.Wrapf(errortypes.ErrInsufficientFunds, err.Error())
		}
	}
	return ack, nil
}

func GetDenomByIBCPacket(sourcePort, sourceChannel, destPort, destChannel, packetDenom string) string {
	var denom string

	if transfertypes.ReceiverChainIsSource(sourcePort, sourceChannel, packetDenom) {
		voucherPrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
		unPrefixedDenom := packetDenom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom = unPrefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unPrefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}
	} else {
		// since SendPacket did not prefix the denomination, we must prefix denomination here
		sourcePrefix := transfertypes.GetDenomPrefix(destPort, destChannel)
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + packetDenom

		// construct the denomination trace from the full raw denomination
		denom = transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()
	}
	return denom
}
