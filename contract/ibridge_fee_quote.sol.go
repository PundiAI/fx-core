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

// IBridgeFeeQuoteQuoteInfo is an auto generated low-level Go binding around an user-defined struct.
type IBridgeFeeQuoteQuoteInfo struct {
	Id        *big.Int
	ChainName string
	Token     common.Address
	Oracle    common.Address
	Fee       *big.Int
	GasLimit  *big.Int
	Expiry    *big.Int
}

// IBridgeFeeQuoteQuoteInput is an auto generated low-level Go binding around an user-defined struct.
type IBridgeFeeQuoteQuoteInput struct {
	ChainName  string
	Token      common.Address
	Oracle     common.Address
	QuoteIndex *big.Int
	Fee        *big.Int
	GasLimit   *big.Int
	Expiry     *big.Int
	Signature  []byte
}

// IBridgeFeeQuoteMetaData contains all meta data concerning the IBridgeFeeQuote contract.
var IBridgeFeeQuoteMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getQuote\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"internalType\":\"structIBridgeFeeQuote.QuoteInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getQuoteById\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"internalType\":\"structIBridgeFeeQuote.QuoteInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getQuoteByToken\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"internalType\":\"structIBridgeFeeQuote.QuoteInfo[]\",\"name\":\"quotes\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chainName\",\"type\":\"string\"}],\"name\":\"getQuoteList\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"}],\"internalType\":\"structIBridgeFeeQuote.QuoteInfo[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"quoteIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structIBridgeFeeQuote.QuoteInput[]\",\"name\":\"_inputs\",\"type\":\"tuple[]\"}],\"name\":\"quote\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IBridgeFeeQuoteABI is the input ABI used to generate the binding from.
// Deprecated: Use IBridgeFeeQuoteMetaData.ABI instead.
var IBridgeFeeQuoteABI = IBridgeFeeQuoteMetaData.ABI

// IBridgeFeeQuote is an auto generated Go binding around an Ethereum contract.
type IBridgeFeeQuote struct {
	IBridgeFeeQuoteCaller     // Read-only binding to the contract
	IBridgeFeeQuoteTransactor // Write-only binding to the contract
	IBridgeFeeQuoteFilterer   // Log filterer for contract events
}

// IBridgeFeeQuoteCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBridgeFeeQuoteCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeFeeQuoteTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBridgeFeeQuoteTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeFeeQuoteFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBridgeFeeQuoteFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeFeeQuoteSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBridgeFeeQuoteSession struct {
	Contract     *IBridgeFeeQuote  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBridgeFeeQuoteCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBridgeFeeQuoteCallerSession struct {
	Contract *IBridgeFeeQuoteCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IBridgeFeeQuoteTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBridgeFeeQuoteTransactorSession struct {
	Contract     *IBridgeFeeQuoteTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IBridgeFeeQuoteRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBridgeFeeQuoteRaw struct {
	Contract *IBridgeFeeQuote // Generic contract binding to access the raw methods on
}

// IBridgeFeeQuoteCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBridgeFeeQuoteCallerRaw struct {
	Contract *IBridgeFeeQuoteCaller // Generic read-only contract binding to access the raw methods on
}

// IBridgeFeeQuoteTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBridgeFeeQuoteTransactorRaw struct {
	Contract *IBridgeFeeQuoteTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBridgeFeeQuote creates a new instance of IBridgeFeeQuote, bound to a specific deployed contract.
func NewIBridgeFeeQuote(address common.Address, backend bind.ContractBackend) (*IBridgeFeeQuote, error) {
	contract, err := bindIBridgeFeeQuote(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBridgeFeeQuote{IBridgeFeeQuoteCaller: IBridgeFeeQuoteCaller{contract: contract}, IBridgeFeeQuoteTransactor: IBridgeFeeQuoteTransactor{contract: contract}, IBridgeFeeQuoteFilterer: IBridgeFeeQuoteFilterer{contract: contract}}, nil
}

// NewIBridgeFeeQuoteCaller creates a new read-only instance of IBridgeFeeQuote, bound to a specific deployed contract.
func NewIBridgeFeeQuoteCaller(address common.Address, caller bind.ContractCaller) (*IBridgeFeeQuoteCaller, error) {
	contract, err := bindIBridgeFeeQuote(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeFeeQuoteCaller{contract: contract}, nil
}

// NewIBridgeFeeQuoteTransactor creates a new write-only instance of IBridgeFeeQuote, bound to a specific deployed contract.
func NewIBridgeFeeQuoteTransactor(address common.Address, transactor bind.ContractTransactor) (*IBridgeFeeQuoteTransactor, error) {
	contract, err := bindIBridgeFeeQuote(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeFeeQuoteTransactor{contract: contract}, nil
}

// NewIBridgeFeeQuoteFilterer creates a new log filterer instance of IBridgeFeeQuote, bound to a specific deployed contract.
func NewIBridgeFeeQuoteFilterer(address common.Address, filterer bind.ContractFilterer) (*IBridgeFeeQuoteFilterer, error) {
	contract, err := bindIBridgeFeeQuote(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBridgeFeeQuoteFilterer{contract: contract}, nil
}

// bindIBridgeFeeQuote binds a generic wrapper to an already deployed contract.
func bindIBridgeFeeQuote(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBridgeFeeQuoteMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridgeFeeQuote *IBridgeFeeQuoteRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridgeFeeQuote.Contract.IBridgeFeeQuoteCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridgeFeeQuote *IBridgeFeeQuoteRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridgeFeeQuote.Contract.IBridgeFeeQuoteTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridgeFeeQuote *IBridgeFeeQuoteRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridgeFeeQuote.Contract.IBridgeFeeQuoteTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridgeFeeQuote *IBridgeFeeQuoteCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridgeFeeQuote.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridgeFeeQuote *IBridgeFeeQuoteTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridgeFeeQuote.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridgeFeeQuote *IBridgeFeeQuoteTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridgeFeeQuote.Contract.contract.Transact(opts, method, params...)
}

// GetQuote is a free data retrieval call binding the contract method 0xb02e61a5.
//
// Solidity: function getQuote(string _chainName, address _token, address _oracle, uint256 _index) view returns((uint256,string,address,address,uint256,uint256,uint256))
func (_IBridgeFeeQuote *IBridgeFeeQuoteCaller) GetQuote(opts *bind.CallOpts, _chainName string, _token common.Address, _oracle common.Address, _index *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	var out []interface{}
	err := _IBridgeFeeQuote.contract.Call(opts, &out, "getQuote", _chainName, _token, _oracle, _index)

	if err != nil {
		return *new(IBridgeFeeQuoteQuoteInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IBridgeFeeQuoteQuoteInfo)).(*IBridgeFeeQuoteQuoteInfo)

	return out0, err

}

// GetQuote is a free data retrieval call binding the contract method 0xb02e61a5.
//
// Solidity: function getQuote(string _chainName, address _token, address _oracle, uint256 _index) view returns((uint256,string,address,address,uint256,uint256,uint256))
func (_IBridgeFeeQuote *IBridgeFeeQuoteSession) GetQuote(_chainName string, _token common.Address, _oracle common.Address, _index *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuote(&_IBridgeFeeQuote.CallOpts, _chainName, _token, _oracle, _index)
}

// GetQuote is a free data retrieval call binding the contract method 0xb02e61a5.
//
// Solidity: function getQuote(string _chainName, address _token, address _oracle, uint256 _index) view returns((uint256,string,address,address,uint256,uint256,uint256))
func (_IBridgeFeeQuote *IBridgeFeeQuoteCallerSession) GetQuote(_chainName string, _token common.Address, _oracle common.Address, _index *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuote(&_IBridgeFeeQuote.CallOpts, _chainName, _token, _oracle, _index)
}

// GetQuoteById is a free data retrieval call binding the contract method 0xa8541c17.
//
// Solidity: function getQuoteById(uint256 _id) view returns((uint256,string,address,address,uint256,uint256,uint256))
func (_IBridgeFeeQuote *IBridgeFeeQuoteCaller) GetQuoteById(opts *bind.CallOpts, _id *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	var out []interface{}
	err := _IBridgeFeeQuote.contract.Call(opts, &out, "getQuoteById", _id)

	if err != nil {
		return *new(IBridgeFeeQuoteQuoteInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IBridgeFeeQuoteQuoteInfo)).(*IBridgeFeeQuoteQuoteInfo)

	return out0, err

}

// GetQuoteById is a free data retrieval call binding the contract method 0xa8541c17.
//
// Solidity: function getQuoteById(uint256 _id) view returns((uint256,string,address,address,uint256,uint256,uint256))
func (_IBridgeFeeQuote *IBridgeFeeQuoteSession) GetQuoteById(_id *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuoteById(&_IBridgeFeeQuote.CallOpts, _id)
}

// GetQuoteById is a free data retrieval call binding the contract method 0xa8541c17.
//
// Solidity: function getQuoteById(uint256 _id) view returns((uint256,string,address,address,uint256,uint256,uint256))
func (_IBridgeFeeQuote *IBridgeFeeQuoteCallerSession) GetQuoteById(_id *big.Int) (IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuoteById(&_IBridgeFeeQuote.CallOpts, _id)
}

// GetQuoteByToken is a free data retrieval call binding the contract method 0x38fbcf5b.
//
// Solidity: function getQuoteByToken(string _chainName, address _token) view returns((uint256,string,address,address,uint256,uint256,uint256)[] quotes)
func (_IBridgeFeeQuote *IBridgeFeeQuoteCaller) GetQuoteByToken(opts *bind.CallOpts, _chainName string, _token common.Address) ([]IBridgeFeeQuoteQuoteInfo, error) {
	var out []interface{}
	err := _IBridgeFeeQuote.contract.Call(opts, &out, "getQuoteByToken", _chainName, _token)

	if err != nil {
		return *new([]IBridgeFeeQuoteQuoteInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]IBridgeFeeQuoteQuoteInfo)).(*[]IBridgeFeeQuoteQuoteInfo)

	return out0, err

}

// GetQuoteByToken is a free data retrieval call binding the contract method 0x38fbcf5b.
//
// Solidity: function getQuoteByToken(string _chainName, address _token) view returns((uint256,string,address,address,uint256,uint256,uint256)[] quotes)
func (_IBridgeFeeQuote *IBridgeFeeQuoteSession) GetQuoteByToken(_chainName string, _token common.Address) ([]IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuoteByToken(&_IBridgeFeeQuote.CallOpts, _chainName, _token)
}

// GetQuoteByToken is a free data retrieval call binding the contract method 0x38fbcf5b.
//
// Solidity: function getQuoteByToken(string _chainName, address _token) view returns((uint256,string,address,address,uint256,uint256,uint256)[] quotes)
func (_IBridgeFeeQuote *IBridgeFeeQuoteCallerSession) GetQuoteByToken(_chainName string, _token common.Address) ([]IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuoteByToken(&_IBridgeFeeQuote.CallOpts, _chainName, _token)
}

// GetQuoteList is a free data retrieval call binding the contract method 0x398a0e6b.
//
// Solidity: function getQuoteList(string _chainName) view returns((uint256,string,address,address,uint256,uint256,uint256)[])
func (_IBridgeFeeQuote *IBridgeFeeQuoteCaller) GetQuoteList(opts *bind.CallOpts, _chainName string) ([]IBridgeFeeQuoteQuoteInfo, error) {
	var out []interface{}
	err := _IBridgeFeeQuote.contract.Call(opts, &out, "getQuoteList", _chainName)

	if err != nil {
		return *new([]IBridgeFeeQuoteQuoteInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]IBridgeFeeQuoteQuoteInfo)).(*[]IBridgeFeeQuoteQuoteInfo)

	return out0, err

}

// GetQuoteList is a free data retrieval call binding the contract method 0x398a0e6b.
//
// Solidity: function getQuoteList(string _chainName) view returns((uint256,string,address,address,uint256,uint256,uint256)[])
func (_IBridgeFeeQuote *IBridgeFeeQuoteSession) GetQuoteList(_chainName string) ([]IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuoteList(&_IBridgeFeeQuote.CallOpts, _chainName)
}

// GetQuoteList is a free data retrieval call binding the contract method 0x398a0e6b.
//
// Solidity: function getQuoteList(string _chainName) view returns((uint256,string,address,address,uint256,uint256,uint256)[])
func (_IBridgeFeeQuote *IBridgeFeeQuoteCallerSession) GetQuoteList(_chainName string) ([]IBridgeFeeQuoteQuoteInfo, error) {
	return _IBridgeFeeQuote.Contract.GetQuoteList(&_IBridgeFeeQuote.CallOpts, _chainName)
}

// Quote is a paid mutator transaction binding the contract method 0x71a141c6.
//
// Solidity: function quote((string,address,address,uint256,uint256,uint256,uint256,bytes)[] _inputs) returns(bool)
func (_IBridgeFeeQuote *IBridgeFeeQuoteTransactor) Quote(opts *bind.TransactOpts, _inputs []IBridgeFeeQuoteQuoteInput) (*types.Transaction, error) {
	return _IBridgeFeeQuote.contract.Transact(opts, "quote", _inputs)
}

// Quote is a paid mutator transaction binding the contract method 0x71a141c6.
//
// Solidity: function quote((string,address,address,uint256,uint256,uint256,uint256,bytes)[] _inputs) returns(bool)
func (_IBridgeFeeQuote *IBridgeFeeQuoteSession) Quote(_inputs []IBridgeFeeQuoteQuoteInput) (*types.Transaction, error) {
	return _IBridgeFeeQuote.Contract.Quote(&_IBridgeFeeQuote.TransactOpts, _inputs)
}

// Quote is a paid mutator transaction binding the contract method 0x71a141c6.
//
// Solidity: function quote((string,address,address,uint256,uint256,uint256,uint256,bytes)[] _inputs) returns(bool)
func (_IBridgeFeeQuote *IBridgeFeeQuoteTransactorSession) Quote(_inputs []IBridgeFeeQuoteQuoteInput) (*types.Transaction, error) {
	return _IBridgeFeeQuote.Contract.Quote(&_IBridgeFeeQuote.TransactOpts, _inputs)
}
