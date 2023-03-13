// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package crosschain_test

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CrosschainTestABI is the input ABI used to generate the binding from.
const CrosschainTestABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txid\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// CrosschainTestBin is the compiled bytecode used for deploying new contracts.
var CrosschainTestBin = "0x608060405234801561001057600080fd5b50610986806100206000396000f3fe6080604052600436106100295760003560e01c80630b56c1901461002e578063160d7c7314610062575b600080fd5b34801561003a57600080fd5b5061004e61004936600461067d565b610075565b604051901515815260200160405180910390f35b61004e6100703660046106c2565b61008a565b600061008183836101d5565b90505b92915050565b60006001600160a01b038716156101bc576001600160a01b0387166323b872dd33306100b6888a610760565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af115801561010a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061012e9190610781565b506001600160a01b03871663095ea7b361100461014b8789610760565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015610196573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101ba9190610781565b505b6101ca87878787878761028c565b979650505050505050565b600080806110046101e68686610481565b6040516101f391906107c7565b6000604051808303816000865af19150503d8060008114610230576040519150601f19603f3d011682016040523d82523d6000602084013e610235565b606091505b509150915061027a82826040518060400160405280601e81526020017f63616e63656c2073656e6420746f2065787465726e616c206661696c656400008152506104c8565b61028381610547565b95945050505050565b60006001600160a01b0387161561037557604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e90604401602060405180830381865afa1580156102ec573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031091906107e3565b905061031c8587610760565b811461036f5760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b506103cd565b61037f8486610760565b34146103cd5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b206665656044820152606401610366565b600080611004346103e28b8b8b8b8b8b610565565b6040516103ef91906107c7565b60006040518083038185875af1925050503d806000811461042c576040519150601f19603f3d011682016040523d82523d6000602084013e610431565b606091505b509150915061046b82826040518060400160405280601281526020017118dc9bdcdccb58da185a5b8819985a5b195960721b8152506104c8565b61047481610547565b9998505050505050505050565b60608282604051602401610496929190610828565b60408051601f198184030181529190526020810180516001600160e01b031663eeb3593d60e01b179052905092915050565b82610542576000828060200190518101906104e3919061084a565b9050600182511015610509578060405162461bcd60e51b815260040161036691906108b8565b818160405160200161051c9291906108cb565b60408051601f198184030181529082905262461bcd60e51b8252610366916004016108b8565b505050565b6000808280602001905181019061055e9190610781565b9392505050565b606086868686868660405160240161058296959493929190610908565b60408051601f198184030181529190526020810180516001600160e01b031663160d7c7360e01b17905290509695505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156105f7576105f76105b8565b604052919050565b600067ffffffffffffffff821115610619576106196105b8565b50601f01601f191660200190565b600082601f83011261063857600080fd5b813561064b610646826105ff565b6105ce565b81815284602083860101111561066057600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561069057600080fd5b823567ffffffffffffffff8111156106a757600080fd5b6106b385828601610627565b95602094909401359450505050565b60008060008060008060c087890312156106db57600080fd5b86356001600160a01b03811681146106f257600080fd5b9550602087013567ffffffffffffffff8082111561070f57600080fd5b61071b8a838b01610627565b965060408901359550606089013594506080890135935060a089013591508082111561074657600080fd5b5061075389828a01610627565b9150509295509295509295565b8082018082111561008457634e487b7160e01b600052601160045260246000fd5b60006020828403121561079357600080fd5b8151801515811461055e57600080fd5b60005b838110156107be5781810151838201526020016107a6565b50506000910152565b600082516107d98184602087016107a3565b9190910192915050565b6000602082840312156107f557600080fd5b5051919050565b600081518084526108148160208601602086016107a3565b601f01601f19169290920160200192915050565b60408152600061083b60408301856107fc565b90508260208301529392505050565b60006020828403121561085c57600080fd5b815167ffffffffffffffff81111561087357600080fd5b8201601f8101841361088457600080fd5b8051610892610646826105ff565b8181528560208385010111156108a757600080fd5b6102838260208301602086016107a3565b60208152600061008160208301846107fc565b600083516108dd8184602088016107a3565b6101d160f51b90830190815283516108fc8160028401602088016107a3565b01600201949350505050565b6001600160a01b038716815260c06020820181905260009061092c908301886107fc565b86604084015285606084015284608084015282810360a084015261047481856107fc56fea2646970667358221220ed776e583be707d2c7ace1659d741142a5dc508bc50375ba1c3152a551ada14664736f6c63430008130033"

