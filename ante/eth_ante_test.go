package ante_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethante "github.com/evmos/ethermint/app/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/ante"
	"github.com/pundiai/fx-core/v8/app"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func TestEthPubkeyParse(t *testing.T) {
	chainConfig := evmtypes.DefaultChainConfig().EthereumConfig(big.NewInt(530))
	ethSigner := ethtypes.MakeSigner(chainConfig, big.NewInt(100))
	privateKey := helpers.NewEthPrivKey()
	addr := privateKey.PubKey().Address()

	tx := buildLegacyTx(t, privateKey, ethSigner)
	pub, err := ante.EthPubkeyParse(tx, ethSigner, addr.Bytes())
	require.NoError(t, err)
	require.Equal(t, addr, pub.Address())

	tx = buildAccessListTx(t, privateKey, ethSigner)
	pub, err = ante.EthPubkeyParse(tx, ethSigner, addr.Bytes())
	require.NoError(t, err)
	require.Equal(t, addr, pub.Address())

	tx = buildDynamicFeeTx(t, privateKey, ethSigner)
	pub, err = ante.EthPubkeyParse(tx, ethSigner, addr.Bytes())
	require.NoError(t, err)
	require.Equal(t, addr, pub.Address())
}

func TestCheckAndSetEthSenderNonce(t *testing.T) {
	chainConfig := evmtypes.DefaultChainConfig().EthereumConfig(big.NewInt(530))
	ethSigner := ethtypes.MakeSigner(chainConfig, big.NewInt(100))
	privateKey := helpers.NewEthPrivKey()
	addr := privateKey.PubKey().Address()

	myApp, ctx := helpers.NewAppWithValNumber(t, 1)
	acc := myApp.AccountKeeper.GetAccount(ctx, addr.Bytes())
	require.Empty(t, acc)
	acc = myApp.AccountKeeper.NewAccountWithAddress(ctx, addr.Bytes())
	require.NoError(t, acc.SetSequence(1))
	myApp.AccountKeeper.SetAccount(ctx, acc)
	accountGetter := ethante.NewCachedAccountGetter(ctx, myApp.AccountKeeper)

	ethTx := buildLegacyTx(t, privateKey, ethSigner)
	tx := buildTx(t, myApp, ethTx, ethSigner)
	require.NoError(t, ante.CheckAndSetEthSenderNonce(ctx, tx, myApp.AccountKeeper, false, accountGetter, ethSigner))
	newAcc := myApp.AccountKeeper.GetAccount(ctx, addr.Bytes())
	require.NotEmpty(t, newAcc)
	require.Equal(t, privateKey.PubKey(), newAcc.GetPubKey())

	require.NoError(t, newAcc.SetPubKey(nil))
	myApp.AccountKeeper.SetAccount(ctx, acc)

	ethTx = buildAccessListTx(t, privateKey, ethSigner)
	tx = buildTx(t, myApp, ethTx, ethSigner)
	require.NoError(t, ante.CheckAndSetEthSenderNonce(ctx, tx, myApp.AccountKeeper, false, accountGetter, ethSigner))
	newAcc = myApp.AccountKeeper.GetAccount(ctx, addr.Bytes())
	require.NotEmpty(t, newAcc)
	require.Equal(t, privateKey.PubKey(), newAcc.GetPubKey())

	require.NoError(t, newAcc.SetPubKey(nil))
	myApp.AccountKeeper.SetAccount(ctx, acc)

	ethTx = buildDynamicFeeTx(t, privateKey, ethSigner)
	tx = buildTx(t, myApp, ethTx, ethSigner)
	require.NoError(t, ante.CheckAndSetEthSenderNonce(ctx, tx, myApp.AccountKeeper, false, accountGetter, ethSigner))
	newAcc = myApp.AccountKeeper.GetAccount(ctx, addr.Bytes())
	require.NotEmpty(t, newAcc)
	require.Equal(t, privateKey.PubKey(), newAcc.GetPubKey())
}

func buildTx(t *testing.T, app *app.App, ethTx *ethtypes.Transaction, signer ethtypes.Signer) sdk.Tx {
	t.Helper()

	msg := &evmtypes.MsgEthereumTx{}
	require.NoError(t, msg.FromSignedEthereumTx(ethTx, signer))
	require.NoError(t, msg.ValidateBasic())
	option, err := codectypes.NewAnyWithValue(&evmtypes.ExtensionOptionsEthereumTx{})
	require.NoError(t, err)
	fees := make(sdk.Coins, 0)
	fee := msg.GetFee()
	feeAmt := sdkmath.NewIntFromBigInt(fee)
	if feeAmt.Sign() > 0 {
		fees = append(fees, sdk.NewCoin(fxtypes.DefaultDenom, feeAmt))
	}
	b := app.GetTxConfig().NewTxBuilder()
	builder, ok := b.(authtx.ExtensionOptionsTxBuilder)
	require.True(t, ok)
	builder.SetExtensionOptions(option)
	require.NoError(t, builder.SetMsgs(&evmtypes.MsgEthereumTx{From: msg.From, Raw: msg.Raw}))
	builder.SetFeeAmount(fees)
	builder.SetGasLimit(msg.GetGas())
	return builder.GetTx()
}

func buildLegacyTx(t *testing.T, privateKey cryptotypes.PrivKey, signer ethtypes.Signer) *ethtypes.Transaction {
	t.Helper()

	to := helpers.GenHexAddress()
	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    1,
		To:       &to,
		Value:    big.NewInt(1),
		Gas:      21000,
		GasPrice: big.NewInt(1),
	})
	sig, err := privateKey.Sign(signer.Hash(tx).Bytes())
	require.NoError(t, err)
	tx, err = tx.WithSignature(signer, sig)
	require.NoError(t, err)
	return tx
}

func buildAccessListTx(t *testing.T, privateKey cryptotypes.PrivKey, signer ethtypes.Signer) *ethtypes.Transaction {
	t.Helper()

	to := helpers.GenHexAddress()
	tx := ethtypes.NewTx(&ethtypes.AccessListTx{
		ChainID:    big.NewInt(530),
		Nonce:      2,
		To:         &to,
		Gas:        21000,
		GasPrice:   big.NewInt(1),
		AccessList: ethtypes.AccessList{{Address: helpers.GenHexAddress(), StorageKeys: []common.Hash{{0}}}},
	})
	sig, err := privateKey.Sign(signer.Hash(tx).Bytes())
	require.NoError(t, err)
	tx, err = tx.WithSignature(signer, sig)
	require.NoError(t, err)
	return tx
}

func buildDynamicFeeTx(t *testing.T, privateKey cryptotypes.PrivKey, signer ethtypes.Signer) *ethtypes.Transaction {
	t.Helper()

	to := helpers.GenHexAddress()
	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		To:        &to,
		Nonce:     3,
		Value:     big.NewInt(1),
		Gas:       21000,
		GasTipCap: big.NewInt(1),
		GasFeeCap: big.NewInt(1),
	})
	sig, err := privateKey.Sign(signer.Hash(tx).Bytes())
	require.NoError(t, err)
	tx, err = tx.WithSignature(signer, sig)
	require.NoError(t, err)
	return tx
}
