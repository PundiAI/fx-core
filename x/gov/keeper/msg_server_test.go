package keeper_test

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/gov/types"
)

func (suite *KeeperTestSuite) TestSubmitProposal() {
	errInitCoins := []sdk.Coin{{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(100).Mul(sdkmath.NewInt(1e18))}}
	suite.True(types.GetInitialDeposit().IsAllGT(errInitCoins))
	TestProposal := govv1beta1.NewTextProposal("Test", "description")
	legacyContent, err := govv1.NewLegacyContent(TestProposal, suite.govAcct)
	suite.NoError(err)
	errProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, errInitCoins, suite.newAddress().String(), "")
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), errProposalMsg)
	suite.Error(err)
	suite.EqualValues(fmt.Sprintf("%v is smaller than %v: initial amount too low", errInitCoins, types.GetInitialDeposit()), err.Error())

	successInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}
	suite.True(types.GetInitialDeposit().IsAllLTE(successInitCoins))
	successProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, successInitCoins, suite.newAddress().String(), "")
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), successProposalMsg)
	suite.NoError(err)

	testCases := []struct {
		testName       string
		content        govv1beta1.Content
		initialDeposit sdk.Coin
		status         govv1.ProposalStatus
		expectedErr    error
	}{
		{
			testName:       "the deposit is less than the minimum amount",
			content:        TestProposal,
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)},
			status:         govv1.StatusDepositPeriod,
			expectedErr:    nil,
		},
		{
			testName:       "The deposit is greater than the minimum amount",
			content:        TestProposal,
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(10 * 1e3).MulRaw(1e18)},
			status:         govv1.StatusVotingPeriod,
			expectedErr:    nil,
		},
		{
			testName: "The deposit is greater than the minimum amount",
			content: distributiontypes.NewCommunityPoolSpendProposal(
				"community Pool Spend Proposal",
				"description",
				sdk.AccAddress(helpers.GenerateAddress().Bytes()),
				sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(2000 * 1e3).MulRaw(1e18)}},
			),
			initialDeposit: sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(10 * 1e3).MulRaw(1e18)},
			status:         govv1.StatusDepositPeriod,
			expectedErr:    nil,
		},
	}
	for _, tc := range testCases {
		if tc.content.ProposalType() == distributiontypes.ProposalTypeCommunityPoolSpend {
			suite.addFundCommunityPool()
		}
		legacyContent, err = govv1.NewLegacyContent(tc.content, suite.govAcct)
		suite.NoError(err)
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, sdk.NewCoins(tc.initialDeposit), suite.newAddress().String(), "")
		suite.NoError(err)
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
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(10 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(11 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}},
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(11 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(11 * 1e2).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(200 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(16 * 1e3).MulRaw(1e18)}},
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "",
			amount:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(200 * 1e3).MulRaw(1e18)}},
			expect:       sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(20 * 1e3).MulRaw(1e18)}},
			votingPeriod: true,
			expectedErr:  nil,
		},
	}
	for _, tc := range testCases {
		LegacyContentMsg, err := govv1.NewLegacyContent(distributiontypes.NewCommunityPoolSpendProposal(
			"community Pool Spend Proposal",
			"description",
			sdk.AccAddress(helpers.GenerateAddress().Bytes()),
			tc.amount), suite.govAcct)
		suite.NoError(err)
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{LegacyContentMsg}, sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}, suite.newAddress().String(), "")
		suite.NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.True(found)
		if tc.votingPeriod {
			suite.True(tc.expect.IsAllGTE(types.EGFProposalMinDeposit(tc.amount)))
			continue
		}
		suite.True(sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}.IsEqual(proposal.TotalDeposit))
	}
}