// DeployCrosschainTest deploys a new Ethereum contract, binding an instance of CrosschainTest to it.
func DeployCrosschainTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrosschainTest, error) {
	parsed, err := abi.JSON(strings.NewReader(CrosschainTestABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CrosschainTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrosschainTest{CrosschainTestCaller: CrosschainTestCaller{contract: contract}, CrosschainTestTransactor: CrosschainTestTransactor{contract: contract}, CrosschainTestFilterer: CrosschainTestFilterer{contract: contract}}, nil
}

// CrosschainTest is an auto generated Go binding around an Ethereum contract.
type CrosschainTest struct {
	CrosschainTestCaller     // Read-only binding to the contract
	CrosschainTestTransactor // Write-only binding to the contract
	CrosschainTestFilterer   // Log filterer for contract events
}

// CrosschainTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrosschainTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrosschainTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrosschainTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrosschainTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrosschainTestSession struct {
	Contract     *CrosschainTest   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrosschainTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrosschainTestCallerSession struct {
	Contract *CrosschainTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CrosschainTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrosschainTestTransactorSession struct {
	Contract     *CrosschainTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CrosschainTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrosschainTestRaw struct {
	Contract *CrosschainTest // Generic contract binding to access the raw methods on
}

// CrosschainTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrosschainTestCallerRaw struct {
	Contract *CrosschainTestCaller // Generic read-only contract binding to access the raw methods on
}

// CrosschainTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrosschainTestTransactorRaw struct {
	Contract *CrosschainTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrosschainTest creates a new instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTest(address common.Address, backend bind.ContractBackend) (*CrosschainTest, error) {
	contract, err := bindCrosschainTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrosschainTest{CrosschainTestCaller: CrosschainTestCaller{contract: contract}, CrosschainTestTransactor: CrosschainTestTransactor{contract: contract}, CrosschainTestFilterer: CrosschainTestFilterer{contract: contract}}, nil
}

// NewCrosschainTestCaller creates a new read-only instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTestCaller(address common.Address, caller bind.ContractCaller) (*CrosschainTestCaller, error) {
	contract, err := bindCrosschainTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestCaller{contract: contract}, nil
}

// NewCrosschainTestTransactor creates a new write-only instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTestTransactor(address common.Address, transactor bind.ContractTransactor) (*CrosschainTestTransactor, error) {
	contract, err := bindCrosschainTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestTransactor{contract: contract}, nil
}

// NewCrosschainTestFilterer creates a new log filterer instance of CrosschainTest, bound to a specific deployed contract.
func NewCrosschainTestFilterer(address common.Address, filterer bind.ContractFilterer) (*CrosschainTestFilterer, error) {
	contract, err := bindCrosschainTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestFilterer{contract: contract}, nil
}

// bindCrosschainTest binds a generic wrapper to an already deployed contract.
func bindCrosschainTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CrosschainTestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainTest *CrosschainTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainTest.Contract.CrosschainTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainTest *CrosschainTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrosschainTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainTest *CrosschainTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrosschainTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrosschainTest *CrosschainTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrosschainTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrosschainTest *CrosschainTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrosschainTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrosschainTest *CrosschainTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrosschainTest.Contract.contract.Transact(opts, method, params...)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txid) returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) CancelSendToExternal(opts *bind.TransactOpts, _chain string, _txid *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "cancelSendToExternal", _chain, _txid)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txid) returns(bool)
func (_CrosschainTest *CrosschainTestSession) CancelSendToExternal(_chain string, _txid *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CancelSendToExternal(&_CrosschainTest.TransactOpts, _chain, _txid)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txid) returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) CancelSendToExternal(_chain string, _txid *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CancelSendToExternal(&_CrosschainTest.TransactOpts, _chain, _txid)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) CrossChain(opts *bind.TransactOpts, _token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "crossChain", _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrosschainTest *CrosschainTestSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrossChain(&_CrosschainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CrossChain(&_CrosschainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}
