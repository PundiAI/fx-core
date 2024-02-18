package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

// GetParams returns the total set of erc20 parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetEnableErc20 returns the EnableErc20 parameter.
func (k Keeper) GetEnableErc20(ctx sdk.Context) bool {
	return k.GetParams(ctx).EnableErc20
}

// GetEnableEVMHook returns the EnableEVMHook parameter.
func (k Keeper) GetEnableEVMHook(ctx sdk.Context) bool {
	return k.GetParams(ctx).EnableEVMHook
}

// GetIbcTimeout returns the IbcTimeout parameter.
func (k Keeper) GetIbcTimeout(ctx sdk.Context) time.Duration {
	return k.GetParams(ctx).IbcTimeout
}
