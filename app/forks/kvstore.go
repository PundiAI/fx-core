package forks

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/types"
)

func ClearEvmV0KVStores(ctx sdk.Context, keys map[string]*types.KVStoreKey) {
	logger := ctx.Logger()
	logger.Info("clear evm v0 kv store", "stores", fxtypes.ClearKVStores())

	multiStore := ctx.MultiStore()
	for _, storeName := range fxtypes.ClearKVStores() {
		if kvStoreKey, ok := keys[storeName]; ok {
			kvStore := multiStore.GetKVStore(kvStoreKey)
			if err := deleteKVStore(kvStore); err != nil {
				panic(fmt.Sprintf("failed to delete store %s: %s", storeName, err.Error()))
			}
		} else {
			panic(fmt.Sprintf("%s store not found", storeName))
		}
	}
}

func deleteKVStore(kv types.KVStore) error {
	// Note that we cannot write while iterating, so load all keys here, delete below
	var keys [][]byte
	itr := kv.Iterator(nil, nil)
	defer itr.Close()

	for itr.Valid() {
		keys = append(keys, itr.Key())
		itr.Next()
	}

	for _, k := range keys {
		kv.Delete(k)
	}
	return nil
}
