package store

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RemoveStoreKeys(ctx sdk.Context, storeKey storetypes.StoreKey, prefixKeys [][]byte) {
	store := ctx.KVStore(storeKey)

	deleteFn := func(prefixKey []byte) {
		iterator := storetypes.KVStorePrefixIterator(store, prefixKey)
		defer iterator.Close()
		deleteCount := 0
		for ; iterator.Valid(); iterator.Next() {
			store.Delete(iterator.Key())
			deleteCount++
		}
		if deleteCount > 0 {
			ctx.Logger().Info("remove store key", "kvStore", storeKey.Name(),
				"prefix", prefixKey[0], "deleteKeyCount", deleteCount)
		}
	}

	for _, key := range prefixKeys {
		deleteFn(key)
	}
}
