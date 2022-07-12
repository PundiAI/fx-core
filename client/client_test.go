package client_test

import (
	"fmt"
	"reflect"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/app/helpers"
	"github.com/functionx/fx-core/client/grpc"
	"github.com/functionx/fx-core/client/jsonrpc"
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
	EstimatingGas(txBody *tx.TxBody, authInfo *tx.AuthInfo, sign []byte) (*sdk.GasInfo, error)
	BroadcastTx(txRaw *tx.TxRaw, mode ...tx.BroadcastMode) (*sdk.TxResponse, error)
	TxByHash(txHash string) (*sdk.TxResponse, error)
}

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
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
	grpcClient, err := grpc.NewGRPCClient(fmt.Sprintf("http://%s", validator.AppConfig.GRPC.Address))
	suite.NoError(err)
	return []TestClient{
		jsonrpc.NewJsonRPC(jsonrpc.NewFastClient(validator.RPCAddress)),
		grpcClient,
	}
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.T().Log("setting up integration test suite")

	suite.cfg = helpers.DefaultConfig()
	suite.cfg.NumValidators = 1

	suite.network = network.New(suite.T(), suite.cfg)

	_, err := suite.network.WaitForHeight(1)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	suite.T().Log("tearing down integration test suite")

	// This is important and must be called to ensure other tests can create
	// a network!
	suite.network.Cleanup()
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
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			clients := suite.GetClients()
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
