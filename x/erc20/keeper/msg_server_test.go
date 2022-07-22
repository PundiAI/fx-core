package keeper_test

import (
	"fmt"
	"math/big"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v2/x/erc20/types"
)

func (suite *KeeperTestSuite) TestConvertCoinNativeCoin() {
	testCases := []struct {
		name           string
		mint           int64
		burn           int64
		malleate       func(common.Address)
		extra          func()
		expPass        bool
		selfdestructed bool
	}{
		{"ok - sufficient funds",
			100,
			10,
			func(common.Address) {},
			func() {},
			true,
			false},
		{"ok - equal funds",
			10,
			10,
			func(common.Address) {},
			func() {},
			true,
			false,
		},
		{
			"ok - suicided contract",
			10,
			10,
			func(erc20 common.Address) {
				stateDb := suite.StateDB()
				ok := stateDb.Suicide(erc20)
				suite.Require().True(ok)
				suite.Require().NoError(stateDb.Commit())
			},
			func() {},
			true,
			true,
		},
		{"fail - insufficient funds",
			0,
			10,
			func(common.Address) {},
			func() {},
			false,
			false},
		{
			"fail - minting disabled",
			100,
			10,
			func(common.Address) {
				params := types.DefaultParams()
				params.EnableErc20 = false
				suite.app.Erc20Keeper.SetParams(suite.ctx, params)
			},
			func() {},
			false,
			false,
		},
		{
			"fail - deleted module account - force fail",
			100,
			10,
			func(common.Address) {},
			func() {
				acc := suite.app.AccountKeeper.GetAccount(suite.ctx, types.ModuleAddress.Bytes())
				suite.app.AccountKeeper.RemoveAccount(suite.ctx, acc)
			},
			false,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			erc20 := pair.GetERC20Contract()
			tc.malleate(erc20)
			suite.Commit()

			ctx := sdk.WrapSDKContext(suite.ctx)
			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenBase, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.address.Bytes())
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenBase, sdk.NewInt(tc.burn)),
				suite.address,
				sender,
			)

			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))

			tc.extra()
			res, err := suite.app.Erc20Keeper.ConvertCoin(ctx, msg)
			expRes := &types.MsgConvertCoinResponse{}
			suite.Commit()
			balance := suite.BalanceOf(common.HexToAddress(pair.Erc20Address), suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				acc := suite.app.EvmKeeper.GetAccountWithoutBalance(suite.ctx, erc20)
				if tc.selfdestructed {
					suite.Require().Nil(acc, "expected contract to be destroyed")
				} else {
					suite.Require().NotNil(acc)
				}

				if tc.selfdestructed || !acc.IsContract() {
					id := suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, erc20.String())
					_, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, id)
					suite.Require().False(found)
				} else {
					suite.Require().Equal(expRes, res)
					suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
					suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.burn).Int64())
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
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
				acc := suite.app.AccountKeeper.GetAccount(suite.ctx, types.ModuleAddress.Bytes())
				suite.app.AccountKeeper.RemoveAccount(suite.ctx, acc)
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
			suite.Require().NotNil(pair)

			// Precondition: Convert Coin to ERC20
			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenBase, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.address.Bytes())
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenBase, sdk.NewInt(tc.burn)),
				suite.address,
				sender,
			)

			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.app.Erc20Keeper.ConvertCoin(ctx, msg)
			suite.Require().NoError(err, tc.name)
			suite.Commit()
			balance := suite.BalanceOf(common.HexToAddress(pair.Erc20Address), suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)
			suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
			suite.Require().Equal(balance, big.NewInt(tc.burn))

			// Convert ERC20s back to Coins
			ctx = sdk.WrapSDKContext(suite.ctx)
			contractAddr := common.HexToAddress(pair.Erc20Address)
			msgConvertERC20 := types.NewMsgConvertERC20(
				sdk.NewInt(tc.reconvert),
				sender,
				contractAddr,
				suite.address,
			)

			tc.malleate()
			//set pubkey before covert erc20
			//acc := suite.app.AccountKeeper.GetAccount(suite.ctx, suite.address.Bytes())
			//err = acc.SetPubKey(suite.privateKey.PubKey())
			//suite.Require().NoError(err)
			//suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			res, err := suite.app.Erc20Keeper.ConvertERC20(ctx, msgConvertERC20)
			expRes := &types.MsgConvertERC20Response{}
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

