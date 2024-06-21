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
	ABI: "[{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"increaseBridgeFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610907806100206000396000f3fe60806040526004361061004a5760003560e01c80630b56c1901461004f578063160d7c73146100845780638fefb76514610097578063c79a6b7b146100c5578063f7356475146100d8575b600080fd5b34801561005b57600080fd5b5061006f61006a3660046105fa565b610106565b60405190151581526020015b60405180910390f35b61006f61009236600461065b565b610179565b3480156100a357600080fd5b506100b76100b23660046106eb565b610469565b60405190815260200161007b565b61006f6100d3366004610715565b6104dd565b3480156100e457600080fd5b506100ee61100481565b6040516001600160a01b03909116815260200161007b565b60405162b56c1960e41b815260009061100490630b56c1909061012f90869086906004016107c1565b6020604051808303816000875af115801561014e573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017291906107e3565b9392505050565b60006001600160a01b038716156102ab576001600160a01b0387166323b872dd33306101a5888a610805565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af11580156101f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021d91906107e3565b506001600160a01b03871663095ea7b361100461023a8789610805565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015610285573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102a991906107e3565b505b6001600160a01b0387161561039257604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e90604401602060405180830381865afa158015610309573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032d919061082b565b90506103398587610805565b811461038c5760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b506103ea565b61039c8486610805565b34146103ea5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b206665656044820152606401610383565b60405163160d7c7360e01b81526110049063160d7c7390349061041b908b908b908b908b908b908b90600401610844565b60206040518083038185885af1158015610439573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061045e91906107e3565b979650505050505050565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa1580156104b9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610172919061082b565b60405163c79a6b7b60e01b81526000906110049063c79a6b7b9061050b908890889088908890600401610899565b6020604051808303816000875af115801561052a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061054e91906107e3565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261057e57600080fd5b813567ffffffffffffffff8082111561059957610599610557565b604051601f8301601f19908116603f011681019082821181831017156105c1576105c1610557565b816040528381528660208588010111156105da57600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000806040838503121561060d57600080fd5b823567ffffffffffffffff81111561062457600080fd5b6106308582860161056d565b95602094909401359450505050565b80356001600160a01b038116811461065657600080fd5b919050565b60008060008060008060c0878903121561067457600080fd5b61067d8761063f565b9550602087013567ffffffffffffffff8082111561069a57600080fd5b6106a68a838b0161056d565b965060408901359550606089013594506080890135935060a08901359150808211156106d157600080fd5b506106de89828a0161056d565b9150509295509295509295565b600080604083850312156106fe57600080fd5b6107078361063f565b946020939093013593505050565b6000806000806080858703121561072b57600080fd5b843567ffffffffffffffff81111561074257600080fd5b61074e8782880161056d565b945050602085013592506107646040860161063f565b9396929550929360600135925050565b6000815180845260005b8181101561079a5760208185018101518683018201520161077e565b818111156107ac576000602083870101525b50601f01601f19169290920160200192915050565b6040815260006107d46040830185610774565b90508260208301529392505050565b6000602082840312156107f557600080fd5b8151801515811461017257600080fd5b6000821982111561082657634e487b7160e01b600052601160045260246000fd5b500190565b60006020828403121561083d57600080fd5b5051919050565b6001600160a01b038716815260c06020820181905260009061086890830188610774565b86604084015285606084015284608084015282810360a084015261088c8185610774565b9998505050505050505050565b6080815260006108ac6080830187610774565b6020830195909552506001600160a01b0392909216604083015260609091015291905056fea2646970667358221220d2bf11cd0f68982053846dd6c9d7c9b5241feb670b1603cccea9a5e82736e47264736f6c634300080a0033",
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

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrossChainTest *CrossChainTestTransactor) CancelSendToExternal(opts *bind.TransactOpts, _chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.contract.Transact(opts, "cancelSendToExternal", _chain, _txID)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrossChainTest *CrossChainTestSession) CancelSendToExternal(_chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CancelSendToExternal(&_CrossChainTest.TransactOpts, _chain, _txID)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrossChainTest *CrossChainTestTransactorSession) CancelSendToExternal(_chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CancelSendToExternal(&_CrossChainTest.TransactOpts, _chain, _txID)
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

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactor) IncreaseBridgeFee(opts *bind.TransactOpts, _chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.contract.Transact(opts, "increaseBridgeFee", _chain, _txID, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrossChainTest *CrossChainTestSession) IncreaseBridgeFee(_chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.IncreaseBridgeFee(&_CrossChainTest.TransactOpts, _chain, _txID, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactorSession) IncreaseBridgeFee(_chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.IncreaseBridgeFee(&_CrossChainTest.TransactOpts, _chain, _txID, _token, _fee)
}
