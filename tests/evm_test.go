package tests

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/client"
	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/types"
	fxevmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

func (suite *IntegrationTest) WFXTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
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
	key := helpers.NewEthPrivKey()
	suite.Send(key.PubKey().Address().Bytes(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	suite.evm.WFXDeposit(key, proxy, new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+81), big.NewInt(1e18)))
	suite.evm.WFXWithdraw(key, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(tmrand.Int63n(29)+1), big.NewInt(1e18)))
	suite.evm.TransferERC20(key, proxy, helpers.GenerateAddress(), new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18)))
	spenderKey := helpers.NewEthPrivKey()
	suite.Send(spenderKey.PubKey().Address().Bytes(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	suite.evm.ApproveERC20(key, proxy, common.BytesToAddress(spenderKey.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(tmrand.Int63n(10)+20), big.NewInt(1e18)))
	suite.evm.TransferFromERC20(spenderKey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), helpers.GenerateAddress(), new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18)))
}

func (suite *IntegrationTest) ERC20TokenTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	tx, err := client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, nil, nil, types.GetFIP20().Bin)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	logic := crypto.CreateAddress(common.BytesToAddress(suite.evm.privKey.PubKey().Address().Bytes()), tx.Nonce())
	proxy := suite.deployProxy(suite.evm.privKey, logic, []byte{})
	pack, err := types.GetFIP20().ABI.Pack("initialize", "Test ERC20", "ERC20", uint8(18), common.BytesToAddress([]byte(evmtypes.ModuleName)))
	suite.NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, &proxy, nil, pack)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)

	key := helpers.NewEthPrivKey()
	suite.Send(key.PubKey().Address().Bytes(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	suite.evm.MintERC20(suite.evm.privKey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)))
	suite.evm.CheckBalanceOf(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)))

	transferAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18))
	recipient := helpers.GenerateAddress()

	suite.evm.TransferERC20(key, proxy, recipient, transferAmount)
	suite.evm.CheckBalanceOf(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), new(big.Int).Sub(new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), transferAmount))
	suite.evm.CheckBalanceOf(proxy, recipient, transferAmount)

	spenderKey := helpers.NewEthPrivKey()
	approveAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(10)+20), big.NewInt(1e18))
	transferFromAmount := new(big.Int).Mul(big.NewInt(tmrand.Int63n(19)+1), big.NewInt(1e18))
	suite.Send(spenderKey.PubKey().Address().Bytes(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	suite.evm.ApproveERC20(key, proxy, common.BytesToAddress(spenderKey.PubKey().Address().Bytes()), approveAmount)
	suite.evm.CheckAllowance(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), common.BytesToAddress(spenderKey.PubKey().Address().Bytes()), approveAmount)
	suite.evm.TransferFromERC20(spenderKey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), helpers.GenerateAddress(), transferFromAmount)
	suite.evm.CheckAllowance(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), common.BytesToAddress(spenderKey.PubKey().Address().Bytes()), new(big.Int).Sub(approveAmount, transferFromAmount))
}

func (suite *IntegrationTest) ERC721Test() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
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

	key := helpers.NewEthPrivKey()
	suite.Send(key.PubKey().Address().Bytes(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	suite.evm.SafeMintERC721(suite.evm.privKey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()))
	suite.evm.CheckBalanceOfERC721(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), big.NewInt(1))
	approvekey := helpers.NewEthPrivKey()
	suite.Send(approvekey.PubKey().Address().Bytes(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))

	suite.evm.ApproveERC721(key, proxy, common.BytesToAddress(approvekey.PubKey().Address().Bytes()), big.NewInt(0))
	suite.evm.SafeTransferFrom(approvekey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), helpers.GenerateAddress(), big.NewInt(0))
	suite.evm.CheckBalanceOfERC721(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), big.NewInt(0))

	suite.evm.SafeMintERC721(suite.evm.privKey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()))
	suite.evm.CheckBalanceOfERC721(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), big.NewInt(1))
	suite.evm.SafeTransferFrom(key, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), helpers.GenerateAddress(), big.NewInt(1))
	suite.evm.CheckBalanceOfERC721(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), big.NewInt(0))

	suite.evm.SafeMintERC721(suite.evm.privKey, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()))
	suite.evm.CheckBalanceOfERC721(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), big.NewInt(1))
	suite.evm.SetApprovalForAll(key, proxy, common.BytesToAddress(approvekey.PubKey().Address().Bytes()), true)
	suite.True(suite.evm.IsApprovedForAll(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), common.BytesToAddress(approvekey.PubKey().Address().Bytes())))

	suite.evm.SafeTransferFrom(key, proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), helpers.GenerateAddress(), big.NewInt(2))
	suite.evm.CheckBalanceOfERC721(proxy, common.BytesToAddress(key.PubKey().Address().Bytes()), big.NewInt(0))
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
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))

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

