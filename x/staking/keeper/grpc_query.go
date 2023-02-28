package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v3/x/staking/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	*Keeper
}

var _ types.QueryServer = Querier{}

func (k Querier) ValidatorLPToken(c context.Context, req *types.QueryValidatorLPTokenRequest) (*types.QueryValidatorLPTokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.ValidatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	lpTokenContract, found := k.GetValidatorLPToken(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "lp token %s not found", req.ValidatorAddr)
	}

	val, found := k.GetValidator(ctx, valAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "lp token %s not found", req.ValidatorAddr)
	}

	lpToken := types.LPToken{
		ValidatorAddr: req.ValidatorAddr,
		Address:       lpTokenContract.String(),
		Name:          req.ValidatorAddr,
		Symbol:        types.LPTokenSymbol,
		Decimal:       uint32(types.LPTokenDecimals),
		TotalSupply:   sdkmath.NewIntFromBigInt(val.DelegatorShares.BigInt()),
	}

	return &types.QueryValidatorLPTokenResponse{LpToken: &lpToken}, nil
}
