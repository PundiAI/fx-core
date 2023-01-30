package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	helpers2 "github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
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
		// {
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
		// },
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
		// {
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
		// },
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

			suite.MintERC20Token(suite.signer, contractAddr, suite.signer.Address(), big.NewInt(tc.mint))

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
			suite.MintERC20Token(suite.signer, contractAddr, erc20ModuleAddr, big.NewInt(tc.mint))
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

func (suite *KeeperTestSuite) TestToTargetDenom() {
	testCases := []struct {
		name     string
		malleate func() (string, string, []string, fxtypes.FxTarget, string)
	}{
		{
			name: "empty target, expect base",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				denom := helpers2.NewRandDenom()
				base := denom
				return denom, base, []string{}, fxtypes.ParseFxTarget(""), denom
			},
		},
		{
			name: "erc20 target, expect base",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				denom := helpers2.NewRandDenom()
				base := denom
				return denom, base, []string{}, fxtypes.ParseFxTarget("erc20"), denom
			},
		},
		{
			name: "base, empty alias, expect denom",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				denom := helpers2.NewRandDenom()
				return denom, "", []string{}, fxtypes.ParseFxTarget("eth"), denom
			},
		},
		{
			name: "base denom, math alias ibc",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				portID, channelID := suite.RandTransferChannel()
				ibcDenom := suite.AddIBCToken(portID, channelID)
				denom := helpers2.NewRandDenom()
				aliases := make([]string, 0)
				keepers := suite.CrossChainKeepers()
				for module := range keepers {
					aliases = append(aliases, fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module)))
				}
				base := denom
				return denom, base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(fmt.Sprintf("ibc/%s/px", strings.TrimPrefix(channelID, ibcchanneltypes.ChannelPrefix))), ibcDenom
			},
		},
		{
			name: "base denom, not alias ibc",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				_, channelID := suite.RandTransferChannel()
				ibcDenom := fmt.Sprintf("ibc/%s", strings.ToUpper(hex.EncodeToString(tmrand.Bytes(32))))
				denom := helpers2.NewRandDenom()
				aliases := make([]string, 0)
				keepers := suite.CrossChainKeepers()
				for module := range keepers {
					aliases = append(aliases, fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module)))
				}
				base := denom
				return denom, base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(fmt.Sprintf("ibc/%s/px", strings.TrimPrefix(channelID, ibcchanneltypes.ChannelPrefix))), denom
			},
		},
		{
			name: "base denom, math alias, expected ibc",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				portID, channelID := suite.RandTransferChannel()
				ibcDenom := suite.AddIBCToken(portID, channelID)
				denom := helpers2.NewRandDenom()
				keepers := suite.CrossChainKeepers()
				i, idx, idxModule, idxDenom := 0, tmrand.Intn(len(keepers)), "", ""
				aliases := make([]string, 0)
				for module := range keepers {
					randToken := fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module))
					aliases = append(aliases, randToken)
					if i == idx {
						idxModule = module
						idxDenom = randToken
					}
					i++
				}
				base := denom
				return denom, base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(idxModule), idxDenom
			},
		},
		{
			name: "base denom, not math alias",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				portID, channelID := suite.RandTransferChannel()
				ibcDenom := suite.AddIBCToken(portID, channelID)
				denom := helpers2.NewRandDenom()
				keepers := suite.CrossChainKeepers()
				i, idx, idxModule := 0, tmrand.Intn(len(keepers)), ""
				aliases := make([]string, 0)
				for module := range keepers {
					if i == idx {
						idxModule = module
					} else {
						randToken := fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module))
						aliases = append(aliases, randToken)
					}
					i++
				}
				base := denom
				return denom, base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(idxModule), denom
			},
		},
		{
			name: "alias denom, math alias ibc",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				portID, channelID := suite.RandTransferChannel()
				ibcDenom := suite.AddIBCToken(portID, channelID)
				aliases := make([]string, 0)
				keepers := suite.CrossChainKeepers()
				for module := range keepers {
					aliases = append(aliases, fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module)))
				}
				base := helpers2.NewRandDenom()
				return aliases[0], base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(fmt.Sprintf("ibc/%s/px", strings.TrimPrefix(channelID, ibcchanneltypes.ChannelPrefix))), ibcDenom
			},
		},
		{
			name: "alias denom, not alias ibc",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				_, channelID := suite.RandTransferChannel()
				ibcDenom := fmt.Sprintf("ibc/%s", strings.ToUpper(hex.EncodeToString(tmrand.Bytes(32))))
				aliases := make([]string, 0)
				keepers := suite.CrossChainKeepers()

				for module := range keepers {
					aliases = append(aliases, fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module)))
				}
				base := helpers2.NewRandDenom()
				return aliases[0], base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(fmt.Sprintf("ibc/%s/px", strings.TrimPrefix(channelID, ibcchanneltypes.ChannelPrefix))), aliases[0]
			},
		},
		{
			name: "alias denom, math alias, expected ibc",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				portID, channelID := suite.RandTransferChannel()
				ibcDenom := suite.AddIBCToken(portID, channelID)
				keepers := suite.CrossChainKeepers()

				i, idx, idxModule, idxDenom := 0, tmrand.Intn(len(keepers)), "", ""
				if idx == 0 {
					idx = 1
				}
				aliases := make([]string, 0)
				for module := range keepers {
					randToken := fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module))
					aliases = append(aliases, randToken)
					if i == idx {
						idxModule = module
						idxDenom = randToken
					}
					i++
				}
				base := helpers2.NewRandDenom()
				return aliases[0], base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(idxModule), idxDenom
			},
		},
		{
			name: "alias denom, not math alias",
			malleate: func() (string, string, []string, fxtypes.FxTarget, string) {
				portID, channelID := suite.RandTransferChannel()
				ibcDenom := suite.AddIBCToken(portID, channelID)
				keepers := suite.CrossChainKeepers()

				i, idx, idxModule := 0, tmrand.Intn(len(keepers)), ""
				if idx == 0 {
					idx = 1
				}
				aliases := make([]string, 0)
				for module := range keepers {
					if i == idx {
						idxModule = module
					} else {
						randToken := fmt.Sprintf("%s%s", module, helpers2.GenerateAddressByModule(module))
						aliases = append(aliases, randToken)
					}
					i++
				}
				base := helpers2.NewRandDenom()
				return aliases[0], base, append(aliases, ibcDenom), fxtypes.ParseFxTarget(idxModule), aliases[0]
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			denom, base, aliases, fxTarget, expDenom := tc.malleate()
			targetDenom := suite.app.Erc20Keeper.ToTargetDenom(suite.ctx, denom, base, aliases, fxTarget)
			suite.Require().EqualValues(expDenom, targetDenom)
		})
	}
}

