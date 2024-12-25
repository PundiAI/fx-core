package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/migrate/types"
)

func (suite *KeeperTestSuite) TestMigrateAccount() {
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
	suite.Require().NotEmpty(b1)

	_, found := suite.App.MigrateKeeper.GetMigrateRecord(suite.Ctx, acc)
	suite.Require().False(found)

	_, found = suite.App.MigrateKeeper.GetMigrateRecord(suite.Ctx, ethAcc.Bytes())
	suite.Require().False(found)

	found = suite.App.MigrateKeeper.HasMigratedDirectionFrom(suite.Ctx, acc)
	suite.Require().False(found)

	found = suite.App.MigrateKeeper.HasMigratedDirectionTo(suite.Ctx, ethAcc)
	suite.Require().False(found)

	_, err := suite.App.MigrateKeeper.MigrateAccount(suite.Ctx, &types.MsgMigrateAccount{
		From:      acc.String(),
		To:        ethAcc.String(),
		Signature: "",
	})
	suite.Require().NoError(err)

	record, found := suite.App.MigrateKeeper.GetMigrateRecord(suite.Ctx, acc)
	suite.Require().True(found)
	suite.Require().Equal(record.From, acc.String())

	record, found = suite.App.MigrateKeeper.GetMigrateRecord(suite.Ctx, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(record.To, ethAcc.String())

	found = suite.App.MigrateKeeper.HasMigratedDirectionFrom(suite.Ctx, acc)
	suite.Require().True(found)

	found = suite.App.MigrateKeeper.HasMigratedDirectionTo(suite.Ctx, ethAcc)
	suite.Require().True(found)

	bb1 := suite.App.BankKeeper.GetAllBalances(suite.Ctx, acc)
	suite.Require().True(bb1.Empty())
	bb2 := suite.App.BankKeeper.GetAllBalances(suite.Ctx, ethAcc.Bytes())
	suite.Require().Equal(b1, bb2.Sub(b2...))

	// expect 1 record
	recordCount := 0
	suite.App.MigrateKeeper.IterateMigrateRecords(suite.Ctx, func(record types.MigrateRecord) bool {
		suite.Require().Equal(record.From, acc.String())
		suite.Require().Equal(record.To, ethAcc.String())
		recordCount++
		return false
	})
	suite.Require().Equal(recordCount, 1)
}
