package keeper

import (
	"math"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) DeleteOutgoingBridgeCallRecord(ctx sdk.Context, bridgeCallNonce uint64) error {
	// 1. delete bridge call
	k.DeleteOutgoingBridgeCall(ctx, bridgeCallNonce)

	// 2. delete bridge call confirm
	k.DeleteBridgeCallConfirm(ctx, bridgeCallNonce)

	// 3. delete cache origin amount
	return k.erc20Keeper.DeleteCache(ctx, types.NewOriginTokenKey(k.moduleName, bridgeCallNonce))
}

func (k Keeper) SetOutgoingBridgeCall(ctx sdk.Context, outCall *types.OutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingBridgeCallNonceKey(outCall.Nonce), k.cdc.MustMarshal(outCall))
	// value is just a placeholder
	store.Set(
		types.GetOutgoingBridgeCallAddressAndNonceKey(outCall.Sender, outCall.Nonce),
		k.cdc.MustMarshal(&gogotypes.BoolValue{Value: true}),
	)
}

func (k Keeper) HasOutgoingBridgeCall(ctx sdk.Context, nonce uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingBridgeCallNonceKey(nonce))
}

func (k Keeper) HasOutgoingBridgeCallAddressAndNonce(ctx sdk.Context, sender string, nonce uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingBridgeCallAddressAndNonceKey(sender, nonce))
}

func (k Keeper) GetOutgoingBridgeCallByNonce(ctx sdk.Context, nonce uint64) (*types.OutgoingBridgeCall, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingBridgeCallNonceKey(nonce))
	if bz == nil {
		return nil, false
	}
	var outCall types.OutgoingBridgeCall
	k.cdc.MustUnmarshal(bz, &outCall)
	return &outCall, true
}

func (k Keeper) DeleteOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	outCall, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
	if !found {
		return
	}
	store.Delete(types.GetOutgoingBridgeCallNonceKey(nonce))
	store.Delete(types.GetOutgoingBridgeCallAddressAndNonceKey(outCall.Sender, outCall.Nonce))
}

func (k Keeper) IterateOutgoingBridgeCalls(ctx sdk.Context, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.OutgoingBridgeCallNonceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var outCall types.OutgoingBridgeCall
		k.cdc.MustUnmarshal(iterator.Value(), &outCall)
		if cb(&outCall) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallsByAddress(ctx sdk.Context, senderAddr string, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.GetOutgoingBridgeCallAddressKey(senderAddr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nonce := types.ParseOutgoingBridgeCallNonce(iterator.Key(), senderAddr)
		outCall, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
		if !found {
			continue
		}
		if cb(outCall) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallByNonce(ctx sdk.Context, startNonce uint64, cb func(outCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		outCall := new(types.OutgoingBridgeCall)
		k.cdc.MustUnmarshal(iter.Value(), outCall)
		if cb(outCall) {
			break
		}
	}
}
