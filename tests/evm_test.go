package tests

import (
	"fmt"
	"math/big"
	"reflect"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/client"
	"github.com/functionx/fx-core/v3/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *IntegrationTest) WFXTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))
	tx, err := client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, nil, nil, types.GetWFX().Bin)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	logic := crypto.CreateAddress(common.BytesToAddress(suite.evm.privKey.PubKey().Address().Bytes()), tx.Nonce())
	proxy := suite.deployProxy(suite.evm.privKey, logic, []byte{})
	pack, err := types.GetWFX().ABI.Pack("initialize", "Wrapped Function X", "WFX", uint8(18), common.BytesToAddress([]byte(evmtypes.ModuleName)))
	suite.NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, &proxy, nil, pack)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	testKey := helpers.NewEthPrivKey()
	suite.Send(testKey.PubKey().Address().Bytes(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))

	suite.evm.WFXDeposit(testKey, proxy, new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+81), big.NewInt(1e18)))
	suite.evm.WFXWithdraw(testKey, proxy, helpers.GenerateAddress(), new(big.Int).Mul(big.NewInt(tmrand.Int63n(29)+1), big.NewInt(1e18)))
	suite.evm.TransferERC20(testKey, proxy, helpers.GenerateAddress(), new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18)))

	spenderKey := helpers.NewEthPrivKey()
	suite.Send(spenderKey.PubKey().Address().Bytes(), suite.NewCoin(sdk.NewInt(1000).MulRaw(1e18)))
	suite.evm.ApproveERC20(testKey, proxy, common.BytesToAddress(spenderKey.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(tmrand.Int63n(10)+20), big.NewInt(1e18)))
	suite.evm.TransferFromERC20(spenderKey, proxy, common.BytesToAddress(testKey.PubKey().Address().Bytes()), helpers.GenerateAddress(), new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18)))
	suite.evm.WFXDeposit(spenderKey, proxy, new(big.Int).Mul(big.NewInt(500), big.NewInt(1e18)))
	fxAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(200)+200), big.NewInt(1e18))
	totalAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(450)+200), big.NewInt(1e18))
	feeAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(30)+1), big.NewInt(1e18))
	suite.evm.WFXTransferCrossChain(spenderKey, proxy, helpers.GenerateAddress().Hex(), totalAmount, fxAmount, feeAmount, "eth")
}

func (suite *IntegrationTest) ERC20TokenTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))
	tx, err := client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, nil, nil, types.GetERC20().Bin)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	logic := crypto.CreateAddress(common.BytesToAddress(suite.evm.privKey.PubKey().Address().Bytes()), tx.Nonce())
	proxy := suite.deployProxy(suite.evm.privKey, logic, []byte{})
	pack, err := types.GetERC20().ABI.Pack("initialize", "Test ERC20", "ERC20", uint8(18), common.BytesToAddress([]byte(evmtypes.ModuleName)))
	suite.NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, &proxy, nil, pack)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)

	testKey := helpers.NewEthPrivKey()
	testAddress := common.BytesToAddress(testKey.PubKey().Address().Bytes())
	suite.Send(testKey.PubKey().Address().Bytes(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))

	suite.evm.MintERC20(suite.evm.privKey, proxy, testAddress, new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)))
	transferAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18))
	suite.evm.TransferERC20(testKey, proxy, helpers.GenerateAddress(), transferAmount)

	spenderKey := helpers.NewEthPrivKey()
	spenderAddress := common.BytesToAddress(spenderKey.PubKey().Address().Bytes())
	approveAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(10)+20), big.NewInt(1e18))
	transferFromAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18))
	suite.Send(spenderKey.PubKey().Address().Bytes(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))
	suite.evm.ApproveERC20(testKey, proxy, spenderAddress, approveAmount)
	suite.evm.TransferFromERC20(spenderKey, proxy, testAddress, helpers.GenerateAddress(), transferFromAmount)
}

