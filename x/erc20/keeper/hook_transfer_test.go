package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookTransferNativeToken() {
	testCases := []struct {
		name     string
		malleate func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - transfer to module",
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
				// todo execution failed detail
				return "execution reverted: evm transaction execution failed"
			},
			result: false,
		},
		{
			name: "failed - module insufficient funds",
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			name: "failed - undefined owner",
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			md := suite.GenerateCrossChainDenoms()
			pair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, md.GetMetadata())
			suite.NoError(err)
			// mint lock token
			randMint := big.NewInt(int64(tmrand.Uint32() + 10))
			totalMint := suite.MintLockNativeTokenToModule(md.GetMetadata(), sdk.NewIntFromBigInt(randMint))
			// relay event
			relays, errArgs := tc.malleate(*pair, signer.Address(), totalMint)
			// hook transfer
			beforeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
			beforeModuleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.Erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
			err = suite.app.Erc20Keeper.EVMHooks().HookTransferEvent(suite.ctx, relays)
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
				moduleBalance := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.Erc20Keeper.ModuleAddress().Bytes(), pair.GetDenom())
				suite.Require().Equal(relayAmt, beforeModuleBalance.Sub(moduleBalance).Amount.BigInt())
			} else {
				suite.Require().Error(err)
				// check error msg
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHookTransferFX() {
	testCases := []struct {
		name     string
		malleate func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - fx - transfer to module",
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			signer := suite.RandSigner()
			// token pair
			md, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, fxtypes.DefaultDenom)
			suite.Require().True(found)
			pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, fxtypes.DefaultDenom)
			suite.Require().True(found)

			// mint lock token
			totalMint := suite.MintLockNativeTokenToModule(md, sdk.NewIntFromBigInt(big.NewInt(int64(tmrand.Uint32()+1))))
			// relay event
			relays, errArgs := tc.malleate(pair, signer.Address(), totalMint)
			// hook transfer
			beforeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
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
				contractBalance := suite.app.BankKeeper.GetBalance(suite.ctx, pair.GetERC20Contract().Bytes(), pair.GetDenom())
				suite.Require().Equal(relayAmt, beforeContractBalance.Sub(contractBalance).Amount.BigInt())
			} else {
				suite.Require().Error(err)
				// check error msg
				suite.Require().EqualError(err, tc.error(errArgs))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestHookTransferERC20Token() {
	testCases := []struct {
		name     string
		malleate func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string)
		error    func(args []string) string
		result   bool
	}{
		{
			name: "ok - transfer to module",
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string) {
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
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string) {
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
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalMint *big.Int) ([]types.RelayTransfer, []string) {
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
			malleate: func(pair types.TokenPair, singerAddr common.Address, totalCanMint *big.Int) ([]types.RelayTransfer, []string) {
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
			contract, err := suite.DeployContract(suite.signer.Address())
			suite.Require().NoError(err)
			pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, contract)
			suite.Require().NoError(err)
			_, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, pair.Denom)
			suite.Require().True(found)
			// mint token
			totalMint := big.NewInt(int64(tmrand.Uint32() + 10))
			suite.MintERC20Token(suite.signer, pair.GetERC20Contract(), signer.Address(), totalMint)
			// relay event
			relays, errArgs := tc.malleate(*pair, signer.Address(), totalMint)
			// hook transfer
			beforeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, signer.AccAddress(), pair.GetDenom())
			err = suite.app.Erc20Keeper.EVMHooks().HookTransferEvent(suite.ctx, relays)
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
	nativeMetadata := suite.GenerateCrossChainDenoms()
	nativePair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, nativeMetadata.GetMetadata())
	suite.NoError(err)

	nativeTotalMint := suite.MintLockNativeTokenToModule(nativeMetadata.GetMetadata(), sdk.NewIntFromBigInt(big.NewInt(int64(tmrand.Uint32()+1))))
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
	erc20Pair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, erc20Contract)
	suite.Require().NoError(err)
	_, found := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, erc20Pair.Denom)
	suite.Require().True(found)

	erc20TotalMint := big.NewInt(int64(tmrand.Uint32() + 1))
	suite.MintERC20Token(suite.signer, erc20Pair.GetERC20Contract(), signer.Address(), erc20TotalMint)
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
