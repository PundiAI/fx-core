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
	ABI: "[{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_denom\",\"type\":\"bytes32\"}],\"name\":\"getERC20Token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_enable\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610f29806100206000396000f3fe60806040526004361061007b5760003560e01c8063d5147e6d1161004e578063d5147e6d14610116578063e6d69ede14610136578063e8e1f20514610149578063f73564751461018857600080fd5b8063160d7c73146100805780634ac3bdc3146100a85780638fefb765146100c8578063a5df3875146100f6575b600080fd5b61009361008e366004610905565b6101b6565b60405190151581526020015b60405180910390f35b3480156100b457600080fd5b506100936100c3366004610997565b6103c8565b3480156100d457600080fd5b506100e86100e33660046109dc565b61043c565b60405190815260200161009f565b34801561010257600080fd5b50610093610111366004610a08565b6104b0565b34801561012257600080fd5b50610093610131366004610a08565b610501565b6100e8610144366004610b2b565b610539565b34801561015557600080fd5b50610169610164366004610c32565b6107b7565b604080516001600160a01b03909316835290151560208301520161009f565b34801561019457600080fd5b5061019e61100481565b6040516001600160a01b03909116815260200161009f565b60006001600160a01b038716156102ec576001600160a01b0387166323b872dd33306101e2888a610c61565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af1158015610236573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025a9190610c89565b506001600160a01b03871663095ea7b36110046102778789610c61565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af11580156102c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102e69190610c89565b50610349565b6102f68486610c61565b34146103495760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b60405163160d7c7360e01b81526110049063160d7c7390349061037a908b908b908b908b908b908b90600401610cf1565b60206040518083038185885af1158015610398573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906103bd9190610c89565b979650505050505050565b604051634ac3bdc360e01b815260009061100490634ac3bdc3906103f29086908690600401610d46565b6020604051808303816000875af1158015610411573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104359190610c89565b9392505050565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa15801561048c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104359190610d68565b60405163a5df387560e01b8152600481018390526001600160a01b03821660248201526000906110049063a5df3875906044015b602060405180830381865afa158015610411573d6000803e3d6000fd5b60405163d5147e6d60e01b8152600481018390526001600160a01b03821660248201526000906110049063d5147e6d906044016104e4565b600086518851146105965760405162461bcd60e51b815260206004820152602160248201527f746f6b656e20616e6420616d6f756e74206c656e677468206e6f7420657175616044820152601b60fa1b6064820152608401610340565b60005b885181101561072e578881815181106105b4576105b4610d81565b60200260200101516001600160a01b03166323b872dd33308b85815181106105de576105de610d81565b60209081029190910101516040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af115801561063d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106619190610c89565b5088818151811061067457610674610d81565b60200260200101516001600160a01b031663095ea7b36110048a848151811061069f5761069f610d81565b60200260200101516040518363ffffffff1660e01b81526004016106d89291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af11580156106f7573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061071b9190610c89565b508061072681610d97565b915050610599565b5060405163736b4f6f60e11b81526110049063e6d69ede903490610766908e908e908e908e908e908e908e908e908e90600401610ded565b60206040518083038185885af1158015610784573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906107a99190610d68565b9a9950505050505050505050565b60405163e8e1f20560e01b81526004810182905260009081906110049063e8e1f205906024016040805180830381865afa1580156107f9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061081d9190610ebe565b91509150915091565b6001600160a01b038116811461083b57600080fd5b50565b803561084981610826565b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561088d5761088d61084e565b604052919050565b600082601f8301126108a657600080fd5b813567ffffffffffffffff8111156108c0576108c061084e565b6108d3601f8201601f1916602001610864565b8181528460208386010111156108e857600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c0878903121561091e57600080fd5b863561092981610826565b9550602087013567ffffffffffffffff8082111561094657600080fd5b6109528a838b01610895565b965060408901359550606089013594506080890135935060a089013591508082111561097d57600080fd5b5061098a89828a01610895565b9150509295509295509295565b600080604083850312156109aa57600080fd5b823567ffffffffffffffff8111156109c157600080fd5b6109cd85828601610895565b95602094909401359450505050565b600080604083850312156109ef57600080fd5b82356109fa81610826565b946020939093013593505050565b60008060408385031215610a1b57600080fd5b823591506020830135610a2d81610826565b809150509250929050565b600067ffffffffffffffff821115610a5257610a5261084e565b5060051b60200190565b600082601f830112610a6d57600080fd5b81356020610a82610a7d83610a38565b610864565b82815260059290921b84018101918181019086841115610aa157600080fd5b8286015b84811015610ac5578035610ab881610826565b8352918301918301610aa5565b509695505050505050565b600082601f830112610ae157600080fd5b81356020610af1610a7d83610a38565b82815260059290921b84018101918181019086841115610b1057600080fd5b8286015b84811015610ac55780358352918301918301610b14565b60008060008060008060008060006101208a8c031215610b4a57600080fd5b893567ffffffffffffffff80821115610b6257600080fd5b610b6e8d838e01610895565b9a50610b7c60208d0161083e565b995060408c0135915080821115610b9257600080fd5b610b9e8d838e01610a5c565b985060608c0135915080821115610bb457600080fd5b610bc08d838e01610ad0565b9750610bce60808d0161083e565b965060a08c0135915080821115610be457600080fd5b610bf08d838e01610895565b955060c08c0135945060e08c013593506101008c0135915080821115610c1557600080fd5b50610c228c828d01610895565b9150509295985092959850929598565b600060208284031215610c4457600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610c7457610c74610c4b565b500190565b8051801515811461084957600080fd5b600060208284031215610c9b57600080fd5b61043582610c79565b6000815180845260005b81811015610cca57602081850181015186830182015201610cae565b81811115610cdc576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038716815260c060208201819052600090610d1590830188610ca4565b86604084015285606084015284608084015282810360a0840152610d398185610ca4565b9998505050505050505050565b604081526000610d596040830185610ca4565b90508260208301529392505050565b600060208284031215610d7a57600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415610dab57610dab610c4b565b5060010190565b600081518084526020808501945080840160005b83811015610de257815187529582019590820190600101610dc6565b509495945050505050565b6000610120808352610e018184018d610ca4565b6001600160a01b038c811660208681019190915285830360408701528c518084528d82019450909283019060005b81811015610e4d578551841683529484019491840191600101610e2f565b50508581036060870152610e61818d610db2565b9350505050610e7b60808401896001600160a01b03169052565b82810360a0840152610e8d8188610ca4565b90508560c08401528460e0840152828103610100840152610eae8185610ca4565b9c9b505050505050505050505050565b60008060408385031215610ed157600080fd5b8251610edc81610826565b9150610eea60208401610c79565b9050925092905056fea26469706673582212207e483c982a79c7c22dd4c7b8e427b8fadc975c9d8c6d3115d9ff8281a9ba6ded64736f6c634300080a0033",
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

