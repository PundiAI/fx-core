// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// CrossChainTestMetaData contains all meta data concerning the CrossChainTest contract.
var CrossChainTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506106da806100206000396000f3fe6080604052600436106100345760003560e01c8063160d7c73146100395780638fefb76514610061578063f73564751461008f575b600080fd5b61004c6100473660046104e7565b6100bd565b60405190151581526020015b60405180910390f35b34801561006d57600080fd5b5061008161007c366004610577565b6103ad565b604051908152602001610058565b34801561009b57600080fd5b506100a561100481565b6040516001600160a01b039091168152602001610058565b60006001600160a01b038716156101ef576001600160a01b0387166323b872dd33306100e9888a6105a1565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af115801561013d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061016191906105c7565b506001600160a01b03871663095ea7b361100461017e87896105a1565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af11580156101c9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101ed91906105c7565b505b6001600160a01b038716156102d657604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e90604401602060405180830381865afa15801561024d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027191906105e9565b905061027d85876105a1565b81146102d05760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b5061032e565b6102e084866105a1565b341461032e5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b2066656560448201526064016102c7565b60405163160d7c7360e01b81526110049063160d7c7390349061035f908b908b908b908b908b908b9060040161064f565b60206040518083038185885af115801561037d573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906103a291906105c7565b979650505050505050565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa1580156103fd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061042191906105e9565b9392505050565b80356001600160a01b038116811461043f57600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261046b57600080fd5b813567ffffffffffffffff8082111561048657610486610444565b604051601f8301601f19908116603f011681019082821181831017156104ae576104ae610444565b816040528381528660208588010111156104c757600080fd5b836020870160208301376000602085830101528094505050505092915050565b60008060008060008060c0878903121561050057600080fd5b61050987610428565b9550602087013567ffffffffffffffff8082111561052657600080fd5b6105328a838b0161045a565b965060408901359550606089013594506080890135935060a089013591508082111561055d57600080fd5b5061056a89828a0161045a565b9150509295509295509295565b6000806040838503121561058a57600080fd5b61059383610428565b946020939093013593505050565b600082198211156105c257634e487b7160e01b600052601160045260246000fd5b500190565b6000602082840312156105d957600080fd5b8151801515811461042157600080fd5b6000602082840312156105fb57600080fd5b5051919050565b6000815180845260005b818110156106285760208185018101518683018201520161060c565b8181111561063a576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038716815260c06020820181905260009061067390830188610602565b86604084015285606084015284608084015282810360a08401526106978185610602565b999850505050505050505056fea2646970667358221220a7820ef250f22731ab00f08d0d3f7a29e1d9270f6947e9ea594e944b105882e964736f6c634300080a0033",
}

// CrossChainTestABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainTestMetaData.ABI instead.
var CrossChainTestABI = CrossChainTestMetaData.ABI

// CrossChainTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainTestMetaData.Bin instead.
var CrossChainTestBin = CrossChainTestMetaData.Bin

// DeployCrossChainTest deploys a new Ethereum contract, binding an instance of CrossChainTest to it.
func DeployCrossChainTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChainTest, error) {
	parsed, err := CrossChainTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainTest{CrossChainTestCaller: CrossChainTestCaller{contract: contract}, CrossChainTestTransactor: CrossChainTestTransactor{contract: contract}, CrossChainTestFilterer: CrossChainTestFilterer{contract: contract}}, nil
}

// CrossChainTest is an auto generated Go binding around an Ethereum contract.
type CrossChainTest struct {
	CrossChainTestCaller     // Read-only binding to the contract
	CrossChainTestTransactor // Write-only binding to the contract
	CrossChainTestFilterer   // Log filterer for contract events
}

// CrossChainTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainTestSession struct {
	Contract     *CrossChainTest   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrossChainTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainTestCallerSession struct {
	Contract *CrossChainTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CrossChainTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainTestTransactorSession struct {
	Contract     *CrossChainTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CrossChainTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainTestRaw struct {
	Contract *CrossChainTest // Generic contract binding to access the raw methods on
}

// CrossChainTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainTestCallerRaw struct {
	Contract *CrossChainTestCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainTestTransactorRaw struct {
	Contract *CrossChainTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChainTest creates a new instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTest(address common.Address, backend bind.ContractBackend) (*CrossChainTest, error) {
	contract, err := bindCrossChainTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainTest{CrossChainTestCaller: CrossChainTestCaller{contract: contract}, CrossChainTestTransactor: CrossChainTestTransactor{contract: contract}, CrossChainTestFilterer: CrossChainTestFilterer{contract: contract}}, nil
}

// NewCrossChainTestCaller creates a new read-only instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTestCaller(address common.Address, caller bind.ContractCaller) (*CrossChainTestCaller, error) {
	contract, err := bindCrossChainTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainTestCaller{contract: contract}, nil
}

// NewCrossChainTestTransactor creates a new write-only instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTestTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainTestTransactor, error) {
	contract, err := bindCrossChainTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainTestTransactor{contract: contract}, nil
}

// NewCrossChainTestFilterer creates a new log filterer instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTestFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainTestFilterer, error) {
	contract, err := bindCrossChainTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainTestFilterer{contract: contract}, nil
}

// bindCrossChainTest binds a generic wrapper to an already deployed contract.
func bindCrossChainTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrossChainTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainTest *CrossChainTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainTest.Contract.CrossChainTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainTest *CrossChainTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChainTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainTest *CrossChainTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChainTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainTest *CrossChainTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainTest *CrossChainTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainTest *CrossChainTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainTest.Contract.contract.Transact(opts, method, params...)
}

// CROSSCHAINADDRESS is a free data retrieval call binding the contract method 0xf7356475.
//
// Solidity: function CROSS_CHAIN_ADDRESS() view returns(address)
func (_CrossChainTest *CrossChainTestCaller) CROSSCHAINADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrossChainTest.contract.Call(opts, &out, "CROSS_CHAIN_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CROSSCHAINADDRESS is a free data retrieval call binding the contract method 0xf7356475.
//
// Solidity: function CROSS_CHAIN_ADDRESS() view returns(address)
func (_CrossChainTest *CrossChainTestSession) CROSSCHAINADDRESS() (common.Address, error) {
	return _CrossChainTest.Contract.CROSSCHAINADDRESS(&_CrossChainTest.CallOpts)
}

// CROSSCHAINADDRESS is a free data retrieval call binding the contract method 0xf7356475.
//
// Solidity: function CROSS_CHAIN_ADDRESS() view returns(address)
func (_CrossChainTest *CrossChainTestCallerSession) CROSSCHAINADDRESS() (common.Address, error) {
	return _CrossChainTest.Contract.CROSSCHAINADDRESS(&_CrossChainTest.CallOpts)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrossChainTest *CrossChainTestCaller) BridgeCoinAmount(opts *bind.CallOpts, _token common.Address, _target [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _CrossChainTest.contract.Call(opts, &out, "bridgeCoinAmount", _token, _target)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrossChainTest *CrossChainTestSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _CrossChainTest.Contract.BridgeCoinAmount(&_CrossChainTest.CallOpts, _token, _target)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrossChainTest *CrossChainTestCallerSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _CrossChainTest.Contract.BridgeCoinAmount(&_CrossChainTest.CallOpts, _token, _target)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactor) CrossChain(opts *bind.TransactOpts, _token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrossChainTest.contract.Transact(opts, "crossChain", _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrossChainTest *CrossChainTestSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChain(&_CrossChainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactorSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChain(&_CrossChainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}
