package keeper_test

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
)

func (suite *KeeperTestSuite) TestSubmitProposal() {
	errInitCoins := []sdk.Coin{{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(100).Mul(sdkmath.NewInt(1e18))}}
	TestProposal := govv1beta1.NewTextProposal("Test", "description")
	legacyContent, err := govv1.NewLegacyContent(TestProposal, suite.govAcct)
	suite.NoError(err)
	initialDeposit := suite.app.GovKeeper.GetMinInitialDeposit(suite.ctx, legacyContent.Content.TypeUrl)
	suite.True(sdk.NewCoins(initialDeposit).IsAllGT(errInitCoins))
	errProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, errInitCoins, suite.newAddress().String(),
		"", TestProposal.GetTitle(), TestProposal.GetDescription())
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), errProposalMsg)
	suite.Error(err)
	suite.EqualValues(fmt.Sprintf("%v is smaller than %v: initial amount too low", errInitCoins, initialDeposit), err.Error())

	differentMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent, &erc20types.MsgUpdateParams{Authority: suite.govAcct, Params: erc20types.DefaultParams()}}, errInitCoins, suite.newAddress().String(),
		"", TestProposal.GetTitle(), TestProposal.GetDescription())
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), differentMsg)
	suite.Error(err)
	suite.EqualValues("proposal MsgTypeURL is different: invalid proposal type", err.Error())

	successInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}
	suite.True(sdk.NewCoins(initialDeposit).IsAllLTE(successInitCoins))
	successProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, successInitCoins, suite.newAddress().String(),
		"", TestProposal.GetTitle(), TestProposal.GetDescription())
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
	}
	for _, tc := range testCases {
		legacyContent, err = govv1.NewLegacyContent(tc.content, suite.govAcct)
		suite.NoError(err)
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, sdk.NewCoins(tc.initialDeposit), suite.newAddress().String(),
			"", tc.content.GetTitle(), tc.content.GetDescription())
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
		spendProposal := &distributiontypes.MsgCommunityPoolSpend{
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			Recipient: helpers.GenAccAddress().String(),
			Amount:    tc.amount,
		}
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal},
			sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}},
			suite.newAddress().String(),
			"", "community Pool Spend Proposal", "description")
		suite.NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.True(found)
		if tc.votingPeriod {
			suite.True(tc.expect.IsAllGTE(suite.app.GovKeeper.EGFProposalMinDeposit(suite.ctx, sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}), tc.amount)))
			manyProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal, spendProposal, spendProposal},
				sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}},
				suite.newAddress().String(),
				"", "community Pool Spend Proposal", "description")
			suite.NoError(err)
			proposalResponse, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), manyProposalMsg)
			suite.NoError(err)
			proposal, found = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
			suite.True(found)
			suite.Require().EqualValues(proposal.Status, govv1.ProposalStatus_PROPOSAL_STATUS_DEPOSIT_PERIOD)
			continue
		}
		suite.True(sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}.IsEqual(proposal.TotalDeposit))
	}
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
		_, err := suite.MsgServer.UpdateStore(sdk.WrapSDKContext(suite.ctx), msg)
		if tc.pass {
			suite.NoError(err)
		} else {
			suite.Error(err)
		}

	}
}

