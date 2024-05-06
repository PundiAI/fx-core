package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) HasBridgeCallConfirm(ctx sdk.Context, nonce uint64, oracleAddr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBridgeCallConfirmKey(nonce, oracleAddr))
}

func (k Keeper) SetBridgeCallConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, msg *types.MsgBridgeCallConfirm) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeCallConfirmKey(msg.Nonce, oracleAddr), k.cdc.MustMarshal(msg))
}

func (k Keeper) IterBridgeCallConfirmByNonce(ctx sdk.Context, nonce uint64, cb func(msg *types.MsgBridgeCallConfirm) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBridgeCallConfirmNonceKey(nonce))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgBridgeCallConfirm)
		k.cdc.MustUnmarshal(iter.Value(), confirm)
		// cb returns true to stop early
		if cb(confirm) {
			break
		}
	}
}

func (k Keeper) DeleteBridgeCallConfirm(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetBridgeCallConfirmNonceKey(nonce))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

func (k Keeper) SetLastSlashedBridgeCallNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedBridgeCallNonce, sdk.Uint64ToBigEndian(nonce))
}

func (k Keeper) GetLastSlashedBridgeCallNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(types.LastSlashedBridgeCallNonce))
}

func (k Keeper) GetUnSlashedBridgeCalls(ctx sdk.Context, height uint64) []types.OutgoingBridgeCall {
	nonce := k.GetLastSlashedBridgeCallNonce(ctx)
	var bridgeCalls []types.OutgoingBridgeCall
	k.IterateOutgoingBridgeCallByNonce(ctx, nonce, func(bridgeCall *types.OutgoingBridgeCall) bool {
		if bridgeCall.BlockHeight <= height {
			bridgeCalls = append(bridgeCalls, *bridgeCall)
			return false
		}
		return true
	})
	return bridgeCalls
}
