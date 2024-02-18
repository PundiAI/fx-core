package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
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
	msgExecLegacyContent, err := govv1.NewLegacyContent(proposalContent, suite.govAcct)
	suite.NoError(err)
	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx.WithChainID(fxtypes.ChainIdWithEIP155()), []sdk.Msg{msgExecLegacyContent},
		fxgovtypes.NewFXMetadata(proposalContent.GetTitle(), proposalContent.GetDescription(), "").String())
	suite.NoError(err)
	suite.Equal(proposal.Status, govv1.StatusDepositPeriod)

	proposal.Status = govv1.StatusVotingPeriod
	suite.app.GovKeeper.SetProposal(suite.ctx, proposal)

	poolBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.app.StakingKeeper.GetBondedPool(suite.ctx).GetAddress())
	suite.Equal(sdkmath.NewInt(int64(len(suite.valAddr)*100)).MulRaw(1e18).String(), poolBalances[0].Amount.String(), len(suite.valAddr))

	for _, voterAddr := range suite.valAddr {
		err := suite.app.GovKeeper.AddVote(suite.ctx, proposal.Id, voterAddr.Bytes(), govv1.NewNonSplitVoteOption(govv1.OptionYes), "")
		suite.NoError(err)
	}

	proposal, ok := suite.app.GovKeeper.GetProposal(suite.ctx, proposal.Id)
	suite.True(ok)

	passes, burnDeposits, tallyResults := suite.app.GovKeeper.Tally(suite.ctx, proposal)
	suite.True(passes)
	suite.False(burnDeposits)
	suite.Equal(sdkmath.ZeroInt().String(), tallyResults.AbstainCount)
	suite.Equal(sdkmath.ZeroInt().String(), tallyResults.NoCount)
	suite.Equal(sdkmath.ZeroInt().String(), tallyResults.NoWithVetoCount)
	suite.Equal(poolBalances[0].Amount.String(), tallyResults.YesCount)
}
