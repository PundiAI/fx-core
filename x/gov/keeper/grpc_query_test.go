package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	govtypes "github.com/functionx/fx-core/v7/x/gov/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient
	response, err := queryClient.Params(suite.ctx, &govtypes.QueryParamsRequest{MsgType: sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{})})
	suite.Require().NoError(err)
	params := response.GetParams()
	suite.Require().EqualValues(params.Quorum, govtypes.DefaultErc20Quorum.String())

	params.Quorum = "0.3"
	_, err = suite.MsgServer.UpdateParams(suite.ctx, &govtypes.MsgUpdateParams{Authority: suite.govAcct, Params: params})
	suite.Require().NoError(err)
	response, err = queryClient.Params(suite.ctx, &govtypes.QueryParamsRequest{MsgType: sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{})})
	suite.Require().NoError(err)
	params = response.GetParams()
	suite.Require().NotEqualValues(params.Quorum, govtypes.DefaultErc20Quorum.String())
	suite.Require().EqualValues(params.Quorum, "0.3")

	response, err = queryClient.Params(suite.ctx, &govtypes.QueryParamsRequest{MsgType: sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{})})
	suite.Require().NoError(err)
	params = response.GetParams()
	suite.Require().EqualValues(params.Quorum, govtypes.DefaultErc20Quorum.String())
	suite.Require().NotEqualValues(params.Quorum, "0.3")
}

func (suite *KeeperTestSuite) TestGRPCQueryEGFParams() {
	queryClient := suite.queryClient
	response, err := queryClient.EGFParams(suite.ctx, &govtypes.QueryEGFParamsRequest{})
	suite.Require().NoError(err)
	params := response.GetParams()
	suite.Require().EqualValues(params.EgfDepositThreshold.String(), sdk.NewCoin(fxtypes.DefaultDenom, govtypes.DefaultEgfDepositThreshold).String())
	suite.Require().EqualValues(params.ClaimRatio, govtypes.DefaultClaimRatio.String())

	params.ClaimRatio = "0.4"
	params.EgfDepositThreshold = sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100000).MulRaw(1e18))
	_, err = suite.MsgServer.UpdateEGFParams(suite.ctx, &govtypes.MsgUpdateEGFParams{Authority: suite.govAcct, Params: params})
	suite.Require().NoError(err)
	response, err = queryClient.EGFParams(suite.ctx, &govtypes.QueryEGFParamsRequest{})
	suite.Require().NoError(err)
	params = response.GetParams()
	suite.Require().EqualValues(params.ClaimRatio, "0.4")
	suite.Require().EqualValues(params.EgfDepositThreshold.String(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100000).MulRaw(1e18)).String())
}
