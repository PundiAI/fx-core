package v045

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
)

// pruneIbcSequenceKey removes the zero balance addresses from balances store.
func MigratePruneIbcSequenceKey(ctx sdk.Context, storeKey sdk.StoreKey) {
	store := ctx.KVStore(storeKey)
	ibcSequenceStore := prefix.NewStore(store, v042.KeyIbcSequenceHeight)
	iterator := ibcSequenceStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ibcSequenceStore.Delete(iterator.Key())
	}
}
