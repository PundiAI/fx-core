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
const StakingTestABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"validatorShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// StakingTestBin is the compiled bytecode used for deploying new contracts.
var StakingTestBin = "0x608060405234801561001057600080fd5b50610aaf806100206000396000f3fe60806040526004361061004a5760003560e01c806331fb67c21461004f5780638dfc8897146100825780639ddb511a146100bd578063bf98d772146100e5578063d5c498eb1461011d575b600080fd5b34801561005b57600080fd5b5061006f61006a366004610784565b61013d565b6040519081526020015b60405180910390f35b34801561008e57600080fd5b506100a261009d3660046107b9565b61014e565b60408051938452602084019290925290820152606001610079565b6100d06100cb366004610784565b6101a7565b60408051928352602083019190915201610079565b3480156100f157600080fd5b5061006f610100366004610784565b805160208183018101805160008252928201919093012091525481565b34801561012957600080fd5b5061006f6101383660046107fe565b6101f7565b60006101488261020a565b92915050565b60008060008060008061016188886102b1565b925092509250866000896040516101789190610880565b9081526020016040518091039020600082825461019591906108b2565b90915550929891975095509350505050565b6000806000806101b685610364565b91509150816000866040516101cb9190610880565b908152602001604051809103902060008282546101e891906108c5565b90915550919590945092505050565b60006102038383610410565b9392505050565b6000808061100361021a856104b9565b6040516102279190610880565b6000604051808303816000865af19150503d8060008114610264576040519150601f19603f3d011682016040523d82523d6000602084013e610269565b606091505b50915091506102a082826040518060400160405280600f81526020016e1dda5d1a191c985dc819985a5b1959608a1b8152506104fc565b6102a981610584565b949350505050565b6000808080806110036102c4888861059b565b6040516102d19190610880565b6000604051808303816000865af19150503d806000811461030e576040519150601f19603f3d011682016040523d82523d6000602084013e610313565b606091505b509150915061034c8282604051806040016040528060118152602001701d5b99195b1959d85d194819985a5b1959607a1b8152506104fc565b610355816105e2565b94509450945050509250925092565b6000808080611003346103768761060f565b6040516103839190610880565b60006040518083038185875af1925050503d80600081146103c0576040519150601f19603f3d011682016040523d82523d6000602084013e6103c5565b606091505b50915091506103fc82826040518060400160405280600f81526020016e19195b1959d85d194819985a5b1959608a1b8152506104fc565b61040581610652565b935093505050915091565b600080806110036104218686610678565b60405161042e9190610880565b600060405180830381855afa9150503d8060008114610469576040519150601f19603f3d011682016040523d82523d6000602084013e61046e565b606091505b50915091506104a782826040518060400160405280601181526020017019195b1959d85d1a5bdb8819985a5b1959607a1b8152506104fc565b6104b081610584565b95945050505050565b6060816040516024016104cc9190610904565b60408051601f198184030181529190526020810180516001600160e01b03166318fdb3e160e11b17905292915050565b8261057f576000828060200190518101906105179190610917565b9050600182511015610546578060405162461bcd60e51b815260040161053d9190610904565b60405180910390fd5b8181604051602001610559929190610985565b60408051601f198184030181529082905262461bcd60e51b825261053d91600401610904565b505050565b6000808280602001905181019061020391906109c2565b606082826040516024016105b09291906109db565b60408051601f198184030181529190526020810180516001600160e01b0316638dfc889760e01b179052905092915050565b600080600080600080868060200190518101906105ff91906109fd565b9199909850909650945050505050565b6060816040516024016106229190610904565b60408051601f198184030181529190526020810180516001600160e01b0316634eeda88d60e11b17905292915050565b6000806000808480602001905181019061066c9190610a2b565b90969095509350505050565b6060828260405160240161068d929190610a4f565b60408051601f198184030181529190526020810180516001600160e01b031663d5c498eb60e01b179052905092915050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156106fe576106fe6106bf565b604052919050565b600067ffffffffffffffff821115610720576107206106bf565b50601f01601f191660200190565b600082601f83011261073f57600080fd5b813561075261074d82610706565b6106d5565b81815284602083860101111561076757600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561079657600080fd5b813567ffffffffffffffff8111156107ad57600080fd5b6102a98482850161072e565b600080604083850312156107cc57600080fd5b823567ffffffffffffffff8111156107e357600080fd5b6107ef8582860161072e565b95602094909401359450505050565b6000806040838503121561081157600080fd5b823567ffffffffffffffff81111561082857600080fd5b6108348582860161072e565b92505060208301356001600160a01b038116811461085157600080fd5b809150509250929050565b60005b8381101561087757818101518382015260200161085f565b50506000910152565b6000825161089281846020870161085c565b9190910192915050565b634e487b7160e01b600052601160045260246000fd5b818103818111156101485761014861089c565b808201808211156101485761014861089c565b600081518084526108f081602086016020860161085c565b601f01601f19169290920160200192915050565b60208152600061020360208301846108d8565b60006020828403121561092957600080fd5b815167ffffffffffffffff81111561094057600080fd5b8201601f8101841361095157600080fd5b805161095f61074d82610706565b81815285602083850101111561097457600080fd5b6104b082602083016020860161085c565b6000835161099781846020880161085c565b6101d160f51b90830190815283516109b681600284016020880161085c565b01600201949350505050565b6000602082840312156109d457600080fd5b5051919050565b6040815260006109ee60408301856108d8565b90508260208301529392505050565b600080600060608486031215610a1257600080fd5b8351925060208401519150604084015190509250925092565b60008060408385031215610a3e57600080fd5b505080516020909101519092909150565b604081526000610a6260408301856108d8565b905060018060a01b0383166020830152939250505056fea2646970667358221220a050d46d590ec70f7735a06dfcfd83e1a707b8ad2196405198e1580202c5ea6a64736f6c63430008130033"

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

// ValidatorShares is a free data retrieval call binding the contract method 0xbf98d772.
//
// Solidity: function validatorShares(string ) view returns(uint256)
func (_StakingTest *StakingTestCaller) ValidatorShares(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "validatorShares", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// ValidatorShares is a free data retrieval call binding the contract method 0xbf98d772.
//
// Solidity: function validatorShares(string ) view returns(uint256)
func (_StakingTest *StakingTestSession) ValidatorShares(arg0 string) (*big.Int, error) {
	return _StakingTest.Contract.ValidatorShares(&_StakingTest.CallOpts, arg0)
}

// ValidatorShares is a free data retrieval call binding the contract method 0xbf98d772.
//
// Solidity: function validatorShares(string ) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) ValidatorShares(arg0 string) (*big.Int, error) {
	return _StakingTest.Contract.ValidatorShares(&_StakingTest.CallOpts, arg0)
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