// GetERC20Token is a free data retrieval call binding the contract method 0xe8e1f205.
//
// Solidity: function getERC20Token(bytes32 _denom) view returns(address _token, bool _enable)
func (_CrosschainTest *CrosschainTestCaller) GetERC20Token(opts *bind.CallOpts, _denom [32]byte) (struct {
	Token  common.Address
	Enable bool
}, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "getERC20Token", _denom)

	outstruct := new(struct {
		Token  common.Address
		Enable bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Token = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Enable = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// GetERC20Token is a free data retrieval call binding the contract method 0xe8e1f205.
//
// Solidity: function getERC20Token(bytes32 _denom) view returns(address _token, bool _enable)
func (_CrosschainTest *CrosschainTestSession) GetERC20Token(_denom [32]byte) (struct {
	Token  common.Address
	Enable bool
}, error) {
	return _CrosschainTest.Contract.GetERC20Token(&_CrosschainTest.CallOpts, _denom)
}

// GetERC20Token is a free data retrieval call binding the contract method 0xe8e1f205.
//
// Solidity: function getERC20Token(bytes32 _denom) view returns(address _token, bool _enable)
func (_CrosschainTest *CrosschainTestCallerSession) GetERC20Token(_denom [32]byte) (struct {
	Token  common.Address
	Enable bool
}, error) {
	return _CrosschainTest.Contract.GetERC20Token(&_CrosschainTest.CallOpts, _denom)
}

// HasOracle is a free data retrieval call binding the contract method 0xa5df3875.
//
// Solidity: function hasOracle(bytes32 _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCaller) HasOracle(opts *bind.CallOpts, _chain [32]byte, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "hasOracle", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasOracle is a free data retrieval call binding the contract method 0xa5df3875.
//
// Solidity: function hasOracle(bytes32 _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestSession) HasOracle(_chain [32]byte, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.HasOracle(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// HasOracle is a free data retrieval call binding the contract method 0xa5df3875.
//
// Solidity: function hasOracle(bytes32 _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCallerSession) HasOracle(_chain [32]byte, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.HasOracle(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0xd5147e6d.
//
// Solidity: function isOracleOnline(bytes32 _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCaller) IsOracleOnline(opts *bind.CallOpts, _chain [32]byte, _externalAddress common.Address) (bool, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "isOracleOnline", _chain, _externalAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracleOnline is a free data retrieval call binding the contract method 0xd5147e6d.
//
// Solidity: function isOracleOnline(bytes32 _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestSession) IsOracleOnline(_chain [32]byte, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.IsOracleOnline(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// IsOracleOnline is a free data retrieval call binding the contract method 0xd5147e6d.
//
// Solidity: function isOracleOnline(bytes32 _chain, address _externalAddress) view returns(bool _result)
func (_CrosschainTest *CrosschainTestCallerSession) IsOracleOnline(_chain [32]byte, _externalAddress common.Address) (bool, error) {
	return _CrosschainTest.Contract.IsOracleOnline(&_CrosschainTest.CallOpts, _chain, _externalAddress)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _receiver, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) payable returns(uint256)
func (_CrosschainTest *CrosschainTestTransactor) BridgeCall(opts *bind.TransactOpts, _dstChain string, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _CrosschainTest.contract.Transact(opts, "bridgeCall", _dstChain, _receiver, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _receiver, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) payable returns(uint256)
func (_CrosschainTest *CrosschainTestSession) BridgeCall(_dstChain string, _receiver common.Address, _tokens []common.Address, _amounts []*big.Int, _to common.Address, _data []byte, _quoteId *big.Int, _gasLimit *big.Int, _memo []byte) (*types.Transaction, error) {
	return _CrosschainTest.Contract.BridgeCall(&_CrosschainTest.TransactOpts, _dstChain, _receiver, _tokens, _amounts, _to, _data, _quoteId, _gasLimit, _memo)
}

// BridgeCall is a paid mutator transaction binding the contract method 0xe6d69ede.
//
// Solidity: function bridgeCall(string _dstChain, address _receiver, address[] _tokens, uint256[] _amounts, address _to, bytes _data, uint256 _quoteId, uint256 _gasLimit, bytes _memo) payable returns(uint256)
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
