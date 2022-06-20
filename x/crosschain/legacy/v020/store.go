package v020

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v010 "github.com/functionx/fx-core/x/crosschain/legacy/v010"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey) {
	store := ctx.KVStore(storeKey)
	store.Delete(v010.OracleTotalDepositKey)

	// pruneIbcSequenceKey
	MigratePruneKey(store, v010.IbcSequenceHeightKey)
}

func MigratePruneKey(store sdk.KVStore, key []byte) {
	prefixStore := prefix.NewStore(store, key)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		prefixStore.Delete(iterator.Key())
	}
}
