package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/app"
	"github.com/functionx/fx-core/v8/client"
	"github.com/functionx/fx-core/v8/client/grpc"
	"github.com/functionx/fx-core/v8/client/jsonrpc"
	fxauth "github.com/functionx/fx-core/v8/server/grpc/auth"
	"github.com/functionx/fx-core/v8/testutil"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/testutil/network"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

type rpcTestClient interface {
	AppVersion() (string, error)
	GetChainId() (chain string, err error)
	GetBlockHeight() (int64, error)
	GetMintDenom() (denom string, err error)
	GetGasPrices() (sdk.Coins, error)
	GetAddressPrefix() (prefix string, err error)
	GetModuleAccounts() ([]sdk.AccountI, error)
	QueryAccount(address string) (sdk.AccountI, error)
	QueryBalance(address, denom string) (sdk.Coin, error)
	QueryBalances(address string) (sdk.Coins, error)
	QuerySupply() (sdk.Coins, error)
	BuildTxRaw(privKey cryptotypes.PrivKey, msgs []sdk.Msg, gasLimit, timeout uint64, memo string) (*tx.TxRaw, error)
	EstimatingGas(raw *tx.TxRaw) (*sdk.GasInfo, error)
	BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error)
	WaitMined(txHash string, timeout, pollInterval time.Duration) (*sdk.TxResponse, error)
	TxByHash(txHash string) (*sdk.TxResponse, error)
}

type rpcTestSuite struct {
	suite.Suite

	network *network.Network
}

func TestRPCSuite(t *testing.T) {
	suite.Run(t, new(rpcTestSuite))
}

func (suite *rpcTestSuite) SetupSuite() {
	suite.T().Log("setting up integration test suite")

	fxtypes.SetConfig(true)
	cfg := testutil.DefaultNetworkConfig(func(config *network.Config) {
		// config.EnableTMLogging = true
	})
	cfg.TimeoutCommit = 100 * time.Millisecond
	cfg.NumValidators = 1
	cfg.Mnemonics = append(cfg.Mnemonics, helpers.NewMnemonic())

	suite.network = network.New(suite.T(), cfg)

	_, err := suite.network.WaitForHeight(1)
	suite.Require().NoError(err)

	suite.FirstValidatorTransferTo(1, sdkmath.NewInt(1_000).MulRaw(1e18))
}

func (suite *rpcTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *rpcTestSuite) GetFirstValidator() *network.Validator {
	return suite.network.Validators[0]
}

func (suite *rpcTestSuite) GetFirstValPrivKey() cryptotypes.PrivKey {
	return suite.GetPrivKeyByIndex(hd.Secp256k1Type, 0)
}

func (suite *rpcTestSuite) GetPrivKeyByIndex(algo hd.PubKeyType, index uint32) cryptotypes.PrivKey {
	privKey, err := helpers.PrivKeyFromMnemonic(suite.network.Config.Mnemonics[0], algo, 0, index)
	suite.Require().NoError(err)
	return privKey
}

func (suite *rpcTestSuite) GetClients() []rpcTestClient {
	validator := suite.GetFirstValidator()
	suite.True(validator.AppConfig.GRPC.Enable)
	grpcClient, err := grpc.DailClient(fmt.Sprintf("http://%s", validator.AppConfig.GRPC.Address))
	suite.Require().NoError(err)
	rpcAddress := validator.Ctx.Config.RPC.ListenAddress
	wsClient, err := jsonrpc.NewWsClient(context.Background(), rpcAddress+"/websocket")
	suite.Require().NoError(err)
	return []rpcTestClient{
		grpcClient,
		jsonrpc.NewNodeRPC(jsonrpc.NewClient(rpcAddress)),
		jsonrpc.NewNodeRPC(wsClient),
	}
}

func (suite *rpcTestSuite) GetClient() rpcTestClient {
	clients := suite.GetClients()
	return clients[tmrand.Int()%len(clients)]
}

