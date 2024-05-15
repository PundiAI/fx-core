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

// IRefundCallbackMetaData contains all meta data concerning the IRefundCallback contract.
var IRefundCallbackMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"name\":\"refundCallback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IRefundCallbackABI is the input ABI used to generate the binding from.
// Deprecated: Use IRefundCallbackMetaData.ABI instead.
var IRefundCallbackABI = IRefundCallbackMetaData.ABI

// IRefundCallback is an auto generated Go binding around an Ethereum contract.
type IRefundCallback struct {
	IRefundCallbackCaller     // Read-only binding to the contract
	IRefundCallbackTransactor // Write-only binding to the contract
	IRefundCallbackFilterer   // Log filterer for contract events
}

// IRefundCallbackCaller is an auto generated read-only Go binding around an Ethereum contract.
type IRefundCallbackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRefundCallbackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IRefundCallbackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRefundCallbackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IRefundCallbackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRefundCallbackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IRefundCallbackSession struct {
	Contract     *IRefundCallback  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IRefundCallbackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IRefundCallbackCallerSession struct {
	Contract *IRefundCallbackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IRefundCallbackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IRefundCallbackTransactorSession struct {
	Contract     *IRefundCallbackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IRefundCallbackRaw is an auto generated low-level Go binding around an Ethereum contract.
type IRefundCallbackRaw struct {
	Contract *IRefundCallback // Generic contract binding to access the raw methods on
}

// IRefundCallbackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IRefundCallbackCallerRaw struct {
	Contract *IRefundCallbackCaller // Generic read-only contract binding to access the raw methods on
}

// IRefundCallbackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IRefundCallbackTransactorRaw struct {
	Contract *IRefundCallbackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIRefundCallback creates a new instance of IRefundCallback, bound to a specific deployed contract.
func NewIRefundCallback(address common.Address, backend bind.ContractBackend) (*IRefundCallback, error) {
	contract, err := bindIRefundCallback(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IRefundCallback{IRefundCallbackCaller: IRefundCallbackCaller{contract: contract}, IRefundCallbackTransactor: IRefundCallbackTransactor{contract: contract}, IRefundCallbackFilterer: IRefundCallbackFilterer{contract: contract}}, nil
}

// NewIRefundCallbackCaller creates a new read-only instance of IRefundCallback, bound to a specific deployed contract.
func NewIRefundCallbackCaller(address common.Address, caller bind.ContractCaller) (*IRefundCallbackCaller, error) {
	contract, err := bindIRefundCallback(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IRefundCallbackCaller{contract: contract}, nil
}

// NewIRefundCallbackTransactor creates a new write-only instance of IRefundCallback, bound to a specific deployed contract.
func NewIRefundCallbackTransactor(address common.Address, transactor bind.ContractTransactor) (*IRefundCallbackTransactor, error) {
	contract, err := bindIRefundCallback(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IRefundCallbackTransactor{contract: contract}, nil
}

// NewIRefundCallbackFilterer creates a new log filterer instance of IRefundCallback, bound to a specific deployed contract.
func NewIRefundCallbackFilterer(address common.Address, filterer bind.ContractFilterer) (*IRefundCallbackFilterer, error) {
	contract, err := bindIRefundCallback(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IRefundCallbackFilterer{contract: contract}, nil
}

// bindIRefundCallback binds a generic wrapper to an already deployed contract.
func bindIRefundCallback(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IRefundCallbackMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRefundCallback *IRefundCallbackRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRefundCallback.Contract.IRefundCallbackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRefundCallback *IRefundCallbackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRefundCallback.Contract.IRefundCallbackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRefundCallback *IRefundCallbackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRefundCallback.Contract.IRefundCallbackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRefundCallback *IRefundCallbackCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRefundCallback.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRefundCallback *IRefundCallbackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRefundCallback.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRefundCallback *IRefundCallbackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRefundCallback.Contract.contract.Transact(opts, method, params...)
}

// RefundCallback is a paid mutator transaction binding the contract method 0x3a37fb2e.
//
// Solidity: function refundCallback(uint256 , address[] , uint256[] ) returns()
func (_IRefundCallback *IRefundCallbackTransactor) RefundCallback(opts *bind.TransactOpts, arg0 *big.Int, arg1 []common.Address, arg2 []*big.Int) (*types.Transaction, error) {
	return _IRefundCallback.contract.Transact(opts, "refundCallback", arg0, arg1, arg2)
}

// RefundCallback is a paid mutator transaction binding the contract method 0x3a37fb2e.
//
// Solidity: function refundCallback(uint256 , address[] , uint256[] ) returns()
func (_IRefundCallback *IRefundCallbackSession) RefundCallback(arg0 *big.Int, arg1 []common.Address, arg2 []*big.Int) (*types.Transaction, error) {
	return _IRefundCallback.Contract.RefundCallback(&_IRefundCallback.TransactOpts, arg0, arg1, arg2)
}

// RefundCallback is a paid mutator transaction binding the contract method 0x3a37fb2e.
//
// Solidity: function refundCallback(uint256 , address[] , uint256[] ) returns()
func (_IRefundCallback *IRefundCallbackTransactorSession) RefundCallback(arg0 *big.Int, arg1 []common.Address, arg2 []*big.Int) (*types.Transaction, error) {
	return _IRefundCallback.Contract.RefundCallback(&_IRefundCallback.TransactOpts, arg0, arg1, arg2)
}
