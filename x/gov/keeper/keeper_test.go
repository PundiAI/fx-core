package keeper_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	"github.com/pundiai/fx-core/v8/x/gov/keeper"
	"github.com/pundiai/fx-core/v8/x/gov/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	govAcct     string
	msgServer   types.MsgServerPro
	queryClient types.QueryClient
}

func TestGovKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.MintValNumber = 1
	suite.BaseSuite.SetupTest()

	suite.govAcct = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.msgServer = keeper.NewMsgServerImpl(suite.GetKeeper())

	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(suite.GetKeeper()))
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) GetKeeper() *keeper.Keeper {
	return suite.App.GovKeeper
}

func (suite *KeeperTestSuite) addFundCommunityPool() {
	sender := suite.AddTestSigner(5 * 1e8)
	balances := suite.Balance(sender.AccAddress())

	err := suite.App.DistrKeeper.FundCommunityPool(suite.Ctx, balances, sender.AccAddress())
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) newAddress() sdk.AccAddress {
	return suite.AddTestSigner(50_000).AccAddress()
}

func (suite *KeeperTestSuite) GetProposal(id uint64) govv1.Proposal {
	proposal, err := suite.GetKeeper().Keeper.Proposals.Get(suite.Ctx, id)
	suite.Require().NoError(err)
	return proposal
}

func (suite *KeeperTestSuite) GetParams() govv1.Params {
	params, err := suite.GetKeeper().Params.Get(suite.Ctx)
	suite.Require().NoError(err)
	return params
}

func (suite *KeeperTestSuite) GetDeposit(proposalID uint64, depositor sdk.AccAddress) govv1.Deposit {
	deposit, err := suite.GetKeeper().Deposits.Get(suite.Ctx, collections.Join(proposalID, depositor))
	suite.Require().NoError(err)
	return deposit
}

func (suite *KeeperTestSuite) NewMsgSubmitProposal(initialDeposit sdk.Coins, proposer sdk.AccAddress, messages ...sdk.Msg) *govv1.MsgSubmitProposal {
	proposalMsg, err := govv1.NewMsgSubmitProposal(messages, initialDeposit, proposer.String(), "", "title", "description", false)
	suite.Require().NoError(err)
	return proposalMsg
}

func (suite *KeeperTestSuite) SubmitProposal(initialDeposit sdk.Coins, proposer sdk.AccAddress, messages ...sdk.Msg) govv1.Proposal {
	proposalMsg := suite.NewMsgSubmitProposal(initialDeposit, proposer, messages...)
	response, err := suite.msgServer.SubmitProposal(suite.Ctx, proposalMsg)
	suite.Require().NoError(err)
	suite.NotZero(response.ProposalId)
	suite.GetDeposit(response.ProposalId, proposer)
	return suite.GetProposal(response.ProposalId)
}

func (suite *KeeperTestSuite) GetMinInitialDeposit(proposal govv1.Proposal) sdk.Coins {
	params, err := suite.GetKeeper().Params.Get(suite.Ctx)
	suite.Require().NoError(err)

	minDepositAmount := proposal.GetMinDepositFromParams(params)
	minDepositAmount, err = suite.GetKeeper().GetMinDepositAmountFromProposalMsgs(suite.Ctx, minDepositAmount, proposal)
	suite.Require().NoError(err)
	return minDepositAmount
}

