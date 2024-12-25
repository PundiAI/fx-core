package keeper_test

import (
	"encoding/hex"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	"github.com/pundiai/fx-core/v8/x/gov/types"
)

func (suite *KeeperTestSuite) TestSubmitProposal() {
	suite.T().Skip()
	/*errInitCoin := helpers.NewStakingCoin(100, 18)
	TestProposal := govv1beta1.NewTextProposal("Test", "description")
	legacyContent, err := govv1.NewLegacyContent(TestProposal, suite.govAcct)
	suite.Require().NoError(err)
	initialDeposit := suite.App.GovKeeper.GetMinInitialDeposit(suite.Ctx, legacyContent.Content.TypeUrl)
	suite.True(initialDeposit.IsGTE(errInitCoin))
	errProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, sdk.NewCoins(errInitCoin), suite.newAddress().String(),
		"", TestProposal.GetTitle(), TestProposal.GetDescription(), false)
	suite.Require().NoError(err)
	_, err = suite.msgServer.SubmitProposal(suite.Ctx, errProposalMsg)
	suite.Error(err)
	suite.EqualValues(fmt.Sprintf("initial deposit must be at least %s: invalid request", initialDeposit), err.Error())

	differentMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent, &erc20types.MsgUpdateParams{Authority: suite.govAcct, Params: erc20types.DefaultParams()}}, sdk.NewCoins(errInitCoin), suite.newAddress().String(),
		"", TestProposal.GetTitle(), TestProposal.GetDescription(), false)
	suite.Require().NoError(err)
	_, err = suite.msgServer.SubmitProposal(suite.Ctx, differentMsg)
	suite.Error(err)
	suite.EqualValues("proposal MsgTypeURL is different: invalid proposal type", err.Error())

	successInitCoin := helpers.NewStakingCoin(1000, 18)
	suite.True(initialDeposit.IsLTE(successInitCoin))
	successProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, sdk.NewCoins(successInitCoin), suite.newAddress().String(),
		"", TestProposal.GetTitle(), TestProposal.GetDescription(), false)
	suite.Require().NoError(err)
	_, err = suite.msgServer.SubmitProposal(suite.Ctx, successProposalMsg)
	suite.Require().NoError(err)

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
			initialDeposit: helpers.NewStakingCoin(1000, 18),
			status:         govv1.StatusDepositPeriod,
			expectedErr:    nil,
		},
		{
			testName:       "The deposit is greater than the minimum amount",
			content:        TestProposal,
			initialDeposit: helpers.NewStakingCoin(10_000, 18),
			status:         govv1.StatusVotingPeriod,
			expectedErr:    nil,
		},
	}
	for _, tc := range testCases {
		legacyContent, err = govv1.NewLegacyContent(tc.content, suite.govAcct)
		suite.Require().NoError(err)
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, sdk.NewCoins(tc.initialDeposit), suite.newAddress().String(),
			"", tc.content.GetTitle(), tc.content.GetDescription(), false)
		suite.Require().NoError(err)
		suite.Require().NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(suite.Ctx, testProposalMsg)
		suite.Require().NoError(err)
		proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
		suite.Require().NoError(err)
		suite.EqualValues(tc.status, proposal.Status)
	}*/
}

