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
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111dc806100206000396000f3fe60806040526004361061004a5760003560e01c806331fb67c21461004f57806351af513a1461008c5780638dfc8897146100c95780639ddb511a14610108578063d5c498eb14610139575b600080fd5b34801561005b57600080fd5b5061007660048036038101906100719190610bd4565b610176565b6040516100839190610c36565b60405180910390f35b34801561009857600080fd5b506100b360048036038101906100ae9190610caf565b610188565b6040516100c09190610c36565b60405180910390f35b3480156100d557600080fd5b506100f060048036038101906100eb9190610d37565b61019c565b6040516100ff93929190610d93565b60405180910390f35b610122600480360381019061011d9190610bd4565b6101b8565b604051610130929190610dca565b60405180910390f35b34801561014557600080fd5b50610160600480360381019061015b9190610caf565b6101cd565b60405161016d9190610c36565b60405180910390f35b6000610181826101e1565b9050919050565b600061019483836102ae565b905092915050565b60008060006101ab858561037b565b9250925092509250925092565b6000806101c483610452565b91509150915091565b60006101d98383610523565b905092915050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff16610208856105f0565b6040516102159190610e64565b6000604051808303816000865af19150503d8060008114610252576040519150601f19603f3d011682016040523d82523d6000602084013e610257565b606091505b509150915061029c82826040518060400160405280600f81526020017f7769746864726177206661696c65640000000000000000000000000000000000815250610687565b6102a58161074e565b92505050919050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff166102d68686610770565b6040516102e39190610e64565b600060405180830381855afa9150503d806000811461031e576040519150601f19603f3d011682016040523d82523d6000602084013e610323565b606091505b509150915061036882826040518060400160405280601881526020017f64656c65676174696f6e52657761726473206661696c65640000000000000000815250610687565b6103718161080a565b9250505092915050565b600080600080600061100373ffffffffffffffffffffffffffffffffffffffff166103a6888861082c565b6040516103b39190610e64565b6000604051808303816000865af19150503d80600081146103f0576040519150601f19603f3d011682016040523d82523d6000602084013e6103f5565b606091505b509150915061043a82826040518060400160405280601181526020017f756e64656c6567617465206661696c6564000000000000000000000000000000815250610687565b610443816108c6565b94509450945050509250925092565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff163461047b876108fc565b6040516104889190610e64565b60006040518083038185875af1925050503d80600081146104c5576040519150601f19603f3d011682016040523d82523d6000602084013e6104ca565b606091505b509150915061050f82826040518060400160405280600f81526020017f64656c6567617465206661696c65640000000000000000000000000000000000815250610687565b61051881610993565b935093505050915091565b600080600061100373ffffffffffffffffffffffffffffffffffffffff1661054b86866109be565b6040516105589190610e64565b600060405180830381855afa9150503d8060008114610593576040519150601f19603f3d011682016040523d82523d6000602084013e610598565b606091505b50915091506105dd82826040518060400160405280601181526020017f64656c65676174696f6e206661696c6564000000000000000000000000000000815250610687565b6105e681610a58565b9250505092915050565b6060816040516024016106039190610ed0565b6040516020818303038152906040527f31fb67c2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b82610749576000828060200190518101906106a29190610f62565b90506001825110156106eb57806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106e29190610ed0565b60405180910390fd5b81816040516020016106fe929190611033565b6040516020818303038152906040526040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107409190610ed0565b60405180910390fd5b505050565b600080828060200190518101906107659190611077565b905080915050919050565b606082826040516024016107859291906110b3565b6040516020818303038152906040527f51af513a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b600080828060200190518101906108219190611077565b905080915050919050565b606082826040516024016108419291906110e3565b6040516020818303038152906040527f8dfc8897000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b600080600080600080868060200190518101906108e39190611113565b9250925092508282829550955095505050509193909250565b60608160405160240161090f9190610ed0565b6040516020818303038152906040527f9ddb511a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b600080600080848060200190518101906109ad9190611166565b915091508181935093505050915091565b606082826040516024016109d39291906110b3565b6040516020818303038152906040527fd5c498eb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b60008082806020019051810190610a6f9190611077565b905080915050919050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610ae182610a98565b810181811067ffffffffffffffff82111715610b0057610aff610aa9565b5b80604052505050565b6000610b13610a7a565b9050610b1f8282610ad8565b919050565b600067ffffffffffffffff821115610b3f57610b3e610aa9565b5b610b4882610a98565b9050602081019050919050565b82818337600083830152505050565b6000610b77610b7284610b24565b610b09565b905082815260208101848484011115610b9357610b92610a93565b5b610b9e848285610b55565b509392505050565b600082601f830112610bbb57610bba610a8e565b5b8135610bcb848260208601610b64565b91505092915050565b600060208284031215610bea57610be9610a84565b5b600082013567ffffffffffffffff811115610c0857610c07610a89565b5b610c1484828501610ba6565b91505092915050565b6000819050919050565b610c3081610c1d565b82525050565b6000602082019050610c4b6000830184610c27565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610c7c82610c51565b9050919050565b610c8c81610c71565b8114610c9757600080fd5b50565b600081359050610ca981610c83565b92915050565b60008060408385031215610cc657610cc5610a84565b5b600083013567ffffffffffffffff811115610ce457610ce3610a89565b5b610cf085828601610ba6565b9250506020610d0185828601610c9a565b9150509250929050565b610d1481610c1d565b8114610d1f57600080fd5b50565b600081359050610d3181610d0b565b92915050565b60008060408385031215610d4e57610d4d610a84565b5b600083013567ffffffffffffffff811115610d6c57610d6b610a89565b5b610d7885828601610ba6565b9250506020610d8985828601610d22565b9150509250929050565b6000606082019050610da86000830186610c27565b610db56020830185610c27565b610dc26040830184610c27565b949350505050565b6000604082019050610ddf6000830185610c27565b610dec6020830184610c27565b9392505050565b600081519050919050565b600081905092915050565b60005b83811015610e27578082015181840152602081019050610e0c565b60008484015250505050565b6000610e3e82610df3565b610e488185610dfe565b9350610e58818560208601610e09565b80840191505092915050565b6000610e708284610e33565b915081905092915050565b600081519050919050565b600082825260208201905092915050565b6000610ea282610e7b565b610eac8185610e86565b9350610ebc818560208601610e09565b610ec581610a98565b840191505092915050565b60006020820190508181036000830152610eea8184610e97565b905092915050565b6000610f05610f0084610b24565b610b09565b905082815260208101848484011115610f2157610f20610a93565b5b610f2c848285610e09565b509392505050565b600082601f830112610f4957610f48610a8e565b5b8151610f59848260208601610ef2565b91505092915050565b600060208284031215610f7857610f77610a84565b5b600082015167ffffffffffffffff811115610f9657610f95610a89565b5b610fa284828501610f34565b91505092915050565b600081905092915050565b6000610fc182610e7b565b610fcb8185610fab565b9350610fdb818560208601610e09565b80840191505092915050565b7f3a20000000000000000000000000000000000000000000000000000000000000600082015250565b600061101d600283610fab565b915061102882610fe7565b600282019050919050565b600061103f8285610fb6565b915061104a82611010565b91506110568284610fb6565b91508190509392505050565b60008151905061107181610d0b565b92915050565b60006020828403121561108d5761108c610a84565b5b600061109b84828501611062565b91505092915050565b6110ad81610c71565b82525050565b600060408201905081810360008301526110cd8185610e97565b90506110dc60208301846110a4565b9392505050565b600060408201905081810360008301526110fd8185610e97565b905061110c6020830184610c27565b9392505050565b60008060006060848603121561112c5761112b610a84565b5b600061113a86828701611062565b935050602061114b86828701611062565b925050604061115c86828701611062565b9150509250925092565b6000806040838503121561117d5761117c610a84565b5b600061118b85828601611062565b925050602061119c85828601611062565b915050925092905056fea2646970667358221220f2a254f67ca17029d03b425e7a4bf8e247bf634bd24b0794ef204493aa6753c064736f6c63430008130033",
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
// Solidity: function delegation(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestCaller) Delegation(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "delegation", _val, _del)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestSession) Delegation(_val string, _del common.Address) (*big.Int, error) {
	return _StakingTest.Contract.Delegation(&_StakingTest.CallOpts, _val, _del)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) Delegation(_val string, _del common.Address) (*big.Int, error) {
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
