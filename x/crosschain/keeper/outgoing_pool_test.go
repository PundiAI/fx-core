package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

func (suite *KeeperTestSuite) TestKeeper_OutgoingPool() {
	sender := helpers.GenHexAddress().Bytes()
	bridgeToken := helpers.GenHexAddress().Hex()
	denom := types.NewBridgeDenom(suite.chainName, bridgeToken)
	suite.Equal(sdk.NewCoin(denom, sdkmath.ZeroInt()), suite.App.BankKeeper.GetSupply(suite.Ctx, denom))

	sendAmount := sdk.NewCoin(denom, sdkmath.NewInt(int64(tmrand.Uint32()*2)))
	err := suite.App.BankKeeper.MintCoins(suite.Ctx, suite.chainName, sdk.NewCoins(sendAmount))
	suite.NoError(err)
	err = suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, suite.chainName, sender, sdk.NewCoins(sendAmount))
	suite.NoError(err)
	suite.Equal(sendAmount, suite.App.BankKeeper.GetSupply(suite.Ctx, denom))

	suite.Keeper().AddBridgeToken(suite.Ctx, bridgeToken, denom)

	suite.Equal(suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender).AmountOf(denom).String(), sendAmount.Amount.String())
	receiver := helpers.GenHexAddress().Hex()
	amount := sdk.NewCoin(denom, sendAmount.Amount.QuoRaw(2))
	txId, err := suite.Keeper().AddToOutgoingPool(suite.Ctx, sender, receiver, amount, amount)
	suite.NoError(err)
	suite.Equal(suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender).AmountOf(denom).String(), sdkmath.NewInt(0).String())

	suite.Equal(sdk.NewCoin(denom, sdkmath.ZeroInt()), suite.App.BankKeeper.GetSupply(suite.Ctx, denom))

	_, err = suite.Keeper().RemoveFromOutgoingPoolAndRefund(suite.Ctx, txId, sender)
	suite.NoError(err)
	suite.Equal(suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender).AmountOf(denom).String(), sendAmount.Amount.String())

	suite.Equal(sendAmount, suite.App.BankKeeper.GetSupply(suite.Ctx, denom))
}

func (suite *KeeperTestSuite) TestKeeper_OutgoingPool2() {
	sender := helpers.GenHexAddress().Bytes()
	bridgeToken := helpers.GenHexAddress().Hex()
	denom := fxtypes.DefaultDenom
	supply := suite.App.BankKeeper.GetSupply(suite.Ctx, denom)

	sendAmount := sdk.NewCoin(denom, sdkmath.NewInt(int64(tmrand.Uint32()*2)))
	err := suite.App.BankKeeper.SendCoinsFromModuleToAccount(suite.Ctx, suite.chainName, sender, sdk.NewCoins(sendAmount))
	if suite.chainName != ethtypes.ModuleName {
		suite.Error(err)
		return
	}
	suite.NoError(err)

	suite.Keeper().AddBridgeToken(suite.Ctx, bridgeToken, denom)

	suite.Equal(suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender).AmountOf(denom).String(), sendAmount.Amount.String())
	receiver := helpers.GenHexAddress().Hex()
	amount := sdk.NewCoin(denom, sendAmount.Amount.QuoRaw(2))
	txId, err := suite.Keeper().AddToOutgoingPool(suite.Ctx, sender, receiver, amount, amount)
	suite.NoError(err)
	suite.Equal(suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender).AmountOf(denom).String(), sdkmath.NewInt(0).String())

	suite.Equal(supply, suite.App.BankKeeper.GetSupply(suite.Ctx, denom))

	_, err = suite.Keeper().RemoveFromOutgoingPoolAndRefund(suite.Ctx, txId, sender)
	suite.NoError(err)
	suite.Equal(suite.App.BankKeeper.GetAllBalances(suite.Ctx, sender).AmountOf(denom).String(), sendAmount.Amount.String())

	suite.Equal(supply, suite.App.BankKeeper.GetSupply(suite.Ctx, denom))
}
