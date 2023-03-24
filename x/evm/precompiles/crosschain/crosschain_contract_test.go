// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package crosschain_test

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CrosschainTestMetaData contains all meta data concerning the CrosschainTest contract.
var CrosschainTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txid\",\"type\":\"uint256\"}],\"name\":\"CancelSendToExternal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receipt\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"CrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"IncreaseBridgeFee\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txid\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"increaseBridgeFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611361806100206000396000f3fe6080604052600436106100345760003560e01c80630b56c19014610039578063160d7c7314610076578063c79a6b7b146100a6575b600080fd5b34801561004557600080fd5b50610060600480360381019061005b9190610ad7565b6100d6565b60405161006d9190610b4e565b60405180910390f35b610090600480360381019061008b9190610bfd565b6100ea565b60405161009d9190610b4e565b60405180910390f35b6100c060048036038101906100bb9190610cc2565b610252565b6040516100cd9190610b4e565b60405180910390f35b60006100e2838361026a565b905092915050565b60008073ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff1614610238578673ffffffffffffffffffffffffffffffffffffffff166323b872dd333087896101499190610d74565b6040518463ffffffff1660e01b815260040161016793929190610dc6565b6020604051808303816000875af1158015610186573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101aa9190610e29565b508673ffffffffffffffffffffffffffffffffffffffff1663095ea7b361100486886101d69190610d74565b6040518363ffffffff1660e01b81526004016101f3929190610e56565b6020604051808303816000875af1158015610212573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102369190610e29565b505b610246878787878787610339565b90509695505050505050565b600061026085858585610567565b9050949350505050565b600080600061100473ffffffffffffffffffffffffffffffffffffffff16610292868661063a565b60405161029f9190610ef0565b6000604051808303816000865af19150503d80600081146102dc576040519150601f19603f3d011682016040523d82523d6000602084013e6102e1565b606091505b509150915061032682826040518060400160405280601e81526020017f63616e63656c2073656e6420746f2065787465726e616c206661696c656400008152506106d4565b61032f8161079b565b9250505092915050565b60008073ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff16146104435760008773ffffffffffffffffffffffffffffffffffffffff1663dd62ed3e306110046040518363ffffffff1660e01b81526004016103ad929190610f07565b602060405180830381865afa1580156103ca573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ee9190610f45565b905084866103fc9190610d74565b811461043d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161043490610fcf565b60405180910390fd5b50610491565b838561044f9190610d74565b3414610490576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104879061103b565b60405180910390fd5b5b60008061100473ffffffffffffffffffffffffffffffffffffffff16346104bc8b8b8b8b8b8b6107bd565b6040516104c99190610ef0565b60006040518083038185875af1925050503d8060008114610506576040519150601f19603f3d011682016040523d82523d6000602084013e61050b565b606091505b509150915061055082826040518060400160405280601281526020017f63726f73732d636861696e206661696c656400000000000000000000000000008152506106d4565b61055981610863565b925050509695505050505050565b600080600061100473ffffffffffffffffffffffffffffffffffffffff1661059188888888610885565b60405161059e9190610ef0565b6000604051808303816000865af19150503d80600081146105db576040519150601f19603f3d011682016040523d82523d6000602084013e6105e0565b606091505b509150915061062582826040518060400160405280601a81526020017f696e6372656173652062726964676520666565206661696c65640000000000008152506106d4565b61062e81610925565b92505050949350505050565b6060828260405160240161064f92919061109f565b6040516020818303038152906040527feeb3593d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b82610796576000828060200190518101906106ef919061113f565b905060018251101561073857806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161072f9190611188565b60405180910390fd5b818160405160200161074b929190611232565b6040516020818303038152906040526040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078d9190611188565b60405180910390fd5b505050565b600080828060200190518101906107b29190610e29565b905080915050919050565b60608686868686866040516024016107da96959493929190611270565b6040516020818303038152906040527f160d7c73000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509695505050505050565b6000808280602001905181019061087a9190610e29565b905080915050919050565b60608484848460405160240161089e94939291906112df565b6040516020818303038152906040527f9b45009d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050949350505050565b6000808280602001905181019061093c9190610e29565b905080915050919050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6109ae82610965565b810181811067ffffffffffffffff821117156109cd576109cc610976565b5b80604052505050565b60006109e0610947565b90506109ec82826109a5565b919050565b600067ffffffffffffffff821115610a0c57610a0b610976565b5b610a1582610965565b9050602081019050919050565b82818337600083830152505050565b6000610a44610a3f846109f1565b6109d6565b905082815260208101848484011115610a6057610a5f610960565b5b610a6b848285610a22565b509392505050565b600082601f830112610a8857610a8761095b565b5b8135610a98848260208601610a31565b91505092915050565b6000819050919050565b610ab481610aa1565b8114610abf57600080fd5b50565b600081359050610ad181610aab565b92915050565b60008060408385031215610aee57610aed610951565b5b600083013567ffffffffffffffff811115610b0c57610b0b610956565b5b610b1885828601610a73565b9250506020610b2985828601610ac2565b9150509250929050565b60008115159050919050565b610b4881610b33565b82525050565b6000602082019050610b636000830184610b3f565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b9482610b69565b9050919050565b610ba481610b89565b8114610baf57600080fd5b50565b600081359050610bc181610b9b565b92915050565b6000819050919050565b610bda81610bc7565b8114610be557600080fd5b50565b600081359050610bf781610bd1565b92915050565b60008060008060008060c08789031215610c1a57610c19610951565b5b6000610c2889828a01610bb2565b965050602087013567ffffffffffffffff811115610c4957610c48610956565b5b610c5589828a01610a73565b9550506040610c6689828a01610ac2565b9450506060610c7789828a01610ac2565b9350506080610c8889828a01610be8565b92505060a087013567ffffffffffffffff811115610ca957610ca8610956565b5b610cb589828a01610a73565b9150509295509295509295565b60008060008060808587031215610cdc57610cdb610951565b5b600085013567ffffffffffffffff811115610cfa57610cf9610956565b5b610d0687828801610a73565b9450506020610d1787828801610ac2565b9350506040610d2887828801610bb2565b9250506060610d3987828801610ac2565b91505092959194509250565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610d7f82610aa1565b9150610d8a83610aa1565b9250828201905080821115610da257610da1610d45565b5b92915050565b610db181610b89565b82525050565b610dc081610aa1565b82525050565b6000606082019050610ddb6000830186610da8565b610de86020830185610da8565b610df56040830184610db7565b949350505050565b610e0681610b33565b8114610e1157600080fd5b50565b600081519050610e2381610dfd565b92915050565b600060208284031215610e3f57610e3e610951565b5b6000610e4d84828501610e14565b91505092915050565b6000604082019050610e6b6000830185610da8565b610e786020830184610db7565b9392505050565b600081519050919050565b600081905092915050565b60005b83811015610eb3578082015181840152602081019050610e98565b60008484015250505050565b6000610eca82610e7f565b610ed48185610e8a565b9350610ee4818560208601610e95565b80840191505092915050565b6000610efc8284610ebf565b915081905092915050565b6000604082019050610f1c6000830185610da8565b610f296020830184610da8565b9392505050565b600081519050610f3f81610aab565b92915050565b600060208284031215610f5b57610f5a610951565b5b6000610f6984828501610f30565b91505092915050565b600082825260208201905092915050565b7f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b20666565600082015250565b6000610fb9602083610f72565b9150610fc482610f83565b602082019050919050565b60006020820190508181036000830152610fe881610fac565b9050919050565b7f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b20666565600082015250565b6000611025602083610f72565b915061103082610fef565b602082019050919050565b6000602082019050818103600083015261105481611018565b9050919050565b600081519050919050565b60006110718261105b565b61107b8185610f72565b935061108b818560208601610e95565b61109481610965565b840191505092915050565b600060408201905081810360008301526110b98185611066565b90506110c86020830184610db7565b9392505050565b60006110e26110dd846109f1565b6109d6565b9050828152602081018484840111156110fe576110fd610960565b5b611109848285610e95565b509392505050565b600082601f8301126111265761112561095b565b5b81516111368482602086016110cf565b91505092915050565b60006020828403121561115557611154610951565b5b600082015167ffffffffffffffff81111561117357611172610956565b5b61117f84828501611111565b91505092915050565b600060208201905081810360008301526111a28184611066565b905092915050565b600081905092915050565b60006111c08261105b565b6111ca81856111aa565b93506111da818560208601610e95565b80840191505092915050565b7f3a20000000000000000000000000000000000000000000000000000000000000600082015250565b600061121c6002836111aa565b9150611227826111e6565b600282019050919050565b600061123e82856111b5565b91506112498261120f565b915061125582846111b5565b91508190509392505050565b61126a81610bc7565b82525050565b600060c0820190506112856000830189610da8565b81810360208301526112978188611066565b90506112a66040830187610db7565b6112b36060830186610db7565b6112c06080830185611261565b81810360a08301526112d28184611066565b9050979650505050505050565b600060808201905081810360008301526112f98187611066565b90506113086020830186610db7565b6113156040830185610da8565b6113226060830184610db7565b9594505050505056fea2646970667358221220c5dfe4e8f36bd5be12aa2073b254188b5ae15aa456cb41141385d0bbe7fe60a464736f6c63430008130033",
}

