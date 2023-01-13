package ante_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/server/config"

	"github.com/evmos/ethermint/app/ante"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/functionx/fx-core/v3/app/helpers"
)

func (suite *AnteTestSuite) TestEthSigVerificationDecorator() {
	getTx := func() sdk.Tx {
		signer := helpers.NewSigner(helpers.NewEthPrivKey())
		unprotectedTx := evmtypes.NewTxContract(nil, 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
		unprotectedTx.From = signer.Address().String()
		suite.Require().NoError(unprotectedTx.Sign(ethtypes.HomesteadSigner{}, signer))
		return unprotectedTx
	}
	testCases := []struct {
		name                string
		tx                  func() sdk.Tx
		allowUnprotectedTxs bool
		reCheckTx           bool
		expPass             bool
	}{
		{"ReCheckTx", func() sdk.Tx { return &invalidTx{} }, false, true, false},
		{"invalid transaction type", func() sdk.Tx { return &invalidTx{} }, false, false, false},
		{
			"invalid sender",
			func() sdk.Tx {
				addr := helpers.GenerateAddress()
				return evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 1, &addr, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
			},
			true,
			false,
			false,
		},
		{"successful signature verification", func() sdk.Tx {
			signer := helpers.NewSigner(helpers.NewEthPrivKey())

			signedTx := suite.NewTxContract()
			signedTx.From = signer.Address().String()
			ethSigner := ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
			suite.Require().NoError(signedTx.Sign(ethSigner, signer))
			return signedTx
		}, false, false, true},
		{"invalid, reject unprotected txs", getTx, false, false, false},
		{"successful, allow unprotected txs", getTx, true, false, true},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			params := suite.app.EvmKeeper.GetParams(suite.ctx)
			params.AllowUnprotectedTxs = tc.allowUnprotectedTxs
			suite.app.EvmKeeper.SetParams(suite.ctx, params)

			dec := ante.NewEthSigVerificationDecorator(suite.app.EvmKeeper)
			_, err := dec.AnteHandle(suite.ctx.WithIsReCheckTx(tc.reCheckTx), tc.tx(), false, NextFn)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) NewTxContract() *evmtypes.MsgEthereumTx {
	return evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
}

func (suite *AnteTestSuite) TestNewEthAccountVerificationDecorator() {

	addr := helpers.GenerateAddress()
	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		tx       sdk.Tx
		malleate func()
		checkTx  bool
		expPass  bool
	}{
		{"not CheckTx", nil, func() {}, false, true},
		{"invalid transaction type", &invalidTx{}, func() {}, true, false},
		{
			"sender not set to msg",
			evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			func() {},
			true,
			false,
		},
		{
			"sender not EOA",
			tx,
			func() {
				// set not as an EOA
				vmdb.SetCode(addr, []byte("1"))
			},
			true,
			false,
		},
		{
			"not enough balance to cover tx cost",
			tx,
			func() {
				// reset back to EOA
				vmdb.SetCode(addr, nil)
			},
			true,
			false,
		},
		{
			"success new account",
			tx,
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			true,
			true,
		},
		{
			"success existing account",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			true,
			true,
		},
	}

	dec := ante.NewEthAccountVerificationDecorator(suite.app.AccountKeeper, suite.app.EvmKeeper)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vmdb = suite.StateDB()
			tc.malleate()
			suite.Require().NoError(vmdb.Commit())

			_, err := dec.AnteHandle(suite.ctx.WithIsCheckTx(tc.checkTx), tc.tx, false, NextFn)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestEthNonceVerificationDecorator() {

	addr := helpers.GenerateAddress()
	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()

	testCases := []struct {
		name      string
		tx        sdk.Tx
		malleate  func()
		reCheckTx bool
		expPass   bool
	}{
		{"ReCheckTx", &invalidTx{}, func() {}, true, false},
		{"invalid transaction type", &invalidTx{}, func() {}, false, false},
		{"sender account not found", tx, func() {}, false, false},
		{
			"sender nonce missmatch",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			false,
			false,
		},
		{
			"success",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr.Bytes())
				suite.Require().NoError(acc.SetSequence(1))
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			false,
			true,
		},
	}

	suite.SetupTest()
	dec := ante.NewEthIncrementSenderSequenceDecorator(suite.app.AccountKeeper)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()
			_, err := dec.AnteHandle(suite.ctx.WithIsReCheckTx(tc.reCheckTx), tc.tx, false, NextFn)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestEthGasConsumeDecorator() {

	addr := helpers.GenerateAddress()

	txGasLimit := uint64(1000)
	tx := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), txGasLimit, big.NewInt(1), nil, nil, nil, nil)
	tx.From = addr.Hex()

	tx2GasLimit := uint64(1000000)
	tx2 := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), tx2GasLimit, big.NewInt(1), nil, nil, nil, &ethtypes.AccessList{{Address: addr, StorageKeys: nil}})
	tx2.From = addr.Hex()

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		tx       sdk.Tx
		gasLimit uint64
		malleate func()
		expPass  bool
		expPanic bool
	}{
		{"invalid transaction type", &invalidTx{}, 0, func() {}, false, false},
		{
			"sender not found",
			evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 1, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			0,
			func() {},
			false, false,
		},
		{
			"gas limit too low",
			tx,
			0,
			func() {},
			false, false,
		},
		{
			"not enough balance for fees",
			tx2,
			0,
			func() {},
			false, false,
		},
		{
			"not enough tx gas",
			tx2,
			0,
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))
			},
			false, true,
		},
		{
			"not enough block gas",
			tx2,
			0,
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))

				suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1))
			},
			false, true,
		},
		{
			"success",
			tx2,
			tx2GasLimit, // it's capped
			func() {
				vmdb.AddBalance(addr, big.NewInt(1000000))

				suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(10000000000000000000))
			},
			true, false,
		},
	}

	dec := ante.NewEthGasConsumeDecorator(suite.app.EvmKeeper, config.DefaultMaxTxGasWanted)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vmdb = suite.StateDB()
			tc.malleate()
			suite.Require().NoError(vmdb.Commit())

			if tc.expPanic {
				suite.Require().Panics(func() {
					_, _ = dec.AnteHandle(suite.ctx.WithIsCheckTx(true).WithGasMeter(sdk.NewGasMeter(1)), tc.tx, false, NextFn)
				})
				return
			}

			ctx, err := dec.AnteHandle(suite.ctx.WithIsCheckTx(true).WithGasMeter(sdk.NewInfiniteGasMeter()), tc.tx, false, NextFn)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
			suite.Require().Equal(tc.gasLimit, ctx.GasMeter().Limit())
		})
	}
}

