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

// LPTokenABI is the input ABI used to generate the binding from.
const LPTokenABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"ApprovalLock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferLock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Unlock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approveLock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"lockAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"lockAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"}],\"name\":\"lockBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"selfDestruct\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferLock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// LPTokenBin is the compiled bytecode used for deploying new contracts.
var LPTokenBin = "0x60a06040523060601b60805234801561001757600080fd5b5060805160601c6125c7610052600039600081816108a3015281816108e3015281816109d101528181610a110152610aa401526125c76000f3fe6080604052600436106101815760003560e01c806370a08231116100d15780639cb8a26a1161008a578063d50edf8011610064578063d50edf8014610489578063dd62ed3e146104a9578063e2095ab4146104ef578063f2fde38b1461053557600080fd5b80639cb8a26a146104345780639dc29fac14610449578063a9059cbb1461046957600080fd5b806370a0823114610346578063715018a61461037c5780637eee288d1461039157806382c5c93c146103b15780638da5cb5b146103f757806395d89b411461041f57600080fd5b8063282d3fdf1161013e57806340c10f191161011857806340c10f19146102de5780634ef09d79146102fe5780634f1ef2861461031e57806352d1902d1461033157600080fd5b8063282d3fdf1461027c578063313ce5671461029c5780633659cfe6146102be57600080fd5b806306fdde0314610186578063095ea7b3146101b157806310e776ed146101e15780631624f6c61461022557806318160ddd1461024757806323b872dd1461025c575b600080fd5b34801561019257600080fd5b5061019b610555565b6040516101a8919061233d565b60405180910390f35b3480156101bd57600080fd5b506101d16101cc366004612266565b6105e7565b60405190151581526020016101a8565b3480156101ed57600080fd5b506102176101fc366004612180565b6001600160a01b0316600090815260cf602052604090205490565b6040519081526020016101a8565b34801561023157600080fd5b506102456102403660046122a7565b61063e565b005b34801561025357600080fd5b5060cc54610217565b34801561026857600080fd5b506101d16102773660046121cc565b610746565b34801561028857600080fd5b506101d1610297366004612266565b6107f0565b3480156102a857600080fd5b5060cb5460405160ff90911681526020016101a8565b3480156102ca57600080fd5b506102456102d9366004612180565b610898565b3480156102ea57600080fd5b506101d16102f9366004612266565b610978565b34801561030a57600080fd5b506101d16103193660046121cc565b6109b8565b61024561032c366004612207565b6109c6565b34801561033d57600080fd5b50610217610a97565b34801561035257600080fd5b50610217610361366004612180565b6001600160a01b0316600090815260cd602052604090205490565b34801561038857600080fd5b50610245610b4a565b34801561039d57600080fd5b506101d16103ac366004612266565b610b80565b3480156103bd57600080fd5b506102176103cc36600461219a565b6001600160a01b03918216600090815260d16020908152604080832093909416825291909152205490565b34801561040357600080fd5b506097546040516001600160a01b0390911681526020016101a8565b34801561042b57600080fd5b5061019b610b8d565b34801561044057600080fd5b50610245610b9c565b34801561045557600080fd5b506101d1610464366004612266565b610be8565b34801561047557600080fd5b506101d1610484366004612266565b610c1f565b34801561049557600080fd5b506101d16104a4366004612266565b610c2c565b3480156104b557600080fd5b506102176104c436600461219a565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156104fb57600080fd5b5061021761050a36600461219a565b6001600160a01b03918216600090815260d06020908152604080832093909416825291909152205490565b34801561054157600080fd5b50610245610550366004612180565b610c76565b606060c98054610564906124e3565b80601f0160208091040260200160405190810160405280929190818152602001828054610590906124e3565b80156105dd5780601f106105b2576101008083540402835291602001916105dd565b820191906000526020600020905b8154815290600101906020018083116105c057829003601f168201915b5050505050905090565b60006105f4338484610d0e565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925906020015b60405180910390a350600192915050565b600054610100900460ff166106595760005460ff161561065d565b303b155b6106c55760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b600054610100900460ff161580156106e7576000805461ffff19166101011790555b83516106fa9060c9906020870190612036565b50825161070e9060ca906020860190612036565b5060cb805460ff191660ff8416179055610726610d90565b61072e610dbf565b8015610740576000805461ff00191690555b50505050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156107c45760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084016106bc565b6107d885336107d386856124a0565b610d0e565b6107e3858585610de6565b60019150505b9392505050565b6001600160a01b038216600090815260d1602090815260408083203384529091528120548281101561086f5760405162461bcd60e51b815260206004820152602260248201527f6c6f636b20616d6f756e742065786365656473206c6f636b20616c6c6f77616e604482015261636560f01b60648201526084016106bc565b610883843361087e86856124a0565b61101a565b61088e3385856110a7565b5060019392505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156108e15760405162461bcd60e51b81526004016106bc90612370565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661092a60008051602061254b833981519152546001600160a01b031690565b6001600160a01b0316146109505760405162461bcd60e51b81526004016106bc906123bc565b6109598161128c565b60408051600080825260208201909252610975918391906112b6565b50565b6097546000906001600160a01b031633146109a55760405162461bcd60e51b81526004016106bc90612408565b6109af8383611435565b50600192915050565b600061088e33858585611557565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161415610a0f5760405162461bcd60e51b81526004016106bc90612370565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610a5860008051602061254b833981519152546001600160a01b031690565b6001600160a01b031614610a7e5760405162461bcd60e51b81526004016106bc906123bc565b610a878261128c565b610a93828260016112b6565b5050565b6000306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610b375760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016106bc565b5060008051602061254b83398151915290565b6097546001600160a01b03163314610b745760405162461bcd60e51b81526004016106bc90612408565b610b7e60006118be565b565b60006109af338484611910565b606060ca8054610564906124e3565b6097546001600160a01b03163314610bc65760405162461bcd60e51b81526004016106bc90612408565b6000610bda6097546001600160a01b031690565b9050806001600160a01b0316ff5b6097546000906001600160a01b03163314610c155760405162461bcd60e51b81526004016106bc90612408565b6109af8383611b00565b60006109af338484610de6565b6000610c3933848461101a565b6040518281526001600160a01b0384169033907f6a8ad50d47d7e8cfb288b0f40af42ea12bdc6f2ddcfe403f95854afa116fc8ee9060200161062d565b6097546001600160a01b03163314610ca05760405162461bcd60e51b81526004016106bc90612408565b6001600160a01b038116610d055760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016106bc565b610975816118be565b6001600160a01b038316610d645760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016106bc565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b600054610100900460ff16610db75760405162461bcd60e51b81526004016106bc9061243d565b610b7e611d06565b600054610100900460ff16610b7e5760405162461bcd60e51b81526004016106bc9061243d565b6001600160a01b038316610e3c5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016106bc565b6001600160a01b038216610e925760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016106bc565b6001600160a01b038316600090815260cd602052604090205481811015610efb5760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016106bc565b6001600160a01b038416600090815260cf602052604090205482610f1f82846124a0565b1015610f7f5760405162461bcd60e51b815260206004820152602960248201527f7472616e7366657220616d6f756e7420657863656564732072656d61696e696e604482015268672062616c616e636560b81b60648201526084016106bc565b610f8983836124a0565b6001600160a01b03808716600090815260cd60205260408082209390935590861681529081208054859290610fbf908490612488565b92505081905550836001600160a01b0316856001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8560405161100b91815260200190565b60405180910390a35050505050565b6001600160a01b03831661107b5760405162461bcd60e51b815260206004820152602260248201527f617070726f7665206c6f636b2066726f6d20746865207a65726f206164647265604482015261737360f01b60648201526084016106bc565b6001600160a01b03928316600090815260d1602090815260408083209490951682529290925291902055565b6001600160a01b0383166110fd5760405162461bcd60e51b815260206004820152601d60248201527f6c6f636b207370656e64657220746865207a65726f206164647265737300000060448201526064016106bc565b6001600160a01b0382166111535760405162461bcd60e51b815260206004820152601860248201527f6c6f636b20746f20746865207a65726f2061646472657373000000000000000060448201526064016106bc565b6001600160a01b038216600090815260cd6020526040902054818110156111bc5760405162461bcd60e51b815260206004820152601b60248201527f6c6f636b20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016106bc565b6001600160a01b038316600090815260cf6020526040902054826111e082846124a0565b101561123c5760405162461bcd60e51b815260206004820152602560248201527f6c6f636b20616d6f756e7420657863656564732072656d61696e696e672062616044820152646c616e636560d81b60648201526084016106bc565b611247848685611d36565b836001600160a01b0316856001600160a01b03167fec36c0364d931187a76cf66d7eee08fad0ec2e8b7458a8d8b26b36769d4d13f38560405161100b91815260200190565b6097546001600160a01b031633146109755760405162461bcd60e51b81526004016106bc90612408565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156112ee576112e983611da4565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561132757600080fd5b505afa925050508015611357575060408051601f3d908101601f191682019092526113549181019061228f565b60015b6113ba5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016106bc565b60008051602061254b83398151915281146114295760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016106bc565b506112e9838383611e40565b6001600160a01b03821661148b5760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016106bc565b8060cc600082825461149d9190612488565b90915550506001600160a01b038216600090815260cd6020526040812080548392906114ca908490612488565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3816001600160a01b03167f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d41213968858260405161154b91815260200190565b60405180910390a25050565b6001600160a01b0384166115bc5760405162461bcd60e51b815260206004820152602660248201527f7472616e73666572206c6f636b207370656e64657220746865207a65726f206160448201526564647265737360d01b60648201526084016106bc565b6001600160a01b0383166116205760405162461bcd60e51b815260206004820152602560248201527f7472616e73666572206c6f636b206c6f636b657220746865207a65726f206164604482015264647265737360d81b60648201526084016106bc565b6001600160a01b0382166116805760405162461bcd60e51b815260206004820152602160248201527f7472616e73666572206c6f636b20746f20746865207a65726f206164647265736044820152607360f81b60648201526084016106bc565b6001600160a01b038316600090815260cd6020526040902054818110156116f55760405162461bcd60e51b8152602060048201526024808201527f7472616e73666572206c6f636b20616d6f756e7420657863656564732062616c604482015263616e636560e01b60648201526084016106bc565b6001600160a01b038416600090815260cf6020526040902054828110156117705760405162461bcd60e51b815260206004820152602960248201527f7472616e73666572206c6f636b20616d6f756e742065786365656473206c6f636044820152686b2062616c616e636560b81b60648201526084016106bc565b6001600160a01b03808616600090815260d060209081526040808320938a1683529290522054838110156117f75760405162461bcd60e51b815260206004820152602860248201527f7472616e73666572206c6f636b20616d6f756e742065786365656473206c6f636044820152671ac8185b5bdd5b9d60c21b60648201526084016106bc565b611802868886611e65565b6001600160a01b038616600090815260cd60205260408120805486929061182a9084906124a0565b90915550506001600160a01b038516600090815260cd602052604081208054869290611857908490612488565b92505081905550846001600160a01b0316866001600160a01b0316886001600160a01b03167f164e3db520d3c5f437914d09eef8478390d49d797a1b1c9aa5a0eed966183343876040516118ad91815260200190565b60405180910390a450505050505050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b0383166119665760405162461bcd60e51b815260206004820152601f60248201527f756e6c6f636b207370656e64657220746865207a65726f20616464726573730060448201526064016106bc565b6001600160a01b0382166119bc5760405162461bcd60e51b815260206004820152601a60248201527f756e6c6f636b20746f20746865207a65726f206164647265737300000000000060448201526064016106bc565b6001600160a01b038216600090815260cf602052604090205481811015611a305760405162461bcd60e51b815260206004820152602260248201527f756e6c6f636b20616d6f756e742065786365656473206c6f636b2062616c616e604482015261636560f01b60648201526084016106bc565b6001600160a01b03808416600090815260d0602090815260408083209388168352929052205482811015611ab05760405162461bcd60e51b815260206004820152602160248201527f756e6c6f636b20616d6f756e742065786365656473206c6f636b20616d6f756e6044820152601d60fa1b60648201526084016106bc565b611abb848685611e65565b836001600160a01b0316856001600160a01b03167fc1c90b8e0705b212262c0dbd7580efe1862c2f185bf96899226f7596beb2db098560405161100b91815260200190565b6001600160a01b038216611b565760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016106bc565b6001600160a01b038216600090815260cd602052604090205481811015611bbf5760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016106bc565b6001600160a01b038316600090815260cf602052604090205482611be382846124a0565b1015611c3f5760405162461bcd60e51b815260206004820152602560248201527f6275726e20616d6f756e7420657863656564732072656d61696e696e672062616044820152646c616e636560d81b60648201526084016106bc565b611c4983836124a0565b6001600160a01b038516600090815260cd602052604081209190915560cc8054859290611c779084906124a0565b90915550506040518381526000906001600160a01b038616907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3836001600160a01b03167fcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca584604051611cf891815260200190565b60405180910390a250505050565b600054610100900460ff16611d2d5760405162461bcd60e51b81526004016106bc9061243d565b610b7e336118be565b6001600160a01b038316600090815260cf602052604081208054839290611d5e908490612488565b90915550506001600160a01b03808416600090815260d06020908152604080832093861683529290529081208054839290611d9a908490612488565b9091555050505050565b6001600160a01b0381163b611e115760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016106bc565b60008051602061254b83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611e4983611ec9565b600082511180611e565750805b156112e9576107408383611f09565b6001600160a01b038316600090815260cf602052604081208054839290611e8d9084906124a0565b90915550506001600160a01b03808416600090815260d06020908152604080832093861683529290529081208054839290611d9a9084906124a0565b611ed281611da4565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606001600160a01b0383163b611f715760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b60648201526084016106bc565b600080846001600160a01b031684604051611f8c9190612321565b600060405180830381855af49150503d8060008114611fc7576040519150601f19603f3d011682016040523d82523d6000602084013e611fcc565b606091505b5091509150611ff4828260405180606001604052806027815260200161256b60279139611ffd565b95945050505050565b6060831561200c5750816107e9565b82511561201c5782518084602001fd5b8160405162461bcd60e51b81526004016106bc919061233d565b828054612042906124e3565b90600052602060002090601f01602090048101928261206457600085556120aa565b82601f1061207d57805160ff19168380011785556120aa565b828001600101855582156120aa579182015b828111156120aa57825182559160200191906001019061208f565b506120b69291506120ba565b5090565b5b808211156120b657600081556001016120bb565b600067ffffffffffffffff808411156120ea576120ea612534565b604051601f8501601f19908116603f0116810190828211818310171561211257612112612534565b8160405280935085815286868601111561212b57600080fd5b858560208301376000602087830101525050509392505050565b80356001600160a01b038116811461215c57600080fd5b919050565b600082601f830112612171578081fd5b6107e9838335602085016120cf565b600060208284031215612191578081fd5b6107e982612145565b600080604083850312156121ac578081fd5b6121b583612145565b91506121c360208401612145565b90509250929050565b6000806000606084860312156121e0578081fd5b6121e984612145565b92506121f760208501612145565b9150604084013590509250925092565b60008060408385031215612219578182fd5b61222283612145565b9150602083013567ffffffffffffffff81111561223d578182fd5b8301601f8101851361224d578182fd5b61225c858235602084016120cf565b9150509250929050565b60008060408385031215612278578182fd5b61228183612145565b946020939093013593505050565b6000602082840312156122a0578081fd5b5051919050565b6000806000606084860312156122bb578283fd5b833567ffffffffffffffff808211156122d2578485fd5b6122de87838801612161565b945060208601359150808211156122f3578384fd5b5061230086828701612161565b925050604084013560ff81168114612316578182fd5b809150509250925092565b600082516123338184602087016124b7565b9190910192915050565b602081526000825180602084015261235c8160408501602087016124b7565b601f01601f19169190910160400192915050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b6000821982111561249b5761249b61251e565b500190565b6000828210156124b2576124b261251e565b500390565b60005b838110156124d25781810151838201526020016124ba565b838111156107405750506000910152565b600181811c908216806124f757607f821691505b6020821081141561251857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fdfe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122047965ac112b4afd70bfa61dbd0cec8c41c2f4ed7f3e1a67ef2d103b9f71eb25964736f6c63430008040033"

