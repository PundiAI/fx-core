package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

func (suite *KeeperTestSuite) TestEvmHooksRegisterERC20() {
	testCases := []struct {
		name     string
		malleate func(common.Address)
		result   bool
	}{
		{
			"correct execution",
			func(contractAddr common.Address) {
				// pair := types.NewTokenPair(contractAddr, "coinevm", true, types.OWNER_MODULE)
				_, err := suite.app.IntrarelayerKeeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				accAddress := sdk.AccAddress(suite.priKey.PubKey().Address())
				hexAddress := suite.address

				// Mint 10 tokens to suite.address (owner)
				_ = suite.MintERC20Token(contractAddr, hexAddress, hexAddress, big.NewInt(10))
				suite.Commit()

				// Burn the 10 tokens of suite.address (owner)
				suite.RelayERC20Token(contractAddr, hexAddress, common.BytesToAddress(accAddress.Bytes()), big.NewInt(10))
			},
			true,
		},
		{
			"unregistered pair",
			func(contractAddr common.Address) {
				// Mint 10 tokens to suite.address (owner)
				_ = suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(10))
				suite.Commit()

				// Burn the 10 tokens of suite.address (owner)
				msg := suite.BurnERC20Token(contractAddr, suite.address, big.NewInt(10))
				logs := suite.app.EvmKeeper.GetTxLogsTransient(msg.AsTransaction().Hash())

				// Since theres no pair registered, no coins should be minted
				err := suite.app.IntrarelayerKeeper.PostTxProcessing(suite.ctx, msg.AsTransaction().Hash(), logs)
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"wrong event",
			func(contractAddr common.Address) {
				_, err := suite.app.IntrarelayerKeeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				// Mint 10 tokens to suite.address (owner)
				msg := suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(10))
				logs := suite.app.EvmKeeper.GetTxLogsTransient(msg.AsTransaction().Hash())

				// No coins should be minted on cosmos after a mint of the erc20 token
				err = suite.app.IntrarelayerKeeper.PostTxProcessing(suite.ctx, msg.AsTransaction().Hash(), logs)
				suite.Require().NoError(err)
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()

			contractAddr := suite.DeployContract("coin", "token", 18)
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
		{"correct execution", 100, 10, 5, true},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			accAddress := sdk.AccAddress(suite.priKey.PubKey().Address())
			hexAddress := suite.address

			contractAddr := common.HexToAddress(pair.Erc20Address)

			coins := sdk.NewCoins(sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.mint)))
			suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
			suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, accAddress, coins)

			convertCoin := types.NewMsgConvertCoin(
				sdk.NewCoin(cosmosTokenName, sdk.NewInt(tc.burn)),
				hexAddress,
				accAddress,
			)

			ctx := sdk.WrapSDKContext(suite.ctx)
			_, err := suite.app.IntrarelayerKeeper.ConvertCoin(ctx, convertCoin)
			suite.Require().NoError(err, tc.name)
			suite.Commit()

			balance := suite.BalanceOf(common.HexToAddress(pair.Erc20Address), hexAddress)
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, accAddress, metadata.Base)
			suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
			suite.Require().Equal(balance, big.NewInt(tc.burn))

			// relay the 5 tokens of suite.address (owner)
			suite.RelayERC20Token(contractAddr, hexAddress, common.BytesToAddress(accAddress.Bytes()), big.NewInt(tc.reconvert))

			balance = suite.BalanceOf(common.HexToAddress(pair.Erc20Address), hexAddress)
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
