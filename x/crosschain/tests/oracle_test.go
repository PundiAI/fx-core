package tests_test

import (
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestOracleAndBridger() {
	for _, oracle := range suite.oracleAddrs {
		require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	for _, bridger := range suite.bridgerAddrs {
		oracle, found := suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
		require.False(suite.T(), found)
		require.Equal(suite.T(), "", oracle.String())
	}
}
