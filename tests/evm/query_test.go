package evm_test

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/functionx/fx-core/app/fxcore"
	_ "github.com/functionx/fx-core/app/fxcore"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/rpc/client/http"
	"testing"
)

func TestQueryBalance(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	addressBytes, err := sdk.AccAddressFromBech32("fx17ykqect7ee5e9r4l2end78d8gmp6mauzj87cwz")
	require.NoError(t, err)

	address := common.BytesToAddress(addressBytes)
	println(address.Hex())
	balanceRes, err := client.BalanceAt(context.Background(), address, nil)
	require.NoError(t, err)
	println(balanceRes.String())
}

func TestQueryTransaction(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	transactionReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x74a90ed91f42baa375804c22e2fa17087a6060bbca4ffb8f1e0fc1446883a0f7"))
	require.NoError(t, err)
	t.Logf("transactionReceipt:%+#v", transactionReceipt)
}

func TestQueryFxTxByEvmHash(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	transactionReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x74a90ed91f42baa375804c22e2fa17087a6060bbca4ffb8f1e0fc1446883a0f7"))
	require.NoError(t, err)
	t.Logf("transactionReceipt:%+#v", transactionReceipt)

	fxClient, err := http.New("http://0.0.0.0:26657", "/websocket")
	require.NoError(t, err)
	evmHashBlockNumber := transactionReceipt.BlockNumber.Int64()
	blockData, err := fxClient.Block(context.Background(), &evmHashBlockNumber)
	require.NoError(t, err)
	require.True(t, uint(len(blockData.Block.Txs)) > transactionReceipt.TransactionIndex)
	fxTx := blockData.Block.Txs[transactionReceipt.TransactionIndex]
	encodingConfig := fxcore.MakeEncodingConfig()
	tx, err := encodingConfig.TxConfig.TxDecoder()(fxTx)
	require.NoError(t, err)
	txJsonStr, err := encodingConfig.TxConfig.TxJSONEncoder()(tx)
	//marshalIndent, err := json.MarshalIndent(string(txJsonStr), "", "  ")
	//require.NoError(t, err)
	t.Logf("\nTxHash:%x\nData:\n%v", fxTx.Hash(), string(txJsonStr))

}
func TestMnemonicToFxPrivate(t *testing.T) {
	privKey, err := mnemonicToFxPrivKey("december slow blue fury silly bread friend unknown render resource dry buyer brand final abstract gallery slow since hood shadow neglect travel convince foil")
	require.NoError(t, err)
	t.Logf("%x", privKey.Key)
}

func TestEthPrivateKeyToAddress(t *testing.T) {
	//privateKey, err := crypto.GenerateKey()
	//require.NoError(t, err)
	//fromECDSA := crypto.FromECDSA(privateKey)
	//t.Logf("fromEc:%x", fromECDSA)

	// 1ce31354ff0a3f057c9b70ebbbdacb68ace4bf9c008ac722f2b996328ab3ca08
	hexPrivKey := "86b87f127b6e0901f7f00aa77b6c82624847f2628a901bf1833b2d48883b73d3"
	recoverPrivKey, err := crypto.HexToECDSA(hexPrivKey)
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(recoverPrivKey.PublicKey)
	t.Logf("Eth address:%v, FxAddress:%v", address.Hex(), sdk.AccAddress(address.Bytes()).String())
}

func TestEthAddressToFxAddress(t *testing.T) {
	ethAddress := common.HexToAddress("0xf12C0Ce17eCE69928ebf5666Df1Da746c3adf782")
	t.Logf("%o", ethAddress.Bytes())
	t.Logf("EthAddress:%v, FxAddress:%v", ethAddress.Hex(), sdk.AccAddress(ethAddress.Bytes()).String())
}

func TestFxAddressToEthAddress(t *testing.T) {
	fxAddress, err := sdk.AccAddressFromBech32("fx10kg059hhxc2pevxssszunvgc70jpmxsjal4xf6")
	require.NoError(t, err)
	ethAddress := common.BytesToAddress(fxAddress)
	t.Logf("EthAddress:%v, FxAddress:%v", ethAddress.Hex(), sdk.AccAddress(ethAddress.Bytes()).String())
}

func mnemonicToFxPrivKey(mnemonic string) (*secp256k1.PrivKey, error) {
	algo := hd.Secp256k1
	bytes, err := algo.Derive()(mnemonic, "", "m/44'/118'/0'/0/0")
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(bytes)
	priv, ok := privKey.(*secp256k1.PrivKey)
	if !ok {
		return nil, fmt.Errorf("not secp256k1.PrivKey")
	}
	return priv, nil
}
