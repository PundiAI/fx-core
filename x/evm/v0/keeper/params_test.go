package keeper_test

import (
	"github.com/functionx/fx-core/x/evm/v0/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.EvmKeeperV0.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	params.EvmDenom = "inj"
	suite.app.EvmKeeperV0.SetParams(suite.ctx, params)
	newParams := suite.app.EvmKeeperV0.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