// CrosschainTestABI is the input ABI used to generate the binding from.
// Deprecated: Use CrosschainTestMetaData.ABI instead.
var CrosschainTestABI = CrosschainTestMetaData.ABI

// CrosschainTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrosschainTestMetaData.Bin instead.
var CrosschainTestBin = CrosschainTestMetaData.Bin

// DeployCrosschainTest deploys a new Ethereum contract, binding an instance of CrosschainTest to it.
func DeployCrosschainTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrosschainTest, error) {
	parsed, err := CrosschainTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrosschainTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrosschainTest{CrosschainTestCaller: CrosschainTestCaller{contract: contract}, CrosschainTestTransactor: CrosschainTestTransactor{contract: contract}, CrosschainTestFilterer: CrosschainTestFilterer{contract: contract}}, nil
}

// CrosschainTest is an auto generated Go binding around an Ethereum contract.
type CrosschainTest struct {
	CrosschainTestCaller     // Read-only binding to the contract
	CrosschainTestTransactor // Write-only binding to the contract
	CrosschainTestFilterer   // Log filterer for contract events
}

// CrosschainTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrosschainTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrosschainTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrosschainTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrosschainTestSession struct {
	Contract     *CrosschainTest   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrosschainTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrosschainTestCallerSession struct {
	Contract *CrosschainTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CrosschainTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrosschainTestTransactorSession struct {
	Contract     *CrosschainTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CrosschainTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrosschainTestRaw struct {
	Contract *CrosschainTest // Generic contract binding to access the raw methods on
}

// CrosschainTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrosschainTestCallerRaw struct {
	Contract *CrosschainTestCaller // Generic read-only contract binding to access the raw methods on
}

// CrosschainTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrosschainTestTransactorRaw struct {
	Contract *CrosschainTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrosschainTest creates a new instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTest(address common.Address, backend bind.ContractBackend) (*CrosschainTest, error) {
	contract, err := bindCrosschainTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrosschainTest{CrosschainTestCaller: CrosschainTestCaller{contract: contract}, CrosschainTestTransactor: CrosschainTestTransactor{contract: contract}, CrosschainTestFilterer: CrosschainTestFilterer{contract: contract}}, nil
}

// NewCrosschainTestCaller creates a new read-only instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTestCaller(address common.Address, caller bind.ContractCaller) (*CrosschainTestCaller, error) {
	contract, err := bindCrosschainTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestCaller{contract: contract}, nil
}