func (suite *KeeperTestSuite) TestSubmitEGFProposal() {
	suite.T().Skip()
	/*suite.addFundCommunityPool()

	testCases := []struct {
		testName     string
		amount       sdk.Coins
		expect       sdk.Coins
		votingPeriod bool
		expectedErr  error
	}{
		{
			testName:     "1",
			amount:       helpers.NewStakingCoins(10_000, 18),
			expect:       helpers.NewStakingCoins(1000, 18),
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "2",
			amount:       helpers.NewStakingCoins(11_000, 18),
			expect:       helpers.NewStakingCoins(1000, 18),
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "3",
			amount:       helpers.NewStakingCoins(11_000, 18),
			expect:       helpers.NewStakingCoins(1100, 18),
			votingPeriod: true,
			expectedErr:  nil,
		},
		{
			testName:     "4",
			amount:       helpers.NewStakingCoins(200_000, 18),
			expect:       helpers.NewStakingCoins(16_000, 18),
			votingPeriod: false,
			expectedErr:  nil,
		},
		{
			testName:     "5",
			amount:       helpers.NewStakingCoins(200_000, 18),
			expect:       helpers.NewStakingCoins(20_000, 18),
			votingPeriod: true,
			expectedErr:  nil,
		},
	}
	for _, tc := range testCases {
		spendProposal := &distributiontypes.MsgCommunityPoolSpend{
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			Recipient: helpers.GenAccAddress().String(),
			Amount:    tc.amount,
		}
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal},
			helpers.NewStakingCoins(1000, 18),
			suite.newAddress().String(),
			"", "community Pool Spend Proposal", "description", false)
		suite.Require().NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(suite.Ctx, testProposalMsg)
		suite.Require().NoError(err)
		proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
		suite.Require().NoError(err)
		if tc.votingPeriod {
			suite.True(tc.expect.IsAllGTE(suite.App.GovKeeper.EGFProposalMinDeposit(suite.Ctx, sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}), tc.amount)))
			manyProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal, spendProposal, spendProposal},
				helpers.NewStakingCoins(1000, 18),
				suite.newAddress().String(),
				"", "community Pool Spend Proposal", "description", false)
			suite.Require().NoError(err)
			proposalResponse, err = suite.msgServer.SubmitProposal(suite.Ctx, manyProposalMsg)
			suite.Require().NoError(err)
			proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
			suite.Require().NoError(err)
			suite.Require().EqualValues(proposal.Status, govv1.ProposalStatus_PROPOSAL_STATUS_DEPOSIT_PERIOD)
			continue
		}
		suite.True(helpers.NewStakingCoins(1000, 18).Equal(proposal.TotalDeposit))
	}*/
}

