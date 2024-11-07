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

// BridgeFeeOracleMetaData contains all meta data concerning the BridgeFeeOracle contract.
var BridgeFeeOracleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OWNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"QUOTE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"blackOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crosschainContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracleList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_crosschain\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"isOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"oracleStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isBlacklisted\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_defaultOracle\",\"type\":\"address\"}],\"name\":\"setDefaultOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611c7f61004c6000396000818161054601528181610586015281816106260152818161066601526107070152611c7f6000f3fe60806040526004361061012a5760003560e01c806380dce169116100ab578063c4d66de81161006f578063c4d66de814610363578063d10c106114610383578063d547741f146103a3578063e58378bb146103c3578063e863f6a7146103e5578063ec331b2b1461043657600080fd5b806380dce169146102a157806391d14854146102da578063a217fddf146102fa578063b908afa81461030f578063c44014d21461034357600080fd5b80634f1ef286116100f25780634f1ef28614610204578063510c27ad1461021757806352d1902d146102395780635bca74db1461024e5780637c90c9a91461028157600080fd5b806301ffc9a71461012f578063248a9ca3146101645780632f2ff15d146101a257806336568abe146101c45780633659cfe6146101e4575b600080fd5b34801561013b57600080fd5b5061014f61014a3660046116bd565b610457565b60405190151581526020015b60405180910390f35b34801561017057600080fd5b5061019461017f3660046116e7565b600090815260c9602052604090206001015490565b60405190815260200161015b565b3480156101ae57600080fd5b506101c26101bd36600461171c565b61048e565b005b3480156101d057600080fd5b506101c26101df36600461171c565b6104b8565b3480156101f057600080fd5b506101c26101ff366004611748565b61053b565b6101c26102123660046117ef565b61061b565b34801561022357600080fd5b5061022c6106e8565b60405161015b9190611851565b34801561024557600080fd5b506101946106fa565b34801561025a57600080fd5b506101947e0caaa0e08f624de190c2474175cd13784c8c75bbdd1b63ae5fab5540967b3c81565b34801561028d57600080fd5b506101c261029c366004611748565b6107ad565b3480156102ad57600080fd5b5061012e546102c2906001600160a01b031681565b6040516001600160a01b03909116815260200161015b565b3480156102e657600080fd5b5061014f6102f536600461171c565b61087a565b34801561030657600080fd5b50610194600081565b34801561031b57600080fd5b506101947f88aa719609f728b0c5e7fb8dd3608d5c25d497efbb3b9dd64e9251ebba10150881565b34801561034f57600080fd5b506101c261035e366004611748565b6108a5565b34801561036f57600080fd5b506101c261037e366004611748565b610951565b34801561038f57600080fd5b5061014f61039e36600461189e565b610adb565b3480156103af57600080fd5b506101c26103be36600461171c565b610ce0565b3480156103cf57600080fd5b50610194600080516020611be383398151915281565b3480156103f157600080fd5b5061041f610400366004611748565b6101316020526000908152604090205460ff8082169161010090041682565b60408051921515835290151560208301520161015b565b34801561044257600080fd5b5061012d546102c2906001600160a01b031681565b60006001600160e01b03198216637965db0b60e01b148061048857506301ffc9a760e01b6001600160e01b03198316145b92915050565b600082815260c960205260409020600101546104a981610d05565b6104b38383610d0f565b505050565b6001600160a01b038116331461052d5760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b60648201526084015b60405180910390fd5b6105378282610d95565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156105845760405162461bcd60e51b8152600401610524906118f7565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166105cd600080516020611c03833981519152546001600160a01b031690565b6001600160a01b0316146105f35760405162461bcd60e51b815260040161052490611943565b6105fc81610dfc565b6040805160008082526020820190925261061891839190610e26565b50565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106645760405162461bcd60e51b8152600401610524906118f7565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106ad600080516020611c03833981519152546001600160a01b031690565b6001600160a01b0316146106d35760405162461bcd60e51b815260040161052490611943565b6106dc82610dfc565b61053782826001610e26565b60606106f561012f610f91565b905090565b6000306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461079a5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610524565b50600080516020611c0383398151915290565b600080516020611be38339815191526107c581610d05565b6107cd610fa5565b6001600160a01b0382166000908152610131602052604090205460ff16156107f457610870565b6001600160a01b03821660009081526101316020526040902054610100900460ff161561084b576001600160a01b038216600090815261013160205260409020805461ff001916905561084961012f83610fff565b505b6001600160a01b038216600090815261013160205260409020805460ff191660011790555b610537600160fb55565b600091825260c9602090815260408084206001600160a01b0393909316845291905290205460ff1690565b600080516020611be38339815191526108bd81610d05565b6108c961012f8361101b565b61092d576040805180820182526000808252600160208084019182526001600160a01b0387168352610131905292902090518154925161ffff1990931690151561ff001916176101009215159290920291909117905561092b61012f8361103d565b505b5061012e80546001600160a01b0319166001600160a01b0392909216919091179055565b600054610100900460ff16158080156109715750600054600160ff909116105b8061098b5750303b15801561098b575060005460ff166001145b6109ee5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610524565b6000805460ff191660011790558015610a11576000805461ff0019166101001790555b61012d80546001600160a01b0319166001600160a01b038416179055610a38600033610d0f565b610a627f88aa719609f728b0c5e7fb8dd3608d5c25d497efbb3b9dd64e9251ebba10150833610d0f565b610a7a600080516020611be383398151915233610d0f565b610a82611052565b610a8a611052565b610a9261107b565b8015610537576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15050565b60007e0caaa0e08f624de190c2474175cd13784c8c75bbdd1b63ae5fab5540967b3c610b0681610d05565b610b0e610fa5565b6001600160a01b03831660009081526101316020526040902054610100900460ff1615610b3e5760019150610ccf565b6001600160a01b0383166000908152610131602052604090205460ff1615610b695760009150610ccf565b61012d546040516333e7eceb60e11b81526001600160a01b03909116906367cfd9d690610b9c90879087906004016119e7565b602060405180830381865afa158015610bb9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bdd9190611a11565b610bea5760009150610ccf565b61012d54604051630b63ae7d60e11b81526001600160a01b03909116906316c75cfa90610c1d90879087906004016119e7565b602060405180830381865afa158015610c3a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c5e9190611a11565b610c6b5760009150610ccf565b6040805180820182526000808252600160208084019182526001600160a01b0388168352610131905292902090518154925161ffff1990931690151561ff0019161761010092151592909202919091179055610cc961012f8461103d565b50600191505b610cd9600160fb55565b5092915050565b600082815260c96020526040902060010154610cfb81610d05565b6104b38383610d95565b61061881336110aa565b610d19828261087a565b61053757600082815260c9602090815260408083206001600160a01b03851684529091529020805460ff19166001179055610d513390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b610d9f828261087a565b1561053757600082815260c9602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b7f88aa719609f728b0c5e7fb8dd3608d5c25d497efbb3b9dd64e9251ebba10150861053781610d05565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610e59576104b383611103565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610eb3575060408051601f3d908101601f19168201909252610eb091810190611a33565b60015b610f165760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610524565b600080516020611c038339815191528114610f855760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610524565b506104b383838361119f565b60606000610f9e836111ca565b9392505050565b600260fb541415610ff85760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610524565b600260fb55565b6000610f9e836001600160a01b038416611226565b600160fb55565b6001600160a01b03811660009081526001830160205260408120541515610f9e565b6000610f9e836001600160a01b038416611319565b600054610100900460ff166110795760405162461bcd60e51b815260040161052490611a4c565b565b600054610100900460ff166110a25760405162461bcd60e51b815260040161052490611a4c565b611079611368565b6110b4828261087a565b610537576110c18161138f565b6110cc8360206113a1565b6040516020016110dd929190611a97565b60408051601f198184030181529082905262461bcd60e51b825261052491600401611b0c565b6001600160a01b0381163b6111705760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610524565b600080516020611c0383398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b6111a88361153d565b6000825111806111b55750805b156104b3576111c4838361157d565b50505050565b60608160000180548060200260200160405190810160405280929190818152602001828054801561121a57602002820191906000526020600020905b815481526020019060010190808311611206575b50505050509050919050565b6000818152600183016020526040812054801561130f57600061124a600183611b35565b855490915060009061125e90600190611b35565b90508181146112c357600086600001828154811061127e5761127e611b4c565b90600052602060002001549050808760000184815481106112a1576112a1611b4c565b6000918252602080832090910192909255918252600188019052604090208390555b85548690806112d4576112d4611b62565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610488565b6000915050610488565b600081815260018301602052604081205461136057508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610488565b506000610488565b600054610100900460ff166110145760405162461bcd60e51b815260040161052490611a4c565b60606104886001600160a01b03831660145b606060006113b0836002611b78565b6113bb906002611b97565b67ffffffffffffffff8111156113d3576113d3611763565b6040519080825280601f01601f1916602001820160405280156113fd576020820181803683370190505b509050600360fc1b8160008151811061141857611418611b4c565b60200101906001600160f81b031916908160001a905350600f60fb1b8160018151811061144757611447611b4c565b60200101906001600160f81b031916908160001a905350600061146b846002611b78565b611476906001611b97565b90505b60018111156114ee576f181899199a1a9b1b9c1cb0b131b232b360811b85600f16601081106114aa576114aa611b4c565b1a60f81b8282815181106114c0576114c0611b4c565b60200101906001600160f81b031916908160001a90535060049490941c936114e781611baf565b9050611479565b508315610f9e5760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152606401610524565b61154681611103565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060610f9e8383604051806060016040528060278152602001611c23602791396060600080856001600160a01b0316856040516115ba9190611bc6565b600060405180830381855af49150503d80600081146115f5576040519150601f19603f3d011682016040523d82523d6000602084013e6115fa565b606091505b509150915061160b86838387611615565b9695505050505050565b6060831561168157825161167a576001600160a01b0385163b61167a5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610524565b508161168b565b61168b8383611693565b949350505050565b8151156116a35781518083602001fd5b8060405162461bcd60e51b81526004016105249190611b0c565b6000602082840312156116cf57600080fd5b81356001600160e01b031981168114610f9e57600080fd5b6000602082840312156116f957600080fd5b5035919050565b80356001600160a01b038116811461171757600080fd5b919050565b6000806040838503121561172f57600080fd5b8235915061173f60208401611700565b90509250929050565b60006020828403121561175a57600080fd5b610f9e82611700565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff8084111561179457611794611763565b604051601f8501601f19908116603f011681019082821181831017156117bc576117bc611763565b816040528093508581528686860111156117d557600080fd5b858560208301376000602087830101525050509392505050565b6000806040838503121561180257600080fd5b61180b83611700565b9150602083013567ffffffffffffffff81111561182757600080fd5b8301601f8101851361183857600080fd5b61184785823560208401611779565b9150509250929050565b6020808252825182820181905260009190848201906040850190845b818110156118925783516001600160a01b03168352928401929184019160010161186d565b50909695505050505050565b600080604083850312156118b157600080fd5b823567ffffffffffffffff8111156118c857600080fd5b8301601f810185136118d957600080fd5b6118e885823560208401611779565b92505061173f60208401611700565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b60005b838110156119aa578181015183820152602001611992565b838111156111c45750506000910152565b600081518084526119d381602086016020860161198f565b601f01601f19169290920160200192915050565b6040815260006119fa60408301856119bb565b905060018060a01b03831660208301529392505050565b600060208284031215611a2357600080fd5b81518015158114610f9e57600080fd5b600060208284031215611a4557600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000815260008351611acf81601785016020880161198f565b7001034b99036b4b9b9b4b733903937b6329607d1b6017918401918201528351611b0081602884016020880161198f565b01602801949350505050565b602081526000610f9e60208301846119bb565b634e487b7160e01b600052601160045260246000fd5b600082821015611b4757611b47611b1f565b500390565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052603160045260246000fd5b6000816000190483118215151615611b9257611b92611b1f565b500290565b60008219821115611baa57611baa611b1f565b500190565b600081611bbe57611bbe611b1f565b506000190190565b60008251611bd881846020870161198f565b919091019291505056feb19546dff01e856fb3f010c267a7b1c60363cf8a4664e21cc89c26224620214e360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220ae8e57c4c6d7c195977c203a30f1be8e88223ce5ae767b698b15cfaf890080b464736f6c634300080a0033",
}

// BridgeFeeOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeFeeOracleMetaData.ABI instead.
var BridgeFeeOracleABI = BridgeFeeOracleMetaData.ABI

// BridgeFeeOracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BridgeFeeOracleMetaData.Bin instead.
var BridgeFeeOracleBin = BridgeFeeOracleMetaData.Bin

// DeployBridgeFeeOracle deploys a new Ethereum contract, binding an instance of BridgeFeeOracle to it.
func DeployBridgeFeeOracle(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BridgeFeeOracle, error) {
	parsed, err := BridgeFeeOracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BridgeFeeOracleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BridgeFeeOracle{BridgeFeeOracleCaller: BridgeFeeOracleCaller{contract: contract}, BridgeFeeOracleTransactor: BridgeFeeOracleTransactor{contract: contract}, BridgeFeeOracleFilterer: BridgeFeeOracleFilterer{contract: contract}}, nil
}

// BridgeFeeOracle is an auto generated Go binding around an Ethereum contract.
type BridgeFeeOracle struct {
	BridgeFeeOracleCaller     // Read-only binding to the contract
	BridgeFeeOracleTransactor // Write-only binding to the contract
	BridgeFeeOracleFilterer   // Log filterer for contract events
}

// BridgeFeeOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeFeeOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFeeOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeFeeOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFeeOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFeeOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFeeOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeFeeOracleSession struct {
	Contract     *BridgeFeeOracle  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeFeeOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeFeeOracleCallerSession struct {
	Contract *BridgeFeeOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// BridgeFeeOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeFeeOracleTransactorSession struct {
	Contract     *BridgeFeeOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// BridgeFeeOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeFeeOracleRaw struct {
	Contract *BridgeFeeOracle // Generic contract binding to access the raw methods on
}

// BridgeFeeOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeFeeOracleCallerRaw struct {
	Contract *BridgeFeeOracleCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeFeeOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeFeeOracleTransactorRaw struct {
	Contract *BridgeFeeOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeFeeOracle creates a new instance of BridgeFeeOracle, bound to a specific deployed contract.
func NewBridgeFeeOracle(address common.Address, backend bind.ContractBackend) (*BridgeFeeOracle, error) {
	contract, err := bindBridgeFeeOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracle{BridgeFeeOracleCaller: BridgeFeeOracleCaller{contract: contract}, BridgeFeeOracleTransactor: BridgeFeeOracleTransactor{contract: contract}, BridgeFeeOracleFilterer: BridgeFeeOracleFilterer{contract: contract}}, nil
}

// NewBridgeFeeOracleCaller creates a new read-only instance of BridgeFeeOracle, bound to a specific deployed contract.
func NewBridgeFeeOracleCaller(address common.Address, caller bind.ContractCaller) (*BridgeFeeOracleCaller, error) {
	contract, err := bindBridgeFeeOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleCaller{contract: contract}, nil
}

// NewBridgeFeeOracleTransactor creates a new write-only instance of BridgeFeeOracle, bound to a specific deployed contract.
func NewBridgeFeeOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeFeeOracleTransactor, error) {
	contract, err := bindBridgeFeeOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleTransactor{contract: contract}, nil
}

