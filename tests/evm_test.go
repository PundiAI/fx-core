package tests

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

func (suite *EvmTestSuite) TestWeb3Query() {
	tests := []struct {
		name     string
		funcName string
		params   []interface{}
		wantRes  []interface{}
	}{
		{
			name:     "eth_chainId",
			funcName: "ChainID",
			params:   []interface{}{},
			wantRes:  []interface{}{big.NewInt(530), nil},
		},
		{
			name:     "eth_getBlockByNumber",
			funcName: "BlockByNumber",
			params:   []interface{}{big.NewInt(1)},
			wantRes:  []interface{}{nil, nil},
		},
		{
			name:     "eth_getBlockByNumber latest",
			funcName: "BlockByNumber",
			params:   []interface{}{nil},
			wantRes:  []interface{}{nil, nil},
		},
		{
			name:     "eth_blockNumber",
			funcName: "BlockNumber",
			params:   []interface{}{},
			wantRes:  []interface{}{nil, nil},
		},
		{
			name:     "eth_getBlockByNumber",
			funcName: "HeaderByNumber",
			params:   []interface{}{big.NewInt(1)},
			wantRes:  []interface{}{nil, nil},
		},
		{
			name:     "eth_getBlockByNumber latest",
			funcName: "HeaderByNumber",
			params:   []interface{}{nil},
			wantRes:  []interface{}{nil, nil},
		},
		{
			name:     "eth_syncing",
			funcName: "SyncProgress",
			params:   []interface{}{},
			wantRes:  []interface{}{nil, nil},
		},
		{
			name:     "net_version",
			funcName: "NetworkID",
			params:   []interface{}{},
			wantRes:  []interface{}{big.NewInt(530), nil},
		},
		{
			name:     "eth_getBalance latest",
			funcName: "BalanceAt",
			params:   []interface{}{suite.HexAddress(), nil},
			wantRes:  []interface{}{new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), nil},
		},
		{
			name:     "eth_getBalance",
			funcName: "BalanceAt",
			params:   []interface{}{suite.HexAddress(), big.NewInt(1)},
			wantRes:  []interface{}{big.NewInt(0), nil},
		},
		{
			name:     "eth_getStorageAt",
			funcName: "StorageAt",
			params:   []interface{}{suite.HexAddress(), common.Hash{}, nil},
			wantRes:  []interface{}{[32]byte{}, nil},
		},
		{
			name:     "eth_getCode",
			funcName: "CodeAt",
			params:   []interface{}{suite.HexAddress(), nil},
			wantRes:  []interface{}{[]byte{}, nil},
		},
		{
			name:     "eth_getTransactionCount",
			funcName: "NonceAt",
			params:   []interface{}{suite.HexAddress(), nil},
			wantRes:  []interface{}{uint64(0), nil},
		},
		{
			name:     "eth_getBalance pending",
			funcName: "PendingBalanceAt",
			params:   []interface{}{suite.HexAddress()},
			wantRes:  []interface{}{new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), nil},
		},
		{
			name:     "eth_getStorageAt pending",
			funcName: "PendingStorageAt",
			params:   []interface{}{suite.HexAddress(), common.Hash{}},
			wantRes:  []interface{}{[32]byte{}, nil},
		},
		{
			name:     "eth_getCode pending",
			funcName: "PendingCodeAt",
			params:   []interface{}{suite.HexAddress()},
			wantRes:  []interface{}{[]byte{}, nil},
		},
		{
			name:     "eth_getTransactionCount pending",
			funcName: "PendingNonceAt",
			params:   []interface{}{suite.HexAddress()},
			wantRes:  []interface{}{big.NewInt(0), nil},
		},
		{
			name:     "eth_getBlockTransactionCountByNumber pending",
			funcName: "PendingTransactionCount",
			params:   []interface{}{},
			wantRes:  []interface{}{uint64(0), nil},
		},
		{
			name:     "eth_gasPrice",
			funcName: "SuggestGasPrice",
			params:   []interface{}{},
			wantRes:  []interface{}{big.NewInt(562500000000), nil},
		},
		{
			name:     "eth_maxPriorityFeePerGas",
			funcName: "SuggestGasTipCap",
			params:   []interface{}{},
			wantRes:  []interface{}{big.NewInt(62500000000), nil},
		},
	}
	client := suite.EthClient()
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			typeOf := reflect.TypeOf(client)
			method, is := typeOf.MethodByName(tt.funcName)
			suite.True(is)
			params := make([]reflect.Value, len(tt.params)+2)
			params[0] = reflect.ValueOf(client)
			params[1] = reflect.ValueOf(suite.ctx)
			for i := 2; i < len(params); i++ {
				p := tt.params[i-2]
				if p != nil {
					params[i] = reflect.ValueOf(p)
				} else {
					params[i] = reflect.New(reflect.TypeOf(&big.Int{})).Elem()
				}
			}
			results := method.Func.Call(params)
			for i := 0; i < len(results); i++ {
				if i == 0 && tt.wantRes[i] == nil {
					suite.T().Log(results[i])
					continue
				}
				suite.EqualValues(
					fmt.Sprintf("%v", tt.wantRes[i]),
					fmt.Sprintf("%v", results[i]),
				)
			}
		})
	}
}
