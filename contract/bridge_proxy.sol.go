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

// BridgeProxyMetaData contains all meta data concerning the BridgeProxy contract.
var BridgeProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610284806100206000396000f3fe6080604052600436106100225760003560e01c806319ab453c1461003957610031565b366100315761002f610059565b005b61002f610059565b34801561004557600080fd5b5061002f61005436600461021e565b61006b565b6100696100646100d0565b610108565b565b600061009e7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b0316146100c45760405162dc149f60e41b815260040160405180910390fd5b6100cd8161012c565b50565b60006101037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e808015610127573d6000f35b3d6000fd5b6101358161016c565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6001600160a01b0381163b6101dd5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840160405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b60006020828403121561023057600080fd5b81356001600160a01b038116811461024757600080fd5b939250505056fea2646970667358221220f5fd2c6b493d2d8c8fccb10ef0e5ed3a7f4f1914bae22833db4f76dbce6da4fe64736f6c634300080a0033",
}

// BridgeProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeProxyMetaData.ABI instead.
var BridgeProxyABI = BridgeProxyMetaData.ABI

// BridgeProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeProxyMetaData.Bin instead.
var BridgeProxyBin = BridgeProxyMetaData.Bin

// DeployBridgeProxy deploys a new Ethereum contract, binding an instance of BridgeProxy to it.
func DeployBridgeProxy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeProxy, error) {
	parsed, err := BridgeProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeProxyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeProxy{BridgeProxyCaller: BridgeProxyCaller{contract: contract}, BridgeProxyTransactor: BridgeProxyTransactor{contract: contract}, BridgeProxyFilterer: BridgeProxyFilterer{contract: contract}}, nil
}

// BridgeProxy is an auto generated Go binding around an Ethereum contract.
type BridgeProxy struct {
	BridgeProxyCaller     // Read-only binding to the contract
	BridgeProxyTransactor // Write-only binding to the contract
	BridgeProxyFilterer   // Log filterer for contract events
}

// BridgeProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeProxySession struct {
	Contract     *BridgeProxy      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeProxyCallerSession struct {
	Contract *BridgeProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BridgeProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeProxyTransactorSession struct {
	Contract     *BridgeProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BridgeProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeProxyRaw struct {
	Contract *BridgeProxy // Generic contract binding to access the raw methods on
}

// BridgeProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeProxyCallerRaw struct {
	Contract *BridgeProxyCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeProxyTransactorRaw struct {
	Contract *BridgeProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeProxy creates a new instance of BridgeProxy, bound to a specific deployed contract.
func NewBridgeProxy(address common.Address, backend bind.ContractBackend) (*BridgeProxy, error) {
	contract, err := bindBridgeProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeProxy{BridgeProxyCaller: BridgeProxyCaller{contract: contract}, BridgeProxyTransactor: BridgeProxyTransactor{contract: contract}, BridgeProxyFilterer: BridgeProxyFilterer{contract: contract}}, nil
}

// NewBridgeProxyCaller creates a new read-only instance of BridgeProxy, bound to a specific deployed contract.
func NewBridgeProxyCaller(address common.Address, caller bind.ContractCaller) (*BridgeProxyCaller, error) {
	contract, err := bindBridgeProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeProxyCaller{contract: contract}, nil
}

// NewBridgeProxyTransactor creates a new write-only instance of BridgeProxy, bound to a specific deployed contract.
func NewBridgeProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeProxyTransactor, error) {
	contract, err := bindBridgeProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeProxyTransactor{contract: contract}, nil
}

// NewBridgeProxyFilterer creates a new log filterer instance of BridgeProxy, bound to a specific deployed contract.
func NewBridgeProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeProxyFilterer, error) {
	contract, err := bindBridgeProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeProxyFilterer{contract: contract}, nil
}

// bindBridgeProxy binds a generic wrapper to an already deployed contract.
func bindBridgeProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeProxy *BridgeProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeProxy.Contract.BridgeProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeProxy *BridgeProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeProxy.Contract.BridgeProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeProxy *BridgeProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeProxy.Contract.BridgeProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeProxy *BridgeProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeProxy *BridgeProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeProxy *BridgeProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeProxy.Contract.contract.Transact(opts, method, params...)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _logic) returns()
func (_BridgeProxy *BridgeProxyTransactor) Init(opts *bind.TransactOpts, _logic common.Address) (*types.Transaction, error) {
	return _BridgeProxy.contract.Transact(opts, "init", _logic)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _logic) returns()
