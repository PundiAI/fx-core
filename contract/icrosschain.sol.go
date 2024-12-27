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

// ICrosschainMetaData contains all meta data concerning the ICrosschain contract.
var ICrosschainMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_txOrigin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"BridgeCallEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receipt\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"CrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_errReason\",\"type\":\"string\"}],\"name\":\"ExecuteClaimEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_refund\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ICrosschainABI is the input ABI used to generate the binding from.
// Deprecated: Use ICrosschainMetaData.ABI instead.
var ICrosschainABI = ICrosschainMetaData.ABI

// ICrosschain is an auto generated Go binding around an Ethereum contract.
type ICrosschain struct {
	ICrosschainCaller     // Read-only binding to the contract
	ICrosschainTransactor // Write-only binding to the contract
	ICrosschainFilterer   // Log filterer for contract events
}

// ICrosschainCaller is an auto generated read-only Go binding around an Ethereum contract.
type ICrosschainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICrosschainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ICrosschainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICrosschainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ICrosschainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICrosschainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ICrosschainSession struct {
	Contract     *ICrosschain      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ICrosschainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ICrosschainCallerSession struct {
	Contract *ICrosschainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ICrosschainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ICrosschainTransactorSession struct {
	Contract     *ICrosschainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ICrosschainRaw is an auto generated low-level Go binding around an Ethereum contract.
type ICrosschainRaw struct {
	Contract *ICrosschain // Generic contract binding to access the raw methods on
}

// ICrosschainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ICrosschainCallerRaw struct {
	Contract *ICrosschainCaller // Generic read-only contract binding to access the raw methods on
}

// ICrosschainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ICrosschainTransactorRaw struct {
	Contract *ICrosschainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICrosschain creates a new instance of ICrosschain, bound to a specific deployed contract.
func NewICrosschain(address common.Address, backend bind.ContractBackend) (*ICrosschain, error) {
	contract, err := bindICrosschain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICrosschain{ICrosschainCaller: ICrosschainCaller{contract: contract}, ICrosschainTransactor: ICrosschainTransactor{contract: contract}, ICrosschainFilterer: ICrosschainFilterer{contract: contract}}, nil
}

// NewICrosschainCaller creates a new read-only instance of ICrosschain, bound to a specific deployed contract.
func NewICrosschainCaller(address common.Address, caller bind.ContractCaller) (*ICrosschainCaller, error) {
	contract, err := bindICrosschain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICrosschainCaller{contract: contract}, nil
}

// NewICrosschainTransactor creates a new write-only instance of ICrosschain, bound to a specific deployed contract.
func NewICrosschainTransactor(address common.Address, transactor bind.ContractTransactor) (*ICrosschainTransactor, error) {
	contract, err := bindICrosschain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICrosschainTransactor{contract: contract}, nil
}

// NewICrosschainFilterer creates a new log filterer instance of ICrosschain, bound to a specific deployed contract.
func NewICrosschainFilterer(address common.Address, filterer bind.ContractFilterer) (*ICrosschainFilterer, error) {
	contract, err := bindICrosschain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICrosschainFilterer{contract: contract}, nil
}

// bindICrosschain binds a generic wrapper to an already deployed contract.
func bindICrosschain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICrosschainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICrosschain *ICrosschainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICrosschain.Contract.ICrosschainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICrosschain *ICrosschainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICrosschain.Contract.ICrosschainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICrosschain *ICrosschainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICrosschain.Contract.ICrosschainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICrosschain *ICrosschainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICrosschain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICrosschain *ICrosschainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICrosschain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICrosschain *ICrosschainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICrosschain.Contract.contract.Transact(opts, method, params...)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256 _amount)
func (_ICrosschain *ICrosschainCaller) BridgeCoinAmount(opts *bind.CallOpts, _token common.Address, _target [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ICrosschain.contract.Call(opts, &out, "bridgeCoinAmount", _token, _target)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256 _amount)
func (_ICrosschain *ICrosschainSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _ICrosschain.Contract.BridgeCoinAmount(&_ICrosschain.CallOpts, _token, _target)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256 _amount)
func (_ICrosschain *ICrosschainCallerSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _ICrosschain.Contract.BridgeCoinAmount(&_ICrosschain.CallOpts, _token, _target)
}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrosschain *ICrosschainCaller) HasOracle(opts *bind.CallOpts, _chain string, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _ICrosschain.contract.Call(opts, &out, "hasOracle", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrosschain *ICrosschainSession) HasOracle(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrosschain.Contract.HasOracle(&_ICrosschain.CallOpts, _chain, _externalAddress)
}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrosschain *ICrosschainCallerSession) HasOracle(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrosschain.Contract.HasOracle(&_ICrosschain.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrosschain *ICrosschainCaller) IsOracleOnline(opts *bind.CallOpts, _chain string, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _ICrosschain.contract.Call(opts, &out, "isOracleOnline", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrosschain *ICrosschainSession) IsOracleOnline(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrosschain.Contract.IsOracleOnline(&_ICrosschain.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_ICrosschain *ICrosschainCallerSession) IsOracleOnline(_chain string, _externalAddress common.Address) (bool, error) {
	return _ICrosschain.Contract.IsOracleOnline(&_ICrosschain.CallOpts, _chain, _externalAddress)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _refund, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) payable returns(uint256 _eventNonce)
func (_ICrosschain *ICrosschainTransactor) BridgeCall(opts *bind.TransactOpts, _dstChain string, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _ICrosschain.contract.Transact(opts, "bridgeCall", _dstChain, _refund, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _refund, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) payable returns(uint256 _eventNonce)
func (_ICrosschain *ICrosschainSession) BridgeCall(_dstChain string, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _ICrosschain.Contract.BridgeCall(&_ICrosschain.TransactOpts, _dstChain, _refund, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _refund, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) payable returns(uint256 _eventNonce)
func (_ICrosschain *ICrosschainTransactorSession) BridgeCall(_dstChain string, _refund common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _ICrosschain.Contract.BridgeCall(&_ICrosschain.TransactOpts, _dstChain, _refund, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool _result)
func (_ICrosschain *ICrosschainTransactor) CrossChain(opts *bind.TransactOpts, _token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _ICrosschain.contract.Transact(opts, "crossChain", _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool _result)
func (_ICrosschain *ICrosschainSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _ICrosschain.Contract.CrossChain(&_ICrosschain.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool _result)
func (_ICrosschain *ICrosschainTransactorSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _ICrosschain.Contract.CrossChain(&_ICrosschain.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_ICrosschain *ICrosschainTransactor) ExecuteClaim(opts *bind.TransactOpts, _chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _ICrosschain.contract.Transact(opts, "executeClaim", _chain, _eventNonce)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_ICrosschain *ICrosschainSession) ExecuteClaim(_chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _ICrosschain.Contract.ExecuteClaim(&_ICrosschain.TransactOpts, _chain, _eventNonce)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_ICrosschain *ICrosschainTransactorSession) ExecuteClaim(_chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _ICrosschain.Contract.ExecuteClaim(&_ICrosschain.TransactOpts, _chain, _eventNonce)
}

// ICrosschainBridgeCallEventIterator is returned from FilterBridgeCallEvent and is used to iterate over the raw logs and unpacked data for BridgeCallEvent events raised by the ICrosschain contract.
type ICrosschainBridgeCallEventIterator struct {
	Event *ICrosschainBridgeCallEvent // Event containing the contract specifics and raw log

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
func (it *ICrosschainBridgeCallEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICrosschainBridgeCallEvent)
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
		it.Event = new(ICrosschainBridgeCallEvent)
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
func (it *ICrosschainBridgeCallEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICrosschainBridgeCallEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICrosschainBridgeCallEvent represents a BridgeCallEvent event raised by the ICrosschain contract.
type ICrosschainBridgeCallEvent struct {
	Sender     common.Address
	Receiver   common.Address
	To         common.Address
	TxOrigin   common.Address
	EventNonce *big.Int
	DstChain   string
	Tokens     []common.Address
	Amounts    []*big.Int
	Data       []byte
	QuoteId    *big.Int
	GasLimit   *big.Int
	Memo       []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBridgeCallEvent is a free log retrieval operation binding the contract event 0xcaa0e5b7ba998f542b3804184a5d30836451c57f6d1f031c466272e188f4a70f.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address _txOrigin, uint256 _eventNonce, string _dstChain, address[] _tokens, uint256[] _amounts, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo)
func (_ICrosschain *ICrosschainFilterer) FilterBridgeCallEvent(opts *bind.FilterOpts, _sender []common.Address, _receiver []common.Address, _to []common.Address) (*ICrosschainBridgeCallEventIterator, error) {

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

	logs, sub, err := _ICrosschain.contract.FilterLogs(opts, "BridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &ICrosschainBridgeCallEventIterator{contract: _ICrosschain.contract, event: "BridgeCallEvent", logs: logs, sub: sub}, nil
}

// WatchBridgeCallEvent is a free log subscription operation binding the contract event 0xcaa0e5b7ba998f542b3804184a5d30836451c57f6d1f031c466272e188f4a70f.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address _txOrigin, uint256 _eventNonce, string _dstChain, address[] _tokens, uint256[] _amounts, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo)
func (_ICrosschain *ICrosschainFilterer) WatchBridgeCallEvent(opts *bind.WatchOpts, sink chan<- *ICrosschainBridgeCallEvent, _sender []common.Address, _receiver []common.Address, _to []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _ICrosschain.contract.WatchLogs(opts, "BridgeCallEvent", _senderRule, _receiverRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICrosschainBridgeCallEvent)
				if err := _ICrosschain.contract.UnpackLog(event, "BridgeCallEvent", log); err != nil {
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

// ParseBridgeCallEvent is a log parse operation binding the contract event 0xcaa0e5b7ba998f542b3804184a5d30836451c57f6d1f031c466272e188f4a70f.
//
// Solidity: event BridgeCallEvent(address indexed _sender, address indexed _receiver, address indexed _to, address _txOrigin, uint256 _eventNonce, string _dstChain, address[] _tokens, uint256[] _amounts, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo)
func (_ICrosschain *ICrosschainFilterer) ParseBridgeCallEvent(log types.Log) (*ICrosschainBridgeCallEvent, error) {
	event := new(ICrosschainBridgeCallEvent)
	if err := _ICrosschain.contract.UnpackLog(event, "BridgeCallEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICrosschainCrossChainIterator is returned from FilterCrossChain and is used to iterate over the raw logs and unpacked data for CrossChain events raised by the ICrosschain contract.
type ICrosschainCrossChainIterator struct {
	Event *ICrosschainCrossChain // Event containing the contract specifics and raw log

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
func (it *ICrosschainCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICrosschainCrossChain)
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
		it.Event = new(ICrosschainCrossChain)
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
func (it *ICrosschainCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICrosschainCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICrosschainCrossChain represents a CrossChain event raised by the ICrosschain contract.
type ICrosschainCrossChain struct {
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
func (_ICrosschain *ICrosschainFilterer) FilterCrossChain(opts *bind.FilterOpts, sender []common.Address, token []common.Address) (*ICrosschainCrossChainIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ICrosschain.contract.FilterLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &ICrosschainCrossChainIterator{contract: _ICrosschain.contract, event: "CrossChain", logs: logs, sub: sub}, nil
}

// WatchCrossChain is a free log subscription operation binding the contract event 0xb783df819ac99ca709650d67d9237a00b553c6ef941dceabeed6f4bc990d31ba.
//
// Solidity: event CrossChain(address indexed sender, address indexed token, string denom, string receipt, uint256 amount, uint256 fee, bytes32 target, string memo)
func (_ICrosschain *ICrosschainFilterer) WatchCrossChain(opts *bind.WatchOpts, sink chan<- *ICrosschainCrossChain, sender []common.Address, token []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _ICrosschain.contract.WatchLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICrosschainCrossChain)
				if err := _ICrosschain.contract.UnpackLog(event, "CrossChain", log); err != nil {
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
func (_ICrosschain *ICrosschainFilterer) ParseCrossChain(log types.Log) (*ICrosschainCrossChain, error) {
	event := new(ICrosschainCrossChain)
	if err := _ICrosschain.contract.UnpackLog(event, "CrossChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICrosschainExecuteClaimEventIterator is returned from FilterExecuteClaimEvent and is used to iterate over the raw logs and unpacked data for ExecuteClaimEvent events raised by the ICrosschain contract.
type ICrosschainExecuteClaimEventIterator struct {
	Event *ICrosschainExecuteClaimEvent // Event containing the contract specifics and raw log

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
func (it *ICrosschainExecuteClaimEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICrosschainExecuteClaimEvent)
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
		it.Event = new(ICrosschainExecuteClaimEvent)
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
func (it *ICrosschainExecuteClaimEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICrosschainExecuteClaimEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICrosschainExecuteClaimEvent represents a ExecuteClaimEvent event raised by the ICrosschain contract.
type ICrosschainExecuteClaimEvent struct {
	Sender     common.Address
	EventNonce *big.Int
	Chain      string
	ErrReason  string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterExecuteClaimEvent is a free log retrieval operation binding the contract event 0x67ddf3796d30bb96cc50ccd9d322ab53317f5c6cac5860f3137894ee70ed0053.
//
// Solidity: event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain, string _errReason)
func (_ICrosschain *ICrosschainFilterer) FilterExecuteClaimEvent(opts *bind.FilterOpts, _sender []common.Address) (*ICrosschainExecuteClaimEventIterator, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}

	logs, sub, err := _ICrosschain.contract.FilterLogs(opts, "ExecuteClaimEvent", _senderRule)
	if err != nil {
		return nil, err
	}
	return &ICrosschainExecuteClaimEventIterator{contract: _ICrosschain.contract, event: "ExecuteClaimEvent", logs: logs, sub: sub}, nil
}

// WatchExecuteClaimEvent is a free log subscription operation binding the contract event 0x67ddf3796d30bb96cc50ccd9d322ab53317f5c6cac5860f3137894ee70ed0053.
//
// Solidity: event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain, string _errReason)
func (_ICrosschain *ICrosschainFilterer) WatchExecuteClaimEvent(opts *bind.WatchOpts, sink chan<- *ICrosschainExecuteClaimEvent, _sender []common.Address) (event.Subscription, error) {

	var _senderRule []interface{}
	for _, _senderItem := range _sender {
		_senderRule = append(_senderRule, _senderItem)
	}

	logs, sub, err := _ICrosschain.contract.WatchLogs(opts, "ExecuteClaimEvent", _senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICrosschainExecuteClaimEvent)
				if err := _ICrosschain.contract.UnpackLog(event, "ExecuteClaimEvent", log); err != nil {
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

// ParseExecuteClaimEvent is a log parse operation binding the contract event 0x67ddf3796d30bb96cc50ccd9d322ab53317f5c6cac5860f3137894ee70ed0053.
//
// Solidity: event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain, string _errReason)
func (_ICrosschain *ICrosschainFilterer) ParseExecuteClaimEvent(log types.Log) (*ICrosschainExecuteClaimEvent, error) {
	event := new(ICrosschainExecuteClaimEvent)
	if err := _ICrosschain.contract.UnpackLog(event, "ExecuteClaimEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
