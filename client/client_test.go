package client_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	"github.com/gogo/protobuf/proto"

	"github.com/functionx/fx-core/client/jsonrpc"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app/helpers"
	"github.com/functionx/fx-core/client/grpc"
	fxtypes "github.com/functionx/fx-core/types"
)

type TestClient interface {
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

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) GetFirstValidator() *network.Validator {
	return suite.network.Validators[0]
}

func (suite *IntegrationTestSuite) GetClients() []TestClient {
	validator := suite.GetFirstValidator()
	suite.True(validator.AppConfig.GRPC.Enable)
	grpcClient, err := grpc.NewClient(fmt.Sprintf("http://%s", validator.AppConfig.GRPC.Address))
	suite.NoError(err)
	return []TestClient{
		jsonrpc.NewNodeRPC(jsonrpc.NewFastClient(validator.RPCAddress)),
		grpcClient,
	}
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.T().Log("setting up integration test suite")

	cfg := helpers.DefaultConfig()
	cfg.NumValidators = 1
	cfg.Mnemonics = append(cfg.Mnemonics, helpers.NewMnemonic())
	suite.network = network.New(suite.T(), cfg)

	_, err := suite.network.WaitForHeight(1)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
}

func (suite *IntegrationTestSuite) TestClient_Tx() {
	cfg := suite.network.Config
	privKey, err := helpers.PrivKeyFromMnemonic(cfg.Mnemonics[0], hd.Secp256k1Type, 0, 0)
	suite.NoError(err)

	clients := suite.GetClients()
	for i := 0; i < len(clients); i++ {
		client := clients[i]
		toAddress := sdk.AccAddress(helpers.NewPriKey().PubKey().Address())
		txRaw, err := client.BuildTx(privKey, []sdk.Msg{
			banktypes.NewMsgSend(
				privKey.PubKey().Address().Bytes(),
				toAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1))),
			)},
		)
		suite.NoError(err)

		gas, err := client.EstimatingGas(txRaw)
		suite.NoError(err)
		suite.Equal(uint64(76053), gas.GasUsed)
		suite.Equal(uint64(0), gas.GasWanted)

		txResponse, err := client.BroadcastTx(txRaw)
		suite.NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		err = suite.network.WaitForNextBlock()
		suite.NoError(err)

		txRes, err := client.TxByHash(txResponse.TxHash)
		suite.NoError(err)
		txRes.Tx = nil
		txRes.Timestamp = ""
		suite.Equal(txResponse, txRes)

		account, err := client.QueryAccount(toAddress.String())
		suite.NoError(err)
		suite.Equal(authtypes.NewBaseAccount(toAddress, nil, uint64(11+i), 0), account)
	}

	ethPrivKey, err := helpers.PrivKeyFromMnemonic(cfg.Mnemonics[0], hd2.EthSecp256k1Type, 0, 0)
	suite.NoError(err)

	ethAddress := sdk.AccAddress(ethPrivKey.PubKey().Address().Bytes())

	for i := 0; i < len(clients); i++ {
		client := clients[i]
		txRaw, err := client.BuildTx(privKey, []sdk.Msg{
			banktypes.NewMsgSend(
				privKey.PubKey().Address().Bytes(),
				ethAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10).MulRaw(1e18))),
			)},
		)
		suite.NoError(err)

		gas, err := client.EstimatingGas(txRaw)
		suite.NoError(err)
		suite.True(gas.GasUsed == uint64(76823) || gas.GasUsed == uint64(68148))
		suite.Equal(uint64(0), gas.GasWanted)

		txResponse, err := client.BroadcastTx(txRaw)
		suite.NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		err = suite.network.WaitForNextBlock()
		suite.NoError(err)

		account, err := client.QueryAccount(ethAddress.String())
		suite.NoError(err)
		suite.Equal(authtypes.NewBaseAccount(ethAddress, nil, 13, 0), account)
	}

	for i := 0; i < len(clients); i++ {
		client := clients[i]
		toAddress := sdk.AccAddress(helpers.NewPriKey().PubKey().Address())
		txRaw, err := client.BuildTx(ethPrivKey, []sdk.Msg{
			banktypes.NewMsgSend(
				ethPrivKey.PubKey().Address().Bytes(),
				toAddress,
				sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1))),
			)},
		)
		suite.NoError(err)

		gas, err := client.EstimatingGas(txRaw)
		suite.NoError(err)
		suite.True(gas.GasUsed == uint64(76465) || gas.GasUsed == uint64(83152))
		suite.Equal(uint64(0), gas.GasWanted)

		txResponse, err := client.BroadcastTx(txRaw)
		suite.NoError(err)
		suite.Equal(uint32(0), txResponse.Code)

		err = suite.network.WaitForNextBlock()
		suite.NoError(err)

		account, err := client.QueryAccount(ethAddress.String())
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

