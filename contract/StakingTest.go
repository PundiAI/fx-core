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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"ApproveShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"TransferShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Undelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowanceShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approveShares\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFromShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"validatorShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506112fa806100206000396000f3fe6080604052600436106100915760003560e01c80638dfc8897116100595780638dfc88971461016e5780639ddb511a146101a9578063bf98d772146101bc578063d5c498eb146101f4578063dc6ffc7d1461021457610091565b8063161298c11461009657806331fb67c2146100d057806349da433e146100fe57806351af513a1461012e5780637b625c0f1461014e575b600080fd5b3480156100a257600080fd5b506100b66100b1366004610f65565b610234565b604080519283526020830191909152015b60405180910390f35b3480156100dc57600080fd5b506100f06100eb366004610dbc565b610253565b6040519081526020016100c7565b34801561010a57600080fd5b5061011e610119366004610f65565b610268565b60405190151581526020016100c7565b34801561013a57600080fd5b506100f0610149366004610e59565b61027f565b34801561015a57600080fd5b506100f0610169366004610ea5565b610292565b34801561017a57600080fd5b5061018e610189366004610fba565b6102a7565b604080519384526020840192909252908201526060016100c7565b6100b66101b7366004610dbc565b610300565b3480156101c857600080fd5b506100f06101d7366004610dbc565b805160208183018101805160008252928201919093012091525481565b34801561020057600080fd5b506100b661020f366004610e59565b610351565b34801561022057600080fd5b506100b661022f366004610f01565b610369565b60008060008061024587878761038a565b909890975095505050505050565b60008061025f83610441565b9150505b919050565b6000806102768585856104e0565b95945050505050565b600061028b8383610591565b9392505050565b600061029f84848461063d565b949350505050565b6000806000806000806102ba88886106ea565b925092509250866000896040516102d19190611091565b908152602001604051809103902060008282546102ee9190611251565b90915550929891975095509350505050565b600080600080610310853461079d565b91509150816000866040516103259190611091565b908152602001604051809103902060008282546103429190611239565b90915550919350915050915091565b60008061035e848461084b565b915091509250929050565b60008060008061037b888888886108e3565b90999098509650505050505050565b600080808061100361039d8888886109a3565b6040516103aa9190611091565b6000604051808303816000865af19150503d80600081146103e7576040519150601f19603f3d011682016040523d82523d6000602084013e6103ec565b606091505b509150915061042a8282604051806040016040528060168152602001751d1c985b9cd9995c881cda185c995cc819985a5b195960521b8152506109ed565b61043381610a75565b935093505050935093915050565b6000808061100361045185610a9b565b60405161045e9190611091565b6000604051808303816000865af19150503d806000811461049b576040519150601f19603f3d011682016040523d82523d6000602084013e6104a0565b606091505b50915091506104d782826040518060400160405280600f81526020016e1dda5d1a191c985dc819985a5b1959608a1b8152506109ed565b61029f81610adf565b600080806110036104f2878787610af6565b6040516104ff9190611091565b6000604051808303816000865af19150503d806000811461053c576040519150601f19603f3d011682016040523d82523d6000602084013e610541565b606091505b509150915061057e828260405180604001604052806015815260200174185c1c1c9bdd99481cda185c995cc819985a5b1959605a1b8152506109ed565b61058781610b40565b9695505050505050565b600080806110036105a28686610b57565b6040516105af9190611091565b600060405180830381855afa9150503d80600081146105ea576040519150601f19603f3d011682016040523d82523d6000602084013e6105ef565b606091505b509150915061063482826040518060400160405280601881526020017f64656c65676174696f6e52657761726473206661696c656400000000000000008152506109ed565b61027681610adf565b6000808061100361064f878787610b9e565b60405161065c9190611091565b600060405180830381855afa9150503d8060008114610697576040519150601f19603f3d011682016040523d82523d6000602084013e61069c565b606091505b50915091506106e182826040518060400160405280601781526020017f616c6c6f77616e636520736861726573206661696c65640000000000000000008152506109ed565b61058781610adf565b6000808080806110036106fd8888610be8565b60405161070a9190611091565b6000604051808303816000865af19150503d8060008114610747576040519150601f19603f3d011682016040523d82523d6000602084013e61074c565b606091505b50915091506107858282604051806040016040528060118152602001701d5b99195b1959d85d194819985a5b1959607a1b8152506109ed565b61078e81610c2f565b94509450945050509250925092565b6000808080611003856107af88610c5c565b6040516107bc9190611091565b60006040518083038185875af1925050503d80600081146107f9576040519150601f19603f3d011682016040523d82523d6000602084013e6107fe565b606091505b509150915061083582826040518060400160405280600f81526020016e19195b1959d85d194819985a5b1959608a1b8152506109ed565b61083e81610a75565b9350935050509250929050565b600080808061100361085d8787610ca0565b60405161086a9190611091565b600060405180830381855afa9150503d80600081146108a5576040519150601f19603f3d011682016040523d82523d6000602084013e6108aa565b606091505b509150915061083582826040518060400160405280601181526020017019195b1959d85d1a5bdb8819985a5b1959607a1b8152506109ed565b60008080806110036108f789898989610ce7565b6040516109049190611091565b6000604051808303816000865af19150503d8060008114610941576040519150601f19603f3d011682016040523d82523d6000602084013e610946565b606091505b509150915061098b82826040518060400160405280601a81526020017f7472616e7366657246726f6d20736861726573206661696c65640000000000008152506109ed565b61099481610a75565b93509350505094509492505050565b60608383836040516024016109ba93929190611190565b60408051601f198184030181529190526020810180516001600160e01b031663161298c160e01b17905290509392505050565b82610a7057600082806020019051810190610a089190610def565b9050600182511015610a37578060405162461bcd60e51b8152600401610a2e91906110ea565b60405180910390fd5b8181604051602001610a4a9291906110ad565b60408051601f198184030181529082905262461bcd60e51b8252610a2e916004016110ea565b505050565b60008060008084806020019051810190610a8f9190611015565b90945092505050915091565b606081604051602401610aae91906110ea565b60408051601f198184030181529190526020810180516001600160e01b03166318fdb3e160e11b1790529050919050565b6000808280602001905181019061025f9190610ffd565b6060838383604051602401610b0d93929190611190565b60408051601f198184030181529190526020810180516001600160e01b03166324ed219f60e11b17905290509392505050565b6000808280602001905181019061025f9190610d9c565b60608282604051602401610b6c9291906110fd565b60408051601f198184030181529190526020810180516001600160e01b03166328d7a89d60e11b179052905092915050565b6060838383604051602401610bb593929190611127565b60408051601f198184030181529190526020810180516001600160e01b0316637b625c0f60e01b17905290509392505050565b60608282604051602401610bfd9291906111be565b60408051601f198184030181529190526020810180516001600160e01b0316638dfc889760e01b179052905092915050565b60008060008060008086806020019051810190610c4c9190611038565b9199909850909650945050505050565b606081604051602401610c6f91906110ea565b60408051601f198184030181529190526020810180516001600160e01b0316634eeda88d60e11b1790529050919050565b60608282604051602401610cb59291906110fd565b60408051601f198184030181529190526020810180516001600160e01b031663d5c498eb60e01b179052905092915050565b606084848484604051602401610d00949392919061115a565b60408051601f198184030181529190526020810180516001600160e01b031663dc6ffc7d60e01b1790529050949350505050565b80356001600160a01b038116811461026357600080fd5b600082601f830112610d5b578081fd5b8135610d6e610d6982611211565b6111e0565b818152846020838601011115610d82578283fd5b816020850160208301379081016020019190915292915050565b600060208284031215610dad578081fd5b8151801515811461028b578182fd5b600060208284031215610dcd578081fd5b813567ffffffffffffffff811115610de3578182fd5b61029f84828501610d4b565b600060208284031215610e00578081fd5b815167ffffffffffffffff811115610e16578182fd5b8201601f81018413610e26578182fd5b8051610e34610d6982611211565b818152856020838501011115610e48578384fd5b610276826020830160208601611268565b60008060408385031215610e6b578081fd5b823567ffffffffffffffff811115610e81578182fd5b610e8d85828601610d4b565b925050610e9c60208401610d34565b90509250929050565b600080600060608486031215610eb9578081fd5b833567ffffffffffffffff811115610ecf578182fd5b610edb86828701610d4b565b935050610eea60208501610d34565b9150610ef860408501610d34565b90509250925092565b60008060008060808587031215610f16578081fd5b843567ffffffffffffffff811115610f2c578182fd5b610f3887828801610d4b565b945050610f4760208601610d34565b9250610f5560408601610d34565b9396929550929360600135925050565b600080600060608486031215610f79578283fd5b833567ffffffffffffffff811115610f8f578384fd5b610f9b86828701610d4b565b935050610faa60208501610d34565b9150604084013590509250925092565b60008060408385031215610fcc578182fd5b823567ffffffffffffffff811115610fe2578283fd5b610fee85828601610d4b565b95602094909401359450505050565b60006020828403121561100e578081fd5b5051919050565b60008060408385031215611027578182fd5b505080516020909101519092909150565b60008060006060848603121561104c578283fd5b8351925060208401519150604084015190509250925092565b6000815180845261107d816020860160208601611268565b601f01601f19169290920160200192915050565b600082516110a3818460208701611268565b9190910192915050565b600083516110bf818460208801611268565b6101d160f51b90830190815283516110de816002840160208801611268565b01600201949350505050565b60006020825261028b6020830184611065565b6000604082526111106040830185611065565b905060018060a01b03831660208301529392505050565b60006060825261113a6060830186611065565b6001600160a01b0394851660208401529290931660409091015292915050565b60006080825261116d6080830187611065565b6001600160a01b0395861660208401529390941660408201526060015292915050565b6000606082526111a36060830186611065565b6001600160a01b039490941660208301525060400152919050565b6000604082526111d16040830185611065565b90508260208301529392505050565b604051601f8201601f1916810167ffffffffffffffff81118282101715611209576112096112ae565b604052919050565b600067ffffffffffffffff82111561122b5761122b6112ae565b50601f01601f191660200190565b6000821982111561124c5761124c611298565b500190565b60008282101561126357611263611298565b500390565b60005b8381101561128357818101518382015260200161126b565b83811115611292576000848401525b50505050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fdfea26469706673582212205304a5533c4f43ed0d4d5d2c414b816d845350c1321e0f88712c2ff73251285f64736f6c63430008020033",
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

// StakingTestDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the StakingTest contract.
type StakingTestDelegateIterator struct {
	Event *StakingTestDelegate // Event containing the contract specifics and raw log

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
func (it *StakingTestDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestDelegate)
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
		it.Event = new(StakingTestDelegate)
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
func (it *StakingTestDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestDelegate represents a Delegate event raised by the StakingTest contract.
type StakingTestDelegate struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_StakingTest *StakingTestFilterer) FilterDelegate(opts *bind.FilterOpts, delegator []common.Address) (*StakingTestDelegateIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Delegate", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestDelegateIterator{contract: _StakingTest.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_StakingTest *StakingTestFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *StakingTestDelegate, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Delegate", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestDelegate)
				if err := _StakingTest.contract.UnpackLog(event, "Delegate", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseDelegate(log types.Log) (*StakingTestDelegate, error) {
	event := new(StakingTestDelegate)
	if err := _StakingTest.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// StakingTestUndelegateIterator is returned from FilterUndelegate and is used to iterate over the raw logs and unpacked data for Undelegate events raised by the StakingTest contract.
type StakingTestUndelegateIterator struct {
	Event *StakingTestUndelegate // Event containing the contract specifics and raw log

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
func (it *StakingTestUndelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestUndelegate)
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
		it.Event = new(StakingTestUndelegate)
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
func (it *StakingTestUndelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestUndelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestUndelegate represents a Undelegate event raised by the StakingTest contract.
type StakingTestUndelegate struct {
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
func (_StakingTest *StakingTestFilterer) FilterUndelegate(opts *bind.FilterOpts, sender []common.Address) (*StakingTestUndelegateIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Undelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestUndelegateIterator{contract: _StakingTest.contract, event: "Undelegate", logs: logs, sub: sub}, nil
}

// WatchUndelegate is a free log subscription operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchUndelegate(opts *bind.WatchOpts, sink chan<- *StakingTestUndelegate, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Undelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestUndelegate)
				if err := _StakingTest.contract.UnpackLog(event, "Undelegate", log); err != nil {
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
func (_StakingTest *StakingTestFilterer) ParseUndelegate(log types.Log) (*StakingTestUndelegate, error) {
	event := new(StakingTestUndelegate)
	if err := _StakingTest.contract.UnpackLog(event, "Undelegate", log); err != nil {
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