// DeployLPToken deploys a new Ethereum contract, binding an instance of LPToken to it.
func DeployLPToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LPToken, error) {
	parsed, err := abi.JSON(strings.NewReader(LPTokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LPTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LPToken{LPTokenCaller: LPTokenCaller{contract: contract}, LPTokenTransactor: LPTokenTransactor{contract: contract}, LPTokenFilterer: LPTokenFilterer{contract: contract}}, nil
}

// LPToken is an auto generated Go binding around an Ethereum contract.
type LPToken struct {
	LPTokenCaller     // Read-only binding to the contract
	LPTokenTransactor // Write-only binding to the contract
	LPTokenFilterer   // Log filterer for contract events
}

// LPTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type LPTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LPTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LPTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LPTokenSession struct {
	Contract     *LPToken          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LPTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LPTokenCallerSession struct {
	Contract *LPTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// LPTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LPTokenTransactorSession struct {
	Contract     *LPTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LPTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type LPTokenRaw struct {
	Contract *LPToken // Generic contract binding to access the raw methods on
}

// LPTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LPTokenCallerRaw struct {
	Contract *LPTokenCaller // Generic read-only contract binding to access the raw methods on
}

// LPTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LPTokenTransactorRaw struct {
	Contract *LPTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLPToken creates a new instance of LPToken, bound to a specific deployed contract.
func NewLPToken(address common.Address, backend bind.ContractBackend) (*LPToken, error) {
	contract, err := bindLPToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LPToken{LPTokenCaller: LPTokenCaller{contract: contract}, LPTokenTransactor: LPTokenTransactor{contract: contract}, LPTokenFilterer: LPTokenFilterer{contract: contract}}, nil
}

// NewLPTokenCaller creates a new read-only instance of LPToken, bound to a specific deployed contract.
func NewLPTokenCaller(address common.Address, caller bind.ContractCaller) (*LPTokenCaller, error) {
	contract, err := bindLPToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LPTokenCaller{contract: contract}, nil
}

// NewLPTokenTransactor creates a new write-only instance of LPToken, bound to a specific deployed contract.
func NewLPTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*LPTokenTransactor, error) {
	contract, err := bindLPToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LPTokenTransactor{contract: contract}, nil
}

// NewLPTokenFilterer creates a new log filterer instance of LPToken, bound to a specific deployed contract.
func NewLPTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*LPTokenFilterer, error) {
	contract, err := bindLPToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LPTokenFilterer{contract: contract}, nil
}

