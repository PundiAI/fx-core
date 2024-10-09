package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/gov/types"
)

func (keeper Keeper) GetSwitchParams(ctx sdk.Context) (params types.SwitchParams) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.FxSwitchParamsKey)
	if bz == nil {
		return params
	}
	keeper.cdc.MustUnmarshal(bz, &params)
	return params
}

func (keeper Keeper) GetDisabledMsgs(ctx sdk.Context) []string {
	return keeper.GetSwitchParams(ctx).DisableMsgTypes
}

func (keeper Keeper) SetSwitchParams(ctx sdk.Context, params *types.SwitchParams) error {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(types.FxSwitchParamsKey, bz)
	return nil
}

func (keeper Keeper) GetCustomParams(ctx context.Context, msgType string) (types.CustomParams, bool) {
	value, err := keeper.CustomerParams.Get(ctx, msgType)
	return value, err == nil
}
