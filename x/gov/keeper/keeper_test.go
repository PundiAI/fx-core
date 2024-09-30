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
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	evmtypes "github.com/functionx/fx-core/v8/x/evm/types"
	"github.com/functionx/fx-core/v8/x/gov/keeper"
	"github.com/functionx/fx-core/v8/x/gov/types"
)

type KeeperTestSuite struct {
	helpers.BaseSuite

	govAcct     string
	msgServer   types.MsgServerPro
	queryClient types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.BaseSuite.MintValNumber = 1
	suite.BaseSuite.SetupTest()

	suite.govAcct = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	suite.msgServer = keeper.NewMsgServerImpl(govkeeper.NewMsgServerImpl(suite.App.GovKeeper.Keeper), suite.App.GovKeeper)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.App.GovKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *KeeperTestSuite) addFundCommunityPool() {
	sender := suite.AddTestSigner(5 * 1e8)
	balances := suite.Balance(sender.AccAddress())

	err := suite.App.DistrKeeper.FundCommunityPool(suite.Ctx, balances, sender.AccAddress())
	suite.NoError(err)
}

func (suite *KeeperTestSuite) newAddress() sdk.AccAddress {
	return suite.AddTestSigner(50_000).AccAddress()
}

func (suite *KeeperTestSuite) TestDeposits() {
	initCoins, testProposalMsg, err := suite.getTextProposal()
	suite.NoError(err)
	proposalResponse, err := suite.msgServer.SubmitProposal(suite.Ctx, testProposalMsg)
	suite.NoError(err)
	addr := suite.newAddress()
	_, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposalResponse.ProposalId, addr))
	suite.Error(err)
	proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
	suite.NoError(err)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	params, err := suite.App.GovKeeper.Keeper.Params.Get(suite.Ctx)
	suite.NoError(err)
	suite.True(initCoins.IsAllLT(params.MinDeposit))

	// first deposit
	firstCoins := helpers.NewStakingCoins(1000, 18)
	votingStarted, err := suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, addr, firstCoins)
	suite.NoError(err)
	suite.False(votingStarted)
	deposit, err := suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposal.Id, addr))
	suite.NoError(err)
	suite.Equal(firstCoins.String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
	suite.NoError(err)
	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
	suite.Equal(firstCoins, sdk.NewCoins(proposal.TotalDeposit...).Sub(initCoins...))
	suite.True(initCoins.Add(firstCoins...).IsAllLT(params.MinDeposit))

	// second deposit
	secondCoins := helpers.NewStakingCoins(9_000, 18)
	votingStarted, err = suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, addr, secondCoins)
	suite.NoError(err)
	suite.True(votingStarted)
	deposit, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposal.Id, addr))
	suite.NoError(err)
	suite.Equal(firstCoins.Add(secondCoins...).String(), sdk.NewCoins(deposit.Amount...).String())
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
	suite.NoError(err)
	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllGTE(params.MinDeposit))
}

func (suite *KeeperTestSuite) getTextProposal() (sdk.Coins, *govv1.MsgSubmitProposal, error) {
	initCoins := helpers.NewStakingCoins(1000, 18)
	content := govv1beta1.NewTextProposal("Test", "description")
	msgExecLegacyContent, err := govv1.NewLegacyContent(content, suite.govAcct)
	suite.NoError(err)
	testProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{msgExecLegacyContent}, initCoins, suite.newAddress().String(),
		"", content.GetTitle(), content.GetDescription(), false)
	return initCoins, testProposalMsg, err
}

// func (suite *KeeperTestSuite) TestEGFDepositsLessThan1000() {
// 	suite.addFundCommunityPool()
//
// 	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(10 * 1e3).MulRaw(1e18)}}
//
// 	spendProposal := &distributiontypes.MsgCommunityPoolSpend{Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(), Recipient: helpers.GenAccAddress().String(), Amount: egfCoins}
// 	minDeposit := suite.App.GovKeeper.EGFProposalMinDeposit(suite.Ctx, sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}), egfCoins)
// 	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}
// 	suite.True(initCoins.Equal(minDeposit))
//
// 	communityPoolSpendProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal}, initCoins, suite.newAddress().String(),
// 		"", "community Pool Spend Proposal", "description", false)
// 	suite.NoError(err)
// 	proposalResponse, err := suite.msgServer.SubmitProposal(suite.Ctx, communityPoolSpendProposalMsg)
// 	suite.NoError(err)
// 	_, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposalResponse.ProposalId, suite.newAddress()))
// 	suite.Error(err)
// 	proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
// }

