package tests_test

import (
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestUpdateCrossChainOraclesProposal() {
	updateOracle := &types.UpdateChainOraclesProposal{ // nolint:staticcheck
		Title:       "Test UpdateCrossChainOracles",
		Description: "test",
		Oracles:     []string{},
		ChainName:   suite.chainName,
	}
	for _, oracle := range suite.oracleAddrs {
		updateOracle.Oracles = append(updateOracle.Oracles, oracle.String())
	}

	err := suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.NoError(suite.T(), err)
	for _, oracle := range suite.oracleAddrs {
		require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	updateOracle.Oracles = []string{}
	number := tmrand.Intn(100)
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.NoError(suite.T(), err)

	updateOracle.Oracles = []string{}
	number = tmrand.Intn(2) + 101
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.Error(suite.T(), err)
}