func (suite *KeeperTestSuite) TestDeposits() {
	initCoins := helpers.NewStakingCoins(10, 18)
	msgContent, err := govv1.NewLegacyContent(govv1beta1.NewTextProposal("Test", "description"), suite.govAcct)
	suite.Require().NoError(err)
	proposal := suite.SubmitProposal(initCoins, suite.newAddress(), msgContent)

	addr := suite.newAddress()
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	params := suite.GetParams()
	suite.True(initCoins.IsAllLT(params.MinDeposit))

	// first deposit
	firstCoins := helpers.NewStakingCoins(10, 18)
	votingStarted, err := suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, addr, firstCoins)
	suite.Require().NoError(err)
	suite.False(votingStarted)

	deposit := suite.GetDeposit(proposal.Id, addr)
	suite.Equal(firstCoins.String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)

	proposal = suite.GetProposal(proposal.Id)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	suite.Equal(firstCoins, sdk.NewCoins(proposal.TotalDeposit...).Sub(initCoins...))
	suite.True(initCoins.Add(firstCoins...).IsAllLT(params.MinDeposit))

	// second deposit
	secondCoins := helpers.NewStakingCoins(90, 18)
	votingStarted, err = suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, addr, secondCoins)
	suite.Require().NoError(err)
	suite.True(votingStarted)

	deposit = suite.GetDeposit(proposal.Id, addr)
	suite.Equal(firstCoins.Add(secondCoins...).String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)

	proposal = suite.GetProposal(proposal.Id)
	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllGTE(params.MinDeposit))
}

func (suite *KeeperTestSuite) TestEGFDepositsLessThan30() {
	suite.addFundCommunityPool()

	initCoins := helpers.NewStakingCoins(30, 18)
	msg := &distributiontypes.MsgCommunityPoolSpend{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Recipient: helpers.GenAccAddress().String(),
		Amount:    helpers.NewStakingCoins(300, 18),
	}
	proposal := suite.SubmitProposal(initCoins, suite.newAddress(), msg)
	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
	suite.Equal(types.DefaultEGFCustomParamVotingPeriod, proposal.VotingEndTime.Sub(*proposal.VotingStartTime))
	minDeposit := suite.GetMinInitialDeposit(proposal)
	suite.Equal(initCoins.String(), minDeposit.String())
}

func (suite *KeeperTestSuite) TestEGFDepositsMoreThan30() {
	suite.addFundCommunityPool()

	initCoins := helpers.NewStakingCoins(30, 18)
	msg := &distributiontypes.MsgCommunityPoolSpend{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Recipient: helpers.GenAccAddress().String(),
		Amount:    helpers.NewStakingCoins(310, 18),
	}
	proposal := suite.SubmitProposal(initCoins, suite.newAddress(), msg)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	minDeposit := suite.GetMinInitialDeposit(proposal)
	suite.True(initCoins.IsAllLT(minDeposit), minDeposit.String())

	depositCoins := helpers.NewStakingCoins(1, 18)
	votingStarted, err := suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, suite.newAddress(), depositCoins)
	suite.Require().NoError(err)
	suite.True(votingStarted)

	proposal = suite.GetProposal(proposal.Id)
	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
	suite.Equal(sdk.NewCoins(proposal.TotalDeposit...).String(), minDeposit.String())
}

func (suite *KeeperTestSuite) TestEGFDeposits() {
	suite.addFundCommunityPool()

	initCoins := helpers.NewStakingCoins(10, 18)
	msg := &distributiontypes.MsgCommunityPoolSpend{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Recipient: helpers.GenAccAddress().String(),
		Amount:    helpers.NewStakingCoins(1500, 18),
	}
	proposal := suite.SubmitProposal(initCoins, suite.newAddress(), msg)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	minDeposit := suite.GetMinInitialDeposit(proposal)
	suite.True(initCoins.IsAllLT(minDeposit))

	// first deposit
	firstCoins := helpers.NewStakingCoins(10, 18)
	addr := suite.newAddress()
	votingStarted, err := suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, addr, firstCoins)
	suite.Require().NoError(err)
	suite.False(votingStarted)

	deposit := suite.GetDeposit(proposal.Id, addr)
	suite.Equal(firstCoins.String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)

	proposal = suite.GetProposal(proposal.Id)
	suite.Require().NoError(err)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	suite.Equal(firstCoins, sdk.NewCoins(proposal.TotalDeposit...).Sub(initCoins...))
	suite.True(initCoins.Add(firstCoins...).IsAllLT(minDeposit))

	// second deposit
	secondCoins := helpers.NewStakingCoins(90, 18)
	votingStarted, err = suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, addr, secondCoins)
	suite.Require().NoError(err)
	suite.False(votingStarted)

	deposit = suite.GetDeposit(proposal.Id, addr)
	suite.Equal(firstCoins.Add(secondCoins...).String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)

	proposal = suite.GetProposal(proposal.Id)
	suite.Require().NoError(err)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllLT(minDeposit))

	// third deposit
	thirdCoins := helpers.NewStakingCoins(40, 18)
	votingStarted, err = suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, addr, thirdCoins)
	suite.Require().NoError(err)
	suite.True(votingStarted)

	deposit = suite.GetDeposit(proposal.Id, addr)
	suite.Equal(firstCoins.Add(secondCoins...).Add(thirdCoins...).String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)

	proposal = suite.GetProposal(proposal.Id)
	suite.Require().NoError(err)
	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).Add(thirdCoins...).IsAllGTE(minDeposit))
}

