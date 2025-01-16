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

// WarpTokenUpgradableMetaData contains all meta data concerning the WarpTokenUpgradable contract.
var WarpTokenUpgradableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a060405261100260805234801561001657600080fd5b50608051611a8161004e6000396000818161067c015281816106bc01528181610772015281816107b201526108410152611a816000f3fe6080604052600436106101395760003560e01c8063715018a6116100ab578063b86d52981161006f578063b86d529814610366578063d0e30db014610148578063dd62ed3e14610384578063de7ea79d146103ca578063f2fde38b146103ea578063f3fef3a31461040a57610148565b8063715018a6146102ca5780638da5cb5b146102df57806395d89b41146103115780639dc29fac14610326578063a9059cbb1461034657610148565b8063313ce567116100fd578063313ce5671461020a5780633659cfe61461022c57806340c10f191461024c5780634f1ef2861461026c57806352d1902d1461027f57806370a082311461029457610148565b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab57806323b872dd146101ca5780632e1a7d4d146101ea57610148565b366101485761014661042a565b005b61014661042a565b34801561015c57600080fd5b5061016561046b565b60405161017291906115a9565b60405180910390f35b34801561018757600080fd5b5061019b6101963660046115f1565b6104fd565b6040519015158152602001610172565b3480156101b757600080fd5b5060cc545b604051908152602001610172565b3480156101d657600080fd5b5061019b6101e536600461161d565b610553565b3480156101f657600080fd5b5061014661020536600461165e565b610600565b34801561021657600080fd5b5060cb5460405160ff9091168152602001610172565b34801561023857600080fd5b50610146610247366004611677565b610671565b34801561025857600080fd5b506101466102673660046115f1565b610751565b61014661027a366004611720565b610767565b34801561028b57600080fd5b506101bc610834565b3480156102a057600080fd5b506101bc6102af366004611677565b6001600160a01b0316600090815260cd602052604090205490565b3480156102d657600080fd5b506101466108e7565b3480156102eb57600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610172565b34801561031d57600080fd5b506101656108fb565b34801561033257600080fd5b506101466103413660046115f1565b61090a565b34801561035257600080fd5b5061019b6103613660046115f1565b61091c565b34801561037257600080fd5b5060cf546001600160a01b03166102f9565b34801561039057600080fd5b506101bc61039f366004611784565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103d657600080fd5b506101466103e53660046117dd565b610932565b3480156103f657600080fd5b50610146610405366004611677565b610aa1565b34801561041657600080fd5b506101466104253660046115f1565b610b17565b6104343334610b9c565b60405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b606060c9805461047a9061186c565b80601f01602080910402602001604051908101604052809291908181526020018280546104a69061186c565b80156104f35780601f106104c8576101008083540402835291602001916104f3565b820191906000526020600020905b8154815290600101906020018083116104d657829003601f168201915b5050505050905090565b600061050a338484610c74565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156105d65760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b6105ea85336105e586856118bd565b610c74565b6105f5858585610cf6565b506001949350505050565b61060b335b82610ea5565b604051339082156108fc029083906000818181858888f19350505050158015610638573d6000803e3d6000fd5b5060405181815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a250565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106ba5760405162461bcd60e51b81526004016105cd906118d4565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610703600080516020611a05833981519152546001600160a01b031690565b6001600160a01b0316146107295760405162461bcd60e51b81526004016105cd90611920565b61073281610fe7565b6040805160008082526020820190925261074e91839190610fef565b50565b61075961115f565b6107638282610b9c565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156107b05760405162461bcd60e51b81526004016105cd906118d4565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166107f9600080516020611a05833981519152546001600160a01b031690565b6001600160a01b03161461081f5760405162461bcd60e51b81526004016105cd90611920565b61082882610fe7565b61076382826001610fef565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146108d45760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016105cd565b50600080516020611a0583398151915290565b6108ef61115f565b6108f960006111b9565b565b606060ca805461047a9061186c565b61091261115f565b6107638282610ea5565b6000610929338484610cf6565b50600192915050565b600054610100900460ff16158080156109525750600054600160ff909116105b8061096c5750303b15801561096c575060005460ff166001145b6109cf5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016105cd565b6000805460ff1916600117905580156109f2576000805461ff0019166101001790555b8451610a059060c99060208801906114e4565b508351610a199060ca9060208701906114e4565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610a4c61120b565b610a5461123a565b8015610a9a576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b610aa961115f565b6001600160a01b038116610b0e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016105cd565b61074e816111b9565b610b2033610605565b6040516001600160a01b0383169082156108fc029083906000818181858888f19350505050158015610b56573d6000803e3d6000fd5b506040518181526001600160a01b0383169033907f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb906020015b60405180910390a35050565b6001600160a01b038216610bf25760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016105cd565b8060cc6000828254610c04919061196c565b90915550506001600160a01b038216600090815260cd602052604081208054839290610c3190849061196c565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610b90565b6001600160a01b038316610cca5760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016105cd565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610d4c5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016105cd565b6001600160a01b038216610da25760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016105cd565b6001600160a01b038316600090815260cd602052604090205481811015610e0b5760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016105cd565b610e1582826118bd565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610e4b90849061196c565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610e9791815260200190565b60405180910390a350505050565b6001600160a01b038216610efb5760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016105cd565b6001600160a01b038216600090815260cd602052604090205481811015610f645760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016105cd565b610f6e82826118bd565b6001600160a01b038416600090815260cd602052604081209190915560cc8054849290610f9c9084906118bd565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b61074e61115f565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156110275761102283611261565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611081575060408051601f3d908101601f1916820190925261107e91810190611984565b60015b6110e45760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016105cd565b600080516020611a0583398151915281146111535760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016105cd565b506110228383836112fd565b6097546001600160a01b031633146108f95760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105cd565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166112325760405162461bcd60e51b81526004016105cd9061199d565b6108f9611328565b600054610100900460ff166108f95760405162461bcd60e51b81526004016105cd9061199d565b6001600160a01b0381163b6112ce5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016105cd565b600080516020611a0583398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61130683611358565b6000825111806113135750805b15611022576113228383611398565b50505050565b600054610100900460ff1661134f5760405162461bcd60e51b81526004016105cd9061199d565b6108f9336111b9565b61136181611261565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606113bd8383604051806060016040528060278152602001611a25602791396113c4565b9392505050565b6060600080856001600160a01b0316856040516113e191906119e8565b600060405180830381855af49150503d806000811461141c576040519150601f19603f3d011682016040523d82523d6000602084013e611421565b606091505b50915091506114328683838761143c565b9695505050505050565b606083156114a85782516114a1576001600160a01b0385163b6114a15760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016105cd565b50816114b2565b6114b283836114ba565b949350505050565b8151156114ca5781518083602001fd5b8060405162461bcd60e51b81526004016105cd91906115a9565b8280546114f09061186c565b90600052602060002090601f0160209004810192826115125760008555611558565b82601f1061152b57805160ff1916838001178555611558565b82800160010185558215611558579182015b8281111561155857825182559160200191906001019061153d565b50611564929150611568565b5090565b5b808211156115645760008155600101611569565b60005b83811015611598578181015183820152602001611580565b838111156113225750506000910152565b60208152600082518060208401526115c881604085016020870161157d565b601f01601f19169190910160400192915050565b6001600160a01b038116811461074e57600080fd5b6000806040838503121561160457600080fd5b823561160f816115dc565b946020939093013593505050565b60008060006060848603121561163257600080fd5b833561163d816115dc565b9250602084013561164d816115dc565b929592945050506040919091013590565b60006020828403121561167057600080fd5b5035919050565b60006020828403121561168957600080fd5b81356113bd816115dc565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff808411156116c5576116c5611694565b604051601f8501601f19908116603f011681019082821181831017156116ed576116ed611694565b8160405280935085815286868601111561170657600080fd5b858560208301376000602087830101525050509392505050565b6000806040838503121561173357600080fd5b823561173e816115dc565b9150602083013567ffffffffffffffff81111561175a57600080fd5b8301601f8101851361176b57600080fd5b61177a858235602084016116aa565b9150509250929050565b6000806040838503121561179757600080fd5b82356117a2816115dc565b915060208301356117b2816115dc565b809150509250929050565b600082601f8301126117ce57600080fd5b6113bd838335602085016116aa565b600080600080608085870312156117f357600080fd5b843567ffffffffffffffff8082111561180b57600080fd5b611817888389016117bd565b9550602087013591508082111561182d57600080fd5b5061183a878288016117bd565b935050604085013560ff8116811461185157600080fd5b91506060850135611861816115dc565b939692955090935050565b600181811c9082168061188057607f821691505b602082108114156118a157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156118cf576118cf6118a7565b500390565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6000821982111561197f5761197f6118a7565b500190565b60006020828403121561199657600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b600082516119fa81846020870161157d565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122095ea5db6f864c5f06b9809babcbcb37d1d3a8a5650c75771cb6ab9d5b29ede0b64736f6c634300080a0033",
}

