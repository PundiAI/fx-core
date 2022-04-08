package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

func (suite *KeeperTestSuite) TestEvmHooksRegisterFIP20() {
	testCases := []struct {
		name     string
		malleate func(common.Address)
		result   bool
	}{
		{
			"correct execution",
			func(contractAddr common.Address) {
				// pair := types.NewTokenPair(contractAddr, "coinevm", true, types.OWNER_MODULE)
				_, err := suite.app.IntrarelayerKeeper.RegisterFIP20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				accAddress := sdk.AccAddress(suite.priKey.PubKey().Address())
				hexAddress := common.BytesToAddress(accAddress)

				// Mint 10 tokens to suite.address (owner)
				_ = suite.MintFIP20Token(contractAddr, hexAddress, hexAddress, big.NewInt(10))
				suite.Commit()

				// Burn the 10 tokens of suite.address (owner)
				suite.BurnFIP20Token(contractAddr, hexAddress, big.NewInt(10))
			},
			true,
		},
		{
			"unregistered pair",
			func(contractAddr common.Address) {
				// Mint 10 tokens to suite.address (owner)
				_ = suite.MintFIP20Token(contractAddr, suite.address, suite.address, big.NewInt(10))
				suite.Commit()

				// Burn the 10 tokens of suite.address (owner)
				_ = suite.BurnFIP20Token(contractAddr, suite.address, big.NewInt(10))
			},
			false,
		},
		{
			"wrong event",
			func(contractAddr common.Address) {
				_, err := suite.app.IntrarelayerKeeper.RegisterFIP20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				// Mint 10 tokens to suite.address (owner)
				_ = suite.MintFIP20Token(contractAddr, suite.address, suite.address, big.NewInt(10))
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()

			contractAddr := suite.DeployContract(suite.address, "coin", "token", 18)
			suite.Commit()

			tc.malleate(contractAddr)
			suite.Commit()

			balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(suite.priKey.PubKey().Address()), types.CreateDenom(contractAddr.String()))
			suite.Commit()
			if tc.result {
				// Check if the execution was successfull
				suite.Require().Equal(balance.Amount, sdk.NewInt(10))
			} else {
				// Check that no changes were made to the account
				suite.Require().Equal(balance.Amount, sdk.NewInt(0))
			}
		})
	}
	suite.mintFeeCollector = false
}

func (suite *KeeperTestSuite) TestEvmHooksRegisterCoin() {
	testCases := []struct {
		name      string
		mint      int64
		burn      int64
		reconvert int64

		result bool
	}{
		{name: "correct execution", mint: 100, burn: 10, reconvert: 5, result: true},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			accAddress := sdk.AccAddress(suite.priKey.PubKey().Address())
			hexAddress := common.BytesToAddress(accAddress)

			contractAddr := common.HexToAddress(pair.Fip20Address)

			// 1. mint token to accAddress
			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.mint)))
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, accAddress, coins))

			// 2. accAddress convert token to 0xAddress
			_, err := suite.app.IntrarelayerKeeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx), types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.burn)),
				hexAddress,
				accAddress,
			))
			suite.Require().NoError(err, tc.name)
			suite.Commit()

			balance := suite.BalanceOf(common.HexToAddress(pair.Fip20Address), hexAddress)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, accAddress, metadata.Base)
			suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
			suite.Require().Equal(balance, big.NewInt(tc.burn))

			// relay the 5 tokens of suite.address (owner)
			suite.BurnFIP20Token(contractAddr, hexAddress, big.NewInt(tc.reconvert))

			balance = suite.BalanceOf(common.HexToAddress(pair.Fip20Address), hexAddress)
			cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, accAddress, metadata.Base)

			if tc.result {
				// Check if the execution was successfull
				suite.Require().NoError(err)
				suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.mint-tc.burn+tc.reconvert))
			} else {
				// Check that no changes were made to the account
				suite.Require().Error(err)
				suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.mint-tc.burn))
			}
		})
	}
	suite.mintFeeCollector = false
}