func (suite *KeeperTestSuite) TestConvertERC20NativeERC20() {
	var contractAddr common.Address

	testCases := []struct {
		name           string
		mint           int64
		transfer       int64
		malleate       func(common.Address)
		extra          func()
		contractType   int
		expPass        bool
		selfdestructed bool
	}{
		{
			"ok - sufficient funds",
			100,
			10,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			true,
			false,
		},
		{
			"ok - equal funds",
			10,
			10,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			true,
			false,
		},
		{
			"ok - equal funds",
			10,
			10,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			true,
			false,
		},
		{
			"ok - suicided contract",
			10,
			10,
			func(erc20 common.Address) {
				stateDb := suite.StateDB()
				ok := stateDb.Suicide(erc20)
				suite.Require().True(ok)
				suite.Require().NoError(stateDb.Commit())
			},
			func() {},
			contractMinterBurner,
			true,
			true,
		},
		{
			"fail - insufficient funds - callEVM",
			0,
			10,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
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
			func() {},
			contractMinterBurner,
			false,
			false,
		},
		{
			"fail - negative transfer contract",
			10,
			-10,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			false,
			false,
		},
		{
			"fail - no module address",
			100,
			10,
			func(common.Address) {
			},
			func() {
				acc := suite.app.AccountKeeper.GetAccount(suite.ctx, types.ModuleAddress.Bytes())
				suite.app.AccountKeeper.RemoveAccount(suite.ctx, acc)
			},
			contractMinterBurner,
			false,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()

			contractAddr = suite.setupRegisterERC20Pair(tc.contractType)

			tc.malleate(contractAddr)
			suite.Require().NotNil(contractAddr)
			suite.Commit()

			coinName := types.CreateDenom(contractAddr.String())
			sender := sdk.AccAddress(suite.address.Bytes())
			msg := types.NewMsgConvertERC20(
				sdk.NewInt(tc.transfer),
				sender,
				contractAddr,
				suite.address,
			)

			suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(tc.mint))
			suite.Commit()
			ctx := sdk.WrapSDKContext(suite.ctx)

			tc.extra()
			res, err := suite.app.Erc20Keeper.ConvertERC20(ctx, msg)

			expRes := &types.MsgConvertERC20Response{}
			suite.Commit()
			balance := suite.BalanceOf(contractAddr, suite.address)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				acc := suite.app.EvmKeeper.GetAccountWithoutBalance(suite.ctx, contractAddr)
				if tc.selfdestructed {
					suite.Require().Nil(acc, "expected contract to be destroyed")
				} else {
					suite.Require().NotNil(acc)
				}

				if tc.selfdestructed || !acc.IsContract() {
					id := suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, contractAddr.String())
					_, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, id)
					suite.Require().False(found)
				} else {
					suite.Require().Equal(expRes, res)
					suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.transfer))
					suite.Require().Equal(balance.(*big.Int).Int64(), big.NewInt(tc.mint-tc.transfer).Int64())
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
}

