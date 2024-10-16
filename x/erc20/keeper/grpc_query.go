package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	k Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k: k}
}

// TokenPairs return registered pairs
func (s queryServer) TokenPairs(c context.Context, req *types.QueryTokenPairsRequest) (*types.QueryTokenPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	erc20tokens, pageRes, err := query.CollectionPaginate(ctx, s.k.ERC20Token, req.Pagination,
		func(_ string, value types.ERC20Token) (types.ERC20Token, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
	}

	return &types.QueryTokenPairsResponse{
		Erc20Tokens: erc20tokens,
		Pagination:  pageRes,
	}, nil
}

// TokenPair returns a given registered token pair
func (s queryServer) TokenPair(c context.Context, req *types.QueryTokenPairRequest) (*types.QueryTokenPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// check if the token is a hex address, if not, check if it is a valid SDK denom
	if err := contract.ValidateEthereumAddress(req.Token); err != nil {
		if err = sdk.ValidateDenom(req.Token); err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"invalid format for token %s, should be either hex ('0x...') cosmos denom", req.Token,
			)
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	baseDenom, err := s.k.DenomIndex.Get(ctx, req.Token)
	if err != nil {
		baseDenom = req.Token
	}
	erc20Token, err := s.k.ERC20Token.Get(ctx, baseDenom)
	if err != nil {
		return nil, err
	}
	return &types.QueryTokenPairResponse{Erc20Token: erc20Token}, nil
}

// Params return erc20 module param
func (s queryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params, err := s.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}
