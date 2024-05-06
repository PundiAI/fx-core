package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) HandleOutgoingBridgeCallRefund(ctx sdk.Context, data *types.OutgoingBridgeCall) {
	receiveAddr := types.ExternalAddressToAccAddress(k.moduleName, data.GetSender())
	coins, err := k.bridgeCallTransferToSender(ctx, receiveAddr, data.Tokens)
	if err != nil {
		panic(err)
	}

	if k.HasBridgeCallFromMsg(ctx, data.Nonce) {
		return
	}
	// precompile bridge call refund
	if err = k.bridgeCallTransferToReceiver(ctx, receiveAddr, receiveAddr, coins); err != nil {
		panic(err)
	}
}

func (k Keeper) DeleteOutgoingBridgeCallRecord(ctx sdk.Context, bridgeCallNonce uint64) {
	// 1. delete bridge call
	k.DeleteOutgoingBridgeCall(ctx, bridgeCallNonce)

	// 2. delete bridge call confirm
	k.DeleteBridgeCallConfirm(ctx, bridgeCallNonce)

	// 3. delete bridge call from msg
	k.DeleteBridgeCallFromMsg(ctx, bridgeCallNonce)
}
