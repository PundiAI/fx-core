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

// FIP20ABI is the input ABI used to generate the binding from.
const FIP20ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"TransferCrossChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"module_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"target\",\"type\":\"bytes32\"}],\"name\":\"transferCrossChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// FIP20Bin is the compiled bytecode used for deploying new contracts.
var FIP20Bin = "0x60a06040523060601b60805234801561001757600080fd5b5060805160601c611b4061005260003960008181610582015281816105c20152818161069a015281816106da01526107690152611b406000f3fe60806040526004361061011f5760003560e01c8063715018a6116100a0578063b86d529811610064578063b86d52981461031c578063c5cb9b511461033a578063dd62ed3e1461035a578063de7ea79d146103a0578063f2fde38b146103c057600080fd5b8063715018a6146102805780638da5cb5b1461029557806395d89b41146102c75780639dc29fac146102dc578063a9059cbb146102fc57600080fd5b80633659cfe6116100e75780633659cfe6146101e057806340c10f19146102025780634f1ef2861461022257806352d1902d1461023557806370a082311461024a57600080fd5b806306fdde0314610124578063095ea7b31461014f57806318160ddd1461017f57806323b872dd1461019e578063313ce567146101be575b600080fd5b34801561013057600080fd5b506101396103e0565b60405161014691906118a7565b60405180910390f35b34801561015b57600080fd5b5061016f61016a366004611743565b610472565b6040519015158152602001610146565b34801561018b57600080fd5b5060cc545b604051908152602001610146565b3480156101aa57600080fd5b5061016f6101b93660046116a9565b6104c8565b3480156101ca57600080fd5b5060cb5460405160ff9091168152602001610146565b3480156101ec57600080fd5b506102006101fb36600461165d565b610577565b005b34801561020e57600080fd5b5061020061021d366004611743565b610657565b6102006102303660046116e4565b61068f565b34801561024157600080fd5b5061019061075c565b34801561025657600080fd5b5061019061026536600461165d565b6001600160a01b0316600090815260cd602052604090205490565b34801561028c57600080fd5b5061020061080f565b3480156102a157600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610146565b3480156102d357600080fd5b50610139610845565b3480156102e857600080fd5b506102006102f7366004611743565b610854565b34801561030857600080fd5b5061016f610317366004611743565b610888565b34801561032857600080fd5b5060cf546001600160a01b03166102af565b34801561034657600080fd5b5061016f61035536600461180d565b61089e565b34801561036657600080fd5b50610190610375366004611677565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103ac57600080fd5b506102006103bb366004611784565b61090d565b3480156103cc57600080fd5b506102006103db36600461165d565b610a2c565b606060c980546103ef90611a5c565b80601f016020809104026020016040519081016040528092919081815260200182805461041b90611a5c565b80156104685780601f1061043d57610100808354040283529160200191610468565b820191906000526020600020905b81548152906001019060200180831161044b57829003601f168201915b5050505050905090565b600061047f338484610ac4565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce602090815260408083203384529091528120548281101561054b5760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b61055f853361055a8685611a19565b610ac4565b61056a858585610b46565b60019150505b9392505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156105c05760405162461bcd60e51b8152600401610542906118e9565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610609600080516020611ac4833981519152546001600160a01b031690565b6001600160a01b03161461062f5760405162461bcd60e51b815260040161054290611935565b61063881610cf5565b6040805160008082526020820190925261065491839190610d1f565b50565b6097546001600160a01b031633146106815760405162461bcd60e51b815260040161054290611981565b61068b8282610e9e565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106d85760405162461bcd60e51b8152600401610542906118e9565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610721600080516020611ac4833981519152546001600160a01b031690565b6001600160a01b0316146107475760405162461bcd60e51b815260040161054290611935565b61075082610cf5565b61068b82826001610d1f565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107fc5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610542565b50600080516020611ac483398151915290565b6097546001600160a01b031633146108395760405162461bcd60e51b815260040161054290611981565b6108436000610f7d565b565b606060ca80546103ef90611a5c565b6097546001600160a01b0316331461087e5760405162461bcd60e51b815260040161054290611981565b61068b8282610fcf565b6000610895338484610b46565b50600192915050565b600063ffffffff333b16156108f55760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e7472616374000000000000006044820152606401610542565b6109023386868686611111565b506001949350505050565b600054610100900460ff166109285760005460ff161561092c565b303b155b61098f5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610542565b600054610100900460ff161580156109b1576000805461ffff19166101011790555b84516109c49060c9906020880190611513565b5083516109d89060ca906020870190611513565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610a0b611259565b610a13611288565b8015610a25576000805461ff00191690555b5050505050565b6097546001600160a01b03163314610a565760405162461bcd60e51b815260040161054290611981565b6001600160a01b038116610abb5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610542565b61065481610f7d565b6001600160a01b038316610b1a5760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f20616464726573730000006044820152606401610542565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610b9c5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f206164647265737300006044820152606401610542565b6001600160a01b038216610bf25760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f2061646472657373000000006044820152606401610542565b6001600160a01b038316600090815260cd602052604090205481811015610c5b5760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e6365006044820152606401610542565b610c658282611a19565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610c9b908490611a01565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610ce791815260200190565b60405180910390a350505050565b6097546001600160a01b031633146106545760405162461bcd60e51b815260040161054290611981565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610d5757610d52836112af565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610d9057600080fd5b505afa925050508015610dc0575060408051601f3d908101601f19168201909252610dbd9181019061176c565b60015b610e235760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610542565b600080516020611ac48339815191528114610e925760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610542565b50610d5283838361134b565b6001600160a01b038216610ef45760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f206164647265737300000000000000006044820152606401610542565b8060cc6000828254610f069190611a01565b90915550506001600160a01b038216600090815260cd602052604081208054839290610f33908490611a01565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b0382166110255760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f20616464726573730000000000006044820152606401610542565b6001600160a01b038216600090815260cd60205260409020548181101561108e5760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e636500000000006044820152606401610542565b6110988282611a19565b6001600160a01b038416600090815260cd602052604081209190915560cc80548492906110c6908490611a19565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b0385166111675760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f206164647265737300006044820152606401610542565b60008451116111ac5760405162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a59081c9958da5c1a595b9d607a1b6044820152606401610542565b806111ea5760405162461bcd60e51b815260206004820152600e60248201526d1a5b9d985b1a59081d185c99d95d60921b6044820152606401610542565b60cf5461120b9086906001600160a01b03166112068587611a01565b610b46565b846001600160a01b03167f282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d8585858560405161124a94939291906118ba565b60405180910390a25050505050565b600054610100900460ff166112805760405162461bcd60e51b8152600401610542906119b6565b610843611376565b600054610100900460ff166108435760405162461bcd60e51b8152600401610542906119b6565b6001600160a01b0381163b61131c5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610542565b600080516020611ac483398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611354836113a6565b6000825111806113615750805b15610d525761137083836113e6565b50505050565b600054610100900460ff1661139d5760405162461bcd60e51b8152600401610542906119b6565b61084333610f7d565b6113af816112af565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606001600160a01b0383163b61144e5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a2064656c65676174652063616c6c20746f206e6f6e2d636f6044820152651b9d1c9858dd60d21b6064820152608401610542565b600080846001600160a01b031684604051611469919061188b565b600060405180830381855af49150503d80600081146114a4576040519150601f19603f3d011682016040523d82523d6000602084013e6114a9565b606091505b50915091506114d18282604051806060016040528060278152602001611ae4602791396114da565b95945050505050565b606083156114e9575081610570565b8251156114f95782518084602001fd5b8160405162461bcd60e51b815260040161054291906118a7565b82805461151f90611a5c565b90600052602060002090601f0160209004810192826115415760008555611587565b82601f1061155a57805160ff1916838001178555611587565b82800160010185558215611587579182015b8281111561158757825182559160200191906001019061156c565b50611593929150611597565b5090565b5b808211156115935760008155600101611598565b600067ffffffffffffffff808411156115c7576115c7611aad565b604051601f8501601f19908116603f011681019082821181831017156115ef576115ef611aad565b8160405280935085815286868601111561160857600080fd5b858560208301376000602087830101525050509392505050565b80356001600160a01b038116811461163957600080fd5b919050565b600082601f83011261164e578081fd5b610570838335602085016115ac565b60006020828403121561166e578081fd5b61057082611622565b60008060408385031215611689578081fd5b61169283611622565b91506116a060208401611622565b90509250929050565b6000806000606084860312156116bd578081fd5b6116c684611622565b92506116d460208501611622565b9150604084013590509250925092565b600080604083850312156116f6578182fd5b6116ff83611622565b9150602083013567ffffffffffffffff81111561171a578182fd5b8301601f8101851361172a578182fd5b611739858235602084016115ac565b9150509250929050565b60008060408385031215611755578182fd5b61175e83611622565b946020939093013593505050565b60006020828403121561177d578081fd5b5051919050565b60008060008060808587031215611799578081fd5b843567ffffffffffffffff808211156117b0578283fd5b6117bc8883890161163e565b955060208701359150808211156117d1578283fd5b506117de8782880161163e565b935050604085013560ff811681146117f4578182fd5b915061180260608601611622565b905092959194509250565b60008060008060808587031215611822578384fd5b843567ffffffffffffffff811115611838578485fd5b6118448782880161163e565b97602087013597506040870135966060013595509350505050565b60008151808452611877816020860160208601611a30565b601f01601f19169290920160200192915050565b6000825161189d818460208701611a30565b9190910192915050565b602081526000610570602083018461185f565b6080815260006118cd608083018761185f565b6020830195909552506040810192909252606090910152919050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008219821115611a1457611a14611a97565b500190565b600082821015611a2b57611a2b611a97565b500390565b60005b83811015611a4b578181015183820152602001611a33565b838111156113705750506000910152565b600181811c90821680611a7057607f821691505b60208210811415611a9157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fdfe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220fcaf066833d727d01e7acbf759899143615dfd717f285f84b2f7c2fd419211ed64736f6c63430008040033"

