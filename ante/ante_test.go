package ante_test

import (
	"errors"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	fxante "github.com/functionx/fx-core/v7/ante"
	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

type AnteTestSuite struct {
	suite.Suite

	app         *app.App
	anteHandler sdk.AnteHandler
	ctx         sdk.Context
	signer      *helpers.Signer
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, &AnteTestSuite{})
}

const TestGasLimit uint64 = 100000

func (suite *AnteTestSuite) StateDB() *statedb.StateDB {
	return statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))
}

func (suite *AnteTestSuite) SetupTest() {
	valConsPriv := ed25519.GenPrivKey()

	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(1, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		Height:          suite.app.LastBlockHeight(),
		ChainID:         fxtypes.ChainId(),
		ProposerAddress: valConsPriv.PubKey().Address(),
		Time:            time.Now().UTC(),
	})
	suite.ctx = suite.ctx.
		WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt()))).
		WithBlockGasMeter(sdk.NewGasMeter(1e18)).
		WithGasMeter(sdk.NewInfiniteGasMeter())
	suite.signer = helpers.NewSigner(helpers.NewPriKey())

	valAddr := helpers.GenHexAddress().Bytes()
	validator, err := stakingtypes.NewValidator(valAddr, valConsPriv.PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)

	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	suite.Require().NoError(err)
	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	encodingConfig := app.MakeEncodingConfig()
	// We're using TestMsg amino encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	options := fxante.HandlerOptions{
		AccountKeeper:   suite.app.AccountKeeper,
		BankKeeper:      suite.app.BankKeeper,
		EvmKeeper:       suite.app.EvmKeeper,
		FeeMarketKeeper: suite.app.FeeMarketKeeper,
		IbcKeeper:       suite.app.IBCKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  fxante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  0,
	}
	suite.Require().NoError(options.Validate())
	suite.anteHandler = fxante.NewAnteHandler(options)
}

