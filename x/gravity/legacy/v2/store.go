package v2

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/gravity/types"
)

// MigrateStore performs in-place store migrations from v1 to v2. The
// migration includes:
//     Prune ibc sequence
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey) error {
	store := ctx.KVStore(storeKey)

	return pruneIbcSequenceKey(store)
}

// pruneIbcSequenceKey removes the zero balance addresses from balances store.
func pruneIbcSequenceKey(store sdk.KVStore) error {
	ibcSequenceStore := prefix.NewStore(store, types.KeyIbcSequenceHeight)
	iterator := ibcSequenceStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ibcSequenceStore.Delete(iterator.Key())
	}
	return nil
}
