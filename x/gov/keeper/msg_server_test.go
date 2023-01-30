package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/gov/types"
)

func (suite *KeeperTestSuite) TestSubmitProposal() {
	errInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100).MulRaw(1e18)}}
	suite.True(types.GetInitialDeposit().IsAllGT(errInitCoins))
	TestProposal := govtypes.NewTextProposal("Test", "description")
	errProposalMsg, err := govtypes.NewMsgSubmitProposal(
		TestProposal,
		errInitCoins,
		suite.newAddress(),
	)
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), errProposalMsg)
	suite.Error(err)
	suite.EqualValues(fmt.Sprintf("%v is smaller than %v: initial amount too low", errInitCoins, types.GetInitialDeposit()), err.Error())

	successInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	suite.True(types.GetInitialDeposit().IsAllLTE(successInitCoins))
	successProposalMsg, err := govtypes.NewMsgSubmitProposal(
		govtypes.NewTextProposal("Test", "description"),
		successInitCoins,
		suite.newAddress(),
	)
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), successProposalMsg)
	suite.NoError(err)

	testCases := []struct {
		testName       string
		content        govtypes.Content
		initialDeposit sdk.Coin
		status         govtypes.ProposalStatus
		expectedErr    error
	}{
		{
			testName:       "the deposit is less than the minimum amount",
			content:        TestProposal,
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)},
			status:         govtypes.StatusDepositPeriod,
			expectedErr:    nil,
		},
		{
			testName:       "The deposit is greater than the minimum amount",
			content:        TestProposal,
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			status:         govtypes.StatusVotingPeriod,
			expectedErr:    nil,
		},
		{
			testName: "The deposit is greater than the minimum amount",
			content: &distributiontypes.CommunityPoolSpendProposal{
				Title:       "community Pool Spend Proposal",
				Description: "description",
				Recipient:   sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				Amount:      sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(2000 * 1e3).MulRaw(1e18)}},
			},
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			status:         govtypes.StatusDepositPeriod,
			expectedErr:    nil,
		},
	}
	for _, tc := range testCases {
		if tc.content.ProposalType() == distributiontypes.ProposalTypeCommunityPoolSpend {
			suite.addFundCommunityPool()
		}
		testProposalMsg, err := govtypes.NewMsgSubmitProposal(
			tc.content,
			sdk.Coins{tc.initialDeposit},
			suite.newAddress(),
		)
		suite.NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.True(found)
		suite.EqualValues(tc.status, proposal.Status)
	}
}

func (suite *KeeperTestSuite) TestSubmitEGFProposal() {
	suite.addFundCommunityPool()

	testCases := []struct {
		testName     string
		amount       sdk.Coins
		expect       sdk.Coins
		votingPeriod bool
		expectedErr  error
	}{
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(11 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}},
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(11 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(11 * 1e2).MulRaw(1e18)}},
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
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(20 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
	}
	for _, tc := range testCases {
		testProposalMsg, err := govtypes.NewMsgSubmitProposal(
			&distributiontypes.CommunityPoolSpendProposal{
				Title:       "community Pool Spend Proposal",
				Description: "description",
				Recipient:   sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				Amount:      tc.amount,
			},
			sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}},
			suite.newAddress(),
		)
		suite.NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.True(found)
		if tc.votingPeriod {
			suite.True(tc.expect.IsAllGTE(types.EGFProposalMinDeposit(tc.amount)))
			continue
		}
		suite.True(sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}.IsEqual(proposal.TotalDeposit))
	}
}
