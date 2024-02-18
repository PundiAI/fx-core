package keeper_test

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
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
		types.NewFXMetadata(TestProposal.GetTitle(), TestProposal.GetDescription(), "").String())
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), errProposalMsg)
	suite.Error(err)
	suite.EqualValues(fmt.Sprintf("%v is smaller than %v: initial amount too low", errInitCoins, initialDeposit), err.Error())

	differentMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent, &erc20types.MsgUpdateParams{Authority: suite.govAcct, Params: erc20types.DefaultParams()}}, errInitCoins, suite.newAddress().String(),
		types.NewFXMetadata(TestProposal.GetTitle(), TestProposal.GetDescription(), "").String())
	suite.NoError(err)
	_, err = suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), differentMsg)
	suite.Error(err)
	suite.EqualValues("proposal MsgTypeURL is different: invalid proposal content", err.Error())

	successInitCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}
	suite.True(sdk.NewCoins(initialDeposit).IsAllLTE(successInitCoins))
	successProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, successInitCoins, suite.newAddress().String(),
		types.NewFXMetadata(TestProposal.GetTitle(), TestProposal.GetDescription(), "").String())
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
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, sdk.NewCoins(tc.initialDeposit), suite.newAddress().String(),
			types.NewFXMetadata(tc.content.GetTitle(), tc.content.GetDescription(), "").String())
		suite.NoError(err)
		suite.NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.True(found)
		suite.EqualValues(tc.status, proposal.Status)
	}
}