func (suite *rpcTestSuite) FirstValidatorTransferTo(index uint32, amount sdkmath.Int) {
	cli := suite.GetClient()
	valKey := suite.GetFirstValPrivKey()
	toAccountKey := suite.GetPrivKeyByIndex(hd.Secp256k1Type, index)
	from := sdk.AccAddress(valKey.PubKey().Address().Bytes())
	account, chainId, gasPrice, err := client.GetChainInfo(cli, from.String())
	suite.Require().NoError(err)
	msgs := []sdk.Msg{banktypes.NewMsgSend(
		valKey.PubKey().Address().Bytes(),
		toAccountKey.PubKey().Address().Bytes(),
		sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)),
	)}
	txRaw, err := client.BuildTxRaw(chainId, account.GetSequence(), account.GetAccountNumber(), valKey, msgs, gasPrice, 250000, 0, "")
	suite.Require().NoError(err)
	txResponse, err := cli.BroadcastTx(txRaw)
	suite.Require().NoError(err)
	suite.Equal(uint32(0), txResponse.Code)
	suite.Less(txResponse.GasUsed, int64(100000))
	txResponse, err = cli.WaitMined(txResponse.TxHash, time.Second, 100*time.Millisecond)
	suite.Require().NoError(err)
	suite.Equal(uint32(0), txResponse.Code)
}

func (suite *rpcTestSuite) TestClient_Tx() {
	privKey := suite.GetPrivKeyByIndex(hd.Secp256k1Type, 1)

	clients := suite.GetClients()
	for i := 0; i < len(clients); i++ {
		cli := clients[i]
		toAddress := sdk.AccAddress(helpers.NewPriKey().PubKey().Address())
		txRaw, err := cli.BuildTxRaw(
			privKey,
			[]sdk.Msg{banktypes.NewMsgSend(
				privKey.PubKey().Address().Bytes(),
				toAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1))),
			)},
			0, 0, "",
		)
		suite.Require().NoError(err)

		gas, err := cli.EstimatingGas(txRaw)
		suite.Require().NoError(err)
		suite.Less(gas.GasUsed, uint64(100000))

		txResponse, err := cli.BroadcastTx(txRaw)
		suite.Require().NoError(err)
		suite.Equal(uint32(0), txResponse.Code)
		suite.Less(txResponse.GasUsed, int64(100000))

		txResponse, err = cli.WaitMined(txResponse.TxHash, time.Second, 100*time.Millisecond)
		suite.Require().NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		txRes, err := cli.TxByHash(txResponse.TxHash)
		suite.Require().NoError(err)
		suite.Equal(txResponse, txRes)

		account, err := cli.QueryAccount(toAddress.String())
		suite.Require().NoError(err)
		// acconts is
		// 0. initAccount
		// 1.fee_collector + 2.distribution + 3.bonded_tokens_pool + 4.not_bonded_tokens_pool + 5.gov + 6.mint + 7.autytypes.NewModuleAddress(crosschain)
		// 8.evm 9.0x..1001 10.0x..1002 11.erc20 12.wfx-contract
		suite.Equal(authtypes.NewBaseAccount(toAddress, nil, uint64(15+i), 0), account)
	}

	ethPrivKey := suite.GetPrivKeyByIndex(hd2.EthSecp256k1Type, 0)

	ethAddress := sdk.AccAddress(ethPrivKey.PubKey().Address().Bytes())

	for i := 0; i < len(clients); i++ {
		cli := clients[i]
		txRaw, err := cli.BuildTxRaw(
			privKey,
			[]sdk.Msg{banktypes.NewMsgSend(
				privKey.PubKey().Address().Bytes(),
				ethAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10).MulRaw(1e18))),
			)},
			0, 0, "",
		)
		suite.Require().NoError(err)

		gas, err := cli.EstimatingGas(txRaw)
		suite.Require().NoError(err)
		suite.Less(gas.GasUsed, uint64(100000))

		txResponse, err := cli.BroadcastTx(txRaw)
		suite.Require().NoError(err)
		suite.Equal(uint32(0), txResponse.Code)
		suite.Less(txResponse.GasUsed, int64(100000))

		txResponse, err = cli.WaitMined(txResponse.TxHash, time.Second, 100*time.Millisecond)
		suite.Require().NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		account, err := cli.QueryAccount(ethAddress.String())
		suite.Require().NoError(err)
		suite.Equal(authtypes.NewBaseAccount(ethAddress, nil, uint64(18), 0), account)
	}

	for i := 0; i < len(clients); i++ {
		cli := clients[i]
		toAddress := sdk.AccAddress(helpers.NewPriKey().PubKey().Address())
		txRaw, err := cli.BuildTxRaw(
			ethPrivKey,
			[]sdk.Msg{banktypes.NewMsgSend(
				ethPrivKey.PubKey().Address().Bytes(),
				toAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1))),
			)},
			0, 0, "",
		)
		suite.Require().NoError(err)

		gas, err := cli.EstimatingGas(txRaw)
		suite.Require().NoError(err)
		suite.Less(gas.GasUsed, uint64(100000))

		txResponse, err := cli.BroadcastTx(txRaw)
		suite.Require().NoError(err)
		suite.Equal(uint32(0), txResponse.Code)
		suite.Less(txResponse.GasUsed, int64(100000))

		txResponse, err = cli.WaitMined(txResponse.TxHash, time.Second, 100*time.Millisecond)
		suite.Require().NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		account, err := cli.QueryAccount(ethAddress.String())
		suite.Require().NoError(err)
		baseAccount, ok := account.(*authtypes.BaseAccount)
		suite.True(ok)
		if baseAccount.PubKey.TypeUrl != "" {
			pubAny, err := types.NewAnyWithValue(ethPrivKey.PubKey())
			suite.Require().NoError(err)
			suite.Equal("/"+proto.MessageName(&ethsecp256k1.PubKey{}), baseAccount.PubKey.TypeUrl)
			suite.Equal(pubAny, baseAccount.PubKey)
		}
		suite.Equal(uint64(i+1), account.GetSequence())
	}
}

