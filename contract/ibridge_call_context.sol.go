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

// IBridgeCallContextMetaData contains all meta data concerning the IBridgeCallContext contract.
var IBridgeCallContextMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_refund\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"onBridgeCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_msg\",\"type\":\"bytes\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IBridgeCallContextABI is the input ABI used to generate the binding from.
// Deprecated: Use IBridgeCallContextMetaData.ABI instead.
var IBridgeCallContextABI = IBridgeCallContextMetaData.ABI

// IBridgeCallContext is an auto generated Go binding around an Ethereum contract.
type IBridgeCallContext struct {
	IBridgeCallContextCaller     // Read-only binding to the contract
	IBridgeCallContextTransactor // Write-only binding to the contract
	IBridgeCallContextFilterer   // Log filterer for contract events
}

// IBridgeCallContextCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBridgeCallContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeCallContextTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBridgeCallContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeCallContextFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBridgeCallContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeCallContextSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBridgeCallContextSession struct {
	Contract     *IBridgeCallContext // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IBridgeCallContextCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBridgeCallContextCallerSession struct {
	Contract *IBridgeCallContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IBridgeCallContextTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBridgeCallContextTransactorSession struct {
	Contract     *IBridgeCallContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IBridgeCallContextRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBridgeCallContextRaw struct {
	Contract *IBridgeCallContext // Generic contract binding to access the raw methods on
}

// IBridgeCallContextCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBridgeCallContextCallerRaw struct {
	Contract *IBridgeCallContextCaller // Generic read-only contract binding to access the raw methods on
}

// IBridgeCallContextTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBridgeCallContextTransactorRaw struct {
	Contract *IBridgeCallContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBridgeCallContext creates a new instance of IBridgeCallContext, bound to a specific deployed contract.
func NewIBridgeCallContext(address common.Address, backend bind.ContractBackend) (*IBridgeCallContext, error) {
	contract, err := bindIBridgeCallContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallContext{IBridgeCallContextCaller: IBridgeCallContextCaller{contract: contract}, IBridgeCallContextTransactor: IBridgeCallContextTransactor{contract: contract}, IBridgeCallContextFilterer: IBridgeCallContextFilterer{contract: contract}}, nil
}

// NewIBridgeCallContextCaller creates a new read-only instance of IBridgeCallContext, bound to a specific deployed contract.
func NewIBridgeCallContextCaller(address common.Address, caller bind.ContractCaller) (*IBridgeCallContextCaller, error) {
	contract, err := bindIBridgeCallContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallContextCaller{contract: contract}, nil
}

// NewIBridgeCallContextTransactor creates a new write-only instance of IBridgeCallContext, bound to a specific deployed contract.
func NewIBridgeCallContextTransactor(address common.Address, transactor bind.ContractTransactor) (*IBridgeCallContextTransactor, error) {
	contract, err := bindIBridgeCallContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallContextTransactor{contract: contract}, nil
}

// NewIBridgeCallContextFilterer creates a new log filterer instance of IBridgeCallContext, bound to a specific deployed contract.
func NewIBridgeCallContextFilterer(address common.Address, filterer bind.ContractFilterer) (*IBridgeCallContextFilterer, error) {
	contract, err := bindIBridgeCallContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBridgeCallContextFilterer{contract: contract}, nil
}

// bindIBridgeCallContext binds a generic wrapper to an already deployed contract.
func bindIBridgeCallContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBridgeCallContextMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridgeCallContext *IBridgeCallContextRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridgeCallContext.Contract.IBridgeCallContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridgeCallContext *IBridgeCallContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.IBridgeCallContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridgeCallContext *IBridgeCallContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.IBridgeCallContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridgeCallContext *IBridgeCallContextCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridgeCallContext.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridgeCallContext *IBridgeCallContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridgeCallContext *IBridgeCallContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.contract.Transact(opts, method, params...)
}

// OnBridgeCall is a paid mutator transaction binding the contract method 0x57ffc092.
//
// Solidity: function onBridgeCall(address _sender, address _refund, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo) returns()
func (_IBridgeCallContext *IBridgeCallContextTransactor) OnBridgeCall(opts *bind.TransactOpts, _sender common.Address, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _data []byte, _memo []byte) (*types.Transaction, error) {
	return _IBridgeCallContext.contract.Transact(opts, "onBridgeCall", _sender, _refund, _tokens, _amounts, _data, _memo)
}

// OnBridgeCall is a paid mutator transaction binding the contract method 0x57ffc092.
//
// Solidity: function onBridgeCall(address _sender, address _refund, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo) returns()
func (_IBridgeCallContext *IBridgeCallContextSession) OnBridgeCall(_sender common.Address, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _data []byte, _memo []byte) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.OnBridgeCall(&_IBridgeCallContext.TransactOpts, _sender, _refund, _tokens, _amounts, _data, _memo)
}

// OnBridgeCall is a paid mutator transaction binding the contract method 0x57ffc092.
//
// Solidity: function onBridgeCall(address _sender, address _refund, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo) returns()
func (_IBridgeCallContext *IBridgeCallContextTransactorSession) OnBridgeCall(_sender common.Address, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _data []byte, _memo []byte) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.OnBridgeCall(&_IBridgeCallContext.TransactOpts, _sender, _refund, _tokens, _amounts, _data, _memo)
}

// OnRevert is a paid mutator transaction binding the contract method 0x32e1e16e.
//
// Solidity: function onRevert(uint256 nonce, bytes _msg) returns()
func (_IBridgeCallContext *IBridgeCallContextTransactor) OnRevert(opts *bind.TransactOpts, nonce *big.Int, _msg []byte) (*types.Transaction, error) {
	return _IBridgeCallContext.contract.Transact(opts, "onRevert", nonce, _msg)
}

// OnRevert is a paid mutator transaction binding the contract method 0x32e1e16e.
//
// Solidity: function onRevert(uint256 nonce, bytes _msg) returns()
func (_IBridgeCallContext *IBridgeCallContextSession) OnRevert(nonce *big.Int, _msg []byte) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.OnRevert(&_IBridgeCallContext.TransactOpts, nonce, _msg)
}

// OnRevert is a paid mutator transaction binding the contract method 0x32e1e16e.
//
// Solidity: function onRevert(uint256 nonce, bytes _msg) returns()
func (_IBridgeCallContext *IBridgeCallContextTransactorSession) OnRevert(nonce *big.Int, _msg []byte) (*types.Transaction, error) {
	return _IBridgeCallContext.Contract.OnRevert(&_IBridgeCallContext.TransactOpts, nonce, _msg)
}
