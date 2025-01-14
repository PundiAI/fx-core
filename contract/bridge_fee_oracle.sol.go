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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OWNER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"QUOTE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chainName\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"activeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chainName\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"blackOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_bytes32\",\"type\":\"bytes32\"}],\"name\":\"bytes32ToString\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"crosschainContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chainName\",\"type\":\"bytes32\"}],\"name\":\"getOracleList\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_crosschain\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chainName\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"isOnline\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"oracleStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isBlack\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_defaultOracle\",\"type\":\"address\"}],\"name\":\"setDefaultOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b50608051611e9961004c600039600081816106ab015281816106eb0152818161078b015281816107cb015261085a0152611e996000f3fe6080604052600436106101405760003560e01c806380dce169116100b6578063b908afa81161006f578063b908afa8146103f9578063c44014d21461042d578063c4d66de81461044d578063d547741f1461046d578063e58378bb1461048d578063ec331b2b146104af57600080fd5b806380dce169146102e257806391d148541461031b5780639201de551461033b578063a217fddf14610368578063af51ef101461037d578063b338123c146103d957600080fd5b80633659cfe6116101085780633659cfe61461021a5780634f1ef2861461023a57806352d1902d1461024d5780635bca74db146102625780635cbb51791461029557806361166581146102b557600080fd5b806301ffc9a7146101455780632237bdad1461017a578063248a9ca31461019c5780632f2ff15d146101da57806336568abe146101fa575b600080fd5b34801561015157600080fd5b50610165610160366004611926565b6104d0565b60405190151581526020015b60405180910390f35b34801561018657600080fd5b5061019a61019536600461196c565b610507565b005b3480156101a857600080fd5b506101cc6101b7366004611998565b600090815260c9602052604090206001015490565b604051908152602001610171565b3480156101e657600080fd5b5061019a6101f536600461196c565b6105f8565b34801561020657600080fd5b5061019a61021536600461196c565b61061d565b34801561022657600080fd5b5061019a6102353660046119b1565b6106a0565b61019a6102483660046119e2565b610780565b34801561025957600080fd5b506101cc61084d565b34801561026e57600080fd5b506101cc7e0caaa0e08f624de190c2474175cd13784c8c75bbdd1b63ae5fab5540967b3c81565b3480156102a157600080fd5b5061019a6102b036600461196c565b610900565b3480156102c157600080fd5b506102d56102d0366004611998565b6109c3565b6040516101719190611aa4565b3480156102ee57600080fd5b5061012e54610303906001600160a01b031681565b6040516001600160a01b039091168152602001610171565b34801561032757600080fd5b5061016561033636600461196c565b6109de565b34801561034757600080fd5b5061035b610356366004611998565b610a09565b6040516101719190611b49565b34801561037457600080fd5b506101cc600081565b34801561038957600080fd5b506103c261039836600461196c565b61013060209081526000928352604080842090915290825290205460ff8082169161010090041682565b604080519215158352901515602083015201610171565b3480156103e557600080fd5b506101656103f436600461196c565b610b40565b34801561040557600080fd5b506101cc7f88aa719609f728b0c5e7fb8dd3608d5c25d497efbb3b9dd64e9251ebba10150881565b34801561043957600080fd5b5061019a6104483660046119b1565b610db2565b34801561045957600080fd5b5061019a6104683660046119b1565b610dee565b34801561047957600080fd5b5061019a61048836600461196c565b610f78565b34801561049957600080fd5b506101cc600080516020611dfd83398151915281565b3480156104bb57600080fd5b5061012d54610303906001600160a01b031681565b60006001600160e01b03198216637965db0b60e01b148061050157506301ffc9a760e01b6001600160e01b03198316145b92915050565b600080516020611dfd83398151915261051f81610f9d565b6000838152610130602090815260408083206001600160a01b038616845290915290205460ff161561055057505050565b6000838152610130602090815260408083206001600160a01b0386168452909152902054610100900460ff16156105c4576000838152610130602090815260408083206001600160a01b03861684528252808320805461ff001916905585835261012f90915290206105c29083610fa7565b505b6000838152610130602090815260408083206001600160a01b03861684529091529020805460ff191660011790555b505050565b600082815260c9602052604090206001015461061381610f9d565b6105f38383610fbc565b6001600160a01b03811633146106925760405162461bcd60e51b815260206004820152602f60248201527f416363657373436f6e74726f6c3a2063616e206f6e6c792072656e6f756e636560448201526e103937b632b9903337b91039b2b63360891b60648201526084015b60405180910390fd5b61069c8282611042565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156106e95760405162461bcd60e51b815260040161068990611b5c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610732600080516020611e1d833981519152546001600160a01b031690565b6001600160a01b0316146107585760405162461bcd60e51b815260040161068990611ba8565b610761816110a9565b6040805160008082526020820190925261077d918391906110d3565b50565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156107c95760405162461bcd60e51b815260040161068990611b5c565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316610812600080516020611e1d833981519152546001600160a01b031690565b6001600160a01b0316146108385760405162461bcd60e51b815260040161068990611ba8565b610841826110a9565b61069c828260016110d3565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146108ed5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610689565b50600080516020611e1d83398151915290565b600080516020611dfd83398151915261091881610f9d565b6000838152610130602090815260408083206001600160a01b0386168452909152902054610100900460ff161561094e57505050565b60408051808201825260008082526001602080840191825287835261013081528483206001600160a01b0388168452815284832093518454925161ffff1990931690151561ff00191617610100921515929092029190911790925585815261012f909152206109bd908361123e565b50505050565b600081815261012f6020526040902060609061050190611253565b600091825260c9602090815260408084206001600160a01b0393909316845291905290205460ff1690565b606060005b60208160ff16108015610a425750828160ff1660208110610a3157610a31611bf4565b1a60f81b6001600160f81b03191615155b15610a595780610a5181611c20565b915050610a0e565b60008160ff1667ffffffffffffffff811115610a7757610a776119cc565b6040519080825280601f01601f191660200182016040528015610aa1576020820181803683370190505b509050600091505b60208260ff16108015610add5750838260ff1660208110610acc57610acc611bf4565b1a60f81b6001600160f81b03191615155b15610b3957838260ff1660208110610af757610af7611bf4565b1a60f81b818360ff1681518110610b1057610b10611bf4565b60200101906001600160f81b031916908160001a90535081610b3181611c20565b925050610aa9565b9392505050565b60007e0caaa0e08f624de190c2474175cd13784c8c75bbdd1b63ae5fab5540967b3c610b6b81610f9d565b610b73611260565b6000848152610130602090815260408083206001600160a01b0387168452909152902054610100900460ff1615610bad5760019150610da1565b6000848152610130602090815260408083206001600160a01b038716845290915290205460ff1615610be25760009150610da1565b61012e546001600160a01b0384811691161415610c48576000848152610130602090815260408083206001600160a01b03871684528252808320805461ff00191661010017905586835261012f9091529020610c3e908461123e565b5060019150610da1565b61012d546001600160a01b03166367cfd9d6610c6386610a09565b856040518363ffffffff1660e01b8152600401610c81929190611c40565b602060405180830381865afa158015610c9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cc29190611c6a565b610ccf5760009150610da1565b61012d546001600160a01b03166316c75cfa610cea86610a09565b856040518363ffffffff1660e01b8152600401610d08929190611c40565b602060405180830381865afa158015610d25573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d499190611c6a565b610d565760009150610da1565b6000848152610130602090815260408083206001600160a01b03871684528252808320805461ff00191661010017905586835261012f9091529020610d9b908461123e565b50600191505b610dab600160fb55565b5092915050565b600080516020611dfd833981519152610dca81610f9d565b5061012e80546001600160a01b0319166001600160a01b0392909216919091179055565b600054610100900460ff1615808015610e0e5750600054600160ff909116105b80610e285750303b158015610e28575060005460ff166001145b610e8b5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610689565b6000805460ff191660011790558015610eae576000805461ff0019166101001790555b61012d80546001600160a01b0319166001600160a01b038416179055610ed26112c1565b610eda6112c1565b610ee26112ea565b610eed600033610fbc565b610f177f88aa719609f728b0c5e7fb8dd3608d5c25d497efbb3b9dd64e9251ebba10150833610fbc565b610f2f600080516020611dfd83398151915233610fbc565b801561069c576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15050565b600082815260c96020526040902060010154610f9381610f9d565b6105f38383611042565b61077d8133611319565b6000610b39836001600160a01b038416611372565b610fc682826109de565b61069c57600082815260c9602090815260408083206001600160a01b03851684529091529020805460ff19166001179055610ffe3390565b6001600160a01b0316816001600160a01b0316837f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45050565b61104c82826109de565b1561069c57600082815260c9602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b7f88aa719609f728b0c5e7fb8dd3608d5c25d497efbb3b9dd64e9251ebba10150861069c81610f9d565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615611106576105f383611465565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611160575060408051601f3d908101601f1916820190925261115d91810190611c8c565b60015b6111c35760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610689565b600080516020611e1d83398151915281146112325760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610689565b506105f3838383611501565b6000610b39836001600160a01b038416611526565b60606000610b3983611575565b600260fb5414156112b35760405162461bcd60e51b815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610689565b600260fb55565b600160fb55565b600054610100900460ff166112e85760405162461bcd60e51b815260040161068990611ca5565b565b600054610100900460ff166113115760405162461bcd60e51b815260040161068990611ca5565b6112e86115d1565b61132382826109de565b61069c57611330816115f8565b61133b83602061160a565b60405160200161134c929190611cf0565b60408051601f198184030181529082905262461bcd60e51b825261068991600401611b49565b6000818152600183016020526040812054801561145b576000611396600183611d65565b85549091506000906113aa90600190611d65565b905081811461140f5760008660000182815481106113ca576113ca611bf4565b90600052602060002001549050808760000184815481106113ed576113ed611bf4565b6000918252602080832090910192909255918252600188019052604090208390555b855486908061142057611420611d7c565b600190038181906000526020600020016000905590558560010160008681526020019081526020016000206000905560019350505050610501565b6000915050610501565b6001600160a01b0381163b6114d25760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610689565b600080516020611e1d83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61150a836117a6565b6000825111806115175750805b156105f3576109bd83836117e6565b600081815260018301602052604081205461156d57508154600181810184556000848152602080822090930184905584548482528286019093526040902091909155610501565b506000610501565b6060816000018054806020026020016040519081016040528092919081815260200182805480156115c557602002820191906000526020600020905b8154815260200190600101908083116115b1575b50505050509050919050565b600054610100900460ff166112ba5760405162461bcd60e51b815260040161068990611ca5565b60606105016001600160a01b03831660145b60606000611619836002611d92565b611624906002611db1565b67ffffffffffffffff81111561163c5761163c6119cc565b6040519080825280601f01601f191660200182016040528015611666576020820181803683370190505b509050600360fc1b8160008151811061168157611681611bf4565b60200101906001600160f81b031916908160001a905350600f60fb1b816001815181106116b0576116b0611bf4565b60200101906001600160f81b031916908160001a90535060006116d4846002611d92565b6116df906001611db1565b90505b6001811115611757576f181899199a1a9b1b9c1cb0b131b232b360811b85600f166010811061171357611713611bf4565b1a60f81b82828151811061172957611729611bf4565b60200101906001600160f81b031916908160001a90535060049490941c9361175081611dc9565b90506116e2565b508315610b395760405162461bcd60e51b815260206004820181905260248201527f537472696e67733a20686578206c656e67746820696e73756666696369656e746044820152606401610689565b6117af81611465565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6060610b398383604051806060016040528060278152602001611e3d602791396060600080856001600160a01b0316856040516118239190611de0565b600060405180830381855af49150503d806000811461185e576040519150601f19603f3d011682016040523d82523d6000602084013e611863565b606091505b50915091506118748683838761187e565b9695505050505050565b606083156118ea5782516118e3576001600160a01b0385163b6118e35760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610689565b50816118f4565b6118f483836118fc565b949350505050565b81511561190c5781518083602001fd5b8060405162461bcd60e51b81526004016106899190611b49565b60006020828403121561193857600080fd5b81356001600160e01b031981168114610b3957600080fd5b80356001600160a01b038116811461196757600080fd5b919050565b6000806040838503121561197f57600080fd5b8235915061198f60208401611950565b90509250929050565b6000602082840312156119aa57600080fd5b5035919050565b6000602082840312156119c357600080fd5b610b3982611950565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156119f557600080fd5b6119fe83611950565b9150602083013567ffffffffffffffff80821115611a1b57600080fd5b818501915085601f830112611a2f57600080fd5b813581811115611a4157611a416119cc565b604051601f8201601f19908116603f01168101908382118183101715611a6957611a696119cc565b81604052828152886020848701011115611a8257600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6020808252825182820181905260009190848201906040850190845b81811015611ae55783516001600160a01b031683529284019291840191600101611ac0565b50909695505050505050565b60005b83811015611b0c578181015183820152602001611af4565b838111156109bd5750506000910152565b60008151808452611b35816020860160208601611af1565b601f01601f19169290920160200192915050565b602081526000610b396020830184611b1d565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060ff821660ff811415611c3757611c37611c0a565b60010192915050565b604081526000611c536040830185611b1d565b905060018060a01b03831660208301529392505050565b600060208284031215611c7c57600080fd5b81518015158114610b3957600080fd5b600060208284031215611c9e57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b7f416363657373436f6e74726f6c3a206163636f756e7420000000000000000000815260008351611d28816017850160208801611af1565b7001034b99036b4b9b9b4b733903937b6329607d1b6017918401918201528351611d59816028840160208801611af1565b01602801949350505050565b600082821015611d7757611d77611c0a565b500390565b634e487b7160e01b600052603160045260246000fd5b6000816000190483118215151615611dac57611dac611c0a565b500290565b60008219821115611dc457611dc4611c0a565b500190565b600081611dd857611dd8611c0a565b506000190190565b60008251611df2818460208701611af1565b919091019291505056feb19546dff01e856fb3f010c267a7b1c60363cf8a4664e21cc89c26224620214e360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220d58e285e24b4b7831822eff3a5f86a0992d3dd7b410accf0aeae13441167c00664736f6c634300080a0033",
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

