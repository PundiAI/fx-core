package keeper

import (
	"fmt"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
	polygontypes "github.com/pundiai/fx-core/v8/x/polygon/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
)

func (k Keeper) OnRecvPacketWithRouter(ctx sdk.Context, ibcModule porttypes.IBCModule, packet channeltypes.Packet, data types.FungibleTokenPacketData, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	receiver, transferAmount, feeAmount, err := parseReceiveAndAmountByPacketWithRouter(data)
	if err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}

	if err = fxtypes.ValidateExternalAddr(data.Router, data.Receiver); err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}

	receiveAmount := transferAmount.Add(feeAmount)
	packetData := transfertypes.NewFungibleTokenPacketData(data.GetDenom(), receiveAmount.String(), data.GetSender(), receiver.String(), data.Memo)
	packet.Data = packetData.GetBytes()
	onRecvPacketCtxWithNewEvent := ctx.WithEventManager(sdk.NewEventManager())
	ack := ibcModule.OnRecvPacket(onRecvPacketCtxWithNewEvent, packet, relayer)
	if ack == nil || !ack.Success() {
		return ack
	}

	receiveDenom := parseIBCCoinDenom(packet, data.GetDenom())
	receiveCoin := sdk.NewCoin(receiveDenom, receiveAmount)
	receiveCoin, err = k.compatibleWithOldData(ctx, receiver, receiveCoin)
	if err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}

	// bridgeCoinToBase()
	found, baseDenom, err := k.crosschainKeeper.IBCCoinToBaseCoin(ctx, receiver, receiveCoin)
	if err != nil {
		return types.NewAckErrorWithErrorEvent(ctx, err)
	}
	if !found {
		return types.NewAckErrorWithErrorEvent(ctx, fmt.Errorf("token not support"))
	}

	ibcAmount := sdk.NewCoin(baseDenom, transferAmount)
	ibcFee := sdk.NewCoin(baseDenom, feeAmount)

	routerCtxWithNewEvent := ctx.WithEventManager(sdk.NewEventManager())
	ctx = ctx.WithKVGasConfig(storetypes.GasConfig{}).WithTransientKVGasConfig(storetypes.GasConfig{})
	_, err = k.crosschaniRouterMsgServer.SendToExternal(ctx, &crosschaintypes.MsgSendToExternal{
		Sender:    receiver.String(),
		Dest:      data.Receiver,
		Amount:    ibcAmount,
		BridgeFee: ibcFee,
		ChainName: data.Router,
	})

	routerEvent := sdk.NewEvent(types.EventTypeReceiveRoute,
		sdk.NewAttribute(types.AttributeKeyRoute, data.Router),
		sdk.NewAttribute(types.AttributeKeySuccess, fmt.Sprintf("%t", err == nil)),
	)
	if err != nil {
		routerEvent = routerEvent.AppendAttributes(sdk.NewAttribute(types.AttributeKeyError, err.Error()))
		ack = channeltypes.NewErrorAcknowledgement(err)
	} else {
		ctx.EventManager().EmitEvents(onRecvPacketCtxWithNewEvent.EventManager().Events())
		ctx.EventManager().EmitEvents(routerCtxWithNewEvent.EventManager().Events())
	}
	ctx.EventManager().EmitEvent(routerEvent)
	return ack
}

func (k Keeper) compatibleWithOldData(ctx sdk.Context, holder sdk.AccAddress, receiveCoin sdk.Coin) (sdk.Coin, error) {
	// Try to confirm the bridge token
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, receiveCoin.Denom)
	if err != nil {
		return receiveCoin, nil
	}

	// If the denom is a bridge token, convert it to the base coin
	bridgeModule := ""
	if strings.HasPrefix(receiveCoin.Denom, polygontypes.ModuleName) {
		bridgeModule = polygontypes.ModuleName
	} else if strings.HasPrefix(receiveCoin.Denom, trontypes.ModuleName) {
		bridgeModule = trontypes.ModuleName
	}
	if len(bridgeModule) == 0 {
		return receiveCoin, nil
	}

	// if get bridge token fail, return the original coin
	bridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, bridgeModule, baseDenom)
	if err != nil {
		return receiveCoin, nil
	}

	baseCoin, err := k.crosschainKeeper.BridgeTokenToBaseCoin(ctx, holder, receiveCoin.Amount, bridgeToken)
	if err != nil {
		return receiveCoin, err
	}

	// return the base coin
	return sdk.NewCoin(baseCoin.Denom, receiveCoin.Amount), nil
}