func (suite *AnteTestSuite) TestAnteHandler() {
	testCases := []struct {
		name      string
		txFn      func(signer helpers.Signer) sdk.Tx
		checkTx   bool
		reCheckTx bool
		expPass   bool
	}{
		{
			"success - DeliverTx (contract)",
			func(signer helpers.Signer) sdk.Tx {
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
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			false, false, true,
		},
		{
			"success - CheckTx (contract)",
			func(signer helpers.Signer) sdk.Tx {
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
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			true, false, true,
		},
		{
			"success - ReCheckTx (contract)",
			func(signer helpers.Signer) sdk.Tx {
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
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"success - DeliverTx",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			false, false, true,
		},
		{
			"success - CheckTx",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			true, false, true,
		},
		{
			"success - ReCheckTx",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"success - CheckTx (cosmos tx not signed)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(150),
					big.NewInt(200),
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			false, true, true,
		},
		{
			"fail - CheckTx (cosmos tx is not valid)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)
				// bigger than MaxGasWanted
				txBuilder.SetGasLimit(uint64(1 << 63))
				return txBuilder.GetTx()
			},
			true, false, false,
		},
		{
			"fail - CheckTx (memo too long)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)
				txBuilder.SetMemo(strings.Repeat("*", 257))
				return txBuilder.GetTx()
			},
			true, false, false,
		},
		{
			"fail - CheckTx (ExtensionOptionsEthereumTx not set)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false, true)
				return txBuilder.GetTx()
			},
			true, false, false,
		},
		// Based on EVMBackend.SendTransaction, for cosmos tx, forcing null for some fields except ExtensionOptions, Fee, MsgEthereumTx
		// should be part of consensus
		{
			"fail - DeliverTx (cosmos tx signed)",
			func(signer helpers.Signer) sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, signer.AccAddress())
				suite.Require().NoError(err)
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					nonce,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, true)
				return tx
			},
			false, false, false,
		},
		{
			"fail - DeliverTx (cosmos tx with memo)",
			func(signer helpers.Signer) sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, signer.AccAddress())
				suite.Require().NoError(err)
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					nonce,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)
				txBuilder.SetMemo("memo for cosmos tx not allowed")
				return txBuilder.GetTx()
			},
			false, false, false,
		},
		{
			"fail - DeliverTx (cosmos tx with timeoutheight)",
			func(signer helpers.Signer) sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, signer.AccAddress())
				suite.Require().NoError(err)
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					nonce,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)
				txBuilder.SetTimeoutHeight(10)
				return txBuilder.GetTx()
			},
			false, false, false,
		},
		{
			"fail - DeliverTx (invalid fee amount)",
			func(signer helpers.Signer) sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, signer.AccAddress())
				suite.Require().NoError(err)
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					nonce,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)

				txData, err := evmtypes.UnpackTxData(signedTx.Data)
				suite.Require().NoError(err)

				invalidFee := new(big.Int).Add(txData.Fee(), big.NewInt(1))
				invalidFeeAmount := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(invalidFee))
				txBuilder.SetFeeAmount(sdk.NewCoins(invalidFeeAmount))
				return txBuilder.GetTx()
			},
			false, false, false,
		},
		{
			"fail - DeliverTx (invalid fee gaslimit)",
			func(signer helpers.Signer) sdk.Tx {
				nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, signer.AccAddress())
				suite.Require().NoError(err)

				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					nonce,
					&to,
					big.NewInt(10),
					100000,
					big.NewInt(1),
					nil,
					nil,
					nil,
					nil,
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)

				expGasLimit := signedTx.GetGas()
				invalidGasLimit := expGasLimit + 1
				txBuilder.SetGasLimit(invalidGasLimit)
				return txBuilder.GetTx()
			},
			false, false, false,
		},
		{
			"fails - invalid from",
			func(signer helpers.Signer) sdk.Tx {
				msg := evmtypes.NewTxContract(
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
				msg.From = signer.Address().String()
				tx := suite.CreateTestTx(msg, signer.PrivKey(), 1, false)
				msg = tx.GetMsgs()[0].(*evmtypes.MsgEthereumTx)
				msg.From = signer.Address().String()
				return tx
			},
			true, false, false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
			params.MinGasPrice = sdk.ZeroDec()
			err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
			suite.NoError(err)
			signer := helpers.NewSigner(helpers.NewEthPrivKey())

			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.Address().Bytes())
			suite.Require().NoError(acc.SetSequence(1))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			err = suite.app.EvmKeeper.SetBalance(suite.ctx, signer.Address(), big.NewInt(10000000000))
			suite.NoError(err)

			suite.app.FeeMarketKeeper.SetBaseFee(suite.ctx, big.NewInt(100))

			suite.ctx = suite.ctx.WithIsCheckTx(tc.checkTx).WithIsReCheckTx(tc.reCheckTx)

			_, err = suite.anteHandler(suite.ctx, tc.txFn(*signer), false)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestAnteHandlerWithDynamicTxFee() {
	testCases := []struct {
		name           string
		txFn           func(signer helpers.Signer) sdk.Tx
		enableLondonHF bool
		checkTx        bool
		reCheckTx      bool
		expPass        bool
	}{
		{
			"success - DeliverTx (contract)",
			func(signer helpers.Signer) sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			false, false, true,
		},
		{
			"success - CheckTx (contract)",
			func(signer helpers.Signer) sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			true, false, true,
		},
		{
			"success - ReCheckTx (contract)",
			func(signer helpers.Signer) sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			false, true, true,
		},
		{
			"success - DeliverTx",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			false, false, true,
		},
		{
			"success - CheckTx",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			true, false, true,
		},
		{
			"success - ReCheckTx",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			false, true, true,
		},
		{
			"success - CheckTx (cosmos tx not signed)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			true,
			false, true, true,
		},
		{
			"fail - CheckTx (cosmos tx is not valid)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)
				// bigger than MaxGasWanted
				txBuilder.SetGasLimit(uint64(1 << 63))
				return txBuilder.GetTx()
			},
			true,
			true, false, false,
		},
		{
			"fail - CheckTx (memo too long)",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				txBuilder := suite.CreateTestTxBuilder(signedTx, signer.PrivKey(), 1, false)
				txBuilder.SetMemo(strings.Repeat("*", 257))
				return txBuilder.GetTx()
			},
			true,
			true, false, false,
		},
		{
			"fail - DynamicFeeTx without london hark fork",
			func(signer helpers.Signer) sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			false,
			false, false, false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			suite.ctx = suite.ctx.WithBlockHeight(1)
			params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
			params.MinGasPrice = sdk.ZeroDec()
			params.BaseFee = sdkmath.NewIntFromUint64(ethparams.InitialBaseFee)
			params.EnableHeight = 1
			params.NoBaseFee = false
			err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
			suite.NoError(err)
			if !tc.enableLondonHF {
				ethParams := suite.app.EvmKeeper.GetParams(suite.ctx)
				maxInt := sdkmath.NewInt(math.MaxInt64)
				ethParams.ChainConfig.LondonBlock = &maxInt
				ethParams.ChainConfig.ArrowGlacierBlock = &maxInt
				ethParams.ChainConfig.GrayGlacierBlock = &maxInt
				ethParams.ChainConfig.MergeNetsplitBlock = &maxInt
				ethParams.ChainConfig.ShanghaiBlock = &maxInt
				ethParams.ChainConfig.CancunBlock = &maxInt
				err = suite.app.EvmKeeper.SetParams(suite.ctx, ethParams)
				suite.NoError(err)
			}

			signer := helpers.NewSigner(helpers.NewEthPrivKey())
			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.Address().Bytes())
			suite.Require().NoError(acc.SetSequence(1))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			suite.ctx = suite.ctx.WithIsCheckTx(tc.checkTx).WithIsReCheckTx(tc.reCheckTx)
			err = suite.app.EvmKeeper.SetBalance(suite.ctx, signer.Address(), big.NewInt((ethparams.InitialBaseFee+10)*100000))
			suite.NoError(err)
			_, err = suite.anteHandler(suite.ctx, tc.txFn(*signer), false)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *AnteTestSuite) TestAnteHandlerWithParams() {
	testCases := []struct {
		name         string
		txFn         func(signer helpers.Signer) sdk.Tx
		enableCall   bool
		enableCreate bool
		expErr       error
	}{
		{
			"fail - Contract Creation Disabled",
			func(signer helpers.Signer) sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			true, false,
			evmtypes.ErrCreateDisabled,
		},
		{
			"success - Contract Creation Enabled",
			func(signer helpers.Signer) sdk.Tx {
				signedContractTx := evmtypes.NewTxContract(
					suite.app.EvmKeeper.ChainID(),
					1,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedContractTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedContractTx, signer.PrivKey(), 1, false)
				return tx
			},
			true, true,
			nil,
		},
		{
			"fail - EVM Call Disabled",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			false, true,
			evmtypes.ErrCallDisabled,
		},
		{
			"success - EVM Call Enabled",
			func(signer helpers.Signer) sdk.Tx {
				to := helpers.GenHexAddress()
				signedTx := evmtypes.NewTx(
					suite.app.EvmKeeper.ChainID(),
					1,
					&to,
					big.NewInt(10),
					100000,
					nil,
					big.NewInt(ethparams.InitialBaseFee+1),
					big.NewInt(1),
					nil,
					&ethtypes.AccessList{},
				)
				signedTx.From = signer.Address().String()

				tx := suite.CreateTestTx(signedTx, signer.PrivKey(), 1, false)
				return tx
			},
			true, true,
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
			params.MinGasPrice = sdk.ZeroDec()
			params.BaseFee = sdkmath.NewIntFromUint64(ethparams.InitialBaseFee)
			err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
			suite.NoError(err)
			ethParams := suite.app.EvmKeeper.GetParams(suite.ctx)
			ethParams.EnableCall = tc.enableCall
			ethParams.EnableCreate = tc.enableCreate
			err = suite.app.EvmKeeper.SetParams(suite.ctx, ethParams)
			suite.NoError(err)
			signer := helpers.NewSigner(helpers.NewEthPrivKey())
			acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, signer.Address().Bytes())
			suite.Require().NoError(acc.SetSequence(1))
			suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

			suite.ctx = suite.ctx.WithIsCheckTx(true)
			err = suite.app.EvmKeeper.SetBalance(suite.ctx, signer.Address(), big.NewInt((ethparams.InitialBaseFee+10)*100000))
			suite.NoError(err)
			_, err = suite.anteHandler(suite.ctx, tc.txFn(*signer), false)
			if tc.expErr == nil {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().True(errors.Is(err, tc.expErr))
			}
		})
	}
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool, unsetExtensionOptions ...bool) authsigning.Tx {
	return suite.CreateTestTxBuilder(msg, priv, accNum, signCosmosTx, unsetExtensionOptions...).GetTx()
}

