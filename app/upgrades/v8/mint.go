package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func migrateMintParams(ctx sdk.Context, keeper mintkeeper.Keeper) error {
	params, err := keeper.Params.Get(ctx)
	if err != nil {
		return err
	}
	params.MintDenom = fxtypes.DefaultDenom
	return keeper.Params.Set(ctx, params)
}
