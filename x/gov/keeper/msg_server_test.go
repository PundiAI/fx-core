package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"

	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/gov/keeper"
	types3 "github.com/functionx/fx-core/x/gov/types"
)

func (suite *KeeperTestSuite) TestSubmitProposal() {
	errInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e2).MulRaw(1e18)}}
	suite.Require().True(types3.InitialDeposit.IsAllGT(errInitCoins))
	TestProposal := types.NewTextProposal("Test", "description")
	errProposalMsg, err := types.NewMsgSubmitProposal(
		TestProposal,
		errInitCoins,
		suite.addrs[0],
	)
	suite.Require().NoError(err)
	_, err = suite.MsgServer().SubmitProposal(sdk.WrapSDKContext(suite.ctx), errProposalMsg)
	suite.Require().Error(err)
	suite.Require().EqualValues(fmt.Sprintf("%v is smaller than %v: initial amount too low", errInitCoins.String(), types3.InitialDeposit.String()), err.Error())

	successInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	suite.Require().True(types3.InitialDeposit.IsAllLTE(successInitCoins))
	successProposalMsg, err := types.NewMsgSubmitProposal(
		types.NewTextProposal("Test", "description"),
		successInitCoins,
		suite.addrs[0],
	)
	suite.Require().NoError(err)
	_, err = suite.MsgServer().SubmitProposal(sdk.WrapSDKContext(suite.ctx), successProposalMsg)
	suite.Require().NoError(err)

	testCases := []struct {
		testName       string
		content        types.Content
		initialDeposit sdk.Coin
		status         types.ProposalStatus
		expectedErr    error
	}{
		{
			testName:       "the deposit is less than the minimum amount",
			content:        TestProposal,
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)},
			status:         types.StatusDepositPeriod,
			expectedErr:    nil,
		},
		{
			testName:       "The deposit is greater than the minimum amount",
			content:        TestProposal,
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			status:         types.StatusVotingPeriod,
			expectedErr:    nil,
		},
		{
			testName: "The deposit is greater than the minimum amount",
			content: &types2.CommunityPoolSpendProposal{
				Title:       "community Pool Spend Proposal",
				Description: "description", Recipient: suite.addrs[0].String(),
				Amount: sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(2000 * 1e3).MulRaw(1e18)}}},
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			status:         types.StatusDepositPeriod,
			expectedErr:    nil,
		},
	}
	for _, tc := range testCases {
		testProposalMsg, err := types.NewMsgSubmitProposal(
			tc.content,
			sdk.Coins{tc.initialDeposit},
			suite.addrs[0],
		)
		suite.Require().NoError(err)
		proposalResponse, err := suite.MsgServer().SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.Require().NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.Require().True(found)
		suite.Require().EqualValues(tc.status, proposal.Status)
	}
}

func (suite *KeeperTestSuite) TestSubmitEGFProposal() {
	testCases := []struct {
		testName     string
		amount       sdk.Coins
		expect       sdk.Coins
		votingPeriod bool
		expectedErr  error
	}{
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(150 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}},
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(150 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(16 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(200 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(16 * 1e3).MulRaw(1e18)}},
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(200 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(21 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
	}
	for _, tc := range testCases {
		testProposalMsg, err := types.NewMsgSubmitProposal(
			&types2.CommunityPoolSpendProposal{
				Title:       "community Pool Spend Proposal",
				Description: "description", Recipient: suite.addrs[0].String(),
				Amount: tc.amount},
			sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}},
			suite.addrs[0],
		)
		suite.Require().NoError(err)
		proposalResponse, err := suite.MsgServer().SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.Require().NoError(err)
		_, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.Require().True(found)
		if tc.votingPeriod {
			suite.Require().True(tc.expect.IsAllGTE(keeper.SupportEGFProposalTotalDeposit(true, tc.amount)))
			continue
		}
		suite.Require().True(sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}.IsEqual(types3.InitialDeposit))
	}
}
