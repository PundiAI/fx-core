package keeper_test

func (suite *KeeperTestSuite) TestOracleAndBridger() {
	for _, oracle := range suite.oracleAddrs {
		suite.Require().True(suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	for _, bridger := range suite.bridgerAddrs {
		oracle, found := suite.Keeper().GetOracleAddrByBridgerAddr(suite.ctx, bridger)
		suite.Require().False(found)
		suite.Require().Equal("", oracle.String())
	}
}
