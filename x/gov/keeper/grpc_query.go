package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/gov/types"
)

var _ types.QueryServer = Keeper{}

func (keeper Keeper) Params(c context.Context, re *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := keeper.GetParams(ctx, re.MsgType)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (keeper Keeper) EGFParams(c context.Context, _ *types.QueryEGFParamsRequest) (*types.QueryEGFParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := keeper.GetEGFParams(ctx)
	return &types.QueryEGFParamsResponse{Params: params}, nil
}