// CreateTestTxBuilder is a helper function to create a tx builder given multiple inputs.
func (suite *AnteTestSuite) CreateTestTxBuilder(msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool, unsetExtensionOptions ...bool) client.TxBuilder {
	var option *codectypes.Any
	var err error
	if len(unsetExtensionOptions) == 0 {
		option, err = codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
		suite.Require().NoError(err)
	}
	cliCtx := NewClientCtx()
	txBuilder := cliCtx.TxConfig.NewTxBuilder()
	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	suite.Require().True(ok)

	if len(unsetExtensionOptions) == 0 {
		builder.SetExtensionOptions(option)
	}

	ethSigner := ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
	err = msg.Sign(ethSigner, helpers.NewSigner(priv))
	suite.Require().NoError(err)

	msg.From = ""
	err = builder.SetMsgs(msg)
	suite.Require().NoError(err)

	txData, err := evmtypes.UnpackTxData(msg.Data)
	suite.Require().NoError(err)

	fees := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(txData.Fee())))
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())

	if signCosmosTx {
		// First round: we gather all the signer infos. We use the "set empty
		// signature" hack to do that.
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  cliCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: txData.GetNonce(),
		}

		sigsV2 := []signing.SignatureV2{sigV2}

		err = txBuilder.SetSignatures(sigsV2...)
		suite.Require().NoError(err)

		// Second round: all signer infos are set, so each signer can sign.

		signerData := authsigning.SignerData{
			ChainID:       suite.ctx.ChainID(),
			AccountNumber: accNum,
			Sequence:      txData.GetNonce(),
		}
		sigV2, err = clitx.SignWithPrivKey(
			cliCtx.TxConfig.SignModeHandler().DefaultMode(), signerData,
			txBuilder, priv, cliCtx.TxConfig, txData.GetNonce(),
		)
		suite.Require().NoError(err)

		sigsV2 = []signing.SignatureV2{sigV2}

		err = txBuilder.SetSignatures(sigsV2...)
		suite.Require().NoError(err)
	}

	return txBuilder
}