func (suite *KeeperTestSuite) TestConvertDenomToTarget() {
	testCases := []struct {
		name     string
		malleate func(acc sdk.AccAddress) (originCoin sdk.Coin, expCoin sdk.Coin, target fxtypes.FxTarget, errArgs []string)
		expPass  bool
		expErr   func(args []string) string
	}{
		{
			name: "ok - DefaultDenom, not convert",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))

				originCoin := sdk.NewCoin(fxtypes.DefaultDenom, amt)
				expCoin := originCoin
				fxTarget := fxtypes.ParseFxTarget("")

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(originCoin))

				return originCoin, expCoin, fxTarget, nil
			},
			expPass: true,
		},
		{
			name: "ok - register denom and not have alias, not convert",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()
				mdmd := md.GetMetadata()
				mdmd.DenomUnits[0].Aliases = []string{}

				pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, mdmd)
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(pair.GetDenom(), amt)
				expCoin := originCoin
				fxTarget := fxtypes.ParseFxTarget("")

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(originCoin))

				return originCoin, expCoin, fxTarget, nil
			},
			expPass: true,
		},
		{
			name: "ok - register denom, have alias, coin denom equal target coin",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				expCoin := originCoin
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[0]) // or empty

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(originCoin))

				return originCoin, expCoin, fxTarget, nil
			},
			expPass: true,
		},
		{
			name: "failed - register denom, have alias, base to alias, insufficient funds",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(pair.GetDenom(), amt.Add(sdk.NewInt(1)))
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt.Add(sdk.NewInt(1)))
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[0])

				mintAmt := sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), amt))
				helpers2.AddTestAddr(suite.app, suite.ctx, acc, mintAmt)

				return originCoin, expCoin, fxTarget, []string{mintAmt.String(), originCoin.String()}
			},
			expPass: false,
			expErr: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
		},
		{
			name: "failed - register denom, have alias, base to alias, module insufficient funds",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(pair.GetDenom(), amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[0])

				mintAmt := sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), amt))
				helpers2.AddTestAddr(suite.app, suite.ctx, acc, mintAmt)

				return originCoin, expCoin, fxTarget, []string{sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], sdk.NewInt(0)).String(), expCoin.String()}
			},
			expPass: false,
			expErr: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
		},
		{
			name: "ok - register denom, have alias, base to alias",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(pair.GetDenom(), amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[0])

				mintAmt := sdk.NewCoins(sdk.NewCoin(pair.GetDenom(), amt))
				helpers2.AddTestAddr(suite.app, suite.ctx, acc, mintAmt)

				// mint alias token to erc20 module
				err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(expCoin))
				suite.Require().NoError(err)

				return originCoin, expCoin, fxTarget, []string{}
			},
			expPass: true,
		},
		{
			name: "ok - not register",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				randDenom := helpers2.NewRandDenom()

				originCoin := sdk.NewCoin(randDenom, amt)
				expCoin := sdk.NewCoin(randDenom, amt)
				fxTarget := fxtypes.ParseFxTarget("erc20") // any target

				return originCoin, expCoin, fxTarget, []string{}
			},
			expPass: true,
		},
		{
			name: "ok - register alias, alias equal target coin",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[0]) // any target

				return originCoin, expCoin, fxTarget, []string{}
			},
			expPass: true,
		},
		{
			name: "failed - register alias, alias to base, insufficient funds",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				expCoin := sdk.NewCoin(md.GetMetadata().Base, amt)
				fxTarget := fxtypes.ParseFxTarget("erc20") // or empty

				return originCoin, expCoin, fxTarget, []string{sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], sdk.NewInt(0)).String(), originCoin.String()}
			},
			expPass: false,
			expErr: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
		},
		{
			name: "ok - register alias, alias to base",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()

				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				expCoin := sdk.NewCoin(md.GetMetadata().Base, amt)
				fxTarget := fxtypes.ParseFxTarget("erc20") // or empty

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(originCoin))

				return originCoin, expCoin, fxTarget, []string{}
			},
			expPass: true,
		},
		{
			name: "failed - register alias, alias to alias, module insufficient funds",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()
				for len(md.GetMetadata().DenomUnits[0].Aliases) <= 1 {
					md = suite.GenerateCrossChainDenoms()
				}

				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[1], amt)
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[1])

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(originCoin))

				return originCoin, expCoin, fxTarget, []string{sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[1], sdk.NewInt(0)).String(), expCoin.String()}
			},
			expPass: false,
			expErr: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
		},
		{
			name: "ok - register alias, alias to alias",
			malleate: func(acc sdk.AccAddress) (sdk.Coin, sdk.Coin, fxtypes.FxTarget, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				md := suite.GenerateCrossChainDenoms()
				for len(md.GetMetadata().DenomUnits[0].Aliases) <= 1 {
					md = suite.GenerateCrossChainDenoms()
				}

				_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
				suite.Require().NoError(err)

				originCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[1], amt)
				fxTarget := fxtypes.ParseFxTarget(md.GetModules()[1])

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(originCoin))
				err = suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(expCoin))
				suite.Require().NoError(err)

				return originCoin, expCoin, fxTarget, []string{}
			},
			expPass: true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()

			originCoin, expCoin, fxTarget, errArgs := tc.malleate(signer.AccAddress())

			targetCoin, err := suite.app.Erc20Keeper.ConvertDenomToTarget(suite.ctx, signer.AccAddress(), originCoin, fxTarget)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().EqualValues(expCoin, targetCoin, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().EqualError(err, tc.expErr(errArgs), tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestConvertDenom() {
	testCases := []struct {
		name     string
		malleate func(md Metadata, acc, rec sdk.AccAddress) (receiver string, coin, expCoin sdk.Coin, targetStr string, errArgs []string)
		expPass  bool
		expErr   func(args []string) string
	}{
		{
			name: "failed - convert to source denom",
			malleate: func(md Metadata, acc, rec sdk.AccAddress) (string, sdk.Coin, sdk.Coin, string, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				coin := sdk.NewCoin(md.GetMetadata().Base, amt)
				return acc.String(), coin, coin, "", []string{coin.Denom}
			},
			expPass: false,
			expErr: func(args []string) string {
				return fmt.Sprintf("convert to source denom: %s: invalid denom", args[0])
			},
		},
		{
			name: "ok - base to alias",
			malleate: func(md Metadata, acc, rec sdk.AccAddress) (string, sdk.Coin, sdk.Coin, string, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				coin := sdk.NewCoin(md.GetMetadata().Base, amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(coin))
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(expCoin))
				suite.Require().NoError(err)

				return acc.String(), coin, expCoin, md.GetModules()[0], []string{}
			},
			expPass: true,
		},
		{
			name: "ok - base to alias - sender not equal receiver",
			malleate: func(md Metadata, acc, rec sdk.AccAddress) (string, sdk.Coin, sdk.Coin, string, []string) {
				amt := sdk.NewInt(int64(tmrand.Uint32() + 1000))
				coin := sdk.NewCoin(md.GetMetadata().Base, amt)
				expCoin := sdk.NewCoin(md.GetMetadata().DenomUnits[0].Aliases[0], amt)

				helpers2.AddTestAddr(suite.app, suite.ctx, acc, sdk.NewCoins(coin))
				err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, sdk.NewCoins(expCoin))
				suite.Require().NoError(err)

				return rec.String(), coin, expCoin, md.GetModules()[0], []string{}
			},
			expPass: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			receive := suite.RandSigner()

			md := suite.GenerateCrossChainDenoms()
			_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
			suite.Require().NoError(err)

			receiveAddr, coin, expCoin, targetStr, errArgs := tc.malleate(md, signer.AccAddress(), receive.AccAddress())

			msg := types.MsgConvertDenom{
				Sender:   signer.AccAddress().String(),
				Receiver: receiveAddr,
				Coin:     coin,
				Target:   targetStr,
			}

			coinBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), coin.Denom)
			addr, _ := sdk.AccAddressFromBech32(msg.Receiver)

			expCoinBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr, expCoin.Denom)
			_, err = suite.app.Erc20Keeper.ConvertDenom(sdk.WrapSDKContext(suite.ctx), &msg)

			afterCoinBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), coin.Denom)
			afterExpCoinBalance := suite.app.BankKeeper.GetBalance(suite.ctx, addr, expCoin.Denom)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().EqualValues(coinBalance.Sub(afterCoinBalance).Amount, afterExpCoinBalance.Sub(expCoinBalance).Amount)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().EqualError(err, tc.expErr(errArgs), tc.name)
			}
		})
	}
}
