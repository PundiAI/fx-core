package v2

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschainv1 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v1"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey) {
	store := ctx.KVStore(storeKey)
	store.Delete(crosschainv1.OracleTotalDepositKey)

	// pruneIbcSequenceKey
	MigratePruneKey(store, crosschainv1.IbcSequenceHeightKey)
}

func MigratePruneKey(store sdk.KVStore, key []byte) {
	prefixStore := prefix.NewStore(store, key)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		prefixStore.Delete(iterator.Key())
	}
}
