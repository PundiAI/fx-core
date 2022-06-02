package ante_test

import (
	"errors"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"math/big"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/tests"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

func (suite AnteTestSuite) TestAnteHandler() {
	suite.SetupTest() // reset

	addr, privKey := tests.NewAddrKey()
	to := tests.GenerateAddress()

	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
	suite.Require().NoError(acc.SetSequence(1))
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	suite.Require().NoError(suite.app.EvmKeeper.SetBalance(suite.ctx, addr, big.NewInt(10000000000)))

	suite.app.FeeMarketKeeper.SetBaseFee(suite.ctx, big.NewInt(100))

	testCases := []struct {
		name      string
		txFn      func() sdk.Tx
		checkTx   bool
		reCheckTx bool
		expPass   bool
	}{
		{
			"success - DeliverTx (contract)",
			func() sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			false, false, true,
		},
		{
			"success - CheckTx (contract)",
			func() sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					2,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			true, false, true,
		},
		{
			"success - ReCheckTx (contract)",
			func() sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					3,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"success - DeliverTx",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					4,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			false, false, true,
		},
		{
			"success - CheckTx",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					5,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			true, false, true,
		},
		{
			"success - ReCheckTx",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					6,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			}, false, true, true,
		},
		{
			"success - CheckTx (cosmos tx not signed)",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					7,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			}, false, true, true,
		},
		{
			"fail - CheckTx (cosmos tx is not valid)",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 8, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)
				// bigger than MaxGasWanted
				txBuilder.SetGasLimit(uint64(1 << 63))
				return txBuilder.GetTx()
			}, true, false, false,
		},
		{
			"fail - CheckTx (memo too long)",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 5, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)
				txBuilder.SetMemo(strings.Repeat("*", 257))
				return txBuilder.GetTx()
			}, true, false, false,
		},
		{
			"fail - CheckTx (ExtensionOptionsEthereumTx not set)",
			func() sdk.Tx {
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 5, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false, true)
				return txBuilder.GetTx()
			}, true, false, false,
		},
		// Based on EVMBackend.SendTransaction, for cosmos tx, forcing null for some fields except ExtensionOptions, Fee, MsgEthereumTx
		// should be part of consensus
		{
			"fail - DeliverTx (cosmos tx signed)",
			func() sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
				suite.Require().NoError(err)
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), nonce, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, true)
				return tx
			}, false, false, false,
		},
		{
			"fail - DeliverTx (cosmos tx with memo)",
			func() sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
				suite.Require().NoError(err)
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), nonce, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)
				txBuilder.SetMemo("memo for cosmos tx not allowed")
				return txBuilder.GetTx()
			}, false, false, false,
		},
		{
			"fail - DeliverTx (cosmos tx with timeoutheight)",
			func() sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
				suite.Require().NoError(err)
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), nonce, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)
				txBuilder.SetTimeoutHeight(10)
				return txBuilder.GetTx()
			}, false, false, false,
		},
		{
			"fail - DeliverTx (invalid fee amount)",
			func() sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
				suite.Require().NoError(err)
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), nonce, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)

				txData, err := evmtypes.UnpackTxData(signedTx.Data)
				suite.Require().NoError(err)

				expFee := txData.Fee()
				invalidFee := new(big.Int).Add(expFee, big.NewInt(1))
				invalidFeeAmount := sdk.Coins{sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(invalidFee))}
				txBuilder.SetFeeAmount(invalidFeeAmount)
				return txBuilder.GetTx()
			}, false, false, false,
		},
		{
			"fail - DeliverTx (invalid fee gaslimit)",
			func() sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, acc.GetAddress())
				suite.Require().NoError(err)
				signedTx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), nonce, &to, big.NewInt(10), 100000, big.NewInt(1), nil, nil, nil, nil)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)

				expGasLimit := signedTx.GetGas()
				invalidGasLimit := expGasLimit + 1
				txBuilder.SetGasLimit(invalidGasLimit)
				return txBuilder.GetTx()
			}, false, false, false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.ctx = suite.ctx.WithIsCheckTx(tc.checkTx).WithIsReCheckTx(tc.reCheckTx)

			_, err := suite.anteHandler(suite.ctx, tc.txFn(), false)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite AnteTestSuite) TestAnteHandlerWithDynamicTxFee() {
	addr, privKey := tests.NewAddrKey()
	to := tests.GenerateAddress()

	testCases := []struct {
		name      string
		txFn      func() sdk.Tx
		checkTx   bool
		reCheckTx bool
		expPass   bool
	}{
		{
			"success - DeliverTx (contract)",
			func() sdk.Tx {
				signedContractTx :=
					evmtypes.NewTxContract(
						suite.app.EvmKeeper.ChainID(),
						1,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			false, false, true,
		},
		{
			"success - CheckTx (contract)",
			func() sdk.Tx {
				signedContractTx :=
					evmtypes.NewTxContract(
						suite.app.EvmKeeper.ChainID(),
						1,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			true, false, true,
		},
		{
			"success - ReCheckTx (contract)",
			func() sdk.Tx {
				signedContractTx :=
					evmtypes.NewTxContract(
						suite.app.EvmKeeper.ChainID(),
						1,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"success - DeliverTx",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			false, false, true,
		},
		{
			"success - CheckTx",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			true, false, true,
		},
		{
			"success - ReCheckTx",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"success - CheckTx (cosmos tx not signed)",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"fail - CheckTx (cosmos tx is not valid)",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)
				// bigger than MaxGasWanted
				txBuilder.SetGasLimit(uint64(1 << 63))
				return txBuilder.GetTx()
			},
			true, false, false,
		},
		{
			"fail - CheckTx (memo too long)",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						big.NewInt(feemarkettypes.MinBaseFee.Int64()+1),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				txBuilder := suite.CreateTestTxBuilder(signedTx, privKey, 1, false)
				txBuilder.SetMemo(strings.Repeat("*", 257))
				return txBuilder.GetTx()
			},
			true, false, false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
			suite.Require().NoError(acc.SetSequence(1))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			suite.ctx = suite.ctx.WithIsCheckTx(tc.checkTx).WithIsReCheckTx(tc.reCheckTx)
			suite.Require().NoError(suite.app.EvmKeeper.SetBalance(suite.ctx, addr, big.NewInt((feemarkettypes.MinBaseFee.Int64()+10)*100000)))
			_, err := suite.anteHandler(suite.ctx, tc.txFn(), false)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite AnteTestSuite) TestAnteHandlerWithParams() {
	addr, privKey := tests.NewAddrKey()
	to := tests.GenerateAddress()

	testCases := []struct {
		name         string
		txFn         func() sdk.Tx
		enableCall   bool
		enableCreate bool
		expErr       error
	}{
		{
			"fail - Contract Creation Disabled",
			func() sdk.Tx {
				signedContractTx :=
					evmtypes.NewTxContract(
						suite.app.EvmKeeper.ChainID(),
						1,
						big.NewInt(10),
						100000,
						nil,
						feemarkettypes.MinBaseFee.BigInt(),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			true, false,
			evmtypes.ErrCreateDisabled,
		},
		{
			"success - Contract Creation Enabled",
			func() sdk.Tx {
				signedContractTx :=
					evmtypes.NewTxContract(
						suite.app.EvmKeeper.ChainID(),
						1,
						big.NewInt(10),
						100000,
						nil,
						feemarkettypes.MinBaseFee.BigInt(),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedContractTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedContractTx, privKey, 1, false)
				return tx
			},
			true, true,
			nil,
		},
		{
			"fail - EVM Call Disabled",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						feemarkettypes.MinBaseFee.BigInt(),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			false, true,
			evmtypes.ErrCallDisabled,
		},
		{
			"success - EVM Call Enabled",
			func() sdk.Tx {
				signedTx :=
					evmtypes.NewTx(
						suite.app.EvmKeeper.ChainID(),
						1,
						&to,
						big.NewInt(10),
						100000,
						nil,
						feemarkettypes.MinBaseFee.BigInt(),
						big.NewInt(1),
						nil,
						&types.AccessList{},
					)
				signedTx.From = addr.Hex()

				tx := suite.CreateTestTx(signedTx, privKey, 1, false)
				return tx
			},
			true, true,
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			params := suite.app.EvmKeeper.GetParams(suite.ctx)
			params.EnableCreate = tc.enableCreate
			params.EnableCall = tc.enableCall
			suite.app.EvmKeeper.SetParams(suite.ctx, params)

			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
			suite.Require().NoError(acc.SetSequence(1))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			suite.ctx = suite.ctx.WithIsCheckTx(true)
			err := suite.app.EvmKeeper.SetBalance(suite.ctx, addr, feemarkettypes.MinBaseFee.Mul(sdk.NewInt(1000000)).BigInt())
			suite.Require().NoError(err)
			_, err = suite.anteHandler(suite.ctx, tc.txFn(), false)
			if tc.expErr == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().True(errors.Is(err, tc.expErr))
			}
		})
	}
}

func (suite AnteTestSuite) TestAnteHandlerWithEthSecp256k1() {
	var (
		secp256k1Key       = tests.NewPriKey()
		_, ethSecp256k1Key = tests.NewAddrKey()
	)

	testCases := []struct {
		name    string
		txFn    func() sdk.Tx
		expFlag bool
	}{
		{
			"success - evm tx with secp256k1",
			func() sdk.Tx {

				msg := testdata.NewTestMsg(secp256k1Key.PubKey().Address().Bytes())
				suite.Require().NoError(suite.txBuilder.SetMsgs(msg))

				account := suite.app.AccountKeeper.GetAccount(suite.ctx, secp256k1Key.PubKey().Address().Bytes())

				privs, accNums, accSeqs := []cryptotypes.PrivKey{secp256k1Key}, []uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()}
				tx, err := suite.CreateEmptyTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
				suite.Require().NoError(err)
				return tx
			},
			true,
		},
		{
			"success - evm tx with secp256k1",
			func() sdk.Tx {
				msg := testdata.NewTestMsg(secp256k1Key.PubKey().Address().Bytes())
				suite.Require().NoError(suite.txBuilder.SetMsgs(msg))

				account := suite.app.AccountKeeper.GetAccount(suite.ctx, secp256k1Key.PubKey().Address().Bytes())

				privs, accNums, accSeqs := []cryptotypes.PrivKey{secp256k1Key}, []uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()}
				tx, err := suite.CreateEmptyTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
				suite.Require().NoError(err)
				return tx
			},
			true,
		},
		{
			"success - evm tx with eth_secp256k1",
			func() sdk.Tx {
				msg := testdata.NewTestMsg(ethSecp256k1Key.PubKey().Address().Bytes())
				suite.Require().NoError(suite.txBuilder.SetMsgs(msg))

				account := suite.app.AccountKeeper.GetAccount(suite.ctx, ethSecp256k1Key.PubKey().Address().Bytes())

				privs, accNums, accSeqs := []cryptotypes.PrivKey{ethSecp256k1Key}, []uint64{account.GetAccountNumber()}, []uint64{account.GetSequence()}
				tx, err := suite.CreateEmptyTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
				suite.Require().NoError(err)
				return tx
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, secp256k1Key.PubKey().Address().Bytes())
			suite.Require().NoError(acc.SetSequence(0))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			suite.ctx = suite.ctx.WithIsCheckTx(true)
			amount := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1e18).Mul(sdk.NewInt(1000))))
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, amount))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, secp256k1Key.PubKey().Address().Bytes(), amount))

			acc = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, ethSecp256k1Key.PubKey().Address().Bytes())
			suite.Require().NoError(acc.SetSequence(0))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			suite.ctx = suite.ctx.WithIsCheckTx(true)
			suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, amount))
			suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, ethSecp256k1Key.PubKey().Address().Bytes(), amount))

			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			suite.txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin(fxtypes.DefaultDenom, 400)))
			suite.txBuilder.SetGasLimit(testdata.NewTestGasLimit())

			// Set high gas price so standard test fee fails
			feeAmt := sdk.NewDecCoinFromDec(fxtypes.DefaultDenom, sdk.NewDec(200).Quo(sdk.NewDec(100000)))
			minGasPrice := []sdk.DecCoin{feeAmt}
			suite.ctx = suite.ctx.WithMinGasPrices(minGasPrice).WithIsCheckTx(true)

			_, err := suite.anteHandler(suite.ctx, tc.txFn(), false)
			if tc.expFlag {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
