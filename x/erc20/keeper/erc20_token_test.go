package keeper_test

func (suite *KeeperTestSuite) TestKeeper_ToggleTokenConvert() {
	erc20Token, err := suite.GetKeeper().ToggleTokenConvert(suite.Ctx, "test")
	suite.Require().EqualError(err, "collections: not found: key 'test' of type github.com/cosmos/gogoproto/fx.erc20.v1.ERC20Token")
	suite.Require().Empty(erc20Token)
}
