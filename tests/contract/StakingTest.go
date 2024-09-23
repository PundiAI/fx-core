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

// StakingTestMetaData contains all meta data concerning the StakingTest contract.
var StakingTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"ApproveShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valSrc\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valDst\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"RedelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"TransferShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"UndelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"STAKING_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowanceShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approveShares\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"delegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_valSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_valDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"redelegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"slashingInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_jailed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_missed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFromShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"undelegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIStaking.ValidatorSortBy\",\"name\":\"_sortBy\",\"type\":\"uint8\"}],\"name\":\"validatorList\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"validatorShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610f6e806100206000396000f3fe6080604052600436106100dd5760003560e01c80637b625c0f1161007f578063d5c498eb11610059578063d5c498eb1461029b578063dc6ffc7d146102bb578063de2b3451146102db578063ee226c66146102fb57600080fd5b80637b625c0f146102155780638c20570b14610235578063bf98d7721461026357600080fd5b806349da433e116100bb57806349da433e1461017b5780634e94633a146101ab57806351af513a146101e25780636d7880351461020257600080fd5b8063029c0a51146100e2578063161298c11461011857806331fb67c21461014d575b600080fd5b3480156100ee57600080fd5b506101026100fd366004610841565b61031b565b60405161010f91906108be565b60405180910390f35b34801561012457600080fd5b50610138610133366004610a01565b61038e565b6040805192835260208301919091520161010f565b34801561015957600080fd5b5061016d610168366004610a58565b61040a565b60405190815260200161010f565b34801561018757600080fd5b5061019b610196366004610a01565b610475565b604051901515815260200161010f565b3480156101b757600080fd5b506101cb6101c6366004610a58565b6104ec565b60408051921515835260208301919091520161010f565b3480156101ee57600080fd5b5061016d6101fd366004610a8d565b61055f565b61019b610210366004610adb565b6105d1565b34801561022157600080fd5b5061016d610230366004610b20565b610688565b34801561024157600080fd5b5061024b61100381565b6040516001600160a01b03909116815260200161010f565b34801561026f57600080fd5b5061016d61027e366004610a58565b805160208183018101805160008252928201919093012091525481565b3480156102a757600080fd5b506101386102b6366004610a8d565b6106f5565b3480156102c757600080fd5b506101386102d6366004610b7e565b61076c565b3480156102e757600080fd5b5061019b6102f6366004610adb565b6107eb565b34801561030757600080fd5b5061019b610316366004610be4565b610815565b60405163029c0a5160e01b81526060906110039063029c0a5190610343908590600401610c51565b600060405180830381865afa158015610360573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103889190810190610c79565b92915050565b60405163161298c160e01b815260009081906110039063161298c1906103bc90889088908890600401610d78565b60408051808303816000875af11580156103da573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103fe9190610da6565b91509150935093915050565b6040516318fdb3e160e11b8152600090611003906331fb67c290610432908590600401610dca565b6020604051808303816000875af1158015610451573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103889190610ddd565b6040516324ed219f60e11b8152600090611003906349da433e906104a190879087908790600401610d78565b6020604051808303816000875af11580156104c0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e49190610e06565b949350505050565b60405163274a319d60e11b8152600090819061100390634e94633a90610516908690600401610dca565b6040805180830381865afa158015610532573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105569190610e21565b91509150915091565b6040516328d7a89d60e11b8152600090611003906351af513a906105899086908690600401610e4d565b602060405180830381865afa1580156105a6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105ca9190610ddd565b9392505050565b60008147101561061e5760405162461bcd60e51b8152602060048201526014602482015273696e73756666696369656e742062616c616e636560601b604482015260640160405180910390fd5b604051636d78803560e01b815261100390636d788035906106459086908690600401610e77565b6020604051808303816000875af1158015610664573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105ca9190610e06565b604051637b625c0f60e01b815260009061100390637b625c0f906106b490879087908790600401610e99565b602060405180830381865afa1580156106d1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104e49190610ddd565b60405163d5c498eb60e01b815260009081906110039063d5c498eb906107219087908790600401610e4d565b6040805180830381865afa15801561073d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107619190610da6565b915091509250929050565b60405163dc6ffc7d60e01b815260009081906110039063dc6ffc7d9061079c908990899089908990600401610ecc565b60408051808303816000875af11580156107ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107de9190610da6565b9150915094509492505050565b60405163de2b345160e01b81526000906110039063de2b3451906106459086908690600401610e77565b604051637711363360e11b81526000906110039063ee226c66906104a190879087908790600401610f02565b60006020828403121561085357600080fd5b8135600281106105ca57600080fd5b60005b8381101561087d578181015183820152602001610865565b8381111561088c576000848401525b50505050565b600081518084526108aa816020860160208601610862565b601f01601f19169290920160200192915050565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101561091357603f19888603018452610901858351610892565b945092850192908501906001016108e5565b5092979650505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561095f5761095f610920565b604052919050565b600067ffffffffffffffff82111561098157610981610920565b50601f01601f191660200190565b600082601f8301126109a057600080fd5b81356109b36109ae82610967565b610936565b8181528460208386010111156109c857600080fd5b816020850160208301376000918101602001919091529392505050565b80356001600160a01b03811681146109fc57600080fd5b919050565b600080600060608486031215610a1657600080fd5b833567ffffffffffffffff811115610a2d57600080fd5b610a398682870161098f565b935050610a48602085016109e5565b9150604084013590509250925092565b600060208284031215610a6a57600080fd5b813567ffffffffffffffff811115610a8157600080fd5b6104e48482850161098f565b60008060408385031215610aa057600080fd5b823567ffffffffffffffff811115610ab757600080fd5b610ac38582860161098f565b925050610ad2602084016109e5565b90509250929050565b60008060408385031215610aee57600080fd5b823567ffffffffffffffff811115610b0557600080fd5b610b118582860161098f565b95602094909401359450505050565b600080600060608486031215610b3557600080fd5b833567ffffffffffffffff811115610b4c57600080fd5b610b588682870161098f565b935050610b67602085016109e5565b9150610b75604085016109e5565b90509250925092565b60008060008060808587031215610b9457600080fd5b843567ffffffffffffffff811115610bab57600080fd5b610bb78782880161098f565b945050610bc6602086016109e5565b9250610bd4604086016109e5565b9396929550929360600135925050565b600080600060608486031215610bf957600080fd5b833567ffffffffffffffff80821115610c1157600080fd5b610c1d8783880161098f565b94506020860135915080821115610c3357600080fd5b50610c408682870161098f565b925050604084013590509250925092565b6020810160028310610c7357634e487b7160e01b600052602160045260246000fd5b91905290565b60006020808385031215610c8c57600080fd5b825167ffffffffffffffff80821115610ca457600080fd5b818501915085601f830112610cb857600080fd5b815181811115610cca57610cca610920565b8060051b610cd9858201610936565b9182528381018501918581019089841115610cf357600080fd5b86860192505b83831015610d6b57825185811115610d115760008081fd5b8601603f81018b13610d235760008081fd5b878101516040610d356109ae83610967565b8281528d82848601011115610d4a5760008081fd5b610d59838c8301848701610862565b85525050509186019190860190610cf9565b9998505050505050505050565b606081526000610d8b6060830186610892565b6001600160a01b039490941660208301525060400152919050565b60008060408385031215610db957600080fd5b505080516020909101519092909150565b6020815260006105ca6020830184610892565b600060208284031215610def57600080fd5b5051919050565b805180151581146109fc57600080fd5b600060208284031215610e1857600080fd5b6105ca82610df6565b60008060408385031215610e3457600080fd5b610e3d83610df6565b9150602083015190509250929050565b604081526000610e606040830185610892565b905060018060a01b03831660208301529392505050565b604081526000610e8a6040830185610892565b90508260208301529392505050565b606081526000610eac6060830186610892565b6001600160a01b0394851660208401529290931660409091015292915050565b608081526000610edf6080830187610892565b6001600160a01b0395861660208401529390941660408201526060015292915050565b606081526000610f156060830186610892565b8281036020840152610f278186610892565b91505082604083015294935050505056fea26469706673582212209be27dfe6387d6c89a031fefd642735437cb92cd4828c9c844abd957f3f556e964736f6c634300080a0033",
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
	parsed, err := StakingTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// STAKINGADDRESS is a free data retrieval call binding the contract method 0x8c20570b.
