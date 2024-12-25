package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	migratekeeper "github.com/pundiai/fx-core/v8/x/migrate/keeper"
)

func (suite *KeeperTestSuite) TestMigrateBank() {
	suite.MintToken(suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	b1 := suite.App.BankKeeper.GetAllBalances(suite.Ctx, acc)
	suite.Require().NotEmpty(b1)
	b2 := suite.App.BankKeeper.GetAllBalances(suite.Ctx, ethAcc.Bytes())
	suite.Require().NotEmpty(b2)

	m := migratekeeper.NewBankMigrate(suite.App.BankKeeper)
	err := m.Validate(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)
	err = m.Execute(suite.Ctx, suite.App.AppCodec(), acc, ethAcc)
	suite.Require().NoError(err)

	bb1 := suite.App.BankKeeper.GetAllBalances(suite.Ctx, acc)
	suite.Require().Empty(bb1)
	bb2 := suite.App.BankKeeper.GetAllBalances(suite.Ctx, ethAcc.Bytes())
	suite.Require().Equal(b1, bb2.Sub(b2...))
}
