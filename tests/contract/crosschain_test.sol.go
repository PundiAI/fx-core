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
	ABI: "[{\"inputs\":[],\"name\":\"BRIDGE_FEE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CROSS_CHAIN_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_quoteId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_memo\",\"type\":\"bytes\"}],\"name\":\"bridgeCall\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"}],\"name\":\"bridgeCoinAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_receipt\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_target\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"crossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chain\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_eventNonce\",\"type\":\"uint256\"}],\"name\":\"executeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_denom\",\"type\":\"bytes32\"}],\"name\":\"getERC20Token\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_enable\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"hasOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chain\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_externalAddress\",\"type\":\"address\"}],\"name\":\"isOracleOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506111e0806100206000396000f3fe6080604052600436106100865760003560e01c8063c863cee511610059578063c863cee514610121578063d5147e6d1461014f578063e6d69ede1461016f578063e8e1f20514610182578063f7356475146101c157600080fd5b8063160d7c731461008b5780634ac3bdc3146100b35780638fefb765146100d3578063a5df387514610101575b600080fd5b61009e610099366004610b05565b6101d7565b60405190151581526020015b60405180910390f35b3480156100bf57600080fd5b5061009e6100ce366004610b97565b6103e9565b3480156100df57600080fd5b506100f36100ee366004610bdc565b61045d565b6040519081526020016100aa565b34801561010d57600080fd5b5061009e61011c366004610c08565b6104d1565b34801561012d57600080fd5b5061013761100581565b6040516001600160a01b0390911681526020016100aa565b34801561015b57600080fd5b5061009e61016a366004610c08565b610522565b6100f361017d366004610d2b565b61055a565b34801561018e57600080fd5b506101a261019d366004610e32565b6109b7565b604080516001600160a01b0390931683529015156020830152016100aa565b3480156101cd57600080fd5b5061013761100481565b60006001600160a01b0387161561030d576001600160a01b0387166323b872dd3330610203888a610e61565b6040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af1158015610257573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027b9190610e89565b506001600160a01b03871663095ea7b36110046102988789610e61565b6040516001600160e01b031960e085901b1681526001600160a01b03909216600483015260248201526044016020604051808303816000875af11580156102e3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103079190610e89565b5061036a565b6103178486610e61565b341461036a5760405162461bcd60e51b815260206004820181905260248201527f6d73672e76616c7565206e6f7420657175616c20616d6f756e74202b2066656560448201526064015b60405180910390fd5b60405163160d7c7360e01b81526110049063160d7c7390349061039b908b908b908b908b908b908b90600401610ef1565b60206040518083038185885af11580156103b9573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906103de9190610e89565b979650505050505050565b604051634ac3bdc360e01b815260009061100490634ac3bdc3906104139086908690600401610f46565b6020604051808303816000875af1158015610432573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104569190610e89565b9392505050565b604051638fefb76560e01b81526001600160a01b03831660048201526024810182905260009061100490638fefb76590604401602060405180830381865afa1580156104ad573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104569190610f68565b60405163a5df387560e01b8152600481018390526001600160a01b03821660248201526000906110049063a5df3875906044015b602060405180830381865afa158015610432573d6000803e3d6000fd5b60405163d5147e6d60e01b8152600481018390526001600160a01b03821660248201526000906110049063d5147e6d90604401610505565b600086518851146105b75760405162461bcd60e51b815260206004820152602160248201527f746f6b656e20616e6420616d6f756e74206c656e677468206e6f7420657175616044820152601b60fa1b6064820152608401610361565b60005b885181101561074f578881815181106105d5576105d5610f81565b60200260200101516001600160a01b03166323b872dd33308b85815181106105ff576105ff610f81565b60209081029190910101516040516001600160e01b031960e086901b1681526001600160a01b03938416600482015292909116602483015260448201526064016020604051808303816000875af115801561065e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106829190610e89565b5088818151811061069557610695610f81565b60200260200101516001600160a01b031663095ea7b36110048a84815181106106c0576106c0610f81565b60200260200101516040518363ffffffff1660e01b81526004016106f99291906001600160a01b03929092168252602082015260400190565b6020604051808303816000875af1158015610718573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061073c9190610e89565b508061074781610f97565b9150506105ba565b5060405163a8541c1760e01b8152600481018590526000906110059063a8541c179060240160e060405180830381865afa158015610791573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107b59190610fca565b60408101519091503490676170756e6469616960c01b141561083257816080015134101561081c5760405162461bcd60e51b81526020600482015260146024820152730dae6ce5cecc2d8eaca40dcdee840cadcdeeaced60631b6044820152606401610361565b608082015161082b9082611062565b9050610925565b604080830151905163e8e1f20560e01b81526000916110049163e8e1f205916108619160040190815260200190565b6040805180830381865afa15801561087d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108a19190611079565b5060808401516040516323b872dd60e01b815233600482015230602482015260448101919091529091506001600160a01b038216906323b872dd906064016020604051808303816000875af11580156108fe573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109229190610e89565b50505b6110046001600160a01b031663e6d69ede828e8e8e8e8e8e8e8e8e6040518b63ffffffff1660e01b8152600401610964999897969594939291906110e9565b60206040518083038185885af1158015610982573d6000803e3d6000fd5b50505050506040513d601f19601f820116820180604052508101906109a79190610f68565b9c9b505050505050505050505050565b60405163e8e1f20560e01b81526004810182905260009081906110049063e8e1f205906024016040805180830381865afa1580156109f9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a1d9190611079565b91509150915091565b6001600160a01b0381168114610a3b57600080fd5b50565b8035610a4981610a26565b919050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610a8d57610a8d610a4e565b604052919050565b600082601f830112610aa657600080fd5b813567ffffffffffffffff811115610ac057610ac0610a4e565b610ad3601f8201601f1916602001610a64565b818152846020838601011115610ae857600080fd5b816020850160208301376000918101602001919091529392505050565b60008060008060008060c08789031215610b1e57600080fd5b8635610b2981610a26565b9550602087013567ffffffffffffffff80821115610b4657600080fd5b610b528a838b01610a95565b965060408901359550606089013594506080890135935060a0890135915080821115610b7d57600080fd5b50610b8a89828a01610a95565b9150509295509295509295565b60008060408385031215610baa57600080fd5b823567ffffffffffffffff811115610bc157600080fd5b610bcd85828601610a95565b95602094909401359450505050565b60008060408385031215610bef57600080fd5b8235610bfa81610a26565b946020939093013593505050565b60008060408385031215610c1b57600080fd5b823591506020830135610c2d81610a26565b809150509250929050565b600067ffffffffffffffff821115610c5257610c52610a4e565b5060051b60200190565b600082601f830112610c6d57600080fd5b81356020610c82610c7d83610c38565b610a64565b82815260059290921b84018101918181019086841115610ca157600080fd5b8286015b84811015610cc5578035610cb881610a26565b8352918301918301610ca5565b509695505050505050565b600082601f830112610ce157600080fd5b81356020610cf1610c7d83610c38565b82815260059290921b84018101918181019086841115610d1057600080fd5b8286015b84811015610cc55780358352918301918301610d14565b60008060008060008060008060006101208a8c031215610d4a57600080fd5b893567ffffffffffffffff80821115610d6257600080fd5b610d6e8d838e01610a95565b9a50610d7c60208d01610a3e565b995060408c0135915080821115610d9257600080fd5b610d9e8d838e01610c5c565b985060608c0135915080821115610db457600080fd5b610dc08d838e01610cd0565b9750610dce60808d01610a3e565b965060a08c0135915080821115610de457600080fd5b610df08d838e01610a95565b955060c08c0135945060e08c013593506101008c0135915080821115610e1557600080fd5b50610e228c828d01610a95565b9150509295985092959850929598565b600060208284031215610e4457600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610e7457610e74610e4b565b500190565b80518015158114610a4957600080fd5b600060208284031215610e9b57600080fd5b61045682610e79565b6000815180845260005b81811015610eca57602081850181015186830182015201610eae565b81811115610edc576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b038716815260c060208201819052600090610f1590830188610ea4565b86604084015285606084015284608084015282810360a0840152610f398185610ea4565b9998505050505050505050565b604081526000610f596040830185610ea4565b90508260208301529392505050565b600060208284031215610f7a57600080fd5b5051919050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415610fab57610fab610e4b565b5060010190565b805167ffffffffffffffff81168114610a4957600080fd5b600060e08284031215610fdc57600080fd5b60405160e0810181811067ffffffffffffffff82111715610fff57610fff610a4e565b8060405250825181526020830151602082015260408301516040820152606083015161102a81610a26565b60608201526080838101519082015261104560a08401610fb2565b60a082015261105660c08401610fb2565b60c08201529392505050565b60008282101561107457611074610e4b565b500390565b6000806040838503121561108c57600080fd5b825161109781610a26565b91506110a560208401610e79565b90509250929050565b600081518084526020808501945080840160005b838110156110de578151875295820195908201906001016110c2565b509495945050505050565b60006101208083526110fd8184018d610ea4565b6001600160a01b038c811660208681019190915285830360408701528c518084528d82019450909283019060005b8181101561114957855184168352948401949184019160010161112b565b5050858103606087015261115d818d6110ae565b935050505061117760808401896001600160a01b03169052565b82810360a08401526111898188610ea4565b90508560c08401528460e08401528281036101008401526109a78185610ea456fea2646970667358221220270564294290823af031a0023c1cb749df0f4188f466cb8d35b8b5cc33a6114664736f6c634300080a0033",
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
