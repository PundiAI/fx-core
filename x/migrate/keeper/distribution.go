package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func MigrateDistributionFunc(key sdk.StoreKey) MigrateHandler {
	return func(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error {
		store := ctx.KVStore(key)

		iter := sdk.KVStorePrefixIterator(store, distributiontypes.DelegatorStartingInfoPrefix)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			val, del := distributiontypes.GetDelegatorStartingInfoAddresses(iter.Key())
			if del.Equals(from) {
				store.Delete(distributiontypes.GetDelegatorStartingInfoKey(val, del))
				store.Set(distributiontypes.GetDelegatorStartingInfoKey(val, to), iter.Value())
			}
		}
		return nil
	}
}
