package ante_test

import (
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/server/config"

	"github.com/functionx/fx-core/v7/ante"
	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

func (suite *AnteTestSuite) TestRejectValidatorGranted() {
	testCases := []struct {
		name       string
		malleate   func(addr sdk.AccAddress, val sdk.ValAddress, privs []cryptotypes.PrivKey, accNums, accSeqs []uint64) sdk.Tx
		expectPass bool
	}{
		{
			name: "success",
			malleate: func(addr sdk.AccAddress, val sdk.ValAddress, privs []cryptotypes.PrivKey, accNums, accSeqs []uint64) sdk.Tx {
				testMsg := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), fxtypes.DefaultDenom, &testMsg)
				tx, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
				suite.Require().NoError(err)
				return tx
			},
			expectPass: true,
		},
		{
			name: "fail - disable address",
			malleate: func(addr sdk.AccAddress, valAddr sdk.ValAddress, privs []cryptotypes.PrivKey, accNums, accSeqs []uint64) sdk.Tx {
				testMsg := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				err := suite.app.StakingKeeper.DisableValidatorAddress(suite.ctx, valAddr)
				suite.Require().NoError(err)
				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), fxtypes.DefaultDenom, &testMsg)
				tx, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
				suite.Require().NoError(err)
				return tx
			},
			expectPass: false,
		},
		{
			name: "fail with multiple msgs",
			malleate: func(addr sdk.AccAddress, valAddr sdk.ValAddress, privs []cryptotypes.PrivKey, accNums, accSeqs []uint64) sdk.Tx {
				msg1 := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				msg2 := banktypes.MsgSend{
					FromAddress: addr.String(),
					ToAddress:   sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: fxtypes.DefaultDenom}},
				}
				err := suite.app.StakingKeeper.DisableValidatorAddress(suite.ctx, valAddr)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), fxtypes.DefaultDenom, &msg1, &msg2)
				tx, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
				suite.Require().NoError(err)
				return tx
			},
			expectPass: false,
		},
		{
			name: "eth - success",
			malleate: func(addr sdk.AccAddress, val sdk.ValAddress, privs []cryptotypes.PrivKey, accNums, accSeqs []uint64) sdk.Tx {
				from := common.BytesToAddress(addr.Bytes())
				to := helpers.GenerateAddress()
				emptyAccessList := ethtypes.AccessList{}

				helpers.AddTestAddr(suite.app, suite.ctx, addr, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e18))))

				msg := suite.BuildTestEthTx(from, to, nil, make([]byte, 0), big.NewInt(0), big.NewInt(500000000000), big.NewInt(50), &emptyAccessList)
				return suite.CreateTestTx(msg, privs[0], 0, false)
			},
			expectPass: true,
		},
		{
			name: "eth - failed tx",
			malleate: func(addr sdk.AccAddress, val sdk.ValAddress, privs []cryptotypes.PrivKey, accNums, accSeqs []uint64) sdk.Tx {
				from := common.BytesToAddress(addr.Bytes())
				to := helpers.GenerateAddress()
				emptyAccessList := ethtypes.AccessList{}

				err := suite.app.StakingKeeper.DisableValidatorAddress(suite.ctx, sdk.ValAddress(addr))
				suite.Require().NoError(err)

				msg := suite.BuildTestEthTx(from, to, nil, make([]byte, 0), big.NewInt(0), big.NewInt(500000000000), big.NewInt(50), &emptyAccessList)
				return suite.CreateTestTx(msg, privs[0], 0, false)
			},
			expectPass: false,
		},
	}

	ethSigDec := ante.NewEthSigVerificationDecorator(suite.app.EvmKeeper)
	dec := ante.NewSigVerificationDecorator(suite.app.AccountKeeper, app.MakeEncodingConfig().TxConfig.SignModeHandler())
	ethDec := ante.NewEthGasConsumeDecorator(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.EvmKeeper, config.DefaultMaxTxGasWanted)
	var err error
	for _, tc := range testCases {
		val := helpers.NewEthPrivKey()
		valAddr := sdk.ValAddress(val.PubKey().Address())
		addr := sdk.AccAddress(valAddr)
		account := authtypes.NewBaseAccount(addr, val.PubKey(), 0, 0)
		suite.app.AccountKeeper.SetAccount(suite.ctx, account)
		privs, accNums, accSeqs := []cryptotypes.PrivKey{val}, []uint64{0}, []uint64{0}

		suite.Run(tc.name, func() {
			tx := tc.malleate(addr, valAddr, privs, accNums, accSeqs)
			if strings.HasPrefix(tc.name, "eth") {
				handler := sdk.ChainAnteDecorators(ethSigDec, ethDec)
				_, err = handler(suite.ctx, tx, false)
			} else {
				_, err = dec.AnteHandle(suite.ctx, tx, false, NextFn)
			}
			if tc.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.ErrorContains(err, "account disabled: invalid address")
			}
		})
	}
}