// NewBridgeFeeOracleFilterer creates a new log filterer instance of BridgeFeeOracle, bound to a specific deployed contract.
func NewBridgeFeeOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFeeOracleFilterer, error) {
	contract, err := bindBridgeFeeOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleFilterer{contract: contract}, nil
}

// bindBridgeFeeOracle binds a generic wrapper to an already deployed contract.
func bindBridgeFeeOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeFeeOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeFeeOracle *BridgeFeeOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeFeeOracle.Contract.BridgeFeeOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeFeeOracle *BridgeFeeOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.BridgeFeeOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeFeeOracle *BridgeFeeOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.BridgeFeeOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeFeeOracle *BridgeFeeOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeFeeOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeFeeOracle *BridgeFeeOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeFeeOracle *BridgeFeeOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.DEFAULTADMINROLE(&_BridgeFeeOracle.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.DEFAULTADMINROLE(&_BridgeFeeOracle.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) OWNERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "OWNER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleSession) OWNERROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.OWNERROLE(&_BridgeFeeOracle.CallOpts)
}

// OWNERROLE is a free data retrieval call binding the contract method 0xe58378bb.
//
// Solidity: function OWNER_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) OWNERROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.OWNERROLE(&_BridgeFeeOracle.CallOpts)
}

// QUOTEROLE is a free data retrieval call binding the contract method 0x5bca74db.
//
// Solidity: function QUOTE_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) QUOTEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "QUOTE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// QUOTEROLE is a free data retrieval call binding the contract method 0x5bca74db.
//
// Solidity: function QUOTE_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleSession) QUOTEROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.QUOTEROLE(&_BridgeFeeOracle.CallOpts)
}

