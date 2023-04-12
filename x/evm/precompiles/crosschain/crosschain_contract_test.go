// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package crosschain_test

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

// CrosschainTestMetaData contains all meta data concerning the CrosschainTest contract.
var CrosschainTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txID\",\"type\":\"uint256\"}],\"name\":\"CancelSendToExternal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"denom\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receipt\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"CrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"IncreaseBridgeFee\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"}],\"name\":\"cancelSendToExternal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"fip20CrossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_txID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"name\":\"increaseBridgeFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506113b9806100206000396000f3fe60806040526004361061003f5760003560e01c80630b56c19014610044578063160d7c73146100815780633c3e7d77146100b1578063c79a6b7b146100ee575b600080fd5b34801561005057600080fd5b5061006b60048036038101906100669190610b2f565b61011e565b6040516100789190610ba6565b60405180910390f35b61009b60048036038101906100969190610c55565b610132565b6040516100a89190610ba6565b60405180910390f35b3480156100bd57600080fd5b506100d860048036038101906100d39190610c55565b61029a565b6040516100e59190610ba6565b60405180910390f35b61010860048036038101906101039190610d1a565b6102aa565b6040516101159190610ba6565b60405180910390f35b600061012a83836102c2565b905092915050565b60008073ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff1614610280578673ffffffffffffffffffffffffffffffffffffffff166323b872dd333087896101919190610dcc565b6040518463ffffffff1660e01b81526004016101af93929190610e1e565b6020604051808303816000875af11580156101ce573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101f29190610e81565b508673ffffffffffffffffffffffffffffffffffffffff1663095ea7b3611004868861021e9190610dcc565b6040518363ffffffff1660e01b815260040161023b929190610eae565b6020604051808303816000875af115801561025a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027e9190610e81565b505b61028e878787878787610391565b90509695505050505050565b6000600190509695505050505050565b60006102b8858585856105bf565b9050949350505050565b600080600061100473ffffffffffffffffffffffffffffffffffffffff166102ea8686610692565b6040516102f79190610f48565b6000604051808303816000865af19150503d8060008114610334576040519150601f19603f3d011682016040523d82523d6000602084013e610339565b606091505b509150915061037e82826040518060400160405280601e81526020017f63616e63656c2073656e6420746f2065787465726e616c206661696c6564000081525061072c565b610387816107f3565b9250505092915050565b60008073ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff161461049b5760008773ffffffffffffffffffffffffffffffffffffffff1663dd62ed3e306110046040518363ffffffff1660e01b8152600401610405929190610f5f565b602060405180830381865afa158015610422573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104469190610f9d565b905084866104549190610dcc565b8114610495576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161048c90611027565b60405180910390fd5b506104e9565b83856104a79190610dcc565b34146104e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104df90611093565b60405180910390fd5b5b60008061100473ffffffffffffffffffffffffffffffffffffffff16346105148b8b8b8b8b8b610815565b6040516105219190610f48565b60006040518083038185875af1925050503d806000811461055e576040519150601f19603f3d011682016040523d82523d6000602084013e610563565b606091505b50915091506105a882826040518060400160405280601281526020017f63726f73732d636861696e206661696c6564000000000000000000000000000081525061072c565b6105b1816108bb565b925050509695505050505050565b600080600061100473ffffffffffffffffffffffffffffffffffffffff166105e9888888886108dd565b6040516105f69190610f48565b6000604051808303816000865af19150503d8060008114610633576040519150601f19603f3d011682016040523d82523d6000602084013e610638565b606091505b509150915061067d82826040518060400160405280601a81526020017f696e6372656173652062726964676520666565206661696c656400000000000081525061072c565b6106868161097d565b92505050949350505050565b606082826040516024016106a79291906110f7565b6040516020818303038152906040527feeb3593d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b826107ee576000828060200190518101906107479190611197565b905060018251101561079057806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078791906111e0565b60405180910390fd5b81816040516020016107a392919061128a565b6040516020818303038152906040526040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107e591906111e0565b60405180910390fd5b505050565b6000808280602001905181019061080a9190610e81565b905080915050919050565b6060868686868686604051602401610832969594939291906112c8565b6040516020818303038152906040527f160d7c73000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509695505050505050565b600080828060200190518101906108d29190610e81565b905080915050919050565b6060848484846040516024016108f69493929190611337565b6040516020818303038152906040527f9b45009d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050949350505050565b600080828060200190518101906109949190610e81565b905080915050919050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610a06826109bd565b810181811067ffffffffffffffff82111715610a2557610a246109ce565b5b80604052505050565b6000610a3861099f565b9050610a4482826109fd565b919050565b600067ffffffffffffffff821115610a6457610a636109ce565b5b610a6d826109bd565b9050602081019050919050565b82818337600083830152505050565b6000610a9c610a9784610a49565b610a2e565b905082815260208101848484011115610ab857610ab76109b8565b5b610ac3848285610a7a565b509392505050565b600082601f830112610ae057610adf6109b3565b5b8135610af0848260208601610a89565b91505092915050565b6000819050919050565b610b0c81610af9565b8114610b1757600080fd5b50565b600081359050610b2981610b03565b92915050565b60008060408385031215610b4657610b456109a9565b5b600083013567ffffffffffffffff811115610b6457610b636109ae565b5b610b7085828601610acb565b9250506020610b8185828601610b1a565b9150509250929050565b60008115159050919050565b610ba081610b8b565b82525050565b6000602082019050610bbb6000830184610b97565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610bec82610bc1565b9050919050565b610bfc81610be1565b8114610c0757600080fd5b50565b600081359050610c1981610bf3565b92915050565b6000819050919050565b610c3281610c1f565b8114610c3d57600080fd5b50565b600081359050610c4f81610c29565b92915050565b60008060008060008060c08789031215610c7257610c716109a9565b5b6000610c8089828a01610c0a565b965050602087013567ffffffffffffffff811115610ca157610ca06109ae565b5b610cad89828a01610acb565b9550506040610cbe89828a01610b1a565b9450506060610ccf89828a01610b1a565b9350506080610ce089828a01610c40565b92505060a087013567ffffffffffffffff811115610d0157610d006109ae565b5b610d0d89828a01610acb565b9150509295509295509295565b60008060008060808587031215610d3457610d336109a9565b5b600085013567ffffffffffffffff811115610d5257610d516109ae565b5b610d5e87828801610acb565b9450506020610d6f87828801610b1a565b9350506040610d8087828801610c0a565b9250506060610d9187828801610b1a565b91505092959194509250565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610dd782610af9565b9150610de283610af9565b9250828201905080821115610dfa57610df9610d9d565b5b92915050565b610e0981610be1565b82525050565b610e1881610af9565b82525050565b6000606082019050610e336000830186610e00565b610e406020830185610e00565b610e4d6040830184610e0f565b949350505050565b610e5e81610b8b565b8114610e6957600080fd5b50565b600081519050610e7b81610e55565b92915050565b600060208284031215610e9757610e966109a9565b5b6000610ea584828501610e6c565b91505092915050565b6000604082019050610ec36000830185610e00565b610ed06020830184610e0f565b9392505050565b600081519050919050565b600081905092915050565b60005b83811015610f0b578082015181840152602081019050610ef0565b60008484015250505050565b6000610f2282610ed7565b610f2c8185610ee2565b9350610f3c818560208601610eed565b80840191505092915050565b6000610f548284610f17565b915081905092915050565b6000604082019050610f746000830185610e00565b610f816020830184610e00565b9392505050565b600081519050610f9781610b03565b92915050565b600060208284031215610fb357610fb26109a9565b5b6000610fc184828501610f88565b91505092915050565b600082825260208201905092915050565b7f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b20666565600082015250565b6000611011602083610fca565b915061101c82610fdb565b602082019050919050565b6000602082019050818103600083015261104081611004565b9050919050565b7f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b20666565600082015250565b600061107d602083610fca565b915061108882611047565b602082019050919050565b600060208201905081810360008301526110ac81611070565b9050919050565b600081519050919050565b60006110c9826110b3565b6110d38185610fca565b93506110e3818560208601610eed565b6110ec816109bd565b840191505092915050565b6000604082019050818103600083015261111181856110be565b90506111206020830184610e0f565b9392505050565b600061113a61113584610a49565b610a2e565b905082815260208101848484011115611156576111556109b8565b5b611161848285610eed565b509392505050565b600082601f83011261117e5761117d6109b3565b5b815161118e848260208601611127565b91505092915050565b6000602082840312156111ad576111ac6109a9565b5b600082015167ffffffffffffffff8111156111cb576111ca6109ae565b5b6111d784828501611169565b91505092915050565b600060208201905081810360008301526111fa81846110be565b905092915050565b600081905092915050565b6000611218826110b3565b6112228185611202565b9350611232818560208601610eed565b80840191505092915050565b7f3a20000000000000000000000000000000000000000000000000000000000000600082015250565b6000611274600283611202565b915061127f8261123e565b600282019050919050565b6000611296828561120d565b91506112a182611267565b91506112ad828461120d565b91508190509392505050565b6112c281610c1f565b82525050565b600060c0820190506112dd6000830189610e00565b81810360208301526112ef81886110be565b90506112fe6040830187610e0f565b61130b6060830186610e0f565b61131860808301856112b9565b81810360a083015261132a81846110be565b9050979650505050505050565b6000608082019050818103600083015261135181876110be565b90506113606020830186610e0f565b61136d6040830185610e00565b61137a6060830184610e0f565b9594505050505056fea26469706673582212207cd298166473fc4ef10450aadfdcfec8f03ebbcd4a3c7dc2a5458b6179152df164736f6c63430008130033",
}

