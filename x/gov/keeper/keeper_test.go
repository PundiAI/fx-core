package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/gov/keeper"
	"github.com/functionx/fx-core/v3/x/gov/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *app.App
	ctx       sdk.Context
	msgServer govtypes.MsgServer
	valAddr   []sdk.ValAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	valNumber := tmrand.Intn(99) + 1

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address,
	})
	suite.valAddr = make([]sdk.ValAddress, valNumber)
	for i, addr := range valAccounts {
		suite.valAddr[i] = addr.GetAddress().Bytes()
	}
	suite.msgServer = keeper.NewMsgServerImpl(govkeeper.NewMsgServerImpl(suite.app.GovKeeper.Keeper), suite.app.GovKeeper)
}

func (suite *KeeperTestSuite) addFundCommunityPool() {
	sender := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	coin := sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(5 * 1e8).MulRaw(1e18)}
	helpers.AddTestAddr(suite.app, suite.ctx, sender, sdk.NewCoins(coin))
	err := suite.app.DistrKeeper.FundCommunityPool(suite.ctx, sdk.NewCoins(coin), sender)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) newAddress() sdk.AccAddress {
	address := sdk.AccAddress(helpers.GenerateAddress().Bytes())
	coin := sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(50_000).MulRaw(1e18)}
	helpers.AddTestAddr(suite.app, suite.ctx, address, sdk.NewCoins(coin))
	return address
}

func (suite *KeeperTestSuite) TestDeposits() {
	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1e3).MulRaw(1e18)}}
	testProposalMsg, err := govtypes.NewMsgSubmitProposal(
		govtypes.NewTextProposal("Test", "description"),
		initCoins,
		suite.newAddress(),
	)
	suite.NoError(err)
	proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), testProposalMsg)
	suite.NoError(err)
	addr := suite.newAddress()
	_, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposalResponse.ProposalId, addr)
	suite.False(found)
	proposal, ok := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusDepositPeriod, proposal.Status)
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
	suite.True(initCoins.IsAllLT(minDeposit))

	// first deposit
	firstCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1e3).MulRaw(1e18)}}
	votingStarted, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, addr, firstCoins)
	suite.NoError(err)
	suite.False(votingStarted)
	deposit, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, addr)
	suite.True(found)
	suite.Equal(firstCoins, deposit.Amount)
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusDepositPeriod, proposal.Status)
	suite.Equal(firstCoins, proposal.TotalDeposit.Sub(initCoins))
	suite.True(initCoins.Add(firstCoins...).IsAllLT(minDeposit))

	// second deposit
	secondCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(9 * 1e3).MulRaw(1e18)}}
	votingStarted, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, addr, secondCoins)
	suite.NoError(err)
	suite.True(votingStarted)
	deposit, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, addr)
	suite.True(found)
	suite.Equal(firstCoins.Add(secondCoins...), deposit.Amount)
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusVotingPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllGTE(minDeposit))
}

func (suite *KeeperTestSuite) TestEGFDepositsLessThan1000() {
	suite.addFundCommunityPool()

	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)}}
	communityPoolSpendProposal := &distributiontypes.CommunityPoolSpendProposal{
		Title:       "community Pool Spend Proposal",
		Description: "description",
		Recipient:   sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		Amount:      egfCoins,
	}
	minDeposit := types.EGFProposalMinDeposit(egfCoins)
	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	suite.True(initCoins.IsEqual(minDeposit))
	communityPoolSpendProposalMsg, err := govtypes.NewMsgSubmitProposal(
		communityPoolSpendProposal,
		initCoins,
		suite.newAddress(),
	)
	suite.NoError(err)
	proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), communityPoolSpendProposalMsg)
	suite.NoError(err)
	_, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposalResponse.ProposalId, suite.newAddress())
	suite.False(found)
	proposal, ok := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusVotingPeriod, proposal.Status)
}

