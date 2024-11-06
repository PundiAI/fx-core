package integration

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	fxevmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

func (suite *IntegrationTest) WFXTest() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	tokenAddr := suite.GetErc20TokenAddress(fxtypes.DefaultDenom)
	wfxTokenSuite := NewERC20TokenSuite(suite.EthSuite, tokenAddr, signer)

	wfxTokenSuite.Deposit(helpers.NewBigInt(20, 18))
	wfxTokenSuite.Withdraw(signer.Address(), helpers.NewBigInt(10, 18))
	wfxTokenSuite.Transfer(helpers.GenHexAddress(), helpers.NewBigInt(5, 18))

	newSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	approveAmount := helpers.NewBigInt(5, 18)
	wfxTokenSuite.Approve(newSigner.Address(), approveAmount)

	// send tx fee to signer
	suite.Send(newSigner.AccAddress(), suite.NewStakingCoin(2, 18))

	wfxTokenSuite.WithSigner(newSigner).TransferFrom(signer.Address(), helpers.GenHexAddress(), approveAmount)
}

func (suite *IntegrationTest) ERC20TokenTest() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	tokenAddr := suite.DeployERC20(signer, "TEST")
	erc20TokenSuite := NewERC20TokenSuite(suite.EthSuite, tokenAddr, signer)

	erc20TokenSuite.Mint(signer.Address(), helpers.NewBigInt(100, 18))

	erc20TokenSuite.Transfer(helpers.GenHexAddress(), helpers.NewBigInt(20, 18))

	newSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	approveAmount := helpers.NewBigInt(5, 18)
	erc20TokenSuite.Approve(newSigner.Address(), approveAmount)

	// send tx fee to signer
	suite.Send(newSigner.AccAddress(), suite.NewStakingCoin(100, 18))

	erc20TokenSuite.WithSigner(newSigner).TransferFrom(signer.Address(), helpers.GenHexAddress(), approveAmount)
}

func (suite *IntegrationTest) ERC721Test() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	tokenAddr := suite.DeployERC721(signer)
	erc721TokenSuite := NewERC721TokenSuite(suite.EthSuite, tokenAddr, signer)

	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	erc721TokenSuite.SafeMint(signer.Address())

	newSigner := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(newSigner.AccAddress(), suite.NewStakingCoin(100, 18))

	erc721TokenSuite.Approve(newSigner.Address(), big.NewInt(0))
	erc721TokenSuite.WithSigner(newSigner).SafeTransferFrom(signer.Address(), helpers.GenHexAddress(), big.NewInt(0))

	erc721TokenSuite.SafeMint(signer.Address())
	erc721TokenSuite.SafeTransferFrom(signer.Address(), helpers.GenHexAddress(), big.NewInt(1))

	erc721TokenSuite.SafeMint(signer.Address())
	erc721TokenSuite.SetApprovalForAll(newSigner.Address(), true)
	erc721TokenSuite.WithSigner(newSigner).SafeTransferFrom(signer.Address(), helpers.GenHexAddress(), big.NewInt(2))
}

func (suite *IntegrationTest) EVMWeb3Test() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

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
			params:   []interface{}{signer.Address(), nil},
			wantRes:  []interface{}{new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), nil},
		},
		{
			name:     "eth_getBalance",
			funcName: "BalanceAt",
			params:   []interface{}{signer.Address(), big.NewInt(1)},
			wantRes:  []interface{}{big.NewInt(0), nil},
		},
		{
			name:     "eth_getStorageAt",
			funcName: "StorageAt",
			params:   []interface{}{signer.Address(), common.Hash{}, nil},
			wantRes:  []interface{}{[32]byte{}, nil},
		},
		{
			name:     "eth_getCode",
			funcName: "CodeAt",
			params:   []interface{}{signer.Address(), nil},
			wantRes:  []interface{}{[]byte{}, nil},
		},
		{
			name:     "eth_getTransactionCount",
			funcName: "NonceAt",
			params:   []interface{}{signer.Address(), nil},
			wantRes:  []interface{}{uint64(0), nil},
		},
		{
			name:     "eth_getBalance pending",
			funcName: "PendingBalanceAt",
			params:   []interface{}{signer.Address()},
			wantRes:  []interface{}{new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), nil},
		},
		{
			name:     "eth_getStorageAt pending",
			funcName: "PendingStorageAt",
			params:   []interface{}{signer.Address(), common.Hash{}},
			wantRes:  []interface{}{[32]byte{}, nil},
		},
		{
			name:     "eth_getCode pending",
			funcName: "PendingCodeAt",
			params:   []interface{}{signer.Address()},
			wantRes:  []interface{}{[]byte{}, nil},
		},
		{
			name:     "eth_getTransactionCount pending",
			funcName: "PendingNonceAt",
			params:   []interface{}{signer.Address()},
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
	ethClient := suite.EthSuite.ethCli
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
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	tokenAddr := suite.DeployERC20(signer, "TEST")
	erc20TokenSuite := NewERC20TokenSuite(suite.EthSuite, tokenAddr, signer)

	evmModuleAddr := common.BytesToAddress(authtypes.NewModuleAddress(evmtypes.ModuleName))
	erc20TokenSuite.TransferOwnership(evmModuleAddr)

	args, err := helpers.PackERC20Mint(signer.Address(), helpers.NewBigInt(100, 18))
	suite.Require().NoError(err)

	response, proposalId := suite.BroadcastProposalTxV1(
		&fxevmtypes.MsgCallContract{
			Authority:       authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ContractAddress: tokenAddr.String(),
			Data:            common.Bytes2Hex(args),
		},
	)
	suite.Require().EqualValues(response.Code, 0)
	suite.Require().True(proposalId > 0)
}

func (suite *IntegrationTest) ERC20CodeTest() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	erc20 := contract.GetFIP20()
	erc20Addr, _ := suite.DeployContract(signer, erc20.Bin)
	code, err := suite.EthSuite.ethCli.CodeAt(suite.ctx, erc20Addr, nil)
	suite.Require().NoError(err)
	suite.Equal(erc20.Code, code, fmt.Sprintf("erc20 deployed code: %s", common.Bytes2Hex(code)))

	deployedCode := bytes.ReplaceAll(code, erc20.Address.Bytes(), common.Address{}.Bytes())
	suite.True(strings.HasSuffix(contract.FIP20UpgradableMetaData.Bin, common.Bytes2Hex(deployedCode)))
}

func (suite *IntegrationTest) WFXCodeTest() {
	signer := helpers.NewSigner(helpers.NewEthPrivKey())
	suite.Send(signer.AccAddress(), suite.NewStakingCoin(100, 18))

	wfx := contract.GetWFX()
	wfxAddr, _ := suite.DeployContract(signer, wfx.Bin)
	code, err := suite.EthSuite.ethCli.CodeAt(suite.ctx, wfxAddr, nil)
	suite.Require().NoError(err)
	suite.Equal(wfx.Code, code, fmt.Sprintf("wfx deployed code: %s", common.Bytes2Hex(code)))

	deployedCode := bytes.ReplaceAll(code, wfx.Address.Bytes(), common.Address{}.Bytes())
	suite.True(strings.HasSuffix(contract.WFXUpgradableMetaData.Bin, common.Bytes2Hex(deployedCode)))
}
