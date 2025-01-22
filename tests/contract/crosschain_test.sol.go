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
	ABI: "[{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610e03806100206000396000f3fe6080604052600436106100705760003560e01c8063a5df38751161004e578063a5df3875146100eb578063d5147e6d1461010b578063e6d69ede1461012b578063f73564751461013e57600080fd5b8063160d7c73146100755780634ac3bdc31461009d5780638fefb765146100bd575b600080fd5b610088610083366004610840565b61016c565b60405190151581526020015b60405180910390f35b3480156100a957600080fd5b506100886100b83660046108d0565b61037e565b3480156100c957600080fd5b506100dd6100d8366004610915565b6103f2565b604051908152602001610094565b3480156100f757600080fd5b5061008861010636600461093f565b610466565b34801561011757600080fd5b5061008861012636600461093f565b6104b7565b6100dd610139366004610a5c565b6104ef565b34801561014a57600080fd5b5061015461100481565b6040516001600160a01b039091168152602001610094565b60006001600160a01b038716156102a2576001600160a01b0387166323b872dd3330610198888a610b79565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af11580156101ec573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102109190610b91565b506001600160a01b03871663095ea7b361100461022d8789610b79565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af1158015610278573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061029c9190610b91565b506102ff565b6102ac8486610b79565b34146102ff5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b60405163160d7c7360e01b81526110049063160d7c73903490610330908b908b908b908b908b908b90600401610c00565b60206040518083038185885af115801561034e573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906103739190610b91565b979650505050505050565b604051634ac3bdc360e01b815260009061100490634ac3bdc3906103a89086908690600401610c55565b6020604051808303816000875af11580156103c7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103eb9190610b91565b9392505050565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa158015610442573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103eb9190610c77565b60405163a5df387560e01b8152600481018390526001600160a01b03821660248201526000906110049063a5df3875906044015b602060405180830381865afa1580156103c7573d6000803e3d6000fd5b60405163d5147e6d60e01b8152600481018390526001600160a01b03821660248201526000906110049063d5147e6d9060440161049a565b6000865188511461054c5760405162461bcd60e51b815260206004820152602160248201527f746f6b656e20616e6420616d6f756e74206c656e677468206e6f7420657175616044820152601b60fa1b60648201526084016102f6565b60005b88518110156106e45788818151811061056a5761056a610c90565b60200260200101516001600160a01b03166323b872dd33308b858151811061059457610594610c90565b60209081029190910101516040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af11580156105f3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106179190610b91565b5088818151811061062a5761062a610c90565b60200260200101516001600160a01b031663095ea7b36110048a848151811061065557610655610c90565b60200260200101516040518363ffffffff1660e01b815260040161068e9291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af11580156106ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106d19190610b91565b50806106dc81610ca6565b91505061054f565b5060405163736b4f6f60e11b81526110049063e6d69ede90349061071c908e908e908e908e908e908e908e908e908e90600401610cfc565b60206040518083038185885af115801561073a573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019061075f9190610c77565b9a9950505050505050505050565b80356001600160a01b038116811461078457600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156107c8576107c8610789565b604052919050565b600082601f8301126107e157600080fd5b813567ffffffffffffffff8111156107fb576107fb610789565b61080e601f8201601f191660200161079f565b81815284602083860101111561082357600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c0878903121561085957600080fd5b6108628761076d565b9550602087013567ffffffffffffffff8082111561087f57600080fd5b61088b8a838b016107d0565b965060408901359550606089013594506080890135935060a08901359150808211156108b657600080fd5b506108c389828a016107d0565b9150509295509295509295565b600080604083850312156108e357600080fd5b823567ffffffffffffffff8111156108fa57600080fd5b610906858286016107d0565b95602094909401359450505050565b6000806040838503121561092857600080fd5b6109318361076d565b946020939093013593505050565b6000806040838503121561095257600080fd5b823591506109626020840161076d565b90509250929050565b600067ffffffffffffffff82111561098557610985610789565b5060051b60200190565b600082601f8301126109a057600080fd5b813560206109b56109b08361096b565b61079f565b82815260059290921b840181019181810190868411156109d457600080fd5b8286015b848110156109f6576109e98161076d565b83529183019183016109d8565b509695505050505050565b600082601f830112610a1257600080fd5b81356020610a226109b08361096b565b82815260059290921b84018101918181019086841115610a4157600080fd5b8286015b848110156109f65780358352918301918301610a45565b60008060008060008060008060006101208a8c031215610a7b57600080fd5b893567ffffffffffffffff80821115610a9357600080fd5b610a9f8d838e016107d0565b9a50610aad60208d0161076d565b995060408c0135915080821115610ac357600080fd5b610acf8d838e0161098f565b985060608c0135915080821115610ae557600080fd5b610af18d838e01610a01565b9750610aff60808d0161076d565b965060a08c0135915080821115610b1557600080fd5b610b218d838e016107d0565b955060c08c0135945060e08c013593506101008c0135915080821115610b4657600080fd5b50610b538c828d016107d0565b9150509295985092959850929598565b634e487b7160e01b600052601160045260246000fd5b60008219821115610b8c57610b8c610b63565b500190565b600060208284031215610ba357600080fd5b815180151581146103eb57600080fd5b6000815180845260005b81811015610bd957602081850181015186830182015201610bbd565b81811115610beb576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038716815260c060208201819052600090610c2490830188610bb3565b86604084015285606084015284608084015282810360a0840152610c488185610bb3565b9998505050505050505050565b604081526000610c686040830185610bb3565b90508260208301529392505050565b600060208284031215610c8957600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415610cba57610cba610b63565b5060010190565b600081518084526020808501945080840160005b83811015610cf157815187529582019590820190600101610cd5565b509495945050505050565b6000610120808352610d108184018d610bb3565b6001600160a01b038c811660208681019190915285830360408701528c518084528d82019450909283019060005b81811015610d5c578551841683529484019491840191600101610d3e565b50508581036060870152610d70818d610cc1565b9350505050610d8a60808401896001600160a01b03169052565b82810360a0840152610d9c8188610bb3565b90508560c08401528460e0840152828103610100840152610dbd8185610bb3565b9c9b50505050505050505050505056fea264697066735822122084d946cdee812bca08d11d1a8d86f828898736d7201775c604ea2a18a55bdb6964736f6c634300080a0033",
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