// CrosschainTestABI is the input ABI used to generate the binding from.
// Deprecated: Use CrosschainTestMetaData.ABI instead.
var CrosschainTestABI = CrosschainTestMetaData.ABI

// CrosschainTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CrosschainTestMetaData.Bin instead.
var CrosschainTestBin = CrosschainTestMetaData.Bin

// DeployCrosschainTest deploys a new Ethereum contract, binding an instance of CrosschainTest to it.
func DeployCrosschainTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CrosschainTest, error) {
	parsed, err := CrosschainTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CrosschainTestBin), backend)
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
	parsed, err := CrosschainTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) CancelSendToExternal(opts *bind.TransactOpts, _chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "cancelSendToExternal", _chain, _txID)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrosschainTest *CrosschainTestSession) CancelSendToExternal(_chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CancelSendToExternal(&_CrosschainTest.TransactOpts, _chain, _txID)
}

// CancelSendToExternal is a paid mutator transaction binding the contract method 0x0b56c190.
//
// Solidity: function cancelSendToExternal(string _chain, uint256 _txID) returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) CancelSendToExternal(_chain string, _txID *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.CancelSendToExternal(&_CrosschainTest.TransactOpts, _chain, _txID)
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

// Fip20CrossChain is a paid mutator transaction binding the contract method 0x3c3e7d77.
//
// Solidity: function fip20CrossChain(address _sender, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) returns(bool _result)
func (_CrosschainTest *CrosschainTestTransactor) Fip20CrossChain(opts *bind.TransactOpts, _sender common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "fip20CrossChain", _sender, _receipt, _amount, _fee, _target, _memo)
}