func (suite *KeeperTestSuite) TestUpdateParams() {
	testCases := []struct {
		testName    string
		amount      sdk.Coins
		msg         []sdk.Msg
		result      bool
		expectedErr string
	}{
		{
			testName:    "set Erc20Params",
			amount:      sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10*1e3).MulRaw(1e18))),
			msg:         []sdk.Msg{&erc20types.MsgUpdateParams{Authority: "0x1", Params: erc20types.DefaultParams()}},
			result:      false,
			expectedErr: "authority",
		},
		{
			testName:    "set CrosschainParam",
			amount:      sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10*1e3).MulRaw(1e18))),
			msg:         []sdk.Msg{&crosschaintypes.MsgUpdateParams{ChainName: "eth", Authority: suite.govAcct, Params: crosschaintypes.DefaultParams()}},
			result:      true,
			expectedErr: "",
		},
		{
			testName:    "set Erc20Params",
			amount:      sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10*1e3).MulRaw(1e18))),
			msg:         []sdk.Msg{&erc20types.MsgUpdateParams{Authority: suite.govAcct, Params: erc20types.DefaultParams()}},
			result:      true,
			expectedErr: "",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.testName), func() {
			proposal, err := suite.GetKeeper().SubmitProposal(suite.Ctx, tc.msg, "", tc.testName, tc.testName, suite.newAddress(), false)
			if tc.result {
				suite.Require().NoError(err)
				_, err = suite.GetKeeper().AddDeposit(suite.Ctx, proposal.Id, suite.newAddress(), tc.amount)
				suite.Require().NoError(err)
				proposal, err := suite.GetKeeper().Keeper.Proposals.Get(suite.Ctx, proposal.Id)
				suite.Require().NoError(err)
				suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
			} else {
				suite.Require().Error(err)
				suite.Require().True(strings.Contains(err.Error(), tc.expectedErr))
			}
		})
	}
}

func TestCheckContractAddressIsDisabled(t *testing.T) {
	addr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	methodId, _ := hex.DecodeString("abcdef01")

	tests := []struct {
		name                string
		disabledPrecompiles []string
		addr                common.Address
		methodId            []byte
		expectedError       string
	}{
		{
			name:                "Empty disabled list",
			disabledPrecompiles: []string{},
			addr:                addr,
			methodId:            methodId,
			expectedError:       "",
		},
		{
			name:                "Address is disabled",
			disabledPrecompiles: []string{addr.Hex()},
			addr:                addr,
			methodId:            methodId,
			expectedError:       "precompile address is disabled",
		},
		{
			name:                "Address and method are disabled",
			disabledPrecompiles: []string{fmt.Sprintf("%s/%s", addr.Hex(), "abcdef01")},
			addr:                addr,
			methodId:            methodId,
			expectedError:       "precompile method abcdef01 is disabled",
		},
		{
			name:                "Address and method are not disabled",
			disabledPrecompiles: []string{fmt.Sprintf("%s/%s", addr.Hex(), "12345678")},
			addr:                addr,
			methodId:            methodId,
			expectedError:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := keeper.CheckContractAddressIsDisabled(tt.disabledPrecompiles, tt.addr, tt.methodId)
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
