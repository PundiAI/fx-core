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
const CrosschainTestABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txid\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txid\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"increaseBridgeFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// CrosschainTestBin is the compiled bytecode used for deploying new contracts.
var CrosschainTestBin = "0x608060405234801561001057600080fd5b50610b54806100206000396000f3fe6080604052600436106100345760003560e01c80630b56c19014610039578063160d7c731461006d578063c79a6b7b14610080575b600080fd5b34801561004557600080fd5b506100596100543660046107a6565b610093565b604051901515815260200160405180910390f35b61005961007b366004610807565b6100a8565b61005961008e366004610897565b6101f3565b600061009f838361020a565b90505b92915050565b60006001600160a01b038716156101da576001600160a01b0387166323b872dd33306100d4888a6108f6565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af1158015610128573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061014c9190610917565b506001600160a01b03871663095ea7b361100461016987896108f6565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af11580156101b4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101d89190610917565b505b6101e88787878787876102b8565b979650505050505050565b6000610201858585856104ad565b95945050505050565b6000808061100461021b868661055d565b604051610228919061095d565b6000604051808303816000865af19150503d8060008114610265576040519150601f19603f3d011682016040523d82523d6000602084013e61026a565b606091505b50915091506102af82826040518060400160405280601e81526020017f63616e63656c2073656e6420746f2065787465726e616c206661696c656400008152506105a4565b61020181610623565b60006001600160a01b038716156103a157604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e90604401602060405180830381865afa158015610318573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061033c9190610979565b905061034885876108f6565b811461039b5760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b506103f9565b6103ab84866108f6565b34146103f95760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b206665656044820152606401610392565b6000806110043461040e8b8b8b8b8b8b610641565b60405161041b919061095d565b60006040518083038185875af1925050503d8060008114610458576040519150601f19603f3d011682016040523d82523d6000602084013e61045d565b606091505b509150915061049782826040518060400160405280601281526020017118dc9bdcdccb58da185a5b8819985a5b195960721b8152506105a4565b6104a081610623565b9998505050505050505050565b600080806110046104c088888888610694565b6040516104cd919061095d565b6000604051808303816000865af19150503d806000811461050a576040519150601f19603f3d011682016040523d82523d6000602084013e61050f565b606091505b509150915061055482826040518060400160405280601a81526020017f696e6372656173652062726964676520666565206661696c65640000000000008152506105a4565b6101e881610623565b606082826040516024016105729291906109be565b60408051601f198184030181529190526020810180516001600160e01b031663eeb3593d60e01b179052905092915050565b8261061e576000828060200190518101906105bf91906109e0565b90506001825110156105e5578060405162461bcd60e51b81526004016103929190610a4e565b81816040516020016105f8929190610a61565b60408051601f198184030181529082905262461bcd60e51b825261039291600401610a4e565b505050565b6000808280602001905181019061063a9190610917565b9392505050565b606086868686868660405160240161065e96959493929190610a9e565b60408051601f198184030181529190526020810180516001600160e01b031663160d7c7360e01b17905290509695505050505050565b6060848484846040516024016106ad9493929190610ae6565b60408051601f198184030181529190526020810180516001600160e01b0316639b45009d60e01b1790529050949350505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610720576107206106e1565b604052919050565b600067ffffffffffffffff821115610742576107426106e1565b50601f01601f191660200190565b600082601f83011261076157600080fd5b813561077461076f82610728565b6106f7565b81815284602083860101111561078957600080fd5b816020850160208301376000918101602001919091529392505050565b600080604083850312156107b957600080fd5b823567ffffffffffffffff8111156107d057600080fd5b6107dc85828601610750565b95602094909401359450505050565b80356001600160a01b038116811461080257600080fd5b919050565b60008060008060008060c0878903121561082057600080fd5b610829876107eb565b9550602087013567ffffffffffffffff8082111561084657600080fd5b6108528a838b01610750565b965060408901359550606089013594506080890135935060a089013591508082111561087d57600080fd5b5061088a89828a01610750565b9150509295509295509295565b600080600080608085870312156108ad57600080fd5b843567ffffffffffffffff8111156108c457600080fd5b6108d087828801610750565b945050602085013592506108e6604086016107eb565b9396929550929360600135925050565b808201808211156100a257634e487b7160e01b600052601160045260246000fd5b60006020828403121561092957600080fd5b8151801515811461063a57600080fd5b60005b8381101561095457818101518382015260200161093c565b50506000910152565b6000825161096f818460208701610939565b9190910192915050565b60006020828403121561098b57600080fd5b5051919050565b600081518084526109aa816020860160208601610939565b601f01601f19169290920160200192915050565b6040815260006109d16040830185610992565b90508260208301529392505050565b6000602082840312156109f257600080fd5b815167ffffffffffffffff811115610a0957600080fd5b8201601f81018413610a1a57600080fd5b8051610a2861076f82610728565b818152856020838501011115610a3d57600080fd5b610201826020830160208601610939565b60208152600061009f6020830184610992565b60008351610a73818460208801610939565b6101d160f51b9083019081528351610a92816002840160208801610939565b01600201949350505050565b6001600160a01b038716815260c060208201819052600090610ac290830188610992565b86604084015285606084015284608084015282810360a08401526104a08185610992565b608081526000610af96080830187610992565b6020830195909552506001600160a01b0392909216604083015260609091015291905056fea26469706673582212204257c07f0f1337ceb6a97d057ae898670d1b31a12b1765735992285c74b61ef864736f6c63430008130033"

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

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txid, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) IncreaseBridgeFee(opts *bind.TransactOpts, _chain string, _txid *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "increaseBridgeFee", _chain, _txid, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txid, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestSession) IncreaseBridgeFee(_chain string, _txid *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.IncreaseBridgeFee(&_CrosschainTest.TransactOpts, _chain, _txid, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txid, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) IncreaseBridgeFee(_chain string, _txid *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.IncreaseBridgeFee(&_CrosschainTest.TransactOpts, _chain, _txid, _token, _fee)
}