func (suite *KeeperTestSuite) TestConvertCoinNativeERC20() {
	var contractAddr common.Address

	testCases := []struct {
		name         string
		mint         int64
		convert      int64
		malleate     func(common.Address)
		extra        func()
		contractType int
		expPass      bool
	}{
		{
			"ok - sufficient funds",
			100,
			10,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			true,
		},
		{
			"ok - equal funds",
			100,
			100,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			true,
		},
		{
			"fail - insufficient funds",
			100,
			200,
			func(common.Address) {},
			func() {},
			contractMinterBurner,
			false,
		},
		{
			"fail - deleted module address - force fail",
			100,
			10,
			func(common.Address) {},
			func() {
				acc := suite.app.AccountKeeper.GetAccount(suite.ctx, types.ModuleAddress.Bytes())
				suite.app.AccountKeeper.RemoveAccount(suite.ctx, acc)
			},
			contractMinterBurner,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			contractAddr = suite.setupRegisterERC20Pair(tc.contractType)
			suite.Require().NotNil(contractAddr)

			id := suite.app.Erc20Keeper.GetTokenPairID(suite.ctx, contractAddr.String())
			pair, _ := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, id)
			coins := sdk.NewCoins(sdk.NewCoin(pair.Denom, sdk.NewInt(tc.mint)))
			coinName := types.CreateDenom(contractAddr.String())
			sender := sdk.AccAddress(suite.address.Bytes())

			// Precondition: Mint Coins to convert on sender account
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, coinName)
			suite.Require().Equal(sdk.NewInt(tc.mint), cosmosBalance.Amount)

			// Precondition: Mint escrow tokens on module account
			//suite.GrantERC20Token(contractAddr, suite.address, types.ModuleAddress, "MINTER_ROLE")
			suite.MintERC20Token(contractAddr, suite.address, types.ModuleAddress, big.NewInt(tc.mint))
			tokenBalance := suite.BalanceOf(contractAddr, types.ModuleAddress)
			suite.Require().Equal(big.NewInt(tc.mint), tokenBalance)

			tc.malleate(contractAddr)
			suite.Commit()

			// Convert Coins back to ERC20s
			receiver := suite.address
			ctx := sdk.WrapSDKContext(suite.ctx)
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(coinName, sdk.NewInt(tc.convert)),
				receiver,
				sender,
			)

			tc.extra()
			res, err := suite.app.Erc20Keeper.ConvertCoin(ctx, msg)

			expRes := &types.MsgConvertCoinResponse{}
			//suite.Commit()
			tokenBalance = suite.BalanceOf(contractAddr, suite.address)
			cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, sender, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expRes, res)
				suite.Require().Equal(sdk.NewInt(tc.mint-tc.convert), cosmosBalance.Amount)
				suite.Require().Equal(big.NewInt(tc.convert), tokenBalance.(*big.Int))
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
	suite.mintFeeCollector = false
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
			suite.mintFeeCollector = true
			suite.SetupTest()
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			// Precondition: Convert Coin to ERC20
			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenBase, sdk.NewInt(tc.mint)))
			sender := sdk.AccAddress(suite.address.Bytes())
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))
			msg := types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenBase, sdk.NewInt(tc.burn)),
				suite.address,
				sender,
			)

			pair.ContractOwner = types.OWNER_UNSPECIFIED
			suite.app.Erc20Keeper.SetTokenPair(suite.ctx, *pair)

			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.app.Erc20Keeper.ConvertCoin(ctx, msg)
			suite.Require().Error(err, tc.name)

			// Convert ERC20s back to Coins
			ctx = sdk.WrapSDKContext(suite.ctx)
			contractAddr := common.HexToAddress(pair.Erc20Address)
			msgConvertERC20 := types.NewMsgConvertERC20(
				sdk.NewInt(tc.reconvert),
				sender,
				contractAddr,
				suite.address,
			)

			_, err = suite.app.Erc20Keeper.ConvertERC20(ctx, msgConvertERC20)
			suite.Require().Error(err, tc.name)
		})
	}
}

