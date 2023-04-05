package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/gov/types"
)

// FxGovParamsKey FxParamsKey is the key to query all gov params
var FxGovParamsKey = []byte("FxGovParam")

// GetParams gets the gov module's parameters.
func (keeper Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(FxGovParamsKey)
	if bz == nil {
		return params
	}
	keeper.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the gov module's parameters.
func (keeper Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(FxGovParamsKey, bz)

	return nil
}
