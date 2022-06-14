package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	migratekeeper "github.com/functionx/fx-core/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateBank() {
	suite.purseBalance = sdk.NewInt(1000)
	suite.SetupTest()

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := sdk.AccAddress(ethKeys[0].PubKey().Address().Bytes())

	b1 := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc)
	require.NotEmpty(suite.T(), b1)
	b2 := suite.app.BankKeeper.GetAllBalances(suite.ctx, ethAcc)
	require.NotEmpty(suite.T(), b2)

	migrateKeeper := suite.app.MigrateKeeper
	m := migratekeeper.NewBankMigrate(suite.app.BankKeeper)
	err := m.Validate(suite.ctx, migrateKeeper, acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, migrateKeeper, acc, ethAcc)
	suite.Require().NoError(err)

	bb1 := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc)
	suite.Require().Empty(bb1)
	bb2 := suite.app.BankKeeper.GetAllBalances(suite.ctx, ethAcc)
	suite.Require().Equal(b1, bb2.Sub(b2))

}
