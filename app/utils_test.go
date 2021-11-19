package app_test

import (
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/tests"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type AnteTestSuite struct {
	suite.Suite

	ctx          sdk.Context
	app          *fxcore.App
	clientCtx    client.Context
	anteHandler  sdk.AnteHandler
	ethSigner    ethtypes.Signer
	dynamicTxFee bool
}

func (suite *AnteTestSuite) SetupTest() {
	checkTx := false
	suite.app = fxcore.Setup(checkTx)

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 2, ChainID: "fxcore", Time: time.Now().UTC()})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(evmtypes.DefaultEVMDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1000000000000000000))
	suite.app.EvmKeeper.WithContext(suite.ctx)
	suite.app.EvmKeeper.WithChainID(suite.ctx)
	require.NoError(suite.T(), InitEvmModuleParams(suite.ctx, suite.app.EvmKeeper, suite.dynamicTxFee))
	infCtx := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	suite.app.AccountKeeper.SetParams(infCtx, authtypes.DefaultParams())
	suite.app.EvmKeeper.SetParams(infCtx, evmtypes.DefaultParams())

	encodingConfig := fxcore.MakeEncodingConfig()

	// We're using TestMsg amino encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)

	suite.anteHandler = app.NewAnteHandlerWithEVM(
		suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.EvmKeeper, suite.app.FeeMarketKeeper,
		ante.DefaultSigVerificationGasConsumer, encodingConfig.TxConfig.SignModeHandler(),
	)
	suite.ethSigner = ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(
	msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool,
) authsigning.Tx {
	return suite.CreateTestTxBuilder(msg, priv, accNum, signCosmosTx).GetTx()
}

// CreateTestTxBuilder is a helper function to create a tx builder given multiple inputs.
func (suite *AnteTestSuite) CreateTestTxBuilder(
	msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool,
) client.TxBuilder {
	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	suite.Require().NoError(err)

	txBuilder := suite.clientCtx.TxConfig.NewTxBuilder()
	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	suite.Require().True(ok)

	builder.SetExtensionOptions(option)

	err = msg.Sign(suite.ethSigner, tests.NewSigner(priv))
	suite.Require().NoError(err)

	err = builder.SetMsgs(msg)
	suite.Require().NoError(err)

	txData, err := evmtypes.UnpackTxData(msg.Data)
	suite.Require().NoError(err)

	fees := sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdk.NewIntFromBigInt(txData.Fee())))
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())

	if signCosmosTx {
		// First round: we gather all the signer infos. We use the "set empty
		// signature" hack to do that.
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
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
		sigV2, err = tx.SignWithPrivKey(
			suite.clientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData,
			txBuilder, priv, suite.clientCtx.TxConfig, txData.GetNonce(),
		)
		suite.Require().NoError(err)

		sigsV2 = []signing.SignatureV2{sigV2}

		err = txBuilder.SetSignatures(sigsV2...)
		suite.Require().NoError(err)
	}

	return txBuilder
}

var _ sdk.Tx = &invalidTx{}

type invalidTx struct{}

func (invalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (invalidTx) ValidateBasic() error { return nil }

func InitEvmModuleParams(ctx sdk.Context, keeper *evmkeeper.Keeper, dynamicTxFee bool) error {
	defaultEvmParams := evmtypes.DefaultParams()
	defaultFeeMarketParams := feemarkettypes.DefaultParams()

	if dynamicTxFee {
		defaultFeeMarketParams.EnableHeight = 1
		defaultFeeMarketParams.NoBaseFee = false
	}

	if err := keeper.HandleInitEvmParamsProposal(ctx, &evmtypes.InitEvmParamsProposal{
		Title:           "Init evm title",
		Description:     "Init emv module description",
		EvmParams:       &defaultEvmParams,
		FeemarketParams: &defaultFeeMarketParams,
	}); err != nil {
		return err
	}
	keeper.WithChainID(ctx)
	return nil
}