// func (suite *KeeperTestSuite) TestEGFDepositsMoreThan1000() {
// 	suite.addFundCommunityPool()
//
// 	thousand := sdkmath.NewInt(1 * 1e3).MulRaw(1e18)
// 	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: thousand.MulRaw(10).Add(sdkmath.NewInt(10))}}
//
// 	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: thousand}}
// 	spendProposal := &distributiontypes.MsgCommunityPoolSpend{Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(), Recipient: helpers.GenAccAddress().String(), Amount: egfCoins}
// 	minDeposit := suite.App.GovKeeper.EGFProposalMinDeposit(suite.Ctx, sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}), egfCoins)
//
// 	communityPoolSpendProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal}, initCoins, suite.newAddress().String(),
// 		"", "community Pool Spend Proposal", "description", false)
// 	suite.NoError(err)
// 	proposalResponse, err := suite.msgServer.SubmitProposal(suite.Ctx, communityPoolSpendProposalMsg)
// 	suite.NoError(err)
// 	_, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposalResponse.ProposalId, suite.newAddress()))
// 	suite.Error(err)
// 	proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
// 	suite.True(initCoins.IsAllLT(minDeposit))
//
// 	depositCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1)}}
// 	votingStarted, err := suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, suite.newAddress(), depositCoins)
// 	suite.NoError(err)
// 	suite.True(votingStarted)
// 	proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
// 	suite.Equal(sdk.NewCoins(proposal.TotalDeposit...).String(), minDeposit.String())
// }

