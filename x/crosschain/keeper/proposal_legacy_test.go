package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
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
	suite.Require().NoError(err)
	for _, oracle := range suite.oracleAddrs {
		suite.Require().True(suite.Keeper().IsProposalOracle(suite.ctx, oracle.String()))
	}

	updateOracle.Oracles = []string{}
	number := tmrand.Intn(100)
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	suite.Require().NoError(err)

	updateOracle.Oracles = []string{}
	number = tmrand.Intn(2) + 101
	for i := 0; i < number; i++ {
		updateOracle.Oracles = append(updateOracle.Oracles, helpers.GenAccAddress().String())
	}
	err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, updateOracle)
	suite.Require().Error(err)
}
