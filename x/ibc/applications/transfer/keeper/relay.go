package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	coretypes "github.com/cosmos/ibc-go/v6/modules/core/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

// make SendTransfer private
// https://github.com/cosmos/ibc-go/pull/2446
func (k Keeper) sendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.Coin, sender sdk.AccAddress,
	receiver string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin, memo string,
) (uint64, error) {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return 0, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// NOTE: denomination and hex hash correctness checked during msg.ValidateBasic
	fullDenomPath := token.Denom

	var err error

	// deconstruct the token denomination into the denomination trace info
	// to determine if the sender is the source chain
	if strings.HasPrefix(token.Denom, "ibc/") {
		fullDenomPath, err = k.DenomPathFromHash(ctx, token.Denom)
		if err != nil {
			return 0, err
		}
	}

	labels := []metrics.Label{
		telemetry.NewLabel(coretypes.LabelDestinationPort, destinationPort),
		telemetry.NewLabel(coretypes.LabelDestinationChannel, destinationChannel),
	}

	packetData := types.NewFungibleTokenPacketData(
		fullDenomPath, token.Amount.String(), sender.String(), receiver, router, fee.Amount.String(),
	)

	packetData.Memo = memo
	// If the router address is specified, the number of token + fee is deducted
	if router != "" {
		token = token.Add(sdk.NewCoin(token.Denom, fee.Amount))
	}
	// NOTE: SendTransfer simply sends the denomination as it exists on its own
	// chain inside the packet data. The receiving chain will perform denom
	// prefixing as necessary.
	if transfertypes.SenderChainIsSource(sourcePort, sourceChannel, fullDenomPath) {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "true"))

		// create the escrow address for the tokens
		escrowAddress := transfertypes.GetEscrowAddress(sourcePort, sourceChannel)

		// escrow source tokens. It fails if balance insufficient.
		if err = k.bankKeeper.SendCoins(ctx, sender, escrowAddress, sdk.NewCoins(token)); err != nil {
			return 0, err
		}

	} else {
		labels = append(labels, telemetry.NewLabel(coretypes.LabelSource, "false"))

		// transfer the coins to the module account and burn them
		if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, transfertypes.ModuleName, sdk.NewCoins(token)); err != nil {
			return 0, err
		}

		if err = k.bankKeeper.BurnCoins(ctx, transfertypes.ModuleName, sdk.NewCoins(token)); err != nil {
			// NOTE: should not happen as the module account was
			// retrieved on the step above and it has enough balance
			// to burn.
			panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
		}
	}

	sequence, err := k.ics4Wrapper.SendPacket(
		ctx,
		channelCap,
		sourcePort,
		sourceChannel,
		timeoutHeight,
		timeoutTimestamp,
		packetData.GetBytes(),
	)
	if err != nil {
		return 0, err
	}

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"ibc", transfertypes.ModuleName, "send"},
			1,
			labels,
		)
	}()

	return sequence, nil
}

