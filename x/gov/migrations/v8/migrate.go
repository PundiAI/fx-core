package v8

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// Deprecated: do not use, remove in v8
	FxBaseParamsKeyPrefix = []byte("0x90")
	// Deprecated: do not use, remove in v8
	FxEGFParamsKey = []byte("0x91")
)

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{FxBaseParamsKeyPrefix, FxEGFParamsKey}
}

func DeleteOldParamsStore(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
) {
	store := ctx.KVStore(storeKey)
	removeKeys := GetRemovedStoreKeys()
	for _, key := range removeKeys {
		iterator := storetypes.KVStorePrefixIterator(store, key)
		for ; iterator.Valid(); iterator.Next() {
			store.Delete(iterator.Key())
		}
	}
}
