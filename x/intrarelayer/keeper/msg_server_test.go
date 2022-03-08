package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

func (suite *KeeperTestSuite) TestConvertCoinNativeCoin() {
	testCases := []struct {
		name     string
		mint     int64
		burn     int64
		malleate func()
		expPass  bool
	}{
		{"ok - sufficient funds", 100, 10, func() {}, true},
		{"ok - equal funds", 10, 10, func() {}, true},
		{"fail - insufficient funds", 0, 10, func() {}, false},
		{
			"fail - minting disabled",
			100,
			10,
			func() {
				params := types.DefaultParams()
				params.EnableIntrarelayer = false
				suite.app.IntrarelayerKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)

			tc.malleate()
			suite.Commit()

			ctx := sdk.WrapSDKContext(suite.ctx)
			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.address.Bytes())
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.burn)),
				suite.address,
				sender,
			)

			suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins)

			res, err := suite.app.IntrarelayerKeeper.ConvertCoin(ctx, msg)
			expRes := &types.MsgConvertCoinResponse{}
			suite.Commit()
			balance := suite.BalanceOf(common.HexToAddress(pair.Fip20Address), suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, metadata.Base)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expRes, res)
				suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
				suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.burn).Int64())

			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
}

func (suite *KeeperTestSuite) TestConvertECR20NativeCoin() {
	testCases := []struct {
		name      string
		mint      int64
		burn      int64
		reconvert int64
		expPass   bool
	}{
		{"ok - sufficient funds", 100, 10, 5, true},
		{"ok - equal funds", 10, 10, 10, true},
		{"fail - insufficient funds", 10, 1, 5, false},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			// Precondition: Convert Coin to FIP20
			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.priKey.PubKey().Address())
			receiver := sdk.AccAddress(suite.address.Bytes())
			suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins)
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.burn)),
				suite.address,
				sender,
			)

			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.app.IntrarelayerKeeper.ConvertCoin(ctx, msg)
			suite.Require().NoError(err, tc.name)
			suite.Commit()
			balance := suite.BalanceOf(common.HexToAddress(pair.Fip20Address), suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, metadata.Base)
			suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
			suite.Require().Equal(balance, big.NewInt(tc.burn))

			// Convert FIP20s back to Coins
			ctx = sdk.WrapSDKContext(suite.ctx)
			contractAddr := common.HexToAddress(pair.Fip20Address)
			pubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, suite.priKey.PubKey())
			msgConvertFIP20 := types.NewMsgConvertFIP20(
				sdk.NewInt(tc.reconvert),
				sender,
				receiver,
				contractAddr,
				pubKey,
			)

			res, err := suite.app.IntrarelayerKeeper.ConvertFIP20(ctx, msgConvertFIP20)
			expRes := &types.MsgConvertFIP20Response{}
			suite.Commit()
			balance = suite.BalanceOf(contractAddr, suite.address)
			cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expRes, res)
				suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn+tc.reconvert).Int64())
				suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.burn-tc.reconvert).Int64())
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
}

func (suite *KeeperTestSuite) TestConvertECR20NativeFIP20() {
	testCases := []struct {
		name     string
		mint     int64
		burn     int64
		malleate func()
		expPass  bool
	}{
		{"ok - sufficient funds", 100, 10, func() {}, true},
		{"ok - equal funds", 10, 10, func() {}, true},
		{"fail - insufficient funds - callEVM", 0, 10, func() {}, false},
		{
			"fail - minting disabled",
			100,
			10,
			func() {
				params := types.DefaultParams()
				params.EnableIntrarelayer = false
				suite.app.IntrarelayerKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			contractAddr := suite.setupRegisterFIP20Pair()
			suite.Require().NotNil(contractAddr)

			tc.malleate()
			suite.Commit()

			coinName := types.CreateDenom(contractAddr.String())
			receiver := sdk.AccAddress(suite.address.Bytes())
			pubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, suite.priKey.PubKey())
			sender := sdk.AccAddress(suite.priKey.PubKey().Address())
			msg := types.NewMsgConvertFIP20(
				sdk.NewInt(tc.burn),
				sender,
				receiver,
				contractAddr,
				pubKey,
			)
			suite.MintFIP20Token(contractAddr, suite.address, suite.address, big.NewInt(tc.mint))
			suite.Commit()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.app.IntrarelayerKeeper.ConvertFIP20(ctx, msg)
			expRes := &types.MsgConvertFIP20Response{}
			suite.Commit()
			balance := suite.BalanceOf(contractAddr, suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expRes, res)
				suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.burn))
				suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.mint-tc.burn).Int64())

			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
}

func (suite *KeeperTestSuite) TestConvertCoinNativeFIP20() {
	testCases := []struct {
		name      string
		mint      int64
		burn      int64
		reconvert int64
		expPass   bool
	}{
		{"ok - sufficient funds", 100, 10, 5, true},
		{"ok - equal funds", 10, 10, 10, true},
		{"fail - insufficient funds", 10, 1, 5, false},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			contractAddr := suite.setupRegisterFIP20Pair()
			suite.Require().NotNil(contractAddr)

			// Precondition: Convert FIP20 to Coins
			coinName := types.CreateDenom(contractAddr.String())
			pubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, suite.priKey.PubKey())
			sender := sdk.AccAddress(suite.priKey.PubKey().Address())
			receiver := sdk.AccAddress(suite.address.Bytes())
			suite.MintFIP20Token(contractAddr, suite.address, suite.address, big.NewInt(tc.mint))
			suite.Commit()
			msgConvertFIP20 := types.NewMsgConvertFIP20(
				sdk.NewInt(tc.burn),
				sender,
				receiver,
				contractAddr,
				pubKey,
			)

			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.app.IntrarelayerKeeper.ConvertFIP20(ctx, msgConvertFIP20)
			suite.Require().NoError(err)
			suite.Commit()
			balance := suite.BalanceOf(contractAddr, suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, receiver, coinName)
			suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.burn))
			suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.mint-tc.burn).Int64())

			// Convert Coins back to FIP20s
			ctx = sdk.WrapSDKContext(suite.ctx)
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(coinName, sdk.NewInt(tc.reconvert)),
				suite.address,
				sender,
			)

			res, err := suite.app.IntrarelayerKeeper.ConvertCoin(ctx, msg)
			expRes := &types.MsgConvertCoinResponse{}
			suite.Commit()
			balance = suite.BalanceOf(contractAddr, suite.address)
			cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, sender, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expRes, res)
				suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.burn-tc.reconvert).Int64())
				suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.mint-tc.burn+tc.reconvert).Int64())

			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
}

func (suite *KeeperTestSuite) TestConvertNativeIBC() {
	suite.SetupTest()
	validMetadata := banktypes.Metadata{
		Description: "ATOM IBC voucher (channel 14)",
		Base:        "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2",
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2",
				Exponent: 0,
			},
			{
				Denom:    displayCoinName,
				Exponent: uint32(18),
			},
		},
		Display: displayCoinName,
	}
	_, err := suite.app.IntrarelayerKeeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)
	suite.Commit()
}