// WarpTokenUpgradableABI is the input ABI used to generate the binding from.
// Deprecated: Use WarpTokenUpgradableMetaData.ABI instead.
var WarpTokenUpgradableABI = WarpTokenUpgradableMetaData.ABI

// WarpTokenUpgradableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WarpTokenUpgradableMetaData.Bin instead.
var WarpTokenUpgradableBin = WarpTokenUpgradableMetaData.Bin

// DeployWarpTokenUpgradable deploys a new Ethereum contract, binding an instance of WarpTokenUpgradable to it.
func DeployWarpTokenUpgradable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WarpTokenUpgradable, error) {
	parsed, err := WarpTokenUpgradableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WarpTokenUpgradableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WarpTokenUpgradable{WarpTokenUpgradableCaller: WarpTokenUpgradableCaller{contract: contract}, WarpTokenUpgradableTransactor: WarpTokenUpgradableTransactor{contract: contract}, WarpTokenUpgradableFilterer: WarpTokenUpgradableFilterer{contract: contract}}, nil
}

// WarpTokenUpgradable is an auto generated Go binding around an Ethereum contract.
type WarpTokenUpgradable struct {
	WarpTokenUpgradableCaller     // Read-only binding to the contract
	WarpTokenUpgradableTransactor // Write-only binding to the contract
	WarpTokenUpgradableFilterer   // Log filterer for contract events
}

