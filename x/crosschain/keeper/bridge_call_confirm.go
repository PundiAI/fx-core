package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) SetSnapshotOracle(ctx sdk.Context, snapshotOracleKey *types.SnapshotOracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSnapshotOracleKey(snapshotOracleKey.OracleSetNonce), k.cdc.MustMarshal(snapshotOracleKey))
}

func (k Keeper) GetSnapshotOracle(ctx sdk.Context, oracleSetNonce uint64) (*types.SnapshotOracle, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSnapshotOracleKey(oracleSetNonce))
	if bz == nil {
		return nil, false
	}
	snapshotOracle := new(types.SnapshotOracle)
	k.cdc.MustUnmarshal(bz, snapshotOracle)
	return snapshotOracle, true
}

func (k Keeper) HasBridgeCallConfirm(ctx sdk.Context, nonce uint64, addr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBridgeCallConfirmKey(nonce, addr))
}

func (k Keeper) DeleteSnapshotOracle(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetSnapshotOracleKey(nonce))
}

func (k Keeper) RemoveEventSnapshotOracle(ctx sdk.Context, oracleSetNonce, nonce uint64) {
	snapshotOracle, found := k.GetSnapshotOracle(ctx, oracleSetNonce)
	if !found {
		return
	}
	snapshotOracle.RemoveNonce(nonce)
	if len(snapshotOracle.Nonces) == 0 {
		k.DeleteSnapshotOracle(ctx, oracleSetNonce)
	} else {
		k.SetSnapshotOracle(ctx, snapshotOracle)
	}
}

func (k Keeper) GetBridgeCallConfirm(ctx sdk.Context, nonce uint64, addr sdk.AccAddress) (*types.MsgBridgeCallConfirm, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeCallConfirmKey(nonce, addr))
	if bz == nil {
		return nil, false
	}
	var msg types.MsgBridgeCallConfirm
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetBridgeCallConfirm(ctx sdk.Context, addr sdk.AccAddress, msg *types.MsgBridgeCallConfirm) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeCallConfirmKey(msg.Nonce, addr), k.cdc.MustMarshal(msg))
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
	iterator := sdk.KVStorePrefixIterator(store, types.GetBridgeCallConfirmKeyByNonce(nonce))
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
	k.IterateBridgeCallByNonce(ctx, nonce, func(bridgeCall *types.OutgoingBridgeCall) bool {
		if bridgeCall.BlockHeight <= height {
			bridgeCalls = append(bridgeCalls, *bridgeCall)
			return false
		}
		return true
	})
	return bridgeCalls
}