func (suite *KeeperTestSuite) TestSubmitUpdateStoreProposal() {
	testCases := []struct {
		testName     string
		updateStores []types.UpdateStore
		pass         bool
	}{
		{
			testName: "success",
			updateStores: []types.UpdateStore{
				{
					Space:    "eth",
					Key:      hex.EncodeToString(crosschaintypes.LastObservedBlockHeightKey),
					OldValue: "",
					Value:    "01",
				},
			},
			pass: true,
		},
		{
			testName: "invalid store space",
			updateStores: []types.UpdateStore{
				{
					Space:    "eth1",
					Key:      hex.EncodeToString(crosschaintypes.LastObservedBlockHeightKey),
					OldValue: "",
					Value:    "01",
				},
			},
		},
		{
			testName: "invalid old value",
			updateStores: []types.UpdateStore{
				{
					Space:    "eth1",
					Key:      hex.EncodeToString(crosschaintypes.LastObservedBlockHeightKey),
					OldValue: "01",
					Value:    "01",
				},
			},
		},
	}
	for _, tc := range testCases {
		msg := types.NewMsgUpdateStore(authtypes.NewModuleAddress(govtypes.ModuleName).String(), tc.updateStores)
		_, err := suite.msgServer.UpdateStore(suite.Ctx, msg)
		if tc.pass {
			suite.Require().NoError(err)
		} else {
			suite.Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestVoteReq() {
	suite.T().Skip()
	proposer := suite.AddTestSigner(50_000)

	coins := helpers.NewStakingCoins(1000, 18)
	params, err := suite.App.GovKeeper.Keeper.Params.Get(suite.Ctx)
	suite.Require().NoError(err)
	bankMsg := &banktypes.MsgSend{
		FromAddress: suite.govAcct,
		ToAddress:   proposer.AccAddress().String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		params.MinDeposit,
		proposer.AccAddress().String(),
		"", "send", "send", false,
	)
	suite.Require().NoError(err)

	res, err := suite.msgServer.SubmitProposal(suite.Ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProposalId)
	proposalId := res.ProposalId

	cases := map[string]struct {
		preRun    func() uint64
		expErr    bool
		expErrMsg string
		option    govv1.VoteOption
		metadata  string
		voter     sdk.AccAddress
	}{
		"vote on inactive proposal": {
			preRun: func() uint64 {
				msg, err := govv1.NewMsgSubmitProposal(
					[]sdk.Msg{bankMsg},
					coins,
					proposer.AccAddress().String(),
					"", "send", "send", false,
				)
				suite.Require().NoError(err)

				res, err := suite.msgServer.SubmitProposal(suite.Ctx, msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
				return res.ProposalId
			},
			option:    govv1.VoteOption_VOTE_OPTION_YES,
			voter:     proposer.AccAddress(),
			metadata:  "",
			expErr:    true,
			expErrMsg: "inactive proposal",
		},
		"metadata too long": {
			preRun: func() uint64 {
				return proposalId
			},
			option:    govv1.VoteOption_VOTE_OPTION_YES,
			voter:     proposer.AccAddress(),
			metadata:  strings.Repeat("a", 10240),
			expErr:    true,
			expErrMsg: "metadata too long",
		},
		"voter error": {
			preRun: func() uint64 {
				return proposalId
			},
			option:    govv1.VoteOption_VOTE_OPTION_YES,
			voter:     sdk.AccAddress(strings.Repeat("a", 300)),
			metadata:  "",
			expErr:    true,
			expErrMsg: "address max length is 255",
		},
		"all good": {
			preRun: func() uint64 {
				msg, err := govv1.NewMsgSubmitProposal(
					[]sdk.Msg{bankMsg},
					params.MinDeposit,
					proposer.AccAddress().String(),
					"", "send", "send", false,
				)
				suite.Require().NoError(err)

				res, err := suite.msgServer.SubmitProposal(suite.Ctx, msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
				return res.ProposalId
			},
			option:   govv1.VoteOption_VOTE_OPTION_YES,
			voter:    proposer.AccAddress(),
			metadata: "",
			expErr:   false,
		},
	}

	for name, tc := range cases {
		suite.Run(name, func() {
			pId := tc.preRun()
			voteReq := govv1.NewMsgVote(tc.voter, pId, tc.option, tc.metadata)
			_, err := suite.msgServer.Vote(suite.Ctx, voteReq)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestDepositReq() {
	suite.T().Skip()
	proposer := suite.AddTestSigner(50_000)

	coins := helpers.NewStakingCoins(1000, 18)
	params, err := suite.App.GovKeeper.Keeper.Params.Get(suite.Ctx)
	suite.Require().NoError(err)
	bankMsg := &banktypes.MsgSend{
		FromAddress: suite.govAcct,
		ToAddress:   proposer.AccAddress().String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		coins,
		proposer.AccAddress().String(),
		"", "send", "send", false,
	)
	suite.Require().NoError(err)

	res, err := suite.msgServer.SubmitProposal(suite.Ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProposalId)
	pId := res.ProposalId

	cases := map[string]struct {
		preRun     func() uint64
		expErr     bool
		proposalId uint64
		depositor  sdk.AccAddress
		deposit    sdk.Coins
		options    govv1.WeightedVoteOptions
	}{
		"wrong proposal id": {
			preRun: func() uint64 {
				return 0
			},
			depositor: proposer.AccAddress(),
			deposit:   coins,
			expErr:    true,
			options:   govv1.NewNonSplitVoteOption(govv1.OptionYes),
		},
		"all good": {
			preRun: func() uint64 {
				return pId
			},
			depositor: proposer.AccAddress(),
			deposit:   params.MinDeposit,
			expErr:    false,
			options:   govv1.NewNonSplitVoteOption(govv1.OptionYes),
		},
	}

	for name, tc := range cases {
		suite.Run(name, func() {
			proposalId := tc.preRun()
			depositReq := govv1.NewMsgDeposit(tc.depositor, proposalId, tc.deposit)
			_, err := suite.msgServer.Deposit(suite.Ctx, depositReq)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