// Fip20CrossChain is a paid mutator transaction binding the contract method 0x3c3e7d77.
//
// Solidity: function fip20CrossChain(address _sender, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) returns(bool _result)
func (_CrosschainTest *CrosschainTestSession) Fip20CrossChain(_sender common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.Contract.Fip20CrossChain(&_CrosschainTest.TransactOpts, _sender, _receipt, _amount, _fee, _target, _memo)
}

// Fip20CrossChain is a paid mutator transaction binding the contract method 0x3c3e7d77.
//
// Solidity: function fip20CrossChain(address _sender, string _receipt, uint256 _amount, uint256 _fee, bytes32 _target, string _memo) returns(bool _result)
func (_CrosschainTest *CrosschainTestTransactorSession) Fip20CrossChain(_sender common.Address, _receipt string, _amount *big.Int, _fee *big.Int, _target [32]byte, _memo string) (*types.Transaction, error) {
	return _CrosschainTest.Contract.Fip20CrossChain(&_CrosschainTest.TransactOpts, _sender, _receipt, _amount, _fee, _target, _memo)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactor) IncreaseBridgeFee(opts *bind.TransactOpts, _chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "increaseBridgeFee", _chain, _txID, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestSession) IncreaseBridgeFee(_chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.IncreaseBridgeFee(&_CrosschainTest.TransactOpts, _chain, _txID, _token, _fee)
}

