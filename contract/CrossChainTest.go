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

// CrossChainTestMetaData contains all meta data concerning the CrossChainTest contract.
var CrossChainTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"increaseBridgeFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610b77806100206000396000f3fe6080604052600436106100345760003560e01c80630b56c19014610039578063160d7c731461006d578063c79a6b7b14610080575b600080fd5b34801561004557600080fd5b5061005961005436600461088c565b610093565b604051901515815260200160405180910390f35b61005961007b366004610775565b6100a6565b61005961008e3660046108cf565b61035d565b600061009f8383610374565b9392505050565b60006001600160a01b038716156101f6576001600160a01b0387166323b872dd33306100d2888a610ad7565b6040516001600160e01b031960e086901b1681526001600160a01b0393841660048201529290911660248301526044820152606401602060405180830381600087803b15801561012157600080fd5b505af1158015610135573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101599190610802565b506001600160a01b03871663095ea7b36110046101768789610ad7565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401602060405180830381600087803b1580156101bc57600080fd5b505af11580156101d0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101f49190610802565b505b6001600160a01b038716156102ec57604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e9060440160206040518083038186803b15801561024f57600080fd5b505afa158015610263573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610287919061092c565b90506102938587610ad7565b81146102e65760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b50610344565b6102f68486610ad7565b34146103445760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b2066656560448201526064016102dd565b610352878787878787610422565b979650505050505050565b600061036b858585856104d7565b95945050505050565b600080806110046103858686610587565b6040516103929190610970565b6000604051808303816000865af19150503d80600081146103cf576040519150601f19603f3d011682016040523d82523d6000602084013e6103d4565b606091505b509150915061041982826040518060400160405280601e81526020017f63616e63656c2073656e6420746f2065787465726e616c206661696c656400008152506105ce565b61036b8161064d565b60008080611004346104388b8b8b8b8b8b61066d565b6040516104459190610970565b60006040518083038185875af1925050503d8060008114610482576040519150601f19603f3d011682016040523d82523d6000602084013e610487565b606091505b50915091506104c182826040518060400160405280601281526020017118dc9bdcdccb58da185a5b8819985a5b195960721b8152506105ce565b6104ca8161064d565b9998505050505050505050565b600080806110046104ea888888886106c0565b6040516104f79190610970565b6000604051808303816000865af19150503d8060008114610534576040519150601f19603f3d011682016040523d82523d6000602084013e610539565b606091505b509150915061057e82826040518060400160405280601a81526020017f696e6372656173652062726964676520666565206661696c65640000000000008152506105ce565b6103528161064d565b6060828260405160240161059c929190610a24565b60408051601f198184030181529190526020810180516001600160e01b031663eeb3593d60e01b179052905092915050565b82610648576000828060200190518101906105e99190610822565b905060018251101561060f578060405162461bcd60e51b81526004016102dd9190610a11565b818160405160200161062292919061098c565b60408051601f198184030181529082905262461bcd60e51b82526102dd91600401610a11565b505050565b600080828060200190518101906106649190610802565b9150505b919050565b606086868686868660405160240161068a969594939291906109c9565b60408051601f198184030181529190526020810180516001600160e01b031663160d7c7360e01b17905290509695505050505050565b6060848484846040516024016106d99493929190610a46565b60408051601f198184030181529190526020810180516001600160e01b0316639b45009d60e01b1790529050949350505050565b80356001600160a01b038116811461066857600080fd5b600082601f830112610734578081fd5b813561074761074282610aaf565b610a7e565b81815284602083860101111561075b578283fd5b816020850160208301379081016020019190915292915050565b60008060008060008060c0878903121561078d578182fd5b6107968761070d565b9550602087013567ffffffffffffffff808211156107b2578384fd5b6107be8a838b01610724565b965060408901359550606089013594506080890135935060a08901359150808211156107e8578283fd5b506107f589828a01610724565b9150509295509295509295565b600060208284031215610813578081fd5b8151801515811461009f578182fd5b600060208284031215610833578081fd5b815167ffffffffffffffff811115610849578182fd5b8201601f81018413610859578182fd5b805161086761074282610aaf565b81815285602083850101111561087b578384fd5b61036b826020830160208601610afb565b6000806040838503121561089e578182fd5b823567ffffffffffffffff8111156108b4578283fd5b6108c085828601610724565b95602094909401359450505050565b600080600080608085870312156108e4578384fd5b843567ffffffffffffffff8111156108fa578485fd5b61090687828801610724565b9450506020850135925061091c6040860161070d565b9396929550929360600135925050565b60006020828403121561093d578081fd5b5051919050565b6000815180845261095c816020860160208601610afb565b601f01601f19169290920160200192915050565b60008251610982818460208701610afb565b9190910192915050565b6000835161099e818460208801610afb565b6101d160f51b90830190815283516109bd816002840160208801610afb565b01600201949350505050565b6001600160a01b038716815260c0602082018190526000906109ed90830188610944565b86604084015285606084015284608084015282810360a08401526104ca8185610944565b60006020825261009f6020830184610944565b600060408252610a376040830185610944565b90508260208301529392505050565b600060808252610a596080830187610944565b6020830195909552506001600160a01b03929092166040830152606090910152919050565b604051601f8201601f1916810167ffffffffffffffff81118282101715610aa757610aa7610b2b565b604052919050565b600067ffffffffffffffff821115610ac957610ac9610b2b565b50601f01601f191660200190565b60008219821115610af657634e487b7160e01b81526011600452602481fd5b500190565b60005b83811015610b16578181015183820152602001610afe565b83811115610b25576000848401525b50505050565b634e487b7160e01b600052604160045260246000fdfea2646970667358221220b4ceaa85e103e483802e547fcd87df04b31230c37751b81ad785688d218737e164736f6c63430008020033",
}