// func (suite *KeeperTestSuite) TestEGFDeposits() {
// 	suite.addFundCommunityPool()
//
// 	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(150 * 1e3).MulRaw(1e18)}}
//
// 	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}
// 	spendProposal := &distributiontypes.MsgCommunityPoolSpend{Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(), Recipient: helpers.GenAccAddress().String(), Amount: egfCoins}
// 	minDeposit := suite.App.GovKeeper.EGFProposalMinDeposit(suite.Ctx, sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}), egfCoins)
//
// 	communityPoolSpendProposalMsg, err := govv1.NewMsgSubmitProposal([]sdk.Msg{spendProposal}, initCoins, suite.newAddress().String(),
// 		"", "community Pool Spend Proposal", "description", false)
// 	suite.NoError(err)
// 	proposalResponse, err := suite.msgServer.SubmitProposal(suite.Ctx, communityPoolSpendProposalMsg)
// 	suite.NoError(err)
// 	_, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposalResponse.ProposalId, suite.newAddress()))
// 	suite.Error(err)
// 	proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
// 	suite.True(initCoins.IsAllLT(minDeposit))
//
// 	// first deposit
// 	firstCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(1 * 1e3).MulRaw(1e18)}}
// 	addr := suite.newAddress()
// 	votingStarted, err := suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, addr, firstCoins)
// 	suite.NoError(err)
// 	suite.False(votingStarted)
// 	deposit, err := suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposal.Id, addr))
// 	suite.NoError(err)
// 	suite.Equal(firstCoins.String(), sdk.NewCoins(deposit.Amount...).String())
// 	suite.Equal(addr.String(), deposit.Depositor)
// 	proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
// 	suite.Equal(firstCoins, sdk.NewCoins(proposal.TotalDeposit...).Sub(initCoins...))
// 	suite.True(initCoins.Add(firstCoins...).IsAllLT(minDeposit))
//
// 	// second deposit
// 	secondCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(9 * 1e3).MulRaw(1e18)}}
// 	votingStarted, err = suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, addr, secondCoins)
// 	suite.NoError(err)
// 	suite.False(votingStarted)
// 	deposit, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposal.Id, addr))
// 	suite.NoError(err)
// 	suite.Equal(firstCoins.Add(secondCoins...).String(), sdk.NewCoins(deposit.Amount...).String())
// 	suite.Equal(addr.String(), deposit.Depositor)
// 	proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusDepositPeriod, proposal.Status)
// 	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllLT(minDeposit))
//
// 	// third deposit
// 	thirdCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdkmath.NewInt(4 * 1e3).MulRaw(1e18)}}
// 	votingStarted, err = suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, addr, thirdCoins)
// 	suite.NoError(err)
// 	suite.True(votingStarted)
// 	deposit, err = suite.App.GovKeeper.Deposits.Get(suite.Ctx, collections.Join(proposal.Id, addr))
// 	suite.NoError(err)
// 	suite.Equal(firstCoins.Add(secondCoins...).Add(thirdCoins...).String(), sdk.NewCoins(deposit.Amount...).String())
// 	suite.Equal(addr.String(), deposit.Depositor)
// 	proposal, err = suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposalResponse.ProposalId)
// 	suite.NoError(err)
// 	suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
// 	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).Add(thirdCoins...).IsAllGTE(minDeposit))
// }

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
			testName:    "set CrossChainParam",
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
			proposal, err := suite.App.GovKeeper.SubmitProposal(suite.Ctx, tc.msg, "", tc.testName, tc.testName, suite.newAddress(), false)
			if tc.result {
				suite.NoError(err)
				_, err = suite.App.GovKeeper.AddDeposit(suite.Ctx, proposal.Id, suite.newAddress(), tc.amount)
				suite.Require().NoError(err)
				proposal, err := suite.App.GovKeeper.Keeper.Proposals.Get(suite.Ctx, proposal.Id)
				suite.Require().NoError(err)
				suite.Equal(govv1.StatusVotingPeriod, proposal.Status)
			} else {
				suite.Error(err, err)
				suite.Require().True(strings.Contains(err.Error(), tc.expectedErr))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGetParams() {
	govParams, err := suite.App.GovKeeper.Keeper.Params.Get(suite.Ctx)
	suite.NoError(err)

	tallyParamsVetoThreshold := govParams.VetoThreshold
	tallyParamsThreshold := govParams.Threshold
	tallyParamsQuorum := govParams.Quorum

	depositParamsMinDeposit := govParams.MinDeposit
	depositParamsMaxDepositPeriod := govParams.MaxDepositPeriod
	votingParamsVotingPeriod := govParams.VotingPeriod

	// unregistered Msg
	params := suite.App.GovKeeper.GetFXParams(suite.Ctx, sdk.MsgTypeURL(&erc20types.MsgUpdateParams{}))
	suite.Require().EqualValues(params.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultMinInitialDeposit).String())
	suite.Require().EqualValues(params.MinDeposit, depositParamsMinDeposit)
	suite.Require().EqualValues(params.MaxDepositPeriod, depositParamsMaxDepositPeriod)
	suite.Require().EqualValues(params.VotingPeriod, votingParamsVotingPeriod)
	suite.Require().EqualValues(params.VetoThreshold, tallyParamsVetoThreshold)
	suite.Require().EqualValues(params.Threshold, tallyParamsThreshold)
	suite.Require().EqualValues(params.Quorum, tallyParamsQuorum)
	erc20MsgType := []string{
		"/fx.erc20.v1.RegisterCoinProposal",
		"/fx.erc20.v1.RegisterERC20Proposal",
		"/fx.erc20.v1.ToggleTokenConversionProposal",
		"/fx.erc20.v1.UpdateDenomAliasProposal",
		sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
		sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{}),
		sdk.MsgTypeURL(&erc20types.MsgToggleTokenConversion{}),
		sdk.MsgTypeURL(&erc20types.MsgUpdateDenomAlias{}),
	}
	for _, erc20MsgType := range erc20MsgType {
		// registered Msg
		params = suite.App.GovKeeper.GetFXParams(suite.Ctx, erc20MsgType)
		suite.Require().NoError(params.ValidateBasic())
		suite.Require().EqualValues(params.MinDeposit, depositParamsMinDeposit)
		suite.Require().EqualValues(params.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultMinInitialDeposit).String())
		suite.Require().EqualValues(params.MaxDepositPeriod, depositParamsMaxDepositPeriod)
		suite.Require().EqualValues(params.VotingPeriod, votingParamsVotingPeriod)
		suite.Require().EqualValues(params.VetoThreshold, tallyParamsVetoThreshold)
		suite.Require().EqualValues(params.Threshold, tallyParamsThreshold)
		suite.Require().EqualValues(params.Quorum, types.DefaultErc20Quorum.String())
	}

	params = suite.App.GovKeeper.GetFXParams(suite.Ctx, sdk.MsgTypeURL(&evmtypes.MsgCallContract{}))
	suite.Require().NoError(params.ValidateBasic())
	suite.Require().EqualValues(params.MinDeposit, depositParamsMinDeposit)
	suite.Require().EqualValues(params.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultMinInitialDeposit).String())
	suite.Require().EqualValues(params.MaxDepositPeriod, depositParamsMaxDepositPeriod)
	suite.Require().EqualValues(params.VotingPeriod.String(), types.DefaultEvmVotingPeriod.String())
	suite.Require().EqualValues(params.VetoThreshold, tallyParamsVetoThreshold)
	suite.Require().EqualValues(params.Threshold, tallyParamsThreshold)
	suite.Require().EqualValues(params.Quorum, types.DefaultEvmQuorum.String())

	params = suite.App.GovKeeper.GetFXParams(suite.Ctx, "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal")
	egfParams := suite.App.GovKeeper.GetEGFParams(suite.Ctx)
	suite.Require().NoError(params.ValidateBasic())
	suite.Require().NoError(egfParams.ValidateBasic())
	suite.Require().EqualValues(params.MinDeposit, depositParamsMinDeposit)
	suite.Require().EqualValues(params.MinInitialDeposit.String(), sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultMinInitialDeposit).String())
	suite.Require().EqualValues(params.MaxDepositPeriod, depositParamsMaxDepositPeriod)
	suite.Require().EqualValues(params.VotingPeriod.String(), types.DefaultEgfVotingPeriod.String())
	suite.Require().EqualValues(params.VetoThreshold, tallyParamsVetoThreshold)
	suite.Require().EqualValues(params.Threshold, tallyParamsThreshold)
	suite.Require().EqualValues(params.Quorum, tallyParamsQuorum)
	suite.Require().EqualValues(egfParams.EgfDepositThreshold, sdk.NewCoin(fxtypes.DefaultDenom, types.DefaultEgfDepositThreshold))
	suite.Require().EqualValues(egfParams.ClaimRatio, types.DefaultClaimRatio.String())
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
