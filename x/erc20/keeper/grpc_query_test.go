package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/functionx/fx-core/tests"

	"github.com/functionx/fx-core/x/erc20/types"
)

func (suite *KeeperTestSuite) TestTokenPairs() {
	var (
		req    *types.QueryTokenPairsRequest
		expRes *types.QueryTokenPairsResponse
	)

	var (
		fxTokenPair = types.TokenPair{
			Erc20Address:  "0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd",
			Denom:         "FX",
			Enabled:       true,
			ContractOwner: 1,
		}
		pundixTokenPair = types.TokenPair{
			Erc20Address:  "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75",
			Denom:         "eth0x338E7A8687AdA7274Dc87C95D94f920d8F4185AE",
			Enabled:       true,
			ContractOwner: 1,
		}
		purseTokenPair = types.TokenPair{
			Erc20Address:  "0x5FD55A1B9FC24967C4dB09C513C3BA0DFa7FF687",
			Denom:         "ibc/B1861D0C2E4BAFA42A61739291975B7663F278FFAF579F83C9C4AD3890D09CA0",
			Enabled:       true,
			ContractOwner: 1,
		}
		usdtTokenPair = types.TokenPair{
			Erc20Address:  "0xecEEEfCEE421D8062EF8d6b4D814efe4dc898265",
			Denom:         "eth0x1BE1f78d417B1C4A199bb8ad4c946Ca248f7A83e",
			Enabled:       true,
			ContractOwner: 1,
		}
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"4 pairs registered",
			func() {
				req = &types.QueryTokenPairsRequest{}
				expRes = &types.QueryTokenPairsResponse{Pagination: &query.PageResponse{}}
				expRes = &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: 4},
					TokenPairs: []types.TokenPair{fxTokenPair, pundixTokenPair, purseTokenPair, usdtTokenPair},
				}
			},
			true,
		},
		{
			"5 pair registered w/pagination",
			func() {
				req = &types.QueryTokenPairsRequest{
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true},
				}
				pair := types.NewTokenPair(tests.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

				expRes = &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: 5},
					TokenPairs: []types.TokenPair{pair, fxTokenPair, pundixTokenPair, purseTokenPair, usdtTokenPair},
				}
			},
			true,
		},
		{
			"6 pairs registered wo/pagination",
			func() {
				req = &types.QueryTokenPairsRequest{}
				pair := types.NewTokenPair(tests.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				pair2 := types.NewTokenPair(tests.GenerateAddress(), "coin2", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair2)

				expRes = &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: 6},
					TokenPairs: []types.TokenPair{pair, pair2, fxTokenPair, pundixTokenPair, purseTokenPair, usdtTokenPair},
				}
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ctx := sdk.WrapSDKContext(suite.ctx)
			tc.malleate()

			res, err := suite.queryClient.TokenPairs(ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.Pagination, res.Pagination)
				suite.Require().ElementsMatch(expRes.TokenPairs, res.TokenPairs)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	expParams := types.DefaultParams()

	res, err := suite.queryClient.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expParams, res.Params)
}

func (suite *KeeperTestSuite) TestTokenPair() {
	var (
		req    *types.QueryTokenPairRequest
		expRes *types.QueryTokenPairResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"invalid token address",
			func() {
				req = &types.QueryTokenPairRequest{}
				expRes = &types.QueryTokenPairResponse{}
			},
			false,
		},
		{
			"token pair not found",
			func() {
				req = &types.QueryTokenPairRequest{
					Token: tests.GenerateAddress().Hex(),
				}
				expRes = &types.QueryTokenPairResponse{}
			},
			false,
		},
		{
			"token pair found",
			func() {
				addr := tests.GenerateAddress()
				pair := types.NewTokenPair(addr, "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
				suite.app.Erc20Keeper.SetERC20Map(suite.ctx, addr, pair.GetID())
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())

				req = &types.QueryTokenPairRequest{
					Token: pair.Erc20Address,
				}
				expRes = &types.QueryTokenPairResponse{TokenPair: pair}
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ctx := sdk.WrapSDKContext(suite.ctx)
			tc.malleate()

			res, err := suite.queryClient.TokenPair(ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
