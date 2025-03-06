package keeper

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

func (k Keeper) OnRecvPacketWithoutRouter(ctx sdk.Context, ibcModule porttypes.IBCModule, packet channeltypes.Packet, data types.FungibleTokenPacketData, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// parse receive address, compatible with evm addresses
	receiver, _, err := fxtypes.ParseAddress(data.Receiver)
	if err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}

	// Only the receiver is replaced with the fx address, which is compatible with the evm address
	newPacketData := data.ToIBCPacketData()
	newPacketData.Receiver = receiver.String()

	newPacket := packet
	newPacket.Data = newPacketData.GetBytes()
	ack := ibcModule.OnRecvPacket(ctx, newPacket, relayer)
	if ack == nil || !ack.Success() {
		return ack
	}

	// Use the original package to handle ibc to evm
	if err = k.OnRecvPacket(zeroGasConfigCtx(ctx), packet, data.ToIBCPacketData()); err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}

	return ack
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	// parse the transfer amount
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount: %s", data.Amount)
	}

	receiveDenom := parseIBCCoinDenom(packet, data.GetDenom())
	receiveCoin := sdk.NewCoin(receiveDenom, transferAmount)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReceive,
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, data.Receiver),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, receiveCoin.String()),
	))

	if err := k.crosschainKeeper.IBCCoinToEvm(ctx, data.Receiver, receiveCoin); err != nil {
		return err
	}
	return k.HandlerIbcCall(ctx, packet.SourcePort, packet.SourceChannel, data)
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, ibcModule porttypes.IBCModule, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}

	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		ibcPacketData, isZeroAmount, err := UnmarshalAckPacketData(ctx.ChainID(), packet.SourceChannel, packet.Data)
		if err != nil {
			return err
		}
		if isZeroAmount {
			return k.crosschainKeeper.AfterIBCAckSuccess(zeroGasConfigCtx(ctx), packet.SourceChannel, packet.Sequence)
		}
		packet.Data = ibcPacketData.GetBytes()
		if err = ibcModule.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer); err != nil {
			return err
		}
		return k.refundPacketTokenHook(ctx, packet, ibcPacketData)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return k.crosschainKeeper.AfterIBCAckSuccess(zeroGasConfigCtx(ctx), packet.SourceChannel, packet.Sequence)
	}
}

// OnTimeoutPacket refunds the sender since the original packet sent was
// never received and has been timed out.
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, ibcModule porttypes.IBCModule, relayer sdk.AccAddress) error {
	ibcPacketData, isZeroAmount, err := UnmarshalAckPacketData(ctx.ChainID(), packet.SourceChannel, packet.Data)
	if err != nil {
		return err
	}
	if isZeroAmount {
		return k.crosschainKeeper.AfterIBCAckSuccess(zeroGasConfigCtx(ctx), packet.SourceChannel, packet.Sequence)
	}
	packet.Data = ibcPacketData.GetBytes()
	if err = ibcModule.OnTimeoutPacket(ctx, packet, relayer); err != nil {
		return err
	}
	return k.refundPacketTokenHook(ctx, packet, ibcPacketData)
}

// UnmarshalAckPacketData unmarshal ack packet data
// @param packetData []byte
// @return transfertypes.FungibleTokenPacketData
// @return bool isZeroAmount
// @return error
func UnmarshalAckPacketData(chainId, sourceChannel string, packetData []byte) (transfertypes.FungibleTokenPacketData, bool, error) {
	var data types.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packetData, &data); err != nil {
		return transfertypes.FungibleTokenPacketData{}, false, sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}
	amount, fee, err := parseAmountAndFeeByPacket(data)
	if err != nil {
		return transfertypes.FungibleTokenPacketData{}, false, err
	}
	ibcPacketData := data.ToIBCPacketData()
	totalAmount := amount.Add(fee)
	if ibcPacketData.Denom == fxtypes.LegacyFXDenom {
		ibcPacketData.Denom = fxtypes.DefaultDenom
		totalAmount = fxtypes.SwapAmount(totalAmount)
	}
	denomNeedWrap, wrapDenom := fxtypes.AckPacketDenomNeedWrap(chainId, sourceChannel, ibcPacketData.Denom)
	if denomNeedWrap {
		ibcPacketData.Denom = wrapDenom
	}
	ibcPacketData.Amount = totalAmount.String()
	return ibcPacketData, !totalAmount.IsPositive(), nil
}

func (k Keeper) refundPacketTokenHook(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount (%s) into sdkmath.Int", data.Amount)
	}
	// parse the denomination from the full denom path
	trace := transfertypes.ParseDenomTrace(data.Denom)
	token := sdk.NewCoin(trace.IBCDenom(), transferAmount)

	// decode the sender address
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}
	return k.crosschainKeeper.IBCCoinRefund(zeroGasConfigCtx(ctx), sender, token, packet.SourceChannel, packet.Sequence)
}

// zeroGasConfigCtx returns a context with a zero gas meter
// use a zero gas config to avoid extra costs for the relayers
func zeroGasConfigCtx(ctx sdk.Context) sdk.Context {
	return ctx.
		WithKVGasConfig(storetypes.GasConfig{}).
		WithTransientKVGasConfig(storetypes.GasConfig{})
}
