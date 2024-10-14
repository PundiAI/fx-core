package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) SavePendingExecuteClaim(ctx sdk.Context, claim types.ExternalClaim) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalInterface(claim)
	if err != nil {
		return err
	}
	store.Set(types.GetPendingExecuteClaimKey(claim.GetEventNonce()), bz)
	return err
}

func (k Keeper) GetPendingExecuteClaim(ctx sdk.Context, eventNonce uint64) (types.ExternalClaim, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPendingExecuteClaimKey(eventNonce))
	if len(bz) == 0 {
		return nil, fmt.Errorf("claim %d not found", eventNonce)
	}
	var claim types.ExternalClaim
	if err := k.cdc.UnmarshalInterface(bz, &claim); err != nil {
		return nil, err
	}
	return claim, nil
}

func (k Keeper) DeletePendingExecuteClaim(ctx sdk.Context, eventNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPendingExecuteClaimKey(eventNonce))
}
