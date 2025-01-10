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

// CrosschainTestMetaData contains all meta data concerning the CrosschainTest contract.
var CrosschainTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610cf4806100206000396000f3fe6080604052600436106100705760003560e01c806367cfd9d61161004e57806367cfd9d6146100dd5780638fefb765146100fd578063e6d69ede1461012b578063f73564751461014b57600080fd5b8063160d7c731461007557806316c75cfa1461009d5780634ac3bdc3146100bd575b600080fd5b61008861008336600461071e565b610179565b60405190151581526020015b60405180910390f35b3480156100a957600080fd5b506100886100b83660046107ae565b610469565b3480156100c957600080fd5b506100886100d83660046107fc565b6104db565b3480156100e957600080fd5b506100886100f83660046107ae565b610524565b34801561010957600080fd5b5061011d610118366004610841565b61054e565b604051908152602001610094565b34801561013757600080fd5b5061011d61014636600461095c565b6105c2565b34801561015757600080fd5b5061016161100481565b6040516001600160a01b039091168152602001610094565b60006001600160a01b038716156102ab576001600160a01b0387166323b872dd33306101a5888a610a63565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af11580156101f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021d9190610a89565b506001600160a01b03871663095ea7b361100461023a8789610a63565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015610285573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102a99190610a89565b505b6001600160a01b0387161561039257604051636eb1769f60e11b815230600482015261100460248201526000906001600160a01b0389169063dd62ed3e90604401602060405180830381865afa158015610309573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032d9190610aab565b90506103398587610a63565b811461038c5760405162461bcd60e51b815260206004820181905260248201527f616c6c6f77616e6365206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b506103ea565b61039c8486610a63565b34146103ea5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b206665656044820152606401610383565b60405163160d7c7360e01b81526110049063160d7c7390349061041b908b908b908b908b908b908b90600401610b11565b60206040518083038185885af1158015610439573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061045e9190610a89565b979650505050505050565b604051630b63ae7d60e11b8152600090611004906316c75cfa906104939086908690600401610b66565b602060405180830381865afa1580156104b0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d49190610a89565b9392505050565b604051634ac3bdc360e01b815260009061100490634ac3bdc3906105059086908690600401610b90565b6020604051808303816000875af11580156104b0573d6000803e3d6000fd5b6040516333e7eceb60e11b8152600090611004906367cfd9d6906104939086908690600401610b66565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa15801561059e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104d49190610aab565b60405163736b4f6f60e11b81526000906110049063e6d69ede906105fa908d908d908d908d908d908d908d908d908d90600401610bed565b6020604051808303816000875af1158015610619573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061063d9190610aab565b9a9950505050505050505050565b80356001600160a01b038116811461066257600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156106a6576106a6610667565b604052919050565b600082601f8301126106bf57600080fd5b813567ffffffffffffffff8111156106d9576106d9610667565b6106ec601f8201601f191660200161067d565b81815284602083860101111561070157600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c0878903121561073757600080fd5b6107408761064b565b9550602087013567ffffffffffffffff8082111561075d57600080fd5b6107698a838b016106ae565b965060408901359550606089013594506080890135935060a089013591508082111561079457600080fd5b506107a189828a016106ae565b9150509295509295509295565b600080604083850312156107c157600080fd5b823567ffffffffffffffff8111156107d857600080fd5b6107e4858286016106ae565b9250506107f36020840161064b565b90509250929050565b6000806040838503121561080f57600080fd5b823567ffffffffffffffff81111561082657600080fd5b610832858286016106ae565b95602094909401359450505050565b6000806040838503121561085457600080fd5b61085d8361064b565b946020939093013593505050565b600067ffffffffffffffff82111561088557610885610667565b5060051b60200190565b600082601f8301126108a057600080fd5b813560206108b56108b08361086b565b61067d565b82815260059290921b840181019181810190868411156108d457600080fd5b8286015b848110156108f6576108e98161064b565b83529183019183016108d8565b509695505050505050565b600082601f83011261091257600080fd5b813560206109226108b08361086b565b82815260059290921b8401810191818101908684111561094157600080fd5b8286015b848110156108f65780358352918301918301610945565b60008060008060008060008060006101208a8c03121561097b57600080fd5b893567ffffffffffffffff8082111561099357600080fd5b61099f8d838e016106ae565b9a506109ad60208d0161064b565b995060408c01359150808211156109c357600080fd5b6109cf8d838e0161088f565b985060608c01359150808211156109e557600080fd5b6109f18d838e01610901565b97506109ff60808d0161064b565b965060a08c0135915080821115610a1557600080fd5b610a218d838e016106ae565b955060c08c0135945060e08c013593506101008c0135915080821115610a4657600080fd5b50610a538c828d016106ae565b9150509295985092959850929598565b60008219821115610a8457634e487b7160e01b600052601160045260246000fd5b500190565b600060208284031215610a9b57600080fd5b815180151581146104d457600080fd5b600060208284031215610abd57600080fd5b5051919050565b6000815180845260005b81811015610aea57602081850181015186830182015201610ace565b81811115610afc576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038716815260c060208201819052600090610b3590830188610ac4565b86604084015285606084015284608084015282810360a0840152610b598185610ac4565b9998505050505050505050565b604081526000610b796040830185610ac4565b905060018060a01b03831660208301529392505050565b604081526000610ba36040830185610ac4565b90508260208301529392505050565b600081518084526020808501945080840160005b83811015610be257815187529582019590820190600101610bc6565b509495945050505050565b6000610120808352610c018184018d610ac4565b6001600160a01b038c811660208681019190915285830360408701528c518084528d82019450909283019060005b81811015610c4d578551841683529484019491840191600101610c2f565b50508581036060870152610c61818d610bb2565b9350505050610c7b60808401896001600160a01b03169052565b82810360a0840152610c8d8188610ac4565b90508560c08401528460e0840152828103610100840152610cae8185610ac4565b9c9b50505050505050505050505056fea2646970667358221220eaab9d7a3b40a36b8e21e58fdfdd262315f7f3cf34f15838aeb1ef08486af61c64736f6c634300080a0033",
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

