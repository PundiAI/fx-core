package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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

// TallyResult queries the tally of a proposal vote
func (keeper Keeper) TallyResult(ctx context.Context, req *v1.QueryTallyResultRequest) (*v1.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	proposal, err := keeper.Proposals.Get(ctx, req.ProposalId)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "proposal %d doesn't exist", req.ProposalId)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	var tallyResult v1.TallyResult

	switch {
	case proposal.Status == v1.StatusDepositPeriod:
		tallyResult = v1.EmptyTallyResult()

	case proposal.Status == v1.StatusPassed || proposal.Status == v1.StatusRejected || proposal.Status == v1.StatusFailed:
		tallyResult = *proposal.FinalTallyResult

	default:
		// proposal is in voting period
		var err error
		_, _, tallyResult, err = keeper.Tally(ctx, proposal)
		if err != nil {
			return nil, err
		}
	}

	return &v1.QueryTallyResultResponse{Tally: &tallyResult}, nil
}
