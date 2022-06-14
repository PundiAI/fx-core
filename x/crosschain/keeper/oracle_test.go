package keeper_test

import (
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestOracleAndBridger() {

	for _, oracle := range suite.oracles {
		require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	for _, bridger := range suite.bridgers {
		oracle, found := suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
		require.False(suite.T(), found)
		require.Equal(suite.T(), "", oracle.String())
	}
}
