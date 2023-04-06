// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking

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

// StakingMetaData contains all meta data concerning the Staking contract.
var StakingMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Approve\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Undelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_delegateAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_token\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_completionTime\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610c2b806100206000396000f3fe6080604052600436106100865760003560e01c80638dfc8897116100595780638dfc8897146101805780639ddb511a146101bf578063aa2fddc7146101f0578063bd5b93d21461022e578063d5c498eb1461026b57610086565b80631fc38a8b1461008b57806331fb67c2146100c9578063495aff711461010657806351af513a14610143575b600080fd5b34801561009757600080fd5b506100b260048036038101906100ad9190610755565b6102a9565b6040516100c09291906107e7565b60405180910390f35b3480156100d557600080fd5b506100f060048036038101906100eb9190610810565b610326565b6040516100fd9190610859565b60405180910390f35b34801561011257600080fd5b5061012d60048036038101906101289190610874565b610382565b60405161013a9190610859565b60405180910390f35b34801561014f57600080fd5b5061016a600480360381019061016591906108e3565b61038b565b6040516101779190610859565b60405180910390f35b34801561018c57600080fd5b506101a760048036038101906101a2919061093f565b610393565b6040516101b69392919061099b565b60405180910390f35b6101d960048036038101906101d49190610810565b610400565b6040516101e79291906107e7565b60405180910390f35b3480156101fc57600080fd5b50610217600480360381019061021291906109d2565b610463565b6040516102259291906107e7565b60405180910390f35b34801561023a57600080fd5b50610255600480360381019061025091906109d2565b6104e0565b6040516102629190610a5c565b60405180910390f35b34801561027757600080fd5b50610292600480360381019061028d91906108e3565b610556565b6040516102a09291906107e7565b60405180910390f35b6000808373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167f39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a8138860008060405161030e93929190610b3b565b60405180910390a36000809150915094509492505050565b60008073ffffffffffffffffffffffffffffffffffffffff167f901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef836000604051610371929190610b79565b60405180910390a260009050919050565b60009392505050565b600092915050565b60008060008073ffffffffffffffffffffffffffffffffffffffff167fadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f8660008060006040516103e69493929190610ba9565b60405180910390a260008060009250925092509250925092565b600080600073ffffffffffffffffffffffffffffffffffffffff167f5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a8460008060405161044f93929190610b3b565b60405180910390a260008091509150915091565b6000808373ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813876000806040516104c993929190610b3b565b60405180910390a360008091509150935093915050565b60008273ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020866000604051610543929190610b79565b60405180910390a3600090509392505050565b600080600080915091509250929050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105ce82610585565b810181811067ffffffffffffffff821117156105ed576105ec610596565b5b80604052505050565b6000610600610567565b905061060c82826105c5565b919050565b600067ffffffffffffffff82111561062c5761062b610596565b5b61063582610585565b9050602081019050919050565b82818337600083830152505050565b600061066461065f84610611565b6105f6565b9050828152602081018484840111156106805761067f610580565b5b61068b848285610642565b509392505050565b600082601f8301126106a8576106a761057b565b5b81356106b8848260208601610651565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106ec826106c1565b9050919050565b6106fc816106e1565b811461070757600080fd5b50565b600081359050610719816106f3565b92915050565b6000819050919050565b6107328161071f565b811461073d57600080fd5b50565b60008135905061074f81610729565b92915050565b6000806000806080858703121561076f5761076e610571565b5b600085013567ffffffffffffffff81111561078d5761078c610576565b5b61079987828801610693565b94505060206107aa8782880161070a565b93505060406107bb8782880161070a565b92505060606107cc87828801610740565b91505092959194509250565b6107e18161071f565b82525050565b60006040820190506107fc60008301856107d8565b61080960208301846107d8565b9392505050565b60006020828403121561082657610825610571565b5b600082013567ffffffffffffffff81111561084457610843610576565b5b61085084828501610693565b91505092915050565b600060208201905061086e60008301846107d8565b92915050565b60008060006060848603121561088d5761088c610571565b5b600084013567ffffffffffffffff8111156108ab576108aa610576565b5b6108b786828701610693565b93505060206108c88682870161070a565b92505060406108d98682870161070a565b9150509250925092565b600080604083850312156108fa576108f9610571565b5b600083013567ffffffffffffffff81111561091857610917610576565b5b61092485828601610693565b92505060206109358582860161070a565b9150509250929050565b6000806040838503121561095657610955610571565b5b600083013567ffffffffffffffff81111561097457610973610576565b5b61098085828601610693565b925050602061099185828601610740565b9150509250929050565b60006060820190506109b060008301866107d8565b6109bd60208301856107d8565b6109ca60408301846107d8565b949350505050565b6000806000606084860312156109eb576109ea610571565b5b600084013567ffffffffffffffff811115610a0957610a08610576565b5b610a1586828701610693565b9350506020610a268682870161070a565b9250506040610a3786828701610740565b9150509250925092565b60008115159050919050565b610a5681610a41565b82525050565b6000602082019050610a716000830184610a4d565b92915050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610ab1578082015181840152602081019050610a96565b60008484015250505050565b6000610ac882610a77565b610ad28185610a82565b9350610ae2818560208601610a93565b610aeb81610585565b840191505092915050565b6000819050919050565b6000819050919050565b6000610b25610b20610b1b84610af6565b610b00565b61071f565b9050919050565b610b3581610b0a565b82525050565b60006060820190508181036000830152610b558186610abd565b9050610b646020830185610b2c565b610b716040830184610b2c565b949350505050565b60006040820190508181036000830152610b938185610abd565b9050610ba26020830184610b2c565b9392505050565b60006080820190508181036000830152610bc38187610abd565b9050610bd26020830186610b2c565b610bdf6040830185610b2c565b610bec6060830184610b2c565b9594505050505056fea264697066735822122057571f81301a9779102f32dab266182b767913665eed32a7939dcc5068b0216b64736f6c63430008130033",
}

// StakingABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingMetaData.ABI instead.
var StakingABI = StakingMetaData.ABI

// StakingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingMetaData.Bin instead.
var StakingBin = StakingMetaData.Bin

// DeployStaking deploys a new Ethereum contract, binding an instance of Staking to it.
func DeployStaking(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Staking, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0x495aff71.
//
// Solidity: function allowance(string _val, address _owner, address _spender) view returns(uint256 _shares)
func (_Staking *StakingCaller) Allowance(opts *bind.CallOpts, _val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "allowance", _val, _owner, _spender)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0x495aff71.
//
// Solidity: function allowance(string _val, address _owner, address _spender) view returns(uint256 _shares)
func (_Staking *StakingSession) Allowance(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Staking.Contract.Allowance(&_Staking.CallOpts, _val, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0x495aff71.
//
// Solidity: function allowance(string _val, address _owner, address _spender) view returns(uint256 _shares)
func (_Staking *StakingCallerSession) Allowance(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _Staking.Contract.Allowance(&_Staking.CallOpts, _val, _owner, _spender)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256 _shares, uint256 _delegateAmount)
func (_Staking *StakingCaller) Delegation(opts *bind.CallOpts, _val string, _del common.Address) (struct {
	Shares         *big.Int
	DelegateAmount *big.Int
}, error,
) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegation", _val, _del)

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
func (_Staking *StakingSession) Delegation(_val string, _del common.Address) (struct {
	Shares         *big.Int
	DelegateAmount *big.Int
}, error,
) {
	return _Staking.Contract.Delegation(&_Staking.CallOpts, _val, _del)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256 _shares, uint256 _delegateAmount)
func (_Staking *StakingCallerSession) Delegation(_val string, _del common.Address) (struct {
	Shares         *big.Int
	DelegateAmount *big.Int
}, error,
) {
	return _Staking.Contract.Delegation(&_Staking.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256 _reward)
func (_Staking *StakingCaller) DelegationRewards(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Staking.contract.Call(opts, &out, "delegationRewards", _val, _del)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256 _reward)
func (_Staking *StakingSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _Staking.Contract.DelegationRewards(&_Staking.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256 _reward)
func (_Staking *StakingCallerSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _Staking.Contract.DelegationRewards(&_Staking.CallOpts, _val, _del)
}

// Approve is a paid mutator transaction binding the contract method 0xbd5b93d2.
//
// Solidity: function approve(string _val, address _spender, uint256 _shares) returns(bool _result)
func (_Staking *StakingTransactor) Approve(opts *bind.TransactOpts, _val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "approve", _val, _spender, _shares)
}

// Approve is a paid mutator transaction binding the contract method 0xbd5b93d2.
//
// Solidity: function approve(string _val, address _spender, uint256 _shares) returns(bool _result)
func (_Staking *StakingSession) Approve(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Approve(&_Staking.TransactOpts, _val, _spender, _shares)
}

// Approve is a paid mutator transaction binding the contract method 0xbd5b93d2.
//
// Solidity: function approve(string _val, address _spender, uint256 _shares) returns(bool _result)
func (_Staking *StakingTransactorSession) Approve(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Approve(&_Staking.TransactOpts, _val, _spender, _shares)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256 _shares, uint256 _reward)
func (_Staking *StakingTransactor) Delegate(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegate", _val)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256 _shares, uint256 _reward)
func (_Staking *StakingSession) Delegate(_val string) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, _val)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256 _shares, uint256 _reward)
func (_Staking *StakingTransactorSession) Delegate(_val string) (*types.Transaction, error) {
	return _Staking.Contract.Delegate(&_Staking.TransactOpts, _val)
}

// Transfer is a paid mutator transaction binding the contract method 0xaa2fddc7.
//
// Solidity: function transfer(string _val, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_Staking *StakingTransactor) Transfer(opts *bind.TransactOpts, _val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "transfer", _val, _to, _shares)
}

// Transfer is a paid mutator transaction binding the contract method 0xaa2fddc7.
//
// Solidity: function transfer(string _val, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_Staking *StakingSession) Transfer(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Transfer(&_Staking.TransactOpts, _val, _to, _shares)
}

// Transfer is a paid mutator transaction binding the contract method 0xaa2fddc7.
//
// Solidity: function transfer(string _val, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_Staking *StakingTransactorSession) Transfer(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Transfer(&_Staking.TransactOpts, _val, _to, _shares)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x1fc38a8b.
//
// Solidity: function transferFrom(string _val, address _from, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_Staking *StakingTransactor) TransferFrom(opts *bind.TransactOpts, _val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "transferFrom", _val, _from, _to, _shares)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x1fc38a8b.
//
// Solidity: function transferFrom(string _val, address _from, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_Staking *StakingSession) TransferFrom(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.TransferFrom(&_Staking.TransactOpts, _val, _from, _to, _shares)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x1fc38a8b.
//
// Solidity: function transferFrom(string _val, address _from, address _to, uint256 _shares) returns(uint256 _token, uint256 _reward)
func (_Staking *StakingTransactorSession) TransferFrom(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.TransferFrom(&_Staking.TransactOpts, _val, _from, _to, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256 _amount, uint256 _reward, uint256 _completionTime)
func (_Staking *StakingTransactor) Undelegate(opts *bind.TransactOpts, _val string, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "undelegate", _val, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256 _amount, uint256 _reward, uint256 _completionTime)
func (_Staking *StakingSession) Undelegate(_val string, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, _val, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256 _amount, uint256 _reward, uint256 _completionTime)
func (_Staking *StakingTransactorSession) Undelegate(_val string, _shares *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.Undelegate(&_Staking.TransactOpts, _val, _shares)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256 _reward)
func (_Staking *StakingTransactor) Withdraw(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdraw", _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256 _reward)
func (_Staking *StakingSession) Withdraw(_val string) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256 _reward)
func (_Staking *StakingTransactorSession) Withdraw(_val string) (*types.Transaction, error) {
	return _Staking.Contract.Withdraw(&_Staking.TransactOpts, _val)
}

// StakingApproveIterator is returned from FilterApprove and is used to iterate over the raw logs and unpacked data for Approve events raised by the Staking contract.
type StakingApproveIterator struct {
	Event *StakingApprove // Event containing the contract specifics and raw log

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
func (it *StakingApproveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingApprove)
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
		it.Event = new(StakingApprove)
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
func (it *StakingApproveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingApproveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingApprove represents a Approve event raised by the Staking contract.
type StakingApprove struct {
	Owner     common.Address
	Spender   common.Address
	Validator string
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterApprove is a free log retrieval operation binding the contract event 0x757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020.
//
// Solidity: event Approve(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_Staking *StakingFilterer) FilterApprove(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StakingApproveIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Approve", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StakingApproveIterator{contract: _Staking.contract, event: "Approve", logs: logs, sub: sub}, nil
}

// WatchApprove is a free log subscription operation binding the contract event 0x757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020.
//
// Solidity: event Approve(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_Staking *StakingFilterer) WatchApprove(opts *bind.WatchOpts, sink chan<- *StakingApprove, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Approve", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingApprove)
				if err := _Staking.contract.UnpackLog(event, "Approve", log); err != nil {
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

// ParseApprove is a log parse operation binding the contract event 0x757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020.
//
// Solidity: event Approve(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_Staking *StakingFilterer) ParseApprove(log types.Log) (*StakingApprove, error) {
	event := new(StakingApprove)
	if err := _Staking.contract.UnpackLog(event, "Approve", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the Staking contract.
type StakingDelegateIterator struct {
	Event *StakingDelegate // Event containing the contract specifics and raw log

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
func (it *StakingDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegate)
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
		it.Event = new(StakingDelegate)
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
func (it *StakingDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegate represents a Delegate event raised by the Staking contract.
type StakingDelegate struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_Staking *StakingFilterer) FilterDelegate(opts *bind.FilterOpts, delegator []common.Address) (*StakingDelegateIterator, error) {
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Delegate", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegateIterator{contract: _Staking.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_Staking *StakingFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *StakingDelegate, delegator []common.Address) (event.Subscription, error) {
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Delegate", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegate)
				if err := _Staking.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// ParseDelegate is a log parse operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_Staking *StakingFilterer) ParseDelegate(log types.Log) (*StakingDelegate, error) {
	event := new(StakingDelegate)
	if err := _Staking.contract.UnpackLog(event, "Delegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Staking contract.
type StakingTransferIterator struct {
	Event *StakingTransfer // Event containing the contract specifics and raw log

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
func (it *StakingTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTransfer)
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
		it.Event = new(StakingTransfer)
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
func (it *StakingTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTransfer represents a Transfer event raised by the Staking contract.
type StakingTransfer struct {
	From      common.Address
	To        common.Address
	Validator string
	Shares    *big.Int
	Token     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0x39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813.
//
// Solidity: event Transfer(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_Staking *StakingFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingTransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingTransferIterator{contract: _Staking.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0x39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813.
//
// Solidity: event Transfer(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_Staking *StakingFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StakingTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTransfer)
				if err := _Staking.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0x39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813.
//
// Solidity: event Transfer(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_Staking *StakingFilterer) ParseTransfer(log types.Log) (*StakingTransfer, error) {
	event := new(StakingTransfer)
	if err := _Staking.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingUndelegateIterator is returned from FilterUndelegate and is used to iterate over the raw logs and unpacked data for Undelegate events raised by the Staking contract.
type StakingUndelegateIterator struct {
	Event *StakingUndelegate // Event containing the contract specifics and raw log

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
func (it *StakingUndelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingUndelegate)
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
		it.Event = new(StakingUndelegate)
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
func (it *StakingUndelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingUndelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingUndelegate represents a Undelegate event raised by the Staking contract.
type StakingUndelegate struct {
	Sender         common.Address
	Validator      string
	Shares         *big.Int
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUndelegate is a free log retrieval operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_Staking *StakingFilterer) FilterUndelegate(opts *bind.FilterOpts, sender []common.Address) (*StakingUndelegateIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Undelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingUndelegateIterator{contract: _Staking.contract, event: "Undelegate", logs: logs, sub: sub}, nil
}

// WatchUndelegate is a free log subscription operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_Staking *StakingFilterer) WatchUndelegate(opts *bind.WatchOpts, sink chan<- *StakingUndelegate, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Undelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingUndelegate)
				if err := _Staking.contract.UnpackLog(event, "Undelegate", log); err != nil {
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

// ParseUndelegate is a log parse operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_Staking *StakingFilterer) ParseUndelegate(log types.Log) (*StakingUndelegate, error) {
	event := new(StakingUndelegate)
	if err := _Staking.contract.UnpackLog(event, "Undelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Staking contract.
type StakingWithdrawIterator struct {
	Event *StakingWithdraw // Event containing the contract specifics and raw log

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
func (it *StakingWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingWithdraw)
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
		it.Event = new(StakingWithdraw)
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
func (it *StakingWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingWithdraw represents a Withdraw event raised by the Staking contract.
type StakingWithdraw struct {
	Sender    common.Address
	Validator string
	Reward    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_Staking *StakingFilterer) FilterWithdraw(opts *bind.FilterOpts, sender []common.Address) (*StakingWithdrawIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingWithdrawIterator{contract: _Staking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_Staking *StakingFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StakingWithdraw, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingWithdraw)
				if err := _Staking.contract.UnpackLog(event, "Withdraw", log); err != nil {
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
func (_Staking *StakingFilterer) ParseWithdraw(log types.Log) (*StakingWithdraw, error) {
	event := new(StakingWithdraw)
	if err := _Staking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
