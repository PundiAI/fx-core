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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"TransferCrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"transferCrossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a060405261100160805234801561001657600080fd5b50608051611bff61004e60003960008181610580015281816105c001528181610676015281816106b601526107450152611bff6000f3fe60806040526004361061011f5760003560e01c8063715018a6116100a0578063b86d529811610064578063b86d52981461031c578063c5cb9b511461033a578063dd62ed3e1461035a578063de7ea79d146103a0578063f2fde38b146103c057600080fd5b8063715018a6146102805780638da5cb5b1461029557806395d89b41146102c75780639dc29fac146102dc578063a9059cbb146102fc57600080fd5b80633659cfe6116100e75780633659cfe6146101e057806340c10f19146102025780634f1ef2861461022257806352d1902d1461023557806370a082311461024a57600080fd5b806306fdde0314610124578063095ea7b31461014f57806318160ddd1461017f57806323b872dd1461019e578063313ce567146101be575b600080fd5b34801561013057600080fd5b506101396103e0565b6040516101469190611677565b60405180910390f35b34801561015b57600080fd5b5061016f61016a3660046116a6565b610472565b6040519015158152602001610146565b34801561018b57600080fd5b5060cc545b604051908152602001610146565b3480156101aa57600080fd5b5061016f6101b93660046116d0565b6104c8565b3480156101ca57600080fd5b5060cb5460405160ff9091168152602001610146565b3480156101ec57600080fd5b506102006101fb36600461170c565b610575565b005b34801561020e57600080fd5b5061020061021d3660046116a6565b610655565b6102006102303660046117b3565b61066b565b34801561024157600080fd5b50610190610738565b34801561025657600080fd5b5061019061026536600461170c565b6001600160a01b0316600090815260cd602052604090205490565b34801561028c57600080fd5b506102006107eb565b3480156102a157600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610146565b3480156102d357600080fd5b506101396107ff565b3480156102e857600080fd5b506102006102f73660046116a6565b61080e565b34801561030857600080fd5b5061016f6103173660046116a6565b610820565b34801561032857600080fd5b5060cf546001600160a01b03166102af565b34801561034657600080fd5b5061016f610355366004611835565b610836565b34801561036657600080fd5b50610190610375366004611889565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103ac57600080fd5b506102006103bb3660046118bc565b6108ed565b3480156103cc57600080fd5b506102006103db36600461170c565b610a5c565b606060c980546103ef90611949565b80601f016020809104026020016040519081016040528092919081815260200182805461041b90611949565b80156104685780601f1061043d57610100808354040283529160200191610468565b820191906000526020600020905b81548152906001019060200180831161044b57829003601f168201915b5050505050905090565b600061047f338484610ad2565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce602090815260408083203384529091528120548281101561054b5760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b61055f853361055a868561199a565b610ad2565b61056a858585610b54565b506001949350505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156105be5760405162461bcd60e51b8152600401610542906119b1565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610607600080516020611b83833981519152546001600160a01b031690565b6001600160a01b03161461062d5760405162461bcd60e51b8152600401610542906119fd565b61063681610d03565b6040805160008082526020820190925261065291839190610d0b565b50565b61065d610e7b565b6106678282610ed5565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106b45760405162461bcd60e51b8152600401610542906119b1565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106fd600080516020611b83833981519152546001600160a01b031690565b6001600160a01b0316146107235760405162461bcd60e51b8152600401610542906119fd565b61072c82610d03565b61066782826001610d0b565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610542565b50600080516020611b8383398151915290565b6107f3610e7b565b6107fd6000610fb4565b565b606060ca80546103ef90611949565b610816610e7b565b6106678282611006565b600061082d338484610b54565b50600192915050565b600063ffffffff333b161561088d5760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e7472616374000000000000006044820152606401610542565b61089a3386868686611148565b336001600160a01b03167f282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d868686866040516108d99493929190611a49565b60405180910390a25060015b949350505050565b600054610100900460ff161580801561090d5750600054600160ff909116105b806109275750303b158015610927575060005460ff166001145b61098a5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610542565b6000805460ff1916600117905580156109ad576000805461ff0019166101001790555b84516109c09060c9906020880190611586565b5083516109d49060ca906020870190611586565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610a076112ba565b610a0f6112e9565b8015610a55576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b610a64610e7b565b6001600160a01b038116610ac95760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610542565b61065281610fb4565b6001600160a01b038316610b285760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f20616464726573730000006044820152606401610542565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610baa5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f206164647265737300006044820152606401610542565b6001600160a01b038216610c005760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f2061646472657373000000006044820152606401610542565b6001600160a01b038316600090815260cd602052604090205481811015610c695760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e6365006044820152606401610542565b610c73828261199a565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610ca9908490611a78565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610cf591815260200190565b60405180910390a350505050565b610652610e7b565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610d4357610d3e83611310565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610d9d575060408051601f3d908101601f19168201909252610d9a91810190611a90565b60015b610e005760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610542565b600080516020611b838339815191528114610e6f5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610542565b50610d3e8383836113ac565b6097546001600160a01b031633146107fd5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610542565b6001600160a01b038216610f2b5760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f206164647265737300000000000000006044820152606401610542565b8060cc6000828254610f3d9190611a78565b90915550506001600160a01b038216600090815260cd602052604081208054839290610f6a908490611a78565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b03821661105c5760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f20616464726573730000000000006044820152606401610542565b6001600160a01b038216600090815260cd6020526040902054818110156110c55760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e636500000000006044820152606401610542565b6110cf828261199a565b6001600160a01b038416600090815260cd602052604081209190915560cc80548492906110fd90849061199a565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b03851661119e5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f206164647265737300006044820152606401610542565b60008451116111e35760405162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a59081c9958da5c1a595b9d607a1b6044820152606401610542565b806112215760405162461bcd60e51b815260206004820152600e60248201526d1a5b9d985b1a59081d185c99d95d60921b6044820152606401610542565b60cf546112429086906001600160a01b031661123d8587611a78565b610b54565b604051633c3e7d7760e01b815261100490633c3e7d779061126f9088908890889088908890600401611aa9565b6020604051808303816000875af115801561128e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112b29190611af9565b505050505050565b600054610100900460ff166112e15760405162461bcd60e51b815260040161054290611b1b565b6107fd6113d7565b600054610100900460ff166107fd5760405162461bcd60e51b815260040161054290611b1b565b6001600160a01b0381163b61137d5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610542565b600080516020611b8383398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6113b583611407565b6000825111806113c25750805b15610d3e576113d18383611447565b50505050565b600054610100900460ff166113fe5760405162461bcd60e51b815260040161054290611b1b565b6107fd33610fb4565b61141081611310565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061146c8383604051806060016040528060278152602001611ba360279139611473565b9392505050565b6060600080856001600160a01b0316856040516114909190611b66565b600060405180830381855af49150503d80600081146114cb576040519150601f19603f3d011682016040523d82523d6000602084013e6114d0565b606091505b50915091506114e1868383876114eb565b9695505050505050565b60608315611557578251611550576001600160a01b0385163b6115505760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610542565b50816108e5565b6108e5838381511561156c5781518083602001fd5b8060405162461bcd60e51b81526004016105429190611677565b82805461159290611949565b90600052602060002090601f0160209004810192826115b457600085556115fa565b82601f106115cd57805160ff19168380011785556115fa565b828001600101855582156115fa579182015b828111156115fa5782518255916020019190600101906115df565b5061160692915061160a565b5090565b5b80821115611606576000815560010161160b565b60005b8381101561163a578181015183820152602001611622565b838111156113d15750506000910152565b6000815180845261166381602086016020860161161f565b601f01601f19169290920160200192915050565b60208152600061146c602083018461164b565b80356001600160a01b03811681146116a157600080fd5b919050565b600080604083850312156116b957600080fd5b6116c28361168a565b946020939093013593505050565b6000806000606084860312156116e557600080fd5b6116ee8461168a565b92506116fc6020850161168a565b9150604084013590509250925092565b60006020828403121561171e57600080fd5b61146c8261168a565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff8084111561175857611758611727565b604051601f8501601f19908116603f0116810190828211818310171561178057611780611727565b8160405280935085815286868601111561179957600080fd5b858560208301376000602087830101525050509392505050565b600080604083850312156117c657600080fd5b6117cf8361168a565b9150602083013567ffffffffffffffff8111156117eb57600080fd5b8301601f810185136117fc57600080fd5b61180b8582356020840161173d565b9150509250929050565b600082601f83011261182657600080fd5b61146c8383356020850161173d565b6000806000806080858703121561184b57600080fd5b843567ffffffffffffffff81111561186257600080fd5b61186e87828801611815565b97602087013597506040870135966060013595509350505050565b6000806040838503121561189c57600080fd5b6118a58361168a565b91506118b36020840161168a565b90509250929050565b600080600080608085870312156118d257600080fd5b843567ffffffffffffffff808211156118ea57600080fd5b6118f688838901611815565b9550602087013591508082111561190c57600080fd5b5061191987828801611815565b935050604085013560ff8116811461193057600080fd5b915061193e6060860161168a565b905092959194509250565b600181811c9082168061195d57607f821691505b6020821081141561197e57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156119ac576119ac611984565b500390565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b608081526000611a5c608083018761164b565b6020830195909552506040810192909252606090910152919050565b60008219821115611a8b57611a8b611984565b500190565b600060208284031215611aa257600080fd5b5051919050565b6001600160a01b038616815260c060208201819052600090611acd9083018761164b565b6040830195909552506060810192909252608082015280820360a0909101526000815260200192915050565b600060208284031215611b0b57600080fd5b8151801515811461146c57600080fd5b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611b7881846020870161161f565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a264697066735822122070d364539d1f58bd31880a125003eb6bc84f83994cac68e082b695570b528b5d64736f6c634300080a0033",
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

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_FIP20Upgradable *FIP20UpgradableCaller) Module(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FIP20Upgradable.contract.Call(opts, &out, "module")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_FIP20Upgradable *FIP20UpgradableSession) Module() (common.Address, error) {
	return _FIP20Upgradable.Contract.Module(&_FIP20Upgradable.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_FIP20Upgradable *FIP20UpgradableCallerSession) Module() (common.Address, error) {
	return _FIP20Upgradable.Contract.Module(&_FIP20Upgradable.CallOpts)
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

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactor) TransferCrossChain(opts *bind.TransactOpts, recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _FIP20Upgradable.contract.Transact(opts, "transferCrossChain", recipient, amount, fee, target)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableSession) TransferCrossChain(recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.TransferCrossChain(&_FIP20Upgradable.TransactOpts, recipient, amount, fee, target)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) returns(bool)
func (_FIP20Upgradable *FIP20UpgradableTransactorSession) TransferCrossChain(recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _FIP20Upgradable.Contract.TransferCrossChain(&_FIP20Upgradable.TransactOpts, recipient, amount, fee, target)
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

// FIP20UpgradableTransferCrossChainIterator is returned from FilterTransferCrossChain and is used to iterate over the raw logs and unpacked data for TransferCrossChain events raised by the FIP20Upgradable contract.
type FIP20UpgradableTransferCrossChainIterator struct {
	Event *FIP20UpgradableTransferCrossChain // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradableTransferCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20UpgradableTransferCrossChain)
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
		it.Event = new(FIP20UpgradableTransferCrossChain)
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
func (it *FIP20UpgradableTransferCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradableTransferCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20UpgradableTransferCrossChain represents a TransferCrossChain event raised by the FIP20Upgradable contract.
type FIP20UpgradableTransferCrossChain struct {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) FilterTransferCrossChain(opts *bind.FilterOpts, from []common.Address) (*FIP20UpgradableTransferCrossChainIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.FilterLogs(opts, "TransferCrossChain", fromRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradableTransferCrossChainIterator{contract: _FIP20Upgradable.contract, event: "TransferCrossChain", logs: logs, sub: sub}, nil
}

// WatchTransferCrossChain is a free log subscription operation binding the contract event 0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d.
//
// Solidity: event TransferCrossChain(address indexed from, string recipient, uint256 amount, uint256 fee, bytes32 target)
func (_FIP20Upgradable *FIP20UpgradableFilterer) WatchTransferCrossChain(opts *bind.WatchOpts, sink chan<- *FIP20UpgradableTransferCrossChain, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _FIP20Upgradable.contract.WatchLogs(opts, "TransferCrossChain", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20UpgradableTransferCrossChain)
				if err := _FIP20Upgradable.contract.UnpackLog(event, "TransferCrossChain", log); err != nil {
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
func (_FIP20Upgradable *FIP20UpgradableFilterer) ParseTransferCrossChain(log types.Log) (*FIP20UpgradableTransferCrossChain, error) {
	event := new(FIP20UpgradableTransferCrossChain)
	if err := _FIP20Upgradable.contract.UnpackLog(event, "TransferCrossChain", log); err != nil {
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
