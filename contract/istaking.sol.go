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

// IStakingMetaData contains all meta data concerning the IStaking contract.
var IStakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"ApproveShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valSrc\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valDst\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"RedelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"TransferShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"UndelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowanceShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approveShares\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"delegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_delegateAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_valSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_valDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"redelegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"slashingInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_jailed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_missed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFromShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"undelegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIStaking.ValidatorSortBy\",\"name\":\"_sortBy\",\"type\":\"uint8\"}],\"name\":\"validatorList\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use IStakingMetaData.ABI instead.
var IStakingABI = IStakingMetaData.ABI

// IStaking is an auto generated Go binding around an Ethereum contract.
type IStaking struct {
	IStakingCaller     // Read-only binding to the contract
	IStakingTransactor // Write-only binding to the contract
	IStakingFilterer   // Log filterer for contract events
}

// IStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type IStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IStakingSession struct {
	Contract     *IStaking         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IStakingCallerSession struct {
	Contract *IStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IStakingTransactorSession struct {
	Contract     *IStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type IStakingRaw struct {
	Contract *IStaking // Generic contract binding to access the raw methods on
}

// IStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IStakingCallerRaw struct {
	Contract *IStakingCaller // Generic read-only contract binding to access the raw methods on
}

// IStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IStakingTransactorRaw struct {
	Contract *IStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIStaking creates a new instance of IStaking, bound to a specific deployed contract.
func NewIStaking(address common.Address, backend bind.ContractBackend) (*IStaking, error) {
	contract, err := bindIStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IStaking{IStakingCaller: IStakingCaller{contract: contract}, IStakingTransactor: IStakingTransactor{contract: contract}, IStakingFilterer: IStakingFilterer{contract: contract}}, nil
}

// NewIStakingCaller creates a new read-only instance of IStaking, bound to a specific deployed contract.
func NewIStakingCaller(address common.Address, caller bind.ContractCaller) (*IStakingCaller, error) {
	contract, err := bindIStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IStakingCaller{contract: contract}, nil
}

// NewIStakingTransactor creates a new write-only instance of IStaking, bound to a specific deployed contract.
func NewIStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*IStakingTransactor, error) {
	contract, err := bindIStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IStakingTransactor{contract: contract}, nil
}

// NewIStakingFilterer creates a new log filterer instance of IStaking, bound to a specific deployed contract.
func NewIStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*IStakingFilterer, error) {
	contract, err := bindIStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IStakingFilterer{contract: contract}, nil
}

