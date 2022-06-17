package v045

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey) {
	store := ctx.KVStore(storeKey)
	store.Delete(v042.OracleTotalDepositKey)

	// pruneIbcSequenceKey
	MigratePruneKey(store, v042.IbcSequenceHeightKey)
}

func MigratePruneKey(store sdk.KVStore, key []byte) {
	prefixStore := prefix.NewStore(store, key)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		prefixStore.Delete(iterator.Key())
	}
}