// Bytes32ToString is a free data retrieval call binding the contract method 0x9201de55.
//
// Solidity: function bytes32ToString(bytes32 _bytes32) pure returns(string)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) Bytes32ToString(opts *bind.CallOpts, _bytes32 [32]byte) (string, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "bytes32ToString", _bytes32)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Bytes32ToString is a free data retrieval call binding the contract method 0x9201de55.
//
// Solidity: function bytes32ToString(bytes32 _bytes32) pure returns(string)
func (_BridgeFeeOracle *BridgeFeeOracleSession) Bytes32ToString(_bytes32 [32]byte) (string, error) {
	return _BridgeFeeOracle.Contract.Bytes32ToString(&_BridgeFeeOracle.CallOpts, _bytes32)
}

// Bytes32ToString is a free data retrieval call binding the contract method 0x9201de55.
//
// Solidity: function bytes32ToString(bytes32 _bytes32) pure returns(string)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) Bytes32ToString(_bytes32 [32]byte) (string, error) {
	return _BridgeFeeOracle.Contract.Bytes32ToString(&_BridgeFeeOracle.CallOpts, _bytes32)
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

// GetOracleList is a free data retrieval call binding the contract method 0x61166581.
//
// Solidity: function getOracleList(bytes32 _chainName) view returns(address[])
func (_BridgeFeeOracle *BridgeFeeOracleCaller) GetOracleList(opts *bind.CallOpts, _chainName [32]byte) ([]common.Address, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "getOracleList", _chainName)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOracleList is a free data retrieval call binding the contract method 0x61166581.
//
// Solidity: function getOracleList(bytes32 _chainName) view returns(address[])
func (_BridgeFeeOracle *BridgeFeeOracleSession) GetOracleList(_chainName [32]byte) ([]common.Address, error) {
	return _BridgeFeeOracle.Contract.GetOracleList(&_BridgeFeeOracle.CallOpts, _chainName)
}