// QUOTEROLE is a free data retrieval call binding the contract method 0x5bca74db.
//
// Solidity: function QUOTE_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) QUOTEROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.QUOTEROLE(&_BridgeFeeOracle.CallOpts)
}

// UPGRADEROLE is a free data retrieval call binding the contract method 0xb908afa8.
//
// Solidity: function UPGRADE_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) UPGRADEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "UPGRADE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UPGRADEROLE is a free data retrieval call binding the contract method 0xb908afa8.
//
// Solidity: function UPGRADE_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleSession) UPGRADEROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.UPGRADEROLE(&_BridgeFeeOracle.CallOpts)
}

// UPGRADEROLE is a free data retrieval call binding the contract method 0xb908afa8.
//
// Solidity: function UPGRADE_ROLE() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) UPGRADEROLE() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.UPGRADEROLE(&_BridgeFeeOracle.CallOpts)
}

// CrosschainContract is a free data retrieval call binding the contract method 0xec331b2b.
//
// Solidity: function crosschainContract() view returns(address)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) CrosschainContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "crosschainContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CrosschainContract is a free data retrieval call binding the contract method 0xec331b2b.
//
// Solidity: function crosschainContract() view returns(address)
func (_BridgeFeeOracle *BridgeFeeOracleSession) CrosschainContract() (common.Address, error) {
	return _BridgeFeeOracle.Contract.CrosschainContract(&_BridgeFeeOracle.CallOpts)
}

// CrosschainContract is a free data retrieval call binding the contract method 0xec331b2b.
//
// Solidity: function crosschainContract() view returns(address)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) CrosschainContract() (common.Address, error) {
	return _BridgeFeeOracle.Contract.CrosschainContract(&_BridgeFeeOracle.CallOpts)
}

// DefaultOracle is a free data retrieval call binding the contract method 0x80dce169.
//
// Solidity: function defaultOracle() view returns(address)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) DefaultOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "defaultOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultOracle is a free data retrieval call binding the contract method 0x80dce169.
//
// Solidity: function defaultOracle() view returns(address)
func (_BridgeFeeOracle *BridgeFeeOracleSession) DefaultOracle() (common.Address, error) {
	return _BridgeFeeOracle.Contract.DefaultOracle(&_BridgeFeeOracle.CallOpts)
}