// IncreaseBridgeFee is a paid mutator transaction binding the contract method 0xc79a6b7b.
//
// Solidity: function increaseBridgeFee(string _chain, uint256 _txID, address _token, uint256 _fee) payable returns(bool)
func (_CrosschainTest *CrosschainTestTransactorSession) IncreaseBridgeFee(_chain string, _txID *big.Int, _token common.Address, _fee *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.IncreaseBridgeFee(&_CrosschainTest.TransactOpts, _chain, _txID, _token, _fee)
}

// CrosschainTestCancelSendToExternalIterator is returned from FilterCancelSendToExternal and is used to iterate over the raw logs and unpacked data for CancelSendToExternal events raised by the CrosschainTest contract.
type CrosschainTestCancelSendToExternalIterator struct {
	Event *CrosschainTestCancelSendToExternal // Event containing the contract specifics and raw log

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
func (it *CrosschainTestCancelSendToExternalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainTestCancelSendToExternal)
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
		it.Event = new(CrosschainTestCancelSendToExternal)
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
func (it *CrosschainTestCancelSendToExternalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainTestCancelSendToExternalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainTestCancelSendToExternal represents a CancelSendToExternal event raised by the CrosschainTest contract.
type CrosschainTestCancelSendToExternal struct {
	Sender common.Address
	Chain  string
	TxID   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCancelSendToExternal is a free log retrieval operation binding the contract event 0xe2ae965fb5b8e4c7da962424292951c18e0e9c1905b87c78cf0186fa70382535.
//
// Solidity: event CancelSendToExternal(address indexed sender, string chain, uint256 txID)
func (_CrosschainTest *CrosschainTestFilterer) FilterCancelSendToExternal(opts *bind.FilterOpts, sender []common.Address) (*CrosschainTestCancelSendToExternalIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CrosschainTest.contract.FilterLogs(opts, "CancelSendToExternal", senderRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestCancelSendToExternalIterator{contract: _CrosschainTest.contract, event: "CancelSendToExternal", logs: logs, sub: sub}, nil
}

// WatchCancelSendToExternal is a free log subscription operation binding the contract event 0xe2ae965fb5b8e4c7da962424292951c18e0e9c1905b87c78cf0186fa70382535.
//
// Solidity: event CancelSendToExternal(address indexed sender, string chain, uint256 txID)
func (_CrosschainTest *CrosschainTestFilterer) WatchCancelSendToExternal(opts *bind.WatchOpts, sink chan<- *CrosschainTestCancelSendToExternal, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _CrosschainTest.contract.WatchLogs(opts, "CancelSendToExternal", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainTestCancelSendToExternal)
				if err := _CrosschainTest.contract.UnpackLog(event, "CancelSendToExternal", log); err != nil {
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

// ParseCancelSendToExternal is a log parse operation binding the contract event 0xe2ae965fb5b8e4c7da962424292951c18e0e9c1905b87c78cf0186fa70382535.
//
// Solidity: event CancelSendToExternal(address indexed sender, string chain, uint256 txID)
func (_CrosschainTest *CrosschainTestFilterer) ParseCancelSendToExternal(log types.Log) (*CrosschainTestCancelSendToExternal, error) {
	event := new(CrosschainTestCancelSendToExternal)
	if err := _CrosschainTest.contract.UnpackLog(event, "CancelSendToExternal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrosschainTestCrossChainIterator is returned from FilterCrossChain and is used to iterate over the raw logs and unpacked data for CrossChain events raised by the CrosschainTest contract.
type CrosschainTestCrossChainIterator struct {
	Event *CrosschainTestCrossChain // Event containing the contract specifics and raw log

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
func (it *CrosschainTestCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainTestCrossChain)
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
		it.Event = new(CrosschainTestCrossChain)
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
func (it *CrosschainTestCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainTestCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainTestCrossChain represents a CrossChain event raised by the CrosschainTest contract.
type CrosschainTestCrossChain struct {
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
func (_CrosschainTest *CrosschainTestFilterer) FilterCrossChain(opts *bind.FilterOpts, sender []common.Address, token []common.Address) (*CrosschainTestCrossChainIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.FilterLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestCrossChainIterator{contract: _CrosschainTest.contract, event: "CrossChain", logs: logs, sub: sub}, nil
}

// WatchCrossChain is a free log subscription operation binding the contract event 0xb783df819ac99ca709650d67d9237a00b553c6ef941dceabeed6f4bc990d31ba.
//
// Solidity: event CrossChain(address indexed sender, address indexed token, string denom, string receipt, uint256 amount, uint256 fee, bytes32 target, string memo)
func (_CrosschainTest *CrosschainTestFilterer) WatchCrossChain(opts *bind.WatchOpts, sink chan<- *CrosschainTestCrossChain, sender []common.Address, token []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.WatchLogs(opts, "CrossChain", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainTestCrossChain)
				if err := _CrosschainTest.contract.UnpackLog(event, "CrossChain", log); err != nil {
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
func (_CrosschainTest *CrosschainTestFilterer) ParseCrossChain(log types.Log) (*CrosschainTestCrossChain, error) {
	event := new(CrosschainTestCrossChain)
	if err := _CrosschainTest.contract.UnpackLog(event, "CrossChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrosschainTestIncreaseBridgeFeeIterator is returned from FilterIncreaseBridgeFee and is used to iterate over the raw logs and unpacked data for IncreaseBridgeFee events raised by the CrosschainTest contract.
type CrosschainTestIncreaseBridgeFeeIterator struct {
	Event *CrosschainTestIncreaseBridgeFee // Event containing the contract specifics and raw log

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
func (it *CrosschainTestIncreaseBridgeFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrosschainTestIncreaseBridgeFee)
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
		it.Event = new(CrosschainTestIncreaseBridgeFee)
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
func (it *CrosschainTestIncreaseBridgeFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrosschainTestIncreaseBridgeFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrosschainTestIncreaseBridgeFee represents a IncreaseBridgeFee event raised by the CrosschainTest contract.
type CrosschainTestIncreaseBridgeFee struct {
	Sender common.Address
	Token  common.Address
	Chain  string
	TxID   *big.Int
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIncreaseBridgeFee is a free log retrieval operation binding the contract event 0x4b4d0e64eb77c0f61892107908295f09b3e381c50c655f4a73a4ad61c07350a0.
//
// Solidity: event IncreaseBridgeFee(address indexed sender, address indexed token, string chain, uint256 txID, uint256 fee)
func (_CrosschainTest *CrosschainTestFilterer) FilterIncreaseBridgeFee(opts *bind.FilterOpts, sender []common.Address, token []common.Address) (*CrosschainTestIncreaseBridgeFeeIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.FilterLogs(opts, "IncreaseBridgeFee", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &CrosschainTestIncreaseBridgeFeeIterator{contract: _CrosschainTest.contract, event: "IncreaseBridgeFee", logs: logs, sub: sub}, nil
}

// WatchIncreaseBridgeFee is a free log subscription operation binding the contract event 0x4b4d0e64eb77c0f61892107908295f09b3e381c50c655f4a73a4ad61c07350a0.
//
// Solidity: event IncreaseBridgeFee(address indexed sender, address indexed token, string chain, uint256 txID, uint256 fee)
func (_CrosschainTest *CrosschainTestFilterer) WatchIncreaseBridgeFee(opts *bind.WatchOpts, sink chan<- *CrosschainTestIncreaseBridgeFee, sender []common.Address, token []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _CrosschainTest.contract.WatchLogs(opts, "IncreaseBridgeFee", senderRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrosschainTestIncreaseBridgeFee)
				if err := _CrosschainTest.contract.UnpackLog(event, "IncreaseBridgeFee", log); err != nil {
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

// ParseIncreaseBridgeFee is a log parse operation binding the contract event 0x4b4d0e64eb77c0f61892107908295f09b3e381c50c655f4a73a4ad61c07350a0.
//
// Solidity: event IncreaseBridgeFee(address indexed sender, address indexed token, string chain, uint256 txID, uint256 fee)
func (_CrosschainTest *CrosschainTestFilterer) ParseIncreaseBridgeFee(log types.Log) (*CrosschainTestIncreaseBridgeFee, error) {
	event := new(CrosschainTestIncreaseBridgeFee)
	if err := _CrosschainTest.contract.UnpackLog(event, "IncreaseBridgeFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
