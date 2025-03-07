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
	ABI: "[{\"inputs\":[],\"name\":\"BRIDGE_FEE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_denom\",\"type\":\"bytes32\"}],\"name\":\"getERC20Token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_enable\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onBridgeCall\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onRevert\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506113d6806100206000396000f3fe60806040526004361061009c5760003560e01c8063a5df387511610064578063a5df387514610159578063c863cee514610179578063d5147e6d146101a7578063e6d69ede146101c7578063e8e1f205146101da578063f73564751461021957600080fd5b8063160d7c73146100a157806332e1e16e146100c95780634ac3bdc3146100eb57806357ffc0921461010b5780638fefb7651461012b575b600080fd5b6100b46100af366004610c01565b61022f565b60405190151581526020015b60405180910390f35b3480156100d557600080fd5b506100e96100e4366004610c93565b610441565b005b3480156100f757600080fd5b506100b4610106366004610cda565b610491565b34801561011757600080fd5b506100e9610126366004610e12565b610505565b34801561013757600080fd5b5061014b610146366004610ec5565b610559565b6040519081526020016100c0565b34801561016557600080fd5b506100b4610174366004610ef1565b6105cd565b34801561018557600080fd5b5061018f61100581565b6040516001600160a01b0390911681526020016100c0565b3480156101b357600080fd5b506100b46101c2366004610ef1565b61061e565b61014b6101d5366004610f21565b610656565b3480156101e657600080fd5b506101fa6101f5366004611028565b610ab3565b604080516001600160a01b0390931683529015156020830152016100c0565b34801561022557600080fd5b5061018f61100481565b60006001600160a01b03871615610365576001600160a01b0387166323b872dd333061025b888a611057565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af11580156102af573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102d3919061107f565b506001600160a01b03871663095ea7b36110046102f08789611057565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af115801561033b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061035f919061107f565b506103c2565b61036f8486611057565b34146103c25760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b60405163160d7c7360e01b81526110049063160d7c739034906103f3908b908b908b908b908b908b906004016110e7565b60206040518083038185885af1158015610411573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190610436919061107f565b979650505050505050565b336110041461048d5760405162461bcd60e51b81526020600482015260186024820152776f6e6c792063726f73732d636861696e206164647265737360401b60448201526064016103b9565b5050565b604051634ac3bdc360e01b815260009061100490634ac3bdc3906104bb908690869060040161113c565b6020604051808303816000875af11580156104da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104fe919061107f565b9392505050565b33611004146105515760405162461bcd60e51b81526020600482015260186024820152776f6e6c792063726f73732d636861696e206164647265737360401b60448201526064016103b9565b505050505050565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa1580156105a9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104fe919061115e565b60405163a5df387560e01b8152600481018390526001600160a01b03821660248201526000906110049063a5df3875906044015b602060405180830381865afa1580156104da573d6000803e3d6000fd5b60405163d5147e6d60e01b8152600481018390526001600160a01b03821660248201526000906110049063d5147e6d90604401610601565b600086518851146106b35760405162461bcd60e51b815260206004820152602160248201527f746f6b656e20616e6420616d6f756e74206c656e677468206e6f7420657175616044820152601b60fa1b60648201526084016103b9565b60005b885181101561084b578881815181106106d1576106d1611177565b60200260200101516001600160a01b03166323b872dd33308b85815181106106fb576106fb611177565b60209081029190910101516040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af115801561075a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061077e919061107f565b5088818151811061079157610791611177565b60200260200101516001600160a01b031663095ea7b36110048a84815181106107bc576107bc611177565b60200260200101516040518363ffffffff1660e01b81526004016107f59291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af1158015610814573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610838919061107f565b50806108438161118d565b9150506106b6565b5060405163a8541c1760e01b8152600481018590526000906110059063a8541c179060240160e060405180830381865afa15801561088d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108b191906111c0565b60408101519091503490676170756e6469616960c01b141561092e5781608001513410156109185760405162461bcd60e51b81526020600482015260146024820152730dae6ce5cecc2d8eaca40dcdee840cadcdeeaced60631b60448201526064016103b9565b60808201516109279082611258565b9050610a21565b604080830151905163e8e1f20560e01b81526000916110049163e8e1f2059161095d9160040190815260200190565b6040805180830381865afa158015610979573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061099d919061126f565b5060808401516040516323b872dd60e01b815233600482015230602482015260448101919091529091506001600160a01b038216906323b872dd906064016020604051808303816000875af11580156109fa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a1e919061107f565b50505b6110046001600160a01b031663e6d69ede828e8e8e8e8e8e8e8e8e6040518b63ffffffff1660e01b8152600401610a60999897969594939291906112df565b60206040518083038185885af1158015610a7e573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190610aa3919061115e565b9c9b505050505050505050505050565b60405163e8e1f20560e01b81526004810182905260009081906110049063e8e1f205906024016040805180830381865afa158015610af5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b19919061126f565b91509150915091565b6001600160a01b0381168114610b3757600080fd5b50565b8035610b4581610b22565b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610b8957610b89610b4a565b604052919050565b600082601f830112610ba257600080fd5b813567ffffffffffffffff811115610bbc57610bbc610b4a565b610bcf601f8201601f1916602001610b60565b818152846020838601011115610be457600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215610c1a57600080fd5b8635610c2581610b22565b9550602087013567ffffffffffffffff80821115610c4257600080fd5b610c4e8a838b01610b91565b965060408901359550606089013594506080890135935060a0890135915080821115610c7957600080fd5b50610c8689828a01610b91565b9150509295509295509295565b60008060408385031215610ca657600080fd5b82359150602083013567ffffffffffffffff811115610cc457600080fd5b610cd085828601610b91565b9150509250929050565b60008060408385031215610ced57600080fd5b823567ffffffffffffffff811115610d0457600080fd5b610d1085828601610b91565b95602094909401359450505050565b600067ffffffffffffffff821115610d3957610d39610b4a565b5060051b60200190565b600082601f830112610d5457600080fd5b81356020610d69610d6483610d1f565b610b60565b82815260059290921b84018101918181019086841115610d8857600080fd5b8286015b84811015610dac578035610d9f81610b22565b8352918301918301610d8c565b509695505050505050565b600082601f830112610dc857600080fd5b81356020610dd8610d6483610d1f565b82815260059290921b84018101918181019086841115610df757600080fd5b8286015b84811015610dac5780358352918301918301610dfb565b60008060008060008060c08789031215610e2b57600080fd5b610e3487610b3a565b9550610e4260208801610b3a565b9450604087013567ffffffffffffffff80821115610e5f57600080fd5b610e6b8a838b01610d43565b95506060890135915080821115610e8157600080fd5b610e8d8a838b01610db7565b94506080890135915080821115610ea357600080fd5b610eaf8a838b01610b91565b935060a0890135915080821115610c7957600080fd5b60008060408385031215610ed857600080fd5b8235610ee381610b22565b946020939093013593505050565b60008060408385031215610f0457600080fd5b823591506020830135610f1681610b22565b809150509250929050565b60008060008060008060008060006101208a8c031215610f4057600080fd5b893567ffffffffffffffff80821115610f5857600080fd5b610f648d838e01610b91565b9a50610f7260208d01610b3a565b995060408c0135915080821115610f8857600080fd5b610f948d838e01610d43565b985060608c0135915080821115610faa57600080fd5b610fb68d838e01610db7565b9750610fc460808d01610b3a565b965060a08c0135915080821115610fda57600080fd5b610fe68d838e01610b91565b955060c08c0135945060e08c013593506101008c013591508082111561100b57600080fd5b506110188c828d01610b91565b9150509295985092959850929598565b60006020828403121561103a57600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b6000821982111561106a5761106a611041565b500190565b80518015158114610b4557600080fd5b60006020828403121561109157600080fd5b6104fe8261106f565b6000815180845260005b818110156110c0576020818501810151868301820152016110a4565b818111156110d2576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038716815260c06020820181905260009061110b9083018861109a565b86604084015285606084015284608084015282810360a084015261112f818561109a565b9998505050505050505050565b60408152600061114f604083018561109a565b90508260208301529392505050565b60006020828403121561117057600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b60006000198214156111a1576111a1611041565b5060010190565b805167ffffffffffffffff81168114610b4557600080fd5b600060e082840312156111d257600080fd5b60405160e0810181811067ffffffffffffffff821117156111f5576111f5610b4a565b8060405250825181526020830151602082015260408301516040820152606083015161122081610b22565b60608201526080838101519082015261123b60a084016111a8565b60a082015261124c60c084016111a8565b60c08201529392505050565b60008282101561126a5761126a611041565b500390565b6000806040838503121561128257600080fd5b825161128d81610b22565b915061129b6020840161106f565b90509250929050565b600081518084526020808501945080840160005b838110156112d4578151875295820195908201906001016112b8565b509495945050505050565b60006101208083526112f38184018d61109a565b6001600160a01b038c811660208681019190915285830360408701528c518084528d82019450909283019060005b8181101561133f578551841683529484019491840191600101611321565b50508581036060870152611353818d6112a4565b935050505061136d60808401896001600160a01b03169052565b82810360a084015261137f818861109a565b90508560c08401528460e0840152828103610100840152610aa3818561109a56fea2646970667358221220ea1b42836b572269d6e204a472a81d9c2bd77f8ace9dc4a80a0ffe954d350c6964736f6c634300080a0033",
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