func (suite *KeeperTestSuite) TestVoteReq() {
	govAcct := suite.govAcct
	addrs := suite.addrs
	proposer := addrs[0]

	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).MulRaw(1e18)))
	minDeposit := suite.app.GovKeeper.GetParams(suite.ctx).MinDeposit
	bankMsg := &banktypes.MsgSend{
		FromAddress: govAcct,
		ToAddress:   proposer.String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		minDeposit,
		proposer.String(),
		"", "send", "send",
	)
	suite.Require().NoError(err)

	res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
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
					proposer.String(),
					"", "send", "send",
				)
				suite.Require().NoError(err)

				res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
				return res.ProposalId
			},
			option:    govv1.VoteOption_VOTE_OPTION_YES,
			voter:     proposer,
			metadata:  "",
			expErr:    true,
			expErrMsg: "inactive proposal",
		},
		"metadata too long": {
			preRun: func() uint64 {
				return proposalId
			},
			option:    govv1.VoteOption_VOTE_OPTION_YES,
			voter:     proposer,
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
					minDeposit,
					proposer.String(),
					"", "send", "send",
				)
				suite.Require().NoError(err)

				res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
				return res.ProposalId
			},
			option:   govv1.VoteOption_VOTE_OPTION_YES,
			voter:    proposer,
			metadata: "",
			expErr:   false,
		},
	}

	for name, tc := range cases {
		suite.Run(name, func() {
			pId := tc.preRun()
			voteReq := govv1.NewMsgVote(tc.voter, pId, tc.option, tc.metadata)
			_, err := suite.msgServer.Vote(suite.ctx, voteReq)
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
	govAcct := suite.govAcct
	addrs := suite.addrs
	proposer := addrs[0]

	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).MulRaw(1e18)))
	minDeposit := suite.app.GovKeeper.GetParams(suite.ctx).MinDeposit
	bankMsg := &banktypes.MsgSend{
		FromAddress: govAcct,
		ToAddress:   proposer.String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		coins,
		proposer.String(),
		"", "send", "send",
	)
	suite.Require().NoError(err)

	res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
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
			depositor: proposer,
			deposit:   coins,
			expErr:    true,
			options:   govv1.NewNonSplitVoteOption(govv1.OptionYes),
		},
		"all good": {
			preRun: func() uint64 {
				return pId
			},
			depositor: proposer,
			deposit:   minDeposit,
			expErr:    false,
			options:   govv1.NewNonSplitVoteOption(govv1.OptionYes),
		},
	}

	for name, tc := range cases {
		suite.Run(name, func() {
			proposalId := tc.preRun()
			depositReq := govv1.NewMsgDeposit(tc.depositor, proposalId, tc.deposit)
			_, err := suite.msgServer.Deposit(suite.ctx, depositReq)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

// legacy msg server tests
func (suite *KeeperTestSuite) TestLegacyMsgSubmitProposal() {
	addrs := suite.addrs
	proposer := addrs[0]

	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).MulRaw(1e18)))
	initialDeposit := coins
	minDeposit := suite.app.GovKeeper.GetParams(suite.ctx).MinDeposit

	cases := map[string]struct {
		preRun func() (*govv1beta1.MsgSubmitProposal, error)
		expErr bool
	}{
		"all good": {
			preRun: func() (*govv1beta1.MsgSubmitProposal, error) {
				return govv1beta1.NewMsgSubmitProposal(
					govv1beta1.NewTextProposal("test", "I am test"),
					initialDeposit,
					proposer,
				)
			},
			expErr: false,
		},
		"all good with min deposit": {
			preRun: func() (*govv1beta1.MsgSubmitProposal, error) {
				return govv1beta1.NewMsgSubmitProposal(
					govv1beta1.NewTextProposal("test", "I am test"),
					minDeposit,
					proposer,
				)
			},
			expErr: false,
		},
	}

	for name, c := range cases {
		suite.Run(name, func() {
			msg, err := c.preRun()
			suite.Require().NoError(err)
			res, err := suite.legacyMsgServer.SubmitProposal(suite.ctx, msg)
			if c.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestLegacyMsgVote() {
	govAcct := suite.govAcct
	addrs := suite.addrs
	proposer := addrs[0]

	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).MulRaw(1e18)))
	minDeposit := suite.app.GovKeeper.GetParams(suite.ctx).MinDeposit
	bankMsg := &banktypes.MsgSend{
		FromAddress: govAcct,
		ToAddress:   proposer.String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		minDeposit,
		proposer.String(),
		"", "send", "send",
	)
	suite.Require().NoError(err)

	res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProposalId)
	proposalId := res.ProposalId

	cases := map[string]struct {
		preRun    func() uint64
		expErr    bool
		expErrMsg string
		option    govv1beta1.VoteOption
		metadata  string
		voter     sdk.AccAddress
	}{
		"vote on inactive proposal": {
			preRun: func() uint64 {
				msg, err := govv1.NewMsgSubmitProposal(
					[]sdk.Msg{bankMsg},
					coins,
					proposer.String(),
					"", "send", "send",
				)
				suite.Require().NoError(err)

				res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
				return res.ProposalId
			},
			option:    govv1beta1.OptionYes,
			voter:     proposer,
			metadata:  "",
			expErr:    true,
			expErrMsg: "inactive proposal",
		},
		"voter error": {
			preRun: func() uint64 {
				return proposalId
			},
			option:    govv1beta1.OptionYes,
			voter:     sdk.AccAddress(strings.Repeat("a", 300)),
			metadata:  "",
			expErr:    true,
			expErrMsg: "address max length is 255",
		},
		"all good": {
			preRun: func() uint64 {
				msg, err := govv1.NewMsgSubmitProposal(
					[]sdk.Msg{bankMsg},
					minDeposit,
					proposer.String(),
					"", "send", "send",
				)
				suite.Require().NoError(err)

				res, err := suite.msgServer.SubmitProposal(suite.ctx, msg)
				suite.Require().NoError(err)
				suite.Require().NotNil(res.ProposalId)
				return res.ProposalId
			},
			option:   govv1beta1.OptionYes,
			voter:    proposer,
			metadata: "",
			expErr:   false,
		},
	}

	for name, tc := range cases {
		suite.Run(name, func() {
			pId := tc.preRun()
			voteReq := govv1beta1.NewMsgVote(tc.voter, pId, tc.option)
			_, err := suite.legacyMsgServer.Vote(suite.ctx, voteReq)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
