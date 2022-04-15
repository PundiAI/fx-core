package gravity

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	othertypes "github.com/functionx/fx-core/x/other/types"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	"github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"

	"github.com/functionx/fx-core/app"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

const (
	// Fx chain data
	defaultFxGrpcUrl = "localhost:9090"
	defaultFxRpcUrl  = "tcp://localhost:26657"

	// Fx wallet data
	defaultFxMnemonic = "dune antenna hood magic kit blouse film video another pioneer dilemma hobby message rug sail gas culture upgrade twin flag joke people general aunt"
	hdPath            = "m/44'/118'/0'/0/0"

	// Eth data
	ethTokenContract           = "0x30dA8589BFa1E509A319489E014d384b87815D89"
	defaultEthWalletPrivateKey = "b3f8605873861602b62617993fda26c00c057776934931a9d8cfa5d2e78fdc4a"
)

var (
	txModeInfo = &tx.ModeInfo{
		Sum: &tx.ModeInfo_Single_{
			Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
		},
	}
)

type Client struct {
	t                  *testing.T
	ctx                context.Context
	fxRpc              *http.HTTP
	TxClient           tx.ServiceClient
	authQueryClient    authtypes.QueryClient
	gravityQueryClient gravitytypes.QueryClient
	bankQueryClient    banktypes.QueryClient
	otherQueryClient   othertypes.QueryClient
	gasPrice           sdk.Coin
	fxPrivKey          *secp256k1.PrivKey
	encodingConfig     app.EncodingConfig
	ethPrivKey         *ecdsa.PrivateKey
	ethAddress         gethCommon.Address
	mutex              *sync.Mutex
	chainId            string
}

func (c *Client) FxAddress() sdk.AccAddress {
	return sdk.AccAddress(c.fxPrivKey.PubKey().Address())
}

func (c *Client) QueryFxLastEventNonce() uint64 {
	c.t.Helper()
	lastEventNonce, err := c.gravityQueryClient.LastEventNonceByAddr(c.ctx, &gravitytypes.QueryLastEventNonceByAddrRequest{Address: c.FxAddress().String()})
	if err != nil {
		c.t.Fatal(err)
	}
	return lastEventNonce.EventNonce + 1
}

func NewClient(t *testing.T, opts ...ClientOption) *Client {
	clientConn, err := grpcNewClient(defaultFxGrpcUrl)
	if err != nil {
		t.Fatal(err)
	}
	httpClient, err := newHttpClient(defaultFxRpcUrl)
	if err != nil {
		t.Fatal(err)
	}
	fxPrivateKey, err := mnemonicToFxPrivKey(defaultFxMnemonic)
	if err != nil {
		t.Fatal(err)
	}
	grpcClientConn := clientConn
	ethPrivateKey, ethAddress, err := ethPrivateHexKeyToPrivate(defaultEthWalletPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	status, err := httpClient.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	client := &Client{
		t:                  t,
		ctx:                context.Background(),
		fxRpc:              httpClient,
		TxClient:           tx.NewServiceClient(grpcClientConn),
		authQueryClient:    authtypes.NewQueryClient(grpcClientConn),
		bankQueryClient:    banktypes.NewQueryClient(grpcClientConn),
		gravityQueryClient: gravitytypes.NewQueryClient(grpcClientConn),
		otherQueryClient:   othertypes.NewQueryClient(grpcClientConn),
		fxPrivKey:          fxPrivateKey,
		encodingConfig:     app.MakeEncodingConfig(),
		ethPrivKey:         ethPrivateKey,
		ethAddress:         ethAddress,
		mutex:              &sync.Mutex{},
		chainId:            status.NodeInfo.Network,
	}
	gasPrice, err := client.otherQueryClient.GasPrice(client.ctx, &othertypes.GasPriceRequest{})
	require.NoError(t, err)
	if !gasPrice.GetGasPrices().IsZero() {
		client.gasPrice = gasPrice.GetGasPrices()[0]
	}
	for _, opt := range opts {
		opt.apply(client)
	}
	return client
}

func grpcNewClient(grpcUrl string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	return grpc.Dial(grpcUrl, opts...)
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

func ethPrivateHexKeyToPrivate(privateKeyHex string) (*ecdsa.PrivateKey, gethCommon.Address, error) {
	privateKey, err := ethCrypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, gethCommon.Address{}, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, gethCommon.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}
	ethAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, ethAddress, nil
}

func newHttpClient(rpcUrl string) (*http.HTTP, error) {
	return http.New(rpcUrl, "/websocket")
}

func buildTxBodyAndTxAuthInfo(c *Client, msgList *[]sdk.Msg, accountNumber, accountSequence uint64) (txBodyBytes, txAuthInfoBytes []byte) {
	c.t.Helper()

	txBodyMessage := make([]*types.Any, 0)
	for i := 0; i < len(*msgList); i++ {
		msgAnyValue, err := types.NewAnyWithValue((*msgList)[i])
		if err != nil {
			c.t.Fatal(err)
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)

	}

	txBody := &tx.TxBody{
		Messages: txBodyMessage,
	}
	authInfo := &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: nil,
				ModeInfo:  txModeInfo,
				Sequence:  accountSequence,
			},
		},
		Fee: &tx.Fee{
			GasLimit: 0,
		},
	}
	txBodyBytes = mustProtoMarshal(c.t, txBody)

	signResultBytes := sign(c.t, c.fxPrivKey, &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       c.chainId,
		AccountNumber: accountNumber,
	})
	simulateResponse, err := c.TxClient.Simulate(context.Background(), &tx.SimulateRequest{
		Tx: &tx.Tx{
			Body:       txBody,
			AuthInfo:   authInfo,
			Signatures: [][]byte{signResultBytes},
		},
	})
	if err != nil {
		c.t.Fatal(err)
	}

	// adjustment gasLimit 1.3 .
	var gasLimit uint64 = simulateResponse.GasInfo.GasUsed * 13 / 10
	authInfo.Fee = &tx.Fee{
		Amount:   sdk.Coins{},
		GasLimit: gasLimit,
	}
	if !c.gasPrice.IsZero() {
		gasFeeAmount := c.gasPrice.Amount.Mul(sdk.NewInt(int64(gasLimit)))
		authInfo.Fee.Amount = sdk.NewCoins(sdk.NewCoin(c.gasPrice.Denom, gasFeeAmount))
	}
	txAuthInfoBytes = mustProtoMarshal(c.t, authInfo)
	return
}