// BRIDGEFEEADDRESS is a free data retrieval call binding the contract method 0xc863cee5.
//
// Solidity: function BRIDGE_FEE_ADDRESS() view returns(address)
func (_CrosschainTest *CrosschainTestCaller) BRIDGEFEEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "BRIDGE_FEE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BRIDGEFEEADDRESS is a free data retrieval call binding the contract method 0xc863cee5.
//
// Solidity: function BRIDGE_FEE_ADDRESS() view returns(address)
func (_CrosschainTest *CrosschainTestSession) BRIDGEFEEADDRESS() (common.Address, error) {
	return _CrosschainTest.Contract.BRIDGEFEEADDRESS(&_CrosschainTest.CallOpts)
}

// BRIDGEFEEADDRESS is a free data retrieval call binding the contract method 0xc863cee5.
//
// Solidity: function BRIDGE_FEE_ADDRESS() view returns(address)
func (_CrosschainTest *CrosschainTestCallerSession) BRIDGEFEEADDRESS() (common.Address, error) {
	return _CrosschainTest.Contract.BRIDGEFEEADDRESS(&_CrosschainTest.CallOpts)
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

// OnBridgeCall is a free data retrieval call binding the contract method 0x57ffc092.
//
// Solidity: function onBridgeCall(address , address , address[] , uint256[] , bytes , bytes ) view returns()
func (_CrosschainTest *CrosschainTestCaller) OnBridgeCall(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 []common.Address, arg3 []*big.Int, arg4 []byte, arg5 []byte) error {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "onBridgeCall", arg0, arg1, arg2, arg3, arg4, arg5)

	if err != nil {
		return err
	}

	return err

}

