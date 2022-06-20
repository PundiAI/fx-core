package migrate

import (
	"context"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	cryptohd "github.com/evmos/ethermint/crypto/hd"

	"github.com/functionx/fx-core/app"
	othertypes "github.com/functionx/fx-core/x/other/types"
)

const (
	defaultFxGrpcUrl  = "localhost:9090"
	hdPath            = "m/44'/118'/0'/0/0"
	defaultFxMnemonic = "dune antenna hood magic kit blouse film video another pioneer dilemma hobby message rug sail gas culture upgrade twin flag joke people general aunt"
	defaultChainId    = "fxcore"

	toAddressMnemonic = "cloud exclude pass crime pill garment feature cancel affair dream aware reunion reward wide autumn strike edit now crop fever swift easy price about"
	val2Mnemonic      = "position shrug hamster range arena cash uncover execute piece cherry unknown wonder obscure remain coach clump arrest park cover cotton ginger educate radio spawn"
)

var (
	TxModeInfo = &tx.ModeInfo{
		Sum: &tx.ModeInfo_Single_{
			Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
		},
	}
)

type Client struct {
	chainId    string
	privateKey cryptotypes.PrivKey
	grpcConn   *grpc.ClientConn
	t          *testing.T
}

func grpcNewClient(grpcUrl string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	return grpc.Dial(grpcUrl, opts...)
}

func NewClient(t *testing.T) *Client {
	conn, err := grpcNewClient(defaultFxGrpcUrl)
	require.NoError(t, err)

	priKey, err := mnemonicToFxPrivKey(defaultFxMnemonic)
	require.NoError(t, err)

	return &Client{grpcConn: conn, privateKey: priKey, chainId: defaultChainId, t: t}
}

func (c *Client) SetPrivateKey(key cryptotypes.PrivKey) {
	c.privateKey = key
}
func (c *Client) FxAddress() sdk.AccAddress {
	return sdk.AccAddress(c.privateKey.PubKey().Address())
}
func (c *Client) AuthQuery() authtypes.QueryClient {
	return authtypes.NewQueryClient(c.grpcConn)
}
func (c *Client) StakingQuery() stakingtypes.QueryClient {
	return stakingtypes.NewQueryClient(c.grpcConn)
}
func (c *Client) DistrQuery() distrtypes.QueryClient {
	return distrtypes.NewQueryClient(c.grpcConn)
}
func (c *Client) OtherQuery() othertypes.QueryClient {
	return othertypes.NewQueryClient(c.grpcConn)
}
func (c *Client) BankQuery() banktypes.QueryClient {
	return banktypes.NewQueryClient(c.grpcConn)
}
func (c *Client) GovQuery() govtypes.QueryClient {
	return govtypes.NewQueryClient(c.grpcConn)
}
func (c *Client) ServiceClient() tx.ServiceClient {
	return tx.NewServiceClient(c.grpcConn)
}
func (c *Client) BroadcastTx(msgList ...sdk.Msg) string {
	c.t.Helper()
	fxAddress := c.FxAddress()
	accountResponse, err := c.AuthQuery().Account(context.Background(), &authtypes.QueryAccountRequest{Address: fxAddress.String()})
	if err != nil {
		c.t.Fatal(err)
	}

	encodingConfig := app.MakeEncodingConfig()
	var account authtypes.AccountI
	err = encodingConfig.InterfaceRegistry.UnpackAny(accountResponse.GetAccount(), &account)
	if err != nil {
		c.t.Fatal(err)
	}
	//c.t.Logf("BroadcastTx address:%v, number:%v, sequence:%v\n", fxAddress.String(), account.GetAccountNumber(), account.GetSequence())
	//for i, msg := range msgList {
	//	if fmt.Sprintf("%v/%v", msg.Type(), msg.Route()) == "gov/submit_proposal" {
	//		c.t.Logf("gov submit proposal msg...")
	//		continue
	//	}
	//	marshalIndent, err := encodingConfig.Marshaler.MarshalJSON(msg)
	//	require.NoError(c.t, err)
	//c.t.Logf("msg index:[%d], type:[%s], data:[%+v]", i, fmt.Sprintf("%s/%s", msg.Route(), msg.Type()), string(marshalIndent))
	//}

	txBodyBytes, txAuthInfoBytes := c.BuildTx(msgList, account.GetAccountNumber(), account.GetSequence())

	signResultBytes := sign(c.t, c.privateKey, &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       c.chainId,
		AccountNumber: account.GetAccountNumber(),
	})

	return broadcastTx(c.t, c.ServiceClient(), &tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		Signatures:    [][]byte{signResultBytes},
	})
}
func (c *Client) BuildTx(msgList []sdk.Msg, accountNumber, accountSequence uint64) (txBodyBytes, txAuthInfoBytes []byte) {
	c.t.Helper()
	txBodyMessage := make([]*types.Any, 0)
	for i := 0; i < len(msgList); i++ {
		msgAnyValue, err := types.NewAnyWithValue((msgList)[i])
		if err != nil {
			c.t.Fatal(err)
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)

	}
	txBody := &tx.TxBody{
		Messages: txBodyMessage,
	}
	any, err := codectypes.NewAnyWithValue(c.privateKey.PubKey())
	require.NoError(c.t, err)

	authInfo := &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: any,
				ModeInfo:  TxModeInfo,
				Sequence:  accountSequence,
			},
		},
		Fee: &tx.Fee{
			GasLimit: 0,
		},
	}
	txBodyBytes = mustProtoMarshal(c.t, txBody)

	signResultBytes := sign(c.t, c.privateKey, &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       c.chainId,
		AccountNumber: accountNumber,
	})
	simulateResponse, err := c.ServiceClient().Simulate(context.Background(), &tx.SimulateRequest{
		Tx: &tx.Tx{
			Body:       txBody,
			AuthInfo:   authInfo,
			Signatures: [][]byte{signResultBytes},
		},
	})
	require.NoError(c.t, err)

	// adjustment gasLimit 1.3 .
	var gasLimit uint64 = simulateResponse.GasInfo.GasUsed * 13 / 10
	authInfo.Fee = &tx.Fee{
		Amount:   sdk.Coins{},
		GasLimit: gasLimit,
	}
	gasPrice, err := c.OtherQuery().GasPrice(context.Background(), &othertypes.GasPriceRequest{})
	require.NoError(c.t, err)

	if !gasPrice.GasPrices.IsZero() {
		gasFeeAmount := gasPrice.GasPrices[0].Amount.Mul(sdk.NewInt(int64(gasLimit)))
		authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(gasPrice.GasPrices[0].Denom, gasFeeAmount))
	}
	txAuthInfoBytes = mustProtoMarshal(c.t, authInfo)
	return
}