func (suite *IntegrationTest) ERC721Test() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))
	tx, err := client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, nil, nil, GetERC721().Bin)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	logic := crypto.CreateAddress(common.BytesToAddress(suite.evm.privKey.PubKey().Address().Bytes()), tx.Nonce())
	proxy := suite.deployProxy(suite.evm.privKey, logic, []byte{})
	pack, err := GetERC721().ABI.Pack("initialize")
	suite.NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, &proxy, nil, pack)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)

	testKey := helpers.NewEthPrivKey()
	testAddress := common.BytesToAddress(testKey.PubKey().Address().Bytes())
	suite.Send(testKey.PubKey().Address().Bytes(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))

	suite.evm.SafeMintERC721(suite.evm.privKey, proxy, testAddress)

	approveKey := helpers.NewEthPrivKey()
	approveAddress := common.BytesToAddress(approveKey.PubKey().Address().Bytes())
	suite.Send(approveKey.PubKey().Address().Bytes(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))
	suite.evm.ApproveERC721(testKey, proxy, approveAddress, big.NewInt(0))
	suite.evm.SafeTransferFrom(approveKey, proxy, testAddress, helpers.GenerateAddress(), big.NewInt(0))

	suite.evm.SafeMintERC721(suite.evm.privKey, proxy, testAddress)
	suite.evm.SafeTransferFrom(testKey, proxy, testAddress, helpers.GenerateAddress(), big.NewInt(1))

	suite.evm.SafeMintERC721(suite.evm.privKey, proxy, testAddress)
	suite.evm.SetApprovalForAll(testKey, proxy, approveAddress, true)
	suite.evm.SafeTransferFrom(testKey, proxy, testAddress, helpers.GenerateAddress(), big.NewInt(2))
}

func (suite *IntegrationTest) deployProxy(privateKey cryptotypes.PrivKey, logic common.Address, initData []byte) common.Address {
	input, err := types.GetERC1967Proxy().ABI.Pack("", logic, initData)
	suite.NoError(err)
	tx, err := client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), privateKey, nil, nil, append(types.GetERC1967Proxy().Bin, input...))
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	return crypto.CreateAddress(common.BytesToAddress(privateKey.PubKey().Address().Bytes()), tx.Nonce())
}

func (suite *IntegrationTest) EVMWeb3Test() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdk.NewInt(100).MulRaw(1e18)))

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
			params:   []interface{}{suite.evm.HexAddress(), nil},
			wantRes:  []interface{}{new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), nil},
		},
		{
			name:     "eth_getBalance",
			funcName: "BalanceAt",
			params:   []interface{}{suite.evm.HexAddress(), big.NewInt(1)},
			wantRes:  []interface{}{big.NewInt(0), nil},
		},
		{
			name:     "eth_getStorageAt",
			funcName: "StorageAt",
			params:   []interface{}{suite.evm.HexAddress(), common.Hash{}, nil},
			wantRes:  []interface{}{[32]byte{}, nil},
		},
		{
			name:     "eth_getCode",
			funcName: "CodeAt",
			params:   []interface{}{suite.evm.HexAddress(), nil},
			wantRes:  []interface{}{[]byte{}, nil},
		},
		{
			name:     "eth_getTransactionCount",
			funcName: "NonceAt",
			params:   []interface{}{suite.evm.HexAddress(), nil},
			wantRes:  []interface{}{uint64(0), nil},
		},
		{
			name:     "eth_getBalance pending",
			funcName: "PendingBalanceAt",
			params:   []interface{}{suite.evm.HexAddress()},
			wantRes:  []interface{}{new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), nil},
		},
		{
			name:     "eth_getStorageAt pending",
			funcName: "PendingStorageAt",
			params:   []interface{}{suite.evm.HexAddress(), common.Hash{}},
			wantRes:  []interface{}{[32]byte{}, nil},
		},
		{
			name:     "eth_getCode pending",
			funcName: "PendingCodeAt",
			params:   []interface{}{suite.evm.HexAddress()},
			wantRes:  []interface{}{[]byte{}, nil},
		},
		{
			name:     "eth_getTransactionCount pending",
			funcName: "PendingNonceAt",
			params:   []interface{}{suite.evm.HexAddress()},
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
	ethClient := suite.evm.EthClient()
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			typeOf := reflect.TypeOf(ethClient)
			method, is := typeOf.MethodByName(tt.funcName)
			suite.True(is)
			params := make([]reflect.Value, len(tt.params)+2)
			params[0] = reflect.ValueOf(ethClient)
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
					// marshal, _ := json.Marshal(results[i].Interface())
					// suite.T().Log(i, tt.name, string(marshal))
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
