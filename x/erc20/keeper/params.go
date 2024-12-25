package keeper

import (
	"context"

	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) CheckEnableErc20(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	if !params.EnableErc20 {
		return types.ErrDisabled.Wrap("erc20 module is disabled")
	}
	return nil
}