func (suite *rpcTestSuite) TestClient_Query() {
	feeCollectorAddr, err := sdk.AccAddressFromHexUnsafe("f1829676db577682e944fc3493d451b67ff3e29f")
	suite.Require().NoError(err)
	tests := []struct {
		funcName string
		params   []interface{}
		wantRes  []interface{}
	}{
		{
			funcName: "GetChainId",
			params:   []interface{}{},
			wantRes:  []interface{}{fxtypes.MainnetChainId, nil},
		},
		{
			funcName: "GetMintDenom",
			params:   []interface{}{},
			wantRes:  []interface{}{fxtypes.DefaultDenom, nil},
		},
		{
			funcName: "GetAddressPrefix",
			params:   []interface{}{},
			wantRes:  []interface{}{"fx", nil},
		},
		{
			funcName: "AppVersion",
			params:   []interface{}{},
			wantRes:  []interface{}{"", nil},
		},
		{
			funcName: "GetBlockHeight",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(height int64, err error) {
					suite.Require().NoError(err)
					suite.GreaterOrEqual(height, int64(1))
				},
			},
		},
		{
			funcName: "QuerySupply",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(supply sdk.Coins, err error) {
					suite.Require().NoError(err)
					supply.IsAllGTE(
						sdk.Coins{
							sdk.Coin{
								Denom:  fxtypes.DefaultDenom,
								Amount: sdkmath.NewInt(500_000).MulRaw(1e18),
							},
						},
					)
				},
			},
		},
		{
			funcName: "GetGasPrices",
			params:   []interface{}{},
			wantRes: []interface{}{
				sdk.Coins{
					sdk.Coin{
						Denom:  fxtypes.DefaultDenom,
						Amount: sdkmath.NewInt(4).MulRaw(1e12),
					},
				},
				nil,
			},
		},
		{
			funcName: "QueryAccount",
			params: []interface{}{
				suite.GetFirstValidator().Address.String(),
			},
			wantRes: []interface{}{
				authtypes.NewBaseAccount(
					suite.GetFirstValPrivKey().PubKey().Address().Bytes(),
					suite.GetFirstValPrivKey().PubKey(),
					0,
					2,
				),
				nil,
			},
		},
		{
			funcName: "QueryAccount",
			params: []interface{}{
				authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
			},
			wantRes: []interface{}{
				authtypes.NewModuleAccount(
					authtypes.NewBaseAccount(
						feeCollectorAddr,
						nil,
						1,
						0,
					),
					authtypes.FeeCollectorName,
				),
				nil,
			},
		},
		{
			funcName: "QueryBalance",
			params:   []interface{}{suite.GetFirstValidator().Address.String(), fxtypes.DefaultDenom},
			wantRes: []interface{}{
				sdk.Coin{
					Denom:  fxtypes.DefaultDenom,
					Amount: sdkmath.NewInt(488999).MulRaw(1e18),
				},
				nil,
			},
		},
		{
			funcName: "QueryBalance",
			params:   []interface{}{helpers.GenAccAddress().String(), fxtypes.DefaultDenom},
			wantRes: []interface{}{
				sdk.Coin{
					Denom:  fxtypes.DefaultDenom,
					Amount: sdkmath.ZeroInt(),
				},
				nil,
			},
		},
		{
			funcName: "QueryBalances",
			params:   []interface{}{suite.GetFirstValidator().Address.String()},
			wantRes: []interface{}{
				sdk.Coins{
					sdk.Coin{
						Denom:  fxtypes.DefaultDenom,
						Amount: sdkmath.NewInt(488999).MulRaw(1e18),
					},
				},
				nil,
			},
		},
	}
	clients := suite.GetClients()
	for _, tt := range tests {
		suite.Run(tt.funcName, func() {
			for i := 0; i < len(clients); i++ {
				typeOf := reflect.TypeOf(clients[i])
				method, is := typeOf.MethodByName(tt.funcName)
				suite.True(is)
				params := make([]reflect.Value, len(tt.params)+1)
				params[0] = reflect.ValueOf(clients[i])
				for i := 1; i < len(params); i++ {
					params[i] = reflect.ValueOf(tt.params[i-1])
				}
				results := method.Func.Call(params)
				if len(tt.wantRes) == 1 {
					wantResTf := reflect.ValueOf(tt.wantRes[0])
					suite.Equal(reflect.Func, wantResTf.Kind())
					wantResTf.Call(results)
				} else {
					for i := 0; i < len(results); i++ {
						suite.EqualValues(
							fmt.Sprintf("%v", tt.wantRes[i]),
							fmt.Sprintf("%v", results[i]),
						)
					}
				}
			}
		})
	}
}