func (suite *AnteTestSuite) TestCanTransferDecorator() {

	tx := evmtypes.NewTxContract(
		suite.app.EvmKeeper.ChainID(),
		1,
		big.NewInt(10),
		1000,
		big.NewInt(150),
		big.NewInt(200),
		nil,
		nil,
		&ethtypes.AccessList{},
	)
	tx2 := evmtypes.NewTxContract(
		suite.app.EvmKeeper.ChainID(),
		1,
		big.NewInt(10),
		1000,
		big.NewInt(150),
		big.NewInt(200),
		nil,
		nil,
		&ethtypes.AccessList{},
	)

	signer := helpers.NewSigner(helpers.NewEthPrivKey())

	tx.From = signer.Address().Hex()
	ethSigner := ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
	suite.Require().NoError(tx.Sign(ethSigner, signer))

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		tx       sdk.Tx
		malleate func()
		expPass  bool
	}{
		{"invalid transaction type", &invalidTx{}, func() {}, false},
		{"AsMessage failed", tx2, func() {}, false},
		{
			"evm CanTransfer failed",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.Address().Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			false,
		},
		{
			"success",
			tx,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.Address().Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				vmdb.AddBalance(signer.Address(), big.NewInt(1000000))
			},
			true,
		},
	}

	dec := ante.NewCanTransferDecorator(suite.app.EvmKeeper)
	suite.app.FeeMarketKeeper.SetBaseFee(suite.ctx, big.NewInt(100))

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vmdb = suite.StateDB()
			tc.malleate()
			suite.Require().NoError(vmdb.Commit())

			_, err := dec.AnteHandle(suite.ctx.WithIsCheckTx(true), tc.tx, false, NextFn)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestEthIncrementSenderSequenceDecorator() {

	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	to := helpers.GenerateAddress()

	contract := evmtypes.NewTxContract(suite.app.EvmKeeper.ChainID(), 0, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	contract.From = signer.Address().Hex()
	ethSigner := ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
	suite.Require().NoError(contract.Sign(ethSigner, signer))

	tx := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 0, &to, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx.From = signer.Address().Hex()
	suite.Require().NoError(tx.Sign(ethSigner, signer))

	tx2 := evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 1, &to, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil)
	tx2.From = signer.Address().Hex()
	suite.Require().NoError(tx2.Sign(ethSigner, signer))

	testCases := []struct {
		name     string
		tx       sdk.Tx
		malleate func()
		expPass  bool
		expPanic bool
	}{
		{
			"invalid transaction type",
			&invalidTx{},
			func() {},
			false, false,
		},
		{
			"no signers",
			evmtypes.NewTx(suite.app.EvmKeeper.ChainID(), 1, &to, big.NewInt(10), 1000, big.NewInt(1), nil, nil, nil, nil),
			func() {},
			false, false,
		},
		{
			"account not set to store",
			tx,
			func() {},
			false, false,
		},
		{
			"success - create contract",
			contract,
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.Address().Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			true, false,
		},
		{
			"success - call",
			tx2,
			func() {},
			true, false,
		},
	}

	dec := ante.NewEthIncrementSenderSequenceDecorator(suite.app.AccountKeeper)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()

			if tc.expPanic {
				suite.Require().Panics(func() {
					_, _ = dec.AnteHandle(suite.ctx, tc.tx, false, NextFn)
				})
				return
			}

			_, err := dec.AnteHandle(suite.ctx, tc.tx, false, NextFn)

			if tc.expPass {
				suite.Require().NoError(err)
				msg := tc.tx.(*evmtypes.MsgEthereumTx)

				txData, err := evmtypes.UnpackTxData(msg.Data)
				suite.Require().NoError(err)

				nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, signer.Address())
				suite.Require().Equal(txData.GetNonce()+1, nonce)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestEthSetupContextDecorator() {
	testCases := []struct {
		name    string
		tx      func() sdk.Tx
		expPass bool
	}{
		{"invalid transaction type - does not implement GasTx", func() sdk.Tx { return &invalidTx{} }, false},
		{
			"success - transaction implement GasTx",
			func() sdk.Tx {
				return evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					1000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
			},
			true,
		},
	}

	dec := ante.NewEthSetUpContextDecorator(suite.app.EvmKeeper)

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := dec.AnteHandle(suite.ctx, tc.tx(), false, NextFn)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
