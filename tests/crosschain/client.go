package crosschain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"testing"

	fxtypes "github.com/functionx/fx-core/types"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	types2 "github.com/functionx/fx-core/x/ibc/applications/transfer/types"

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
)

const (
	defaultNodeGrpcUrl = "localhost:9090"
	defaultNodeRpcUrl  = "tcp://localhost:26657"

	defaultFxMnemonic = "dune antenna hood magic kit blouse film video another pioneer dilemma hobby message rug sail gas culture upgrade twin flag joke people general aunt"
	hdPath            = "m/44'/118'/0'/0/0"

	//pundixTokenContract        = "0x30dA8589BFa1E509A319489E014d384b87815D89"
	purseTokenContract         = "0x3f6795b8ABE0775a88973469909adE1405f7ac09"
	purseTokenSymbol           = "PRUSE"
	purseTokenChannelIBC       = "transfer/channel-0"
	defaultEthWalletPrivateKey = "b3f8605873861602b62617993fda26c00c057776934931a9d8cfa5d2e78fdc4a"
	chainName                  = "bsc"
)

var (
	purseDenom = types2.DenomTrace{
		Path:      purseTokenChannelIBC,
		BaseDenom: fmt.Sprintf("%s%s", chainName, purseTokenContract),
	}.IBCDenom()
)

var (
	txModeInfo = &tx.ModeInfo{
		Sum: &tx.ModeInfo_Single_{
			Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
		},
	}
)

type Client struct {
	t                     *testing.T
	ctx                   context.Context
	fxRpc                 *http.HTTP
	TxClient              tx.ServiceClient
	authQueryClient       authtypes.QueryClient
	crosschainQueryClient crosschaintypes.QueryClient
	bankQueryClient       banktypes.QueryClient
	fxPrivKey             *secp256k1.PrivKey
	encodingConfig        app.EncodingConfig
	ethPrivKey            *ecdsa.PrivateKey
	ethAddress            gethCommon.Address
	mutex                 *sync.Mutex
	gasFee                sdk.Coin
	chainName             string
	chainId               string
}

func (c *Client) FxAddress() sdk.AccAddress {
	return sdk.AccAddress(c.fxPrivKey.PubKey().Address())
}

func (c *Client) QueryFxLastEventNonce() uint64 {
	c.t.Helper()
	lastEventNonce, err := c.crosschainQueryClient.LastEventNonceByAddr(c.ctx, &crosschaintypes.QueryLastEventNonceByAddrRequest{ChainName: c.chainName, OrchestratorAddress: c.FxAddress().String()})
	if err != nil {
		c.t.Fatal(err)
	}
	return lastEventNonce.EventNonce + 1
}

func (c *Client) QueryObserver() *crosschaintypes.QueryLastObservedBlockHeightResponse {
	c.t.Helper()
	height, err := c.crosschainQueryClient.LastObservedBlockHeight(c.ctx, &crosschaintypes.QueryLastObservedBlockHeightRequest{ChainName: c.chainName})
	if err != nil {
		c.t.Fatal(err)
	}
	return height
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
		marshalIndent, err := c.encodingConfig.Amino.MarshalJSONIndent(msg, "", "\t")
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

func NewClient(t *testing.T, opts ...ClientOption) *Client {
	clientConn, err := grpcNewClient(defaultNodeGrpcUrl)
	if err != nil {
		t.Fatal(err)
	}
	httpClient, err := newHttpClient(defaultNodeRpcUrl)
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
	t.Logf("ethAddress:%v", ethAddress.String())
	status, err := httpClient.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	client := &Client{
		t:                     t,
		ctx:                   context.Background(),
		fxRpc:                 httpClient,
		TxClient:              tx.NewServiceClient(grpcClientConn),
		authQueryClient:       authtypes.NewQueryClient(grpcClientConn),
		bankQueryClient:       banktypes.NewQueryClient(grpcClientConn),
		crosschainQueryClient: crosschaintypes.NewQueryClient(grpcClientConn),
		fxPrivKey:             fxPrivateKey,
		encodingConfig:        app.MakeEncodingConfig(),
		ethPrivKey:            ethPrivateKey,
		ethAddress:            ethAddress,
		mutex:                 &sync.Mutex{},
		gasFee:                sdk.NewCoin("FX", sdk.NewInt(4000000000000)),
		chainName:             chainName,
		chainId:               status.NodeInfo.Network,
	}
	for _, opt := range opts {
		opt.apply(client)
	}
	return client
}

func grpcNewClient(grpcUrl string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	// http use  grpc.WithInsecure()
	opts = append(opts, grpc.WithInsecure())
	// https use grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "*.function.io"))
	//opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, strings.Split(grpcUrl, ":")[0])))
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
	gasPrice, _ := sdk.NewIntFromString("4000000000000")
	gasFeeAmount := gasPrice.Mul(sdk.NewInt(int64(gasLimit)))
	authInfo.Fee = &tx.Fee{
		Amount:   sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, gasFeeAmount)),
		GasLimit: gasLimit,
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