func (suite *IntegrationTestSuite) TestClient_Query() {
	tests := []struct {
		name     string
		funcName string
		params   []interface{}
		wantRes  []interface{}
	}{
		{
			name:     "get chain id",
			funcName: "GetChainId",
			params:   []interface{}{},
			wantRes:  []interface{}{fxtypes.ChainID, nil},
		},
		{
			name:     "get block height",
			funcName: "GetBlockHeight",
			params:   []interface{}{},
			wantRes:  []interface{}{int64(7), nil},
		},
		{
			name:     "get mint denom",
			funcName: "GetMintDenom",
			params:   []interface{}{},
			wantRes:  []interface{}{fxtypes.DefaultDenom, nil},
		},
		{
			name:     "get address prefix",
			funcName: "GetAddressPrefix",
			params:   []interface{}{},
			wantRes:  []interface{}{fxtypes.AddressPrefix, nil},
		},
		{
			name:     "app version",
			funcName: "AppVersion",
			params:   []interface{}{},
			wantRes:  []interface{}{"", nil},
		},
		{
			name:     "get gas price",
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
			name:     "query account",
			funcName: "QueryAccount",
			params:   []interface{}{suite.GetFirstValidator().Address.String()},
			wantRes: []interface{}{authtypes.NewBaseAccount(
				suite.GetFirstValidator().Address,
				suite.GetFirstValidator().PubKey,
				0,
				0,
			),
				nil},
		},
		{
			name:     "query balance",
			funcName: "QueryBalance",
			params:   []interface{}{suite.GetFirstValidator().Address.String(), fxtypes.DefaultDenom},
			wantRes: []interface{}{
				sdk.Coin{
					Denom:  fxtypes.DefaultDenom,
					Amount: sdk.NewInt(40_000).MulRaw(1e18),
				},
				nil,
			},
		},
		{
			name:     "query balances",
			funcName: "QueryBalances",
			params:   []interface{}{suite.GetFirstValidator().Address.String()},
			wantRes: []interface{}{
				sdk.Coins{
					sdk.Coin{
						Denom:  fxtypes.DefaultDenom,
						Amount: sdk.NewInt(40_000).MulRaw(1e18),
					},
					sdk.Coin{
						Denom:  "node0token",
						Amount: sdk.NewInt(100_000).MulRaw(1e18),
					},
				},
				nil,
			},
		},
		{
			name:     "query supply",
			funcName: "QuerySupply",
			params:   []interface{}{},
			wantRes: []interface{}{
				sdk.Coins{
					sdk.Coin{
						Denom:  fxtypes.DefaultDenom,
						Amount: sdk.MustNewDecFromStr("50000019408963423760327").RoundInt(),
					},
					sdk.Coin{
						Denom:  "node0token",
						Amount: sdk.NewInt(100_000).MulRaw(1e18),
					},
				},
				nil,
			},
		},
	}
	clients := suite.GetClients()
	for _, tt := range tests {
		suite.Run(tt.name, func() {
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
				for i := 0; i < len(results); i++ {
					suite.EqualValues(
						fmt.Sprintf("%v", tt.wantRes[i]),
						fmt.Sprintf("%v", results[i]),
					)
				}
			}
		})
	}
}
