package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/gov/types"
)

var _ types.QueryServer = Keeper{}

// Params return hub contract param
func (keeper Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := keeper.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
