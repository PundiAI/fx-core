package keeper_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestUpdateCrossChainOraclesProposal() {
	updateOracle := &types.UpdateChainOraclesProposal{
		Title:       "Test UpdateCrossChainOracles",
		Description: "test",
		Oracles:     []string{},
		ChainName:   suite.chainName,
	}
	for _, oracle := range suite.oracles {
		updateOracle.Oracles = append(updateOracle.Oracles, oracle.String())
	}

	err := suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.NoError(suite.T(), err)
	for _, oracle := range suite.oracles {
		require.True(suite.T(), suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	updateOracle.Oracles = []string{}
	number := rand.Intn(100)
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, sdk.AccAddress(helpers.GenerateAddress().Bytes()).String())
	}
	err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.NoError(suite.T(), err)

	updateOracle.Oracles = []string{}
	number = rand.Intn(2) + 101
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, sdk.AccAddress(helpers.GenerateAddress().Bytes()).String())
	}
	err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	require.Error(suite.T(), err)
}
