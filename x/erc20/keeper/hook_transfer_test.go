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
	testCases := []struct {
		name   string
		pair   func() (types.TokenPair, banktypes.Metadata)
		relays func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string)
		error  func(args []string) string
		result bool
	}{
		{
			name: "ok - transfer to module",
			pair: nativePairFn,
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
				// mint more than total can mint
				moreTotalMint := big.NewInt(0).Add(totalCanMint, big.NewInt(1))
				suite.ModuleMintERC20Token(pair.GetERC20Contract(), singerAddr, moreTotalMint)
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, moreTotalMint)

				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        moreTotalMint,
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
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			// mint lock token
			totalMint := suite.MintLockNativeTokenToModule(md, sdk.NewIntFromBigInt(big.NewInt(int64(tmrand.Uint32()+1))))
			// relay event
			relays, errArgs := tc.relays(pair, signer.Address(), totalMint)
			// hook transfer
			beforeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
			beforeModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.Erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
			beforeContractBalance := suite.app.BankKeeper.GetBalance(suite.ctx, pair.GetERC20Contract().Bytes(), pair.GetDenom())
			err := suite.app.Erc20Keeper.EVMHooks().HookTransferEvent(suite.ctx, relays)
			// check result
			if tc.result {
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
				suite.Require().Error(err)
				// check error msg
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHookTransferERC20Token() {
	erc20PairFn := func() (types.TokenPair, banktypes.Metadata) {
		contract, err := suite.DeployContract(suite.signer.Address())
		suite.Require().NoError(err)
		return suite.DeployERC20RelayToken(contract)
	}
	randTotalMint := func() *big.Int {
		return big.NewInt(int64(tmrand.Uint32() + 1))
	}

	testCases := []struct {
		name   string
		pair   func() (types.TokenPair, banktypes.Metadata)
		relays func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string)
		error  func(args []string) string
		result bool
	}{
		{
			name: "ok - transfer to module",
			pair: erc20PairFn,
			relays: func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string) {
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalMint)
				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        totalMint,
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
			pair: erc20PairFn,
			relays: func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string) {
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalMint)
				relayAmt1 := big.NewInt(0).Div(totalMint, big.NewInt(2))
				relayAmt2 := big.NewInt(0).Sub(totalMint, relayAmt1)
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
			name: "success - transfer amount small than mint denom",
			pair: erc20PairFn,
			relays: func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string) {
				suite.TransferERC20TokenToModuleWithoutHook(pair.GetERC20Contract(), singerAddr, totalMint)
				relay := types.RelayTransfer{
					From:          singerAddr,
					Amount:        big.NewInt(0).Add(totalMint, big.NewInt(1)),
					TokenContract: pair.GetERC20Contract(),
					Denom:         pair.Denom,
					ContractOwner: pair.ContractOwner,
				}
				return []types.RelayTransfer{relay}, []string{pair.Erc20Address}
			},
			result: true,
		},
		{
			name: "failed - undefined owner",
			pair: erc20PairFn,
			relays: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			pair, _ := tc.pair()
			// mint token
			totalMint := randTotalMint()
			suite.MintERC20Token(pair.GetERC20Contract(), suite.signer.Address(), signer.Address(), totalMint)
			// relay event
			relays, errArgs := tc.relays(pair, signer.Address(), totalMint)
			// hook transfer
			beforeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
			err := suite.app.Erc20Keeper.EVMHooks().HookTransferEvent(suite.ctx, relays)
			// check result
			if tc.result {
				suite.Require().NoError(err)
				// check balance
				relayAmt := big.NewInt(0)
				for _, r := range relays {
					relayAmt = relayAmt.Add(relayAmt, r.Amount)
				}
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
				suite.Require().Equal(relayAmt, balance.Sub(beforeBalance).Amount.BigInt())
				moduleBalanceERC20 := suite.BalanceOf(pair.GetERC20Contract(), suite.app.Erc20Keeper.ModuleAddress())

				if relayAmt.Cmp(moduleBalanceERC20) == 1 {
					suite.Require().Equal(big.NewInt(0).Sub(relayAmt, big.NewInt(1)), moduleBalanceERC20)
				} else {
					suite.Require().Equal(relayAmt, moduleBalanceERC20)
				}
			} else {
				suite.Require().Error(err)
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHookTransfer() {
	signer := suite.RandSigner()

	nativePair, nativeMetadata := suite.DeployNativeRelayToken("ABC", suite.GenerateCrossChainDenoms()...)

	nativeTotalMint := suite.MintLockNativeTokenToModule(nativeMetadata, sdk.NewIntFromBigInt(big.NewInt(int64(tmrand.Uint32()+1))))
	suite.ModuleMintERC20Token(nativePair.GetERC20Contract(), signer.Address(), nativeTotalMint)
	suite.TransferERC20TokenToModuleWithoutHook(nativePair.GetERC20Contract(), signer.Address(), nativeTotalMint)
	nativeRelay := types.RelayTransfer{
		From:          signer.Address(),
		Amount:        nativeTotalMint,
		TokenContract: nativePair.GetERC20Contract(),
		Denom:         nativePair.Denom,
		ContractOwner: nativePair.ContractOwner,
	}

	erc20Contract, err := suite.DeployContract(suite.signer.Address())
	suite.Require().NoError(err)
	erc20Pair, _ := suite.DeployERC20RelayToken(erc20Contract)

	erc20TotalMint := big.NewInt(int64(tmrand.Uint32() + 1))
	suite.MintERC20Token(erc20Pair.GetERC20Contract(), suite.signer.Address(), signer.Address(), erc20TotalMint)
	suite.TransferERC20TokenToModuleWithoutHook(erc20Pair.GetERC20Contract(), signer.Address(), erc20TotalMint)
	erc20Relay := types.RelayTransfer{
		From:          signer.Address(),
		Amount:        erc20TotalMint,
		TokenContract: erc20Pair.GetERC20Contract(),
		Denom:         erc20Pair.Denom,
		ContractOwner: erc20Pair.ContractOwner,
	}

	beforeNativeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), nativePair.GetDenom())
	beforeErc20Balance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), erc20Pair.GetDenom())

	err = suite.app.Erc20Keeper.EVMHooks().HookTransferEvent(suite.ctx, []types.RelayTransfer{nativeRelay, erc20Relay})
	suite.Require().NoError(err)

	nativeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), nativePair.GetDenom())
	erc20Balance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), erc20Pair.GetDenom())

	suite.Require().Equal(nativeBalance.Sub(beforeNativeBalance).Amount.BigInt(), nativeTotalMint)
	suite.Require().Equal(erc20Balance.Sub(beforeErc20Balance).Amount.BigInt(), erc20TotalMint)
}