func (suite *IntegrationTest) CallContractTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	tx, err := client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, nil, nil, types.GetFIP20().Bin)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	logic := crypto.CreateAddress(common.BytesToAddress(suite.evm.privKey.PubKey().Address().Bytes()), tx.Nonce())
	proxy := suite.deployProxy(suite.evm.privKey, logic, []byte{})
	pack, err := types.GetFIP20().ABI.Pack("initialize", "Test ERC20", "ERC20", uint8(18), common.BytesToAddress(suite.evm.privKey.PubKey().Address().Bytes()))
	suite.NoError(err)
	tx, err = client.BuildEthTransaction(suite.ctx, suite.evm.EthClient(), suite.evm.privKey, &proxy, nil, pack)
	suite.NoError(err)
	suite.evm.SendTransaction(tx)
	suite.evm.TransferOwnership(suite.evm.privKey, proxy, common.BytesToAddress(authtypes.NewModuleAddress(evmtypes.ModuleName)))
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
	args, err := types.GetFIP20().ABI.Pack("mint", suite.evm.HexAddress(), amount)
	suite.Require().NoError(err)
	response, proposalId := suite.BroadcastProposalTx2([]sdk.Msg{&fxevmtypes.MsgCallContract{
		Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ContractAddress: proxy.String(),
		Data:            common.Bytes2Hex(args),
	}}, "UpdateContractProposal", "UpdateContractProposal")
	//	suite.Require().EqualValues(amount, suite.evm.BalanceOf(proxy, suite.evm.HexAddress()))
	suite.Require().EqualValues(response.Code, 0)
	suite.Require().True(proposalId > 0)
}

func (suite *IntegrationTest) FIP20CodeCheckTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	fip20Addr, _ := suite.evm.DeployContract(suite.evm.privKey, types.GetFIP20().Bin)
	code, err := suite.evm.EthClient().CodeAt(suite.ctx, fip20Addr, nil)
	suite.Require().NoError(err)
	suite.Equal(types.GetFIP20().Code, code, fmt.Sprintf("fip20 deployed code: %s", common.Bytes2Hex(code)))

	deployedCode := bytes.ReplaceAll(code, types.GetFIP20().Address.Bytes(), common.HexToAddress(types.EmptyEvmAddress).Bytes())
	suite.True(strings.HasSuffix(contract.FIP20UpgradableMetaData.Bin, common.Bytes2Hex(deployedCode)))
}

func (suite *IntegrationTest) WFXCodeCheckTest() {
	suite.Send(suite.evm.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
	wfxAddr, _ := suite.evm.DeployContract(suite.evm.privKey, types.GetWFX().Bin)
	code, err := suite.evm.EthClient().CodeAt(suite.ctx, wfxAddr, nil)
	suite.Require().NoError(err)
	suite.Equal(types.GetWFX().Code, code, fmt.Sprintf("wfx deployed code: %s", common.Bytes2Hex(code)))

	deployedCode := bytes.ReplaceAll(code, types.GetWFX().Address.Bytes(), common.HexToAddress(types.EmptyEvmAddress).Bytes())
	suite.True(strings.HasSuffix(contract.WFXUpgradableMetaData.Bin, common.Bytes2Hex(deployedCode)))
}
