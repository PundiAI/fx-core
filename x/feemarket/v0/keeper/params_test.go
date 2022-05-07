package keeper_test

import (
	"github.com/functionx/fx-core/x/feemarket/v0/types"
)

func (suite *KeeperTestSuite) TestSetGetParams() {
	params := suite.app.FeeMarketKeeperV0.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	params.ElasticityMultiplier = 3
	suite.app.FeeMarketKeeperV0.SetParams(suite.ctx, params)
	newParams := suite.app.FeeMarketKeeperV0.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
