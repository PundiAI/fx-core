package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	migratekeeper "github.com/functionx/fx-core/v3/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateGovInactive() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdk.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	content := govtypes.ContentFromProposalType("title", "description", "Text")
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000))))

	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
	suite.Require().NoError(err)

	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, acc, amount)
	suite.Require().NoError(err)

	deposit1, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, acc)
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit1.Amount)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes())
	suite.Require().False(found)

	migrateKeeper := suite.app.MigrateKeeper
	m := migratekeeper.NewGovMigrate(suite.app.GetKey(govtypes.StoreKey), suite.app.GovKeeper)
	err = m.Validate(suite.ctx, migrateKeeper, acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, migrateKeeper, acc, ethAcc)
	suite.Require().NoError(err)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, acc)
	suite.Require().False(found)

	deposit2, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit2.Amount)
}

func (suite *KeeperTestSuite) TestMigrateGovActive() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdk.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(2)
	suite.Require().Equal(len(ethKeys), 2)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())
	toEthAcc := common.BytesToAddress(ethKeys[1].PubKey().Address().Bytes())

	content := govtypes.ContentFromProposalType("title", "description", "Text")
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(5000))))

	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
	suite.Require().NoError(err)

	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, acc, amount)
	suite.Require().NoError(err)

	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes(), amount)
	suite.Require().NoError(err)

	deposit1, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, acc)
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit1.Amount)

	deposit2, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit2.Amount)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, toEthAcc.Bytes())
	suite.Require().False(found)

	migrateKeeper := suite.app.MigrateKeeper
	m := migratekeeper.NewGovMigrate(suite.app.GetKey(govtypes.StoreKey), suite.app.GovKeeper)
	err = m.Validate(suite.ctx, migrateKeeper, acc, toEthAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, migrateKeeper, acc, toEthAcc)
	suite.Require().NoError(err)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, acc)
	suite.Require().False(found)

	deposit3, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit3.Amount)

	deposit4, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, toEthAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit4.Amount)
}

func (suite *KeeperTestSuite) TestMigrateGovActiveAndVote() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdk.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(2)
	suite.Require().Equal(len(ethKeys), 2)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())
	toEthAcc := common.BytesToAddress(ethKeys[1].PubKey().Address().Bytes())

	//add proposal
	content := govtypes.ContentFromProposalType("title", "description", "Text")
	amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(5000))))
	proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, content)
	suite.Require().NoError(err)

	//acc deposit
	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, acc, amount)
	suite.Require().NoError(err)

	//eth acc deposit
	_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes(), amount)
	suite.Require().NoError(err)

	// acc vote
	err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, acc, govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
	suite.Require().NoError(err)

	//check acc deposit
	deposit1, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, acc)
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit1.Amount)
	//check eth acc deposit
	deposit2, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit2.Amount)

	//check acc vote
	vote, found := suite.app.GovKeeper.GetVote(suite.ctx, proposal.ProposalId, acc)
	suite.Require().True(found)
	// nolint
	suite.Require().Equal(vote.Option, govtypes.OptionYes)

	//check to address deposit vote
	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, toEthAcc.Bytes())
	suite.Require().False(found)

	_, found = suite.app.GovKeeper.GetVote(suite.ctx, proposal.ProposalId, toEthAcc.Bytes())
	suite.Require().False(found)

	migrateKeeper := suite.app.MigrateKeeper
	m := migratekeeper.NewGovMigrate(suite.app.GetKey(govtypes.StoreKey), suite.app.GovKeeper)
	err = m.Validate(suite.ctx, migrateKeeper, acc, toEthAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, migrateKeeper, acc, toEthAcc)
	suite.Require().NoError(err)

	_, found = suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, acc)
	suite.Require().False(found)

	deposit3, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit3.Amount)

	deposit4, found := suite.app.GovKeeper.GetDeposit(suite.ctx, proposal.ProposalId, toEthAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(amount, deposit4.Amount)

	_, found = suite.app.GovKeeper.GetVote(suite.ctx, proposal.ProposalId, acc)
	suite.Require().False(found)

	vote, found = suite.app.GovKeeper.GetVote(suite.ctx, proposal.ProposalId, toEthAcc.Bytes())
	suite.Require().True(found)
	// nolint
	suite.Require().Equal(vote.Option, govtypes.OptionYes)

}