// bindIStaking binds a generic wrapper to an already deployed contract.
func bindIStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStaking *IStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStaking.Contract.IStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStaking *IStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStaking.Contract.IStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStaking *IStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStaking.Contract.IStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStaking *IStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStaking *IStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStaking *IStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStaking.Contract.contract.Transact(opts, method, params...)
}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256 _shares)
func (_IStaking *IStakingCaller) AllowanceShares(opts *bind.CallOpts, _val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "allowanceShares", _val, _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256 _shares)
func (_IStaking *IStakingSession) AllowanceShares(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _IStaking.Contract.AllowanceShares(&_IStaking.CallOpts, _val, _owner, _spender)
}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256 _shares)
func (_IStaking *IStakingCallerSession) AllowanceShares(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _IStaking.Contract.AllowanceShares(&_IStaking.CallOpts, _val, _owner, _spender)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256 _shares, uint256 _delegateAmount)
func (_IStaking *IStakingCaller) Delegation(opts *bind.CallOpts, _val string, _del common.Address) (struct {
	Shares         *big.Int
	DelegateAmount *big.Int
}, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "delegation", _val, _del)

	outstruct := new(struct {
		Shares         *big.Int
		DelegateAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Shares = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.DelegateAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256 _shares, uint256 _delegateAmount)
func (_IStaking *IStakingSession) Delegation(_val string, _del common.Address) (struct {
	Shares         *big.Int
	DelegateAmount *big.Int
}, error) {
	return _IStaking.Contract.Delegation(&_IStaking.CallOpts, _val, _del)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256 _shares, uint256 _delegateAmount)
func (_IStaking *IStakingCallerSession) Delegation(_val string, _del common.Address) (struct {
	Shares         *big.Int
	DelegateAmount *big.Int
}, error) {
	return _IStaking.Contract.Delegation(&_IStaking.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256 _reward)
func (_IStaking *IStakingCaller) DelegationRewards(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "delegationRewards", _val, _del)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256 _reward)
func (_IStaking *IStakingSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _IStaking.Contract.DelegationRewards(&_IStaking.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256 _reward)
func (_IStaking *IStakingCallerSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _IStaking.Contract.DelegationRewards(&_IStaking.CallOpts, _val, _del)
}

// SlashingInfo is a free data retrieval call binding the contract method 0x4e94633a.
//
// Solidity: function slashingInfo(string _val) view returns(bool _jailed, uint256 _missed)
func (_IStaking *IStakingCaller) SlashingInfo(opts *bind.CallOpts, _val string) (struct {
	Jailed bool
	Missed *big.Int
}, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "slashingInfo", _val)

	outstruct := new(struct {
		Jailed bool
		Missed *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Jailed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Missed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SlashingInfo is a free data retrieval call binding the contract method 0x4e94633a.
//
// Solidity: function slashingInfo(string _val) view returns(bool _jailed, uint256 _missed)
func (_IStaking *IStakingSession) SlashingInfo(_val string) (struct {
	Jailed bool
	Missed *big.Int
}, error) {
	return _IStaking.Contract.SlashingInfo(&_IStaking.CallOpts, _val)
}

// SlashingInfo is a free data retrieval call binding the contract method 0x4e94633a.
//
// Solidity: function slashingInfo(string _val) view returns(bool _jailed, uint256 _missed)
func (_IStaking *IStakingCallerSession) SlashingInfo(_val string) (struct {
	Jailed bool
	Missed *big.Int
}, error) {
	return _IStaking.Contract.SlashingInfo(&_IStaking.CallOpts, _val)
}

// ValidatorList is a free data retrieval call binding the contract method 0x029c0a51.
//
// Solidity: function validatorList(uint8 _sortBy) view returns(string[])
func (_IStaking *IStakingCaller) ValidatorList(opts *bind.CallOpts, _sortBy uint8) ([]string, error) {
	var out []interface{}
	err := _IStaking.contract.Call(opts, &out, "validatorList", _sortBy)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ValidatorList is a free data retrieval call binding the contract method 0x029c0a51.
//
// Solidity: function validatorList(uint8 _sortBy) view returns(string[])
func (_IStaking *IStakingSession) ValidatorList(_sortBy uint8) ([]string, error) {
	return _IStaking.Contract.ValidatorList(&_IStaking.CallOpts, _sortBy)
}

// ValidatorList is a free data retrieval call binding the contract method 0x029c0a51.
//
// Solidity: function validatorList(uint8 _sortBy) view returns(string[])
func (_IStaking *IStakingCallerSession) ValidatorList(_sortBy uint8) ([]string, error) {
	return _IStaking.Contract.ValidatorList(&_IStaking.CallOpts, _sortBy)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool _result)
func (_IStaking *IStakingTransactor) ApproveShares(opts *bind.TransactOpts, _val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "approveShares", _val, _spender, _shares)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool _result)
func (_IStaking *IStakingSession) ApproveShares(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.ApproveShares(&_IStaking.TransactOpts, _val, _spender, _shares)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool _result)
func (_IStaking *IStakingTransactorSession) ApproveShares(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.ApproveShares(&_IStaking.TransactOpts, _val, _spender, _shares)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) payable returns(bool _result)
func (_IStaking *IStakingTransactor) DelegateV2(opts *bind.TransactOpts, _val string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "delegateV2", _val, _amount)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) payable returns(bool _result)
func (_IStaking *IStakingSession) DelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.DelegateV2(&_IStaking.TransactOpts, _val, _amount)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) payable returns(bool _result)
func (_IStaking *IStakingTransactorSession) DelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.DelegateV2(&_IStaking.TransactOpts, _val, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_IStaking *IStakingTransactor) RedelegateV2(opts *bind.TransactOpts, _valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "redelegateV2", _valSrc, _valDst, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_IStaking *IStakingSession) RedelegateV2(_valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.RedelegateV2(&_IStaking.TransactOpts, _valSrc, _valDst, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_IStaking *IStakingTransactorSession) RedelegateV2(_valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.RedelegateV2(&_IStaking.TransactOpts, _valSrc, _valDst, _amount)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_IStaking *IStakingTransactor) TransferFromShares(opts *bind.TransactOpts, _val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "transferFromShares", _val, _from, _to, _shares)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_IStaking *IStakingSession) TransferFromShares(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.TransferFromShares(&_IStaking.TransactOpts, _val, _from, _to, _shares)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_IStaking *IStakingTransactorSession) TransferFromShares(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.TransferFromShares(&_IStaking.TransactOpts, _val, _from, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_IStaking *IStakingTransactor) TransferShares(opts *bind.TransactOpts, _val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "transferShares", _val, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_IStaking *IStakingSession) TransferShares(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.TransferShares(&_IStaking.TransactOpts, _val, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_IStaking *IStakingTransactorSession) TransferShares(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.TransferShares(&_IStaking.TransactOpts, _val, _to, _shares)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_IStaking *IStakingTransactor) UndelegateV2(opts *bind.TransactOpts, _val string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "undelegateV2", _val, _amount)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_IStaking *IStakingSession) UndelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.UndelegateV2(&_IStaking.TransactOpts, _val, _amount)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_IStaking *IStakingTransactorSession) UndelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _IStaking.Contract.UndelegateV2(&_IStaking.TransactOpts, _val, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256 _reward)
func (_IStaking *IStakingTransactor) Withdraw(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _IStaking.contract.Transact(opts, "withdraw", _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256 _reward)
func (_IStaking *IStakingSession) Withdraw(_val string) (*types.Transaction, error) {
	return _IStaking.Contract.Withdraw(&_IStaking.TransactOpts, _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256 _reward)
func (_IStaking *IStakingTransactorSession) Withdraw(_val string) (*types.Transaction, error) {
	return _IStaking.Contract.Withdraw(&_IStaking.TransactOpts, _val)
}

// IStakingApproveSharesIterator is returned from FilterApproveShares and is used to iterate over the raw logs and unpacked data for ApproveShares events raised by the IStaking contract.
type IStakingApproveSharesIterator struct {
	Event *IStakingApproveShares // Event containing the contract specifics and raw log

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
func (it *IStakingApproveSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingApproveShares)
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
		it.Event = new(IStakingApproveShares)
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
func (it *IStakingApproveSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingApproveSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingApproveShares represents a ApproveShares event raised by the IStaking contract.
type IStakingApproveShares struct {
	Owner     common.Address
	Spender   common.Address
	Validator string
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterApproveShares is a free log retrieval operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_IStaking *IStakingFilterer) FilterApproveShares(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IStakingApproveSharesIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "ApproveShares", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IStakingApproveSharesIterator{contract: _IStaking.contract, event: "ApproveShares", logs: logs, sub: sub}, nil
}

// WatchApproveShares is a free log subscription operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_IStaking *IStakingFilterer) WatchApproveShares(opts *bind.WatchOpts, sink chan<- *IStakingApproveShares, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "ApproveShares", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingApproveShares)
				if err := _IStaking.contract.UnpackLog(event, "ApproveShares", log); err != nil {
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

// ParseApproveShares is a log parse operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_IStaking *IStakingFilterer) ParseApproveShares(log types.Log) (*IStakingApproveShares, error) {
	event := new(IStakingApproveShares)
	if err := _IStaking.contract.UnpackLog(event, "ApproveShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingDelegateV2Iterator is returned from FilterDelegateV2 and is used to iterate over the raw logs and unpacked data for DelegateV2 events raised by the IStaking contract.
type IStakingDelegateV2Iterator struct {
	Event *IStakingDelegateV2 // Event containing the contract specifics and raw log

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
func (it *IStakingDelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingDelegateV2)
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
		it.Event = new(IStakingDelegateV2)
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
func (it *IStakingDelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingDelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingDelegateV2 represents a DelegateV2 event raised by the IStaking contract.
type IStakingDelegateV2 struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegateV2 is a free log retrieval operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_IStaking *IStakingFilterer) FilterDelegateV2(opts *bind.FilterOpts, delegator []common.Address) (*IStakingDelegateV2Iterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "DelegateV2", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &IStakingDelegateV2Iterator{contract: _IStaking.contract, event: "DelegateV2", logs: logs, sub: sub}, nil
}

// WatchDelegateV2 is a free log subscription operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_IStaking *IStakingFilterer) WatchDelegateV2(opts *bind.WatchOpts, sink chan<- *IStakingDelegateV2, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "DelegateV2", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingDelegateV2)
				if err := _IStaking.contract.UnpackLog(event, "DelegateV2", log); err != nil {
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

// ParseDelegateV2 is a log parse operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_IStaking *IStakingFilterer) ParseDelegateV2(log types.Log) (*IStakingDelegateV2, error) {
	event := new(IStakingDelegateV2)
	if err := _IStaking.contract.UnpackLog(event, "DelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingRedelegateV2Iterator is returned from FilterRedelegateV2 and is used to iterate over the raw logs and unpacked data for RedelegateV2 events raised by the IStaking contract.
type IStakingRedelegateV2Iterator struct {
	Event *IStakingRedelegateV2 // Event containing the contract specifics and raw log

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
func (it *IStakingRedelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingRedelegateV2)
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
		it.Event = new(IStakingRedelegateV2)
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
func (it *IStakingRedelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingRedelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingRedelegateV2 represents a RedelegateV2 event raised by the IStaking contract.
type IStakingRedelegateV2 struct {
	Sender         common.Address
	ValSrc         string
	ValDst         string
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRedelegateV2 is a free log retrieval operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) FilterRedelegateV2(opts *bind.FilterOpts, sender []common.Address) (*IStakingRedelegateV2Iterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "RedelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return &IStakingRedelegateV2Iterator{contract: _IStaking.contract, event: "RedelegateV2", logs: logs, sub: sub}, nil
}

// WatchRedelegateV2 is a free log subscription operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) WatchRedelegateV2(opts *bind.WatchOpts, sink chan<- *IStakingRedelegateV2, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "RedelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingRedelegateV2)
				if err := _IStaking.contract.UnpackLog(event, "RedelegateV2", log); err != nil {
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

// ParseRedelegateV2 is a log parse operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) ParseRedelegateV2(log types.Log) (*IStakingRedelegateV2, error) {
	event := new(IStakingRedelegateV2)
	if err := _IStaking.contract.UnpackLog(event, "RedelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingTransferSharesIterator is returned from FilterTransferShares and is used to iterate over the raw logs and unpacked data for TransferShares events raised by the IStaking contract.
type IStakingTransferSharesIterator struct {
	Event *IStakingTransferShares // Event containing the contract specifics and raw log

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
func (it *IStakingTransferSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingTransferShares)
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
		it.Event = new(IStakingTransferShares)
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
func (it *IStakingTransferSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingTransferSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingTransferShares represents a TransferShares event raised by the IStaking contract.
type IStakingTransferShares struct {
	From      common.Address
	To        common.Address
	Validator string
	Shares    *big.Int
	Token     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferShares is a free log retrieval operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_IStaking *IStakingFilterer) FilterTransferShares(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IStakingTransferSharesIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "TransferShares", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IStakingTransferSharesIterator{contract: _IStaking.contract, event: "TransferShares", logs: logs, sub: sub}, nil
}

// WatchTransferShares is a free log subscription operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_IStaking *IStakingFilterer) WatchTransferShares(opts *bind.WatchOpts, sink chan<- *IStakingTransferShares, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "TransferShares", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingTransferShares)
				if err := _IStaking.contract.UnpackLog(event, "TransferShares", log); err != nil {
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

// ParseTransferShares is a log parse operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_IStaking *IStakingFilterer) ParseTransferShares(log types.Log) (*IStakingTransferShares, error) {
	event := new(IStakingTransferShares)
	if err := _IStaking.contract.UnpackLog(event, "TransferShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingUndelegateV2Iterator is returned from FilterUndelegateV2 and is used to iterate over the raw logs and unpacked data for UndelegateV2 events raised by the IStaking contract.
type IStakingUndelegateV2Iterator struct {
	Event *IStakingUndelegateV2 // Event containing the contract specifics and raw log

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
func (it *IStakingUndelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingUndelegateV2)
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
		it.Event = new(IStakingUndelegateV2)
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
func (it *IStakingUndelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingUndelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingUndelegateV2 represents a UndelegateV2 event raised by the IStaking contract.
type IStakingUndelegateV2 struct {
	Sender         common.Address
	Validator      string
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUndelegateV2 is a free log retrieval operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) FilterUndelegateV2(opts *bind.FilterOpts, sender []common.Address) (*IStakingUndelegateV2Iterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "UndelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return &IStakingUndelegateV2Iterator{contract: _IStaking.contract, event: "UndelegateV2", logs: logs, sub: sub}, nil
}

// WatchUndelegateV2 is a free log subscription operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) WatchUndelegateV2(opts *bind.WatchOpts, sink chan<- *IStakingUndelegateV2, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "UndelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingUndelegateV2)
				if err := _IStaking.contract.UnpackLog(event, "UndelegateV2", log); err != nil {
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

// ParseUndelegateV2 is a log parse operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_IStaking *IStakingFilterer) ParseUndelegateV2(log types.Log) (*IStakingUndelegateV2, error) {
	event := new(IStakingUndelegateV2)
	if err := _IStaking.contract.UnpackLog(event, "UndelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IStakingWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the IStaking contract.
type IStakingWithdrawIterator struct {
	Event *IStakingWithdraw // Event containing the contract specifics and raw log

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
func (it *IStakingWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IStakingWithdraw)
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
		it.Event = new(IStakingWithdraw)
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
func (it *IStakingWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IStakingWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IStakingWithdraw represents a Withdraw event raised by the IStaking contract.
type IStakingWithdraw struct {
	Sender    common.Address
	Validator string
	Reward    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_IStaking *IStakingFilterer) FilterWithdraw(opts *bind.FilterOpts, sender []common.Address) (*IStakingWithdrawIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IStaking.contract.FilterLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return &IStakingWithdrawIterator{contract: _IStaking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_IStaking *IStakingFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *IStakingWithdraw, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IStaking.contract.WatchLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IStakingWithdraw)
				if err := _IStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_IStaking *IStakingFilterer) ParseWithdraw(log types.Log) (*IStakingWithdraw, error) {
	event := new(IStakingWithdraw)
	if err := _IStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
