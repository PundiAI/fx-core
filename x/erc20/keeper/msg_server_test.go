package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestConvertCoinNativeCoin() {
	testCases := []struct {
		name           string
		mint           int64
		burn           int64
		malleate       func(common.Address)
		expPass        bool
		selfdestructed bool
	}{
		{
			"ok - sufficient funds",
			100,
			10,
			func(common.Address) {},
			true,
			false,
		},
		{
			"ok - equal funds",
			10,
			10,
			func(common.Address) {},
			true,
			false,
		},
		//{
		//	"ok - suicided contract",
		//	10,
		//	10,
		//	func(erc20 common.Address) {
		//		stateDb := suite.StateDB()
		//		ok := stateDb.Suicide(erc20)
		//		suite.Require().True(ok)
		//		suite.Require().NoError(stateDb.Commit())
		//	},
		//	true,
		//	true,
		//},
		{
			"fail - insufficient funds",
			0,
			10,
			func(common.Address) {},
			false,
			false,
		},
		{
			"fail - minting disabled",
			100,
			10,
			func(common.Address) {
				params := types.DefaultParams()
				params.EnableErc20 = false
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			false,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			erc20 := pair.GetERC20Contract()

			tc.malleate(erc20)

			coins := sdk.NewCoins(sdk.NewCoin(metadata.Base, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.signer.Address().Bytes())

			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))

			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(metadata.Base, sdk.NewInt(tc.burn)),
				suite.signer.Address(),
				sender,
			)
			res, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				balance := suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
				cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)

				acc := suite.app.EvmKeeper.GetAccountWithoutBalance(suite.ctx, erc20)
				if tc.selfdestructed {
					suite.Require().Nil(acc, "expected contract to be destroyed")
				} else {
					suite.Require().NotNil(acc)
				}

				if tc.selfdestructed || !acc.IsContract() {
					_, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, erc20.String())
					suite.Require().False(found)
				} else {
					suite.Require().Equal(&types.MsgConvertCoinResponse{}, res)
					suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
					suite.Require().Equal(balance.Int64(), big.NewInt(tc.burn).Int64())
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConvertERC20NativeCoin() {
	testCases := []struct {
		name      string
		mint      int64
		burn      int64
		reconvert int64
		malleate  func()
		expPass   bool
	}{
		{"ok - sufficient funds", 100, 10, 5, func() {}, true},
		{"ok - equal funds", 10, 10, 10, func() {}, true},
		{"fail - insufficient funds", 10, 1, 5, func() {}, false},
		{"fail ", 10, 1, -5, func() {}, false},
		{
			"fail - deleted module account - force fail", 100, 10, 5,
			func() {
				acc := suite.app.AccountKeeper.GetAccount(suite.ctx, authtypes.NewModuleAddress(types.ModuleName))
				suite.app.AccountKeeper.RemoveAccount(suite.ctx, acc)
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			// Precondition: Convert Coin to ERC20
			coins := sdk.NewCoins(sdk.NewCoin(metadata.Base, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.signer.Address().Bytes())
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(metadata.Base, sdk.NewInt(tc.burn)),
				suite.signer.Address(),
				sender,
			)
			_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), msg)
			suite.Require().NoError(err, tc.name)

			// suite.Commit()
			balance := suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)
			suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
			suite.Require().Equal(balance, big.NewInt(tc.burn))

			tc.malleate()

			contractAddr := pair.GetERC20Contract()
			msgConvertERC20 := types.NewMsgConvertERC20(
				sdk.NewInt(tc.reconvert),
				sender,
				contractAddr,
				suite.signer.Address(),
			)
			res, err := suite.app.Erc20Keeper.ConvertERC20(sdk.WrapSDKContext(suite.ctx), msgConvertERC20)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				balance = suite.BalanceOf(contractAddr, suite.signer.Address())
				cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)

				suite.Require().Equal(&types.MsgConvertERC20Response{}, res)
				suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn+tc.reconvert).Int64())
				suite.Require().Equal(balance.Int64(), big.NewInt(tc.burn-tc.reconvert).Int64())
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConvertERC20NativeERC20() {
	testCases := []struct {
		name           string
		mint           int64
		transfer       int64
		malleate       func(common.Address)
		expPass        bool
		selfdestructed bool
	}{
		{
			"ok - sufficient funds",
			100,
			10,
			func(common.Address) {},
			true,
			false,
		},
		{
			"ok - equal funds",
			10,
			10,
			func(common.Address) {},
			true,
			false,
		},
		{
			"ok - equal funds",
			10,
			10,
			func(common.Address) {},
			true,
			false,
		},
		//{
		//	"ok - suicided contract",
		//	10,
		//	10,
		//	func(erc20 common.Address) {
		//		stateDb := suite.StateDB()
		//		ok := stateDb.Suicide(erc20)
		//		suite.Require().True(ok)
		//		suite.Require().NoError(stateDb.Commit())
		//	},
		//	true,
		//	true,
		//},
		{
			"fail - insufficient funds - callEVM",
			0,
			10,
			func(common.Address) {},
			false,
			false,
		},
		{
			"fail - minting disabled",
			100,
			10,
			func(common.Address) {
				params := types.DefaultParams()
				params.EnableErc20 = false
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			false,
			false,
		},
		{
			"fail - negative transfer contract",
			10,
			-10,
			func(common.Address) {},
			false,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			contractAddr := suite.setupRegisterERC20Pair()

			tc.malleate(contractAddr)
			suite.Require().NotNil(contractAddr)
			// suite.Commit()

			suite.MintERC20Token(contractAddr, suite.signer.Address(), suite.signer.Address(), big.NewInt(tc.mint))

			receiver := sdk.AccAddress(suite.signer.Address().Bytes())
			msg := types.NewMsgConvertERC20(
				sdk.NewInt(tc.transfer),
				receiver,
				contractAddr,
				suite.signer.Address(),
			)
			res, err := suite.app.Erc20Keeper.ConvertERC20(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				balance := suite.BalanceOf(contractAddr, suite.signer.Address())
				cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, receiver, "test")

				acc := suite.app.EvmKeeper.GetAccountWithoutBalance(suite.ctx, contractAddr)
				if tc.selfdestructed {
					suite.Require().Nil(acc, "expected contract to be destroyed")
				} else {
					suite.Require().NotNil(acc)
				}

				if tc.selfdestructed || !acc.IsContract() {
					_, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, contractAddr.String())
					suite.Require().False(found)
				} else {
					suite.Require().Equal(&types.MsgConvertERC20Response{}, res)
					suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.transfer))
					suite.Require().Equal(balance.Int64(), big.NewInt(tc.mint-tc.transfer).Int64())
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConvertCoinNativeERC20() {
	testCases := []struct {
		name    string
		mint    int64
		convert int64
		expPass bool
	}{
		{
			"ok - sufficient funds",
			100,
			10,
			true,
		},
		{
			"ok - equal funds",
			100,
			100,
			true,
		},
		{
			"fail - insufficient funds",
			100,
			200,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			contractAddr := suite.setupRegisterERC20Pair()
			suite.Require().NotNil(contractAddr)

			pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, contractAddr.String())
			suite.True(found)
			coins := sdk.NewCoins(sdk.NewCoin(pair.Denom, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.signer.Address().Bytes())

			// Precondition: Mint Coins to convert on sender account
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))
			denom := "test"
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, denom)
			suite.Require().Equal(sdk.NewInt(tc.mint), cosmosBalance.Amount)

			// Precondition: Mint escrow tokens on module account
			// suite.GrantERC20Token(contractAddr, suite.signer.Address(), types.ModuleAddress, "MINTER_ROLE")
			erc20ModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
			suite.MintERC20Token(contractAddr, suite.signer.Address(), erc20ModuleAddr, big.NewInt(tc.mint))
			tokenBalance := suite.BalanceOf(contractAddr, erc20ModuleAddr)
			suite.Require().Equal(big.NewInt(tc.mint), tokenBalance)

			// Convert Coins back to ERC20s
			receiver := suite.signer.Address()
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(denom, sdk.NewInt(tc.convert)),
				receiver,
				sender,
			)
			res, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), msg)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				tokenBalance = suite.BalanceOf(contractAddr, suite.signer.Address())
				cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, sender, denom)

				suite.Require().Equal(&types.MsgConvertCoinResponse{}, res)
				suite.Require().Equal(sdk.NewInt(tc.mint-tc.convert), cosmosBalance.Amount)
				suite.Require().Equal(big.NewInt(tc.convert), tokenBalance)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestWrongPairOwnerERC20NativeCoin() {
	testCases := []struct {
		name      string
		mint      int64
		burn      int64
		reconvert int64
		expPass   bool
	}{
		{"ok - sufficient funds", 100, 10, 5, true},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			// Precondition: Convert Coin to ERC20
			coins := sdk.NewCoins(sdk.NewCoin(metadata.Base, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.signer.Address().Bytes())
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))

			pair.ContractOwner = types.OWNER_UNSPECIFIED
			suite.app.Erc20Keeper.SetTokenPair(suite.ctx, *pair)

			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(metadata.Base, sdk.NewInt(tc.burn)),
				suite.signer.Address(),
				sender,
			)
			_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), msg)
			suite.Require().Error(err, tc.name)

			// Convert ERC20s back to Coins
			msgConvertERC20 := types.NewMsgConvertERC20(
				sdk.NewInt(tc.reconvert),
				sender,
				pair.GetERC20Contract(),
				suite.signer.Address(),
			)
			_, err = suite.app.Erc20Keeper.ConvertERC20(sdk.WrapSDKContext(suite.ctx), msgConvertERC20)
			suite.Require().Error(err, tc.name)
		})
	}
}
