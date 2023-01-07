package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/client/grpc"
	"github.com/functionx/fx-core/v3/client/jsonrpc"
	"github.com/functionx/fx-core/v3/testutil"
	"github.com/functionx/fx-core/v3/testutil/network"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

type rpcTestClient interface {
	AppVersion() (string, error)
	GetChainId() (chain string, err error)
	GetBlockHeight() (int64, error)
	GetMintDenom() (denom string, err error)
	GetGasPrices() (sdk.Coins, error)
	GetAddressPrefix() (prefix string, err error)
	QueryAccount(address string) (authtypes.AccountI, error)
	QueryBalance(address string, denom string) (sdk.Coin, error)
	QueryBalances(address string) (sdk.Coins, error)
	QuerySupply() (sdk.Coins, error)
	BuildTx(privKey cryptotypes.PrivKey, msgs []sdk.Msg) (*tx.TxRaw, error)
	EstimatingGas(raw *tx.TxRaw) (*sdk.GasInfo, error)
	BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error)
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

	cfg := testutil.DefaultNetworkConfig()
	cfg.TimeoutCommit = time.Millisecond
	cfg.NumValidators = 1
	cfg.Mnemonics = append(cfg.Mnemonics, helpers.NewMnemonic())

	baseDir, err := os.MkdirTemp(suite.T().TempDir(), cfg.ChainID)
	suite.Require().NoError(err)
	suite.network, err = network.New(suite.T(), baseDir, cfg)
	suite.Require().NoError(err)

	suite.FirstValidatorTransferTo(1, sdk.NewInt(1_000).MulRaw(1e18))
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

func (suite *rpcTestSuite) GetFirstValiPrivKey() cryptotypes.PrivKey {
	return suite.GetPrivKeyByIndex(hd.Secp256k1Type, 0)
}

func (suite *rpcTestSuite) GetPrivKeyByIndex(algo hd.PubKeyType, index uint32) cryptotypes.PrivKey {
	privKey, err := helpers.PrivKeyFromMnemonic(suite.network.Config.Mnemonics[0], algo, 0, index)
	suite.NoError(err)
	return privKey
}

func (suite *rpcTestSuite) GetClients() []rpcTestClient {
	validator := suite.GetFirstValidator()
	suite.True(validator.AppConfig.GRPC.Enable)
	grpcClient, err := grpc.NewClient(fmt.Sprintf("http://%s", validator.AppConfig.GRPC.Address))
	suite.NoError(err)
	return []rpcTestClient{
		grpcClient,
		jsonrpc.NewNodeRPC(jsonrpc.NewFastClient(validator.RPCAddress)),
	}
}

func (suite *rpcTestSuite) FirstValidatorTransferTo(index uint32, amount sdk.Int) {
	validator := suite.GetFirstValidator()
	suite.True(validator.AppConfig.GRPC.Enable)
	grpcClient, err := grpc.NewClient(fmt.Sprintf("http://%s", validator.AppConfig.GRPC.Address))
	suite.NoError(err)
	valKey := suite.GetFirstValiPrivKey()
	nextValKey := suite.GetPrivKeyByIndex(hd.Secp256k1Type, index)
	txRaw, err := grpcClient.BuildTxV2(valKey,
		[]sdk.Msg{
			banktypes.NewMsgSend(
				valKey.PubKey().Address().Bytes(),
				nextValKey.PubKey().Address().Bytes(),
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amount)),
			),
		},
		500000,
		"",
		0,
	)
	suite.NoError(err)
	txResponse, err := grpcClient.BroadcastTx(txRaw)
	suite.NoError(err)
	suite.Equal(uint32(0), txResponse.Code)
}

