package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
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
				_, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				// Mint 10 tokens
				suite.MintERC20Token(suite.signer, contractAddr, suite.signer.Address(), big.NewInt(10))

				// Burn the 10 tokens
				suite.TransferERC20TokenToModule(suite.signer, contractAddr, big.NewInt(10))
			},
			true,
		},
		{
			"unregistered pair",
			func(contractAddr common.Address) {
				// Mint 10 tokens
				suite.MintERC20Token(suite.signer, contractAddr, suite.signer.Address(), big.NewInt(10))

				// Burn the 10 tokens
				suite.TransferERC20TokenToModule(suite.signer, contractAddr, big.NewInt(10))
			},
			false,
		},
		{
			"wrong event",
			func(contractAddr common.Address) {
				_, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				// Mint 10 tokens
				suite.MintERC20Token(suite.signer, contractAddr, suite.signer.Address(), big.NewInt(10))
			},
			false,
		},
		{
			"Pair is incorrectly loaded",
			func(contractAddr common.Address) {
				pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				suite.app.Erc20Keeper.RemoveTokenPair(suite.ctx, *pair)

				// Mint 10 tokens
				suite.MintERC20Token(suite.signer, contractAddr, suite.signer.Address(), big.NewInt(10))

				// Burn the 10 tokens
				suite.TransferERC20TokenToModule(suite.signer, contractAddr, big.NewInt(10))
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			contractAddr, err := suite.DeployContract(suite.signer.Address())
			suite.Require().NoError(err)

			tc.malleate(contractAddr)

			balance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.signer.Address().Bytes(), "test")
			if tc.result {
				// Check if the execution was successful
				suite.Require().Equal(int64(10), balance.Amount.Int64())
			} else {
				// Check that no changes were made to the account
				suite.Require().Equal(int64(0), balance.Amount.Int64())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestEvmHooksRegisterCoin() {
	testCases := []struct {
		name      string
		mint      int64
		burn      int64
		reconvert int64
		result    bool
	}{
		{"correct execution", 100, 10, 5, true},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			metadata, pair := suite.setupRegisterCoin()
			suite.Require().NotNil(metadata)
			suite.Require().NotNil(pair)

			sender := sdk.AccAddress(suite.signer.Address().Bytes())

			coins := sdk.NewCoins(sdk.NewCoin(metadata.Base, sdk.NewInt(tc.mint)))
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, sender, coins))

			_, err := suite.app.Erc20Keeper.ConvertCoin(sdk.WrapSDKContext(suite.ctx),
				types.NewMsgConvertCoin(
					sdk.NewCoin(metadata.Base, sdk.NewInt(tc.burn)),
					suite.signer.Address(),
					sender,
				),
			)
			suite.Require().NoError(err, tc.name)

			balance := suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
			cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)
			suite.Require().Equal(cosmosBalance.Amount.Int64(), sdk.NewInt(tc.mint-tc.burn).Int64())
			suite.Require().Equal(balance, big.NewInt(tc.burn))

			suite.TransferERC20TokenToModule(suite.signer, pair.GetERC20Contract(), big.NewInt(tc.reconvert))
			balance = suite.BalanceOf(pair.GetERC20Contract(), suite.signer.Address())
			cosmosBalance = suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)

			if tc.result {
				// Check if the execution was successful
				suite.Require().NoError(err)
				suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.mint-tc.burn+tc.reconvert))
				suite.Require().Equal(balance, big.NewInt(tc.burn-tc.reconvert))
			} else {
				// Check that no changes were made to the account
				suite.Require().Error(err)
				suite.Require().Equal(cosmosBalance.Amount, sdk.NewInt(tc.mint-tc.burn))
				suite.Require().Equal(balance, big.NewInt(tc.burn))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestPostTxProcessing() {
	erc20ModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
	msg := ethtypes.NewMessage(
		erc20ModuleAddr,
		&common.Address{},
		0,
		big.NewInt(0), // amount
		uint64(0),     // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		[]byte{},
		ethtypes.AccessList{}, // AccessList
		true,                  // checkNonce
	)

	account := helpers.GenerateAddress()

	transferData := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	transferData[31] = uint8(10)
	erc20 := fxtypes.GetERC20().ABI

	transferEvent := erc20.Events["Transfer"]

	testCases := []struct {
		name string
		test func()
	}{
		{
			"correct transfer (non burn)",
			func() {
				contractAddr, err := suite.DeployContract(suite.signer.Address())
				suite.Require().NoError(err)

				_, err = suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				topics := []common.Hash{transferEvent.ID, account.Hash(), account.Hash()}
				log := ethtypes.Log{
					Topics:  topics,
					Data:    transferData,
					Address: contractAddr,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err = suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"correct burn",
			func() {
				contractAddr, err := suite.DeployContract(suite.signer.Address())
				suite.Require().NoError(err)

				pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				topics := []common.Hash{transferEvent.ID, account.Hash(), erc20ModuleAddr.Hash()}
				log := ethtypes.Log{
					Topics:  topics,
					Data:    transferData,
					Address: contractAddr,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err = suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
				sender := sdk.AccAddress(account.Bytes())
				cosmosBalance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, pair.Denom)

				transferEvent, err := erc20.Unpack("Transfer", transferData)
				suite.Require().NoError(err)

				tokens, _ := transferEvent[0].(*big.Int)
				suite.Require().Equal(cosmosBalance.Amount.String(), tokens.String())
			},
		},
		{
			"Unspecified Owner",
			func() {
				contractAddr, err := suite.DeployContract(suite.signer.Address())
				suite.Require().NoError(err)

				pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				pair.ContractOwner = types.OWNER_UNSPECIFIED
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, *pair)

				topics := []common.Hash{transferEvent.ID, account.Hash(), erc20ModuleAddr.Hash()}
				log := ethtypes.Log{
					Topics:  topics,
					Data:    transferData,
					Address: contractAddr,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err = suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().Error(err)
			},
		},
		{
			"Fail Evm",
			func() {
				contractAddr, err := suite.DeployContract(suite.signer.Address())
				suite.Require().NoError(err)

				pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contractAddr)
				suite.Require().NoError(err)

				pair.ContractOwner = types.OWNER_MODULE
				suite.app.Erc20Keeper.SetTokenPair(suite.ctx, *pair)

				topics := []common.Hash{transferEvent.ID, account.Hash(), erc20ModuleAddr.Hash()}
				log := ethtypes.Log{
					Topics:  topics,
					Data:    transferData,
					Address: contractAddr,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err = suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().Error(err)
			},
		},
		{
			"No log address",
			func() {
				topics := []common.Hash{transferEvent.ID, account.Hash(), erc20ModuleAddr.Hash()}
				log := ethtypes.Log{
					Topics: topics,
					Data:   transferData,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err := suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"No data on topic",
			func() {
				topics := []common.Hash{transferEvent.ID}
				log := ethtypes.Log{
					Topics: topics,
					Data:   transferData,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err := suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"Empty logs",
			func() {
				log := ethtypes.Log{}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err := suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"No log data",
			func() {
				topics := []common.Hash{transferEvent.ID, account.Hash(), erc20ModuleAddr.Hash()}
				log := ethtypes.Log{
					Topics: topics,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err := suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"Non transfer event",
			func() {
				approvalEvent := erc20.Events["Approval"]
				topics := []common.Hash{approvalEvent.ID, account.Hash(), account.Hash()}
				log := ethtypes.Log{
					Topics: topics,
					Data:   transferData,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err := suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"Non recognized event",
			func() {
				topics := []common.Hash{{}, account.Hash(), account.Hash()}
				log := ethtypes.Log{
					Topics: topics,
					Data:   transferData,
				}
				receipt := &ethtypes.Receipt{
					Logs: []*ethtypes.Log{&log},
				}

				err := suite.app.Erc20Keeper.EVMHooks().PostTxProcessing(suite.ctx, msg, receipt)
				suite.Require().NoError(err)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.test()
		})
	}
}
