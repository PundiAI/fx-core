package ante_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	fxante "github.com/functionx/fx-core/v3/ante"
	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

type AnteTestSuite struct {
	suite.Suite

	chainId     string
	app         *app.App
	anteHandler sdk.AnteHandler

	privateKey  cryptotypes.PrivKey
	consAddress sdk.ConsAddress
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, &AnteTestSuite{chainId: fxtypes.Name})
}

//func (suite *AnteTestSuite) StateDB() *statedb.StateDB {
//	return statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))
//}

func (suite *AnteTestSuite) GetContext(height int64) sdk.Context {
	ctx := suite.app.BaseApp.NewContext(false, tmproto.Header{Height: height, ChainID: suite.chainId, ProposerAddress: suite.consAddress, Time: time.Now().UTC()})
	context, _ := ctx.CacheContext()
	return context
}

func (suite *AnteTestSuite) SetupTest() {

	suite.app = helpers.Setup(false, false)

	// account key
	suite.privateKey = secp256k1.GenPrivKey()

	// consensus key
	valConsPriv := ed25519.GenPrivKey()
	suite.consAddress = sdk.ConsAddress(valConsPriv.PubKey().Address())

	ctx := suite.GetContext(1)
	ctx = ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdk.OneInt())))
	ctx = ctx.WithBlockGasMeter(sdk.NewGasMeter(1e18))
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	suite.app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	valAddr := sdk.ValAddress(suite.GetAccAddress().Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, valConsPriv.PubKey(), stakingtypes.Description{})
	suite.Require().NoError(err)

	err = suite.app.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	suite.Require().NoError(err)
	suite.app.StakingKeeper.SetValidator(ctx, validator)

	encodingConfig := app.MakeEncodingConfig()
	// We're using TestMsg amino encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	options := fxante.HandlerOptions{
		AccountKeeper:   suite.app.AccountKeeper,
		BankKeeper:      suite.app.BankKeeper,
		EvmKeeper:       suite.app.EvmKeeper,
		FeeMarketKeeper: suite.app.FeeMarketKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  fxante.DefaultSigVerificationGasConsumer,
	}
	suite.Require().NoError(options.Validate())
	suite.anteHandler = fxante.NewAnteHandler(options)
}

func (suite *AnteTestSuite) GetAccAddress() sdk.AccAddress {
	return suite.privateKey.PubKey().Address().Bytes()
}

func NewClientCtx() client.Context {
	encodingConfig := app.MakeEncodingConfig()
	return client.Context{}.WithTxConfig(encodingConfig.TxConfig)
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(cliCtx client.Context, msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool, unsetExtensionOptions ...bool) authsigning.Tx {
	return suite.CreateTestTxBuilder(cliCtx, msg, priv, accNum, signCosmosTx, unsetExtensionOptions...).GetTx()
}

// CreateTestTxBuilder is a helper function to create a tx builder given multiple inputs.
func (suite *AnteTestSuite) CreateTestTxBuilder(cliCtx client.Context, msg *evmtypes.MsgEthereumTx, priv cryptotypes.PrivKey, accNum uint64, signCosmosTx bool, unsetExtensionOptions ...bool) client.TxBuilder {
	var option *codectypes.Any
	var err error
	if len(unsetExtensionOptions) == 0 {
		option, err = codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
		suite.Require().NoError(err)
	}

	txBuilder := cliCtx.TxConfig.NewTxBuilder()
	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	suite.Require().True(ok)

	if len(unsetExtensionOptions) == 0 {
		builder.SetExtensionOptions(option)
	}

	ethSigner := ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())
	err = msg.Sign(ethSigner, helpers.NewSigner(priv))
	suite.Require().NoError(err)

	err = builder.SetMsgs(msg)
	suite.Require().NoError(err)

	txData, err := evmtypes.UnpackTxData(msg.Data)
	suite.Require().NoError(err)

	fees := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(txData.Fee())))
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
			ChainID:       suite.chainId,
			AccountNumber: accNum,
			Sequence:      txData.GetNonce(),
		}
		sigV2, err = tx.SignWithPrivKey(
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
	var sigsV2 []signing.SignatureV2
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
			ChainID:       suite.chainId,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
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

var _ sdk.Tx = &invalidTx{}

type invalidTx struct{}

func (invalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (invalidTx) ValidateBasic() error { return nil }