// WarpTokenUpgradableCaller is an auto generated read-only Go binding around an Ethereum contract.
type WarpTokenUpgradableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WarpTokenUpgradableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WarpTokenUpgradableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WarpTokenUpgradableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WarpTokenUpgradableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WarpTokenUpgradableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WarpTokenUpgradableSession struct {
	Contract     *WarpTokenUpgradable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// WarpTokenUpgradableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WarpTokenUpgradableCallerSession struct {
	Contract *WarpTokenUpgradableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// WarpTokenUpgradableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WarpTokenUpgradableTransactorSession struct {
	Contract     *WarpTokenUpgradableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// WarpTokenUpgradableRaw is an auto generated low-level Go binding around an Ethereum contract.
type WarpTokenUpgradableRaw struct {
	Contract *WarpTokenUpgradable // Generic contract binding to access the raw methods on
}

// WarpTokenUpgradableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WarpTokenUpgradableCallerRaw struct {
	Contract *WarpTokenUpgradableCaller // Generic read-only contract binding to access the raw methods on
}

// WarpTokenUpgradableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WarpTokenUpgradableTransactorRaw struct {
	Contract *WarpTokenUpgradableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWarpTokenUpgradable creates a new instance of WarpTokenUpgradable, bound to a specific deployed contract.
func NewWarpTokenUpgradable(address common.Address, backend bind.ContractBackend) (*WarpTokenUpgradable, error) {
	contract, err := bindWarpTokenUpgradable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradable{WarpTokenUpgradableCaller: WarpTokenUpgradableCaller{contract: contract}, WarpTokenUpgradableTransactor: WarpTokenUpgradableTransactor{contract: contract}, WarpTokenUpgradableFilterer: WarpTokenUpgradableFilterer{contract: contract}}, nil
}

// NewWarpTokenUpgradableCaller creates a new read-only instance of WarpTokenUpgradable, bound to a specific deployed contract.
func NewWarpTokenUpgradableCaller(address common.Address, caller bind.ContractCaller) (*WarpTokenUpgradableCaller, error) {
	contract, err := bindWarpTokenUpgradable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableCaller{contract: contract}, nil
}

// NewWarpTokenUpgradableTransactor creates a new write-only instance of WarpTokenUpgradable, bound to a specific deployed contract.
func NewWarpTokenUpgradableTransactor(address common.Address, transactor bind.ContractTransactor) (*WarpTokenUpgradableTransactor, error) {
	contract, err := bindWarpTokenUpgradable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableTransactor{contract: contract}, nil
}

