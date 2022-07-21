package crosschain_test

import (
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v2/x/crosschain"
	"github.com/functionx/fx-core/v2/x/crosschain/types"
)

func (suite *IntegrationTestSuite) TestUpdateCrossChainOraclesProposal() {
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

	err := crosschain.HandleUpdateChainOraclesProposal(suite.ctx, suite.MsgServer(), updateOracle)
	require.NoError(suite.T(), err)
	require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, suite.oracles[0].String()))
	require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, suite.oracles[1].String()))
	require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, suite.oracles[2].String()))
}