//
// Solidity: function STAKING_ADDRESS() view returns(address)
func (_StakingTest *StakingTestCaller) STAKINGADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "STAKING_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// STAKINGADDRESS is a free data retrieval call binding the contract method 0x8c20570b.
//
// Solidity: function STAKING_ADDRESS() view returns(address)
func (_StakingTest *StakingTestSession) STAKINGADDRESS() (common.Address, error) {
	return _StakingTest.Contract.STAKINGADDRESS(&_StakingTest.CallOpts)
}

// STAKINGADDRESS is a free data retrieval call binding the contract method 0x8c20570b.
//
// Solidity: function STAKING_ADDRESS() view returns(address)
func (_StakingTest *StakingTestCallerSession) STAKINGADDRESS() (common.Address, error) {
	return _StakingTest.Contract.STAKINGADDRESS(&_StakingTest.CallOpts)
}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestCaller) AllowanceShares(opts *bind.CallOpts, _val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "allowanceShares", _val, _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestSession) AllowanceShares(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StakingTest.Contract.AllowanceShares(&_StakingTest.CallOpts, _val, _owner, _spender)
}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) AllowanceShares(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StakingTest.Contract.AllowanceShares(&_StakingTest.CallOpts, _val, _owner, _spender)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestCaller) Delegation(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "delegation", _val, _del)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Delegation(_val string, _del common.Address) (*big.Int, *big.Int, error) {
	return _StakingTest.Contract.Delegation(&_StakingTest.CallOpts, _val, _del)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestCallerSession) Delegation(_val string, _del common.Address) (*big.Int, *big.Int, error) {
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

// SlashingInfo is a free data retrieval call binding the contract method 0x4e94633a.
//
// Solidity: function slashingInfo(string _val) view returns(bool _jailed, uint256 _missed)
func (_StakingTest *StakingTestCaller) SlashingInfo(opts *bind.CallOpts, _val string) (struct {
	Jailed bool
	Missed *big.Int
}, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "slashingInfo", _val)

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
func (_StakingTest *StakingTestSession) SlashingInfo(_val string) (struct {
	Jailed bool
	Missed *big.Int
}, error) {
	return _StakingTest.Contract.SlashingInfo(&_StakingTest.CallOpts, _val)
}

// SlashingInfo is a free data retrieval call binding the contract method 0x4e94633a.
//
// Solidity: function slashingInfo(string _val) view returns(bool _jailed, uint256 _missed)
func (_StakingTest *StakingTestCallerSession) SlashingInfo(_val string) (struct {
	Jailed bool
	Missed *big.Int
}, error) {
	return _StakingTest.Contract.SlashingInfo(&_StakingTest.CallOpts, _val)
}

// ValidatorList is a free data retrieval call binding the contract method 0x029c0a51.
//
// Solidity: function validatorList(uint8 _sortBy) view returns(string[])
func (_StakingTest *StakingTestCaller) ValidatorList(opts *bind.CallOpts, _sortBy uint8) ([]string, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "validatorList", _sortBy)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// ValidatorList is a free data retrieval call binding the contract method 0x029c0a51.
//
// Solidity: function validatorList(uint8 _sortBy) view returns(string[])
func (_StakingTest *StakingTestSession) ValidatorList(_sortBy uint8) ([]string, error) {
	return _StakingTest.Contract.ValidatorList(&_StakingTest.CallOpts, _sortBy)
}

// ValidatorList is a free data retrieval call binding the contract method 0x029c0a51.
//
// Solidity: function validatorList(uint8 _sortBy) view returns(string[])
func (_StakingTest *StakingTestCallerSession) ValidatorList(_sortBy uint8) ([]string, error) {
	return _StakingTest.Contract.ValidatorList(&_StakingTest.CallOpts, _sortBy)
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

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestTransactor) ApproveShares(opts *bind.TransactOpts, _val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "approveShares", _val, _spender, _shares)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestSession) ApproveShares(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.ApproveShares(&_StakingTest.TransactOpts, _val, _spender, _shares)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestTransactorSession) ApproveShares(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.ApproveShares(&_StakingTest.TransactOpts, _val, _spender, _shares)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) payable returns(bool _result)
func (_StakingTest *StakingTestTransactor) DelegateV2(opts *bind.TransactOpts, _val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "delegateV2", _val, _amount)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) payable returns(bool _result)
func (_StakingTest *StakingTestSession) DelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.DelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) payable returns(bool _result)
func (_StakingTest *StakingTestTransactorSession) DelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.DelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactor) RedelegateV2(opts *bind.TransactOpts, _valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "redelegateV2", _valSrc, _valDst, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestSession) RedelegateV2(_valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.RedelegateV2(&_StakingTest.TransactOpts, _valSrc, _valDst, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactorSession) RedelegateV2(_valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.RedelegateV2(&_StakingTest.TransactOpts, _valSrc, _valDst, _amount)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) TransferFromShares(opts *bind.TransactOpts, _val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "transferFromShares", _val, _from, _to, _shares)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) TransferFromShares(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferFromShares(&_StakingTest.TransactOpts, _val, _from, _to, _shares)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) TransferFromShares(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferFromShares(&_StakingTest.TransactOpts, _val, _from, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) TransferShares(opts *bind.TransactOpts, _val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "transferShares", _val, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) TransferShares(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferShares(&_StakingTest.TransactOpts, _val, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) TransferShares(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferShares(&_StakingTest.TransactOpts, _val, _to, _shares)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactor) UndelegateV2(opts *bind.TransactOpts, _val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "undelegateV2", _val, _amount)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestSession) UndelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.UndelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactorSession) UndelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.UndelegateV2(&_StakingTest.TransactOpts, _val, _amount)
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

