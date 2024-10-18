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

// IBridgeCallbackMetaData contains all meta data concerning the IBridgeCallback contract.
var IBridgeCallbackMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_refund\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IBridgeCallbackABI is the input ABI used to generate the binding from.
// Deprecated: Use IBridgeCallbackMetaData.ABI instead.
var IBridgeCallbackABI = IBridgeCallbackMetaData.ABI

// IBridgeCallback is an auto generated Go binding around an Ethereum contract.
type IBridgeCallback struct {
	IBridgeCallbackCaller     // Read-only binding to the contract
	IBridgeCallbackTransactor // Write-only binding to the contract
	IBridgeCallbackFilterer   // Log filterer for contract events
}

// IBridgeCallbackCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBridgeCallbackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeCallbackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBridgeCallbackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeCallbackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBridgeCallbackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeCallbackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBridgeCallbackSession struct {
	Contract     *IBridgeCallback  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBridgeCallbackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBridgeCallbackCallerSession struct {
	Contract *IBridgeCallbackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IBridgeCallbackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBridgeCallbackTransactorSession struct {
	Contract     *IBridgeCallbackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IBridgeCallbackRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBridgeCallbackRaw struct {
	Contract *IBridgeCallback // Generic contract binding to access the raw methods on
}

// IBridgeCallbackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBridgeCallbackCallerRaw struct {
	Contract *IBridgeCallbackCaller // Generic read-only contract binding to access the raw methods on
}

// IBridgeCallbackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBridgeCallbackTransactorRaw struct {
	Contract *IBridgeCallbackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBridgeCallback creates a new instance of IBridgeCallback, bound to a specific deployed contract.
func NewIBridgeCallback(address common.Address, backend bind.ContractBackend) (*IBridgeCallback, error) {
	contract, err := bindIBridgeCallback(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallback{IBridgeCallbackCaller: IBridgeCallbackCaller{contract: contract}, IBridgeCallbackTransactor: IBridgeCallbackTransactor{contract: contract}, IBridgeCallbackFilterer: IBridgeCallbackFilterer{contract: contract}}, nil
}

// NewIBridgeCallbackCaller creates a new read-only instance of IBridgeCallback, bound to a specific deployed contract.
func NewIBridgeCallbackCaller(address common.Address, caller bind.ContractCaller) (*IBridgeCallbackCaller, error) {
	contract, err := bindIBridgeCallback(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallbackCaller{contract: contract}, nil
}

// NewIBridgeCallbackTransactor creates a new write-only instance of IBridgeCallback, bound to a specific deployed contract.
func NewIBridgeCallbackTransactor(address common.Address, transactor bind.ContractTransactor) (*IBridgeCallbackTransactor, error) {
	contract, err := bindIBridgeCallback(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallbackTransactor{contract: contract}, nil
}

// NewIBridgeCallbackFilterer creates a new log filterer instance of IBridgeCallback, bound to a specific deployed contract.
func NewIBridgeCallbackFilterer(address common.Address, filterer bind.ContractFilterer) (*IBridgeCallbackFilterer, error) {
	contract, err := bindIBridgeCallback(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallbackFilterer{contract: contract}, nil
}

// bindIBridgeCallback binds a generic wrapper to an already deployed contract.
func bindIBridgeCallback(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBridgeCallbackMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridgeCallback *IBridgeCallbackRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridgeCallback.Contract.IBridgeCallbackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridgeCallback *IBridgeCallbackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridgeCallback.Contract.IBridgeCallbackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridgeCallback *IBridgeCallbackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridgeCallback.Contract.IBridgeCallbackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridgeCallback *IBridgeCallbackCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridgeCallback.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridgeCallback *IBridgeCallbackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridgeCallback.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridgeCallback *IBridgeCallbackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridgeCallback.Contract.contract.Transact(opts, method, params...)
}

// BridgeCallback is a paid mutator transaction binding the contract method 0x13997566.
//
// Solidity: function bridgeCallback(address _sender, address _refund, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo) returns()
func (_IBridgeCallback *IBridgeCallbackTransactor) BridgeCallback(opts *bind.TransactOpts, _sender common.Address, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _data []byte, _memo []byte) (*types.Transaction, error) {
	return _IBridgeCallback.contract.Transact(opts, "bridgeCallback", _sender, _refund, _tokens, _amounts, _data, _memo)
}

// BridgeCallback is a paid mutator transaction binding the contract method 0x13997566.
//
// Solidity: function bridgeCallback(address _sender, address _refund, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo) returns()
func (_IBridgeCallback *IBridgeCallbackSession) BridgeCallback(_sender common.Address, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _data []byte, _memo []byte) (*types.Transaction, error) {
	return _IBridgeCallback.Contract.BridgeCallback(&_IBridgeCallback.TransactOpts, _sender, _refund, _tokens, _amounts, _data, _memo)
}

// BridgeCallback is a paid mutator transaction binding the contract method 0x13997566.
//
// Solidity: function bridgeCallback(address _sender, address _refund, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo) returns()
func (_IBridgeCallback *IBridgeCallbackTransactorSession) BridgeCallback(_sender common.Address, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _data []byte, _memo []byte) (*types.Transaction, error) {
	return _IBridgeCallback.Contract.BridgeCallback(&_IBridgeCallback.TransactOpts, _sender, _refund, _tokens, _amounts, _data, _memo)
}
