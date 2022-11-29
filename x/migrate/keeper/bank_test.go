package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	migratekeeper "github.com/functionx/fx-core/v3/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateBank() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdk.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	b1 := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc)
	require.NotEmpty(suite.T(), b1)
	b2 := suite.app.BankKeeper.GetAllBalances(suite.ctx, ethAcc.Bytes())
	require.NotEmpty(suite.T(), b2)

	migrateKeeper := suite.app.MigrateKeeper
	m := migratekeeper.NewBankMigrate(suite.app.BankKeeper)
	err := m.Validate(suite.ctx, migrateKeeper, acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.ctx, migrateKeeper, acc, ethAcc)
	suite.Require().NoError(err)

	bb1 := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc)
	suite.Require().Empty(bb1)
	bb2 := suite.app.BankKeeper.GetAllBalances(suite.ctx, ethAcc.Bytes())
	suite.Require().Equal(b1, bb2.Sub(b2))

}
