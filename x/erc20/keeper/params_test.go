package keeper_test

import (
	"time"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.Erc20Keeper.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)

	enableErc20 := suite.app.Erc20Keeper.GetEnableErc20(suite.ctx)
	suite.Require().Equal(params.EnableErc20, enableErc20)

	enableEVMHook := suite.app.Erc20Keeper.GetEnableEVMHook(suite.ctx)
	suite.Require().Equal(params.EnableEVMHook, enableEVMHook)

	ibcTimeout := suite.app.Erc20Keeper.GetIbcTimeout(suite.ctx)
	suite.Require().Equal(params.IbcTimeout, ibcTimeout)

	params.EnableErc20 = false
	params.EnableEVMHook = false
	params.IbcTimeout = 24 * time.Hour
	err := suite.app.Erc20Keeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	newParams := suite.app.Erc20Keeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)

	enableErc20 = suite.app.Erc20Keeper.GetEnableErc20(suite.ctx)
	suite.Require().Equal(newParams.EnableErc20, enableErc20)

	enableEVMHook = suite.app.Erc20Keeper.GetEnableEVMHook(suite.ctx)
	suite.Require().Equal(newParams.EnableEVMHook, enableEVMHook)

	ibcTimeout = suite.app.Erc20Keeper.GetIbcTimeout(suite.ctx)
	suite.Require().Equal(newParams.IbcTimeout, ibcTimeout)
}
