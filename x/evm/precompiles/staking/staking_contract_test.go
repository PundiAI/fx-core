// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking_test

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
)

// StakingTestMetaData contains all meta data concerning the StakingTest contract.
var StakingTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111ee806100206000396000f3fe60806040526004361061004a5760003560e01c806331fb67c21461004f57806351af513a1461008c5780638dfc8897146100c95780639ddb511a14610108578063d5c498eb14610139575b600080fd5b34801561005b57600080fd5b5061007660048036038101906100719190610be6565b610177565b6040516100839190610c48565b60405180910390f35b34801561009857600080fd5b506100b360048036038101906100ae9190610cc1565b610189565b6040516100c09190610c48565b60405180910390f35b3480156100d557600080fd5b506100f060048036038101906100eb9190610d49565b61019d565b6040516100ff93929190610da5565b60405180910390f35b610122600480360381019061011d9190610be6565b6101b9565b604051610130929190610ddc565b60405180910390f35b34801561014557600080fd5b50610160600480360381019061015b9190610cc1565b6101ce565b60405161016e929190610ddc565b60405180910390f35b6000610182826101e6565b9050919050565b600061019583836102b3565b905092915050565b60008060006101ac8585610380565b9250925092509250925092565b6000806101c583610457565b91509150915091565b6000806101db8484610528565b915091509250929050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff1661020d856105f9565b60405161021a9190610e76565b6000604051808303816000865af19150503d8060008114610257576040519150601f19603f3d011682016040523d82523d6000602084013e61025c565b606091505b50915091506102a182826040518060400160405280600f81526020017f7769746864726177206661696c65640000000000000000000000000000000000815250610690565b6102aa81610757565b92505050919050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff166102db8686610779565b6040516102e89190610e76565b600060405180830381855afa9150503d8060008114610323576040519150601f19603f3d011682016040523d82523d6000602084013e610328565b606091505b509150915061036d82826040518060400160405280601881526020017f64656c65676174696f6e52657761726473206661696c65640000000000000000815250610690565b61037681610813565b9250505092915050565b600080600080600061100373ffffffffffffffffffffffffffffffffffffffff166103ab8888610835565b6040516103b89190610e76565b6000604051808303816000865af19150503d80600081146103f5576040519150601f19603f3d011682016040523d82523d6000602084013e6103fa565b606091505b509150915061043f82826040518060400160405280601181526020017f756e64656c6567617465206661696c6564000000000000000000000000000000815250610690565b610448816108cf565b94509450945050509250925092565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff163461048087610905565b60405161048d9190610e76565b60006040518083038185875af1925050503d80600081146104ca576040519150601f19603f3d011682016040523d82523d6000602084013e6104cf565b606091505b509150915061051482826040518060400160405280600f81526020017f64656c6567617465206661696c65640000000000000000000000000000000000815250610690565b61051d8161099c565b935093505050915091565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff1661055187876109c7565b60405161055e9190610e76565b600060405180830381855afa9150503d8060008114610599576040519150601f19603f3d011682016040523d82523d6000602084013e61059e565b606091505b50915091506105e382826040518060400160405280601181526020017f64656c65676174696f6e206661696c6564000000000000000000000000000000815250610690565b6105ec81610a61565b9350935050509250929050565b60608160405160240161060c9190610ee2565b6040516020818303038152906040527f31fb67c2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b82610752576000828060200190518101906106ab9190610f74565b90506001825110156106f457806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106eb9190610ee2565b60405180910390fd5b8181604051602001610707929190611045565b6040516020818303038152906040526040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107499190610ee2565b60405180910390fd5b505050565b6000808280602001905181019061076e9190611089565b905080915050919050565b6060828260405160240161078e9291906110c5565b6040516020818303038152906040527f51af513a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b6000808280602001905181019061082a9190611089565b905080915050919050565b6060828260405160240161084a9291906110f5565b6040516020818303038152906040527f8dfc8897000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b600080600080600080868060200190518101906108ec9190611125565b9250925092508282829550955095505050509193909250565b6060816040516024016109189190610ee2565b6040516020818303038152906040527f9ddb511a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b600080600080848060200190518101906109b69190611178565b915091508181935093505050915091565b606082826040516024016109dc9291906110c5565b6040516020818303038152906040527fd5c498eb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b60008060008084806020019051810190610a7b9190611178565b915091508181935093505050915091565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610af382610aaa565b810181811067ffffffffffffffff82111715610b1257610b11610abb565b5b80604052505050565b6000610b25610a8c565b9050610b318282610aea565b919050565b600067ffffffffffffffff821115610b5157610b50610abb565b5b610b5a82610aaa565b9050602081019050919050565b82818337600083830152505050565b6000610b89610b8484610b36565b610b1b565b905082815260208101848484011115610ba557610ba4610aa5565b5b610bb0848285610b67565b509392505050565b600082601f830112610bcd57610bcc610aa0565b5b8135610bdd848260208601610b76565b91505092915050565b600060208284031215610bfc57610bfb610a96565b5b600082013567ffffffffffffffff811115610c1a57610c19610a9b565b5b610c2684828501610bb8565b91505092915050565b6000819050919050565b610c4281610c2f565b82525050565b6000602082019050610c5d6000830184610c39565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610c8e82610c63565b9050919050565b610c9e81610c83565b8114610ca957600080fd5b50565b600081359050610cbb81610c95565b92915050565b60008060408385031215610cd857610cd7610a96565b5b600083013567ffffffffffffffff811115610cf657610cf5610a9b565b5b610d0285828601610bb8565b9250506020610d1385828601610cac565b9150509250929050565b610d2681610c2f565b8114610d3157600080fd5b50565b600081359050610d4381610d1d565b92915050565b60008060408385031215610d6057610d5f610a96565b5b600083013567ffffffffffffffff811115610d7e57610d7d610a9b565b5b610d8a85828601610bb8565b9250506020610d9b85828601610d34565b9150509250929050565b6000606082019050610dba6000830186610c39565b610dc76020830185610c39565b610dd46040830184610c39565b949350505050565b6000604082019050610df16000830185610c39565b610dfe6020830184610c39565b9392505050565b600081519050919050565b600081905092915050565b60005b83811015610e39578082015181840152602081019050610e1e565b60008484015250505050565b6000610e5082610e05565b610e5a8185610e10565b9350610e6a818560208601610e1b565b80840191505092915050565b6000610e828284610e45565b915081905092915050565b600081519050919050565b600082825260208201905092915050565b6000610eb482610e8d565b610ebe8185610e98565b9350610ece818560208601610e1b565b610ed781610aaa565b840191505092915050565b60006020820190508181036000830152610efc8184610ea9565b905092915050565b6000610f17610f1284610b36565b610b1b565b905082815260208101848484011115610f3357610f32610aa5565b5b610f3e848285610e1b565b509392505050565b600082601f830112610f5b57610f5a610aa0565b5b8151610f6b848260208601610f04565b91505092915050565b600060208284031215610f8a57610f89610a96565b5b600082015167ffffffffffffffff811115610fa857610fa7610a9b565b5b610fb484828501610f46565b91505092915050565b600081905092915050565b6000610fd382610e8d565b610fdd8185610fbd565b9350610fed818560208601610e1b565b80840191505092915050565b7f3a20000000000000000000000000000000000000000000000000000000000000600082015250565b600061102f600283610fbd565b915061103a82610ff9565b600282019050919050565b60006110518285610fc8565b915061105c82611022565b91506110688284610fc8565b91508190509392505050565b60008151905061108381610d1d565b92915050565b60006020828403121561109f5761109e610a96565b5b60006110ad84828501611074565b91505092915050565b6110bf81610c83565b82525050565b600060408201905081810360008301526110df8185610ea9565b90506110ee60208301846110b6565b9392505050565b6000604082019050818103600083015261110f8185610ea9565b905061111e6020830184610c39565b9392505050565b60008060006060848603121561113e5761113d610a96565b5b600061114c86828701611074565b935050602061115d86828701611074565b925050604061116e86828701611074565b9150509250925092565b6000806040838503121561118f5761118e610a96565b5b600061119d85828601611074565b92505060206111ae85828601611074565b915050925092905056fea2646970667358221220fd8b5f00f8ffbb5709f2f53fcda504bd8dc79bdf93978162a7e7db837e490e1364736f6c63430008130033",
}

// StakingTestABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingTestMetaData.ABI instead.
var StakingTestABI = StakingTestMetaData.ABI

// StakingTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingTestMetaData.Bin instead.
var StakingTestBin = StakingTestMetaData.Bin

// DeployStakingTest deploys a new Ethereum contract, binding an instance of StakingTest to it.
func DeployStakingTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingTest, error) {
	parsed, err := StakingTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingTest{StakingTestCaller: StakingTestCaller{contract: contract}, StakingTestTransactor: StakingTestTransactor{contract: contract}, StakingTestFilterer: StakingTestFilterer{contract: contract}}, nil
}

// StakingTest is an auto generated Go binding around an Ethereum contract.
type StakingTest struct {
	StakingTestCaller     // Read-only binding to the contract
	StakingTestTransactor // Write-only binding to the contract
	StakingTestFilterer   // Log filterer for contract events
}

// StakingTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingTestSession struct {
	Contract     *StakingTest      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingTestCallerSession struct {
	Contract *StakingTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakingTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTestTransactorSession struct {
	Contract     *StakingTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakingTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingTestRaw struct {
	Contract *StakingTest // Generic contract binding to access the raw methods on
}

// StakingTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingTestCallerRaw struct {
	Contract *StakingTestCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTestTransactorRaw struct {
	Contract *StakingTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingTest creates a new instance of StakingTest, bound to a specific deployed contract.
func NewStakingTest(address common.Address, backend bind.ContractBackend) (*StakingTest, error) {
	contract, err := bindStakingTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingTest{StakingTestCaller: StakingTestCaller{contract: contract}, StakingTestTransactor: StakingTestTransactor{contract: contract}, StakingTestFilterer: StakingTestFilterer{contract: contract}}, nil
}

// NewStakingTestCaller creates a new read-only instance of StakingTest, bound to a specific deployed contract.
func NewStakingTestCaller(address common.Address, caller bind.ContractCaller) (*StakingTestCaller, error) {
	contract, err := bindStakingTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTestCaller{contract: contract}, nil
}

// NewStakingTestTransactor creates a new write-only instance of StakingTest, bound to a specific deployed contract.
func NewStakingTestTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTestTransactor, error) {
	contract, err := bindStakingTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTestTransactor{contract: contract}, nil
}

// NewStakingTestFilterer creates a new log filterer instance of StakingTest, bound to a specific deployed contract.
func NewStakingTestFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingTestFilterer, error) {
	contract, err := bindStakingTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingTestFilterer{contract: contract}, nil
}

// bindStakingTest binds a generic wrapper to an already deployed contract.
func bindStakingTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingTestABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingTest *StakingTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingTest.Contract.StakingTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingTest *StakingTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingTest.Contract.StakingTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingTest *StakingTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingTest.Contract.StakingTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingTest *StakingTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingTest *StakingTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingTest *StakingTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingTest.Contract.contract.Transact(opts, method, params...)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestCaller) Delegation(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "delegation", _val, _del)
	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Delegation(_val string, _del common.Address) (*big.Int, *big.Int, error) {
	return _StakingTest.Contract.Delegation(&_StakingTest.CallOpts, _val, _del)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestCallerSession) Delegation(_val string, _del common.Address) (*big.Int, *big.Int, error) {
	return _StakingTest.Contract.Delegation(&_StakingTest.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestCaller) DelegationRewards(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "delegationRewards", _val, _del)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _StakingTest.Contract.DelegationRewards(&_StakingTest.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _StakingTest.Contract.DelegationRewards(&_StakingTest.CallOpts, _val, _del)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) Delegate(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "delegate", _val)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Delegate(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Delegate(&_StakingTest.TransactOpts, _val)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Delegate(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Delegate(&_StakingTest.TransactOpts, _val)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestTransactor) Undelegate(opts *bind.TransactOpts, _val string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "undelegate", _val, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestSession) Undelegate(_val string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Undelegate(&_StakingTest.TransactOpts, _val, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Undelegate(_val string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Undelegate(&_StakingTest.TransactOpts, _val, _shares)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256)
func (_StakingTest *StakingTestTransactor) Withdraw(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "withdraw", _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256)
func (_StakingTest *StakingTestSession) Withdraw(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Withdraw(&_StakingTest.TransactOpts, _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256)
func (_StakingTest *StakingTestTransactorSession) Withdraw(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Withdraw(&_StakingTest.TransactOpts, _val)
}