func (suite *rpcTestSuite) TestClient_Tx() {
	privKey := suite.GetPrivKeyByIndex(hd.Secp256k1Type, 1)

	clients := suite.GetClients()
	for i := 0; i < len(clients); i++ {
		cli := clients[i]
		toAddress := sdk.AccAddress(helpers.NewPriKey().PubKey().Address())
		txRaw, err := cli.BuildTx(privKey, []sdk.Msg{
			banktypes.NewMsgSend(
				privKey.PubKey().Address().Bytes(),
				toAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1))),
			),
		},
		)
		suite.NoError(err)

		gas, err := cli.EstimatingGas(txRaw)
		suite.NoError(err)
		suite.True(gas.GasUsed < 90000)
		suite.Equal(uint64(0), gas.GasWanted)

		txResponse, err := cli.BroadcastTx(txRaw)
		suite.NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		err = suite.network.WaitForNextBlock()
		suite.NoError(err)

		txRes, err := cli.TxByHash(txResponse.TxHash)
		suite.NoError(err)
		txRes.Tx = nil
		txRes.Timestamp = ""
		suite.Equal(txResponse, txRes)

		account, err := cli.QueryAccount(toAddress.String())
		suite.NoError(err)
		suite.Equal(authtypes.NewBaseAccount(toAddress, nil, uint64(12+i), 0), account)
	}

	ethPrivKey := suite.GetPrivKeyByIndex(hd2.EthSecp256k1Type, 0)

	ethAddress := sdk.AccAddress(ethPrivKey.PubKey().Address().Bytes())

	for i := 0; i < len(clients); i++ {
		cli := clients[i]
		txRaw, err := cli.BuildTx(privKey, []sdk.Msg{
			banktypes.NewMsgSend(
				privKey.PubKey().Address().Bytes(),
				ethAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10).MulRaw(1e18))),
			),
		},
		)
		suite.NoError(err)

		gas, err := cli.EstimatingGas(txRaw)
		suite.NoError(err)
		suite.True(gas.GasUsed < 90000)
		suite.Equal(uint64(0), gas.GasWanted)

		txResponse, err := cli.BroadcastTx(txRaw)
		suite.NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		err = suite.network.WaitForNextBlock()
		suite.NoError(err)

		account, err := cli.QueryAccount(ethAddress.String())
		suite.NoError(err)
		suite.Equal(authtypes.NewBaseAccount(ethAddress, nil, 14, 0), account)
	}

	for i := 0; i < len(clients); i++ {
		cli := clients[i]
		toAddress := sdk.AccAddress(helpers.NewPriKey().PubKey().Address())
		txRaw, err := cli.BuildTx(ethPrivKey, []sdk.Msg{
			banktypes.NewMsgSend(
				ethPrivKey.PubKey().Address().Bytes(),
				toAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1))),
			),
		},
		)
		suite.NoError(err)

		gas, err := cli.EstimatingGas(txRaw)
		suite.NoError(err)
		suite.True(gas.GasUsed < 90000)
		suite.Equal(uint64(0), gas.GasWanted)

		txResponse, err := cli.BroadcastTx(txRaw)
		suite.NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		err = suite.network.WaitForNextBlock()
		suite.NoError(err)

		account, err := cli.QueryAccount(ethAddress.String())
		suite.NoError(err)
		baseAccount, ok := account.(*authtypes.BaseAccount)
		suite.True(ok)
		if baseAccount.PubKey.TypeUrl != "" {
			pubAny, err := types.NewAnyWithValue(ethPrivKey.PubKey())
			suite.NoError(err)
			suite.Equal("/"+proto.MessageName(&ethsecp256k1.PubKey{}), baseAccount.PubKey.TypeUrl)
			suite.Equal(pubAny, baseAccount.PubKey)
		}
		suite.Equal(uint64(i+1), account.GetSequence())
	}
}

func (suite *rpcTestSuite) TestClient_Query() {
	feeCollectorAddr, err := sdk.AccAddressFromHex("f1829676db577682e944fc3493d451b67ff3e29f")
	suite.NoError(err)
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
					suite.NoError(err)
					suite.True(height >= int64(1))
				},
			},
		},
		{
			funcName: "QuerySupply",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(supply sdk.Coins, err error) {
					suite.NoError(err)
					supply.IsAllGTE(
						sdk.Coins{
							sdk.Coin{
								Denom:  fxtypes.DefaultDenom,
								Amount: sdk.NewInt(500_000).MulRaw(1e18),
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
						Amount: sdk.NewInt(4).MulRaw(1e12),
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
					suite.GetFirstValidator().Address,
					suite.GetFirstValidator().PubKey,
					0,
					0,
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
					Amount: sdk.NewInt(488998).MulRaw(1e18),
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
						Amount: sdk.NewInt(488998).MulRaw(1e18),
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
					suite.Equal(wantResTf.Kind(), reflect.Func)
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

	nodeRPC := jsonrpc.NewNodeRPC(jsonrpc.NewFastClient(validator.RPCAddress))
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
					suite.NoError(err1)
					suite.NoError(err2)
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
					suite.NoError(err1)
					suite.NoError(err2)
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
					suite.NoError(err1)
					suite.NoError(err2)
					suite.EqualValues(len(res1.Peers), len(res2.Peers))
				},
			},
		},
		{
			funcName: "ConsensusState",
			params:   []interface{}{},
			wantRes: []interface{}{
				func(_ *ctypes.ResultConsensusState, err1 error, _ *ctypes.ResultConsensusState, err2 error) {
					suite.NoError(err1)
					suite.NoError(err2)
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
				suite.Equal(wantResTf.Kind(), reflect.Func)
				wantResTf.Call(append(result1, result2...))
			} else {
				for i := 0; i < len(result1); i++ {
					data1, err1 := json.Marshal(reflect.Indirect(result1[i]).Interface())
					suite.NoError(err1)
					data2, err2 := json.Marshal(reflect.Indirect(result2[i]).Interface())
					suite.NoError(err2)
					suite.JSONEq(string(data1), string(data2))
				}
			}
			close(resultChan)
		})
	}
}

func (suite *rpcTestSuite) TestJsonRPC_ABCI_Query() {
	// GetStakeValidators
	validator := suite.GetFirstValidator()
	nodeRPC := jsonrpc.NewNodeRPC(jsonrpc.NewFastClient(validator.RPCAddress))
	validators, err := nodeRPC.GetStakeValidators(stakingtypes.Bonded)
	suite.Require().NoError(err)
	suite.Require().Len(validators, 1)

	// QueryBalanceByHeight
	nextValKey := suite.GetPrivKeyByIndex(hd.Secp256k1Type, 1)
	nodeRPC.WithHeight(0)
	balances, err := nodeRPC.QueryBalances(sdk.AccAddress(nextValKey.PubKey().Address().Bytes()).String())
	suite.NoError(err)
	suite.True(balances.IsAllPositive())

	nodeRPC.WithHeight(1)
	balances, err = nodeRPC.QueryBalances(sdk.AccAddress(nextValKey.PubKey().Address().Bytes()).String())
	suite.NoError(err)
	suite.False(balances.IsAllPositive())
}