func (suite *KeeperTestSuite) TestConvertDenom() {
	suite.supportManyToOneBlock = true
	priv1, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	addr1 := common.BytesToAddress(priv1.PubKey().Address().Bytes())

	priv2, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	addr2 := common.BytesToAddress(priv2.PubKey().Address().Bytes())

	tronUSDT := sdk.NewCoin(tronDenom, sdk.NewInt(100))

	testCases := []struct {
		name     string
		register func()
		malleate func(*types.TokenPair) error
		expPass  bool
		errMsg   string
	}{
		{
			"ok",
			func() {
				usdtMatedata, pair = suite.setupRegisterCoinUSDT()
				suite.Require().NotNil(usdtMatedata)

				md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.Denom)
				suite.Require().True(found)
				suite.Require().True(types.IsManyToOneMetadata(md))

				denom := suite.app.Erc20Keeper.GetAliasDenom(suite.ctx, usdtMatedata.DenomUnits[0].Aliases[0])
				suite.Require().True(len(denom) > 0)
			},
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     tronUSDT,
					Target:   "",
				})
				return err
			},
			true,
			"",
		},
		{
			"denom already registered",
			func() {
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, tronDenom, []byte{})
			},
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     tronUSDT,
					Target:   "",
				})
				return err
			},
			false,
			"denom tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t already registered: invalid denom",
		},
		{
			"alias not registered",
			func() {},
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     tronUSDT,
					Target:   "",
				})
				return err
			},
			false,
			"alias tronTR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t not registered: invalid denom",
		},
		{
			"alias denom not registered",
			func() {
				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "usdt", tronDenom, polygonDenom)
			},
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     tronUSDT,
					Target:   "",
				})
				return err
			},
			false,
			"denom usdt not registered: invalid denom",
		},
		{
			"denom not support many to one",
			func() {
				usdtMatedata, pair = suite.setupRegisterCoinUSDTWithOutAlias()
				suite.Require().NotNil(usdtMatedata)

				md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.Denom)
				suite.Require().True(found)
				suite.Require().False(types.IsManyToOneMetadata(md))

				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "usdt", tronDenom, polygonDenom)
			},
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     tronUSDT,
					Target:   "",
				})
				return err
			},
			false,
			"not support with usdt: invalid metadata",
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			if tc.register != nil {
				tc.register()
			}

			//mint and transfer
			err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(tronUSDT))
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr1.Bytes(), sdk.NewCoins(tronUSDT))
			suite.Require().NoError(err)

			beforeBalanceManyResp := suite.app.BankKeeper.GetBalance(suite.ctx, addr1.Bytes(), tronDenom)
			beforeBalanceOneResp := suite.app.BankKeeper.GetBalance(suite.ctx, addr2.Bytes(), "usdt")

			tcErr := tc.malleate(pair)

			afterBalanceManyResp := suite.app.BankKeeper.GetBalance(suite.ctx, addr1.Bytes(), tronDenom)
			afterBalanceOneResp := suite.app.BankKeeper.GetBalance(suite.ctx, addr2.Bytes(), "usdt")

			if tc.expPass {
				suite.Require().NoError(tcErr)
				suite.Require().Equal(beforeBalanceManyResp.Amount.Sub(afterBalanceManyResp.Amount), tronUSDT.Amount)
				suite.Require().Equal(afterBalanceOneResp.Amount.Sub(beforeBalanceOneResp.Amount), tronUSDT.Amount)
			} else {
				suite.Require().Error(tcErr, tc.name)
				suite.Require().EqualError(tcErr, tc.errMsg, tc.name)
			}
		})
	}

	suite.supportManyToOneBlock = false
}

