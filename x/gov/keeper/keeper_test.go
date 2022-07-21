package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v2/app"
	"github.com/functionx/fx-core/v2/app/helpers"
	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/x/gov/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	app   *app.App
	ctx   sdk.Context
	addrs []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) MsgServer() types.MsgServer {
	return keeper.NewMsgServerImpl(govkeeper.NewMsgServerImpl(suite.app.GovKeeper.Keeper), suite.app.GovKeeper)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = helpers.Setup(false, false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.addrs = helpers.AddTestAddrs(suite.app, suite.ctx, 3, sdk.NewInt(5000*1e8).MulRaw(1e18))
	// init communityPool
	err := suite.app.DistrKeeper.FundCommunityPool(suite.ctx, sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(5000 * 1e5).MulRaw(1e18)}}, suite.addrs[0])
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestNormalDeposits() {
	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	TestProposal := types.NewTextProposal("Test", "description")
	testProposalMsg, err := types.NewMsgSubmitProposal(
		TestProposal,
		initCoins,
		suite.addrs[0],
	)
	suite.Require().NoError(err)
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
	proposalResponse, err := suite.MsgServer().SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
	suite.Require().NoError(err)
	_, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposalResponse.ProposalId, suite.addrs[1])
	suite.Require().False(found)
	proposal, ok := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusDepositPeriod, proposal.Status)
	suite.Require().True(initCoins.IsAllLT(minDeposit))

	// first deposit
	firstCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	votingStarted, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1], firstCoins)
	suite.Require().NoError(err)
	suite.Require().False(votingStarted)
	deposit, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1])
	suite.Require().True(found)
	suite.Require().Equal(firstCoins, deposit.Amount)
	suite.Require().Equal(suite.addrs[1].String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusDepositPeriod, proposal.Status)
	suite.Require().Equal(firstCoins, proposal.TotalDeposit.Sub(initCoins))
	suite.Require().True(initCoins.Add(firstCoins...).IsAllLT(minDeposit))

	// second deposit
	secondCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(9 * 1e3).MulRaw(1e18)}}
	votingStarted, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1], secondCoins)
	suite.Require().NoError(err)
	suite.Require().True(votingStarted)
	deposit, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1])
	suite.Require().True(found)
	suite.Require().Equal(firstCoins.Add(secondCoins...), deposit.Amount)
	suite.Require().Equal(suite.addrs[1].String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusVotingPeriod, proposal.Status)
	suite.Require().True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllGTE(minDeposit))

}

func (suite *KeeperTestSuite) TestEGFDeposits() {
	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(150 * 1e3).MulRaw(1e18)}}
	communityPoolSpendProposal := &types2.CommunityPoolSpendProposal{
		Title:       "community Pool Spend Proposal",
		Description: "description", Recipient: suite.addrs[0].String(),
		Amount: egfCoins,
	}
	minDeposit := keeper.SupportEGFProposalTotalDeposit(false, egfCoins)
	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	communityPoolSpendProposalMsg, err := types.NewMsgSubmitProposal(
		communityPoolSpendProposal,
		initCoins,
		suite.addrs[0],
	)
	suite.Require().NoError(err)
	proposalResponse, err := suite.MsgServer().SubmitProposal(sdk.WrapSDKContext(suite.ctx), communityPoolSpendProposalMsg)
	suite.Require().NoError(err)
	_, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposalResponse.ProposalId, suite.addrs[1])
	suite.Require().False(found)
	proposal, ok := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusDepositPeriod, proposal.Status)
	suite.Require().True(initCoins.IsAllLT(minDeposit))

	// first deposit
	firstCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	votingStarted, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1], firstCoins)
	suite.Require().NoError(err)
	suite.Require().False(votingStarted)
	deposit, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1])
	suite.Require().True(found)
	suite.Require().Equal(firstCoins, deposit.Amount)
	suite.Require().Equal(suite.addrs[1].String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusDepositPeriod, proposal.Status)
	suite.Require().Equal(firstCoins, proposal.TotalDeposit.Sub(initCoins))
	suite.Require().True(initCoins.Add(firstCoins...).IsAllLT(minDeposit))

	// second deposit
	secondCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(9 * 1e3).MulRaw(1e18)}}
	votingStarted, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1], secondCoins)
	suite.Require().NoError(err)
	suite.Require().False(votingStarted)
	deposit, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1])
	suite.Require().True(found)
	suite.Require().Equal(firstCoins.Add(secondCoins...), deposit.Amount)
	suite.Require().Equal(suite.addrs[1].String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusDepositPeriod, proposal.Status)
	suite.Require().True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllLT(minDeposit))

	// third deposit
	thirdCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(5 * 1e3).MulRaw(1e18)}}
	votingStarted, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1], thirdCoins)
	suite.Require().NoError(err)
	suite.Require().True(votingStarted)
	deposit, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, suite.addrs[1])
	suite.Require().True(found)
	suite.Require().Equal(firstCoins.Add(secondCoins...).Add(thirdCoins...), deposit.Amount)
	suite.Require().Equal(suite.addrs[1].String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.Require().True(ok)
	suite.Require().Equal(types.StatusVotingPeriod, proposal.Status)
	suite.Require().True(initCoins.Add(firstCoins...).Add(secondCoins...).Add(thirdCoins...).IsAllGTE(minDeposit))
}