// CROSSCHAINADDRESS is a free data retrieval call binding the contract method 0xf7356475.
//
// Solidity: function CROSS_CHAIN_ADDRESS() view returns(address)
func (_CrosschainTest *CrosschainTestCaller) CROSSCHAINADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "CROSS_CHAIN_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CROSSCHAINADDRESS is a free data retrieval call binding the contract method 0xf7356475.
//
// Solidity: function CROSS_CHAIN_ADDRESS() view returns(address)
func (_CrosschainTest *CrosschainTestSession) CROSSCHAINADDRESS() (common.Address, error) {
	return _CrosschainTest.Contract.CROSSCHAINADDRESS(&_CrosschainTest.CallOpts)
}

// CROSSCHAINADDRESS is a free data retrieval call binding the contract method 0xf7356475.
//
// Solidity: function CROSS_CHAIN_ADDRESS() view returns(address)
func (_CrosschainTest *CrosschainTestCallerSession) CROSSCHAINADDRESS() (common.Address, error) {
	return _CrosschainTest.Contract.CROSSCHAINADDRESS(&_CrosschainTest.CallOpts)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrosschainTest *CrosschainTestCaller) BridgeCoinAmount(opts *bind.CallOpts, _token common.Address, _target [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "bridgeCoinAmount", _token, _target)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrosschainTest *CrosschainTestSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _CrosschainTest.Contract.BridgeCoinAmount(&_CrosschainTest.CallOpts, _token, _target)
}

// BridgeCoinAmount is a free data retrieval call binding the contract method 0x8fefb765.
//
// Solidity: function bridgeCoinAmount(address _token, bytes32 _target) view returns(uint256)
func (_CrosschainTest *CrosschainTestCallerSession) BridgeCoinAmount(_token common.Address, _target [32]byte) (*big.Int, error) {
	return _CrosschainTest.Contract.BridgeCoinAmount(&_CrosschainTest.CallOpts, _token, _target)
}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCaller) HasOracle(opts *bind.CallOpts, _chain string, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "hasOracle", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestSession) HasOracle(_chain string, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.HasOracle(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// HasOracle is a free data retrieval call binding the contract method 0x67cfd9d6.
//
// Solidity: function hasOracle(string _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCallerSession) HasOracle(_chain string, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.HasOracle(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCaller) IsOracleOnline(opts *bind.CallOpts, _chain string, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "isOracleOnline", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestSession) IsOracleOnline(_chain string, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.IsOracleOnline(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0x16c75cfa.
//
// Solidity: function isOracleOnline(string _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCallerSession) IsOracleOnline(_chain string, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.IsOracleOnline(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _receiver, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) returns(uint256)
func (_CrosschainTest *CrosschainTestTransactor) BridgeCall(opts *bind.TransactOpts, _dstChain string, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "bridgeCall", _dstChain, _receiver, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _receiver, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) returns(uint256)
func (_CrosschainTest *CrosschainTestSession) BridgeCall(_dstChain string, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _CrosschainTest.Contract.BridgeCall(&_CrosschainTest.TransactOpts, _dstChain, _receiver, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _receiver, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) returns(uint256)
func (_CrosschainTest *CrosschainTestTransactorSession) BridgeCall(_dstChain string, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _CrosschainTest.Contract.BridgeCall(&_CrosschainTest.TransactOpts, _dstChain, _receiver, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
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

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_CrosschainTest *CrosschainTestTransactor) ExecuteClaim(opts *bind.TransactOpts, _chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "executeClaim", _chain, _eventNonce)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_CrosschainTest *CrosschainTestSession) ExecuteClaim(_chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.ExecuteClaim(&_CrosschainTest.TransactOpts, _chain, _eventNonce)
}

// ExecuteClaim is a paid mutator transaction binding the contract method 0x4ac3bdc3.
//
// Solidity: function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)
func (_CrosschainTest *CrosschainTestTransactorSession) ExecuteClaim(_chain string, _eventNonce *big.Int) (*types.Transaction, error) {
	return _CrosschainTest.Contract.ExecuteClaim(&_CrosschainTest.TransactOpts, _chain, _eventNonce)
}