func (suite *KeeperTestSuite) TestConvertDenomWithTarget() {
	suite.supportManyToOneBlock = true
	priv1, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	addr1 := common.BytesToAddress(priv1.PubKey().Address().Bytes())

	priv2, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	addr2 := common.BytesToAddress(priv2.PubKey().Address().Bytes())

	registerFn := func() {
		usdtMatedata, pair = suite.setupRegisterCoinUSDT()
		suite.Require().NotNil(usdtMatedata)

		md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.Denom)
		suite.Require().True(found)
		suite.Require().True(types.IsManyToOneMetadata(md))
	}

	tronUSDT := sdk.NewCoin(tronDenom, sdk.NewInt(100))
	polygonUSDT := sdk.NewCoin(polygonDenom, sdk.NewInt(100))
	usdt := sdk.NewCoin("usdt", sdk.NewInt(1))

	testCases := []struct {
		name     string
		register func()
		malleate func(*types.TokenPair) error
		expPass  bool
		errMsg   string
	}{
		{
			"ok",
			registerFn,
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "tron",
				})
				if err != nil {
					return err
				}
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "polygon",
				})
				return err
			},
			true,
			"",
		},
		{
			"denom not registered",
			nil,
			func(pair *types.TokenPair) error {
				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "tron",
				})
				return err
			},
			false,
			"denom usdt not registered: invalid denom",
		},
		{
			"metadata not found",
			nil,
			func(pair *types.TokenPair) error {
				suite.app.Erc20Keeper.SetDenomMap(suite.ctx, "usdt", []byte{})

				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "tron",
				})
				return err
			},
			false,
			"denom usdt not found: invalid metadata",
		},
		{
			"metadata not support many to one",
			nil,
			func(_ *types.TokenPair) error {
				usdtMatedata, pair = suite.setupRegisterCoinUSDTWithOutAlias()
				suite.Require().NotNil(usdtMatedata)
				md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.Denom)
				suite.Require().True(found)
				suite.Require().False(types.IsManyToOneMetadata(md))

				suite.app.Erc20Keeper.SetAliasesDenom(suite.ctx, "usdt", tronDenom, polygonDenom)

				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "tron",
				})
				return err
			},
			false,
			"denom usdt metadata not support: invalid metadata",
		},
		{
			"target denom not exist",
			registerFn,
			func(pair *types.TokenPair) error {
				usdtCopy := usdtMatedata
				usdtCopy.DenomUnits[0].Aliases = []string{polygonDenom}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, usdtCopy)

				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "tron",
				})
				return err
			},
			false,
			"target tron denom not exist: invalid target",
		},
		{
			"alias not registered",
			registerFn,
			func(pair *types.TokenPair) error {
				usdtCopy := usdtMatedata
				usdtCopy.DenomUnits[0].Aliases = append(usdtCopy.DenomUnits[0].Aliases, "bsc0x0000000000000000000000000000000000000000")
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, usdtCopy)

				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender:   sdk.AccAddress(addr1.Bytes()).String(),
					Receiver: sdk.AccAddress(addr2.Bytes()).String(),
					Coin:     usdt,
					Target:   "bsc",
				})
				return err
			},
			false,
			"alias bsc0x0000000000000000000000000000000000000000 not registered: invalid denom",
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			//mint and transfer
			err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(tronUSDT))
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr1.Bytes(), sdk.NewCoins(tronUSDT))
			suite.Require().NoError(err)

			err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(polygonUSDT))
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr1.Bytes(), sdk.NewCoins(polygonUSDT))
			suite.Require().NoError(err)

			if tc.register != nil {
				tc.register()

				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender: sdk.AccAddress(addr1.Bytes()).String(), Receiver: sdk.AccAddress(addr1.Bytes()).String(), Coin: tronUSDT, Target: ""})
				suite.Require().NoError(err)

				_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &types.MsgConvertDenom{
					Sender: sdk.AccAddress(addr1.Bytes()).String(), Receiver: sdk.AccAddress(addr1.Bytes()).String(), Coin: polygonUSDT, Target: ""})
				suite.Require().NoError(err)

				usdtBalanceResp, err := suite.app.BankKeeper.Balance(sdk.WrapSDKContext(suite.ctx),
					&banktypes.QueryBalanceRequest{Address: sdk.AccAddress(addr1.Bytes()).String(), Denom: "usdt"})
				suite.Require().NoError(err)
				suite.Require().Equal(usdtBalanceResp.Balance.Amount, tronUSDT.Amount.Add(polygonUSDT.Amount))
			}

			beforeAddr1UsdtBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr1.Bytes(), "usdt")
			beforeAddr2TronBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr2.Bytes(), tronDenom)
			beforeAddr2PolygonBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr2.Bytes(), polygonDenom)

			// malleate
			tcErr := tc.malleate(pair)

			afterAddr1UsdtBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr1.Bytes(), "usdt")
			afterAddr2TronBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr2.Bytes(), tronDenom)
			afterAddr2PolygonBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr2.Bytes(), polygonDenom)

			if tc.expPass {
				suite.Require().NoError(tcErr, tc.name)
				suite.Require().Equal(beforeAddr1UsdtBalance, afterAddr1UsdtBalance.Add(usdt).Add(usdt))
				suite.Require().Equal(afterAddr2TronBalance.Sub(beforeAddr2TronBalance).Amount, usdt.Amount)
				suite.Require().Equal(afterAddr2PolygonBalance.Sub(beforeAddr2PolygonBalance).Amount, usdt.Amount)
			} else {
				suite.Require().Error(tcErr, tc.name)
				suite.Require().EqualError(tcErr, tc.errMsg, tc.name)
			}
		})
	}

	suite.supportManyToOneBlock = false
}
