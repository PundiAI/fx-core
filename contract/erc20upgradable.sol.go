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

// ERC20UpgradableMetaData contains all meta data concerning the ERC20Upgradable contract.
var ERC20UpgradableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a060405261100160805234801561001657600080fd5b5060805161191461004e600039600081816104ef01528181610538015281816105f80152818161063801526106c701526119146000f3fe60806040526004361061011f5760003560e01c806370a08231116100a05780639dc29fac116100645780639dc29fac14610312578063a9059cbb14610332578063dd62ed3e14610352578063de7ea79d14610398578063f2fde38b146103b857600080fd5b806370a082311461026a578063715018a6146102a057806379cc6790146102b55780638da5cb5b146102d557806395d89b41146102fd57600080fd5b80633659cfe6116100e75780633659cfe6146101e057806340c10f191461020257806342966c68146102225780634f1ef2861461024257806352d1902d1461025557600080fd5b806306fdde0314610124578063095ea7b31461014f57806318160ddd1461017f57806323b872dd1461019e578063313ce567146101be575b600080fd5b34801561013057600080fd5b506101396103d8565b6040516101469190611448565b60405180910390f35b34801561015b57600080fd5b5061016f61016a366004611497565b61046a565b6040519015158152602001610146565b34801561018b57600080fd5b5060cc545b604051908152602001610146565b3480156101aa57600080fd5b5061016f6101b93660046114c1565b6104c0565b3480156101ca57600080fd5b5060cb5460405160ff9091168152602001610146565b3480156101ec57600080fd5b506102006101fb3660046114fd565b6104e4565b005b34801561020e57600080fd5b5061020061021d366004611497565b6105cd565b34801561022e57600080fd5b5061020061023d366004611518565b6105e3565b6102006102503660046115bd565b6105ed565b34801561026157600080fd5b506101906106ba565b34801561027657600080fd5b506101906102853660046114fd565b6001600160a01b0316600090815260cd602052604090205490565b3480156102ac57600080fd5b5061020061076d565b3480156102c157600080fd5b506102006102d0366004611497565b610781565b3480156102e157600080fd5b506097546040516001600160a01b039091168152602001610146565b34801561030957600080fd5b50610139610796565b34801561031e57600080fd5b5061020061032d366004611497565b6107a5565b34801561033e57600080fd5b5061016f61034d366004611497565b6107ad565b34801561035e57600080fd5b5061019061036d36600461161f565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103a457600080fd5b506102006103b3366004611672565b6107c3565b3480156103c457600080fd5b506102006103d33660046114fd565b610932565b606060c980546103e7906116ff565b80601f0160208091040260200160405190810160405280929190818152602001828054610413906116ff565b80156104605780601f1061043557610100808354040283529160200191610460565b820191906000526020600020905b81548152906001019060200180831161044357829003601f168201915b5050505050905090565b60006104773384846109a8565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6000336104ce858285610a2a565b6104d9858585610abc565b506001949350505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156105365760405162461bcd60e51b815260040161052d9061173a565b60405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661057f600080516020611898833981519152546001600160a01b031690565b6001600160a01b0316146105a55760405162461bcd60e51b815260040161052d90611786565b6105ae81610c6b565b604080516000808252602082019092526105ca91839190610c73565b50565b6105d5610de3565b6105df8282610e3d565b5050565b6105ca3382610f1c565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106365760405162461bcd60e51b815260040161052d9061173a565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661067f600080516020611898833981519152546001600160a01b031690565b6001600160a01b0316146106a55760405162461bcd60e51b815260040161052d90611786565b6106ae82610c6b565b6105df82826001610c73565b6000306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461075a5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c0000000000000000606482015260840161052d565b5060008051602061189883398151915290565b610775610de3565b61077f600061105e565b565b61078c823383610a2a565b6105df8282610f1c565b606060ca80546103e7906116ff565b61078c610de3565b60006107ba338484610abc565b50600192915050565b600054610100900460ff16158080156107e35750600054600160ff909116105b806107fd5750303b1580156107fd575060005460ff166001145b6108605760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b606482015260840161052d565b6000805460ff191660011790558015610883576000805461ff0019166101001790555b84516108969060c9906020880190611383565b5083516108aa9060ca906020870190611383565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b0384161790556108dd6110b0565b6108e56110df565b801561092b576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b61093a610de3565b6001600160a01b03811661099f5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161052d565b6105ca8161105e565b6001600160a01b0383166109fe5760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f2061646472657373000000604482015260640161052d565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b03808416600090815260ce60209081526040808320938616835292905220546000198114610ab65781811015610aa95760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161052d565b610ab684848484036109a8565b50505050565b6001600160a01b038316610b125760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f20616464726573730000604482015260640161052d565b6001600160a01b038216610b685760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f206164647265737300000000604482015260640161052d565b6001600160a01b038316600090815260cd602052604090205481811015610bd15760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e636500604482015260640161052d565b610bdb82826117e8565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610c119084906117ff565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610c5d91815260200190565b60405180910390a350505050565b6105ca610de3565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610cab57610ca683611106565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610d05575060408051601f3d908101601f19168201909252610d0291810190611817565b60015b610d685760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b606482015260840161052d565b6000805160206118988339815191528114610dd75760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b606482015260840161052d565b50610ca68383836111a2565b6097546001600160a01b0316331461077f5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161052d565b6001600160a01b038216610e935760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f20616464726573730000000000000000604482015260640161052d565b8060cc6000828254610ea591906117ff565b90915550506001600160a01b038216600090815260cd602052604081208054839290610ed29084906117ff565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6001600160a01b038216610f725760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f2061646472657373000000000000604482015260640161052d565b6001600160a01b038216600090815260cd602052604090205481811015610fdb5760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e63650000000000604482015260640161052d565b610fe582826117e8565b6001600160a01b038416600090815260cd602052604081209190915560cc80548492906110139084906117e8565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166110d75760405162461bcd60e51b815260040161052d90611830565b61077f6111c7565b600054610100900460ff1661077f5760405162461bcd60e51b815260040161052d90611830565b6001600160a01b0381163b6111735760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840161052d565b60008051602061189883398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6111ab836111f7565b6000825111806111b85750805b15610ca657610ab68383611237565b600054610100900460ff166111ee5760405162461bcd60e51b815260040161052d90611830565b61077f3361105e565b61120081611106565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061125c83836040518060600160405280602781526020016118b860279139611263565b9392505050565b6060600080856001600160a01b031685604051611280919061187b565b600060405180830381855af49150503d80600081146112bb576040519150601f19603f3d011682016040523d82523d6000602084013e6112c0565b606091505b50915091506112d1868383876112db565b9695505050505050565b60608315611347578251611340576001600160a01b0385163b6113405760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161052d565b5081611351565b6113518383611359565b949350505050565b8151156113695781518083602001fd5b8060405162461bcd60e51b815260040161052d9190611448565b82805461138f906116ff565b90600052602060002090601f0160209004810192826113b157600085556113f7565b82601f106113ca57805160ff19168380011785556113f7565b828001600101855582156113f7579182015b828111156113f75782518255916020019190600101906113dc565b50611403929150611407565b5090565b5b808211156114035760008155600101611408565b60005b8381101561143757818101518382015260200161141f565b83811115610ab65750506000910152565b602081526000825180602084015261146781604085016020870161141c565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461149257600080fd5b919050565b600080604083850312156114aa57600080fd5b6114b38361147b565b946020939093013593505050565b6000806000606084860312156114d657600080fd5b6114df8461147b565b92506114ed6020850161147b565b9150604084013590509250925092565b60006020828403121561150f57600080fd5b61125c8261147b565b60006020828403121561152a57600080fd5b5035919050565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff8084111561156257611562611531565b604051601f8501601f19908116603f0116810190828211818310171561158a5761158a611531565b816040528093508581528686860111156115a357600080fd5b858560208301376000602087830101525050509392505050565b600080604083850312156115d057600080fd5b6115d98361147b565b9150602083013567ffffffffffffffff8111156115f557600080fd5b8301601f8101851361160657600080fd5b61161585823560208401611547565b9150509250929050565b6000806040838503121561163257600080fd5b61163b8361147b565b91506116496020840161147b565b90509250929050565b600082601f83011261166357600080fd5b61125c83833560208501611547565b6000806000806080858703121561168857600080fd5b843567ffffffffffffffff808211156116a057600080fd5b6116ac88838901611652565b955060208701359150808211156116c257600080fd5b506116cf87828801611652565b935050604085013560ff811681146116e657600080fd5b91506116f46060860161147b565b905092959194509250565b600181811c9082168061171357607f821691505b6020821081141561173457634e487b7160e01b600052602260045260246000fd5b50919050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b6000828210156117fa576117fa6117d2565b500390565b60008219821115611812576118126117d2565b500190565b60006020828403121561182957600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b6000825161188d81846020870161141c565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212207d16f75d563e9ff3a3d12a47c196b0e2726b12d3b5790c5820241066679d31a464736f6c634300080a0033",
}

