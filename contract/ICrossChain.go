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

// ICrossChainMetaData contains all meta data concerning the ICrossChain contract.
var ICrossChainMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_txOrigin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"BridgeCallEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receipt\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"CrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"}],\"name\":\"ExecuteClaimEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_refund\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ICrossChainABI is the input ABI used to generate the binding from.
// Deprecated: Use ICrossChainMetaData.ABI instead.
var ICrossChainABI = ICrossChainMetaData.ABI

// ICrossChain is an auto generated Go binding around an Ethereum contract.
type ICrossChain struct {
	ICrossChainCaller     // Read-only binding to the contract
	ICrossChainTransactor // Write-only binding to the contract
	ICrossChainFilterer   // Log filterer for contract events
}

// ICrossChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type ICrossChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICrossChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ICrossChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICrossChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ICrossChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICrossChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ICrossChainSession struct {
	Contract     *ICrossChain      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ICrossChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ICrossChainCallerSession struct {
	Contract *ICrossChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ICrossChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ICrossChainTransactorSession struct {
	Contract     *ICrossChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ICrossChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type ICrossChainRaw struct {
	Contract *ICrossChain // Generic contract binding to access the raw methods on
}

// ICrossChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ICrossChainCallerRaw struct {
	Contract *ICrossChainCaller // Generic read-only contract binding to access the raw methods on
}

// ICrossChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ICrossChainTransactorRaw struct {
	Contract *ICrossChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICrossChain creates a new instance of ICrossChain, bound to a specific deployed contract.
func NewICrossChain(address common.Address, backend bind.ContractBackend) (*ICrossChain, error) {
	contract, err := bindICrossChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICrossChain{ICrossChainCaller: ICrossChainCaller{contract: contract}, ICrossChainTransactor: ICrossChainTransactor{contract: contract}, ICrossChainFilterer: ICrossChainFilterer{contract: contract}}, nil
}

// NewICrossChainCaller creates a new read-only instance of ICrossChain, bound to a specific deployed contract.
func NewICrossChainCaller(address common.Address, caller bind.ContractCaller) (*ICrossChainCaller, error) {
	contract, err := bindICrossChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICrossChainCaller{contract: contract}, nil
}

// NewICrossChainTransactor creates a new write-only instance of ICrossChain, bound to a specific deployed contract.
func NewICrossChainTransactor(address common.Address, transactor bind.ContractTransactor) (*ICrossChainTransactor, error) {
	contract, err := bindICrossChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICrossChainTransactor{contract: contract}, nil
}

// NewICrossChainFilterer creates a new log filterer instance of ICrossChain, bound to a specific deployed contract.
func NewICrossChainFilterer(address common.Address, filterer bind.ContractFilterer) (*ICrossChainFilterer, error) {
	contract, err := bindICrossChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICrossChainFilterer{contract: contract}, nil
}

// bindICrossChain binds a generic wrapper to an already deployed contract.
func bindICrossChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICrossChainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICrossChain *ICrossChainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICrossChain.Contract.ICrossChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICrossChain *ICrossChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICrossChain.Contract.ICrossChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICrossChain *ICrossChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICrossChain.Contract.ICrossChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICrossChain *ICrossChainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICrossChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICrossChain *ICrossChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICrossChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICrossChain *ICrossChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICrossChain.Contract.contract.Transact(opts, method, params...)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256 _amount)
func (_ICrossChain *ICrossChainCaller) BridgeCoinAmount(opts *bind.CallOpts, _token common.Address, _target [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ICrossChain.contract.Call(opts, &out, "bridgeCoinAmount", _token, _target)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256 _amount)
func (_ICrossChain *ICrossChainSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _ICrossChain.Contract.BridgeCoinAmount(&_ICrossChain.CallOpts, _token, _target)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256 _amount)
func (_ICrossChain *ICrossChainCallerSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _ICrossChain.Contract.BridgeCoinAmount(&_ICrossChain.CallOpts, _token, _target)
}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrossChain *ICrossChainCaller) HasOracle(opts *bind.CallOpts, _chain string, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _ICrossChain.contract.Call(opts, &out, "hasOracle", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrossChain *ICrossChainSession) HasOracle(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrossChain.Contract.HasOracle(&_ICrossChain.CallOpts, _chain, _externalAddress)
}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrossChain *ICrossChainCallerSession) HasOracle(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrossChain.Contract.HasOracle(&_ICrossChain.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrossChain *ICrossChainCaller) IsOracleOnline(opts *bind.CallOpts, _chain string, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _ICrossChain.contract.Call(opts, &out, "isOracleOnline", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrossChain *ICrossChainSession) IsOracleOnline(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrossChain.Contract.IsOracleOnline(&_ICrossChain.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrossChain *ICrossChainCallerSession) IsOracleOnline(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrossChain.Contract.IsOracleOnline(&_ICrossChain.CallOpts, _chain, _externalAddress)
}

// BridgeCall is a paid mutator transaction binding the contract method 0x851c42ee.
//
// Solidity: function bridgeCall(string _dstChain, address _refund, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _value, bytes _memo) payable returns(uint256 _eventNonce)
func (_ICrossChain *ICrossChainTransactor) BridgeCall(opts *bind.TransactOpts, _dstChain string, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _value *big.Int, _memo []byte) (*types.Transaction, error) {
	return _ICrossChain.contract.Transact(opts, "bridgeCall", _dstChain, _refund, _tokens, _amounts, _to, _data, _value, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0x851c42ee.
//
// Solidity: function bridgeCall(string _dstChain, address _refund, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _value, bytes _memo) payable returns(uint256 _eventNonce)
func (_ICrossChain *ICrossChainSession) BridgeCall(_dstChain string, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _value *big.Int, _memo []byte) (*types.Transaction, error) {
	return _ICrossChain.Contract.BridgeCall(&_ICrossChain.TransactOpts, _dstChain, _refund, _tokens, _amounts, _to, _data, _value, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0x851c42ee.
//
// Solidity: function bridgeCall(string _dstChain, address _refund, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _value, bytes _memo) payable returns(uint256 _eventNonce)
func (_ICrossChain *ICrossChainTransactorSession) BridgeCall(_dstChain string, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _value *big.Int, _memo []byte) (*types.Transaction, error) {
	return _ICrossChain.Contract.BridgeCall(&_ICrossChain.TransactOpts, _dstChain, _refund, _tokens, _amounts, _to, _data, _value, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool _result)
func (_ICrossChain *ICrossChainTransactor) CrossChain(opts *bind.TransactOpts, _token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _ICrossChain.contract.Transact(opts, "crossChain", _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool _result)
func (_ICrossChain *ICrossChainSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _ICrossChain.Contract.CrossChain(&_ICrossChain.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool _result)
func (_ICrossChain *ICrossChainTransactorSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _ICrossChain.Contract.CrossChain(&_ICrossChain.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_ICrossChain *ICrossChainTransactor) ExecuteClaim(opts *bind.TransactOpts, _chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _ICrossChain.contract.Transact(opts, "executeClaim", _chain, _eventNonce)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_ICrossChain *ICrossChainSession) ExecuteClaim(_chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _ICrossChain.Contract.ExecuteClaim(&_ICrossChain.TransactOpts, _chain, _eventNonce)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_ICrossChain *ICrossChainTransactorSession) ExecuteClaim(_chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _ICrossChain.Contract.ExecuteClaim(&_ICrossChain.TransactOpts, _chain, _eventNonce)
}

// ICrossChainBridgeCallEventIterator is returned from FilterBridgeCallEvent and is used to iterate over the raw logs and unpacked data for BridgeCallEvent events raised by the ICrossChain contract.
type ICrossChainBridgeCallEventIterator struct {
	Event *ICrossChainBridgeCallEvent // Event containing the contract specifics and raw log

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
func (it *ICrossChainBridgeCallEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICrossChainBridgeCallEvent)
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
		it.Event = new(ICrossChainBridgeCallEvent)
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
func (it *ICrossChainBridgeCallEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICrossChainBridgeCallEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICrossChainBridgeCallEvent represents a BridgeCallEvent event raised by the ICrossChain contract.
type ICrossChainBridgeCallEvent struct {
	Sender     common.Address
	Receiver   common.Address
	To         common.Address
	TxOrigin   common.Address
	Value      *big.Int
	EventNonce *big.Int
	DstChain   string
	Tokens     []common.Address
	Amounts    []*big.Int
	Data       []byte
	Memo       []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBridgeCallEvent is a free log retrieval operation binding the contract event 0x4a9b24da6150ef33e7c41038842b7c94fe89a4fff22dccb2c3fd79f0176062c6.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address _txOrigin, uint256 _value, uint256 _eventNonce, string _dstChain, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo)
func (_ICrossChain *ICrossChainFilterer) FilterBridgeCallEvent(opts *bind.FilterOpts, _sender []common.Address, _receiver []common.Address, _to []common.Address) (*ICrossChainBridgeCallEventIterator, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _ICrossChain.contract.FilterLogs(opts, "BridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &ICrossChainBridgeCallEventIterator{contract: _ICrossChain.contract, event: "BridgeCallEvent", logs: logs, sub: sub}, nil
}

// WatchBridgeCallEvent is a free log subscription operation binding the contract event 0x4a9b24da6150ef33e7c41038842b7c94fe89a4fff22dccb2c3fd79f0176062c6.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address _txOrigin, uint256 _value, uint256 _eventNonce, string _dstChain, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo)
func (_ICrossChain *ICrossChainFilterer) WatchBridgeCallEvent(opts *bind.WatchOpts, sink chan<- *ICrossChainBridgeCallEvent, _sender []common.Address, _receiver []common.Address, _to []common.Address) (event.Subscription, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}
	var _receiverRule []interface{}
	for _, _receiverItem := range _receiver {
		_receiverRule = append(_receiverRule, _receiverItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _ICrossChain.contract.WatchLogs(opts, "BridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICrossChainBridgeCallEvent)
				if err := _ICrossChain.contract.UnpackLog(event, "BridgeCallEvent", log); err != nil {
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

// ParseBridgeCallEvent is a log parse operation binding the contract event 0x4a9b24da6150ef33e7c41038842b7c94fe89a4fff22dccb2c3fd79f0176062c6.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address _txOrigin, uint256 _value, uint256 _eventNonce, string _dstChain, address[] _tokens, uint256[] _amounts, bytes _data, bytes _memo)
func (_ICrossChain *ICrossChainFilterer) ParseBridgeCallEvent(log types.Log) (*ICrossChainBridgeCallEvent, error) {
	event := new(ICrossChainBridgeCallEvent)
	if err := _ICrossChain.contract.UnpackLog(event, "BridgeCallEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICrossChainCrossChainIterator is returned from FilterCrossChain and is used to iterate over the raw logs and unpacked data for CrossChain events raised by the ICrossChain contract.
type ICrossChainCrossChainIterator struct {
	Event *ICrossChainCrossChain // Event containing the contract specifics and raw log

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
func (it *ICrossChainCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICrossChainCrossChain)
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
		it.Event = new(ICrossChainCrossChain)
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
func (it *ICrossChainCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICrossChainCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICrossChainCrossChain represents a CrossChain event raised by the ICrossChain contract.
type ICrossChainCrossChain struct {
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
func (_ICrossChain *ICrossChainFilterer) FilterCrossChain(opts *bind.FilterOpts, sender []common.Address, token []common.Address) (*ICrossChainCrossChainIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ICrossChain.contract.FilterLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &ICrossChainCrossChainIterator{contract: _ICrossChain.contract, event: "CrossChain", logs: logs, sub: sub}, nil
}

// WatchCrossChain is a free log subscription operation binding the contract event 0xb783df819ac99ca709650d67d9237a00b553c6ef941dceabeed6f4bc990d31ba.
//
// Solidity: event CrossChain(address indexed sender, address indexed token, string denom, string receipt, uint256 amount, uint256 fee, bytes32 target, string memo)
func (_ICrossChain *ICrossChainFilterer) WatchCrossChain(opts *bind.WatchOpts, sink chan<- *ICrossChainCrossChain, sender []common.Address, token []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ICrossChain.contract.WatchLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICrossChainCrossChain)
				if err := _ICrossChain.contract.UnpackLog(event, "CrossChain", log); err != nil {
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
func (_ICrossChain *ICrossChainFilterer) ParseCrossChain(log types.Log) (*ICrossChainCrossChain, error) {
	event := new(ICrossChainCrossChain)
	if err := _ICrossChain.contract.UnpackLog(event, "CrossChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICrossChainExecuteClaimEventIterator is returned from FilterExecuteClaimEvent and is used to iterate over the raw logs and unpacked data for ExecuteClaimEvent events raised by the ICrossChain contract.
type ICrossChainExecuteClaimEventIterator struct {
	Event *ICrossChainExecuteClaimEvent // Event containing the contract specifics and raw log

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
func (it *ICrossChainExecuteClaimEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICrossChainExecuteClaimEvent)
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
		it.Event = new(ICrossChainExecuteClaimEvent)
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
func (it *ICrossChainExecuteClaimEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICrossChainExecuteClaimEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICrossChainExecuteClaimEvent represents a ExecuteClaimEvent event raised by the ICrossChain contract.
type ICrossChainExecuteClaimEvent struct {
	Sender     common.Address
	EventNonce *big.Int
	Chain      string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterExecuteClaimEvent is a free log retrieval operation binding the contract event 0xa45a8d344c26216c8d81958a3688ec20b5f2e5af820e03433537687e94667a78.
//
// Solidity: event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain)
func (_ICrossChain *ICrossChainFilterer) FilterExecuteClaimEvent(opts *bind.FilterOpts, _sender []common.Address) (*ICrossChainExecuteClaimEventIterator, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}

	logs, sub, err := _ICrossChain.contract.FilterLogs(opts, "ExecuteClaimEvent", _senderRule)
	if err != nil {
		return nil, err
	}
	return &ICrossChainExecuteClaimEventIterator{contract: _ICrossChain.contract, event: "ExecuteClaimEvent", logs: logs, sub: sub}, nil
}

// WatchExecuteClaimEvent is a free log subscription operation binding the contract event 0xa45a8d344c26216c8d81958a3688ec20b5f2e5af820e03433537687e94667a78.
//
// Solidity: event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain)
func (_ICrossChain *ICrossChainFilterer) WatchExecuteClaimEvent(opts *bind.WatchOpts, sink chan<- *ICrossChainExecuteClaimEvent, _sender []common.Address) (event.Subscription, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}

	logs, sub, err := _ICrossChain.contract.WatchLogs(opts, "ExecuteClaimEvent", _senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICrossChainExecuteClaimEvent)
				if err := _ICrossChain.contract.UnpackLog(event, "ExecuteClaimEvent", log); err != nil {
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

// ParseExecuteClaimEvent is a log parse operation binding the contract event 0xa45a8d344c26216c8d81958a3688ec20b5f2e5af820e03433537687e94667a78.
//
// Solidity: event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain)
func (_ICrossChain *ICrossChainFilterer) ParseExecuteClaimEvent(log types.Log) (*ICrossChainExecuteClaimEvent, error) {
	event := new(ICrossChainExecuteClaimEvent)
	if err := _ICrossChain.contract.UnpackLog(event, "ExecuteClaimEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
