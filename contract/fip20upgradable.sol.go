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

// FIP20UpgradableMetaData contains all meta data concerning the FIP20Upgradable contract.
var FIP20UpgradableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a060405261100160805234801561001657600080fd5b5060805161188461004e6000396000818161052201528181610562015281816106180152818161065801526106e701526118846000f3fe6080604052600436106101095760003560e01c806370a08231116100955780639dc29fac116100645780639dc29fac146102bc578063a9059cbb146102dc578063dd62ed3e146102fc578063de7ea79d14610342578063f2fde38b1461036257600080fd5b806370a0823114610234578063715018a61461026a5780638da5cb5b1461027f57806395d89b41146102a757600080fd5b8063313ce567116100dc578063313ce567146101a85780633659cfe6146101ca57806340c10f19146101ec5780634f1ef2861461020c57806352d1902d1461021f57600080fd5b806306fdde031461010e578063095ea7b31461013957806318160ddd1461016957806323b872dd14610188575b600080fd5b34801561011a57600080fd5b50610123610382565b60405161013091906113d1565b60405180910390f35b34801561014557600080fd5b50610159610154366004611420565b610414565b6040519015158152602001610130565b34801561017557600080fd5b5060cc545b604051908152602001610130565b34801561019457600080fd5b506101596101a336600461144a565b61046a565b3480156101b457600080fd5b5060cb5460405160ff9091168152602001610130565b3480156101d657600080fd5b506101ea6101e5366004611486565b610517565b005b3480156101f857600080fd5b506101ea610207366004611420565b6105f7565b6101ea61021a36600461152d565b61060d565b34801561022b57600080fd5b5061017a6106da565b34801561024057600080fd5b5061017a61024f366004611486565b6001600160a01b0316600090815260cd602052604090205490565b34801561027657600080fd5b506101ea61078d565b34801561028b57600080fd5b506097546040516001600160a01b039091168152602001610130565b3480156102b357600080fd5b506101236107a1565b3480156102c857600080fd5b506101ea6102d7366004611420565b6107b0565b3480156102e857600080fd5b506101596102f7366004611420565b6107c2565b34801561030857600080fd5b5061017a61031736600461158f565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b34801561034e57600080fd5b506101ea61035d3660046115e2565b6107d8565b34801561036e57600080fd5b506101ea61037d366004611486565b610947565b606060c980546103919061166f565b80601f01602080910402602001604051908101604052809291908181526020018280546103bd9061166f565b801561040a5780601f106103df5761010080835404028352916020019161040a565b820191906000526020600020905b8154815290600101906020018083116103ed57829003601f168201915b5050505050905090565b60006104213384846109bd565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156104ed5760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b61050185336104fc86856116c0565b6109bd565b61050c858585610a3f565b506001949350505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156105605760405162461bcd60e51b81526004016104e4906116d7565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166105a9600080516020611808833981519152546001600160a01b031690565b6001600160a01b0316146105cf5760405162461bcd60e51b81526004016104e490611723565b6105d881610bee565b604080516000808252602082019092526105f491839190610bf6565b50565b6105ff610d66565b6106098282610dc0565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106565760405162461bcd60e51b81526004016104e4906116d7565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661069f600080516020611808833981519152546001600160a01b031690565b6001600160a01b0316146106c55760405162461bcd60e51b81526004016104e490611723565b6106ce82610bee565b61060982826001610bf6565b6000306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461077a5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016104e4565b5060008051602061180883398151915290565b610795610d66565b61079f6000610e9f565b565b606060ca80546103919061166f565b6107b8610d66565b6106098282610ef1565b60006107cf338484610a3f565b50600192915050565b600054610100900460ff16158080156107f85750600054600160ff909116105b806108125750303b158015610812575060005460ff166001145b6108755760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016104e4565b6000805460ff191660011790558015610898576000805461ff0019166101001790555b84516108ab9060c990602088019061130c565b5083516108bf9060ca90602087019061130c565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b0384161790556108f2611033565b6108fa611062565b8015610940576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b61094f610d66565b6001600160a01b0381166109b45760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016104e4565b6105f481610e9f565b6001600160a01b038316610a135760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016104e4565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610a955760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016104e4565b6001600160a01b038216610aeb5760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016104e4565b6001600160a01b038316600090815260cd602052604090205481811015610b545760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016104e4565b610b5e82826116c0565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610b9490849061176f565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610be091815260200190565b60405180910390a350505050565b6105f4610d66565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610c2e57610c2983611089565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610c88575060408051601f3d908101601f19168201909252610c8591810190611787565b60015b610ceb5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016104e4565b6000805160206118088339815191528114610d5a5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016104e4565b50610c29838383611125565b6097546001600160a01b0316331461079f5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016104e4565b6001600160a01b038216610e165760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016104e4565b8060cc6000828254610e28919061176f565b90915550506001600160a01b038216600090815260cd602052604081208054839290610e5590849061176f565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b038216610f475760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016104e4565b6001600160a01b038216600090815260cd602052604090205481811015610fb05760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016104e4565b610fba82826116c0565b6001600160a01b038416600090815260cd602052604081209190915560cc8054849290610fe89084906116c0565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b600054610100900460ff1661105a5760405162461bcd60e51b81526004016104e4906117a0565b61079f611150565b600054610100900460ff1661079f5760405162461bcd60e51b81526004016104e4906117a0565b6001600160a01b0381163b6110f65760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016104e4565b60008051602061180883398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61112e83611180565b60008251118061113b5750805b15610c295761114a83836111c0565b50505050565b600054610100900460ff166111775760405162461bcd60e51b81526004016104e4906117a0565b61079f33610e9f565b61118981611089565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606111e58383604051806060016040528060278152602001611828602791396111ec565b9392505050565b6060600080856001600160a01b03168560405161120991906117eb565b600060405180830381855af49150503d8060008114611244576040519150601f19603f3d011682016040523d82523d6000602084013e611249565b606091505b509150915061125a86838387611264565b9695505050505050565b606083156112d05782516112c9576001600160a01b0385163b6112c95760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016104e4565b50816112da565b6112da83836112e2565b949350505050565b8151156112f25781518083602001fd5b8060405162461bcd60e51b81526004016104e491906113d1565b8280546113189061166f565b90600052602060002090601f01602090048101928261133a5760008555611380565b82601f1061135357805160ff1916838001178555611380565b82800160010185558215611380579182015b82811115611380578251825591602001919060010190611365565b5061138c929150611390565b5090565b5b8082111561138c5760008155600101611391565b60005b838110156113c05781810151838201526020016113a8565b8381111561114a5750506000910152565b60208152600082518060208401526113f08160408501602087016113a5565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461141b57600080fd5b919050565b6000806040838503121561143357600080fd5b61143c83611404565b946020939093013593505050565b60008060006060848603121561145f57600080fd5b61146884611404565b925061147660208501611404565b9150604084013590509250925092565b60006020828403121561149857600080fd5b6111e582611404565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff808411156114d2576114d26114a1565b604051601f8501601f19908116603f011681019082821181831017156114fa576114fa6114a1565b8160405280935085815286868601111561151357600080fd5b858560208301376000602087830101525050509392505050565b6000806040838503121561154057600080fd5b61154983611404565b9150602083013567ffffffffffffffff81111561156557600080fd5b8301601f8101851361157657600080fd5b611585858235602084016114b7565b9150509250929050565b600080604083850312156115a257600080fd5b6115ab83611404565b91506115b960208401611404565b90509250929050565b600082601f8301126115d357600080fd5b6111e5838335602085016114b7565b600080600080608085870312156115f857600080fd5b843567ffffffffffffffff8082111561161057600080fd5b61161c888389016115c2565b9550602087013591508082111561163257600080fd5b5061163f878288016115c2565b935050604085013560ff8116811461165657600080fd5b915061166460608601611404565b905092959194509250565b600181811c9082168061168357607f821691505b602082108114156116a457634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156116d2576116d26116aa565b500390565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b60008219821115611782576117826116aa565b500190565b60006020828403121561179957600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b600082516117fd8184602087016113a5565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220c24d636966a60df9636041ea92ff12a963fc67961fabb8564b1631506b7785e864736f6c634300080a0033",
}

