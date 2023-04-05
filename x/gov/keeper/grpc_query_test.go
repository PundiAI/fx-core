package keeper_test

import (
	"github.com/functionx/fx-core/v3/x/gov/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryFxParams() {
	queryClient := suite.queryClient
	params := types.DefaultParams()
	_, err := suite.MsgServer.UpdateParams(suite.ctx, &types.MsgUpdateParams{Authority: suite.govAcct, Params: *params})
	suite.Require().NoError(err)
	response, err := queryClient.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().EqualValues(response.Params.MinInitialDeposit, params.MinInitialDeposit)
	suite.Require().EqualValues(response.Params.EgfDepositThreshold, params.EgfDepositThreshold)
	suite.Require().EqualValues(response.Params.ClaimRatio, params.ClaimRatio)
	suite.Require().EqualValues(response.Params.Erc20Quorum, params.Erc20Quorum)
	suite.Require().EqualValues(response.Params.EvmQuorum, params.EvmQuorum)
	suite.Require().EqualValues(response.Params.EgfVotingPeriod, params.EgfVotingPeriod)
	suite.Require().EqualValues(response.Params.EgfVotingPeriod, params.EgfVotingPeriod)
}
