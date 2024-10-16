package keeper

import (
	"context"
)

func (k Keeper) GetEnableErc20(ctx context.Context) (bool, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}
	return params.EnableErc20, nil
}
