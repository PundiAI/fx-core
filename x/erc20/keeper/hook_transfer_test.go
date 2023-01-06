package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookTransferNativeToken() {
	nativePairFn := func() (types.TokenPair, banktypes.Metadata) {
		denoms := suite.GenerateCrossChainDenoms()
		return suite.DeployNativeRelayToken("TEST", denoms...)
	}
	randTotalMint := func(md banktypes.Metadata) *big.Int {
		amt := big.NewInt(int64(tmrand.Uint32() + 1))
		return suite.MintLockNativeTokenToModule(md, sdk.NewIntFromBigInt(amt))
	}

	testCases := []struct {
		name   string
		pair   func() (types.TokenPair, banktypes.Metadata)
		relays func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string)
		error  func(args []string) string
		result bool
	}{
		{
			name: "ok - transfer to module",
			pair: nativePairFn,
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)
				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, totalCanMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalCanMint)
				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        totalCanMint,
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relay}, []string{}
			},
			result: true,
		},
		{
			name: "ok - transfer multiple to module",
			pair: nativePairFn,
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)
				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, totalCanMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalCanMint)
				relayAmt1 := big.NewInt(0).Div(totalCanMint, big.NewInt(2))
				relayAmt2 := big.NewInt(0).Sub(totalCanMint, relayAmt1)
				relayEvent1 := types.RelayTransfer{
					From:          singerAddr,
					Amount:        relayAmt1,
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				relayEvent2 := types.RelayTransfer{
					From:          singerAddr,
					Amount:        relayAmt2,
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relayEvent1, relayEvent2}, []string{}
			},
			result: true,
		},
		{
			name: "failed - burn amount exceeds balance",
			pair: nativePairFn,
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)
				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, totalCanMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalCanMint)
				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        big.NewInt(0).Add(totalCanMint, big.NewInt(1)),
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relay}, []string{pair.Erc20Address}
			},
			error: func(args []string) string {
				return fmt.Sprintf("contract call failed: method 'burn', contract '%s': execution reverted: burn amount exceeds balance", args[0])
			},
			result: false,
		},
		{
			name: "failed - module insufficient funds",
			pair: nativePairFn,
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)
				// mint more than total can mint
				moreTotalMint := big.NewInt(0).Add(totalCanMint, big.NewInt(1))
				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, moreTotalMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, moreTotalMint)

				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        big.NewInt(0).Add(totalCanMint, big.NewInt(1)),
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relay}, []string{
					sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(totalCanMint)).String(),
					sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(moreTotalMint)).String(),
				}
			},
			error: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
			result: false,
		},
		{
			name: "ok - fx - transfer to module",
			pair: func() (types.TokenPair, banktypes.Metadata) {
				return suite.DeployFXRelayToken()
			},
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)
				// transfer FX to contract address
				coin := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(totalCanMint))
				err := suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, pair.GetERC20Contract().Bytes(), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, totalCanMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalCanMint)
				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        totalCanMint,
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relay}, []string{}
			},
			result: true,
		},
		{
			name: "failed - fx - contact insufficient funds",
			pair: func() (types.TokenPair, banktypes.Metadata) {
				return suite.DeployFXRelayToken()
			},
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)

				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, totalCanMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalCanMint)

				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        totalCanMint,
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relay}, []string{
					sdk.NewCoin(pair.Denom, sdk.NewInt(0)).String(),
					sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(totalCanMint)).String(),
				}
			},
			error: func(args []string) string {
				return fmt.Sprintf("%s is smaller than %s: insufficient funds", args[0], args[1])
			},
			result: false,
		},
		{
			name: "failed - undefined owner",
			pair: nativePairFn,
			relays: func(md banktypes.Metadata, pair types.TokenPair, singerAddr common.Address) ([]types.RelayTransfer, []string) {
				totalCanMint := randTotalMint(md)
				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, totalCanMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalCanMint)
				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        totalCanMint,
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: types.OWNER_UNSPECIFIED,
				}
				return []types.RelayTransfer{relay}, []string{}
			},
			error: func(args []string) string {
				return "undefined owner of contract pair"
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			pair, md := tc.pair()
			// relay event
			relays, errArgs := tc.relays(md, pair, signer.Address())
			// hook transfer
			beforeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
			beforeModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.Erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
			beforeContractBalance := suite.app.BankKeeper.GetBalance(suite.ctx, pair.GetERC20Contract().Bytes(), pair.GetDenom())
			err := suite.app.Erc20Keeper.EVMHooks().HookTransferEvent(suite.ctx, relays)
			// check result
			if tc.result {
				// check error
				suite.Require().NoError(err)
				// check balance
				relayAmt := big.NewInt(0)
				for _, r := range relays {
					relayAmt = relayAmt.Add(relayAmt, r.Amount)
				}
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
				suite.Require().Equal(relayAmt, balance.Sub(beforeBalance).Amount.BigInt())
				// check module and contract balance
				if pair.GetDenom() == fxtypes.DefaultDenom {
					contractBalance := suite.app.BankKeeper.GetBalance(suite.ctx, pair.GetERC20Contract().Bytes(), pair.GetDenom())
					suite.Require().Equal(relayAmt, beforeContractBalance.Sub(contractBalance).Amount.BigInt())
				} else {
					moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.Erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
					suite.Require().Equal(relayAmt, beforeModuleBalance.Sub(moduleBalance).Amount.BigInt())
				}
			} else {
				// check error
				suite.Require().Error(err)
				// check error msg
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}
