package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v8/x/gov/types"
)

var _ types.QueryServer = Keeper{}

func (keeper Keeper) SwitchParams(ctx context.Context, _ *types.QuerySwitchParamsRequest) (*types.QuerySwitchParamsResponse, error) {
	params := keeper.GetSwitchParams(sdk.UnwrapSDKContext(ctx))
	return &types.QuerySwitchParamsResponse{Params: params}, nil
}

func (keeper Keeper) CustomParams(ctx context.Context, req *types.QueryCustomParamsRequest) (*types.QueryCustomParamsResponse, error) {
	if req == nil || req.MsgUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	customParams, err := keeper.CustomerParams.Get(ctx, req.MsgUrl)
	if err == nil {
		return &types.QueryCustomParamsResponse{Params: customParams}, nil
	}

	if !errors.IsOf(err, collections.ErrNotFound) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	params, err := keeper.Params.Get(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCustomParamsResponse{Params: types.CustomParams{
		DepositRatio: sdkmath.LegacyZeroDec().String(),
		VotingPeriod: params.VotingPeriod,
		Quorum:       params.Quorum,
	}}, nil
}