// CrossChainTestABI is the input ABI used to generate the binding from.
// Deprecated: Use CrossChainTestMetaData.ABI instead.
var CrossChainTestABI = CrossChainTestMetaData.ABI

// CrossChainTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrossChainTestMetaData.Bin instead.
var CrossChainTestBin = CrossChainTestMetaData.Bin

// DeployCrossChainTest deploys a new Ethereum contract, binding an instance of CrossChainTest to it.
func DeployCrossChainTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrossChainTest, error) {
	parsed, err := CrossChainTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrossChainTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CrossChainTest{CrossChainTestCaller: CrossChainTestCaller{contract: contract}, CrossChainTestTransactor: CrossChainTestTransactor{contract: contract}, CrossChainTestFilterer: CrossChainTestFilterer{contract: contract}}, nil
}

// CrossChainTest is an auto generated Go binding around an Ethereum contract.
type CrossChainTest struct {
	CrossChainTestCaller     // Read-only binding to the contract
	CrossChainTestTransactor // Write-only binding to the contract
	CrossChainTestFilterer   // Log filterer for contract events
}

// CrossChainTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrossChainTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrossChainTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrossChainTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrossChainTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrossChainTestSession struct {
	Contract     *CrossChainTest   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrossChainTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrossChainTestCallerSession struct {
	Contract *CrossChainTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// CrossChainTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrossChainTestTransactorSession struct {
	Contract     *CrossChainTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CrossChainTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrossChainTestRaw struct {
	Contract *CrossChainTest // Generic contract binding to access the raw methods on
}

// CrossChainTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrossChainTestCallerRaw struct {
	Contract *CrossChainTestCaller // Generic read-only contract binding to access the raw methods on
}

// CrossChainTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrossChainTestTransactorRaw struct {
	Contract *CrossChainTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrossChainTest creates a new instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTest(address common.Address, backend bind.ContractBackend) (*CrossChainTest, error) {
	contract, err := bindCrossChainTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrossChainTest{CrossChainTestCaller: CrossChainTestCaller{contract: contract}, CrossChainTestTransactor: CrossChainTestTransactor{contract: contract}, CrossChainTestFilterer: CrossChainTestFilterer{contract: contract}}, nil
}

// NewCrossChainTestCaller creates a new read-only instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTestCaller(address common.Address, caller bind.ContractCaller) (*CrossChainTestCaller, error) {
	contract, err := bindCrossChainTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainTestCaller{contract: contract}, nil
}

// NewCrossChainTestTransactor creates a new write-only instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTestTransactor(address common.Address, transactor bind.ContractTransactor) (*CrossChainTestTransactor, error) {
	contract, err := bindCrossChainTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrossChainTestTransactor{contract: contract}, nil
}

// NewCrossChainTestFilterer creates a new log filterer instance of CrossChainTest, bound to a specific deployed contract.
func NewCrossChainTestFilterer(address common.Address, filterer bind.ContractFilterer) (*CrossChainTestFilterer, error) {
	contract, err := bindCrossChainTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrossChainTestFilterer{contract: contract}, nil
}

// bindCrossChainTest binds a generic wrapper to an already deployed contract.
func bindCrossChainTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrossChainTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainTest *CrossChainTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainTest.Contract.CrossChainTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainTest *CrossChainTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChainTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainTest *CrossChainTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChainTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrossChainTest *CrossChainTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrossChainTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrossChainTest *CrossChainTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrossChainTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrossChainTest *CrossChainTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrossChainTest.Contract.contract.Transact(opts, method, params...)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrossChainTest *CrossChainTestTransactor) CancelSendToExternal(opts *bind.TransactOpts, _chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.contract.Transact(opts, "cancelSendToExternal", _chain, _txID)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrossChainTest *CrossChainTestSession) CancelSendToExternal(_chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CancelSendToExternal(&_CrossChainTest.TransactOpts, _chain, _txID)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrossChainTest *CrossChainTestTransactorSession) CancelSendToExternal(_chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CancelSendToExternal(&_CrossChainTest.TransactOpts, _chain, _txID)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactor) CrossChain(opts *bind.TransactOpts, _token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrossChainTest.contract.Transact(opts, "crossChain", _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrossChainTest *CrossChainTestSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChain(&_CrossChainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// CrossChain is a paid mutator transaction binding the contract method 0x160d7c73.
//
// Solidity: function crossChain(address _token, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactorSession) CrossChain(_token common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrossChainTest.Contract.CrossChain(&_CrossChainTest.TransactOpts, _token, _receipt, _amount, _fee, _target, _memo)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactor) IncreaseBridgeFee(opts *bind.TransactOpts, _chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.contract.Transact(opts, "increaseBridgeFee", _chain, _txID, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrossChainTest *CrossChainTestSession) IncreaseBridgeFee(_chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.IncreaseBridgeFee(&_CrossChainTest.TransactOpts, _chain, _txID, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrossChainTest *CrossChainTestTransactorSession) IncreaseBridgeFee(_chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrossChainTest.Contract.IncreaseBridgeFee(&_CrossChainTest.TransactOpts, _chain, _txID, _token, _fee)
}