func (suite *KeeperTestSuite) TestEGFDepositsMoreThan1000() {
	suite.addFundCommunityPool()

	thousand := sdk.NewInt(1 * 1e3).MulRaw(1e18)
	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: thousand.MulRaw(10).Add(sdk.NewInt(10))}}
	communityPoolSpendProposal := &distributiontypes.CommunityPoolSpendProposal{
		Title:       "community Pool Spend Proposal",
		Description: "description",
		Recipient:   sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		Amount:      egfCoins,
	}
	minDeposit := types.EGFProposalMinDeposit(egfCoins)

	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: thousand}}
	communityPoolSpendProposalMsg, err := govtypes.NewMsgSubmitProposal(
		communityPoolSpendProposal,
		initCoins,
		suite.newAddress(),
	)
	suite.NoError(err)
	proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), communityPoolSpendProposalMsg)
	suite.NoError(err)
	_, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposalResponse.ProposalId, suite.newAddress())
	suite.False(found)
	proposal, ok := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusDepositPeriod, proposal.Status)
	suite.True(initCoins.IsAllLT(minDeposit))

	depositCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1)}}
	votingStarted, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.newAddress(), depositCoins)
	suite.NoError(err)
	suite.True(votingStarted)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusVotingPeriod, proposal.Status)
	suite.Equal(proposal.TotalDeposit, minDeposit)
}

func (suite *KeeperTestSuite) TestEGFDeposits() {
	suite.addFundCommunityPool()

	egfCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(150 * 1e3).MulRaw(1e18)}}
	communityPoolSpendProposal := &distributiontypes.CommunityPoolSpendProposal{
		Title:       "community Pool Spend Proposal",
		Description: "description",
		Recipient:   sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		Amount:      egfCoins,
	}
	minDeposit := types.EGFProposalMinDeposit(egfCoins)
	initCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	communityPoolSpendProposalMsg, err := govtypes.NewMsgSubmitProposal(
		communityPoolSpendProposal,
		initCoins,
		suite.newAddress(),
	)
	suite.NoError(err)
	proposalResponse, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), communityPoolSpendProposalMsg)
	suite.NoError(err)
	_, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposalResponse.ProposalId, suite.newAddress())
	suite.False(found)
	proposal, ok := suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusDepositPeriod, proposal.Status)
	suite.True(initCoins.IsAllLT(minDeposit))

	// first deposit
	firstCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(1 * 1e3).MulRaw(1e18)}}
	addr := suite.newAddress()
	votingStarted, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, addr, firstCoins)
	suite.NoError(err)
	suite.False(votingStarted)
	deposit, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, addr)
	suite.True(found)
	suite.Equal(firstCoins, deposit.Amount)
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusDepositPeriod, proposal.Status)
	suite.Equal(firstCoins, proposal.TotalDeposit.Sub(initCoins))
	suite.True(initCoins.Add(firstCoins...).IsAllLT(minDeposit))

	// second deposit
	secondCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(9 * 1e3).MulRaw(1e18)}}
	votingStarted, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, addr, secondCoins)
	suite.NoError(err)
	suite.False(votingStarted)
	deposit, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, addr)
	suite.True(found)
	suite.Equal(firstCoins.Add(secondCoins...), deposit.Amount)
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusDepositPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).IsAllLT(minDeposit))

	// third deposit
	thirdCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(4 * 1e3).MulRaw(1e18)}}
	votingStarted, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, addr, thirdCoins)
	suite.NoError(err)
	suite.True(votingStarted)
	deposit, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, addr)
	suite.True(found)
	suite.Equal(firstCoins.Add(secondCoins...).Add(thirdCoins...), deposit.Amount)
	suite.Equal(addr.String(), deposit.Depositor)
	proposal, ok = suite.app.GovKeeper.Keeper.GetProposal(suite.ctx, proposalResponse.ProposalId)
	suite.True(ok)
	suite.Equal(govtypes.StatusVotingPeriod, proposal.Status)
	suite.True(initCoins.Add(firstCoins...).Add(secondCoins...).Add(thirdCoins...).IsAllGTE(minDeposit))
}
