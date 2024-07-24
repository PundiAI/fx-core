package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) SavePendingExecuteClaim(ctx sdk.Context, claim types.ExternalClaim) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalInterface(claim)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetPendingExecuteClaimKey(claim.GetEventNonce()), bz)
}

func (k Keeper) GetPendingExecuteClaim(ctx sdk.Context, eventNonce uint64) (types.ExternalClaim, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingExecuteClaimKey(eventNonce))
	if len(bz) == 0 {
		return nil, false
	}
	var claim types.ExternalClaim
	if err := k.cdc.UnmarshalInterface(bz, &claim); err != nil {
		panic(err)
	}
	return claim, true
}

func (k Keeper) DeletePendingExecuteClaim(ctx sdk.Context, eventNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPendingExecuteClaimKey(eventNonce))
}