// ERC20UpgradableABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20UpgradableMetaData.ABI instead.
var ERC20UpgradableABI = ERC20UpgradableMetaData.ABI

// ERC20UpgradableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC20UpgradableMetaData.Bin instead.
var ERC20UpgradableBin = ERC20UpgradableMetaData.Bin

// DeployERC20Upgradable deploys a new Ethereum contract, binding an instance of ERC20Upgradable to it.
func DeployERC20Upgradable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20Upgradable, error) {
	parsed, err := ERC20UpgradableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC20UpgradableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Upgradable{ERC20UpgradableCaller: ERC20UpgradableCaller{contract: contract}, ERC20UpgradableTransactor: ERC20UpgradableTransactor{contract: contract}, ERC20UpgradableFilterer: ERC20UpgradableFilterer{contract: contract}}, nil
}

// ERC20Upgradable is an auto generated Go binding around an Ethereum contract.
type ERC20Upgradable struct {
	ERC20UpgradableCaller     // Read-only binding to the contract
	ERC20UpgradableTransactor // Write-only binding to the contract
	ERC20UpgradableFilterer   // Log filterer for contract events
}

// ERC20UpgradableCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20UpgradableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20UpgradableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20UpgradableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20UpgradableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20UpgradableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20UpgradableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20UpgradableSession struct {
	Contract     *ERC20Upgradable  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20UpgradableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20UpgradableCallerSession struct {
	Contract *ERC20UpgradableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ERC20UpgradableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20UpgradableTransactorSession struct {
	Contract     *ERC20UpgradableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ERC20UpgradableRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20UpgradableRaw struct {
	Contract *ERC20Upgradable // Generic contract binding to access the raw methods on
}

// ERC20UpgradableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20UpgradableCallerRaw struct {
	Contract *ERC20UpgradableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20UpgradableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20UpgradableTransactorRaw struct {
	Contract *ERC20UpgradableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Upgradable creates a new instance of ERC20Upgradable, bound to a specific deployed contract.
func NewERC20Upgradable(address common.Address, backend bind.ContractBackend) (*ERC20Upgradable, error) {
	contract, err := bindERC20Upgradable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Upgradable{ERC20UpgradableCaller: ERC20UpgradableCaller{contract: contract}, ERC20UpgradableTransactor: ERC20UpgradableTransactor{contract: contract}, ERC20UpgradableFilterer: ERC20UpgradableFilterer{contract: contract}}, nil
}

// NewERC20UpgradableCaller creates a new read-only instance of ERC20Upgradable, bound to a specific deployed contract.
func NewERC20UpgradableCaller(address common.Address, caller bind.ContractCaller) (*ERC20UpgradableCaller, error) {
	contract, err := bindERC20Upgradable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableCaller{contract: contract}, nil
}

// NewERC20UpgradableTransactor creates a new write-only instance of ERC20Upgradable, bound to a specific deployed contract.
func NewERC20UpgradableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20UpgradableTransactor, error) {
	contract, err := bindERC20Upgradable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableTransactor{contract: contract}, nil
}

// NewERC20UpgradableFilterer creates a new log filterer instance of ERC20Upgradable, bound to a specific deployed contract.
func NewERC20UpgradableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20UpgradableFilterer, error) {
	contract, err := bindERC20Upgradable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableFilterer{contract: contract}, nil
}

// bindERC20Upgradable binds a generic wrapper to an already deployed contract.
func bindERC20Upgradable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC20UpgradableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Upgradable *ERC20UpgradableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Upgradable.Contract.ERC20UpgradableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Upgradable *ERC20UpgradableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.ERC20UpgradableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Upgradable *ERC20UpgradableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.ERC20UpgradableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Upgradable *ERC20UpgradableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Upgradable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Upgradable *ERC20UpgradableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Upgradable *ERC20UpgradableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Upgradable.Contract.Allowance(&_ERC20Upgradable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Upgradable.Contract.Allowance(&_ERC20Upgradable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Upgradable.Contract.BalanceOf(&_ERC20Upgradable.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Upgradable.Contract.BalanceOf(&_ERC20Upgradable.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Upgradable *ERC20UpgradableCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Upgradable *ERC20UpgradableSession) Decimals() (uint8, error) {
	return _ERC20Upgradable.Contract.Decimals(&_ERC20Upgradable.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) Decimals() (uint8, error) {
	return _ERC20Upgradable.Contract.Decimals(&_ERC20Upgradable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Upgradable *ERC20UpgradableCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Upgradable *ERC20UpgradableSession) Name() (string, error) {
	return _ERC20Upgradable.Contract.Name(&_ERC20Upgradable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) Name() (string, error) {
	return _ERC20Upgradable.Contract.Name(&_ERC20Upgradable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20Upgradable *ERC20UpgradableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20Upgradable *ERC20UpgradableSession) Owner() (common.Address, error) {
	return _ERC20Upgradable.Contract.Owner(&_ERC20Upgradable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) Owner() (common.Address, error) {
	return _ERC20Upgradable.Contract.Owner(&_ERC20Upgradable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ERC20Upgradable *ERC20UpgradableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ERC20Upgradable *ERC20UpgradableSession) ProxiableUUID() ([32]byte, error) {
	return _ERC20Upgradable.Contract.ProxiableUUID(&_ERC20Upgradable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ERC20Upgradable.Contract.ProxiableUUID(&_ERC20Upgradable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Upgradable *ERC20UpgradableCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Upgradable *ERC20UpgradableSession) Symbol() (string, error) {
	return _ERC20Upgradable.Contract.Symbol(&_ERC20Upgradable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) Symbol() (string, error) {
	return _ERC20Upgradable.Contract.Symbol(&_ERC20Upgradable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Upgradable.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableSession) TotalSupply() (*big.Int, error) {
	return _ERC20Upgradable.Contract.TotalSupply(&_ERC20Upgradable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Upgradable *ERC20UpgradableCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Upgradable.Contract.TotalSupply(&_ERC20Upgradable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Approve(&_ERC20Upgradable.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Approve(&_ERC20Upgradable.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Burn(&_ERC20Upgradable.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Burn(&_ERC20Upgradable.TransactOpts, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "burn0", account, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Burn0(&_ERC20Upgradable.TransactOpts, account, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Burn0(&_ERC20Upgradable.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.BurnFrom(&_ERC20Upgradable.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.BurnFrom(&_ERC20Upgradable.TransactOpts, account, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "initialize", name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Initialize(&_ERC20Upgradable.TransactOpts, name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Initialize(&_ERC20Upgradable.TransactOpts, name_, symbol_, decimals_, module_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Mint(&_ERC20Upgradable.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Mint(&_ERC20Upgradable.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20Upgradable *ERC20UpgradableSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.RenounceOwnership(&_ERC20Upgradable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.RenounceOwnership(&_ERC20Upgradable.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Transfer(&_ERC20Upgradable.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.Transfer(&_ERC20Upgradable.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.TransferFrom(&_ERC20Upgradable.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.TransferFrom(&_ERC20Upgradable.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.TransferOwnership(&_ERC20Upgradable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.TransferOwnership(&_ERC20Upgradable.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ERC20Upgradable *ERC20UpgradableSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.UpgradeTo(&_ERC20Upgradable.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.UpgradeTo(&_ERC20Upgradable.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ERC20Upgradable *ERC20UpgradableTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ERC20Upgradable.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ERC20Upgradable *ERC20UpgradableSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.UpgradeToAndCall(&_ERC20Upgradable.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ERC20Upgradable *ERC20UpgradableTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ERC20Upgradable.Contract.UpgradeToAndCall(&_ERC20Upgradable.TransactOpts, newImplementation, data)
}

// ERC20UpgradableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC20Upgradable contract.
type ERC20UpgradableAdminChangedIterator struct {
	Event *ERC20UpgradableAdminChanged // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableAdminChanged)
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
		it.Event = new(ERC20UpgradableAdminChanged)
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
func (it *ERC20UpgradableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableAdminChanged represents a AdminChanged event raised by the ERC20Upgradable contract.
type ERC20UpgradableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC20UpgradableAdminChangedIterator, error) {

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableAdminChangedIterator{contract: _ERC20Upgradable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableAdminChanged)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseAdminChanged(log types.Log) (*ERC20UpgradableAdminChanged, error) {
	event := new(ERC20UpgradableAdminChanged)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20Upgradable contract.
type ERC20UpgradableApprovalIterator struct {
	Event *ERC20UpgradableApproval // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableApproval)
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
		it.Event = new(ERC20UpgradableApproval)
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
func (it *ERC20UpgradableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableApproval represents a Approval event raised by the ERC20Upgradable contract.
type ERC20UpgradableApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20UpgradableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableApprovalIterator{contract: _ERC20Upgradable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableApproval)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseApproval(log types.Log) (*ERC20UpgradableApproval, error) {
	event := new(ERC20UpgradableApproval)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC20Upgradable contract.
type ERC20UpgradableBeaconUpgradedIterator struct {
	Event *ERC20UpgradableBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableBeaconUpgraded)
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
		it.Event = new(ERC20UpgradableBeaconUpgraded)
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
func (it *ERC20UpgradableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableBeaconUpgraded represents a BeaconUpgraded event raised by the ERC20Upgradable contract.
type ERC20UpgradableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC20UpgradableBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableBeaconUpgradedIterator{contract: _ERC20Upgradable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableBeaconUpgraded)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseBeaconUpgraded(log types.Log) (*ERC20UpgradableBeaconUpgraded, error) {
	event := new(ERC20UpgradableBeaconUpgraded)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ERC20Upgradable contract.
type ERC20UpgradableInitializedIterator struct {
	Event *ERC20UpgradableInitialized // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableInitialized)
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
		it.Event = new(ERC20UpgradableInitialized)
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
func (it *ERC20UpgradableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableInitialized represents a Initialized event raised by the ERC20Upgradable contract.
type ERC20UpgradableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ERC20UpgradableInitializedIterator, error) {

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableInitializedIterator{contract: _ERC20Upgradable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableInitialized) (event.Subscription, error) {

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableInitialized)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseInitialized(log types.Log) (*ERC20UpgradableInitialized, error) {
	event := new(ERC20UpgradableInitialized)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ERC20Upgradable contract.
type ERC20UpgradableOwnershipTransferredIterator struct {
	Event *ERC20UpgradableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableOwnershipTransferred)
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
		it.Event = new(ERC20UpgradableOwnershipTransferred)
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
func (it *ERC20UpgradableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableOwnershipTransferred represents a OwnershipTransferred event raised by the ERC20Upgradable contract.
type ERC20UpgradableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ERC20UpgradableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableOwnershipTransferredIterator{contract: _ERC20Upgradable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableOwnershipTransferred)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseOwnershipTransferred(log types.Log) (*ERC20UpgradableOwnershipTransferred, error) {
	event := new(ERC20UpgradableOwnershipTransferred)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Upgradable contract.
type ERC20UpgradableTransferIterator struct {
	Event *ERC20UpgradableTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableTransfer)
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
		it.Event = new(ERC20UpgradableTransfer)
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
func (it *ERC20UpgradableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableTransfer represents a Transfer event raised by the ERC20Upgradable contract.
type ERC20UpgradableTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20UpgradableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableTransferIterator{contract: _ERC20Upgradable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableTransfer)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseTransfer(log types.Log) (*ERC20UpgradableTransfer, error) {
	event := new(ERC20UpgradableTransfer)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC20Upgradable contract.
type ERC20UpgradableUpgradedIterator struct {
	Event *ERC20UpgradableUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradableUpgraded)
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
		it.Event = new(ERC20UpgradableUpgraded)
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
func (it *ERC20UpgradableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradableUpgraded represents a Upgraded event raised by the ERC20Upgradable contract.
type ERC20UpgradableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC20Upgradable *ERC20UpgradableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC20UpgradableUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradableUpgradedIterator{contract: _ERC20Upgradable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC20Upgradable *ERC20UpgradableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC20UpgradableUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC20Upgradable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradableUpgraded)
				if err := _ERC20Upgradable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ERC20Upgradable *ERC20UpgradableFilterer) ParseUpgraded(log types.Log) (*ERC20UpgradableUpgraded, error) {
	event := new(ERC20UpgradableUpgraded)
	if err := _ERC20Upgradable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
