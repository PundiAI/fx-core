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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a060405261100260805234801561001657600080fd5b50608051611af761004e6000396000818161064901528181610692015281816107510152818161079101526108200152611af76000f3fe60806040526004361061014f5760003560e01c8063715018a6116100b6578063b86d52981161006f578063b86d5298146103bc578063d0e30db01461015e578063dd62ed3e146103da578063de7ea79d14610420578063f2fde38b14610440578063f3fef3a3146104605761015e565b8063715018a61461030057806379cc6790146103155780638da5cb5b1461033557806395d89b41146103675780639dc29fac1461037c578063a9059cbb1461039c5761015e565b80633659cfe6116101085780633659cfe61461024257806340c10f191461026257806342966c68146102825780634f1ef286146102a257806352d1902d146102b557806370a08231146102ca5761015e565b806306fdde0314610166578063095ea7b31461019157806318160ddd146101c157806323b872dd146101e05780632e1a7d4d14610200578063313ce567146102205761015e565b3661015e5761015c610480565b005b61015c610480565b34801561017257600080fd5b5061017b6104c1565b604051610188919061161f565b60405180910390f35b34801561019d57600080fd5b506101b16101ac366004611667565b610553565b6040519015158152602001610188565b3480156101cd57600080fd5b5060cc545b604051908152602001610188565b3480156101ec57600080fd5b506101b16101fb366004611693565b6105a9565b34801561020c57600080fd5b5061015c61021b3660046116d4565b6105cd565b34801561022c57600080fd5b5060cb5460405160ff9091168152602001610188565b34801561024e57600080fd5b5061015c61025d3660046116ed565b61063e565b34801561026e57600080fd5b5061015c61027d366004611667565b610727565b34801561028e57600080fd5b5061015c61029d3660046116d4565b61073d565b61015c6102b0366004611796565b610746565b3480156102c157600080fd5b506101d2610813565b3480156102d657600080fd5b506101d26102e53660046116ed565b6001600160a01b0316600090815260cd602052604090205490565b34801561030c57600080fd5b5061015c6108c6565b34801561032157600080fd5b5061015c610330366004611667565b6108da565b34801561034157600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610188565b34801561037357600080fd5b5061017b6108ef565b34801561038857600080fd5b5061015c610397366004611667565b6108fe565b3480156103a857600080fd5b506101b16103b7366004611667565b610906565b3480156103c857600080fd5b5060cf546001600160a01b031661034f565b3480156103e657600080fd5b506101d26103f53660046117fa565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b34801561042c57600080fd5b5061015c61043b366004611853565b61091c565b34801561044c57600080fd5b5061015c61045b3660046116ed565b610a8b565b34801561046c57600080fd5b5061015c61047b366004611667565b610b01565b61048a3334610b86565b60405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b606060c980546104d0906118e2565b80601f01602080910402602001604051908101604052809291908181526020018280546104fc906118e2565b80156105495780601f1061051e57610100808354040283529160200191610549565b820191906000526020600020905b81548152906001019060200180831161052c57829003601f168201915b5050505050905090565b6000610560338484610c5e565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6000336105b7858285610ce0565b6105c2858585610d72565b506001949350505050565b6105d8335b82610f21565b604051339082156108fc029083906000818181858888f19350505050158015610605573d6000803e3d6000fd5b5060405181815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a250565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106905760405162461bcd60e51b81526004016106879061191d565b60405180910390fd5b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106d9600080516020611a7b833981519152546001600160a01b031690565b6001600160a01b0316146106ff5760405162461bcd60e51b815260040161068790611969565b61070881611063565b604080516000808252602082019092526107249183919061106b565b50565b61072f6111db565b6107398282610b86565b5050565b610724336105d2565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016141561078f5760405162461bcd60e51b81526004016106879061191d565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166107d8600080516020611a7b833981519152546001600160a01b031690565b6001600160a01b0316146107fe5760405162461bcd60e51b815260040161068790611969565b61080782611063565b6107398282600161106b565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146108b35760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610687565b50600080516020611a7b83398151915290565b6108ce6111db565b6108d86000611235565b565b6108e5823383610ce0565b6107398282610f21565b606060ca80546104d0906118e2565b6108e56111db565b6000610913338484610d72565b50600192915050565b600054610100900460ff161580801561093c5750600054600160ff909116105b806109565750303b158015610956575060005460ff166001145b6109b95760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610687565b6000805460ff1916600117905580156109dc576000805461ff0019166101001790555b84516109ef9060c990602088019061155a565b508351610a039060ca90602087019061155a565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610a36611287565b610a3e6112b6565b8015610a84576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b610a936111db565b6001600160a01b038116610af85760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610687565b61072481611235565b610b0a336105d2565b6040516001600160a01b0383169082156108fc029083906000818181858888f19350505050158015610b40573d6000803e3d6000fd5b506040518181526001600160a01b0383169033907f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb906020015b60405180910390a35050565b6001600160a01b038216610bdc5760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f206164647265737300000000000000006044820152606401610687565b8060cc6000828254610bee91906119cb565b90915550506001600160a01b038216600090815260cd602052604081208054839290610c1b9084906119cb565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610b7a565b6001600160a01b038316610cb45760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f20616464726573730000006044820152606401610687565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b03808416600090815260ce60209081526040808320938616835292905220546000198114610d6c5781811015610d5f5760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e63650000006044820152606401610687565b610d6c8484848403610c5e565b50505050565b6001600160a01b038316610dc85760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f206164647265737300006044820152606401610687565b6001600160a01b038216610e1e5760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f2061646472657373000000006044820152606401610687565b6001600160a01b038316600090815260cd602052604090205481811015610e875760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e6365006044820152606401610687565b610e9182826119e3565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610ec79084906119cb565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610f1391815260200190565b60405180910390a350505050565b6001600160a01b038216610f775760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f20616464726573730000000000006044820152606401610687565b6001600160a01b038216600090815260cd602052604090205481811015610fe05760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e636500000000006044820152606401610687565b610fea82826119e3565b6001600160a01b038416600090815260cd602052604081209190915560cc80548492906110189084906119e3565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6107246111db565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156110a35761109e836112dd565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa9250505080156110fd575060408051601f3d908101601f191682019092526110fa918101906119fa565b60015b6111605760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610687565b600080516020611a7b83398151915281146111cf5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610687565b5061109e838383611379565b6097546001600160a01b031633146108d85760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610687565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166112ae5760405162461bcd60e51b815260040161068790611a13565b6108d861139e565b600054610100900460ff166108d85760405162461bcd60e51b815260040161068790611a13565b6001600160a01b0381163b61134a5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610687565b600080516020611a7b83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611382836113ce565b60008251118061138f5750805b1561109e57610d6c838361140e565b600054610100900460ff166113c55760405162461bcd60e51b815260040161068790611a13565b6108d833611235565b6113d7816112dd565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606114338383604051806060016040528060278152602001611a9b6027913961143a565b9392505050565b6060600080856001600160a01b0316856040516114579190611a5e565b600060405180830381855af49150503d8060008114611492576040519150601f19603f3d011682016040523d82523d6000602084013e611497565b606091505b50915091506114a8868383876114b2565b9695505050505050565b6060831561151e578251611517576001600160a01b0385163b6115175760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610687565b5081611528565b6115288383611530565b949350505050565b8151156115405781518083602001fd5b8060405162461bcd60e51b8152600401610687919061161f565b828054611566906118e2565b90600052602060002090601f01602090048101928261158857600085556115ce565b82601f106115a157805160ff19168380011785556115ce565b828001600101855582156115ce579182015b828111156115ce5782518255916020019190600101906115b3565b506115da9291506115de565b5090565b5b808211156115da57600081556001016115df565b60005b8381101561160e5781810151838201526020016115f6565b83811115610d6c5750506000910152565b602081526000825180602084015261163e8160408501602087016115f3565b601f01601f19169190910160400192915050565b6001600160a01b038116811461072457600080fd5b6000806040838503121561167a57600080fd5b823561168581611652565b946020939093013593505050565b6000806000606084860312156116a857600080fd5b83356116b381611652565b925060208401356116c381611652565b929592945050506040919091013590565b6000602082840312156116e657600080fd5b5035919050565b6000602082840312156116ff57600080fd5b813561143381611652565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff8084111561173b5761173b61170a565b604051601f8501601f19908116603f011681019082821181831017156117635761176361170a565b8160405280935085815286868601111561177c57600080fd5b858560208301376000602087830101525050509392505050565b600080604083850312156117a957600080fd5b82356117b481611652565b9150602083013567ffffffffffffffff8111156117d057600080fd5b8301601f810185136117e157600080fd5b6117f085823560208401611720565b9150509250929050565b6000806040838503121561180d57600080fd5b823561181881611652565b9150602083013561182881611652565b809150509250929050565b600082601f83011261184457600080fd5b61143383833560208501611720565b6000806000806080858703121561186957600080fd5b843567ffffffffffffffff8082111561188157600080fd5b61188d88838901611833565b955060208701359150808211156118a357600080fd5b506118b087828801611833565b935050604085013560ff811681146118c757600080fd5b915060608501356118d781611652565b939692955090935050565b600181811c908216806118f657607f821691505b6020821081141561191757634e487b7160e01b600052602260045260246000fd5b50919050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b600082198211156119de576119de6119b5565b500190565b6000828210156119f5576119f56119b5565b500390565b600060208284031215611a0c57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611a708184602087016115f3565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122028fec7582786fbf176e212f9db246ec5eff003a785551a5f26e94f7f1e8e3e3864736f6c634300080a0033",
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

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Burn(&_WarpTokenUpgradable.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Burn(&_WarpTokenUpgradable.TransactOpts, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) Burn0(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "burn0", account, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Burn0(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// Burn0 is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) Burn0(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.Burn0(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.BurnFrom(&_WarpTokenUpgradable.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_WarpTokenUpgradable *WarpTokenUpgradableTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WarpTokenUpgradable.Contract.BurnFrom(&_WarpTokenUpgradable.TransactOpts, account, amount)
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
