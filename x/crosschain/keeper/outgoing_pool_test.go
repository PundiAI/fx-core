package keeper_test

import (
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
)

// todo need test FX token

func (suite *KeeperTestSuite) TestKeeper_OutgoingPool() {
	sender := helpers.GenHexAddress().Bytes()
	amount := helpers.NewRandAmount()
	fee := helpers.NewRandAmount()
	baseDenom, bridgeDenom, _ := suite.AddRandomBaseToken(false)
	suite.MintBaseToken(sender, baseDenom, bridgeDenom, amount.Add(fee))

	// 1. test add to outgoing pool
	txId, err := suite.Keeper().AddToOutgoingPool(suite.Ctx, sender, tmrand.Str(20), sdk.NewCoin(baseDenom, amount), sdk.NewCoin(baseDenom, fee))
	suite.Require().NoError(err)
	suite.Require().Greater(txId, uint64(0))
	suite.CheckAllBalance(sender, sdk.NewCoins()...)

	// 2. test cancel outgoing pool
	_, err = suite.Keeper().RemoveFromOutgoingPoolAndRefund(suite.Ctx, txId, sender)
	suite.NoError(err)
	suite.CheckAllBalance(sender, sdk.NewCoin(baseDenom, amount.Add(fee)))
}