// FIP20UpgradableABI is the input ABI used to generate the binding from.
// Deprecated: Use FIP20UpgradableMetaData.ABI instead.
var FIP20UpgradableABI = FIP20UpgradableMetaData.ABI

// FIP20UpgradableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FIP20UpgradableMetaData.Bin instead.
var FIP20UpgradableBin = FIP20UpgradableMetaData.Bin

// DeployFIP20Upgradable deploys a new Ethereum contract, binding an instance of FIP20Upgradable to it.
func DeployFIP20Upgradable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FIP20Upgradable, error) {
	parsed, err := FIP20UpgradableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FIP20UpgradableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FIP20Upgradable{FIP20UpgradableCaller: FIP20UpgradableCaller{contract: contract}, FIP20UpgradableTransactor: FIP20UpgradableTransactor{contract: contract}, FIP20UpgradableFilterer: FIP20UpgradableFilterer{contract: contract}}, nil
}

// FIP20Upgradable is an auto generated Go binding around an Ethereum contract.
type FIP20Upgradable struct {
	FIP20UpgradableCaller     // Read-only binding to the contract
	FIP20UpgradableTransactor // Write-only binding to the contract
	FIP20UpgradableFilterer   // Log filterer for contract events
}

// FIP20UpgradableCaller is an auto generated read-only Go binding around an Ethereum contract.
type FIP20UpgradableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIP20UpgradableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FIP20UpgradableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIP20UpgradableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FIP20UpgradableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIP20UpgradableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FIP20UpgradableSession struct {
	Contract     *FIP20Upgradable  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FIP20UpgradableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FIP20UpgradableCallerSession struct {
	Contract *FIP20UpgradableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// FIP20UpgradableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FIP20UpgradableTransactorSession struct {
	Contract     *FIP20UpgradableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// FIP20UpgradableRaw is an auto generated low-level Go binding around an Ethereum contract.
type FIP20UpgradableRaw struct {
	Contract *FIP20Upgradable // Generic contract binding to access the raw methods on
}

// FIP20UpgradableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FIP20UpgradableCallerRaw struct {
	Contract *FIP20UpgradableCaller // Generic read-only contract binding to access the raw methods on
}

// FIP20UpgradableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FIP20UpgradableTransactorRaw struct {
	Contract *FIP20UpgradableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFIP20Upgradable creates a new instance of FIP20Upgradable, bound to a specific deployed contract.
func NewFIP20Upgradable(address common.Address, backend bind.ContractBackend) (*FIP20Upgradable, error) {
	contract, err := bindFIP20Upgradable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FIP20Upgradable{FIP20UpgradableCaller: FIP20UpgradableCaller{contract: contract}, FIP20UpgradableTransactor: FIP20UpgradableTransactor{contract: contract}, FIP20UpgradableFilterer: FIP20UpgradableFilterer{contract: contract}}, nil
}

// NewFIP20UpgradableCaller creates a new read-only instance of FIP20Upgradable, bound to a specific deployed contract.
func NewFIP20UpgradableCaller(address common.Address, caller bind.ContractCaller) (*FIP20UpgradableCaller, error) {
	contract, err := bindFIP20Upgradable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableCaller{contract: contract}, nil
}

// NewFIP20UpgradableTransactor creates a new write-only instance of FIP20Upgradable, bound to a specific deployed contract.
func NewFIP20UpgradableTransactor(address common.Address, transactor bind.ContractTransactor) (*FIP20UpgradableTransactor, error) {
	contract, err := bindFIP20Upgradable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableTransactor{contract: contract}, nil
}

// NewFIP20UpgradableFilterer creates a new log filterer instance of FIP20Upgradable, bound to a specific deployed contract.
func NewFIP20UpgradableFilterer(address common.Address, filterer bind.ContractFilterer) (*FIP20UpgradableFilterer, error) {
	contract, err := bindFIP20Upgradable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableFilterer{contract: contract}, nil
}

// bindFIP20Upgradable binds a generic wrapper to an already deployed contract.
func bindFIP20Upgradable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FIP20UpgradableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FIP20Upgradable *FIP20UpgradableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FIP20Upgradable.Contract.FIP20UpgradableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FIP20Upgradable *FIP20UpgradableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.FIP20UpgradableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FIP20Upgradable *FIP20UpgradableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.FIP20UpgradableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FIP20Upgradable *FIP20UpgradableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FIP20Upgradable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FIP20Upgradable *FIP20UpgradableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FIP20Upgradable *FIP20UpgradableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FIP20Upgradable.Contract.Allowance(&_FIP20Upgradable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FIP20Upgradable.Contract.Allowance(&_FIP20Upgradable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FIP20Upgradable.Contract.BalanceOf(&_FIP20Upgradable.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FIP20Upgradable.Contract.BalanceOf(&_FIP20Upgradable.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FIP20Upgradable *FIP20UpgradableCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FIP20Upgradable *FIP20UpgradableSession) Decimals() (uint8, error) {
	return _FIP20Upgradable.Contract.Decimals(&_FIP20Upgradable.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) Decimals() (uint8, error) {
	return _FIP20Upgradable.Contract.Decimals(&_FIP20Upgradable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FIP20Upgradable *FIP20UpgradableCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FIP20Upgradable *FIP20UpgradableSession) Name() (string, error) {
	return _FIP20Upgradable.Contract.Name(&_FIP20Upgradable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) Name() (string, error) {
	return _FIP20Upgradable.Contract.Name(&_FIP20Upgradable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FIP20Upgradable *FIP20UpgradableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FIP20Upgradable *FIP20UpgradableSession) Owner() (common.Address, error) {
	return _FIP20Upgradable.Contract.Owner(&_FIP20Upgradable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) Owner() (common.Address, error) {
	return _FIP20Upgradable.Contract.Owner(&_FIP20Upgradable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FIP20Upgradable *FIP20UpgradableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FIP20Upgradable *FIP20UpgradableSession) ProxiableUUID() ([32]byte, error) {
	return _FIP20Upgradable.Contract.ProxiableUUID(&_FIP20Upgradable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _FIP20Upgradable.Contract.ProxiableUUID(&_FIP20Upgradable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FIP20Upgradable *FIP20UpgradableCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FIP20Upgradable *FIP20UpgradableSession) Symbol() (string, error) {
	return _FIP20Upgradable.Contract.Symbol(&_FIP20Upgradable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) Symbol() (string, error) {
	return _FIP20Upgradable.Contract.Symbol(&_FIP20Upgradable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableSession) TotalSupply() (*big.Int, error) {
	return _FIP20Upgradable.Contract.TotalSupply(&_FIP20Upgradable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) TotalSupply() (*big.Int, error) {
	return _FIP20Upgradable.Contract.TotalSupply(&_FIP20Upgradable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Approve(&_FIP20Upgradable.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Approve(&_FIP20Upgradable.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "burn", account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_FIP20Upgradable *FIP20UpgradableSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Burn(&_FIP20Upgradable.TransactOpts, account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Burn(&_FIP20Upgradable.TransactOpts, account, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "initialize", name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_FIP20Upgradable *FIP20UpgradableSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Initialize(&_FIP20Upgradable.TransactOpts, name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Initialize(&_FIP20Upgradable.TransactOpts, name_, symbol_, decimals_, module_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_FIP20Upgradable *FIP20UpgradableSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Mint(&_FIP20Upgradable.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Mint(&_FIP20Upgradable.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FIP20Upgradable *FIP20UpgradableSession) RenounceOwnership() (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.RenounceOwnership(&_FIP20Upgradable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.RenounceOwnership(&_FIP20Upgradable.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Transfer(&_FIP20Upgradable.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.Transfer(&_FIP20Upgradable.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.TransferFrom(&_FIP20Upgradable.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.TransferFrom(&_FIP20Upgradable.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FIP20Upgradable *FIP20UpgradableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.TransferOwnership(&_FIP20Upgradable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.TransferOwnership(&_FIP20Upgradable.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FIP20Upgradable *FIP20UpgradableSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.UpgradeTo(&_FIP20Upgradable.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.UpgradeTo(&_FIP20Upgradable.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FIP20Upgradable *FIP20UpgradableTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FIP20Upgradable *FIP20UpgradableSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.UpgradeToAndCall(&_FIP20Upgradable.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.UpgradeToAndCall(&_FIP20Upgradable.TransactOpts, newImplementation, data)
}

// FIP20UpgradableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the FIP20Upgradable contract.
type FIP20UpgradableAdminChangedIterator struct {
	Event *FIP20UpgradableAdminChanged // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableAdminChanged)
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
		it.Event = new(FIP20UpgradableAdminChanged)
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
func (it *FIP20UpgradableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableAdminChanged represents a AdminChanged event raised by the FIP20Upgradable contract.
type FIP20UpgradableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*FIP20UpgradableAdminChangedIterator, error) {

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableAdminChangedIterator{contract: _FIP20Upgradable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableAdminChanged) (event.Subscription, error) {

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableAdminChanged)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseAdminChanged(log types.Log) (*FIP20UpgradableAdminChanged, error) {
	event := new(FIP20UpgradableAdminChanged)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the FIP20Upgradable contract.
type FIP20UpgradableApprovalIterator struct {
	Event *FIP20UpgradableApproval // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableApproval)
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
		it.Event = new(FIP20UpgradableApproval)
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
func (it *FIP20UpgradableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableApproval represents a Approval event raised by the FIP20Upgradable contract.
type FIP20UpgradableApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*FIP20UpgradableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableApprovalIterator{contract: _FIP20Upgradable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableApproval)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseApproval(log types.Log) (*FIP20UpgradableApproval, error) {
	event := new(FIP20UpgradableApproval)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the FIP20Upgradable contract.
type FIP20UpgradableBeaconUpgradedIterator struct {
	Event *FIP20UpgradableBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableBeaconUpgraded)
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
		it.Event = new(FIP20UpgradableBeaconUpgraded)
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
func (it *FIP20UpgradableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableBeaconUpgraded represents a BeaconUpgraded event raised by the FIP20Upgradable contract.
type FIP20UpgradableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*FIP20UpgradableBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableBeaconUpgradedIterator{contract: _FIP20Upgradable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableBeaconUpgraded)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseBeaconUpgraded(log types.Log) (*FIP20UpgradableBeaconUpgraded, error) {
	event := new(FIP20UpgradableBeaconUpgraded)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FIP20Upgradable contract.
type FIP20UpgradableInitializedIterator struct {
	Event *FIP20UpgradableInitialized // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableInitialized)
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
		it.Event = new(FIP20UpgradableInitialized)
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
func (it *FIP20UpgradableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableInitialized represents a Initialized event raised by the FIP20Upgradable contract.
type FIP20UpgradableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterInitialized(opts *bind.FilterOpts) (*FIP20UpgradableInitializedIterator, error) {

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableInitializedIterator{contract: _FIP20Upgradable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableInitialized) (event.Subscription, error) {

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableInitialized)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseInitialized(log types.Log) (*FIP20UpgradableInitialized, error) {
	event := new(FIP20UpgradableInitialized)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FIP20Upgradable contract.
type FIP20UpgradableOwnershipTransferredIterator struct {
	Event *FIP20UpgradableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableOwnershipTransferred)
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
		it.Event = new(FIP20UpgradableOwnershipTransferred)
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
func (it *FIP20UpgradableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableOwnershipTransferred represents a OwnershipTransferred event raised by the FIP20Upgradable contract.
type FIP20UpgradableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FIP20UpgradableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableOwnershipTransferredIterator{contract: _FIP20Upgradable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableOwnershipTransferred)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseOwnershipTransferred(log types.Log) (*FIP20UpgradableOwnershipTransferred, error) {
	event := new(FIP20UpgradableOwnershipTransferred)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the FIP20Upgradable contract.
type FIP20UpgradableTransferIterator struct {
	Event *FIP20UpgradableTransfer // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableTransfer)
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
		it.Event = new(FIP20UpgradableTransfer)
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
func (it *FIP20UpgradableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableTransfer represents a Transfer event raised by the FIP20Upgradable contract.
type FIP20UpgradableTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FIP20UpgradableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableTransferIterator{contract: _FIP20Upgradable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableTransfer)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseTransfer(log types.Log) (*FIP20UpgradableTransfer, error) {
	event := new(FIP20UpgradableTransfer)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FIP20Upgradable contract.
type FIP20UpgradableUpgradedIterator struct {
	Event *FIP20UpgradableUpgraded // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableUpgraded)
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
		it.Event = new(FIP20UpgradableUpgraded)
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
func (it *FIP20UpgradableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableUpgraded represents a Upgraded event raised by the FIP20Upgradable contract.
type FIP20UpgradableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FIP20UpgradableUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableUpgradedIterator{contract: _FIP20Upgradable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableUpgraded)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseUpgraded(log types.Log) (*FIP20UpgradableUpgraded, error) {
	event := new(FIP20UpgradableUpgraded)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