func sign(t *testing.T, fxPrivKey *secp256k1.PrivKey, signDoc *tx.SignDoc) []byte {
	t.Helper()
	signDataBytes := mustProtoMarshal(t, signDoc)
	signResultBytes, err := fxPrivKey.Sign(signDataBytes)
	if err != nil {
		t.Fatal(err)
	}
	return signResultBytes
}

func (c *Client) BroadcastTx(msgList *[]sdk.Msg) string {
	c.t.Helper()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	fxAddress := c.FxAddress()
	accountResponse, err := c.authQueryClient.Account(c.ctx, &authtypes.QueryAccountRequest{Address: fxAddress.String()})
	if err != nil {
		c.t.Fatal(err)
	}
	var account authtypes.AccountI
	err = c.encodingConfig.InterfaceRegistry.UnpackAny(accountResponse.GetAccount(), &account)
	if err != nil {
		c.t.Fatal(err)
	}
	c.t.Logf("BroadcastTx address:%v, number:%v, sequence:%v\n", fxAddress.String(), account.GetAccountNumber(), account.GetSequence())
	c.t.Logf("msgs")
	for i, msg := range *msgList {
		if fmt.Sprintf("%v/%v", msg.Type(), msg.Route()) == "gov/submit_proposal" {
			c.t.Logf("gov submit proposal msg...")
			continue
		}
		marshalIndent, err := c.encodingConfig.Marshaler.MarshalJSON(msg)
		if err != nil {
			c.t.Fatal(err)
		}
		c.t.Logf("msg index:[%d], type:[%s], data:[%+v]", i, fmt.Sprintf("%s/%s", msg.Route(), msg.Type()), string(marshalIndent))
	}

	txBodyBytes, txAuthInfoBytes := buildTxBodyAndTxAuthInfo(c, msgList, account.GetAccountNumber(), account.GetSequence())

	signResultBytes := sign(c.t, c.fxPrivKey, &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       c.chainId,
		AccountNumber: account.GetAccountNumber(),
	})

	return broadcastTx(c.t, c.ctx, c.TxClient, &tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		Signatures:    [][]byte{signResultBytes},
	})

}

func broadcastTx(t *testing.T, ctx context.Context, txClient tx.ServiceClient, data *tx.TxRaw) string {
	t.Helper()
	broadcastData := mustProtoMarshal(t, data)
	broadcastTxResponse, err := txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
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
	t.Logf("broadcast tx success! height:%v txHash:%v gasUsed:%v\n", txResponse.Height, txResponse.TxHash, txResponse.GasUsed)
	return txResponse.TxHash
}

func mustProtoMarshal(t *testing.T, pb proto.Message) (bytes []byte) {
	t.Helper()
	bytes, err := proto.Marshal(pb)
	if err != nil {
		t.Fatal(err)
	}
	return
}

// ClientOption configures how we set up the connection.
type ClientOption interface {
	apply(*Client)
}

type funcDialOption struct {
	f func(*Client)
}

func (fdo *funcDialOption) apply(do *Client) {
	fdo.f(do)
}

func newFuncDialOption(f func(*Client)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

func WithEthPrivateKey(privateKey string) ClientOption {
	return newFuncDialOption(func(o *Client) {
		ethPrivateKey, ethAddress, err := ethPrivateHexKeyToPrivate(privateKey)
		if err != nil {
			panic(err)
		}
		o.ethPrivKey = ethPrivateKey
		o.ethAddress = ethAddress
	})
}

func WithFxMnemonic(mnemonic string) ClientOption {
	return newFuncDialOption(func(o *Client) {
		fxPrivKey, err := mnemonicToFxPrivKey(mnemonic)
		if err != nil {
			panic(err)
		}
		o.fxPrivKey = fxPrivKey
	})
}