// DefaultOracle is a free data retrieval call binding the contract method 0x80dce169.
//
// Solidity: function defaultOracle() view returns(address)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) DefaultOracle() (common.Address, error) {
	return _BridgeFeeOracle.Contract.DefaultOracle(&_BridgeFeeOracle.CallOpts)
}

// GetOracleList is a free data retrieval call binding the contract method 0x510c27ad.
//
// Solidity: function getOracleList() view returns(address[])
func (_BridgeFeeOracle *BridgeFeeOracleCaller) GetOracleList(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "getOracleList")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOracleList is a free data retrieval call binding the contract method 0x510c27ad.
//
// Solidity: function getOracleList() view returns(address[])
func (_BridgeFeeOracle *BridgeFeeOracleSession) GetOracleList() ([]common.Address, error) {
	return _BridgeFeeOracle.Contract.GetOracleList(&_BridgeFeeOracle.CallOpts)
}

// GetOracleList is a free data retrieval call binding the contract method 0x510c27ad.
//
// Solidity: function getOracleList() view returns(address[])
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) GetOracleList() ([]common.Address, error) {
	return _BridgeFeeOracle.Contract.GetOracleList(&_BridgeFeeOracle.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BridgeFeeOracle.Contract.GetRoleAdmin(&_BridgeFeeOracle.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BridgeFeeOracle.Contract.GetRoleAdmin(&_BridgeFeeOracle.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BridgeFeeOracle.Contract.HasRole(&_BridgeFeeOracle.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BridgeFeeOracle.Contract.HasRole(&_BridgeFeeOracle.CallOpts, role, account)
}

// OracleStatus is a free data retrieval call binding the contract method 0xe863f6a7.
//
// Solidity: function oracleStatus(address ) view returns(bool isBlacklisted, bool isActive)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) OracleStatus(opts *bind.CallOpts, arg0 common.Address) (struct {
	IsBlacklisted bool
	IsActive      bool
}, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "oracleStatus", arg0)

	outstruct := new(struct {
		IsBlacklisted bool
		IsActive      bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsBlacklisted = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.IsActive = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// OracleStatus is a free data retrieval call binding the contract method 0xe863f6a7.
//
// Solidity: function oracleStatus(address ) view returns(bool isBlacklisted, bool isActive)
func (_BridgeFeeOracle *BridgeFeeOracleSession) OracleStatus(arg0 common.Address) (struct {
	IsBlacklisted bool
	IsActive      bool
}, error) {
	return _BridgeFeeOracle.Contract.OracleStatus(&_BridgeFeeOracle.CallOpts, arg0)
}

// OracleStatus is a free data retrieval call binding the contract method 0xe863f6a7.
//
// Solidity: function oracleStatus(address ) view returns(bool isBlacklisted, bool isActive)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) OracleStatus(arg0 common.Address) (struct {
	IsBlacklisted bool
	IsActive      bool
}, error) {
	return _BridgeFeeOracle.Contract.OracleStatus(&_BridgeFeeOracle.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleSession) ProxiableUUID() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.ProxiableUUID(&_BridgeFeeOracle.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) ProxiableUUID() ([32]byte, error) {
	return _BridgeFeeOracle.Contract.ProxiableUUID(&_BridgeFeeOracle.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BridgeFeeOracle.Contract.SupportsInterface(&_BridgeFeeOracle.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BridgeFeeOracle.Contract.SupportsInterface(&_BridgeFeeOracle.CallOpts, interfaceId)
}

// BlackOracle is a paid mutator transaction binding the contract method 0x7c90c9a9.
//
// Solidity: function blackOracle(address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) BlackOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "blackOracle", _oracle)
}

// BlackOracle is a paid mutator transaction binding the contract method 0x7c90c9a9.
//
// Solidity: function blackOracle(address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) BlackOracle(_oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.BlackOracle(&_BridgeFeeOracle.TransactOpts, _oracle)
}

// BlackOracle is a paid mutator transaction binding the contract method 0x7c90c9a9.
//
// Solidity: function blackOracle(address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) BlackOracle(_oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.BlackOracle(&_BridgeFeeOracle.TransactOpts, _oracle)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.GrantRole(&_BridgeFeeOracle.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.GrantRole(&_BridgeFeeOracle.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _crosschain) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) Initialize(opts *bind.TransactOpts, _crosschain common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "initialize", _crosschain)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _crosschain) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) Initialize(_crosschain common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.Initialize(&_BridgeFeeOracle.TransactOpts, _crosschain)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _crosschain) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) Initialize(_crosschain common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.Initialize(&_BridgeFeeOracle.TransactOpts, _crosschain)
}

// IsOnline is a paid mutator transaction binding the contract method 0xd10c1061.
//
// Solidity: function isOnline(string _chainName, address _oracle) returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) IsOnline(opts *bind.TransactOpts, _chainName string, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "isOnline", _chainName, _oracle)
}

// IsOnline is a paid mutator transaction binding the contract method 0xd10c1061.
//
// Solidity: function isOnline(string _chainName, address _oracle) returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleSession) IsOnline(_chainName string, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.IsOnline(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
}