func (suite *KeeperTestSuite) TestSubmitProposalFXMetadata() {
	chars := "1234567890"
	tooLongTitle := func() string {
		title := ""
		for uint64(len(title)) <= suite.app.GovKeeper.Config().MaxTitleLen {
			title = title + chars
		}
		return title
	}
	tooLongSummary := func() string {
		summary := ""
		for uint64(len(summary)) <= suite.app.GovKeeper.Config().MaxSummaryLen {
			summary = summary + chars
		}
		return summary
	}
	tooLongMetadata := func() string {
		metadata := ""
		for uint64(len(metadata)) <= suite.app.GovKeeper.Config().MaxSummaryLen {
			metadata = metadata + chars
		}
		return base64.StdEncoding.EncodeToString([]byte(metadata))
	}

	testCases := []struct {
		name        string
		fxMetadata  func() string
		result      bool
		expectedErr error
	}{
		{
			name: "ok",
			fxMetadata: func() string {
				return types.NewFXMetadata("test", "test", "").String()
			},
			result: true,
		},
		{
			name: "error -  empty fx metadata",
			fxMetadata: func() string {
				return ""
			},
			result:      false,
			expectedErr: errorsmod.Wrap(errortypes.ErrInvalidRequest, "invalid fx metadata content: fx metadata cannot be empty"),
		},
		{
			name: "error -  invalid fx metadata",
			fxMetadata: func() string {
				return "abc"
			},
			result: false,

			expectedErr: errortypes.ErrInvalidRequest.Wrap("invalid fx metadata content: illegal base64 data at input byte 0"),
		},
		{
			name: "error -  invalid fx metadata",
			fxMetadata: func() string {
				return "eyJ0aXRsZSI6InRlc3QiLCJzdW1tYXJ5IjoidGVzdCIsIm1ldGFkYXRhIjoiIn0xMTEx" // {"title":"test","summary":"test","metadata":""}1111
			},
			result:      false,
			expectedErr: errortypes.ErrInvalidRequest.Wrap("invalid fx metadata content: invalid character '1' after top-level value"),
		},
		{
			name: "error -  empty title",
			fxMetadata: func() string {
				return types.NewFXMetadata("", "test", "").String()
			},
			result:      false,
			expectedErr: errorsmod.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			name: "error -  too long title",
			fxMetadata: func() string {
				return types.NewFXMetadata(tooLongTitle(), "test", "").String()
			},
			result:      false,
			expectedErr: errorsmod.Wrap(govtypes.ErrInvalidProposalContent, fmt.Sprintf("proposal title is longer than max length of %d", suite.app.GovKeeper.Config().MaxTitleLen)),
		},
		{
			name: "error -  empty summary",
			fxMetadata: func() string {
				return types.NewFXMetadata("test", "", "").String()
			},
			result:      false,
			expectedErr: errorsmod.Wrap(govtypes.ErrInvalidProposalContent, "proposal summary cannot be blank"),
		},
		{
			name: "error -  too long summary",
			fxMetadata: func() string {
				return types.NewFXMetadata("test", tooLongSummary(), "").String()
			},
			result:      false,
			expectedErr: errorsmod.Wrap(govtypes.ErrInvalidProposalContent, fmt.Sprintf("proposal summary is longer than max length of %d", suite.app.GovKeeper.Config().MaxSummaryLen)),
		},
		{
			name: "error -  too long metadata",
			fxMetadata: func() string {
				return types.NewFXMetadata("test", "test", tooLongMetadata()).String()
			},
			result:      false,
			expectedErr: errorsmod.Wrap(govtypes.ErrInvalidProposalContent, fmt.Sprintf("proposal metadata is longer than max length of %d", suite.app.GovKeeper.Config().MaxMetadataLen)),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			initCoins := []sdk.Coin{{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))}}
			updateErc20Params := erc20types.MsgUpdateParams{
				Authority: suite.govAcct,
				Params:    erc20types.Params{EnableErc20: true, EnableEVMHook: true, IbcTimeout: 12 * time.Hour},
			}

			fxMD := tc.fxMetadata()

			errProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{&updateErc20Params}, initCoins, suite.newAddress().String(), fxMD)
			suite.Require().NoError(err)

			proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), errProposalMsg)
			if tc.result {
				suite.NoError(err)
				proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
				suite.True(found)
				suite.Require().Equal(proposal.Metadata, fxMD)
			} else {
				suite.Equal(err.Error(), tc.expectedErr.Error())
			}
		})
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
		spendProposal := distributiontypes.NewCommunityPoolSpendProposal(
			"community Pool Spend Proposal", "description",
			sdk.AccAddress(helpers.GenerateAddress().Bytes()), tc.amount)
		LegacyContentMsg, err := govv1.NewLegacyContent(spendProposal, suite.govAcct)
		suite.NoError(err)
		testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{LegacyContentMsg},
			sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}},
			suite.newAddress().String(),
			types.NewFXMetadata(spendProposal.GetTitle(), spendProposal.GetDescription(), "").String())
		suite.NoError(err)
		proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
		suite.NoError(err)
		proposal, found := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
		suite.True(found)
		if tc.votingPeriod {
			suite.True(tc.expect.IsAllGTE(suite.app.GovKeeper.EGFProposalMinDeposit(suite.ctx, LegacyContentMsg.Content.TypeUrl, tc.amount)))
			manyProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{LegacyContentMsg, LegacyContentMsg, LegacyContentMsg},
				sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}},
				suite.newAddress().String(),
				types.NewFXMetadata(spendProposal.GetTitle(), spendProposal.GetDescription(), "").String())
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

func (suite *KeeperTestSuite) TestVoteReq() {
	govAcct := suite.govAcct
	addrs := suite.addrs
	proposer := addrs[0]

	coins := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).MulRaw(1e18)))
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
	bankMsg := &banktypes.MsgSend{
		FromAddress: govAcct,
		ToAddress:   proposer.String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		minDeposit,
		proposer.String(),
		types.NewFXMetadata("send", "send", "").String(),
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
					types.NewFXMetadata("send", "send", "").String(),
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
			metadata:  strings.Repeat("a", 300),
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
					types.NewFXMetadata("send", "send", "").String(),
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
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
	bankMsg := &banktypes.MsgSend{
		FromAddress: govAcct,
		ToAddress:   proposer.String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		coins,
		proposer.String(),
		types.NewFXMetadata("send", "send", "").String(),
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
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit

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
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
	bankMsg := &banktypes.MsgSend{
		FromAddress: govAcct,
		ToAddress:   proposer.String(),
		Amount:      coins,
	}

	msg, err := govv1.NewMsgSubmitProposal(
		[]sdk.Msg{bankMsg},
		minDeposit,
		proposer.String(),
		types.NewFXMetadata("send", "send", "").String(),
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
					types.NewFXMetadata("send", "send", "").String(),
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
					types.NewFXMetadata("send", "send", "").String(),
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
