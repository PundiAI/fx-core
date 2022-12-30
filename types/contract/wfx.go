// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// WFXABI is the input ABI used to generate the binding from.
const WFXABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"TransferCrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"transferCrossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"

// WFXBin is the compiled bytecode used for deploying new contracts.
var WFXBin = "0x60a06040523060601b60805234801561001757600080fd5b5060805160601c611c626100526000396000818161060001528181610640015281816107180152818161075801526107e70152611c626000f3fe6080604052600436106101395760003560e01c80638da5cb5b116100ab578063c5cb9b511161006f578063c5cb9b5114610364578063d0e30db014610148578063dd62ed3e14610377578063de7ea79d146103bd578063f2fde38b146103dd578063f3fef3a3146103fd57610148565b80638da5cb5b146102bf57806395d89b41146102f15780639dc29fac14610306578063a9059cbb14610326578063b86d52981461034657610148565b80633659cfe6116100fd5780633659cfe61461020c57806340c10f191461022c5780634f1ef2861461024c57806352d1902d1461025f57806370a0823114610274578063715018a6146102aa57610148565b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab57806323b872dd146101ca578063313ce567146101ea57610148565b366101485761014661041d565b005b61014661041d565b34801561015c57600080fd5b5061016561045e565b60405161017291906119b4565b60405180910390f35b34801561018757600080fd5b5061019b610196366004611865565b6104f0565b6040519015158152602001610172565b3480156101b757600080fd5b5060cc545b604051908152602001610172565b3480156101d657600080fd5b5061019b6101e53660046117c4565b610546565b3480156101f657600080fd5b5060cb5460405160ff9091168152602001610172565b34801561021857600080fd5b50610146610227366004611745565b6105f5565b34801561023857600080fd5b50610146610247366004611865565b6106d5565b61014661025a366004611804565b61070d565b34801561026b57600080fd5b506101bc6107da565b34801561028057600080fd5b506101bc61028f366004611745565b6001600160a01b0316600090815260cd602052604090205490565b3480156102b657600080fd5b5061014661088d565b3480156102cb57600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610172565b3480156102fd57600080fd5b506101656108c3565b34801561031257600080fd5b50610146610321366004611865565b6108d2565b34801561033257600080fd5b5061019b610341366004611865565b610906565b34801561035257600080fd5b5060cf546001600160a01b03166102d9565b61019b61037236600461191a565b61091c565b34801561038357600080fd5b506101bc61039236600461178c565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103c957600080fd5b506101466103d836600461188f565b6109e0565b3480156103e957600080fd5b506101466103f8366004611745565b610aff565b34801561040957600080fd5b50610146610418366004611761565b610b97565b6104273334610c1d565b60405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b606060c9805461046d90611b69565b80601f016020809104026020016040519081016040528092919081815260200182805461049990611b69565b80156104e65780601f106104bb576101008083540402835291602001916104e6565b820191906000526020600020905b8154815290600101906020018083116104c957829003601f168201915b5050505050905090565b60006104fd338484610cf5565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156105c95760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b6105dd85336105d88685611b26565b610cf5565b6105e8858585610d77565b60019150505b9392505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016141561063e5760405162461bcd60e51b81526004016105c0906119f6565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610687600080516020611be6833981519152546001600160a01b031690565b6001600160a01b0316146106ad5760405162461bcd60e51b81526004016105c090611a42565b6106b681610f26565b604080516000808252602082019092526106d291839190610f50565b50565b6097546001600160a01b031633146106ff5760405162461bcd60e51b81526004016105c090611a8e565b6107098282610c1d565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156107565760405162461bcd60e51b81526004016105c0906119f6565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661079f600080516020611be6833981519152546001600160a01b031690565b6001600160a01b0316146107c55760405162461bcd60e51b81526004016105c090611a42565b6107ce82610f26565b61070982826001610f50565b6000306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461087a5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016105c0565b50600080516020611be683398151915290565b6097546001600160a01b031633146108b75760405162461bcd60e51b81526004016105c090611a8e565b6108c160006110cf565b565b606060ca805461046d90611b69565b6097546001600160a01b031633146108fc5760405162461bcd60e51b81526004016105c090611a8e565b6107098282611121565b6000610913338484610d77565b50600192915050565b600063ffffffff333b16156109735760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e74726163740000000000000060448201526064016105c0565b34156109815761098161041d565b61098e3386868686611263565b336001600160a01b03167f282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d868686866040516109cd94939291906119c7565b60405180910390a2506001949350505050565b600054610100900460ff166109fb5760005460ff16156109ff565b303b155b610a625760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016105c0565b600054610100900460ff16158015610a84576000805461ffff19166101011790555b8451610a979060c9906020880190611617565b508351610aab9060ca906020870190611617565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610ade61135d565b610ae661138c565b8015610af8576000805461ff00191690555b5050505050565b6097546001600160a01b03163314610b295760405162461bcd60e51b81526004016105c090611a8e565b6001600160a01b038116610b8e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016105c0565b6106d2816110cf565b610ba13382611121565b6040516001600160a01b0383169082156108fc029083906000818181858888f19350505050158015610bd7573d6000803e3d6000fd5b506040518181526001600160a01b0383169033907f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb906020015b60405180910390a35050565b6001600160a01b038216610c735760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016105c0565b8060cc6000828254610c859190611b0e565b90915550506001600160a01b038216600090815260cd602052604081208054839290610cb2908490611b0e565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610c11565b6001600160a01b038316610d4b5760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016105c0565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610dcd5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016105c0565b6001600160a01b038216610e235760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016105c0565b6001600160a01b038316600090815260cd602052604090205481811015610e8c5760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016105c0565b610e968282611b26565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610ecc908490611b0e565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610f1891815260200190565b60405180910390a350505050565b6097546001600160a01b031633146106d25760405162461bcd60e51b81526004016105c090611a8e565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610f8857610f83836113b3565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610fc157600080fd5b505afa925050508015610ff1575060408051601f3d908101601f19168201909252610fee91810190611877565b60015b6110545760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016105c0565b600080516020611be683398151915281146110c35760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016105c0565b50610f8383838361144f565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b0382166111775760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016105c0565b6001600160a01b038216600090815260cd6020526040902054818110156111e05760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016105c0565b6111ea8282611b26565b6001600160a01b038416600090815260cd602052604081209190915560cc8054849290611218908490611b26565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b0385166112b95760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016105c0565b60008451116112fe5760405162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a59081c9958da5c1a595b9d607a1b60448201526064016105c0565b8061133c5760405162461bcd60e51b815260206004820152600e60248201526d1a5b9d985b1a59081d185c99d95d60921b60448201526064016105c0565b60cf54610af89086906001600160a01b03166113588587611b0e565b610d77565b600054610100900460ff166113845760405162461bcd60e51b81526004016105c090611ac3565b6108c161147a565b600054610100900460ff166108c15760405162461bcd60e51b81526004016105c090611ac3565b6001600160a01b0381163b6114205760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016105c0565b600080516020611be683398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611458836114aa565b6000825111806114655750805b15610f835761147483836114ea565b50505050565b600054610100900460ff166114a15760405162461bcd60e51b81526004016105c090611ac3565b6108c1336110cf565b6114b3816113b3565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606001600160a01b0383163b6115525760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b60648201526084016105c0565b600080846001600160a01b03168460405161156d9190611998565b600060405180830381855af49150503d80600081146115a8576040519150601f19603f3d011682016040523d82523d6000602084013e6115ad565b606091505b50915091506115d58282604051806060016040528060278152602001611c06602791396115de565b95945050505050565b606083156115ed5750816105ee565b8251156115fd5782518084602001fd5b8160405162461bcd60e51b81526004016105c091906119b4565b82805461162390611b69565b90600052602060002090601f016020900481019282611645576000855561168b565b82601f1061165e57805160ff191683800117855561168b565b8280016001018555821561168b579182015b8281111561168b578251825591602001919060010190611670565b5061169792915061169b565b5090565b5b80821115611697576000815560010161169c565b600067ffffffffffffffff808411156116cb576116cb611bba565b604051601f8501601f19908116603f011681019082821181831017156116f3576116f3611bba565b8160405280935085815286868601111561170c57600080fd5b858560208301376000602087830101525050509392505050565b600082601f830112611736578081fd5b6105ee838335602085016116b0565b600060208284031215611756578081fd5b81356105ee81611bd0565b60008060408385031215611773578081fd5b823561177e81611bd0565b946020939093013593505050565b6000806040838503121561179e578182fd5b82356117a981611bd0565b915060208301356117b981611bd0565b809150509250929050565b6000806000606084860312156117d8578081fd5b83356117e381611bd0565b925060208401356117f381611bd0565b929592945050506040919091013590565b60008060408385031215611816578182fd5b823561182181611bd0565b9150602083013567ffffffffffffffff81111561183c578182fd5b8301601f8101851361184c578182fd5b61185b858235602084016116b0565b9150509250929050565b60008060408385031215611773578182fd5b600060208284031215611888578081fd5b5051919050565b600080600080608085870312156118a4578081fd5b843567ffffffffffffffff808211156118bb578283fd5b6118c788838901611726565b955060208701359150808211156118dc578283fd5b506118e987828801611726565b935050604085013560ff811681146118ff578182fd5b9150606085013561190f81611bd0565b939692955090935050565b6000806000806080858703121561192f578384fd5b843567ffffffffffffffff811115611945578485fd5b61195187828801611726565b97602087013597506040870135966060013595509350505050565b60008151808452611984816020860160208601611b3d565b601f01601f19169290920160200192915050565b600082516119aa818460208701611b3d565b9190910192915050565b6020815260006105ee602083018461196c565b6080815260006119da608083018761196c565b6020830195909552506040810192909252606090910152919050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008219821115611b2157611b21611ba4565b500190565b600082821015611b3857611b38611ba4565b500390565b60005b83811015611b58578181015183820152602001611b40565b838111156114745750506000910152565b600181811c90821680611b7d57607f821691505b60208210811415611b9e57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b03811681146106d257600080fdfe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212203c27e83d0b8c9b8527404ed6d99dbdd9ee3310b5e1ae03b9506a430d642661b564736f6c63430008040033"

