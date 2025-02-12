package keeper_test

import sdkmath "cosmossdk.io/math"

func (suite *KeeperTestSuite) TestKeeper_SetCache() {
	suite.Require().NoError(suite.GetKeeper().SetCache(suite.Ctx, "key", sdkmath.NewInt(1)))
	suite.Require().Error(suite.GetKeeper().SetCache(suite.Ctx, "key", sdkmath.NewInt(1)))
}