// OnRecvPacket processes a cross chain fungible token transfer. If the
// sender chain is the source of minted tokens then vouchers will be minted
// and sent to the receiving address. Otherwise if the sender chain is sending
// back tokens this chain originally transferred to it, the tokens are
// unescrowed and sent to the receiving address.
func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData) error {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	receiver, isEvmAddr, transferAmount, feeAmount, err := parseReceiveAndAmountByPacket(data)
	if err != nil {
		return err
	}

	receiveAmount := transferAmount.Add(feeAmount)
	packetData := transfertypes.NewFungibleTokenPacketData(data.GetDenom(), receiveAmount.String(), data.GetSender(), receiver.String(), "")
	packetData.Memo = data.Memo
	onRecvPacketCtxWithNewEvent := ctx.WithEventManager(sdk.NewEventManager())
	if err = k.Keeper.OnRecvPacket(onRecvPacketCtxWithNewEvent, packet, packetData); err != nil {
		return err
	}

	receiveDenom := parseIBCCoinDenom(packet, data.GetDenom())

	receiveCoin := sdk.NewCoin(receiveDenom, receiveAmount)
	onRecvPacketCtxWithNewEvent.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReceive,
		sdk.NewAttribute(transfertypes.AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(transfertypes.AttributeKeyAmount, receiveCoin.String()),
	))

	if data.Router == "" || k.router == nil {
		// try to send to evm module
		if receiveCoin.GetDenom() != fxtypes.DefaultDenom && isEvmAddr {
			// convert to base denom
			receiveCoin, err = k.erc20Keeper.ConvertDenomToTarget(onRecvPacketCtxWithNewEvent, receiver, receiveCoin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
			if err != nil {
				return err
			}
			// convert to erc20 token
			_, err = k.erc20Keeper.ConvertCoin(sdk.WrapSDKContext(onRecvPacketCtxWithNewEvent), &erc20types.MsgConvertCoin{
				Coin:     receiveCoin,
				Receiver: common.BytesToAddress(receiver).String(),
				Sender:   receiver.String(),
			})
			if err != nil {
				return err
			}
		}

		// NOTE: if not router, emit onRecvPacketCtx event, only error is nil emit
		ctx.EventManager().EmitEvents(onRecvPacketCtxWithNewEvent.EventManager().Events())

		// ibc call
		if len(data.Memo) > 0 {
			if err = k.HandlerIbcCall(ctx, packet.SourcePort, packet.SourceChannel, data); err != nil {
				return err
			}
		}
		return nil
	}
	route, exists := k.router.GetRoute(data.Router)
	if !exists {
		return errorsmod.Wrap(types.ErrRouterNotFound, data.Router)
	}

	fxTarget := fxtypes.ParseFxTarget(data.Router)
	targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(onRecvPacketCtxWithNewEvent, receiver, receiveCoin, fxTarget)
	if err != nil {
		return err
	}

	ibcAmount := sdk.NewCoin(targetCoin.GetDenom(), transferAmount)
	ibcFee := sdk.NewCoin(targetCoin.GetDenom(), feeAmount)

	routerCtxWithNewEvent := ctx.WithEventManager(sdk.NewEventManager())
	err = route.TransferAfter(routerCtxWithNewEvent, receiver, data.Receiver, ibcAmount, ibcFee, true)
	routerEvent := sdk.NewEvent(types.EventTypeReceiveRoute,
		sdk.NewAttribute(types.AttributeKeyRoute, data.Router),
		sdk.NewAttribute(types.AttributeKeyRouteSuccess, fmt.Sprintf("%t", err == nil)),
	)
	if err != nil {
		routerEvent = routerEvent.AppendAttributes(sdk.NewAttribute(types.AttributeKeyRouteError, err.Error()))
	} else {
		ctx.EventManager().EmitEvents(onRecvPacketCtxWithNewEvent.EventManager().Events())
		ctx.EventManager().EmitEvents(routerCtxWithNewEvent.EventManager().Events())
	}
	ctx.EventManager().EmitEvent(routerEvent)
	return err
}

// OnAcknowledgementPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain. If the acknowledgement
// was a success then nothing occurs. If the acknowledgement failed, then
// the sender is refunded their tokens using the refundPacketToken function.
func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		amount, fee, err := parseAmountAndFeeByPacket(data)
		if err != nil {
			return err
		}
		ibcPacketData := data.ToIBCPacketData()
		ibcPacketData.Amount = amount.Add(fee).String()
		if err = k.Keeper.OnAcknowledgementPacket(ctx, packet, ibcPacketData, ack); err != nil {
			return err
		}
		return k.refundPacketTokenHook(ctx, packet, data, amount, fee)
	default:
		if k.refundHook != nil {
			if err := k.refundHook.AckAfter(ctx, packet.SourceChannel, packet.Sequence); err != nil {
				k.Logger(ctx).Error("acknowledgement packet hook error", "sourceChannel", packet.GetSourceChannel(), "destChannel", packet.GetDestChannel(), "sequence", packet.GetSequence(), "error", err)
			}
		}
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

// OnTimeoutPacket refunds the sender since the original packet sent was
// never received and has been timed out.
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData) error {
	amount, fee, err := parseAmountAndFeeByPacket(data)
	if err != nil {
		return err
	}
	ibcPacketData := data.ToIBCPacketData()
	ibcPacketData.Amount = amount.Add(fee).String()
	if err = k.Keeper.OnTimeoutPacket(ctx, packet, ibcPacketData); err != nil {
		return err
	}
	return k.refundPacketTokenHook(ctx, packet, data, amount, fee)
}

// refundPacketToken will unescrow and send back the tokens back to sender
// if the sending chain was the source chain. Otherwise, the sent tokens
// were burnt in the original send so new tokens are minted and sent to
// the sending address.
func (k Keeper) refundPacketTokenHook(ctx sdk.Context, packet channeltypes.Packet, data types.FungibleTokenPacketData, amount sdkmath.Int, fee sdkmath.Int) error {
	// parse the denomination from the full denom path
	trace := transfertypes.ParseDenomTrace(data.Denom)

	amount = amount.Add(fee)
	token := sdk.NewCoin(trace.IBCDenom(), amount)

	// decode the sender address
	sender, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	if k.refundHook != nil {
		if err = k.refundHook.RefundAfter(ctx, packet.SourceChannel, packet.Sequence, sender, token); err != nil {
			k.Logger(ctx).Info("refundPacketToken hook err", "sourceChannel", packet.GetSourceChannel(), "destChannel", packet.GetDestChannel(), "sequence", packet.GetSequence(), "error", err)
		}
	}
	return nil
}
