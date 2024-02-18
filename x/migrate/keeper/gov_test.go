package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
	migratekeeper "github.com/functionx/fx-core/v7/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateGovInactive() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	content, bl := govv1beta1.ContentFromProposalType("title", "description", "Text")
	suite.Require().True(bl)
	legacyContent, err := govv1.NewLegacyContent(content, suite.govAddr)
	suite.Require().NoError(err)

	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(1000))))

	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, []sdk.Msg{legacyContent},
		fxgovtypes.NewFXMetadata(content.GetTitle(), content.GetDescription(), "").String())
	suite.Require().NoError(err)

	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.Id, acc, amount)
	suite.Require().NoError(err)

	deposit1, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, acc)
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit1.Amount...))

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, ethAcc.Bytes())
	suite.Require().False(found)

	m := migratekeeper.NewGovMigrate(suite.app.GetKey(govtypes.StoreKey), suite.app.GovKeeper)
	err = m.Validate(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, suite.app.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, acc)
	suite.Require().False(found)

	deposit2, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit2.Amount...))
}

func (suite *KeeperTestSuite) TestMigrateGovActive() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(2)
	suite.Require().Equal(len(ethKeys), 2)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())
	toEthAcc := common.BytesToAddress(ethKeys[1].PubKey().Address().Bytes())

	content, bl := govv1beta1.ContentFromProposalType("title", "description", "Text")
	suite.Require().True(bl)
	legacyContent, err := govv1.NewLegacyContent(content, suite.govAddr)
	suite.Require().NoError(err)

	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(5000))))

	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, []sdk.Msg{legacyContent},
		fxgovtypes.NewFXMetadata(content.GetTitle(), content.GetDescription(), "").String())
	suite.Require().NoError(err)

	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.Id, acc, amount)
	suite.Require().NoError(err)

	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.Id, ethAcc.Bytes(), amount)
	suite.Require().NoError(err)

	deposit1, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, acc)
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit1.Amount...))

	deposit2, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit2.Amount...))

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, toEthAcc.Bytes())
	suite.Require().False(found)

	m := migratekeeper.NewGovMigrate(suite.app.GetKey(govtypes.StoreKey), suite.app.GovKeeper)
	err = m.Validate(suite.ctx, suite.app.AppCodec(), acc, toEthAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, suite.app.AppCodec(), acc, toEthAcc)
	suite.Require().NoError(err)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, acc)
	suite.Require().False(found)

	deposit3, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit3.Amount...))

	deposit4, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, toEthAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit4.Amount...))
}

func (suite *KeeperTestSuite) TestMigrateGovActiveAndVote() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(2)
	suite.Require().Equal(len(ethKeys), 2)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())
	toEthAcc := common.BytesToAddress(ethKeys[1].PubKey().Address().Bytes())

	// add proposal
	content, bl := govv1beta1.ContentFromProposalType("title", "description", "Text")
	suite.Require().True(bl)
	legacyContent, err := govv1.NewLegacyContent(content, suite.govAddr)
	suite.Require().NoError(err)

	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromUint64(1e18).Mul(sdkmath.NewInt(5000))))
	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, []sdk.Msg{legacyContent},
		fxgovtypes.NewFXMetadata(content.GetTitle(), content.GetDescription(), "").String())
	suite.Require().NoError(err)

	// acc deposit
	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.Id, acc, amount)
	suite.Require().NoError(err)

	// eth acc deposit
	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.Id, ethAcc.Bytes(), amount)
	suite.Require().NoError(err)

	// acc vote
	err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.Id, acc, govv1.NewNonSplitVoteOption(govv1.OptionYes), "")
	suite.Require().NoError(err)

	// check acc deposit
	deposit1, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, acc)
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit1.Amount...))
	// check eth acc deposit
	deposit2, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit2.Amount...))

	// check acc vote
	vote, found := suite.app.GovKeeper.GetVote(suite.ctx, proposal.Id, acc)
	suite.Require().True(found)
	suite.Require().Equal(vote.Options, []*govv1.WeightedVoteOption{{Option: govv1.VoteOption_VOTE_OPTION_YES, Weight: "1.000000000000000000"}})

	// check to address deposit vote
	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, toEthAcc.Bytes())
	suite.Require().False(found)

	_, found = suite.app.GovKeeper.GetVote(suite.ctx, proposal.Id, toEthAcc.Bytes())
	suite.Require().False(found)

	m := migratekeeper.NewGovMigrate(suite.app.GetKey(govtypes.StoreKey), suite.app.GovKeeper)
	err = m.Validate(suite.ctx, suite.app.AppCodec(), acc, toEthAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, suite.app.AppCodec(), acc, toEthAcc)
	suite.Require().NoError(err)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, acc)
	suite.Require().False(found)

	deposit3, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit3.Amount...))

	deposit4, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.Id, toEthAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, sdk.NewCoins(deposit4.Amount...))

	_, found = suite.app.GovKeeper.GetVote(suite.ctx, proposal.Id, acc)
	suite.Require().False(found)

	vote, found = suite.app.GovKeeper.GetVote(suite.ctx, proposal.Id, toEthAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(vote.Options, []*govv1.WeightedVoteOption{{Option: govv1.VoteOption_VOTE_OPTION_YES, Weight: "1.000000000000000000"}})
}
