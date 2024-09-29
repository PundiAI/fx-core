package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	"github.com/functionx/fx-core/v8/x/ibc/middleware/types"
)

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

	if receiveCoin.GetDenom() != fxtypes.DefaultDenom && isEvmAddr {
		// convert to base denom
		receiveCoin, err = k.erc20Keeper.ConvertDenomToTarget(ctx, receiver, receiveCoin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
		if err != nil {
			return err
		}
		// convert to erc20 token
		_, err = k.erc20Keeper.ConvertCoin(ctx, &erc20types.MsgConvertCoin{
			Coin:     receiveCoin,
			Receiver: common.BytesToAddress(receiver).String(),
			Sender:   receiver.String(),
		})
		if err != nil {
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

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketTokenHook(ctx, packet, data)
	default:
		if k.refundHook != nil {
			k.refundHook.AckAfter(ctx, packet.SourceChannel, packet.Sequence)
		}
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

// OnTimeoutPacket refunds the sender since the original packet sent was
// never received and has been timed out.
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data transfertypes.FungibleTokenPacketData) error {
	return k.refundPacketTokenHook(ctx, packet, data)
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

	if k.refundHook != nil {
		k.refundHook.RefundAfter(ctx, packet.SourceChannel, packet.Sequence, sender, token)
	}
	return nil
}
