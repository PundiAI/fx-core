package keeper_test

import (
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestUpdateCrossChainOraclesProposal() {
	updateOracle := &types.UpdateChainOraclesProposal{
		Title:       "Test UpdateCrossChainOracles",
		Description: "test",
		Oracles: []string{
			suite.oracles[0].String(),
			suite.oracles[1].String(),
			suite.oracles[2].String(),
		},
		ChainName: suite.chainName,
	}

	err := suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.NoError(suite.T(), err)
	require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, suite.oracles[0].String()))
	require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, suite.oracles[1].String()))
	require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, suite.oracles[2].String()))
}