func (suite *rpcTestSuite) TestClient_GetModuleAccounts() {
	clients := suite.GetClients()
	for i := 0; i < len(clients); i++ {
		accounts, err := clients[i].GetModuleAccounts()
		suite.Require().NoError(err)
		suite.Len(accounts, 18)
		suite.Equal(len(app.GetMaccPerms()), len(accounts))
	}
}

func (suite *rpcTestSuite) TestTmClient() {
	validator := suite.GetFirstValidator()
	tmRPC := validator.RPCClient
	callTmRPC := func(funcName string, params []interface{}) []reflect.Value {
		typeOf := reflect.TypeOf(tmRPC)
		method, is := typeOf.MethodByName(funcName)
		suite.True(is)
		callParams := make([]reflect.Value, len(params))
		for i, param := range params {
			callParams[i] = reflect.ValueOf(param)
		}
		callParams = append([]reflect.Value{reflect.ValueOf(tmRPC), reflect.ValueOf(context.Background())}, callParams...)
		results := method.Func.Call(callParams)
		return results
	}

	nodeRPC := jsonrpc.NewNodeRPC(jsonrpc.NewClient(validator.Ctx.Config.RPC.ListenAddress))
	callNodeRPC := func(funcName string, params []interface{}) []reflect.Value {
		typeOf := reflect.TypeOf(nodeRPC)
		method, is := typeOf.MethodByName(funcName)
		suite.True(is)
		callParams := make([]reflect.Value, len(params))
		for i, param := range params {
			callParams[i] = reflect.Indirect(reflect.ValueOf(param))
		}
		callParams = append([]reflect.Value{reflect.ValueOf(nodeRPC)}, callParams...)
		results := method.Func.Call(callParams)
		return results
	}

	height := int64(1)
	limit := 1
	tests := []struct {
		funcName string
		params   []interface{}
		wantRes  []interface{}
	}{
		// ABCIClient
		{
			funcName: "ABCIInfo",
			params:   []interface{}{},
		},
		// HistoryClient
		{
			funcName: "Genesis",
			params:   []interface{}{},
		},
		{
			funcName: "BlockchainInfo",
			params:   []interface{}{int64(1), int64(1)},
			wantRes: []interface{}{
				func(res1 *ctypes.ResultBlockchainInfo, err1 error, res2 *ctypes.ResultBlockchainInfo, err2 error) {
					suite.Require().NoError(err1)
					suite.Require().NoError(err2)
					data1, _ := json.Marshal(res1.BlockMetas)
					data2, _ := json.Marshal(res2.BlockMetas)
					suite.Equal(data1, data2)
				},
			},
		},
		// StatusClient
		{
			funcName: "Status",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(res1 *ctypes.ResultStatus, err1 error, res2 *ctypes.ResultStatus, err2 error) {
					suite.Require().NoError(err1)
					suite.Require().NoError(err2)
					suite.EqualValues(res1.NodeInfo, res2.NodeInfo)
					suite.EqualValues(res1.ValidatorInfo, res2.ValidatorInfo)
				},
			},
		},
		// NetworkClient
		{
			funcName: "NetInfo",
			params:   []interface{}{},
		},
		{
			funcName: "DumpConsensusState",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(res1 *ctypes.ResultDumpConsensusState, err1 error, res2 *ctypes.ResultDumpConsensusState, err2 error) {
					suite.Require().NoError(err1)
					suite.Require().NoError(err2)
					suite.EqualValues(len(res1.Peers), len(res2.Peers))
				},
			},
		},
		{
			funcName: "ConsensusState",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(_ *ctypes.ResultConsensusState, err1 error, _ *ctypes.ResultConsensusState, err2 error) {
					suite.Require().NoError(err1)
					suite.Require().NoError(err2)
				},
			},
		},
		{
			funcName: "ConsensusParams",
			params:   []interface{}{&height},
		},
		{
			funcName: "Health",
			params:   []interface{}{},
		},
		// MempoolClient
		{
			funcName: "UnconfirmedTxs",
			params:   []interface{}{&limit},
		},
		{
			funcName: "NumUnconfirmedTxs",
			params:   []interface{}{},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.funcName, func() {
			wg := sync.WaitGroup{}
			wg.Add(1)
			resultChan := make(chan []reflect.Value, 2)
			go func() {
				defer wg.Done()
				resultChan <- callTmRPC(tt.funcName, tt.params)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				resultChan <- callNodeRPC(tt.funcName, tt.params)
			}()
			wg.Wait()

			result1 := <-resultChan
			result2 := <-resultChan
			suite.Equal(len(result1), len(result2))
			if len(tt.wantRes) == 1 {
				wantResTf := reflect.ValueOf(tt.wantRes[0])
				suite.Equal(reflect.Func, wantResTf.Kind())
				wantResTf.Call(append(result1, result2...))
			} else {
				for i := 0; i < len(result1); i++ {
					data1, err1 := json.Marshal(reflect.Indirect(result1[i]).Interface())
					suite.Require().NoError(err1)
					data2, err2 := json.Marshal(reflect.Indirect(result2[i]).Interface())
					suite.Require().NoError(err2)
					suite.JSONEq(string(data1), string(data2))
				}
			}
			close(resultChan)
		})
	}
}

