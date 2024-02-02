package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	bsctypes "github.com/functionx/fx-core/v7/x/bsc/types"
	"github.com/functionx/fx-core/v7/x/migrate/types"
)

func (suite *KeeperTestSuite) TestMigrateAccount() {
	suite.mintToken(bsctypes.ModuleName, suite.secp256k1PrivKey.PubKey().Address().Bytes(), sdk.NewCoin("ibc/ABC", sdkmath.NewInt(1000)))

	keys := suite.GenerateAcc(1)
	suite.Require().Equal(len(keys), 1)
	acc := sdk.AccAddress(keys[0].PubKey().Address().Bytes())
	ethKeys := suite.GenerateEthAcc(1)
	suite.Require().Equal(len(ethKeys), 1)
	ethAcc := common.BytesToAddress(ethKeys[0].PubKey().Address().Bytes())

	b1 := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc)
	suite.Require().NotEmpty(b1)

	b2 := suite.app.BankKeeper.GetAllBalances(suite.ctx, ethAcc.Bytes())
	suite.Require().NotEmpty(b1)

	_, found := suite.app.MigrateKeeper.GetMigrateRecord(suite.ctx, acc)
	suite.Require().False(found)

	_, found = suite.app.MigrateKeeper.GetMigrateRecord(suite.ctx, ethAcc.Bytes())
	suite.Require().False(found)

	found = suite.app.MigrateKeeper.HasMigratedDirectionFrom(suite.ctx, acc)
	suite.Require().False(found)

	found = suite.app.MigrateKeeper.HasMigratedDirectionTo(suite.ctx, ethAcc)
	suite.Require().False(found)

	_, err := suite.app.MigrateKeeper.MigrateAccount(sdk.WrapSDKContext(suite.ctx), &types.MsgMigrateAccount{
		From:      acc.String(),
		To:        ethAcc.String(),
		Signature: "",
	})
	suite.Require().NoError(err)

	record, found := suite.app.MigrateKeeper.GetMigrateRecord(suite.ctx, acc)
	suite.Require().True(found)
	suite.Require().Equal(record.From, acc.String())

	record, found = suite.app.MigrateKeeper.GetMigrateRecord(suite.ctx, ethAcc.Bytes())
	suite.Require().True(found)
	suite.Require().Equal(record.To, ethAcc.String())

	found = suite.app.MigrateKeeper.HasMigratedDirectionFrom(suite.ctx, acc)
	suite.Require().True(found)

	found = suite.app.MigrateKeeper.HasMigratedDirectionTo(suite.ctx, ethAcc)
	suite.Require().True(found)

	bb1 := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc)
	suite.Require().True(bb1.Empty())
	bb2 := suite.app.BankKeeper.GetAllBalances(suite.ctx, ethAcc.Bytes())
	suite.Require().Equal(b1, bb2.Sub(b2...))
}