// bindLPToken binds a generic wrapper to an already deployed contract.
func bindLPToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LPTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LPToken *LPTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LPToken.Contract.LPTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LPToken *LPTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPToken.Contract.LPTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LPToken *LPTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LPToken.Contract.LPTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LPToken *LPTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LPToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LPToken *LPTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LPToken *LPTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LPToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LPToken *LPTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LPToken *LPTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LPToken.Contract.Allowance(&_LPToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_LPToken *LPTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _LPToken.Contract.Allowance(&_LPToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LPToken *LPTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LPToken *LPTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LPToken.Contract.BalanceOf(&_LPToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_LPToken *LPTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _LPToken.Contract.BalanceOf(&_LPToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LPToken *LPTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LPToken *LPTokenSession) Decimals() (uint8, error) {
	return _LPToken.Contract.Decimals(&_LPToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_LPToken *LPTokenCallerSession) Decimals() (uint8, error) {
	return _LPToken.Contract.Decimals(&_LPToken.CallOpts)
}

// LockAllowance is a free data retrieval call binding the contract method 0x82c5c93c.
//
// Solidity: function lockAllowance(address locker, address spender) view returns(uint256)
func (_LPToken *LPTokenCaller) LockAllowance(opts *bind.CallOpts, locker common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "lockAllowance", locker, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockAllowance is a free data retrieval call binding the contract method 0x82c5c93c.
//
// Solidity: function lockAllowance(address locker, address spender) view returns(uint256)
func (_LPToken *LPTokenSession) LockAllowance(locker common.Address, spender common.Address) (*big.Int, error) {
	return _LPToken.Contract.LockAllowance(&_LPToken.CallOpts, locker, spender)
}

// LockAllowance is a free data retrieval call binding the contract method 0x82c5c93c.
//
// Solidity: function lockAllowance(address locker, address spender) view returns(uint256)
func (_LPToken *LPTokenCallerSession) LockAllowance(locker common.Address, spender common.Address) (*big.Int, error) {
	return _LPToken.Contract.LockAllowance(&_LPToken.CallOpts, locker, spender)
}

// LockAmount is a free data retrieval call binding the contract method 0xe2095ab4.
//
// Solidity: function lockAmount(address locker, address spender) view returns(uint256)
func (_LPToken *LPTokenCaller) LockAmount(opts *bind.CallOpts, locker common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "lockAmount", locker, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockAmount is a free data retrieval call binding the contract method 0xe2095ab4.
//
// Solidity: function lockAmount(address locker, address spender) view returns(uint256)
func (_LPToken *LPTokenSession) LockAmount(locker common.Address, spender common.Address) (*big.Int, error) {
	return _LPToken.Contract.LockAmount(&_LPToken.CallOpts, locker, spender)
}

// LockAmount is a free data retrieval call binding the contract method 0xe2095ab4.
//
// Solidity: function lockAmount(address locker, address spender) view returns(uint256)
func (_LPToken *LPTokenCallerSession) LockAmount(locker common.Address, spender common.Address) (*big.Int, error) {
	return _LPToken.Contract.LockAmount(&_LPToken.CallOpts, locker, spender)
}

// LockBalanceOf is a free data retrieval call binding the contract method 0x10e776ed.
//
// Solidity: function lockBalanceOf(address locker) view returns(uint256)
func (_LPToken *LPTokenCaller) LockBalanceOf(opts *bind.CallOpts, locker common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "lockBalanceOf", locker)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockBalanceOf is a free data retrieval call binding the contract method 0x10e776ed.
//
// Solidity: function lockBalanceOf(address locker) view returns(uint256)
func (_LPToken *LPTokenSession) LockBalanceOf(locker common.Address) (*big.Int, error) {
	return _LPToken.Contract.LockBalanceOf(&_LPToken.CallOpts, locker)
}

// LockBalanceOf is a free data retrieval call binding the contract method 0x10e776ed.
//
// Solidity: function lockBalanceOf(address locker) view returns(uint256)
func (_LPToken *LPTokenCallerSession) LockBalanceOf(locker common.Address) (*big.Int, error) {
	return _LPToken.Contract.LockBalanceOf(&_LPToken.CallOpts, locker)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LPToken *LPTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LPToken *LPTokenSession) Name() (string, error) {
	return _LPToken.Contract.Name(&_LPToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LPToken *LPTokenCallerSession) Name() (string, error) {
	return _LPToken.Contract.Name(&_LPToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LPToken *LPTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LPToken *LPTokenSession) Owner() (common.Address, error) {
	return _LPToken.Contract.Owner(&_LPToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LPToken *LPTokenCallerSession) Owner() (common.Address, error) {
	return _LPToken.Contract.Owner(&_LPToken.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LPToken *LPTokenCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LPToken *LPTokenSession) ProxiableUUID() ([32]byte, error) {
	return _LPToken.Contract.ProxiableUUID(&_LPToken.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LPToken *LPTokenCallerSession) ProxiableUUID() ([32]byte, error) {
	return _LPToken.Contract.ProxiableUUID(&_LPToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LPToken *LPTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LPToken *LPTokenSession) Symbol() (string, error) {
	return _LPToken.Contract.Symbol(&_LPToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LPToken *LPTokenCallerSession) Symbol() (string, error) {
	return _LPToken.Contract.Symbol(&_LPToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LPToken *LPTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LPToken *LPTokenSession) TotalSupply() (*big.Int, error) {
	return _LPToken.Contract.TotalSupply(&_LPToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LPToken *LPTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _LPToken.Contract.TotalSupply(&_LPToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Approve(&_LPToken.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Approve(&_LPToken.TransactOpts, spender, amount)
}

// ApproveLock is a paid mutator transaction binding the contract method 0xd50edf80.
//
// Solidity: function approveLock(address spender, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) ApproveLock(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "approveLock", spender, amount)
}

// ApproveLock is a paid mutator transaction binding the contract method 0xd50edf80.
//
// Solidity: function approveLock(address spender, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) ApproveLock(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.ApproveLock(&_LPToken.TransactOpts, spender, amount)
}

// ApproveLock is a paid mutator transaction binding the contract method 0xd50edf80.
//
// Solidity: function approveLock(address spender, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) ApproveLock(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.ApproveLock(&_LPToken.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "burn", account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Burn(&_LPToken.TransactOpts, account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Burn(&_LPToken.TransactOpts, account, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0x1624f6c6.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_) returns()
func (_LPToken *LPTokenTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "initialize", name_, symbol_, decimals_)
}

// Initialize is a paid mutator transaction binding the contract method 0x1624f6c6.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_) returns()
func (_LPToken *LPTokenSession) Initialize(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _LPToken.Contract.Initialize(&_LPToken.TransactOpts, name_, symbol_, decimals_)
}

// Initialize is a paid mutator transaction binding the contract method 0x1624f6c6.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_) returns()
func (_LPToken *LPTokenTransactorSession) Initialize(name_ string, symbol_ string, decimals_ uint8) (*types.Transaction, error) {
	return _LPToken.Contract.Initialize(&_LPToken.TransactOpts, name_, symbol_, decimals_)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address locker, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) Lock(opts *bind.TransactOpts, locker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "lock", locker, amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address locker, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) Lock(locker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Lock(&_LPToken.TransactOpts, locker, amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address locker, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) Lock(locker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Lock(&_LPToken.TransactOpts, locker, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Mint(&_LPToken.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Mint(&_LPToken.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LPToken *LPTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LPToken *LPTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _LPToken.Contract.RenounceOwnership(&_LPToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LPToken *LPTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _LPToken.Contract.RenounceOwnership(&_LPToken.TransactOpts)
}

// SelfDestruct is a paid mutator transaction binding the contract method 0x9cb8a26a.
//
// Solidity: function selfDestruct() returns()
func (_LPToken *LPTokenTransactor) SelfDestruct(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "selfDestruct")
}

// SelfDestruct is a paid mutator transaction binding the contract method 0x9cb8a26a.
//
// Solidity: function selfDestruct() returns()
func (_LPToken *LPTokenSession) SelfDestruct() (*types.Transaction, error) {
	return _LPToken.Contract.SelfDestruct(&_LPToken.TransactOpts)
}

// SelfDestruct is a paid mutator transaction binding the contract method 0x9cb8a26a.
//
// Solidity: function selfDestruct() returns()
func (_LPToken *LPTokenTransactorSession) SelfDestruct() (*types.Transaction, error) {
	return _LPToken.Contract.SelfDestruct(&_LPToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Transfer(&_LPToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Transfer(&_LPToken.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.TransferFrom(&_LPToken.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.TransferFrom(&_LPToken.TransactOpts, sender, recipient, amount)
}

// TransferLock is a paid mutator transaction binding the contract method 0x4ef09d79.
//
// Solidity: function transferLock(address locker, address to, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) TransferLock(opts *bind.TransactOpts, locker common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "transferLock", locker, to, amount)
}

// TransferLock is a paid mutator transaction binding the contract method 0x4ef09d79.
//
// Solidity: function transferLock(address locker, address to, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) TransferLock(locker common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.TransferLock(&_LPToken.TransactOpts, locker, to, amount)
}

// TransferLock is a paid mutator transaction binding the contract method 0x4ef09d79.
//
// Solidity: function transferLock(address locker, address to, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) TransferLock(locker common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.TransferLock(&_LPToken.TransactOpts, locker, to, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LPToken *LPTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LPToken *LPTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LPToken.Contract.TransferOwnership(&_LPToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LPToken *LPTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LPToken.Contract.TransferOwnership(&_LPToken.TransactOpts, newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address locker, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactor) Unlock(opts *bind.TransactOpts, locker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "unlock", locker, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address locker, uint256 amount) returns(bool)
func (_LPToken *LPTokenSession) Unlock(locker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Unlock(&_LPToken.TransactOpts, locker, amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address locker, uint256 amount) returns(bool)
func (_LPToken *LPTokenTransactorSession) Unlock(locker common.Address, amount *big.Int) (*types.Transaction, error) {
	return _LPToken.Contract.Unlock(&_LPToken.TransactOpts, locker, amount)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_LPToken *LPTokenTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_LPToken *LPTokenSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _LPToken.Contract.UpgradeTo(&_LPToken.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_LPToken *LPTokenTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _LPToken.Contract.UpgradeTo(&_LPToken.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LPToken *LPTokenTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LPToken.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LPToken *LPTokenSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LPToken.Contract.UpgradeToAndCall(&_LPToken.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LPToken *LPTokenTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LPToken.Contract.UpgradeToAndCall(&_LPToken.TransactOpts, newImplementation, data)
}

// LPTokenAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the LPToken contract.
type LPTokenAdminChangedIterator struct {
	Event *LPTokenAdminChanged // Event containing the contract specifics and raw log

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
func (it *LPTokenAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenAdminChanged)
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
		it.Event = new(LPTokenAdminChanged)
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
func (it *LPTokenAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenAdminChanged represents a AdminChanged event raised by the LPToken contract.
type LPTokenAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_LPToken *LPTokenFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*LPTokenAdminChangedIterator, error) {

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &LPTokenAdminChangedIterator{contract: _LPToken.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_LPToken *LPTokenFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *LPTokenAdminChanged) (event.Subscription, error) {

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenAdminChanged)
				if err := _LPToken.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_LPToken *LPTokenFilterer) ParseAdminChanged(log types.Log) (*LPTokenAdminChanged, error) {
	event := new(LPTokenAdminChanged)
	if err := _LPToken.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LPToken contract.
type LPTokenApprovalIterator struct {
	Event *LPTokenApproval // Event containing the contract specifics and raw log

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
func (it *LPTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenApproval)
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
		it.Event = new(LPTokenApproval)
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
func (it *LPTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenApproval represents a Approval event raised by the LPToken contract.
type LPTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LPToken *LPTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LPTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenApprovalIterator{contract: _LPToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LPToken *LPTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LPTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenApproval)
				if err := _LPToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_LPToken *LPTokenFilterer) ParseApproval(log types.Log) (*LPTokenApproval, error) {
	event := new(LPTokenApproval)
	if err := _LPToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenApprovalLockIterator is returned from FilterApprovalLock and is used to iterate over the raw logs and unpacked data for ApprovalLock events raised by the LPToken contract.
type LPTokenApprovalLockIterator struct {
	Event *LPTokenApprovalLock // Event containing the contract specifics and raw log

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
func (it *LPTokenApprovalLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenApprovalLock)
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
		it.Event = new(LPTokenApprovalLock)
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
func (it *LPTokenApprovalLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenApprovalLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenApprovalLock represents a ApprovalLock event raised by the LPToken contract.
type LPTokenApprovalLock struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApprovalLock is a free log retrieval operation binding the contract event 0x6a8ad50d47d7e8cfb288b0f40af42ea12bdc6f2ddcfe403f95854afa116fc8ee.
//
// Solidity: event ApprovalLock(address indexed owner, address indexed spender, uint256 value)
func (_LPToken *LPTokenFilterer) FilterApprovalLock(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LPTokenApprovalLockIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "ApprovalLock", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenApprovalLockIterator{contract: _LPToken.contract, event: "ApprovalLock", logs: logs, sub: sub}, nil
}

// WatchApprovalLock is a free log subscription operation binding the contract event 0x6a8ad50d47d7e8cfb288b0f40af42ea12bdc6f2ddcfe403f95854afa116fc8ee.
//
// Solidity: event ApprovalLock(address indexed owner, address indexed spender, uint256 value)
func (_LPToken *LPTokenFilterer) WatchApprovalLock(opts *bind.WatchOpts, sink chan<- *LPTokenApprovalLock, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "ApprovalLock", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenApprovalLock)
				if err := _LPToken.contract.UnpackLog(event, "ApprovalLock", log); err != nil {
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

// ParseApprovalLock is a log parse operation binding the contract event 0x6a8ad50d47d7e8cfb288b0f40af42ea12bdc6f2ddcfe403f95854afa116fc8ee.
//
// Solidity: event ApprovalLock(address indexed owner, address indexed spender, uint256 value)
func (_LPToken *LPTokenFilterer) ParseApprovalLock(log types.Log) (*LPTokenApprovalLock, error) {
	event := new(LPTokenApprovalLock)
	if err := _LPToken.contract.UnpackLog(event, "ApprovalLock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the LPToken contract.
type LPTokenBeaconUpgradedIterator struct {
	Event *LPTokenBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *LPTokenBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenBeaconUpgraded)
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
		it.Event = new(LPTokenBeaconUpgraded)
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
func (it *LPTokenBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenBeaconUpgraded represents a BeaconUpgraded event raised by the LPToken contract.
type LPTokenBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_LPToken *LPTokenFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*LPTokenBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenBeaconUpgradedIterator{contract: _LPToken.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_LPToken *LPTokenFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *LPTokenBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenBeaconUpgraded)
				if err := _LPToken.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_LPToken *LPTokenFilterer) ParseBeaconUpgraded(log types.Log) (*LPTokenBeaconUpgraded, error) {
	event := new(LPTokenBeaconUpgraded)
	if err := _LPToken.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the LPToken contract.
type LPTokenBurnIterator struct {
	Event *LPTokenBurn // Event containing the contract specifics and raw log

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
func (it *LPTokenBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenBurn)
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
		it.Event = new(LPTokenBurn)
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
func (it *LPTokenBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenBurn represents a Burn event raised by the LPToken contract.
type LPTokenBurn struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed from, uint256 value)
func (_LPToken *LPTokenFilterer) FilterBurn(opts *bind.FilterOpts, from []common.Address) (*LPTokenBurnIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Burn", fromRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenBurnIterator{contract: _LPToken.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed from, uint256 value)
func (_LPToken *LPTokenFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *LPTokenBurn, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Burn", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenBurn)
				if err := _LPToken.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5.
//
// Solidity: event Burn(address indexed from, uint256 value)
func (_LPToken *LPTokenFilterer) ParseBurn(log types.Log) (*LPTokenBurn, error) {
	event := new(LPTokenBurn)
	if err := _LPToken.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenLockIterator is returned from FilterLock and is used to iterate over the raw logs and unpacked data for Lock events raised by the LPToken contract.
type LPTokenLockIterator struct {
	Event *LPTokenLock // Event containing the contract specifics and raw log

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
func (it *LPTokenLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenLock)
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
		it.Event = new(LPTokenLock)
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
func (it *LPTokenLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenLock represents a Lock event raised by the LPToken contract.
type LPTokenLock struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLock is a free log retrieval operation binding the contract event 0xec36c0364d931187a76cf66d7eee08fad0ec2e8b7458a8d8b26b36769d4d13f3.
//
// Solidity: event Lock(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) FilterLock(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LPTokenLockIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Lock", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenLockIterator{contract: _LPToken.contract, event: "Lock", logs: logs, sub: sub}, nil
}

// WatchLock is a free log subscription operation binding the contract event 0xec36c0364d931187a76cf66d7eee08fad0ec2e8b7458a8d8b26b36769d4d13f3.
//
// Solidity: event Lock(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) WatchLock(opts *bind.WatchOpts, sink chan<- *LPTokenLock, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Lock", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenLock)
				if err := _LPToken.contract.UnpackLog(event, "Lock", log); err != nil {
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

// ParseLock is a log parse operation binding the contract event 0xec36c0364d931187a76cf66d7eee08fad0ec2e8b7458a8d8b26b36769d4d13f3.
//
// Solidity: event Lock(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) ParseLock(log types.Log) (*LPTokenLock, error) {
	event := new(LPTokenLock)
	if err := _LPToken.contract.UnpackLog(event, "Lock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the LPToken contract.
type LPTokenMintIterator struct {
	Event *LPTokenMint // Event containing the contract specifics and raw log

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
func (it *LPTokenMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenMint)
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
		it.Event = new(LPTokenMint)
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
func (it *LPTokenMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenMint represents a Mint event raised by the LPToken contract.
type LPTokenMint struct {
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) FilterMint(opts *bind.FilterOpts, to []common.Address) (*LPTokenMintIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Mint", toRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenMintIterator{contract: _LPToken.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *LPTokenMint, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Mint", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenMint)
				if err := _LPToken.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
//
// Solidity: event Mint(address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) ParseMint(log types.Log) (*LPTokenMint, error) {
	event := new(LPTokenMint)
	if err := _LPToken.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the LPToken contract.
type LPTokenOwnershipTransferredIterator struct {
	Event *LPTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *LPTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenOwnershipTransferred)
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
		it.Event = new(LPTokenOwnershipTransferred)
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
func (it *LPTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenOwnershipTransferred represents a OwnershipTransferred event raised by the LPToken contract.
type LPTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LPToken *LPTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LPTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenOwnershipTransferredIterator{contract: _LPToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LPToken *LPTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LPTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenOwnershipTransferred)
				if err := _LPToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_LPToken *LPTokenFilterer) ParseOwnershipTransferred(log types.Log) (*LPTokenOwnershipTransferred, error) {
	event := new(LPTokenOwnershipTransferred)
	if err := _LPToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LPToken contract.
type LPTokenTransferIterator struct {
	Event *LPTokenTransfer // Event containing the contract specifics and raw log

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
func (it *LPTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenTransfer)
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
		it.Event = new(LPTokenTransfer)
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
func (it *LPTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenTransfer represents a Transfer event raised by the LPToken contract.
type LPTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LPTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenTransferIterator{contract: _LPToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LPTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenTransfer)
				if err := _LPToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_LPToken *LPTokenFilterer) ParseTransfer(log types.Log) (*LPTokenTransfer, error) {
	event := new(LPTokenTransfer)
	if err := _LPToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenTransferLockIterator is returned from FilterTransferLock and is used to iterate over the raw logs and unpacked data for TransferLock events raised by the LPToken contract.
type LPTokenTransferLockIterator struct {
	Event *LPTokenTransferLock // Event containing the contract specifics and raw log

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
func (it *LPTokenTransferLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenTransferLock)
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
		it.Event = new(LPTokenTransferLock)
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
func (it *LPTokenTransferLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenTransferLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenTransferLock represents a TransferLock event raised by the LPToken contract.
type LPTokenTransferLock struct {
	Sender common.Address
	From   common.Address
	To     common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferLock is a free log retrieval operation binding the contract event 0x164e3db520d3c5f437914d09eef8478390d49d797a1b1c9aa5a0eed966183343.
//
// Solidity: event TransferLock(address indexed sender, address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) FilterTransferLock(opts *bind.FilterOpts, sender []common.Address, from []common.Address, to []common.Address) (*LPTokenTransferLockIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "TransferLock", senderRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenTransferLockIterator{contract: _LPToken.contract, event: "TransferLock", logs: logs, sub: sub}, nil
}

// WatchTransferLock is a free log subscription operation binding the contract event 0x164e3db520d3c5f437914d09eef8478390d49d797a1b1c9aa5a0eed966183343.
//
// Solidity: event TransferLock(address indexed sender, address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) WatchTransferLock(opts *bind.WatchOpts, sink chan<- *LPTokenTransferLock, sender []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "TransferLock", senderRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenTransferLock)
				if err := _LPToken.contract.UnpackLog(event, "TransferLock", log); err != nil {
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

// ParseTransferLock is a log parse operation binding the contract event 0x164e3db520d3c5f437914d09eef8478390d49d797a1b1c9aa5a0eed966183343.
//
// Solidity: event TransferLock(address indexed sender, address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) ParseTransferLock(log types.Log) (*LPTokenTransferLock, error) {
	event := new(LPTokenTransferLock)
	if err := _LPToken.contract.UnpackLog(event, "TransferLock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenUnlockIterator is returned from FilterUnlock and is used to iterate over the raw logs and unpacked data for Unlock events raised by the LPToken contract.
type LPTokenUnlockIterator struct {
	Event *LPTokenUnlock // Event containing the contract specifics and raw log

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
func (it *LPTokenUnlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenUnlock)
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
		it.Event = new(LPTokenUnlock)
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
func (it *LPTokenUnlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenUnlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenUnlock represents a Unlock event raised by the LPToken contract.
type LPTokenUnlock struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUnlock is a free log retrieval operation binding the contract event 0xc1c90b8e0705b212262c0dbd7580efe1862c2f185bf96899226f7596beb2db09.
//
// Solidity: event Unlock(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) FilterUnlock(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LPTokenUnlockIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Unlock", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenUnlockIterator{contract: _LPToken.contract, event: "Unlock", logs: logs, sub: sub}, nil
}

// WatchUnlock is a free log subscription operation binding the contract event 0xc1c90b8e0705b212262c0dbd7580efe1862c2f185bf96899226f7596beb2db09.
//
// Solidity: event Unlock(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) WatchUnlock(opts *bind.WatchOpts, sink chan<- *LPTokenUnlock, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Unlock", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenUnlock)
				if err := _LPToken.contract.UnpackLog(event, "Unlock", log); err != nil {
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

// ParseUnlock is a log parse operation binding the contract event 0xc1c90b8e0705b212262c0dbd7580efe1862c2f185bf96899226f7596beb2db09.
//
// Solidity: event Unlock(address indexed from, address indexed to, uint256 value)
func (_LPToken *LPTokenFilterer) ParseUnlock(log types.Log) (*LPTokenUnlock, error) {
	event := new(LPTokenUnlock)
	if err := _LPToken.contract.UnpackLog(event, "Unlock", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPTokenUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the LPToken contract.
type LPTokenUpgradedIterator struct {
	Event *LPTokenUpgraded // Event containing the contract specifics and raw log

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
func (it *LPTokenUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPTokenUpgraded)
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
		it.Event = new(LPTokenUpgraded)
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
func (it *LPTokenUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPTokenUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPTokenUpgraded represents a Upgraded event raised by the LPToken contract.
type LPTokenUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LPToken *LPTokenFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*LPTokenUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _LPToken.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &LPTokenUpgradedIterator{contract: _LPToken.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LPToken *LPTokenFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *LPTokenUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _LPToken.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPTokenUpgraded)
				if err := _LPToken.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_LPToken *LPTokenFilterer) ParseUpgraded(log types.Log) (*LPTokenUpgraded, error) {
	event := new(LPTokenUpgraded)
	if err := _LPToken.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
