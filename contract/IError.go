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

// IErrorMetaData contains all meta data concerning the IError contract.
var IErrorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"Error\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IErrorABI is the input ABI used to generate the binding from.
// Deprecated: Use IErrorMetaData.ABI instead.
var IErrorABI = IErrorMetaData.ABI

// IError is an auto generated Go binding around an Ethereum contract.
type IError struct {
	IErrorCaller     // Read-only binding to the contract
	IErrorTransactor // Write-only binding to the contract
	IErrorFilterer   // Log filterer for contract events
}

// IErrorCaller is an auto generated read-only Go binding around an Ethereum contract.
type IErrorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IErrorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IErrorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IErrorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IErrorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IErrorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IErrorSession struct {
	Contract     *IError           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IErrorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IErrorCallerSession struct {
	Contract *IErrorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IErrorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IErrorTransactorSession struct {
	Contract     *IErrorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IErrorRaw is an auto generated low-level Go binding around an Ethereum contract.
type IErrorRaw struct {
	Contract *IError // Generic contract binding to access the raw methods on
}

// IErrorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IErrorCallerRaw struct {
	Contract *IErrorCaller // Generic read-only contract binding to access the raw methods on
}

// IErrorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IErrorTransactorRaw struct {
	Contract *IErrorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIError creates a new instance of IError, bound to a specific deployed contract.
func NewIError(address common.Address, backend bind.ContractBackend) (*IError, error) {
	contract, err := bindIError(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IError{IErrorCaller: IErrorCaller{contract: contract}, IErrorTransactor: IErrorTransactor{contract: contract}, IErrorFilterer: IErrorFilterer{contract: contract}}, nil
}

// NewIErrorCaller creates a new read-only instance of IError, bound to a specific deployed contract.
func NewIErrorCaller(address common.Address, caller bind.ContractCaller) (*IErrorCaller, error) {
	contract, err := bindIError(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IErrorCaller{contract: contract}, nil
}

// NewIErrorTransactor creates a new write-only instance of IError, bound to a specific deployed contract.
func NewIErrorTransactor(address common.Address, transactor bind.ContractTransactor) (*IErrorTransactor, error) {
	contract, err := bindIError(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IErrorTransactor{contract: contract}, nil
}

// NewIErrorFilterer creates a new log filterer instance of IError, bound to a specific deployed contract.
func NewIErrorFilterer(address common.Address, filterer bind.ContractFilterer) (*IErrorFilterer, error) {
	contract, err := bindIError(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IErrorFilterer{contract: contract}, nil
}

// bindIError binds a generic wrapper to an already deployed contract.
func bindIError(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IErrorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IError *IErrorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IError.Contract.IErrorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IError *IErrorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IError.Contract.IErrorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IError *IErrorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IError.Contract.IErrorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IError *IErrorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IError.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IError *IErrorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IError.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IError *IErrorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IError.Contract.contract.Transact(opts, method, params...)
}

// Error is a paid mutator transaction binding the contract method 0x08c379a0.
//
// Solidity: function Error(string ) returns()
func (_IError *IErrorTransactor) Error(opts *bind.TransactOpts, arg0 string) (*types.Transaction, error) {
	return _IError.contract.Transact(opts, "Error", arg0)
}

// Error is a paid mutator transaction binding the contract method 0x08c379a0.
//
// Solidity: function Error(string ) returns()
func (_IError *IErrorSession) Error(arg0 string) (*types.Transaction, error) {
	return _IError.Contract.Error(&_IError.TransactOpts, arg0)
}

// Error is a paid mutator transaction binding the contract method 0x08c379a0.
//
// Solidity: function Error(string ) returns()
func (_IError *IErrorTransactorSession) Error(arg0 string) (*types.Transaction, error) {
	return _IError.Contract.Error(&_IError.TransactOpts, arg0)
}