// NewCrosschainTestTransactor creates a new write-only instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTestTransactor(address common.Address, transactor bind.ContractTransactor) (*CrosschainTestTransactor, error) {
	contract, err := bindCrosschainTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestTransactor{contract: contract}, nil
}

// NewCrosschainTestFilterer creates a new log filterer instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTestFilterer(address common.Address, filterer bind.ContractFilterer) (*CrosschainTestFilterer, error) {
	contract, err := bindCrosschainTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestFilterer{contract: contract}, nil
}

// bindCrosschainTest binds a generic wrapper to an already deployed contract.
func bindCrosschainTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrosschainTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainTest *CrosschainTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainTest.Contract.CrosschainTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainTest *CrosschainTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrosschainTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainTest *CrosschainTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrosschainTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainTest *CrosschainTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainTest *CrosschainTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainTest *CrosschainTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainTest.Contract.contract.Transact(opts, method, params...)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txid) returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) CancelSendToExternal(opts *bind.TransactOpts, _chain string, _txid *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "cancelSendToExternal", _chain, _txid)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txid) returns(bool)
func (_CrosschainTest *CrosschainTestSession) CancelSendToExternal(_chain string, _txid *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CancelSendToExternal(&_CrosschainTest.TransactOpts, _chain, _txid)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txid) returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) CancelSendToExternal(_chain string, _txid *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CancelSendToExternal(&_CrosschainTest.TransactOpts, _chain, _txid)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) CrossChain(opts *bind.TransactOpts, _token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "crossChain", _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrosschainTest *CrosschainTestSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrossChain(&_CrosschainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrossChain(&_CrosschainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txid, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) IncreaseBridgeFee(opts *bind.TransactOpts, _chain string, _txid *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "increaseBridgeFee", _chain, _txid, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txid, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestSession) IncreaseBridgeFee(_chain string, _txid *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.IncreaseBridgeFee(&_CrosschainTest.TransactOpts, _chain, _txid, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txid, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) IncreaseBridgeFee(_chain string, _txid *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.IncreaseBridgeFee(&_CrosschainTest.TransactOpts, _chain, _txid, _token, _fee)
}

// CrosschainTestCancelSendToExternalIterator is returned from FilterCancelSendToExternal and is used to iterate over the raw logs and unpacked data for CancelSendToExternal events raised by the CrosschainTest contract.
type CrosschainTestCancelSendToExternalIterator struct {
	Event *CrosschainTestCancelSendToExternal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrosschainTestCancelSendToExternalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainTestCancelSendToExternal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrosschainTestCancelSendToExternal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrosschainTestCancelSendToExternalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainTestCancelSendToExternalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainTestCancelSendToExternal represents a CancelSendToExternal event raised by the CrosschainTest contract.
type CrosschainTestCancelSendToExternal struct {
	Sender common.Address
	Chain  string
	Txid   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCancelSendToExternal is a free log retrieval operation binding the contract event 0xe2ae965fb5b8e4c7da962424292951c18e0e9c1905b87c78cf0186fa70382535.
//
// Solidity: event CancelSendToExternal(address indexed sender, string chain, uint256 txid)
func (_CrosschainTest *CrosschainTestFilterer) FilterCancelSendToExternal(opts *bind.FilterOpts, sender []common.Address) (*CrosschainTestCancelSendToExternalIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CrosschainTest.contract.FilterLogs(opts, "CancelSendToExternal", senderRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestCancelSendToExternalIterator{contract: _CrosschainTest.contract, event: "CancelSendToExternal", logs: logs, sub: sub}, nil
}

// WatchCancelSendToExternal is a free log subscription operation binding the contract event 0xe2ae965fb5b8e4c7da962424292951c18e0e9c1905b87c78cf0186fa70382535.
//
// Solidity: event CancelSendToExternal(address indexed sender, string chain, uint256 txid)
func (_CrosschainTest *CrosschainTestFilterer) WatchCancelSendToExternal(opts *bind.WatchOpts, sink chan<- *CrosschainTestCancelSendToExternal, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CrosschainTest.contract.WatchLogs(opts, "CancelSendToExternal", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainTestCancelSendToExternal)
				if err := _CrosschainTest.contract.UnpackLog(event, "CancelSendToExternal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancelSendToExternal is a log parse operation binding the contract event 0xe2ae965fb5b8e4c7da962424292951c18e0e9c1905b87c78cf0186fa70382535.
//
// Solidity: event CancelSendToExternal(address indexed sender, string chain, uint256 txid)
func (_CrosschainTest *CrosschainTestFilterer) ParseCancelSendToExternal(log types.Log) (*CrosschainTestCancelSendToExternal, error) {
	event := new(CrosschainTestCancelSendToExternal)
	if err := _CrosschainTest.contract.UnpackLog(event, "CancelSendToExternal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrosschainTestCrossChainIterator is returned from FilterCrossChain and is used to iterate over the raw logs and unpacked data for CrossChain events raised by the CrosschainTest contract.
type CrosschainTestCrossChainIterator struct {
	Event *CrosschainTestCrossChain // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrosschainTestCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainTestCrossChain)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrosschainTestCrossChain)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrosschainTestCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainTestCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainTestCrossChain represents a CrossChain event raised by the CrosschainTest contract.
type CrosschainTestCrossChain struct {
	Sender  common.Address
	Token   common.Address
	Denom   string
	Receipt string
	Amount  *big.Int
	Fee     *big.Int
	Target  [32]byte
	Memo    string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCrossChain is a free log retrieval operation binding the contract event 0xb783df819ac99ca709650d67d9237a00b553c6ef941dceabeed6f4bc990d31ba.
//
// Solidity: event CrossChain(address indexed sender, address indexed token, string denom, string receipt, uint256 amount, uint256 fee, bytes32 target, string memo)
func (_CrosschainTest *CrosschainTestFilterer) FilterCrossChain(opts *bind.FilterOpts, sender []common.Address, token []common.Address) (*CrosschainTestCrossChainIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.FilterLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestCrossChainIterator{contract: _CrosschainTest.contract, event: "CrossChain", logs: logs, sub: sub}, nil
}

// WatchCrossChain is a free log subscription operation binding the contract event 0xb783df819ac99ca709650d67d9237a00b553c6ef941dceabeed6f4bc990d31ba.
//
// Solidity: event CrossChain(address indexed sender, address indexed token, string denom, string receipt, uint256 amount, uint256 fee, bytes32 target, string memo)
func (_CrosschainTest *CrosschainTestFilterer) WatchCrossChain(opts *bind.WatchOpts, sink chan<- *CrosschainTestCrossChain, sender []common.Address, token []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.WatchLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainTestCrossChain)
				if err := _CrosschainTest.contract.UnpackLog(event, "CrossChain", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCrossChain is a log parse operation binding the contract event 0xb783df819ac99ca709650d67d9237a00b553c6ef941dceabeed6f4bc990d31ba.
//
// Solidity: event CrossChain(address indexed sender, address indexed token, string denom, string receipt, uint256 amount, uint256 fee, bytes32 target, string memo)
func (_CrosschainTest *CrosschainTestFilterer) ParseCrossChain(log types.Log) (*CrosschainTestCrossChain, error) {
	event := new(CrosschainTestCrossChain)
	if err := _CrosschainTest.contract.UnpackLog(event, "CrossChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrosschainTestIncreaseBridgeFeeIterator is returned from FilterIncreaseBridgeFee and is used to iterate over the raw logs and unpacked data for IncreaseBridgeFee events raised by the CrosschainTest contract.
type CrosschainTestIncreaseBridgeFeeIterator struct {
	Event *CrosschainTestIncreaseBridgeFee // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrosschainTestIncreaseBridgeFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainTestIncreaseBridgeFee)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrosschainTestIncreaseBridgeFee)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrosschainTestIncreaseBridgeFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainTestIncreaseBridgeFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainTestIncreaseBridgeFee represents a IncreaseBridgeFee event raised by the CrosschainTest contract.
type CrosschainTestIncreaseBridgeFee struct {
	Sender common.Address
	Token  common.Address
	Chain  string
	Txid   *big.Int
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIncreaseBridgeFee is a free log retrieval operation binding the contract event 0x4b4d0e64eb77c0f61892107908295f09b3e381c50c655f4a73a4ad61c07350a0.
//
// Solidity: event IncreaseBridgeFee(address indexed sender, address indexed token, string chain, uint256 txid, uint256 fee)
func (_CrosschainTest *CrosschainTestFilterer) FilterIncreaseBridgeFee(opts *bind.FilterOpts, sender []common.Address, token []common.Address) (*CrosschainTestIncreaseBridgeFeeIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.FilterLogs(opts, "IncreaseBridgeFee", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestIncreaseBridgeFeeIterator{contract: _CrosschainTest.contract, event: "IncreaseBridgeFee", logs: logs, sub: sub}, nil
}

// WatchIncreaseBridgeFee is a free log subscription operation binding the contract event 0x4b4d0e64eb77c0f61892107908295f09b3e381c50c655f4a73a4ad61c07350a0.
//
// Solidity: event IncreaseBridgeFee(address indexed sender, address indexed token, string chain, uint256 txid, uint256 fee)
func (_CrosschainTest *CrosschainTestFilterer) WatchIncreaseBridgeFee(opts *bind.WatchOpts, sink chan<- *CrosschainTestIncreaseBridgeFee, sender []common.Address, token []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.WatchLogs(opts, "IncreaseBridgeFee", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainTestIncreaseBridgeFee)
				if err := _CrosschainTest.contract.UnpackLog(event, "IncreaseBridgeFee", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIncreaseBridgeFee is a log parse operation binding the contract event 0x4b4d0e64eb77c0f61892107908295f09b3e381c50c655f4a73a4ad61c07350a0.
//
// Solidity: event IncreaseBridgeFee(address indexed sender, address indexed token, string chain, uint256 txid, uint256 fee)
func (_CrosschainTest *CrosschainTestFilterer) ParseIncreaseBridgeFee(log types.Log) (*CrosschainTestIncreaseBridgeFee, error) {
	event := new(CrosschainTestIncreaseBridgeFee)
	if err := _CrosschainTest.contract.UnpackLog(event, "IncreaseBridgeFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