// DeployFIP20 deploys a new Ethereum contract, binding an instance of FIP20 to it.
func DeployFIP20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FIP20, error) {
	parsed, err := abi.JSON(strings.NewReader(FIP20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FIP20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FIP20{FIP20Caller: FIP20Caller{contract: contract}, FIP20Transactor: FIP20Transactor{contract: contract}, FIP20Filterer: FIP20Filterer{contract: contract}}, nil
}

// FIP20 is an auto generated Go binding around an Ethereum contract.
type FIP20 struct {
	FIP20Caller     // Read-only binding to the contract
	FIP20Transactor // Write-only binding to the contract
	FIP20Filterer   // Log filterer for contract events
}

// FIP20Caller is an auto generated read-only Go binding around an Ethereum contract.
type FIP20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIP20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type FIP20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIP20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FIP20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIP20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FIP20Session struct {
	Contract     *FIP20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FIP20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FIP20CallerSession struct {
	Contract *FIP20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FIP20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FIP20TransactorSession struct {
	Contract     *FIP20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FIP20Raw is an auto generated low-level Go binding around an Ethereum contract.
type FIP20Raw struct {
	Contract *FIP20 // Generic contract binding to access the raw methods on
}

// FIP20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FIP20CallerRaw struct {
	Contract *FIP20Caller // Generic read-only contract binding to access the raw methods on
}

// FIP20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FIP20TransactorRaw struct {
	Contract *FIP20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewFIP20 creates a new instance of FIP20, bound to a specific deployed contract.
func NewFIP20(address common.Address, backend bind.ContractBackend) (*FIP20, error) {
	contract, err := bindFIP20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FIP20{FIP20Caller: FIP20Caller{contract: contract}, FIP20Transactor: FIP20Transactor{contract: contract}, FIP20Filterer: FIP20Filterer{contract: contract}}, nil
}

// NewFIP20Caller creates a new read-only instance of FIP20, bound to a specific deployed contract.
func NewFIP20Caller(address common.Address, caller bind.ContractCaller) (*FIP20Caller, error) {
	contract, err := bindFIP20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FIP20Caller{contract: contract}, nil
}

// NewFIP20Transactor creates a new write-only instance of FIP20, bound to a specific deployed contract.
func NewFIP20Transactor(address common.Address, transactor bind.ContractTransactor) (*FIP20Transactor, error) {
	contract, err := bindFIP20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FIP20Transactor{contract: contract}, nil
}

// NewFIP20Filterer creates a new log filterer instance of FIP20, bound to a specific deployed contract.
func NewFIP20Filterer(address common.Address, filterer bind.ContractFilterer) (*FIP20Filterer, error) {
	contract, err := bindFIP20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FIP20Filterer{contract: contract}, nil
}

// bindFIP20 binds a generic wrapper to an already deployed contract.
func bindFIP20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FIP20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FIP20 *FIP20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FIP20.Contract.FIP20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FIP20 *FIP20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIP20.Contract.FIP20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FIP20 *FIP20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FIP20.Contract.FIP20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FIP20 *FIP20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FIP20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FIP20 *FIP20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIP20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FIP20 *FIP20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FIP20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FIP20 *FIP20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "allowance", owner, spender)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FIP20 *FIP20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FIP20.Contract.Allowance(&_FIP20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FIP20 *FIP20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FIP20.Contract.Allowance(&_FIP20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FIP20 *FIP20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "balanceOf", account)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FIP20 *FIP20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _FIP20.Contract.BalanceOf(&_FIP20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FIP20 *FIP20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FIP20.Contract.BalanceOf(&_FIP20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FIP20 *FIP20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FIP20 *FIP20Session) Decimals() (uint8, error) {
	return _FIP20.Contract.Decimals(&_FIP20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FIP20 *FIP20CallerSession) Decimals() (uint8, error) {
	return _FIP20.Contract.Decimals(&_FIP20.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_FIP20 *FIP20Caller) Module(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "module")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_FIP20 *FIP20Session) Module() (common.Address, error) {
	return _FIP20.Contract.Module(&_FIP20.CallOpts)
}

// Module is a free data retrieval call binding the contract method 0xb86d5298.
//
// Solidity: function module() view returns(address)
func (_FIP20 *FIP20CallerSession) Module() (common.Address, error) {
	return _FIP20.Contract.Module(&_FIP20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FIP20 *FIP20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "name")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FIP20 *FIP20Session) Name() (string, error) {
	return _FIP20.Contract.Name(&_FIP20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FIP20 *FIP20CallerSession) Name() (string, error) {
	return _FIP20.Contract.Name(&_FIP20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FIP20 *FIP20Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FIP20 *FIP20Session) Owner() (common.Address, error) {
	return _FIP20.Contract.Owner(&_FIP20.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FIP20 *FIP20CallerSession) Owner() (common.Address, error) {
	return _FIP20.Contract.Owner(&_FIP20.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FIP20 *FIP20Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FIP20 *FIP20Session) ProxiableUUID() ([32]byte, error) {
	return _FIP20.Contract.ProxiableUUID(&_FIP20.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FIP20 *FIP20CallerSession) ProxiableUUID() ([32]byte, error) {
	return _FIP20.Contract.ProxiableUUID(&_FIP20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FIP20 *FIP20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "symbol")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FIP20 *FIP20Session) Symbol() (string, error) {
	return _FIP20.Contract.Symbol(&_FIP20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FIP20 *FIP20CallerSession) Symbol() (string, error) {
	return _FIP20.Contract.Symbol(&_FIP20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FIP20 *FIP20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FIP20.contract.Call(opts, &out, "totalSupply")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FIP20 *FIP20Session) TotalSupply() (*big.Int, error) {
	return _FIP20.Contract.TotalSupply(&_FIP20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FIP20 *FIP20CallerSession) TotalSupply() (*big.Int, error) {
	return _FIP20.Contract.TotalSupply(&_FIP20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FIP20 *FIP20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FIP20 *FIP20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Approve(&_FIP20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FIP20 *FIP20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Approve(&_FIP20.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_FIP20 *FIP20Transactor) Burn(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "burn", account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_FIP20 *FIP20Session) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Burn(&_FIP20.TransactOpts, account, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address account, uint256 amount) returns()
func (_FIP20 *FIP20TransactorSession) Burn(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Burn(&_FIP20.TransactOpts, account, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_FIP20 *FIP20Transactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "initialize", name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_FIP20 *FIP20Session) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _FIP20.Contract.Initialize(&_FIP20.TransactOpts, name_, symbol_, decimals_, module_)
}

// Initialize is a paid mutator transaction binding the contract method 0xde7ea79d.
//
// Solidity: function initialize(string name_, string symbol_, uint8 decimals_, address module_) returns()
func (_FIP20 *FIP20TransactorSession) Initialize(name_ string, symbol_ string, decimals_ uint8, module_ common.Address) (*types.Transaction, error) {
	return _FIP20.Contract.Initialize(&_FIP20.TransactOpts, name_, symbol_, decimals_, module_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_FIP20 *FIP20Transactor) Mint(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "mint", account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_FIP20 *FIP20Session) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Mint(&_FIP20.TransactOpts, account, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address account, uint256 amount) returns()
func (_FIP20 *FIP20TransactorSession) Mint(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Mint(&_FIP20.TransactOpts, account, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FIP20 *FIP20Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FIP20 *FIP20Session) RenounceOwnership() (*types.Transaction, error) {
	return _FIP20.Contract.RenounceOwnership(&_FIP20.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FIP20 *FIP20TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FIP20.Contract.RenounceOwnership(&_FIP20.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_FIP20 *FIP20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_FIP20 *FIP20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Transfer(&_FIP20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_FIP20 *FIP20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.Transfer(&_FIP20.TransactOpts, recipient, amount)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) returns(bool)
func (_FIP20 *FIP20Transactor) TransferCrossChain(opts *bind.TransactOpts, recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "transferCrossChain", recipient, amount, fee, target)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) returns(bool)
func (_FIP20 *FIP20Session) TransferCrossChain(recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _FIP20.Contract.TransferCrossChain(&_FIP20.TransactOpts, recipient, amount, fee, target)
}

// TransferCrossChain is a paid mutator transaction binding the contract method 0xc5cb9b51.
//
// Solidity: function transferCrossChain(string recipient, uint256 amount, uint256 fee, bytes32 target) returns(bool)
func (_FIP20 *FIP20TransactorSession) TransferCrossChain(recipient string, amount *big.Int, fee *big.Int, target [32]byte) (*types.Transaction, error) {
	return _FIP20.Contract.TransferCrossChain(&_FIP20.TransactOpts, recipient, amount, fee, target)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_FIP20 *FIP20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_FIP20 *FIP20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.TransferFrom(&_FIP20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_FIP20 *FIP20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FIP20.Contract.TransferFrom(&_FIP20.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FIP20 *FIP20Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FIP20 *FIP20Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FIP20.Contract.TransferOwnership(&_FIP20.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FIP20 *FIP20TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FIP20.Contract.TransferOwnership(&_FIP20.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FIP20 *FIP20Transactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FIP20 *FIP20Session) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FIP20.Contract.UpgradeTo(&_FIP20.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_FIP20 *FIP20TransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _FIP20.Contract.UpgradeTo(&_FIP20.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FIP20 *FIP20Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FIP20.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FIP20 *FIP20Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FIP20.Contract.UpgradeToAndCall(&_FIP20.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FIP20 *FIP20TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FIP20.Contract.UpgradeToAndCall(&_FIP20.TransactOpts, newImplementation, data)
}

// FIP20AdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the FIP20 contract.
type FIP20AdminChangedIterator struct {
	Event *FIP20AdminChanged // Event containing the contract specifics and raw log

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
func (it *FIP20AdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20AdminChanged)
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
		it.Event = new(FIP20AdminChanged)
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
func (it *FIP20AdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20AdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20AdminChanged represents a AdminChanged event raised by the FIP20 contract.
type FIP20AdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FIP20 *FIP20Filterer) FilterAdminChanged(opts *bind.FilterOpts) (*FIP20AdminChangedIterator, error) {
	logs, sub, err := _FIP20.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &FIP20AdminChangedIterator{contract: _FIP20.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_FIP20 *FIP20Filterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *FIP20AdminChanged) (event.Subscription, error) {
	logs, sub, err := _FIP20.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20AdminChanged)
				if err := _FIP20.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseAdminChanged(log types.Log) (*FIP20AdminChanged, error) {
	event := new(FIP20AdminChanged)
	if err := _FIP20.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the FIP20 contract.
type FIP20ApprovalIterator struct {
	Event *FIP20Approval // Event containing the contract specifics and raw log

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
func (it *FIP20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20Approval)
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
		it.Event = new(FIP20Approval)
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
func (it *FIP20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20Approval represents a Approval event raised by the FIP20 contract.
type FIP20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FIP20 *FIP20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*FIP20ApprovalIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FIP20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &FIP20ApprovalIterator{contract: _FIP20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FIP20 *FIP20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *FIP20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FIP20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20Approval)
				if err := _FIP20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseApproval(log types.Log) (*FIP20Approval, error) {
	event := new(FIP20Approval)
	if err := _FIP20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20BeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the FIP20 contract.
type FIP20BeaconUpgradedIterator struct {
	Event *FIP20BeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *FIP20BeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20BeaconUpgraded)
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
		it.Event = new(FIP20BeaconUpgraded)
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
func (it *FIP20BeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20BeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20BeaconUpgraded represents a BeaconUpgraded event raised by the FIP20 contract.
type FIP20BeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FIP20 *FIP20Filterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*FIP20BeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FIP20.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &FIP20BeaconUpgradedIterator{contract: _FIP20.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_FIP20 *FIP20Filterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *FIP20BeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _FIP20.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20BeaconUpgraded)
				if err := _FIP20.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseBeaconUpgraded(log types.Log) (*FIP20BeaconUpgraded, error) {
	event := new(FIP20BeaconUpgraded)
	if err := _FIP20.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FIP20 contract.
type FIP20OwnershipTransferredIterator struct {
	Event *FIP20OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FIP20OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20OwnershipTransferred)
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
		it.Event = new(FIP20OwnershipTransferred)
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
func (it *FIP20OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20OwnershipTransferred represents a OwnershipTransferred event raised by the FIP20 contract.
type FIP20OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FIP20 *FIP20Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FIP20OwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FIP20.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FIP20OwnershipTransferredIterator{contract: _FIP20.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FIP20 *FIP20Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FIP20OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FIP20.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20OwnershipTransferred)
				if err := _FIP20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseOwnershipTransferred(log types.Log) (*FIP20OwnershipTransferred, error) {
	event := new(FIP20OwnershipTransferred)
	if err := _FIP20.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the FIP20 contract.
type FIP20TransferIterator struct {
	Event *FIP20Transfer // Event containing the contract specifics and raw log

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
func (it *FIP20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20Transfer)
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
		it.Event = new(FIP20Transfer)
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
func (it *FIP20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20Transfer represents a Transfer event raised by the FIP20 contract.
type FIP20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FIP20 *FIP20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FIP20TransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FIP20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FIP20TransferIterator{contract: _FIP20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FIP20 *FIP20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *FIP20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FIP20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20Transfer)
				if err := _FIP20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseTransfer(log types.Log) (*FIP20Transfer, error) {
	event := new(FIP20Transfer)
	if err := _FIP20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20TransferCrossChainIterator is returned from FilterTransferCrossChain and is used to iterate over the raw logs and unpacked data for TransferCrossChain events raised by the FIP20 contract.
type FIP20TransferCrossChainIterator struct {
	Event *FIP20TransferCrossChain // Event containing the contract specifics and raw log

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
func (it *FIP20TransferCrossChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20TransferCrossChain)
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
		it.Event = new(FIP20TransferCrossChain)
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
func (it *FIP20TransferCrossChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20TransferCrossChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20TransferCrossChain represents a TransferCrossChain event raised by the FIP20 contract.
type FIP20TransferCrossChain struct {
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
func (_FIP20 *FIP20Filterer) FilterTransferCrossChain(opts *bind.FilterOpts, from []common.Address) (*FIP20TransferCrossChainIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _FIP20.contract.FilterLogs(opts, "TransferCrossChain", fromRule)
	if err != nil {
		return nil, err
	}
	return &FIP20TransferCrossChainIterator{contract: _FIP20.contract, event: "TransferCrossChain", logs: logs, sub: sub}, nil
}

// WatchTransferCrossChain is a free log subscription operation binding the contract event 0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d.
//
// Solidity: event TransferCrossChain(address indexed from, string recipient, uint256 amount, uint256 fee, bytes32 target)
func (_FIP20 *FIP20Filterer) WatchTransferCrossChain(opts *bind.WatchOpts, sink chan<- *FIP20TransferCrossChain, from []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _FIP20.contract.WatchLogs(opts, "TransferCrossChain", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20TransferCrossChain)
				if err := _FIP20.contract.UnpackLog(event, "TransferCrossChain", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseTransferCrossChain(log types.Log) (*FIP20TransferCrossChain, error) {
	event := new(FIP20TransferCrossChain)
	if err := _FIP20.contract.UnpackLog(event, "TransferCrossChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FIP20UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FIP20 contract.
type FIP20UpgradedIterator struct {
	Event *FIP20Upgraded // Event containing the contract specifics and raw log

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
func (it *FIP20UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FIP20Upgraded)
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
		it.Event = new(FIP20Upgraded)
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
func (it *FIP20UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FIP20UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FIP20Upgraded represents a Upgraded event raised by the FIP20 contract.
type FIP20Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FIP20 *FIP20Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FIP20UpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FIP20.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FIP20UpgradedIterator{contract: _FIP20.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FIP20 *FIP20Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FIP20Upgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FIP20.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FIP20Upgraded)
				if err := _FIP20.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_FIP20 *FIP20Filterer) ParseUpgraded(log types.Log) (*FIP20Upgraded, error) {
	event := new(FIP20Upgraded)
	if err := _FIP20.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
