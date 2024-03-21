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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"increaseBridgeFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610cea806100206000396000f3fe60806040526004361061003f5760003560e01c80630b56c19014610044578063160d7c73146100795780638fefb7651461008c578063c79a6b7b146100ba575b600080fd5b34801561005057600080fd5b5061006461005f3660046109ff565b6100cd565b60405190151581526020015b60405180910390f35b6100646100873660046108e8565b6100e0565b34801561009857600080fd5b506100ac6100a73660046108bf565b610397565b604051908152602001610070565b6100646100c8366004610a42565b6103a3565b60006100d983836103ba565b9392505050565b60006001600160a01b03871615610230576001600160a01b0387166323b872dd333061010c888a610c4a565b6040516001600160e01b031960e086901b1681526001600160a01b0393841660048201529290911660248301526044820152606401602060405180830381600087803b15801561015b57600080fd5b505af115801561016f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101939190610975565b506001600160a01b03871663095ea7b36110046101b08789610c4a565b6040516001600160e01b031960e085901b1681526001600160a01b0390921660048301526024820152604401602060405180830381600087803b1580156101f657600080fd5b505af115801561020a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061022e9190610975565b505b6001600160a01b0387161561032657604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e9060440160206040518083038186803b15801561028957600080fd5b505afa15801561029d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c19190610a9f565b90506102cd8587610c4a565b81146103205760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b5061037e565b6103308486610c4a565b341461037e5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b206665656044820152606401610317565b61038c878787878787610468565b979650505050505050565b60006100d9838361051d565b60006103b18585858561060a565b95945050505050565b600080806110046103cb86866106ba565b6040516103d89190610ae3565b6000604051808303816000865af19150503d8060008114610415576040519150601f19603f3d011682016040523d82523d6000602084013e61041a565b606091505b509150915061045f82826040518060400160405280601e81526020017f63616e63656c2073656e6420746f2065787465726e616c206661696c65640000815250610701565b6103b181610780565b600080806110043461047e8b8b8b8b8b8b6107a0565b60405161048b9190610ae3565b60006040518083038185875af1925050503d80600081146104c8576040519150601f19603f3d011682016040523d82523d6000602084013e6104cd565b606091505b509150915061050782826040518060400160405280601281526020017118dc9bdcdccb58da185a5b8819985a5b195960721b815250610701565b61051081610780565b9998505050505050505050565b6000808061100461057a86866040516001600160a01b03831660248201526044810182905260609060640160408051601f198184030181529190526020810180516001600160e01b0316638fefb76560e01b179052905092915050565b6040516105879190610ae3565b600060405180830381855afa9150503d80600081146105c2576040519150601f19603f3d011682016040523d82523d6000602084013e6105c7565b606091505b5091509150610601828260405180604001604052806012815260200171189c9a5919d94818dbda5b8819985a5b195960721b815250610701565b6103b1816107f3565b6000808061100461061d8888888861080a565b60405161062a9190610ae3565b6000604051808303816000865af19150503d8060008114610667576040519150601f19603f3d011682016040523d82523d6000602084013e61066c565b606091505b50915091506106b182826040518060400160405280601a81526020017f696e6372656173652062726964676520666565206661696c6564000000000000815250610701565b61038c81610780565b606082826040516024016106cf929190610b97565b60408051601f198184030181529190526020810180516001600160e01b031663eeb3593d60e01b179052905092915050565b8261077b5760008280602001905181019061071c9190610995565b9050600182511015610742578060405162461bcd60e51b81526004016103179190610b84565b8181604051602001610755929190610aff565b60408051601f198184030181529082905262461bcd60e51b825261031791600401610b84565b505050565b600080828060200190518101906107979190610975565b9150505b919050565b60608686868686866040516024016107bd96959493929190610b3c565b60408051601f198184030181529190526020810180516001600160e01b031663160d7c7360e01b17905290509695505050505050565b600080828060200190518101906107979190610a9f565b6060848484846040516024016108239493929190610bb9565b60408051601f198184030181529190526020810180516001600160e01b0316639b45009d60e01b1790529050949350505050565b80356001600160a01b038116811461079b57600080fd5b600082601f83011261087e578081fd5b813561089161088c82610c22565b610bf1565b8181528460208386010111156108a5578283fd5b816020850160208301379081016020019190915292915050565b600080604083850312156108d1578182fd5b6108da83610857565b946020939093013593505050565b60008060008060008060c08789031215610900578182fd5b61090987610857565b9550602087013567ffffffffffffffff80821115610925578384fd5b6109318a838b0161086e565b965060408901359550606089013594506080890135935060a089013591508082111561095b578283fd5b5061096889828a0161086e565b9150509295509295509295565b600060208284031215610986578081fd5b815180151581146100d9578182fd5b6000602082840312156109a6578081fd5b815167ffffffffffffffff8111156109bc578182fd5b8201601f810184136109cc578182fd5b80516109da61088c82610c22565b8181528560208385010111156109ee578384fd5b6103b1826020830160208601610c6e565b60008060408385031215610a11578182fd5b823567ffffffffffffffff811115610a27578283fd5b610a338582860161086e565b95602094909401359450505050565b60008060008060808587031215610a57578384fd5b843567ffffffffffffffff811115610a6d578485fd5b610a798782880161086e565b94505060208501359250610a8f60408601610857565b9396929550929360600135925050565b600060208284031215610ab0578081fd5b5051919050565b60008151808452610acf816020860160208601610c6e565b601f01601f19169290920160200192915050565b60008251610af5818460208701610c6e565b9190910192915050565b60008351610b11818460208801610c6e565b6101d160f51b9083019081528351610b30816002840160208801610c6e565b01600201949350505050565b6001600160a01b038716815260c060208201819052600090610b6090830188610ab7565b86604084015285606084015284608084015282810360a08401526105108185610ab7565b6000602082526100d96020830184610ab7565b600060408252610baa6040830185610ab7565b90508260208301529392505050565b600060808252610bcc6080830187610ab7565b6020830195909552506001600160a01b03929092166040830152606090910152919050565b604051601f8201601f1916810167ffffffffffffffff81118282101715610c1a57610c1a610c9e565b604052919050565b600067ffffffffffffffff821115610c3c57610c3c610c9e565b50601f01601f191660200190565b60008219821115610c6957634e487b7160e01b81526011600452602481fd5b500190565b60005b83811015610c89578181015183820152602001610c71565b83811115610c98576000848401525b50505050565b634e487b7160e01b600052604160045260246000fdfea2646970667358221220967af813a3fbc39dcf52c663908ee2066ea3d3dace71c7a0ca0858e247f4ce8c64736f6c63430008020033",
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

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrossChainTest *CrossChainTestCaller) BridgeCoinAmount(opts *bind.CallOpts, _token common.Address, _target [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _CrossChainTest.contract.Call(opts, &out, "bridgeCoinAmount", _token, _target)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrossChainTest *CrossChainTestSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _CrossChainTest.Contract.BridgeCoinAmount(&_CrossChainTest.CallOpts, _token, _target)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrossChainTest *CrossChainTestCallerSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _CrossChainTest.Contract.BridgeCoinAmount(&_CrossChainTest.CallOpts, _token, _target)
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
