package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) HandleOutgoingBridgeCallRefund(ctx sdk.Context, data *types.OutgoingBridgeCall) sdk.Coins {
	refund := types.ExternalAddrToAccAddr(k.moduleName, data.GetRefund())
	coins, err := k.bridgeCallTransferCoins(ctx, refund, data.Tokens)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallRefund,
		sdk.NewAttribute(types.AttributeKeyRefund, refund.String()),
	))

	// precompile bridge call, refund to evm
	if err = k.bridgeCallTransferTokens(ctx, refund, refund, coins); err != nil {
		panic(err)
	}
	return coins
}

func (k Keeper) DeleteOutgoingBridgeCallRecord(ctx sdk.Context, bridgeCallNonce uint64) {
	// 1. delete bridge call
	k.DeleteOutgoingBridgeCall(ctx, bridgeCallNonce)

	// 2. delete bridge call confirm
	k.DeleteBridgeCallConfirm(ctx, bridgeCallNonce)
}
