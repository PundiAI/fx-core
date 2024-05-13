package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	expParams := types.DefaultParams()

	res, err := suite.queryClient.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expParams, res.Params)
}

func (suite *KeeperTestSuite) TestTokenPairs() {
	testCases := []struct {
		name     string
		malleate func() (*types.QueryTokenPairsRequest, *types.QueryTokenPairsResponse)
		expPass  bool
	}{
		{
			"metadata pairs registered",
			func() (*types.QueryTokenPairsRequest, *types.QueryTokenPairsResponse) {
				fxPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)
				return &types.QueryTokenPairsRequest{}, &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: 1},
					TokenPairs: []types.TokenPair{fxPair},
				}
			},
			true,
		},
		{
			"metadata +1 pair registered w/pagination",
			func() (*types.QueryTokenPairsRequest, *types.QueryTokenPairsResponse) {
				req := &types.QueryTokenPairsRequest{
					Pagination: &query.PageRequest{Limit: 10, CountTotal: true},
				}
				pair := types.NewTokenPair(helpers.GenHexAddress(), "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)

				fxPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)
				expRes := &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: 2},
					TokenPairs: []types.TokenPair{fxPair, pair},
				}
				return req, expRes
			},
			true,
		},
		{
			"metadata +2 pairs registered wo/pagination",
			func() (*types.QueryTokenPairsRequest, *types.QueryTokenPairsResponse) {
				req := &types.QueryTokenPairsRequest{}
				pair := types.NewTokenPair(helpers.GenHexAddress(), "coin", true, types.OWNER_MODULE)
				pair2 := types.NewTokenPair(helpers.GenHexAddress(), "coin2", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair)
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, pair2)

				fxPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
				suite.Require().True(found)
				expRes := &types.QueryTokenPairsResponse{
					Pagination: &query.PageResponse{Total: 3},
					TokenPairs: []types.TokenPair{fxPair, pair, pair2},
				}
				return req, expRes
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			req, expRes := tc.malleate()

			res, err := suite.queryClient.TokenPairs(sdk.WrapSDKContext(suite.ctx), req)
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

func (suite *KeeperTestSuite) TestTokenPair() {
	testCases := []struct {
		name     string
		malleate func() (*types.QueryTokenPairRequest, *types.QueryTokenPairResponse)
		expPass  bool
	}{
		{
			"invalid token address",
			func() (*types.QueryTokenPairRequest, *types.QueryTokenPairResponse) {
				return &types.QueryTokenPairRequest{}, &types.QueryTokenPairResponse{}
			},
			false,
		},
		{
			"token pair not found",
			func() (*types.QueryTokenPairRequest, *types.QueryTokenPairResponse) {
				req := &types.QueryTokenPairRequest{
					Token: helpers.GenHexAddress().Hex(),
				}
				expRes := &types.QueryTokenPairResponse{}
				return req, expRes
			},
			false,
		},
		{
			"token pair found",
			func() (*types.QueryTokenPairRequest, *types.QueryTokenPairResponse) {
				addr := helpers.GenHexAddress()
				pair := types.NewTokenPair(addr, "coin", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

				req := &types.QueryTokenPairRequest{
					Token: pair.Erc20Address,
				}
				expRes := &types.QueryTokenPairResponse{TokenPair: pair}
				return req, expRes
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			req, expRes := tc.malleate()

			res, err := suite.queryClient.TokenPair(sdk.WrapSDKContext(suite.ctx), req)
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
	testCases := []struct {
		name     string
		malleate func() (*types.QueryDenomAliasesRequest, *types.QueryDenomAliasesResponse)
		expPass  bool
	}{
		{
			"invalid format for denom",
			func() (*types.QueryDenomAliasesRequest, *types.QueryDenomAliasesResponse) {
				req := &types.QueryDenomAliasesRequest{}
				expRes := &types.QueryDenomAliasesResponse{}
				return req, expRes
			},
			false,
		},
		{
			"not registered with denom",
			func() (*types.QueryDenomAliasesRequest, *types.QueryDenomAliasesResponse) {
				req := &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes := &types.QueryDenomAliasesResponse{}
				return req, expRes
			},
			false,
		},
		{
			"metadata not found",
			func() (*types.QueryDenomAliasesRequest, *types.QueryDenomAliasesResponse) {
				req := &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes := &types.QueryDenomAliasesResponse{}
				pair := types.NewTokenPair(helpers.GenHexAddress(), "usdt", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)
				return req, expRes
			},
			false,
		},
		{
			"metadata not support many to one",
			func() (*types.QueryDenomAliasesRequest, *types.QueryDenomAliasesResponse) {
				req := &types.QueryDenomAliasesRequest{Denom: "usdt"}
				expRes := &types.QueryDenomAliasesResponse{}

				pair := types.NewTokenPair(helpers.GenHexAddress(), "usdt", true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

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
				return req, expRes
			},
			true,
		},
		{
			"ok",
			func() (*types.QueryDenomAliasesRequest, *types.QueryDenomAliasesResponse) {
				metadata := newMetadata()

				req := &types.QueryDenomAliasesRequest{Denom: metadata.Base}
				expRes := &types.QueryDenomAliasesResponse{Aliases: metadata.DenomUnits[0].Aliases}

				pair := types.NewTokenPair(helpers.GenHexAddress(), metadata.Base, true, types.OWNER_MODULE)
				suite.app.Erc20Keeper.AddTokenPair(suite.ctx, pair)

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, metadata)
				return req, expRes
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			req, expRes := tc.malleate()

			res, err := suite.queryClient.DenomAliases(sdk.WrapSDKContext(suite.ctx), req)
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
	testCases := []struct {
		name     string
		malleate func() (*types.QueryAliasDenomRequest, *types.QueryAliasDenomResponse)
		expPass  bool
	}{
		{
			"invalid format for alias",
			func() (*types.QueryAliasDenomRequest, *types.QueryAliasDenomResponse) {
				req := &types.QueryAliasDenomRequest{}
				expRes := &types.QueryAliasDenomResponse{}
				return req, expRes
			},
			false,
		},
		{
			"ok without denom alias",
			func() (*types.QueryAliasDenomRequest, *types.QueryAliasDenomResponse) {
				denom := fmt.Sprintf("test%s", helpers.GenHexAddress().String())
				req := &types.QueryAliasDenomRequest{Alias: denom}
				expRes := &types.QueryAliasDenomResponse{}
				return req, expRes
			},
			false,
		},
		{
			"ok",
			func() (*types.QueryAliasDenomRequest, *types.QueryAliasDenomResponse) {
				denom := fmt.Sprintf("test%s", helpers.GenHexAddress().String())
				req := &types.QueryAliasDenomRequest{Alias: denom}
				expRes := &types.QueryAliasDenomResponse{Denom: "usdt"}

				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "usdt", denom)
				return req, expRes
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			req, expRes := tc.malleate()

			res, err := suite.queryClient.AliasDenom(sdk.WrapSDKContext(suite.ctx), req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