// IsOnline is a paid mutator transaction binding the contract method 0xd10c1061.
//
// Solidity: function isOnline(string _chainName, address _oracle) returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) IsOnline(_chainName string, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.IsOnline(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.RenounceRole(&_BridgeFeeOracle.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.RenounceRole(&_BridgeFeeOracle.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.RevokeRole(&_BridgeFeeOracle.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.RevokeRole(&_BridgeFeeOracle.TransactOpts, role, account)
}

// SetDefaultOracle is a paid mutator transaction binding the contract method 0xc44014d2.
//
// Solidity: function setDefaultOracle(address _defaultOracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) SetDefaultOracle(opts *bind.TransactOpts, _defaultOracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "setDefaultOracle", _defaultOracle)
}

// SetDefaultOracle is a paid mutator transaction binding the contract method 0xc44014d2.
//
// Solidity: function setDefaultOracle(address _defaultOracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) SetDefaultOracle(_defaultOracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.SetDefaultOracle(&_BridgeFeeOracle.TransactOpts, _defaultOracle)
}

// SetDefaultOracle is a paid mutator transaction binding the contract method 0xc44014d2.
//
// Solidity: function setDefaultOracle(address _defaultOracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) SetDefaultOracle(_defaultOracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.SetDefaultOracle(&_BridgeFeeOracle.TransactOpts, _defaultOracle)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.UpgradeTo(&_BridgeFeeOracle.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.UpgradeTo(&_BridgeFeeOracle.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.UpgradeToAndCall(&_BridgeFeeOracle.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.UpgradeToAndCall(&_BridgeFeeOracle.TransactOpts, newImplementation, data)
}

// BridgeFeeOracleAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleAdminChangedIterator struct {
	Event *BridgeFeeOracleAdminChanged // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleAdminChanged)
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
		it.Event = new(BridgeFeeOracleAdminChanged)
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
func (it *BridgeFeeOracleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleAdminChanged represents a AdminChanged event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*BridgeFeeOracleAdminChangedIterator, error) {

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleAdminChangedIterator{contract: _BridgeFeeOracle.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleAdminChanged) (event.Subscription, error) {

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleAdminChanged)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseAdminChanged(log types.Log) (*BridgeFeeOracleAdminChanged, error) {
	event := new(BridgeFeeOracleAdminChanged)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeFeeOracleBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleBeaconUpgradedIterator struct {
	Event *BridgeFeeOracleBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleBeaconUpgraded)
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
		it.Event = new(BridgeFeeOracleBeaconUpgraded)
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
func (it *BridgeFeeOracleBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleBeaconUpgraded represents a BeaconUpgraded event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*BridgeFeeOracleBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleBeaconUpgradedIterator{contract: _BridgeFeeOracle.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleBeaconUpgraded)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseBeaconUpgraded(log types.Log) (*BridgeFeeOracleBeaconUpgraded, error) {
	event := new(BridgeFeeOracleBeaconUpgraded)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeFeeOracleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleInitializedIterator struct {
	Event *BridgeFeeOracleInitialized // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleInitialized)
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
		it.Event = new(BridgeFeeOracleInitialized)
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
func (it *BridgeFeeOracleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleInitialized represents a Initialized event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterInitialized(opts *bind.FilterOpts) (*BridgeFeeOracleInitializedIterator, error) {

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleInitializedIterator{contract: _BridgeFeeOracle.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleInitialized) (event.Subscription, error) {

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleInitialized)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseInitialized(log types.Log) (*BridgeFeeOracleInitialized, error) {
	event := new(BridgeFeeOracleInitialized)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeFeeOracleRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleRoleAdminChangedIterator struct {
	Event *BridgeFeeOracleRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleRoleAdminChanged)
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
		it.Event = new(BridgeFeeOracleRoleAdminChanged)
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
func (it *BridgeFeeOracleRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleRoleAdminChanged represents a RoleAdminChanged event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BridgeFeeOracleRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleRoleAdminChangedIterator{contract: _BridgeFeeOracle.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleRoleAdminChanged)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseRoleAdminChanged(log types.Log) (*BridgeFeeOracleRoleAdminChanged, error) {
	event := new(BridgeFeeOracleRoleAdminChanged)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeFeeOracleRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleRoleGrantedIterator struct {
	Event *BridgeFeeOracleRoleGranted // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleRoleGranted)
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
		it.Event = new(BridgeFeeOracleRoleGranted)
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
func (it *BridgeFeeOracleRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleRoleGranted represents a RoleGranted event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BridgeFeeOracleRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleRoleGrantedIterator{contract: _BridgeFeeOracle.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleRoleGranted)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseRoleGranted(log types.Log) (*BridgeFeeOracleRoleGranted, error) {
	event := new(BridgeFeeOracleRoleGranted)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeFeeOracleRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleRoleRevokedIterator struct {
	Event *BridgeFeeOracleRoleRevoked // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleRoleRevoked)
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
		it.Event = new(BridgeFeeOracleRoleRevoked)
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
func (it *BridgeFeeOracleRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleRoleRevoked represents a RoleRevoked event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BridgeFeeOracleRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleRoleRevokedIterator{contract: _BridgeFeeOracle.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleRoleRevoked)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseRoleRevoked(log types.Log) (*BridgeFeeOracleRoleRevoked, error) {
	event := new(BridgeFeeOracleRoleRevoked)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeFeeOracleUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the BridgeFeeOracle contract.
type BridgeFeeOracleUpgradedIterator struct {
	Event *BridgeFeeOracleUpgraded // Event containing the contract specifics and raw log

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
func (it *BridgeFeeOracleUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeFeeOracleUpgraded)
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
		it.Event = new(BridgeFeeOracleUpgraded)
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
func (it *BridgeFeeOracleUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeFeeOracleUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeFeeOracleUpgraded represents a Upgraded event raised by the BridgeFeeOracle contract.
type BridgeFeeOracleUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BridgeFeeOracleUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BridgeFeeOracleUpgradedIterator{contract: _BridgeFeeOracle.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BridgeFeeOracleUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _BridgeFeeOracle.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeFeeOracleUpgraded)
				if err := _BridgeFeeOracle.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_BridgeFeeOracle *BridgeFeeOracleFilterer) ParseUpgraded(log types.Log) (*BridgeFeeOracleUpgraded, error) {
	event := new(BridgeFeeOracleUpgraded)
	if err := _BridgeFeeOracle.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
