package tests

import (
	"fmt"
	"testing"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

type ERC20ProposalTestSuite struct {
	CrosschainERC20TestSuite
}

func TestERC20ProposalTestSuite(t *testing.T) {
	testSuite := NewTestSuite()
	erc20ProposalTestSuite := &ERC20ProposalTestSuite{
		CrosschainERC20TestSuite: NewCrosschainERC20TestSuite(testSuite),
	}
	suite.Run(t, erc20ProposalTestSuite)
}

func (suite *ERC20ProposalTestSuite) TestERC20Proposal() {
	suite.InitCrossChain()
	suite.InitRegisterCoinUSDT()

	// add eth usd denom
	ethUSDDenom := fmt.Sprintf("eth%s", ethUSDToken)
	proposalId := suite.ERC20.UpdateDenomAliasProposal("usdt", ethUSDDenom)
	suite.NoError(suite.network.WaitForNextBlock())
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	//check add
	aliasesResp, err := suite.ERC20.GRPCClient().ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: "usdt"})
	suite.NoError(err)
	suite.Equal(3, len(aliasesResp.Aliases))
	denomResp, err := suite.ERC20.GRPCClient().ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: ethUSDDenom})
	suite.NoError(err)
	suite.Equal("usdt", denomResp.Denom)

	//remove eth usd
	proposalId = suite.ERC20.UpdateDenomAliasProposal("usdt", ethUSDDenom)
	suite.NoError(suite.network.WaitForNextBlock())
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	//check remove
	aliasesResp, err = suite.ERC20.GRPCClient().ERC20Query().DenomAliases(suite.ctx, &erc20types.QueryDenomAliasesRequest{Denom: "usdt"})
	suite.NoError(err)
	suite.Equal(2, len(aliasesResp.Aliases))
	_, err = suite.ERC20.GRPCClient().ERC20Query().AliasDenom(suite.ctx, &erc20types.QueryAliasDenomRequest{Alias: ethUSDDenom})
	suite.Error(err)

	proposalId = suite.ERC20.ToggleTokenConversionProposal("usdt")
	suite.NoError(suite.network.WaitForNextBlock())
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	tokenPairResp, err := suite.ERC20.GRPCClient().ERC20Query().TokenPair(suite.ctx, &erc20types.QueryTokenPairRequest{Token: "usdt"})
	suite.NoError(err)
	suite.False(tokenPairResp.TokenPair.Enabled)
}