func (suite *rpcTestSuite) TestClient_WithBlockHeight() {
	key := suite.GetPrivKeyByIndex(hd.Secp256k1Type, 1)
	clients := suite.GetClients()
	for _, cli := range clients {
		balances, err := cli.QueryBalances(sdk.AccAddress(key.PubKey().Address().Bytes()).String())
		suite.Require().NoError(err)
		suite.True(balances.IsAllPositive())

		if rpc, ok := cli.(*jsonrpc.NodeRPC); ok {
			cli = rpc.WithBlockHeight(1)
		}
		if rpc, ok := cli.(*grpc.Client); ok {
			cli = rpc.WithBlockHeight(1)
		}

		balances, err = cli.QueryBalances(sdk.AccAddress(key.PubKey().Address().Bytes()).String())
		suite.Require().NoError(err)
		suite.False(balances.IsAllPositive())
	}
}

func (suite *rpcTestSuite) TestGRPCClient_ConvertAddress() {
	validator := suite.GetFirstValidator()
	cli := fxauth.NewQueryClient(validator.ClientCtx)
	res, err := cli.ConvertAddress(context.Background(), &fxauth.ConvertAddressRequest{
		Address: validator.Address.String(),
		Prefix:  sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
	})
	suite.Require().NoError(err)
	suite.Equal(res.Address, validator.ValAddress.String())
}
