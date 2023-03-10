// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking_test

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

// StakingTestABI is the input ABI used to generate the binding from.
const StakingTestABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_acc\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amt\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StakingTestBin is the compiled bytecode used for deploying new contracts.
var StakingTestBin = "0x608060405234801561001057600080fd5b50610a11806100206000396000f3fe60806040526004361061004a5760003560e01c806303f24de11461004f57806331fb67c2146100755780638dfc889714610095578063bf98d772146100ca578063d5c498eb14610102575b600080fd5b61006261005d366004610714565b610122565b6040519081526020015b60405180910390f35b34801561008157600080fd5b50610062610090366004610759565b61016e565b3480156100a157600080fd5b506100b56100b0366004610714565b610179565b6040805192835260208301919091520161006c565b3480156100d657600080fd5b506100626100e5366004610759565b805160208183018101805160008252928201919093012091525481565b34801561010e57600080fd5b5061006261011d36600461078e565b6101cb565b60008061012f84846101de565b9050806000856040516101429190610810565b9081526020016040518091039020600082825461015f9190610842565b90915550909150505b92915050565b600061016882610287565b600080600080610189868661032d565b915091508460008760405161019e9190610810565b908152602001604051809103902060008282546101bb9190610855565b9091555091969095509350505050565b60006101d783836103dc565b9392505050565b600080806064846101ef8782610472565b6040516101fc9190610810565b60006040518083038185875af1925050503d8060008114610239576040519150601f19603f3d011682016040523d82523d6000602084013e61023e565b606091505b509150915061027582826040518060400160405280600f81526020016e19195b1959d85d194819985a5b1959608a1b8152506104b9565b61027e81610541565b95945050505050565b60008080606461029685610558565b6040516102a39190610810565b6000604051808303816000865af19150503d80600081146102e0576040519150601f19603f3d011682016040523d82523d6000602084013e6102e5565b606091505b509150915061031c82826040518060400160405280600f81526020016e1dda5d1a191c985dc819985a5b1959608a1b8152506104b9565b61032581610541565b949350505050565b6000808080606461033e878761059b565b60405161034b9190610810565b6000604051808303816000865af19150503d8060008114610388576040519150601f19603f3d011682016040523d82523d6000602084013e61038d565b606091505b50915091506103c68282604051806040016040528060118152602001701d5b99195b1959d85d194819985a5b1959607a1b8152506104b9565b6103cf816105e2565b9350935050509250929050565b6000808060646103ec8686610608565b6040516103f99190610810565b600060405180830381855afa9150503d8060008114610434576040519150601f19603f3d011682016040523d82523d6000602084013e610439565b606091505b509150915061027582826040518060400160405280601181526020017019195b1959d85d1a5bdb8819985a5b1959607a1b8152506104b9565b60608282604051602401610487929190610894565b60408051601f198184030181529190526020810180516001600160e01b03166303f24de160e01b179052905092915050565b8261053c576000828060200190518101906104d491906108b6565b9050600182511015610503578060405162461bcd60e51b81526004016104fa9190610924565b60405180910390fd5b8181604051602001610516929190610937565b60408051601f198184030181529082905262461bcd60e51b82526104fa91600401610924565b505050565b600080828060200190518101906101d79190610974565b60608160405160240161056b9190610924565b60408051601f198184030181529190526020810180516001600160e01b03166318fdb3e160e11b17905292915050565b606082826040516024016105b0929190610894565b60408051601f198184030181529190526020810180516001600160e01b0316638dfc889760e01b179052905092915050565b600080600080848060200190518101906105fc919061098d565b90969095509350505050565b6060828260405160240161061d9291906109b1565b60408051601f198184030181529190526020810180516001600160e01b031663d5c498eb60e01b179052905092915050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561068e5761068e61064f565b604052919050565b600067ffffffffffffffff8211156106b0576106b061064f565b50601f01601f191660200190565b600082601f8301126106cf57600080fd5b81356106e26106dd82610696565b610665565b8181528460208386010111156106f757600080fd5b816020850160208301376000918101602001919091529392505050565b6000806040838503121561072757600080fd5b823567ffffffffffffffff81111561073e57600080fd5b61074a858286016106be565b95602094909401359450505050565b60006020828403121561076b57600080fd5b813567ffffffffffffffff81111561078257600080fd5b610325848285016106be565b600080604083850312156107a157600080fd5b823567ffffffffffffffff8111156107b857600080fd5b6107c4858286016106be565b92505060208301356001600160a01b03811681146107e157600080fd5b809150509250929050565b60005b838110156108075781810151838201526020016107ef565b50506000910152565b600082516108228184602087016107ec565b9190910192915050565b634e487b7160e01b600052601160045260246000fd5b808201808211156101685761016861082c565b818103818111156101685761016861082c565b600081518084526108808160208601602086016107ec565b601f01601f19169290920160200192915050565b6040815260006108a76040830185610868565b90508260208301529392505050565b6000602082840312156108c857600080fd5b815167ffffffffffffffff8111156108df57600080fd5b8201601f810184136108f057600080fd5b80516108fe6106dd82610696565b81815285602083850101111561091357600080fd5b61027e8260208301602086016107ec565b6020815260006101d76020830184610868565b600083516109498184602088016107ec565b6101d160f51b90830190815283516109688160028401602088016107ec565b01600201949350505050565b60006020828403121561098657600080fd5b5051919050565b600080604083850312156109a057600080fd5b505080516020909101519092909150565b6040815260006109c46040830185610868565b905060018060a01b0383166020830152939250505056fea264697066735822122012f07c3985898b09913966a54466d0d458938c48c6e76a0400814bfe38dac77064736f6c63430008130033"

// DeployStakingTest deploys a new Ethereum contract, binding an instance of StakingTest to it.
func DeployStakingTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingTest, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingTestABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StakingTestBin), backend)
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

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _acc) view returns(uint256)
func (_StakingTest *StakingTestCaller) BalanceOf(opts *bind.CallOpts, _acc common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "balanceOf", _acc)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _acc) view returns(uint256)
func (_StakingTest *StakingTestSession) BalanceOf(_acc common.Address) (*big.Int, error) {
	return _StakingTest.Contract.BalanceOf(&_StakingTest.CallOpts, _acc)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _acc) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) BalanceOf(_acc common.Address) (*big.Int, error) {
	return _StakingTest.Contract.BalanceOf(&_StakingTest.CallOpts, _acc)
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

// Delegate is a paid mutator transaction binding the contract method 0x03f24de1.
//
// Solidity: function delegate(string _val, uint256 _amt) payable returns(uint256)
func (_StakingTest *StakingTestTransactor) Delegate(opts *bind.TransactOpts, _val string, _amt *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "delegate", _val, _amt)
}

// Delegate is a paid mutator transaction binding the contract method 0x03f24de1.
//
// Solidity: function delegate(string _val, uint256 _amt) payable returns(uint256)
func (_StakingTest *StakingTestSession) Delegate(_val string, _amt *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Delegate(&_StakingTest.TransactOpts, _val, _amt)
}

// Delegate is a paid mutator transaction binding the contract method 0x03f24de1.
//
// Solidity: function delegate(string _val, uint256 _amt) payable returns(uint256)
func (_StakingTest *StakingTestTransactorSession) Delegate(_val string, _amt *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Delegate(&_StakingTest.TransactOpts, _val, _amt)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) Undelegate(opts *bind.TransactOpts, _val string, shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "undelegate", _val, shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Undelegate(_val string, shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Undelegate(&_StakingTest.TransactOpts, _val, shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Undelegate(_val string, shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Undelegate(&_StakingTest.TransactOpts, _val, shares)
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
