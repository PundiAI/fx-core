package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/migrations/v3"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v8/x/gov/types"
)

var (
	_ types.QueryServer = QueryServer{}
	_ v1.QueryServer    = QueryServer{}
)

type QueryServer struct {
	v1.QueryServer
	k *Keeper
}

func NewQueryServer(k *Keeper) QueryServer {
	return QueryServer{
		QueryServer: keeper.NewQueryServer(k.Keeper),
		k:           k,
	}
}

func (q QueryServer) SwitchParams(ctx context.Context, _ *types.QuerySwitchParamsRequest) (*types.QuerySwitchParamsResponse, error) {
	params := q.k.GetSwitchParams(sdk.UnwrapSDKContext(ctx))
	return &types.QuerySwitchParamsResponse{Params: params}, nil
}

func (q QueryServer) CustomParams(ctx context.Context, req *types.QueryCustomParamsRequest) (*types.QueryCustomParamsResponse, error) {
	if req == nil || req.MsgUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	customParams, err := q.k.CustomerParams.Get(ctx, req.MsgUrl)
	if err == nil {
		return &types.QueryCustomParamsResponse{Params: customParams}, nil
	}

	if !errors.IsOf(err, collections.ErrNotFound) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	params, err := q.k.Params.Get(ctx)
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
func (q QueryServer) TallyResult(ctx context.Context, req *v1.QueryTallyResultRequest) (*v1.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	proposal, err := q.k.Proposals.Get(ctx, req.ProposalId)
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
		_, _, tallyResult, err = q.k.Tally(ctx, proposal)
		if err != nil {
			return nil, err
		}
	}

	return &v1.QueryTallyResultResponse{Tally: &tallyResult}, nil
}

type legacyQueryServer struct {
	v1beta1.QueryServer
	qs v1.QueryServer
}

func NewLegacyQueryServer(qs v1.QueryServer, k *Keeper) v1beta1.QueryServer {
	return &legacyQueryServer{
		QueryServer: keeper.NewLegacyQueryServer(k.Keeper),
		qs:          qs,
	}
}

func (q legacyQueryServer) TallyResult(ctx context.Context, req *v1beta1.QueryTallyResultRequest) (*v1beta1.QueryTallyResultResponse, error) {
	resp, err := q.qs.TallyResult(ctx, &v1.QueryTallyResultRequest{
		ProposalId: req.ProposalId,
	})
	if err != nil {
		return nil, err
	}

	tally, err := v3.ConvertToLegacyTallyResult(resp.Tally)
	if err != nil {
		return nil, err
	}

	return &v1beta1.QueryTallyResultResponse{Tally: tally}, nil
}