func (suite *AnteTestSuite) CreateEmptyTestTx(txBuilder client.TxBuilder, privs []cryptotypes.PrivKey, accNums []uint64, accSeqs []uint64) (authsigning.Tx, error) {
	cliCtx := NewClientCtx()
	signMode := cliCtx.TxConfig.SignModeHandler().DefaultMode()
	sigsV2 := make([]signing.SignatureV2, 0, len(privs))
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signMode,
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	if err := txBuilder.SetSignatures(sigsV2...); err != nil {
		return nil, err
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := authsigning.SignerData{
			ChainID:       suite.ctx.ChainID(),
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := clitx.SignWithPrivKey(
			signMode,
			signerData,
			txBuilder,
			priv,
			cliCtx.TxConfig,
			accSeqs[i],
		)
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	if err := txBuilder.SetSignatures(sigsV2...); err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}

func (suite *AnteTestSuite) CreateTestCosmosTxBuilder(gasPrice sdkmath.Int, denom string, msgs ...sdk.Msg) client.TxBuilder {
	txBuilder := NewClientCtx().TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(TestGasLimit)
	fees := &sdk.Coins{{Denom: denom, Amount: gasPrice.MulRaw(int64(TestGasLimit))}}
	txBuilder.SetFeeAmount(*fees)
	err := txBuilder.SetMsgs(msgs...)
	suite.Require().NoError(err)
	return txBuilder
}

func (suite *AnteTestSuite) BuildTestEthTx(
	from common.Address,
	to common.Address,
	amount *big.Int,
	input []byte,
	gasPrice *big.Int,
	gasFeeCap *big.Int,
	gasTipCap *big.Int,
	accesses *ethtypes.AccessList,
) *evmtypes.MsgEthereumTx {
	chainID := suite.app.EvmKeeper.ChainID()
	nonce := suite.app.EvmKeeper.GetNonce(
		suite.ctx,
		common.BytesToAddress(from.Bytes()),
	)

	msgEthereumTx := evmtypes.NewTx(
		chainID,
		nonce,
		&to,
		amount,
		TestGasLimit,
		gasPrice,
		gasFeeCap,
		gasTipCap,
		input,
		accesses,
	)
	msgEthereumTx.From = from.String()
	return msgEthereumTx
}

var _ sdk.Tx = &invalidTx{}

type invalidTx struct{}

func (invalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (invalidTx) ValidateBasic() error { return nil }

func NextFn(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}

func NewClientCtx() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	return client.Context{}.WithTxConfig(encodingConfig.TxConfig)
}