// OnBridgeCall is a free data retrieval call binding the contract method 0x57ffc092.
//
// Solidity: function onBridgeCall(address , address , address[] , uint256[] , bytes , bytes ) view returns()
func (_CrosschainTest *CrosschainTestSession) OnBridgeCall(arg0 common.Address, arg1 common.Address, arg2 []common.Address, arg3 []*big.Int, arg4 []byte, arg5 []byte) error {
	return _CrosschainTest.Contract.OnBridgeCall(&_CrosschainTest.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5)
}

// OnBridgeCall is a free data retrieval call binding the contract method 0x57ffc092.
//
// Solidity: function onBridgeCall(address , address , address[] , uint256[] , bytes , bytes ) view returns()
func (_CrosschainTest *CrosschainTestCallerSession) OnBridgeCall(arg0 common.Address, arg1 common.Address, arg2 []common.Address, arg3 []*big.Int, arg4 []byte, arg5 []byte) error {
	return _CrosschainTest.Contract.OnBridgeCall(&_CrosschainTest.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5)
}

// OnRevert is a free data retrieval call binding the contract method 0x32e1e16e.
//
// Solidity: function onRevert(uint256 , bytes ) view returns()
func (_CrosschainTest *CrosschainTestCaller) OnRevert(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _CrosschainTest.contract.Call(opts, &out, "onRevert", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnRevert is a free data retrieval call binding the contract method 0x32e1e16e.
//
// Solidity: function onRevert(uint256 , bytes ) view returns()
func (_CrosschainTest *CrosschainTestSession) OnRevert(arg0 *big.Int, arg1 []byte) error {
	return _CrosschainTest.Contract.OnRevert(&_CrosschainTest.CallOpts, arg0, arg1)
}

// OnRevert is a free data retrieval call binding the contract method 0x32e1e16e.
//
// Solidity: function onRevert(uint256 , bytes ) view returns()
func (_CrosschainTest *CrosschainTestCallerSession) OnRevert(arg0 *big.Int, arg1 []byte) error {
	return _CrosschainTest.Contract.OnRevert(&_CrosschainTest.CallOpts, arg0, arg1)
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
