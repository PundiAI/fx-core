package contract

import (
	"context"
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

const (
	FIP20LogicAddress = "0x0000000000000000000000000000000000001001"
	WFXLogicAddress   = "0x0000000000000000000000000000000000001002"

	StakingAddress         = "0x0000000000000000000000000000000000001003"
	CrosschainAddress      = "0x0000000000000000000000000000000000001004"
	BridgeFeeAddress       = "0x0000000000000000000000000000000000001005"
	BridgeFeeOracleAddress = "0x0000000000000000000000000000000000001006"
)

const DefaultGasCap uint64 = 30000000

var (
	fip20Init = Contract{
		Address: common.HexToAddress(FIP20LogicAddress),
		ABI:     MustABIJson(FIP20UpgradableMetaData.ABI),
		Bin:     MustDecodeHex(FIP20UpgradableMetaData.Bin),
		// deploy code from solidity/contracts/fip20/FIP20Upgradable.sol
		Code: MustDecodeHex("0x6080604052600436106101095760003560e01c806370a08231116100955780639dc29fac116100645780639dc29fac146102bc578063a9059cbb146102dc578063dd62ed3e146102fc578063de7ea79d14610342578063f2fde38b1461036257600080fd5b806370a0823114610234578063715018a61461026a5780638da5cb5b1461027f57806395d89b41146102a757600080fd5b8063313ce567116100dc578063313ce567146101a85780633659cfe6146101ca57806340c10f19146101ec5780634f1ef2861461020c57806352d1902d1461021f57600080fd5b806306fdde031461010e578063095ea7b31461013957806318160ddd1461016957806323b872dd14610188575b600080fd5b34801561011a57600080fd5b50610123610382565b60405161013091906113d1565b60405180910390f35b34801561014557600080fd5b50610159610154366004611420565b610414565b6040519015158152602001610130565b34801561017557600080fd5b5060cc545b604051908152602001610130565b34801561019457600080fd5b506101596101a336600461144a565b61046a565b3480156101b457600080fd5b5060cb5460405160ff9091168152602001610130565b3480156101d657600080fd5b506101ea6101e5366004611486565b610517565b005b3480156101f857600080fd5b506101ea610207366004611420565b6105f7565b6101ea61021a36600461152d565b61060d565b34801561022b57600080fd5b5061017a6106da565b34801561024057600080fd5b5061017a61024f366004611486565b6001600160a01b0316600090815260cd602052604090205490565b34801561027657600080fd5b506101ea61078d565b34801561028b57600080fd5b506097546040516001600160a01b039091168152602001610130565b3480156102b357600080fd5b506101236107a1565b3480156102c857600080fd5b506101ea6102d7366004611420565b6107b0565b3480156102e857600080fd5b506101596102f7366004611420565b6107c2565b34801561030857600080fd5b5061017a61031736600461158f565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b34801561034e57600080fd5b506101ea61035d3660046115e2565b6107d8565b34801561036e57600080fd5b506101ea61037d366004611486565b610947565b606060c980546103919061166f565b80601f01602080910402602001604051908101604052809291908181526020018280546103bd9061166f565b801561040a5780601f106103df5761010080835404028352916020019161040a565b820191906000526020600020905b8154815290600101906020018083116103ed57829003601f168201915b5050505050905090565b60006104213384846109bd565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156104ed5760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b61050185336104fc86856116c0565b6109bd565b61050c858585610a3f565b506001949350505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010011614156105605760405162461bcd60e51b81526004016104e4906116d7565b7f00000000000000000000000000000000000000000000000000000000000010016001600160a01b03166105a9600080516020611808833981519152546001600160a01b031690565b6001600160a01b0316146105cf5760405162461bcd60e51b81526004016104e490611723565b6105d881610bee565b604080516000808252602082019092526105f491839190610bf6565b50565b6105ff610d66565b6106098282610dc0565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010011614156106565760405162461bcd60e51b81526004016104e4906116d7565b7f00000000000000000000000000000000000000000000000000000000000010016001600160a01b031661069f600080516020611808833981519152546001600160a01b031690565b6001600160a01b0316146106c55760405162461bcd60e51b81526004016104e490611723565b6106ce82610bee565b61060982826001610bf6565b6000306001600160a01b037f0000000000000000000000000000000000000000000000000000000000001001161461077a5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016104e4565b5060008051602061180883398151915290565b610795610d66565b61079f6000610e9f565b565b606060ca80546103919061166f565b6107b8610d66565b6106098282610ef1565b60006107cf338484610a3f565b50600192915050565b600054610100900460ff16158080156107f85750600054600160ff909116105b806108125750303b158015610812575060005460ff166001145b6108755760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016104e4565b6000805460ff191660011790558015610898576000805461ff0019166101001790555b84516108ab9060c990602088019061130c565b5083516108bf9060ca90602087019061130c565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b0384161790556108f2611033565b6108fa611062565b8015610940576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b61094f610d66565b6001600160a01b0381166109b45760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016104e4565b6105f481610e9f565b6001600160a01b038316610a135760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016104e4565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610a955760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016104e4565b6001600160a01b038216610aeb5760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016104e4565b6001600160a01b038316600090815260cd602052604090205481811015610b545760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016104e4565b610b5e82826116c0565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610b9490849061176f565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610be091815260200190565b60405180910390a350505050565b6105f4610d66565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610c2e57610c2983611089565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610c88575060408051601f3d908101601f19168201909252610c8591810190611787565b60015b610ceb5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016104e4565b6000805160206118088339815191528114610d5a5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016104e4565b50610c29838383611125565b6097546001600160a01b0316331461079f5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016104e4565b6001600160a01b038216610e165760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016104e4565b8060cc6000828254610e28919061176f565b90915550506001600160a01b038216600090815260cd602052604081208054839290610e5590849061176f565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b038216610f475760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016104e4565b6001600160a01b038216600090815260cd602052604090205481811015610fb05760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016104e4565b610fba82826116c0565b6001600160a01b038416600090815260cd602052604081209190915560cc8054849290610fe89084906116c0565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b600054610100900460ff1661105a5760405162461bcd60e51b81526004016104e4906117a0565b61079f611150565b600054610100900460ff1661079f5760405162461bcd60e51b81526004016104e4906117a0565b6001600160a01b0381163b6110f65760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016104e4565b60008051602061180883398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61112e83611180565b60008251118061113b5750805b15610c295761114a83836111c0565b50505050565b600054610100900460ff166111775760405162461bcd60e51b81526004016104e4906117a0565b61079f33610e9f565b61118981611089565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606111e58383604051806060016040528060278152602001611828602791396111ec565b9392505050565b6060600080856001600160a01b03168560405161120991906117eb565b600060405180830381855af49150503d8060008114611244576040519150601f19603f3d011682016040523d82523d6000602084013e611249565b606091505b509150915061125a86838387611264565b9695505050505050565b606083156112d05782516112c9576001600160a01b0385163b6112c95760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016104e4565b50816112da565b6112da83836112e2565b949350505050565b8151156112f25781518083602001fd5b8060405162461bcd60e51b81526004016104e491906113d1565b8280546113189061166f565b90600052602060002090601f01602090048101928261133a5760008555611380565b82601f1061135357805160ff1916838001178555611380565b82800160010185558215611380579182015b82811115611380578251825591602001919060010190611365565b5061138c929150611390565b5090565b5b8082111561138c5760008155600101611391565b60005b838110156113c05781810151838201526020016113a8565b8381111561114a5750506000910152565b60208152600082518060208401526113f08160408501602087016113a5565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461141b57600080fd5b919050565b6000806040838503121561143357600080fd5b61143c83611404565b946020939093013593505050565b60008060006060848603121561145f57600080fd5b61146884611404565b925061147660208501611404565b9150604084013590509250925092565b60006020828403121561149857600080fd5b6111e582611404565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff808411156114d2576114d26114a1565b604051601f8501601f19908116603f011681019082821181831017156114fa576114fa6114a1565b8160405280935085815286868601111561151357600080fd5b858560208301376000602087830101525050509392505050565b6000806040838503121561154057600080fd5b61154983611404565b9150602083013567ffffffffffffffff81111561156557600080fd5b8301601f8101851361157657600080fd5b611585858235602084016114b7565b9150509250929050565b600080604083850312156115a257600080fd5b6115ab83611404565b91506115b960208401611404565b90509250929050565b600082601f8301126115d357600080fd5b6111e5838335602085016114b7565b600080600080608085870312156115f857600080fd5b843567ffffffffffffffff8082111561161057600080fd5b61161c888389016115c2565b9550602087013591508082111561163257600080fd5b5061163f878288016115c2565b935050604085013560ff8116811461165657600080fd5b915061166460608601611404565b905092959194509250565b600181811c9082168061168357607f821691505b602082108114156116a457634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156116d2576116d26116aa565b500390565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b60008219821115611782576117826116aa565b500190565b60006020828403121561179957600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b600082516117fd8184602087016113a5565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220c24d636966a60df9636041ea92ff12a963fc67961fabb8564b1631506b7785e864736f6c634300080a0033"),
	}
	wfxInit = Contract{
		Address: common.HexToAddress(WFXLogicAddress),
		ABI:     MustABIJson(WFXUpgradableMetaData.ABI),
		Bin:     MustDecodeHex(WFXUpgradableMetaData.Bin),
		// deploy code from solidity/contracts/fip20/WFXUpgradable.sol
		Code: MustDecodeHex("0x6080604052600436106101395760003560e01c8063715018a6116100ab578063b86d52981161006f578063b86d529814610366578063d0e30db014610148578063dd62ed3e14610384578063de7ea79d146103ca578063f2fde38b146103ea578063f3fef3a31461040a57610148565b8063715018a6146102ca5780638da5cb5b146102df57806395d89b41146103115780639dc29fac14610326578063a9059cbb1461034657610148565b8063313ce567116100fd578063313ce5671461020a5780633659cfe61461022c57806340c10f191461024c5780634f1ef2861461026c57806352d1902d1461027f57806370a082311461029457610148565b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab57806323b872dd146101ca5780632e1a7d4d146101ea57610148565b366101485761014661042a565b005b61014661042a565b34801561015c57600080fd5b5061016561046b565b60405161017291906115a9565b60405180910390f35b34801561018757600080fd5b5061019b6101963660046115f1565b6104fd565b6040519015158152602001610172565b3480156101b757600080fd5b5060cc545b604051908152602001610172565b3480156101d657600080fd5b5061019b6101e536600461161d565b610553565b3480156101f657600080fd5b5061014661020536600461165e565b610600565b34801561021657600080fd5b5060cb5460405160ff9091168152602001610172565b34801561023857600080fd5b50610146610247366004611677565b610671565b34801561025857600080fd5b506101466102673660046115f1565b610751565b61014661027a366004611720565b610767565b34801561028b57600080fd5b506101bc610834565b3480156102a057600080fd5b506101bc6102af366004611677565b6001600160a01b0316600090815260cd602052604090205490565b3480156102d657600080fd5b506101466108e7565b3480156102eb57600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610172565b34801561031d57600080fd5b506101656108fb565b34801561033257600080fd5b506101466103413660046115f1565b61090a565b34801561035257600080fd5b5061019b6103613660046115f1565b61091c565b34801561037257600080fd5b5060cf546001600160a01b03166102f9565b34801561039057600080fd5b506101bc61039f366004611784565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103d657600080fd5b506101466103e53660046117dd565b610932565b3480156103f657600080fd5b50610146610405366004611677565b610aa1565b34801561041657600080fd5b506101466104253660046115f1565b610b17565b6104343334610b9c565b60405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b606060c9805461047a9061186c565b80601f01602080910402602001604051908101604052809291908181526020018280546104a69061186c565b80156104f35780601f106104c8576101008083540402835291602001916104f3565b820191906000526020600020905b8154815290600101906020018083116104d657829003601f168201915b5050505050905090565b600061050a338484610c74565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156105d65760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b6105ea85336105e586856118bd565b610c74565b6105f5858585610cf6565b506001949350505050565b61060b335b82610ea5565b604051339082156108fc029083906000818181858888f19350505050158015610638573d6000803e3d6000fd5b5060405181815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a250565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010021614156106ba5760405162461bcd60e51b81526004016105cd906118d4565b7f00000000000000000000000000000000000000000000000000000000000010026001600160a01b0316610703600080516020611a05833981519152546001600160a01b031690565b6001600160a01b0316146107295760405162461bcd60e51b81526004016105cd90611920565b61073281610fe7565b6040805160008082526020820190925261074e91839190610fef565b50565b61075961115f565b6107638282610b9c565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010021614156107b05760405162461bcd60e51b81526004016105cd906118d4565b7f00000000000000000000000000000000000000000000000000000000000010026001600160a01b03166107f9600080516020611a05833981519152546001600160a01b031690565b6001600160a01b03161461081f5760405162461bcd60e51b81526004016105cd90611920565b61082882610fe7565b61076382826001610fef565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000100216146108d45760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016105cd565b50600080516020611a0583398151915290565b6108ef61115f565b6108f960006111b9565b565b606060ca805461047a9061186c565b61091261115f565b6107638282610ea5565b6000610929338484610cf6565b50600192915050565b600054610100900460ff16158080156109525750600054600160ff909116105b8061096c5750303b15801561096c575060005460ff166001145b6109cf5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016105cd565b6000805460ff1916600117905580156109f2576000805461ff0019166101001790555b8451610a059060c99060208801906114e4565b508351610a199060ca9060208701906114e4565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610a4c61120b565b610a5461123a565b8015610a9a576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b610aa961115f565b6001600160a01b038116610b0e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016105cd565b61074e816111b9565b610b2033610605565b6040516001600160a01b0383169082156108fc029083906000818181858888f19350505050158015610b56573d6000803e3d6000fd5b506040518181526001600160a01b0383169033907f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb906020015b60405180910390a35050565b6001600160a01b038216610bf25760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016105cd565b8060cc6000828254610c04919061196c565b90915550506001600160a01b038216600090815260cd602052604081208054839290610c3190849061196c565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610b90565b6001600160a01b038316610cca5760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016105cd565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610d4c5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016105cd565b6001600160a01b038216610da25760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016105cd565b6001600160a01b038316600090815260cd602052604090205481811015610e0b5760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016105cd565b610e1582826118bd565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610e4b90849061196c565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610e9791815260200190565b60405180910390a350505050565b6001600160a01b038216610efb5760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016105cd565b6001600160a01b038216600090815260cd602052604090205481811015610f645760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016105cd565b610f6e82826118bd565b6001600160a01b038416600090815260cd602052604081209190915560cc8054849290610f9c9084906118bd565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b61074e61115f565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff16156110275761102283611261565b505050565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015611081575060408051601f3d908101601f1916820190925261107e91810190611984565b60015b6110e45760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016105cd565b600080516020611a0583398151915281146111535760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016105cd565b506110228383836112fd565b6097546001600160a01b031633146108f95760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105cd565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166112325760405162461bcd60e51b81526004016105cd9061199d565b6108f9611328565b600054610100900460ff166108f95760405162461bcd60e51b81526004016105cd9061199d565b6001600160a01b0381163b6112ce5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016105cd565b600080516020611a0583398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61130683611358565b6000825111806113135750805b15611022576113228383611398565b50505050565b600054610100900460ff1661134f5760405162461bcd60e51b81526004016105cd9061199d565b6108f9336111b9565b61136181611261565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606113bd8383604051806060016040528060278152602001611a25602791396113c4565b9392505050565b6060600080856001600160a01b0316856040516113e191906119e8565b600060405180830381855af49150503d806000811461141c576040519150601f19603f3d011682016040523d82523d6000602084013e611421565b606091505b50915091506114328683838761143c565b9695505050505050565b606083156114a85782516114a1576001600160a01b0385163b6114a15760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016105cd565b50816114b2565b6114b283836114ba565b949350505050565b8151156114ca5781518083602001fd5b8060405162461bcd60e51b81526004016105cd91906115a9565b8280546114f09061186c565b90600052602060002090601f0160209004810192826115125760008555611558565b82601f1061152b57805160ff1916838001178555611558565b82800160010185558215611558579182015b8281111561155857825182559160200191906001019061153d565b50611564929150611568565b5090565b5b808211156115645760008155600101611569565b60005b83811015611598578181015183820152602001611580565b838111156113225750506000910152565b60208152600082518060208401526115c881604085016020870161157d565b601f01601f19169190910160400192915050565b6001600160a01b038116811461074e57600080fd5b6000806040838503121561160457600080fd5b823561160f816115dc565b946020939093013593505050565b60008060006060848603121561163257600080fd5b833561163d816115dc565b9250602084013561164d816115dc565b929592945050506040919091013590565b60006020828403121561167057600080fd5b5035919050565b60006020828403121561168957600080fd5b81356113bd816115dc565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff808411156116c5576116c5611694565b604051601f8501601f19908116603f011681019082821181831017156116ed576116ed611694565b8160405280935085815286868601111561170657600080fd5b858560208301376000602087830101525050509392505050565b6000806040838503121561173357600080fd5b823561173e816115dc565b9150602083013567ffffffffffffffff81111561175a57600080fd5b8301601f8101851361176b57600080fd5b61177a858235602084016116aa565b9150509250929050565b6000806040838503121561179757600080fd5b82356117a2816115dc565b915060208301356117b2816115dc565b809150509250929050565b600082601f8301126117ce57600080fd5b6113bd838335602085016116aa565b600080600080608085870312156117f357600080fd5b843567ffffffffffffffff8082111561180b57600080fd5b611817888389016117bd565b9550602087013591508082111561182d57600080fd5b5061183a878288016117bd565b935050604085013560ff8116811461185157600080fd5b91506060850135611861816115dc565b939692955090935050565b600181811c9082168061188057607f821691505b602082108114156118a157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156118cf576118cf6118a7565b500390565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6000821982111561197f5761197f6118a7565b500190565b60006020828403121561199657600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b600082516119fa81846020870161157d565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212203491f4b9433597e502b7d2b2ea2025687fa17dcc511ea87b77ea73cb9c0109d464736f6c634300080a0033"),
	}

	erc1967Proxy = Contract{
		Address: common.Address{},
		ABI:     MustABIJson(ERC1967ProxyMetaData.ABI),
		Bin:     MustDecodeHex(ERC1967ProxyMetaData.Bin),
		Code:    []byte{},
	}

	bridgeProxy = Contract{
		Address: common.Address{},
		ABI:     MustABIJson(BridgeProxyMetaData.ABI),
		Bin:     MustDecodeHex(BridgeProxyMetaData.Bin),
		// deploy code from solidity/contracts/bridge/BridgeProxy.sol
		Code: MustDecodeHex("0x6080604052600436106100225760003560e01c806319ab453c1461003957610031565b366100315761002f610059565b005b61002f610059565b34801561004557600080fd5b5061002f61005436600461021e565b61006b565b6100696100646100d0565b610108565b565b600061009e7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b0316146100c45760405162dc149f60e41b815260040160405180910390fd5b6100cd8161012c565b50565b60006101037f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e808015610127573d6000f35b3d6000fd5b6101358161016c565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b6001600160a01b0381163b6101dd5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840160405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b60006020828403121561023057600080fd5b81356001600160a01b038116811461024757600080fd5b939250505056fea2646970667358221220f5fd2c6b493d2d8c8fccb10ef0e5ed3a7f4f1914bae22833db4f76dbce6da4fe64736f6c634300080a0033"),
	}

	bridgeFeeQuote = Contract{
		Address: common.Address{},
		ABI:     MustABIJson(BridgeFeeQuoteMetaData.ABI),
		Bin:     MustDecodeHex(BridgeFeeQuoteMetaData.Bin),
		Code:    []byte{},
	}

	bridgeFeeOracle = Contract{
		Address: common.Address{},
		ABI:     MustABIJson(BridgeFeeOracleMetaData.ABI),
		Bin:     MustDecodeHex(BridgeFeeOracleMetaData.Bin),
		Code:    []byte{},
	}

	fxBridgeABI          = MustABIJson(IFxBridgeLogicMetaData.ABI)
	bridgeCallContextABI = MustABIJson(IBridgeCallContextMetaData.ABI)
	errorABI             = MustABIJson(IErrorMetaData.ABI)
)

type Caller interface {
	QueryContract(ctx context.Context, from, contract common.Address, abi abi.ABI, method string, res interface{}, args ...interface{}) error
	ApplyContract(ctx context.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
}

type Contract struct {
	Address common.Address
	ABI     abi.ABI
	Bin     []byte
	Code    []byte
}

func (c Contract) CodeHash() common.Hash {
	return crypto.Keccak256Hash(c.Code)
}

func GetFIP20() Contract {
	return fip20Init
}

func GetWFX() Contract {
	return wfxInit
}

func GetERC1967Proxy() Contract {
	return erc1967Proxy
}

func GetBridgeProxy() Contract {
	return bridgeProxy
}

func GetBridgeFeeQuote() Contract {
	return bridgeFeeQuote
}

func GetBridgeFeeOracle() Contract {
	return bridgeFeeOracle
}

func MustDecodeHex(str string) []byte {
	bz, err := hexutil.Decode(str)
	if err != nil {
		panic(err)
	}
	return bz
}

func MustABIJson(str string) abi.ABI {
	j, err := abi.JSON(strings.NewReader(str))
	if err != nil {
		panic(err)
	}
	return j
}

func PackRetErrV2(err error) ([]byte, error) {
	pack, _ := errorABI.Pack("Error", err.Error())
	return pack, err
}

func PackOnBridgeCall(sender, receiver common.Address, tokens []common.Address, amounts []*big.Int, data, memo []byte) ([]byte, error) {
	return bridgeCallContextABI.Pack("onBridgeCall",
		sender,
		receiver,
		tokens,
		amounts,
		data,
		memo,
	)
}

func PackOracleSetCheckpoint(gravityID, methodName [32]byte, nonce *big.Int, memberAddresses []common.Address, convertedPowers []*big.Int) ([]byte, error) {
	return fxBridgeABI.Pack("oracleSetCheckpoint",
		gravityID,
		methodName,
		nonce,
		memberAddresses,
		convertedPowers,
	)
}

func PackSubmitBatchCheckpoint(gravityID, methodName [32]byte, amounts []*big.Int, destinations []common.Address, fees []*big.Int, batchNonce *big.Int, tokenContract common.Address, batchTimeout *big.Int, feeReceive common.Address) ([]byte, error) {
	return fxBridgeABI.Pack("submitBatchCheckpoint",
		gravityID,
		methodName,
		amounts,
		destinations,
		fees,
		batchNonce,
		tokenContract,
		batchTimeout,
		feeReceive,
	)
}

func PackBridgeCallCheckpoint(gravityID, methodName [32]byte, sender, refund common.Address, tokens []common.Address, amounts []*big.Int, to common.Address, data, memo []byte, nonce, timeout, eventNonce *big.Int) ([]byte, error) {
	return fxBridgeABI.Pack("bridgeCallCheckpoint",
		gravityID,
		methodName,
		sender,
		refund,
		tokens,
		amounts,
		to,
		data,
		memo,
		nonce,
		timeout,
		eventNonce,
	)
}

func unpackRetIsOk(abi abi.ABI, method string, res *evmtypes.MsgEthereumTxResponse) (*evmtypes.MsgEthereumTxResponse, error) {
	var ret struct{ Value bool }
	if err := abi.UnpackIntoInterface(&ret, method, res.Ret); err != nil {
		return res, sdkerrors.ErrInvalidType.Wrapf("failed to unpack %s: %s", method, err.Error())
	}
	if !ret.Value {
		return res, sdkerrors.ErrLogic.Wrapf("failed to execute %s", method)
	}
	return res, nil
}
