package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

var _ types.QueryServer = Keeper{}

// TokenPairs return registered pairs
func (k Keeper) TokenPairs(c context.Context, req *types.QueryTokenPairsRequest) (*types.QueryTokenPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var pairs []types.TokenPair
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var pair types.TokenPair
		if err := k.cdc.Unmarshal(value, &pair); err != nil {
			return err
		}
		pairs = append(pairs, pair)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryTokenPairsResponse{
		TokenPairs: pairs,
		Pagination: pageRes,
	}, nil
}

// TokenPair returns a given registered token pair
func (k Keeper) TokenPair(c context.Context, req *types.QueryTokenPairRequest) (*types.QueryTokenPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// check if the token is a hex address, if not, check if it is a valid SDK denom
	if err := fxtypes.ValidateEthereumAddress(req.Token); err != nil {
		if err := sdk.ValidateDenom(req.Token); err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"invalid format for token %s, should be either hex ('0x...') cosmos denom", req.Token,
			)
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	pair, found := k.GetTokenPair(ctx, req.Token)
	if !found {
		return nil, status.Errorf(codes.NotFound, "token pair with token '%s'", req.Token)
	}

	return &types.QueryTokenPairResponse{TokenPair: pair}, nil
}

// Params return hub contract param
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// DenomAliases returns denom aliases
func (k Keeper) DenomAliases(c context.Context, req *types.QueryDenomAliasesRequest) (*types.QueryDenomAliasesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// check if it is a valid SDK denom
	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format for denom %s", req.Denom)
	}

	ctx := sdk.UnwrapSDKContext(c)
	if !k.IsDenomRegistered(ctx, req.Denom) {
		return nil, status.Errorf(codes.NotFound, "not registered with denom '%s'", req.Denom)
	}

	md, found := k.bankKeeper.GetDenomMetaData(ctx, req.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "metadata not found with denom '%s'", req.Denom)
	}

	if len(md.DenomUnits) == 0 {
		return nil, status.Errorf(codes.NotFound, "not found alias with denom '%s'", req.Denom)
	}

	return &types.QueryDenomAliasesResponse{Aliases: md.DenomUnits[0].Aliases}, nil
}

// AliasDenom returns alias denom
func (k Keeper) AliasDenom(c context.Context, req *types.QueryAliasDenomRequest) (*types.QueryAliasDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// check if it is a valid SDK denom
	if err := sdk.ValidateDenom(req.Alias); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid format for alias %s", req.Alias)
	}

	ctx := sdk.UnwrapSDKContext(c)
	denom, found := k.GetAliasDenom(ctx, req.Alias)
	if !found {
		return nil, status.Errorf(codes.NotFound, "denom not found with alias '%s'", req.Alias)
	}

	return &types.QueryAliasDenomResponse{Denom: denom}, nil
}