// GetOracleList is a free data retrieval call binding the contract method 0x61166581.
//
// Solidity: function getOracleList(bytes32 _chainName) view returns(address[])
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) GetOracleList(_chainName [32]byte) ([]common.Address, error) {
	return _BridgeFeeOracle.Contract.GetOracleList(&_BridgeFeeOracle.CallOpts, _chainName)
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

// OracleStatus is a free data retrieval call binding the contract method 0xaf51ef10.
//
// Solidity: function oracleStatus(bytes32 , address ) view returns(bool isBlack, bool isActive)
func (_BridgeFeeOracle *BridgeFeeOracleCaller) OracleStatus(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (struct {
	IsBlack  bool
	IsActive bool
}, error) {
	var out []interface{}
	err := _BridgeFeeOracle.contract.Call(opts, &out, "oracleStatus", arg0, arg1)

	outstruct := new(struct {
		IsBlack  bool
		IsActive bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsBlack = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.IsActive = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// OracleStatus is a free data retrieval call binding the contract method 0xaf51ef10.
//
// Solidity: function oracleStatus(bytes32 , address ) view returns(bool isBlack, bool isActive)
func (_BridgeFeeOracle *BridgeFeeOracleSession) OracleStatus(arg0 [32]byte, arg1 common.Address) (struct {
	IsBlack  bool
	IsActive bool
}, error) {
	return _BridgeFeeOracle.Contract.OracleStatus(&_BridgeFeeOracle.CallOpts, arg0, arg1)
}

// OracleStatus is a free data retrieval call binding the contract method 0xaf51ef10.
//
// Solidity: function oracleStatus(bytes32 , address ) view returns(bool isBlack, bool isActive)
func (_BridgeFeeOracle *BridgeFeeOracleCallerSession) OracleStatus(arg0 [32]byte, arg1 common.Address) (struct {
	IsBlack  bool
	IsActive bool
}, error) {
	return _BridgeFeeOracle.Contract.OracleStatus(&_BridgeFeeOracle.CallOpts, arg0, arg1)
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

// ActiveOracle is a paid mutator transaction binding the contract method 0x5cbb5179.
//
// Solidity: function activeOracle(bytes32 _chainName, address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) ActiveOracle(opts *bind.TransactOpts, _chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "activeOracle", _chainName, _oracle)
}

// ActiveOracle is a paid mutator transaction binding the contract method 0x5cbb5179.
//
// Solidity: function activeOracle(bytes32 _chainName, address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) ActiveOracle(_chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.ActiveOracle(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
}

// ActiveOracle is a paid mutator transaction binding the contract method 0x5cbb5179.
//
// Solidity: function activeOracle(bytes32 _chainName, address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) ActiveOracle(_chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.ActiveOracle(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
}

// BlackOracle is a paid mutator transaction binding the contract method 0x2237bdad.
//
// Solidity: function blackOracle(bytes32 _chainName, address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) BlackOracle(opts *bind.TransactOpts, _chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "blackOracle", _chainName, _oracle)
}

// BlackOracle is a paid mutator transaction binding the contract method 0x2237bdad.
//
// Solidity: function blackOracle(bytes32 _chainName, address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleSession) BlackOracle(_chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.BlackOracle(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
}

// BlackOracle is a paid mutator transaction binding the contract method 0x2237bdad.
//
// Solidity: function blackOracle(bytes32 _chainName, address _oracle) returns()
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) BlackOracle(_chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.BlackOracle(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
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

// IsOnline is a paid mutator transaction binding the contract method 0xb338123c.
//
// Solidity: function isOnline(bytes32 _chainName, address _oracle) returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleTransactor) IsOnline(opts *bind.TransactOpts, _chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.contract.Transact(opts, "isOnline", _chainName, _oracle)
}

// IsOnline is a paid mutator transaction binding the contract method 0xb338123c.
//
// Solidity: function isOnline(bytes32 _chainName, address _oracle) returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleSession) IsOnline(_chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
	return _BridgeFeeOracle.Contract.IsOnline(&_BridgeFeeOracle.TransactOpts, _chainName, _oracle)
}

// IsOnline is a paid mutator transaction binding the contract method 0xb338123c.
//
// Solidity: function isOnline(bytes32 _chainName, address _oracle) returns(bool)
func (_BridgeFeeOracle *BridgeFeeOracleTransactorSession) IsOnline(_chainName [32]byte, _oracle common.Address) (*types.Transaction, error) {
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
