package keeper_test

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestBaseDenomToBridgeDenom() {
	testCases := []struct {
		Name      string
		BaseDenom string
		Target    string
		Success   bool
		Error     string
	}{
		{
			Name:      "success",
			BaseDenom: "usdt",
			Target:    suite.chainName,
			Success:   true,
		},
		{
			Name:      "failed - bridge denom not found",
			BaseDenom: "usdt1",
			Target:    suite.chainName,
			Success:   false,
			Error:     "denom usdt1 not found",
		},
		{
			Name:      "failed - target not found",
			BaseDenom: "usdt",
			Target:    "abc",
			Success:   false,
			Error:     "not found bridge denom: usdt, abc: invalid coins",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			moduleBridgeDenom := types.NewBridgeDenom(suite.chainName, helpers.GenExternalAddr(suite.chainName))
			suite.SetToken("USDT", moduleBridgeDenom)

			bridgeDenom, err := suite.Keeper().BaseDenomToBridgeDenom(suite.Ctx, tc.BaseDenom, tc.Target)
			if tc.Success {
				suite.NoError(err, tc.Name)
				suite.Equal(moduleBridgeDenom, bridgeDenom, tc.Name)
			} else {
				suite.Error(err, tc.Name)
				suite.Equal(tc.Error, err.Error(), tc.Name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestManyToOne() {
	bridgeDenom := types.NewBridgeDenom(suite.chainName, helpers.GenExternalAddr(suite.chainName))

	testCases := []struct {
		Name         string
		ConvertDenom string
		Target       string
		Success      bool
		ExpectDenom  string
	}{
		{
			Name:         "success - FX",
			ConvertDenom: fxtypes.DefaultDenom,
			Target:       "eth",
			Success:      true,
			ExpectDenom:  fxtypes.DefaultDenom,
		},
		{
			Name:         "success - base to bridge denom",
			ConvertDenom: "usdt",
			Target:       suite.chainName,
			Success:      true,
			ExpectDenom:  bridgeDenom,
		},
		{
			Name:         "success - bridge to base denom",
			ConvertDenom: bridgeDenom,
			Target:       "",
			Success:      true,
			ExpectDenom:  "usdt",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			suite.SetToken("USDT", bridgeDenom)

			targetDenom, err := suite.Keeper().ManyToOne(suite.Ctx, tc.ConvertDenom, tc.Target)
			if tc.Success {
				suite.NoError(err)
				suite.Equal(tc.ExpectDenom, targetDenom, tc.Name)
			} else {
				suite.Error(err, tc.Name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConversionCoin() {
	amount := helpers.NewRandAmount()
	baseCoin := suite.NewCoin(amount)
	bridgeCoin, _ := suite.NewBridgeCoin(suite.chainName, amount)

	testCases := []struct {
		Name             string
		Malleate         func(expCoin sdk.Coin)
		ConvertCoin      sdk.Coin
		TargetDenom      string
		Success          bool
		ModuleExpBalance sdk.Coins
		Error            string
	}{
		{
			Name:        "success - convert FX",
			ConvertCoin: sdk.NewCoin(fxtypes.DefaultDenom, amount),
			TargetDenom: fxtypes.DefaultDenom,
			Success:     true,
		},
		{
			Name: "success - convert native coin base to bridge",
			Malleate: func(expCoin sdk.Coin) {
				suite.AddTokenPair(baseCoin.Denom, true)
				suite.MintTokenToModule(suite.chainName, expCoin)
			},
			ConvertCoin:      baseCoin,
			TargetDenom:      bridgeCoin.Denom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
		{
			Name: "success - native erc20 base to bridge",
			Malleate: func(expCoin sdk.Coin) {
				suite.AddTokenPair(baseCoin.Denom, false)
			},
			ConvertCoin:      baseCoin,
			TargetDenom:      bridgeCoin.Denom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
		{
			Name: "success - native coin bridge to base",
			Malleate: func(expCoin sdk.Coin) {
				suite.AddTokenPair(baseCoin.Denom, true)
			},
			ConvertCoin:      bridgeCoin,
			TargetDenom:      baseCoin.Denom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(bridgeCoin),
		},
		{
			Name: "success - native erc20 bridge to base",
			Malleate: func(expCoin sdk.Coin) {
				suite.AddTokenPair(baseCoin.Denom, false)
			},
			ConvertCoin:      bridgeCoin,
			TargetDenom:      baseCoin.Denom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
		{
			Name:        "failed - denom not found",
			ConvertCoin: bridgeCoin,
			TargetDenom: baseCoin.Denom,
			Success:     false,
			Error:       fmt.Sprintf("token pair not found %s: invalid coins", baseCoin.Denom),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			acc := helpers.GenAccAddress()
			suite.MintToken(acc, tc.ConvertCoin)

			tc.ModuleExpBalance = tc.ModuleExpBalance.Add(suite.Balance(suite.ModuleAddress())...)

			expCoin := sdk.NewCoin(tc.TargetDenom, amount)
			if tc.Malleate != nil {
				tc.Malleate(expCoin)
			}

			err := suite.Keeper().ConversionCoin(suite.Ctx, acc, tc.ConvertCoin, baseCoin.Denom, tc.TargetDenom)
			if tc.Success {
				suite.NoError(err)
				suite.CheckAllBalance(acc, expCoin)
				suite.CheckAllBalance(suite.ModuleAddress(), tc.ModuleExpBalance...)
			} else {
				suite.Error(err)
				suite.Equal(tc.Error, err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestDepositBridgeToken() {
	amount := helpers.NewRandAmount()
	baseCoin := suite.NewCoin(amount)
	birdgeCoin, _ := suite.NewBridgeCoin(suite.chainName, amount)

	testCases := []struct {
		Name         string
		BridgeToken  sdk.Coin
		IsNativeCoin bool
		Success      bool
		Error        string
	}{
		{
			Name:        "success - deposit FX",
			BridgeToken: sdk.NewCoin(fxtypes.DefaultDenom, amount),
			Success:     true,
		},
		{
			Name:         "success - deposit native coin bridge token",
			BridgeToken:  birdgeCoin,
			IsNativeCoin: true,
			Success:      true,
		},
		{
			Name:        "success - deposit native erc20 bridge token",
			BridgeToken: birdgeCoin,
			Success:     true,
		},
		{
			Name:        "failed - bridge denom not found",
			BridgeToken: sdk.NewCoin("aaa", amount),
			Success:     false,
			Error:       "alias aaa not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			acc := helpers.GenAccAddress()
			suite.SetToken(strings.ToUpper(baseCoin.Denom), birdgeCoin.Denom)
			suite.AddTokenPair(baseCoin.Denom, tc.IsNativeCoin)

			moduleBalance := suite.Balance(suite.ModuleAddress())

			if !tc.IsNativeCoin {
				suite.MintTokenToModule(suite.chainName, tc.BridgeToken)
			}

			err := suite.Keeper().DepositBridgeToken(suite.Ctx, tc.BridgeToken, acc)
			if tc.Success {
				suite.NoError(err)
				suite.CheckAllBalance(acc, tc.BridgeToken)
				suite.CheckAllBalance(suite.ModuleAddress(), moduleBalance...)
			} else {
				suite.Error(err)
				suite.Equal(tc.Error, err.Error())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestWithdrawBridgeToken() {
	amount := helpers.NewRandAmount()
	baseCoin := suite.NewCoin(amount)
	bridgeCoin, _ := suite.NewBridgeCoin(suite.chainName, amount)

	testCases := []struct {
		Name             string
		BridgeToken      sdk.Coin
		IsNativeCoin     bool
		Success          bool
		ModuleExpBalance sdk.Coins
	}{
		{
			Name:             "success - withdraw FX",
			BridgeToken:      sdk.NewCoin(fxtypes.DefaultDenom, amount),
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)),
		},
		{
			Name:             "success - withdraw native coin",
			BridgeToken:      bridgeCoin,
			IsNativeCoin:     true,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
		{
			Name:             "success - withdraw native erc20",
			BridgeToken:      bridgeCoin,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(bridgeCoin),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			acc := helpers.GenAccAddress()

			suite.SetToken(strings.ToUpper(baseCoin.Denom), bridgeCoin.Denom)
			suite.AddTokenPair(baseCoin.Denom, tc.IsNativeCoin)
			suite.MintToken(acc, tc.BridgeToken)

			tc.ModuleExpBalance = tc.ModuleExpBalance.Add(suite.Balance(suite.ModuleAddress())...)

			err := suite.Keeper().WithdrawBridgeToken(suite.Ctx, tc.BridgeToken, acc)
			if tc.Success {
				suite.NoError(err)
				suite.CheckAllBalance(acc, sdk.NewCoins()...)
				suite.CheckAllBalance(suite.ModuleAddress(), tc.ModuleExpBalance...)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBridgeTokenToBaseCoin() {
	amount := helpers.NewRandAmount()
	baseCoin := suite.NewCoin(amount)
	bridgeCoin, tokenAddr := suite.NewBridgeCoin(suite.chainName, amount)

	testCases := []struct {
		Name             string
		BaseDenom        string
		BridgeDenom      string
		TokenAddr        string
		IsNativeCoin     bool
		Success          bool
		ModuleExpBalance sdk.Coins
	}{
		{
			Name:             "success - FX",
			TokenAddr:        helpers.GenHexAddress().String(),
			BaseDenom:        fxtypes.DefaultDenom,
			BridgeDenom:      fxtypes.DefaultDenom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
		{
			Name:             "success - native coin",
			TokenAddr:        tokenAddr,
			BaseDenom:        baseCoin.Denom,
			BridgeDenom:      bridgeCoin.Denom,
			IsNativeCoin:     true,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(bridgeCoin),
		},
		{
			Name:             "success - native erc20",
			TokenAddr:        tokenAddr,
			BaseDenom:        baseCoin.Denom,
			BridgeDenom:      bridgeCoin.Denom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			acc := helpers.GenAccAddress()
			suite.SetToken(strings.ToUpper(tc.BaseDenom), tc.BridgeDenom)
			suite.AddTokenPair(tc.BaseDenom, tc.IsNativeCoin)
			suite.AddBridgeToken(tc.TokenAddr, strings.ToUpper(tc.BaseDenom))

			tc.ModuleExpBalance = tc.ModuleExpBalance.Add(suite.Balance(suite.ModuleAddress())...)

			if !tc.IsNativeCoin {
				suite.MintTokenToModule(suite.chainName, sdk.NewCoin(tc.BridgeDenom, amount))
			}

			_, err := suite.Keeper().BridgeTokenToBaseCoin(suite.Ctx, tc.TokenAddr, amount.BigInt(), acc)
			if tc.Success {
				suite.NoError(err)
				suite.CheckAllBalance(acc, sdk.NewCoin(tc.BaseDenom, amount))
				suite.CheckAllBalance(suite.ModuleAddress(), tc.ModuleExpBalance...)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBaseCoinToBridgeToken() {
	amount := helpers.NewRandAmount()
	baseCoin := suite.NewCoin(amount)
	bridgeCoin, _ := suite.NewBridgeCoin(suite.chainName, amount)

	testCases := []struct {
		Name             string
		Coin             sdk.Coin
		BridgeDenom      string
		IsNativeCoin     bool
		Success          bool
		ModuleExpBalance sdk.Coins
	}{
		{
			Name:             "success - FX",
			Coin:             sdk.NewCoin(fxtypes.DefaultDenom, amount),
			BridgeDenom:      fxtypes.DefaultDenom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)),
		},
		{
			Name:             "success - native coin",
			Coin:             baseCoin,
			BridgeDenom:      bridgeCoin.Denom,
			IsNativeCoin:     true,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(),
		},
		{
			Name:             "success - native erc20",
			Coin:             baseCoin,
			BridgeDenom:      bridgeCoin.Denom,
			Success:          true,
			ModuleExpBalance: sdk.NewCoins(bridgeCoin),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			acc := helpers.GenAccAddress()
			suite.MintToken(acc, tc.Coin)
			suite.SetToken(strings.ToLower(tc.Coin.Denom), tc.BridgeDenom)
			suite.AddTokenPair(strings.ToLower(tc.Coin.Denom), tc.IsNativeCoin)

			tc.ModuleExpBalance = tc.ModuleExpBalance.Add(suite.Balance(suite.ModuleAddress())...)

			if tc.IsNativeCoin {
				suite.MintTokenToModule(suite.chainName, sdk.NewCoin(tc.BridgeDenom, amount))
			}

			_, _, err := suite.Keeper().BaseCoinToBridgeToken(suite.Ctx, suite.chainName, tc.Coin, acc)
			if tc.Success {
				suite.NoError(err)
				suite.CheckAllBalance(acc, sdk.NewCoins()...)
				suite.CheckAllBalance(suite.ModuleAddress(), tc.ModuleExpBalance...)
			} else {
				suite.Error(err)
			}
		})
	}
}