func (_BridgeProxy *BridgeProxySession) Init(_logic common.Address) (*types.Transaction, error) {
	return _BridgeProxy.Contract.Init(&_BridgeProxy.TransactOpts, _logic)
}

// Init is a paid mutator transaction binding the contract method 0x19ab453c.
//
// Solidity: function init(address _logic) returns()
func (_BridgeProxy *BridgeProxyTransactorSession) Init(_logic common.Address) (*types.Transaction, error) {
	return _BridgeProxy.Contract.Init(&_BridgeProxy.TransactOpts, _logic)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BridgeProxy *BridgeProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _BridgeProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BridgeProxy *BridgeProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BridgeProxy.Contract.Fallback(&_BridgeProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BridgeProxy *BridgeProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BridgeProxy.Contract.Fallback(&_BridgeProxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeProxy *BridgeProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeProxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeProxy *BridgeProxySession) Receive() (*types.Transaction, error) {
	return _BridgeProxy.Contract.Receive(&_BridgeProxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BridgeProxy *BridgeProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _BridgeProxy.Contract.Receive(&_BridgeProxy.TransactOpts)
}

// BridgeProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the BridgeProxy contract.
type BridgeProxyAdminChangedIterator struct {
	Event *BridgeProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *BridgeProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeProxyAdminChanged)
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
		it.Event = new(BridgeProxyAdminChanged)
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
func (it *BridgeProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeProxyAdminChanged represents a AdminChanged event raised by the BridgeProxy contract.
type BridgeProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BridgeProxy *BridgeProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*BridgeProxyAdminChangedIterator, error) {

	logs, sub, err := _BridgeProxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &BridgeProxyAdminChangedIterator{contract: _BridgeProxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BridgeProxy *BridgeProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *BridgeProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _BridgeProxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeProxyAdminChanged)
				if err := _BridgeProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BridgeProxy *BridgeProxyFilterer) ParseAdminChanged(log types.Log) (*BridgeProxyAdminChanged, error) {
	event := new(BridgeProxyAdminChanged)
	if err := _BridgeProxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeProxyBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the BridgeProxy contract.
type BridgeProxyBeaconUpgradedIterator struct {
	Event *BridgeProxyBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *BridgeProxyBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeProxyBeaconUpgraded)
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
		it.Event = new(BridgeProxyBeaconUpgraded)
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
func (it *BridgeProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeProxyBeaconUpgraded represents a BeaconUpgraded event raised by the BridgeProxy contract.
type BridgeProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BridgeProxy *BridgeProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*BridgeProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _BridgeProxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &BridgeProxyBeaconUpgradedIterator{contract: _BridgeProxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BridgeProxy *BridgeProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *BridgeProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _BridgeProxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeProxyBeaconUpgraded)
				if err := _BridgeProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BridgeProxy *BridgeProxyFilterer) ParseBeaconUpgraded(log types.Log) (*BridgeProxyBeaconUpgraded, error) {
	event := new(BridgeProxyBeaconUpgraded)
	if err := _BridgeProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the BridgeProxy contract.
type BridgeProxyUpgradedIterator struct {
	Event *BridgeProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *BridgeProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeProxyUpgraded)
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
		it.Event = new(BridgeProxyUpgraded)
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
func (it *BridgeProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeProxyUpgraded represents a Upgraded event raised by the BridgeProxy contract.
type BridgeProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BridgeProxy *BridgeProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BridgeProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BridgeProxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BridgeProxyUpgradedIterator{contract: _BridgeProxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BridgeProxy *BridgeProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BridgeProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BridgeProxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeProxyUpgraded)
				if err := _BridgeProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BridgeProxy *BridgeProxyFilterer) ParseUpgraded(log types.Log) (*BridgeProxyUpgraded, error) {
	event := new(BridgeProxyUpgraded)
	if err := _BridgeProxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