// NewWarpTokenUpgradableFilterer creates a new log filterer instance of WarpTokenUpgradable, bound to a specific deployed contract.
func NewWarpTokenUpgradableFilterer(address common.Address, filterer bind.ContractFilterer) (*WarpTokenUpgradableFilterer, error) {
	contract, err := bindWarpTokenUpgradable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableFilterer{contract: contract}, nil
}

// bindWarpTokenUpgradable binds a generic wrapper to an already deployed contract.
func bindWarpTokenUpgradable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WarpTokenUpgradableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WarpTokenUpgradable *WarpTokenUpgradableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WarpTokenUpgradable.Contract.WarpTokenUpgradableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WarpTokenUpgradable *WarpTokenUpgradableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.WarpTokenUpgradableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WarpTokenUpgradable *WarpTokenUpgradableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.WarpTokenUpgradableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WarpTokenUpgradable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WarpTokenUpgradable.Contract.Allowance(&_WarpTokenUpgradable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WarpTokenUpgradable.Contract.Allowance(&_WarpTokenUpgradable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WarpTokenUpgradable.Contract.BalanceOf(&_WarpTokenUpgradable.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WarpTokenUpgradable.Contract.BalanceOf(&_WarpTokenUpgradable.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Decimals() (uint8, error) {
	return _WarpTokenUpgradable.Contract.Decimals(&_WarpTokenUpgradable.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) Decimals() (uint8, error) {
	return _WarpTokenUpgradable.Contract.Decimals(&_WarpTokenUpgradable.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) Module(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "module")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Module() (common.Address, error) {
	return _WarpTokenUpgradable.Contract.Module(&_WarpTokenUpgradable.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) Module() (common.Address, error) {
	return _WarpTokenUpgradable.Contract.Module(&_WarpTokenUpgradable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Name() (string, error) {
	return _WarpTokenUpgradable.Contract.Name(&_WarpTokenUpgradable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) Name() (string, error) {
	return _WarpTokenUpgradable.Contract.Name(&_WarpTokenUpgradable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Owner() (common.Address, error) {
	return _WarpTokenUpgradable.Contract.Owner(&_WarpTokenUpgradable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) Owner() (common.Address, error) {
	return _WarpTokenUpgradable.Contract.Owner(&_WarpTokenUpgradable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) ProxiableUUID() ([32]byte, error) {
	return _WarpTokenUpgradable.Contract.ProxiableUUID(&_WarpTokenUpgradable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _WarpTokenUpgradable.Contract.ProxiableUUID(&_WarpTokenUpgradable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Symbol() (string, error) {
	return _WarpTokenUpgradable.Contract.Symbol(&_WarpTokenUpgradable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) Symbol() (string, error) {
	return _WarpTokenUpgradable.Contract.Symbol(&_WarpTokenUpgradable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WarpTokenUpgradable.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) TotalSupply() (*big.Int, error) {
	return _WarpTokenUpgradable.Contract.TotalSupply(&_WarpTokenUpgradable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WarpTokenUpgradable *WarpTokenUpgradableCallerSession) TotalSupply() (*big.Int, error) {
	return _WarpTokenUpgradable.Contract.TotalSupply(&_WarpTokenUpgradable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Approve(&_WarpTokenUpgradable.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Approve(&_WarpTokenUpgradable.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "burn", account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Burn(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Burn(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Deposit() (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Deposit(&_WarpTokenUpgradable.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Deposit() (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Deposit(&_WarpTokenUpgradable.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "initialize", name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Initialize(&_WarpTokenUpgradable.TransactOpts, name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Initialize(&_WarpTokenUpgradable.TransactOpts, name_, symbol_, decimals_, module_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Mint(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Mint(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) RenounceOwnership() (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.RenounceOwnership(&_WarpTokenUpgradable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.RenounceOwnership(&_WarpTokenUpgradable.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Transfer(&_WarpTokenUpgradable.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Transfer(&_WarpTokenUpgradable.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.TransferFrom(&_WarpTokenUpgradable.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.TransferFrom(&_WarpTokenUpgradable.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.TransferOwnership(&_WarpTokenUpgradable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.TransferOwnership(&_WarpTokenUpgradable.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.UpgradeTo(&_WarpTokenUpgradable.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.UpgradeTo(&_WarpTokenUpgradable.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.UpgradeToAndCall(&_WarpTokenUpgradable.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.UpgradeToAndCall(&_WarpTokenUpgradable.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Withdraw(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "withdraw", value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Withdraw(&_WarpTokenUpgradable.TransactOpts, value)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 value) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Withdraw(value *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Withdraw(&_WarpTokenUpgradable.TransactOpts, value)
}

// Withdraw0 is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address to, uint256 value) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Withdraw0(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "withdraw0", to, value)
}

// Withdraw0 is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address to, uint256 value) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Withdraw0(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Withdraw0(&_WarpTokenUpgradable.TransactOpts, to, value)
}

// Withdraw0 is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address to, uint256 value) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Withdraw0(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Withdraw0(&_WarpTokenUpgradable.TransactOpts, to, value)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Fallback(&_WarpTokenUpgradable.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Fallback(&_WarpTokenUpgradable.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Receive() (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Receive(&_WarpTokenUpgradable.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Receive() (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Receive(&_WarpTokenUpgradable.TransactOpts)
}

// WarpTokenUpgradableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableAdminChangedIterator struct {
	Event *WarpTokenUpgradableAdminChanged // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableAdminChanged)
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
		it.Event = new(WarpTokenUpgradableAdminChanged)
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
func (it *WarpTokenUpgradableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableAdminChanged represents a AdminChanged event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*WarpTokenUpgradableAdminChangedIterator, error) {

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableAdminChangedIterator{contract: _WarpTokenUpgradable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableAdminChanged) (event.Subscription, error) {

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableAdminChanged)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseAdminChanged(log types.Log) (*WarpTokenUpgradableAdminChanged, error) {
	event := new(WarpTokenUpgradableAdminChanged)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableApprovalIterator struct {
	Event *WarpTokenUpgradableApproval // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableApproval)
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
		it.Event = new(WarpTokenUpgradableApproval)
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
func (it *WarpTokenUpgradableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableApproval represents a Approval event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WarpTokenUpgradableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableApprovalIterator{contract: _WarpTokenUpgradable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableApproval)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseApproval(log types.Log) (*WarpTokenUpgradableApproval, error) {
	event := new(WarpTokenUpgradableApproval)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableBeaconUpgradedIterator struct {
	Event *WarpTokenUpgradableBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableBeaconUpgraded)
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
		it.Event = new(WarpTokenUpgradableBeaconUpgraded)
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
func (it *WarpTokenUpgradableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableBeaconUpgraded represents a BeaconUpgraded event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*WarpTokenUpgradableBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableBeaconUpgradedIterator{contract: _WarpTokenUpgradable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableBeaconUpgraded)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseBeaconUpgraded(log types.Log) (*WarpTokenUpgradableBeaconUpgraded, error) {
	event := new(WarpTokenUpgradableBeaconUpgraded)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableDepositIterator struct {
	Event *WarpTokenUpgradableDeposit // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableDeposit)
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
		it.Event = new(WarpTokenUpgradableDeposit)
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
func (it *WarpTokenUpgradableDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableDeposit represents a Deposit event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableDeposit struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed from, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterDeposit(opts *bind.FilterOpts, from []common.Address) (*WarpTokenUpgradableDepositIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Deposit", fromRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableDepositIterator{contract: _WarpTokenUpgradable.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed from, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableDeposit, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Deposit", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableDeposit)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Deposit", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseDeposit(log types.Log) (*WarpTokenUpgradableDeposit, error) {
	event := new(WarpTokenUpgradableDeposit)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableInitializedIterator struct {
	Event *WarpTokenUpgradableInitialized // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableInitialized)
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
		it.Event = new(WarpTokenUpgradableInitialized)
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
func (it *WarpTokenUpgradableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableInitialized represents a Initialized event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterInitialized(opts *bind.FilterOpts) (*WarpTokenUpgradableInitializedIterator, error) {

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableInitializedIterator{contract: _WarpTokenUpgradable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableInitialized) (event.Subscription, error) {

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableInitialized)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseInitialized(log types.Log) (*WarpTokenUpgradableInitialized, error) {
	event := new(WarpTokenUpgradableInitialized)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableOwnershipTransferredIterator struct {
	Event *WarpTokenUpgradableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableOwnershipTransferred)
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
		it.Event = new(WarpTokenUpgradableOwnershipTransferred)
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
func (it *WarpTokenUpgradableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableOwnershipTransferred represents a OwnershipTransferred event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WarpTokenUpgradableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableOwnershipTransferredIterator{contract: _WarpTokenUpgradable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableOwnershipTransferred)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseOwnershipTransferred(log types.Log) (*WarpTokenUpgradableOwnershipTransferred, error) {
	event := new(WarpTokenUpgradableOwnershipTransferred)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableTransferIterator struct {
	Event *WarpTokenUpgradableTransfer // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableTransfer)
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
		it.Event = new(WarpTokenUpgradableTransfer)
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
func (it *WarpTokenUpgradableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableTransfer represents a Transfer event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WarpTokenUpgradableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableTransferIterator{contract: _WarpTokenUpgradable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableTransfer)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseTransfer(log types.Log) (*WarpTokenUpgradableTransfer, error) {
	event := new(WarpTokenUpgradableTransfer)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableUpgradedIterator struct {
	Event *WarpTokenUpgradableUpgraded // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableUpgraded)
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
		it.Event = new(WarpTokenUpgradableUpgraded)
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
func (it *WarpTokenUpgradableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableUpgraded represents a Upgraded event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*WarpTokenUpgradableUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableUpgradedIterator{contract: _WarpTokenUpgradable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableUpgraded)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseUpgraded(log types.Log) (*WarpTokenUpgradableUpgraded, error) {
	event := new(WarpTokenUpgradableUpgraded)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableWithdrawIterator struct {
	Event *WarpTokenUpgradableWithdraw // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableWithdraw)
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
		it.Event = new(WarpTokenUpgradableWithdraw)
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
func (it *WarpTokenUpgradableWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableWithdraw represents a Withdraw event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableWithdraw struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: event Withdraw(address indexed from, address indexed to, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WarpTokenUpgradableWithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Withdraw", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableWithdrawIterator{contract: _WarpTokenUpgradable.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: event Withdraw(address indexed from, address indexed to, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableWithdraw, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Withdraw", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableWithdraw)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Withdraw", log); err != nil {
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
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseWithdraw(log types.Log) (*WarpTokenUpgradableWithdraw, error) {
	event := new(WarpTokenUpgradableWithdraw)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WarpTokenUpgradableWithdraw0Iterator is returned from FilterWithdraw0 and is used to iterate over the raw logs and unpacked data for Withdraw0 events raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableWithdraw0Iterator struct {
	Event *WarpTokenUpgradableWithdraw0 // Event containing the contract specifics and raw log

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
func (it *WarpTokenUpgradableWithdraw0Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WarpTokenUpgradableWithdraw0)
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
		it.Event = new(WarpTokenUpgradableWithdraw0)
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
func (it *WarpTokenUpgradableWithdraw0Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WarpTokenUpgradableWithdraw0Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WarpTokenUpgradableWithdraw0 represents a Withdraw0 event raised by the WarpTokenUpgradable contract.
type WarpTokenUpgradableWithdraw0 struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWithdraw0 is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed from, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) FilterWithdraw0(opts *bind.FilterOpts, from []common.Address) (*WarpTokenUpgradableWithdraw0Iterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.FilterLogs(opts, "Withdraw0", fromRule)
	if err != nil {
		return nil, err
	}
	return &WarpTokenUpgradableWithdraw0Iterator{contract: _WarpTokenUpgradable.contract, event: "Withdraw0", logs: logs, sub: sub}, nil
}

// WatchWithdraw0 is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed from, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) WatchWithdraw0(opts *bind.WatchOpts, sink chan<- *WarpTokenUpgradableWithdraw0, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _WarpTokenUpgradable.contract.WatchLogs(opts, "Withdraw0", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WarpTokenUpgradableWithdraw0)
				if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Withdraw0", log); err != nil {
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

// ParseWithdraw0 is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed from, uint256 value)
func (_WarpTokenUpgradable *WarpTokenUpgradableFilterer) ParseWithdraw0(log types.Log) (*WarpTokenUpgradableWithdraw0, error) {
	event := new(WarpTokenUpgradableWithdraw0)
	if err := _WarpTokenUpgradable.contract.UnpackLog(event, "Withdraw0", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
