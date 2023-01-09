package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestKeeper_Tally() {
	proposalContent := erc20types.NewRegisterCoinProposal("register erc20", "foo",
		banktypes.Metadata{
			Display: "test",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "test", Exponent: 0},
				{Denom: "TEST", Exponent: 18},
			},
			Base:   "test",
			Name:   "test",
			Symbol: "TEST",
		})
	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, proposalContent)
	suite.NoError(err)
	suite.Equal(proposal.Status, types.StatusDepositPeriod)

	proposal.Status = types.StatusVotingPeriod
	suite.app.GovKeeper.SetProposal(suite.ctx, proposal)

	poolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.app.StakingKeeper.GetBondedPool(suite.ctx).GetAddress())
	suite.Equal(sdk.NewInt(int64(len(suite.valAddr)*100)).MulRaw(1e18).String(), poolBalances[0].Amount.String(), len(suite.valAddr))

	for _, voterAddr := range suite.valAddr {
		err := suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, voterAddr.Bytes(), types.NewNonSplitVoteOption(types.OptionYes))
		suite.NoError(err)
	}

	proposal, ok := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.ProposalId)
	suite.True(ok)

	passes, burnDeposits, tallyResults := suite.app.GovKeeper.Tally(suite.ctx, proposal)
	suite.True(passes)
	suite.False(burnDeposits)
	suite.Equal(sdk.ZeroInt(), tallyResults.Abstain)
	suite.Equal(sdk.ZeroInt(), tallyResults.No)
	suite.Equal(sdk.ZeroInt(), tallyResults.NoWithVeto)
	suite.Equal(poolBalances[0].Amount, tallyResults.Yes)
}
