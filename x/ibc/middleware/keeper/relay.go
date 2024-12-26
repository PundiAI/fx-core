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
		return channeltypes.NewErrorAcknowledgement(err)
	}

	newPacketData := data.ToIBCPacketData()
	newPacketData.Receiver = receiver.String()

	newPacket := packet
	newPacket.Data = newPacketData.GetBytes()
	ack := ibcModule.OnRecvPacket(ctx, newPacket, relayer)
	if ack == nil || !ack.Success() {
		return ack
	}

	if err = k.OnRecvPacket(zeroGasConfigCtx(ctx), newPacket, newPacketData); err != nil {
		return channeltypes.NewErrorAcknowledgement(err)
	}

	return ack
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	// parse receive address, compatible with evm addresses
	receiver, isEvmAddr, err := fxtypes.ParseAddress(data.Receiver)
	if err != nil {
		return err
	}

	// parse the transfer amount
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return transfertypes.ErrInvalidAmount.Wrapf("unable to parse transfer amount: %s", data.Amount)
	}

	receiveDenom := parseIBCCoinDenom(packet, data.GetDenom())
	receiveCoin := sdk.NewCoin(receiveDenom, transferAmount)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReceive,
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, receiveCoin.String()),
	))

	if receiveCoin.GetDenom() != fxtypes.DefaultDenom {
		if !isEvmAddr {
			return sdkerrors.ErrInvalidAddress.Wrap("only support hex address")
		}
		if err = k.crosschainKeeper.IBCCoinToEvm(ctx, receiver, receiveCoin); err != nil {
			return err
		}
	}

	// ibc call
	if len(data.Memo) > 0 {
		if err = k.HandlerIbcCall(ctx, packet.SourcePort, packet.SourceChannel, data); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, ibcModule porttypes.IBCModule, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	var ack channeltypes.Acknowledgement
	if err := transfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data types.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		amount, fee, err := parseAmountAndFeeByPacket(data)
		if err != nil {
			return err
		}
		ibcPacketData := data.ToIBCPacketData()
		ibcPacketData.Amount = amount.Add(fee).String()
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
	var data types.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.ErrUnknownRequest.Wrapf("cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	amount, fee, err := parseAmountAndFeeByPacket(data)
	if err != nil {
		return err
	}
	ibcPacketData := data.ToIBCPacketData()
	ibcPacketData.Amount = amount.Add(fee).String()
	packet.Data = ibcPacketData.GetBytes()

	if err = ibcModule.OnTimeoutPacket(ctx, packet, relayer); err != nil {
		return err
	}
	return k.refundPacketTokenHook(ctx, packet, ibcPacketData)
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