// DeployWFX deploys a new Ethereum contract, binding an instance of WFX to it.
func DeployWFX(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WFX, error) {
	parsed, err := abi.JSON(strings.NewReader(WFXABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(WFXBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WFX{WFXCaller: WFXCaller{contract: contract}, WFXTransactor: WFXTransactor{contract: contract}, WFXFilterer: WFXFilterer{contract: contract}}, nil
}

// WFX is an auto generated Go binding around an Ethereum contract.
type WFX struct {
	WFXCaller     // Read-only binding to the contract
	WFXTransactor // Write-only binding to the contract
	WFXFilterer   // Log filterer for contract events
}

// WFXCaller is an auto generated read-only Go binding around an Ethereum contract.
type WFXCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WFXTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WFXTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WFXFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WFXFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WFXSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WFXSession struct {
	Contract     *WFX              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WFXCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WFXCallerSession struct {
	Contract *WFXCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// WFXTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WFXTransactorSession struct {
	Contract     *WFXTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WFXRaw is an auto generated low-level Go binding around an Ethereum contract.
type WFXRaw struct {
	Contract *WFX // Generic contract binding to access the raw methods on
}

// WFXCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WFXCallerRaw struct {
	Contract *WFXCaller // Generic read-only contract binding to access the raw methods on
}

// WFXTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WFXTransactorRaw struct {
	Contract *WFXTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWFX creates a new instance of WFX, bound to a specific deployed contract.
func NewWFX(address common.Address, backend bind.ContractBackend) (*WFX, error) {
	contract, err := bindWFX(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WFX{WFXCaller: WFXCaller{contract: contract}, WFXTransactor: WFXTransactor{contract: contract}, WFXFilterer: WFXFilterer{contract: contract}}, nil
}

// NewWFXCaller creates a new read-only instance of WFX, bound to a specific deployed contract.
func NewWFXCaller(address common.Address, caller bind.ContractCaller) (*WFXCaller, error) {
	contract, err := bindWFX(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WFXCaller{contract: contract}, nil
}

// NewWFXTransactor creates a new write-only instance of WFX, bound to a specific deployed contract.
func NewWFXTransactor(address common.Address, transactor bind.ContractTransactor) (*WFXTransactor, error) {
	contract, err := bindWFX(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WFXTransactor{contract: contract}, nil
}

// NewWFXFilterer creates a new log filterer instance of WFX, bound to a specific deployed contract.
func NewWFXFilterer(address common.Address, filterer bind.ContractFilterer) (*WFXFilterer, error) {
	contract, err := bindWFX(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WFXFilterer{contract: contract}, nil
}

// bindWFX binds a generic wrapper to an already deployed contract.
func bindWFX(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WFXABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WFX *WFXRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WFX.Contract.WFXCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WFX *WFXRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WFX.Contract.WFXTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WFX *WFXRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WFX.Contract.WFXTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WFX *WFXCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WFX.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WFX *WFXTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WFX.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WFX *WFXTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WFX.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WFX *WFXCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WFX *WFXSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WFX.Contract.Allowance(&_WFX.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WFX *WFXCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WFX.Contract.Allowance(&_WFX.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WFX *WFXCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WFX *WFXSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WFX.Contract.BalanceOf(&_WFX.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WFX *WFXCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WFX.Contract.BalanceOf(&_WFX.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WFX *WFXCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WFX *WFXSession) Decimals() (uint8, error) {
	return _WFX.Contract.Decimals(&_WFX.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WFX *WFXCallerSession) Decimals() (uint8, error) {
	return _WFX.Contract.Decimals(&_WFX.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_WFX *WFXCaller) Module(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "module")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_WFX *WFXSession) Module() (common.Address, error) {
	return _WFX.Contract.Module(&_WFX.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_WFX *WFXCallerSession) Module() (common.Address, error) {
	return _WFX.Contract.Module(&_WFX.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WFX *WFXCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WFX *WFXSession) Name() (string, error) {
	return _WFX.Contract.Name(&_WFX.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WFX *WFXCallerSession) Name() (string, error) {
	return _WFX.Contract.Name(&_WFX.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WFX *WFXCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WFX *WFXSession) Owner() (common.Address, error) {
	return _WFX.Contract.Owner(&_WFX.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WFX *WFXCallerSession) Owner() (common.Address, error) {
	return _WFX.Contract.Owner(&_WFX.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_WFX *WFXCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_WFX *WFXSession) ProxiableUUID() ([32]byte, error) {
	return _WFX.Contract.ProxiableUUID(&_WFX.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_WFX *WFXCallerSession) ProxiableUUID() ([32]byte, error) {
	return _WFX.Contract.ProxiableUUID(&_WFX.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WFX *WFXCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WFX *WFXSession) Symbol() (string, error) {
	return _WFX.Contract.Symbol(&_WFX.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WFX *WFXCallerSession) Symbol() (string, error) {
	return _WFX.Contract.Symbol(&_WFX.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WFX *WFXCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WFX.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WFX *WFXSession) TotalSupply() (*big.Int, error) {
	return _WFX.Contract.TotalSupply(&_WFX.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WFX *WFXCallerSession) TotalSupply() (*big.Int, error) {
	return _WFX.Contract.TotalSupply(&_WFX.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WFX *WFXTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WFX *WFXSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Approve(&_WFX.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WFX *WFXTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Approve(&_WFX.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WFX *WFXTransactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "burn", account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WFX *WFXSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Burn(&_WFX.TransactOpts, account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WFX *WFXTransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Burn(&_WFX.TransactOpts, account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WFX *WFXTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WFX *WFXSession) Deposit() (*types.Transaction, error) {
	return _WFX.Contract.Deposit(&_WFX.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WFX *WFXTransactorSession) Deposit() (*types.Transaction, error) {
	return _WFX.Contract.Deposit(&_WFX.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_WFX *WFXTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "initialize", name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_WFX *WFXSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _WFX.Contract.Initialize(&_WFX.TransactOpts, name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_WFX *WFXTransactorSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _WFX.Contract.Initialize(&_WFX.TransactOpts, name_, symbol_, decimals_, module_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_WFX *WFXTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_WFX *WFXSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Mint(&_WFX.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_WFX *WFXTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Mint(&_WFX.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WFX *WFXTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WFX *WFXSession) RenounceOwnership() (*types.Transaction, error) {
	return _WFX.Contract.RenounceOwnership(&_WFX.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WFX *WFXTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WFX.Contract.RenounceOwnership(&_WFX.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_WFX *WFXTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_WFX *WFXSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Transfer(&_WFX.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_WFX *WFXTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Transfer(&_WFX.TransactOpts, recipient, amount)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) payable returns(bool)
func (_WFX *WFXTransactor) TransferCrossChain(opts *bind.TransactOpts, recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "transferCrossChain", recipient, amount, fee, target)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) payable returns(bool)
func (_WFX *WFXSession) TransferCrossChain(recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _WFX.Contract.TransferCrossChain(&_WFX.TransactOpts, recipient, amount, fee, target)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) payable returns(bool)
func (_WFX *WFXTransactorSession) TransferCrossChain(recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _WFX.Contract.TransferCrossChain(&_WFX.TransactOpts, recipient, amount, fee, target)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_WFX *WFXTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_WFX *WFXSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.TransferFrom(&_WFX.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_WFX *WFXTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.TransferFrom(&_WFX.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WFX *WFXTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WFX *WFXSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WFX.Contract.TransferOwnership(&_WFX.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WFX *WFXTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WFX.Contract.TransferOwnership(&_WFX.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_WFX *WFXTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_WFX *WFXSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _WFX.Contract.UpgradeTo(&_WFX.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_WFX *WFXTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _WFX.Contract.UpgradeTo(&_WFX.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_WFX *WFXTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_WFX *WFXSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _WFX.Contract.UpgradeToAndCall(&_WFX.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_WFX *WFXTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _WFX.Contract.UpgradeToAndCall(&_WFX.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address to, uint256 value) returns()
func (_WFX *WFXTransactor) Withdraw(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WFX.contract.Transact(opts, "withdraw", to, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address to, uint256 value) returns()
func (_WFX *WFXSession) Withdraw(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Withdraw(&_WFX.TransactOpts, to, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address to, uint256 value) returns()
func (_WFX *WFXTransactorSession) Withdraw(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WFX.Contract.Withdraw(&_WFX.TransactOpts, to, value)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WFX *WFXTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _WFX.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WFX *WFXSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WFX.Contract.Fallback(&_WFX.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WFX *WFXTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WFX.Contract.Fallback(&_WFX.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_WFX *WFXTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WFX.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_WFX *WFXSession) Receive() (*types.Transaction, error) {
	return _WFX.Contract.Receive(&_WFX.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_WFX *WFXTransactorSession) Receive() (*types.Transaction, error) {
	return _WFX.Contract.Receive(&_WFX.TransactOpts)
}

// WFXAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the WFX contract.
type WFXAdminChangedIterator struct {
	Event *WFXAdminChanged // Event containing the contract specifics and raw log

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
func (it *WFXAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXAdminChanged)
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
		it.Event = new(WFXAdminChanged)
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
func (it *WFXAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXAdminChanged represents a AdminChanged event raised by the WFX contract.
type WFXAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_WFX *WFXFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*WFXAdminChangedIterator, error) {

	logs, sub, err := _WFX.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &WFXAdminChangedIterator{contract: _WFX.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_WFX *WFXFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *WFXAdminChanged) (event.Subscription, error) {

	logs, sub, err := _WFX.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXAdminChanged)
				if err := _WFX.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_WFX *WFXFilterer) ParseAdminChanged(log types.Log) (*WFXAdminChanged, error) {
	event := new(WFXAdminChanged)
	if err := _WFX.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WFX contract.
type WFXApprovalIterator struct {
	Event *WFXApproval // Event containing the contract specifics and raw log

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
func (it *WFXApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXApproval)
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
		it.Event = new(WFXApproval)
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
func (it *WFXApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXApproval represents a Approval event raised by the WFX contract.
type WFXApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WFX *WFXFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WFXApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WFXApprovalIterator{contract: _WFX.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WFX *WFXFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WFXApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXApproval)
				if err := _WFX.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WFX *WFXFilterer) ParseApproval(log types.Log) (*WFXApproval, error) {
	event := new(WFXApproval)
	if err := _WFX.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the WFX contract.
type WFXBeaconUpgradedIterator struct {
	Event *WFXBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *WFXBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXBeaconUpgraded)
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
		it.Event = new(WFXBeaconUpgraded)
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
func (it *WFXBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXBeaconUpgraded represents a BeaconUpgraded event raised by the WFX contract.
type WFXBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_WFX *WFXFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*WFXBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &WFXBeaconUpgradedIterator{contract: _WFX.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_WFX *WFXFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *WFXBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXBeaconUpgraded)
				if err := _WFX.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_WFX *WFXFilterer) ParseBeaconUpgraded(log types.Log) (*WFXBeaconUpgraded, error) {
	event := new(WFXBeaconUpgraded)
	if err := _WFX.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the WFX contract.
type WFXDepositIterator struct {
	Event *WFXDeposit // Event containing the contract specifics and raw log

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
func (it *WFXDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXDeposit)
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
		it.Event = new(WFXDeposit)
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
func (it *WFXDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXDeposit represents a Deposit event raised by the WFX contract.
type WFXDeposit struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed from, uint256 value)
func (_WFX *WFXFilterer) FilterDeposit(opts *bind.FilterOpts, from []common.Address) (*WFXDepositIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "Deposit", fromRule)
	if err != nil {
		return nil, err
	}
	return &WFXDepositIterator{contract: _WFX.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed from, uint256 value)
func (_WFX *WFXFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WFXDeposit, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "Deposit", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXDeposit)
				if err := _WFX.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed from, uint256 value)
func (_WFX *WFXFilterer) ParseDeposit(log types.Log) (*WFXDeposit, error) {
	event := new(WFXDeposit)
	if err := _WFX.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WFX contract.
type WFXOwnershipTransferredIterator struct {
	Event *WFXOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WFXOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXOwnershipTransferred)
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
		it.Event = new(WFXOwnershipTransferred)
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
func (it *WFXOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXOwnershipTransferred represents a OwnershipTransferred event raised by the WFX contract.
type WFXOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WFX *WFXFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WFXOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WFXOwnershipTransferredIterator{contract: _WFX.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WFX *WFXFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WFXOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXOwnershipTransferred)
				if err := _WFX.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WFX *WFXFilterer) ParseOwnershipTransferred(log types.Log) (*WFXOwnershipTransferred, error) {
	event := new(WFXOwnershipTransferred)
	if err := _WFX.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WFX contract.
type WFXTransferIterator struct {
	Event *WFXTransfer // Event containing the contract specifics and raw log

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
func (it *WFXTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXTransfer)
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
		it.Event = new(WFXTransfer)
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
func (it *WFXTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXTransfer represents a Transfer event raised by the WFX contract.
type WFXTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WFX *WFXFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WFXTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WFXTransferIterator{contract: _WFX.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WFX *WFXFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WFXTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXTransfer)
				if err := _WFX.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WFX *WFXFilterer) ParseTransfer(log types.Log) (*WFXTransfer, error) {
	event := new(WFXTransfer)
	if err := _WFX.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXTransferCrossChainIterator is returned from FilterTransferCrossChain and is used to iterate over the raw logs and unpacked data for TransferCrossChain events raised by the WFX contract.
type WFXTransferCrossChainIterator struct {
	Event *WFXTransferCrossChain // Event containing the contract specifics and raw log

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
func (it *WFXTransferCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXTransferCrossChain)
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
		it.Event = new(WFXTransferCrossChain)
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
func (it *WFXTransferCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXTransferCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXTransferCrossChain represents a TransferCrossChain event raised by the WFX contract.
type WFXTransferCrossChain struct {
	From      common.Address
	Recipient string
	Amount    *big.Int
	Fee       *big.Int
	Target    [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferCrossChain is a free log retrieval operation binding the contract event 0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d.
//
// Solidity: event TransferCrossChain(address indexed from, string recipient, uint256 amount, uint256 fee, bytes32 target)
func (_WFX *WFXFilterer) FilterTransferCrossChain(opts *bind.FilterOpts, from []common.Address) (*WFXTransferCrossChainIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "TransferCrossChain", fromRule)
	if err != nil {
		return nil, err
	}
	return &WFXTransferCrossChainIterator{contract: _WFX.contract, event: "TransferCrossChain", logs: logs, sub: sub}, nil
}

// WatchTransferCrossChain is a free log subscription operation binding the contract event 0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d.
//
// Solidity: event TransferCrossChain(address indexed from, string recipient, uint256 amount, uint256 fee, bytes32 target)
func (_WFX *WFXFilterer) WatchTransferCrossChain(opts *bind.WatchOpts, sink chan<- *WFXTransferCrossChain, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "TransferCrossChain", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXTransferCrossChain)
				if err := _WFX.contract.UnpackLog(event, "TransferCrossChain", log); err != nil {
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

// ParseTransferCrossChain is a log parse operation binding the contract event 0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d.
//
// Solidity: event TransferCrossChain(address indexed from, string recipient, uint256 amount, uint256 fee, bytes32 target)
func (_WFX *WFXFilterer) ParseTransferCrossChain(log types.Log) (*WFXTransferCrossChain, error) {
	event := new(WFXTransferCrossChain)
	if err := _WFX.contract.UnpackLog(event, "TransferCrossChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the WFX contract.
type WFXUpgradedIterator struct {
	Event *WFXUpgraded // Event containing the contract specifics and raw log

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
func (it *WFXUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXUpgraded)
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
		it.Event = new(WFXUpgraded)
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
func (it *WFXUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXUpgraded represents a Upgraded event raised by the WFX contract.
type WFXUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_WFX *WFXFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*WFXUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &WFXUpgradedIterator{contract: _WFX.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_WFX *WFXFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *WFXUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXUpgraded)
				if err := _WFX.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_WFX *WFXFilterer) ParseUpgraded(log types.Log) (*WFXUpgraded, error) {
	event := new(WFXUpgraded)
	if err := _WFX.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WFXWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the WFX contract.
type WFXWithdrawIterator struct {
	Event *WFXWithdraw // Event containing the contract specifics and raw log

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
func (it *WFXWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WFXWithdraw)
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
		it.Event = new(WFXWithdraw)
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
func (it *WFXWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WFXWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WFXWithdraw represents a Withdraw event raised by the WFX contract.
type WFXWithdraw struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: event Withdraw(address indexed from, address indexed to, uint256 value)
func (_WFX *WFXFilterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WFXWithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WFX.contract.FilterLogs(opts, "Withdraw", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WFXWithdrawIterator{contract: _WFX.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: event Withdraw(address indexed from, address indexed to, uint256 value)
func (_WFX *WFXFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *WFXWithdraw, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WFX.contract.WatchLogs(opts, "Withdraw", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WFXWithdraw)
				if err := _WFX.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: event Withdraw(address indexed from, address indexed to, uint256 value)
func (_WFX *WFXFilterer) ParseWithdraw(log types.Log) (*WFXWithdraw, error) {
	event := new(WFXWithdraw)
	if err := _WFX.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