// StakingTestApproveSharesIterator is returned from FilterApproveShares and is used to iterate over the raw logs and unpacked data for ApproveShares events raised by the StakingTest contract.
type StakingTestApproveSharesIterator struct {
	Event *StakingTestApproveShares // Event containing the contract specifics and raw log

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
func (it *StakingTestApproveSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestApproveShares)
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
		it.Event = new(StakingTestApproveShares)
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
func (it *StakingTestApproveSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestApproveSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestApproveShares represents a ApproveShares event raised by the StakingTest contract.
type StakingTestApproveShares struct {
	Owner     common.Address
	Spender   common.Address
	Validator string
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterApproveShares is a free log retrieval operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) FilterApproveShares(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StakingTestApproveSharesIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "ApproveShares", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestApproveSharesIterator{contract: _StakingTest.contract, event: "ApproveShares", logs: logs, sub: sub}, nil
}

// WatchApproveShares is a free log subscription operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) WatchApproveShares(opts *bind.WatchOpts, sink chan<- *StakingTestApproveShares, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "ApproveShares", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestApproveShares)
				if err := _StakingTest.contract.UnpackLog(event, "ApproveShares", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseApproveShares(log types.Log) (*StakingTestApproveShares, error) {
	event := new(StakingTestApproveShares)
	if err := _StakingTest.contract.UnpackLog(event, "ApproveShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestDelegateV2Iterator is returned from FilterDelegateV2 and is used to iterate over the raw logs and unpacked data for DelegateV2 events raised by the StakingTest contract.
type StakingTestDelegateV2Iterator struct {
	Event *StakingTestDelegateV2 // Event containing the contract specifics and raw log

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
func (it *StakingTestDelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestDelegateV2)
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
		it.Event = new(StakingTestDelegateV2)
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
func (it *StakingTestDelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestDelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestDelegateV2 represents a DelegateV2 event raised by the StakingTest contract.
type StakingTestDelegateV2 struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegateV2 is a free log retrieval operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_StakingTest *StakingTestFilterer) FilterDelegateV2(opts *bind.FilterOpts, delegator []common.Address) (*StakingTestDelegateV2Iterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "DelegateV2", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestDelegateV2Iterator{contract: _StakingTest.contract, event: "DelegateV2", logs: logs, sub: sub}, nil
}

// WatchDelegateV2 is a free log subscription operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_StakingTest *StakingTestFilterer) WatchDelegateV2(opts *bind.WatchOpts, sink chan<- *StakingTestDelegateV2, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "DelegateV2", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestDelegateV2)
				if err := _StakingTest.contract.UnpackLog(event, "DelegateV2", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseDelegateV2(log types.Log) (*StakingTestDelegateV2, error) {
	event := new(StakingTestDelegateV2)
	if err := _StakingTest.contract.UnpackLog(event, "DelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestRedelegateV2Iterator is returned from FilterRedelegateV2 and is used to iterate over the raw logs and unpacked data for RedelegateV2 events raised by the StakingTest contract.
type StakingTestRedelegateV2Iterator struct {
	Event *StakingTestRedelegateV2 // Event containing the contract specifics and raw log

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
func (it *StakingTestRedelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestRedelegateV2)
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
		it.Event = new(StakingTestRedelegateV2)
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
func (it *StakingTestRedelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestRedelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestRedelegateV2 represents a RedelegateV2 event raised by the StakingTest contract.
type StakingTestRedelegateV2 struct {
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
func (_StakingTest *StakingTestFilterer) FilterRedelegateV2(opts *bind.FilterOpts, sender []common.Address) (*StakingTestRedelegateV2Iterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "RedelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestRedelegateV2Iterator{contract: _StakingTest.contract, event: "RedelegateV2", logs: logs, sub: sub}, nil
}

// WatchRedelegateV2 is a free log subscription operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchRedelegateV2(opts *bind.WatchOpts, sink chan<- *StakingTestRedelegateV2, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "RedelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestRedelegateV2)
				if err := _StakingTest.contract.UnpackLog(event, "RedelegateV2", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseRedelegateV2(log types.Log) (*StakingTestRedelegateV2, error) {
	event := new(StakingTestRedelegateV2)
	if err := _StakingTest.contract.UnpackLog(event, "RedelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestTransferSharesIterator is returned from FilterTransferShares and is used to iterate over the raw logs and unpacked data for TransferShares events raised by the StakingTest contract.
type StakingTestTransferSharesIterator struct {
	Event *StakingTestTransferShares // Event containing the contract specifics and raw log

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
func (it *StakingTestTransferSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestTransferShares)
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
		it.Event = new(StakingTestTransferShares)
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
func (it *StakingTestTransferSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestTransferSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestTransferShares represents a TransferShares event raised by the StakingTest contract.
type StakingTestTransferShares struct {
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
func (_StakingTest *StakingTestFilterer) FilterTransferShares(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingTestTransferSharesIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "TransferShares", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestTransferSharesIterator{contract: _StakingTest.contract, event: "TransferShares", logs: logs, sub: sub}, nil
}

// WatchTransferShares is a free log subscription operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) WatchTransferShares(opts *bind.WatchOpts, sink chan<- *StakingTestTransferShares, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "TransferShares", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestTransferShares)
				if err := _StakingTest.contract.UnpackLog(event, "TransferShares", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseTransferShares(log types.Log) (*StakingTestTransferShares, error) {
	event := new(StakingTestTransferShares)
	if err := _StakingTest.contract.UnpackLog(event, "TransferShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestUndelegateV2Iterator is returned from FilterUndelegateV2 and is used to iterate over the raw logs and unpacked data for UndelegateV2 events raised by the StakingTest contract.
type StakingTestUndelegateV2Iterator struct {
	Event *StakingTestUndelegateV2 // Event containing the contract specifics and raw log

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
func (it *StakingTestUndelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestUndelegateV2)
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
		it.Event = new(StakingTestUndelegateV2)
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
func (it *StakingTestUndelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestUndelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestUndelegateV2 represents a UndelegateV2 event raised by the StakingTest contract.
type StakingTestUndelegateV2 struct {
	Sender         common.Address
	Validator      string
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUndelegateV2 is a free log retrieval operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) FilterUndelegateV2(opts *bind.FilterOpts, sender []common.Address) (*StakingTestUndelegateV2Iterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "UndelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestUndelegateV2Iterator{contract: _StakingTest.contract, event: "UndelegateV2", logs: logs, sub: sub}, nil
}

// WatchUndelegateV2 is a free log subscription operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchUndelegateV2(opts *bind.WatchOpts, sink chan<- *StakingTestUndelegateV2, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "UndelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestUndelegateV2)
				if err := _StakingTest.contract.UnpackLog(event, "UndelegateV2", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseUndelegateV2(log types.Log) (*StakingTestUndelegateV2, error) {
	event := new(StakingTestUndelegateV2)
	if err := _StakingTest.contract.UnpackLog(event, "UndelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the StakingTest contract.
type StakingTestWithdrawIterator struct {
	Event *StakingTestWithdraw // Event containing the contract specifics and raw log

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
func (it *StakingTestWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestWithdraw)
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
		it.Event = new(StakingTestWithdraw)
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
func (it *StakingTestWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestWithdraw represents a Withdraw event raised by the StakingTest contract.
type StakingTestWithdraw struct {
	Sender    common.Address
	Validator string
	Reward    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_StakingTest *StakingTestFilterer) FilterWithdraw(opts *bind.FilterOpts, sender []common.Address) (*StakingTestWithdrawIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestWithdrawIterator{contract: _StakingTest.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_StakingTest *StakingTestFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StakingTestWithdraw, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestWithdraw)
				if err := _StakingTest.contract.UnpackLog(event, "Withdraw", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseWithdraw(log types.Log) (*StakingTestWithdraw, error) {
	event := new(StakingTestWithdraw)
	if err := _StakingTest.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
