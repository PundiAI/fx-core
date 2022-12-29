package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestTokenPairs() {
	var (
		req    *types.QueryTokenPairsRequest
		expRes *types.QueryTokenPairsResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"metadata pairs registered",
			func() {
				req = &types.QueryTokenPairsRequest{}
				expRes = &types.QueryTokenPairsResponse{Pagination: &query.PageResponse{}}

				tokenPairs := getMetadataTokenPairs()
				expRes = &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: uint64(len(tokenPairs))},
					TokenPairs: tokenPairs,
				}
			},
			true,
		},
		{
			"metadata +1 pair registered w/pagination",
			func() {
				req = &types.QueryTokenPairsRequest{
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true},
				}
				pair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

				//clear erc20 address
				pairs := clearTokenPairErc20Address(pair)
				tokenPairs := getMetadataTokenPairs()
				expRes = &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: uint64(len(tokenPairs)) + 1},
					TokenPairs: append(pairs, tokenPairs...),
				}
			},
			true,
		},
		{
			"metadata +2 pairs registered wo/pagination",
			func() {
				req = &types.QueryTokenPairsRequest{}
				pair := types.NewTokenPair(helpers.GenerateAddress(), "coin", true, types.OWNER_MODULE)
				pair2 := types.NewTokenPair(helpers.GenerateAddress(), "coin2", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair2)

				//clear erc20 address
				pairs := clearTokenPairErc20Address(pair, pair2)
				tokenPairs := getMetadataTokenPairs()
				expRes = &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: uint64(len(tokenPairs)) + 2},
					TokenPairs: append(pairs, tokenPairs...),
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
				//clear erc20 address
				newPairs := clearTokenPairErc20Address(res.TokenPairs...)
				suite.Require().ElementsMatch(expRes.TokenPairs, newPairs)
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
					Token: helpers.GenerateAddress().Hex(),
				}
				expRes = &types.QueryTokenPairResponse{}
			},
			false,
		},
		{
			"token pair found",
			func() {
				addr := helpers.GenerateAddress()
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

func (suite *KeeperTestSuite) TestDenomAlias() {
	var (
		req    *types.QueryDenomAliasesRequest
		expRes *types.QueryDenomAliasesResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"invalid format for denom",
			func() {
				req = &types.QueryDenomAliasesRequest{}
				expRes = &types.QueryDenomAliasesResponse{}
			},
			false,
		},
		{
			"not registered with denom",
			func() {
				req = &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes = &types.QueryDenomAliasesResponse{}
			},
			false,
		},
		{
			"metadata not found",
			func() {
				req = &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes = &types.QueryDenomAliasesResponse{}

				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, "usdt", []byte{})
			},
			false,
		},
		{
			"metadata not support many to one",
			func() {
				req = &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes = &types.QueryDenomAliasesResponse{}

				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, "usdt", []byte{})

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, banktypes.Metadata{
					Description: "The cross chain token of the Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "usdt",
							Exponent: 0,
						},
						{
							Denom:    "USDT",
							Exponent: 18,
						},
					},
					Base:    "usdt",
					Display: "usdt",
					Name:    "Tether USD",
					Symbol:  "USDT",
				})
			},
			false,
		},
		{
			"ok",
			func() {
				req = &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes = &types.QueryDenomAliasesResponse{Aliases: []string{bscDenom, polygonDenom}}

				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, "usdt", []byte{})

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, banktypes.Metadata{
					Description: "The cross chain token of the Function X",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "usdt",
							Exponent: 0,
							Aliases:  []string{bscDenom, polygonDenom},
						},
						{
							Denom:    "USDT",
							Exponent: 18,
						},
					},
					Base:    "usdt",
					Display: "usdt",
					Name:    "Tether USD",
					Symbol:  "USDT",
				})
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ctx := sdk.WrapSDKContext(suite.ctx)

			tc.malleate()

			res, err := suite.queryClient.DenomAliases(ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAliasDenom() {
	var (
		req    *types.QueryAliasDenomRequest
		expRes *types.QueryAliasDenomResponse
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"invalid format for alias",
			func() {
				req = &types.QueryAliasDenomRequest{}
				expRes = &types.QueryAliasDenomResponse{}
			},
			false,
		},
		{
			"ok without denom alias",
			func() {
				req = &types.QueryAliasDenomRequest{Alias: bscDenom}
				expRes = &types.QueryAliasDenomResponse{}
			},
			false,
		},
		{
			"ok",
			func() {
				req = &types.QueryAliasDenomRequest{Alias: bscDenom}
				expRes = &types.QueryAliasDenomResponse{Denom: "usdt"}

				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "usdt", bscDenom)
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			ctx := sdk.WrapSDKContext(suite.ctx)

			tc.malleate()

			res, err := suite.queryClient.AliasDenom(ctx, req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