func broadcastTx(t *testing.T, txClient tx.ServiceClient, data *tx.TxRaw) string {
	t.Helper()
	broadcastData := mustProtoMarshal(t, data)
	broadcastTxResponse, err := txClient.BroadcastTx(context.Background(), &tx.BroadcastTxRequest{
		TxBytes: broadcastData,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_BLOCK,
	})
	if err != nil {
		t.Fatal(err)
	}
	txResponse := broadcastTxResponse.TxResponse
	if txResponse.Code != 0 {
		t.Fatalf("broadcast tx fail!!!\ncode:%v, codespace:%v\n%v\n", txResponse.Code, txResponse.Codespace, txResponse.String())
	}
	//t.Logf("broadcast tx success! height:%v txHash:%v gasUsed:%v\n", txResponse.Height, txResponse.TxHash, txResponse.GasUsed)
	return txResponse.TxHash
}
func sign(t *testing.T, fxPrivKey cryptotypes.PrivKey, signDoc *tx.SignDoc) []byte {
	t.Helper()
	signDataBytes := mustProtoMarshal(t, signDoc)
	signResultBytes, err := fxPrivKey.Sign(signDataBytes)
	require.NoError(t, err)
	return signResultBytes
}

func mnemonicToFxPrivKey(mnemonic string) (*secp256k1.PrivKey, error) {
	algo := hd.Secp256k1
	bytes, err := algo.Derive()(mnemonic, "", hdPath)
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

// nolint
func mnemonicToSecp256k1(mnemonic string) (cryptotypes.PrivKey, error) {
	algo := hd.Secp256k1
	bytes, err := algo.Derive()(mnemonic, "", hdPath)
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
func mnemonicToEthSecp256k1(mnemonic string) (cryptotypes.PrivKey, error) {
	algo := cryptohd.EthSecp256k1
	bytes, err := algo.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(bytes)
	priv, ok := privKey.(*ethsecp256k1.PrivKey)
	if !ok {
		return nil, fmt.Errorf("not eth_secp256k1.PrivKey")
	}
	return priv, nil
}
func mnemonicToEd25519(mnemonic string) (cryptotypes.PrivKey, error) {
	priv := ed25519.GenPrivKeyFromSecret([]byte(mnemonic))
	return priv, nil
}

func mustProtoMarshal(t *testing.T, pb proto.Message) (bytes []byte) {
	t.Helper()
	bytes, err := proto.Marshal(pb)
	if err != nil {
		t.Fatal(err)
	}
	return
}
