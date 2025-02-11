package keeper_test

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	"github.com/pundiai/fx-core/v8/x/gov/types"
)

func (suite *KeeperTestSuite) TestSubmitProposal() {
	initCoins := helpers.NewStakingCoins(1, 18)
	textProposal := govv1beta1.NewTextProposal("Test", "description")
	legacyContent, err := govv1.NewLegacyContent(textProposal, suite.govAcct)
	suite.Require().NoError(err)

	proposalMsg := suite.NewMsgSubmitProposal(initCoins, suite.newAddress(), legacyContent)
	_, err = suite.msgServer.SubmitProposal(suite.Ctx, proposalMsg)
	suite.ErrorIs(govtypes.ErrMinDepositTooSmall, err)

	differentMsg := suite.NewMsgSubmitProposal(initCoins, suite.newAddress(),
		legacyContent,
		&erc20types.MsgUpdateParams{
			Authority: suite.govAcct, Params: erc20types.DefaultParams(),
		},
	)
	_, err = suite.msgServer.SubmitProposal(suite.Ctx, differentMsg)
	suite.Error(err)
	suite.EqualValues("proposal MsgTypeURL is different: invalid proposal type", err.Error())

	successInitCoins := helpers.NewStakingCoins(10, 18)
	suite.SubmitProposal(successInitCoins, suite.newAddress(), legacyContent)

	testCases := []struct {
		testName       string
		content        govv1beta1.Content
		initialDeposit sdk.Coins
		status         govv1.ProposalStatus
		expectedErr    error
	}{
		{
			testName:       "the deposit is less than the minimum amount",
			content:        textProposal,
			initialDeposit: helpers.NewStakingCoins(10, 18),
			status:         govv1.StatusDepositPeriod,
			expectedErr:    nil,
		},
		{
			testName:       "The deposit is greater than the minimum amount",
			content:        textProposal,
			initialDeposit: helpers.NewStakingCoins(100, 18),
			status:         govv1.StatusVotingPeriod,
			expectedErr:    nil,
		},
	}
	for _, tc := range testCases {
		legacyContent, err = govv1.NewLegacyContent(tc.content, suite.govAcct)
		suite.Require().NoError(err)
		proposal := suite.SubmitProposal(tc.initialDeposit, suite.newAddress(), legacyContent)
		suite.EqualValues(tc.status, proposal.Status)
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
	suite.addFundCommunityPool()
	for _, tc := range testCases {
		spendProposal := &distributiontypes.MsgCommunityPoolSpend{
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			Recipient: helpers.GenAccAddress().String(),
			Amount:    tc.amount,
		}
		proposal := suite.SubmitProposal(helpers.NewStakingCoins(10, 18), suite.newAddress(), spendProposal)
		if tc.votingPeriod {
			minDepositAmount := suite.GetMinInitialDeposit(proposal)
			suite.True(tc.expect.IsAllGTE(minDepositAmount))

			proposal = suite.SubmitProposal(helpers.NewStakingCoins(10, 18), suite.newAddress(), spendProposal, spendProposal, spendProposal)
			suite.EqualValues(govv1.ProposalStatus_PROPOSAL_STATUS_DEPOSIT_PERIOD, proposal.Status, tc.testName)
			continue
		}
		suite.True(helpers.NewStakingCoins(10, 18).Equal(proposal.TotalDeposit))
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
		_, err := suite.msgServer.UpdateStore(suite.Ctx, msg)
		if tc.pass {
			suite.Require().NoError(err)
		} else {
			suite.Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestVote() {
	proposer := suite.AddTestSigner(50_000)

	coins := helpers.NewStakingCoins(1000, 18)
	params := suite.GetParams()

	bankMsg := &banktypes.MsgSend{
		FromAddress: suite.govAcct,
		ToAddress:   proposer.AccAddress().String(),
		Amount:      coins,
	}
	proposal := suite.SubmitProposal(params.MinDeposit, proposer.AccAddress(), bankMsg)

	voteReq := govv1.NewMsgVote(proposer.AccAddress(), proposal.Id, govv1.VoteOption_VOTE_OPTION_YES, "")
	_, err := suite.msgServer.Vote(suite.Ctx, voteReq)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestDeposit() {
	proposer := suite.AddTestSigner(50_000)
	msg := &banktypes.MsgSend{
		FromAddress: suite.govAcct,
		ToAddress:   proposer.AccAddress().String(),
		Amount:      helpers.NewStakingCoins(10, 18),
	}
	proposal := suite.SubmitProposal(helpers.NewStakingCoins(10, 18), proposer.AccAddress(), msg)

	params := suite.GetParams()

	msgDeposit := govv1.NewMsgDeposit(proposer.AccAddress(), proposal.Id, params.MinDeposit)
	_, err := suite.msgServer.Deposit(suite.Ctx, msgDeposit)
	suite.Require().NoError(err)
}
