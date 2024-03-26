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

// ERC721TokenTestMetaData contains all meta data concerning the ERC721TokenTest contract.
var ERC721TokenTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604080518082018252600f81526e115490cdcc8c551bdad95b95195cdd608a1b60208083019182528351808501909452600384526215151560ea1b9084015281519192916200006491600091620000f3565b5080516200007a906001906020840190620000f3565b50505062000097620000916200009d60201b60201c565b620000a1565b620001d6565b3390565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b828054620001019062000199565b90600052602060002090601f01602090048101928262000125576000855562000170565b82601f106200014057805160ff191683800117855562000170565b8280016001018555821562000170579182015b828111156200017057825182559160200191906001019062000153565b506200017e92915062000182565b5090565b5b808211156200017e576000815560010162000183565b600281046001821680620001ae57607f821691505b60208210811415620001d057634e487b7160e01b600052602260045260246000fd5b50919050565b6115a580620001e66000396000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c806370a08231116100a2578063a22cb46511610071578063a22cb46514610239578063b88d4fde1461024c578063c87b56dd1461025f578063e985e9c514610272578063f2fde38b146102ae57610116565b806370a08231146101f7578063715018a6146102185780638da5cb5b1461022057806395d89b411461023157610116565b806323b872dd116100e957806323b872dd1461019857806340c10f19146101ab57806342842e0e146101be57806342966c68146101d15780636352211e146101e457610116565b806301ffc9a71461011b57806306fdde0314610143578063081812fc14610158578063095ea7b314610183575b600080fd5b61012e6101293660046112fd565b6102c1565b60405190151581526020015b60405180910390f35b61014b610315565b60405161013a91906113e5565b61016b610166366004611335565b6103a7565b6040516001600160a01b03909116815260200161013a565b6101966101913660046112d4565b6103ce565b005b6101966101a636600461118a565b6104e9565b6101966101b93660046112d4565b61051b565b6101966101cc36600461118a565b610531565b6101966101df366004611335565b61054c565b61016b6101f2366004611335565b61057d565b61020a61020536600461113e565b6105dd565b60405190815260200161013a565b610196610663565b6006546001600160a01b031661016b565b61014b610677565b61019661024736600461129a565b610686565b61019661025a3660046111c5565b610691565b61014b61026d366004611335565b6106c9565b61012e610280366004611158565b6001600160a01b03918216600090815260056020908152604080832093909416825291909152205460ff1690565b6101966102bc36600461113e565b610754565b60006001600160e01b031982166380ac58cd60e01b14806102f257506001600160e01b03198216635b5e139f60e01b145b8061030d57506301ffc9a760e01b6001600160e01b03198316145b90505b919050565b60606000805461032490611508565b80601f016020809104026020016040519081016040528092919081815260200182805461035090611508565b801561039d5780601f106103725761010080835404028352916020019161039d565b820191906000526020600020905b81548152906001019060200180831161038057829003601f168201915b5050505050905090565b60006103b2826107ca565b506000908152600460205260409020546001600160a01b031690565b60006103d98261057d565b9050806001600160a01b0316836001600160a01b0316141561044c5760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b60648201526084015b60405180910390fd5b336001600160a01b038216148061046857506104688133610280565b6104da5760405162461bcd60e51b815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c0000006064820152608401610443565b6104e48383610829565b505050565b6104f4335b82610897565b6105105760405162461bcd60e51b8152600401610443906113f8565b6104e4838383610916565b610523610a7a565b61052d8282610ad4565b5050565b6104e483838360405180602001604052806000815250610691565b610555336104ee565b6105715760405162461bcd60e51b8152600401610443906113f8565b61057a81610aee565b50565b6000818152600260205260408120546001600160a01b03168061030d5760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b6044820152606401610443565b60006001600160a01b0382166106475760405162461bcd60e51b815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f7420612076616044820152683634b21037bbb732b960b91b6064820152608401610443565b506001600160a01b031660009081526003602052604090205490565b61066b610a7a565b6106756000610b84565b565b60606001805461032490611508565b61052d338383610bd6565b61069b3383610897565b6106b75760405162461bcd60e51b8152600401610443906113f8565b6106c384848484610ca5565b50505050565b60606106d4826107ca565b600061070260408051808201909152600f81526e1a5c199cce8bcbdd195cdd0b5d5c9b608a1b602082015290565b90506000815111610722576040518060200160405280600081525061074d565b8061072c84610cd8565b60405160200161073d929190611379565b6040516020818303038152906040525b9392505050565b61075c610a7a565b6001600160a01b0381166107c15760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610443565b61057a81610b84565b6000818152600260205260409020546001600160a01b031661057a5760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b6044820152606401610443565b600081815260046020526040902080546001600160a01b0319166001600160a01b038416908117909155819061085e8261057d565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b6000806108a38361057d565b9050806001600160a01b0316846001600160a01b031614806108ea57506001600160a01b0380821660009081526005602090815260408083209388168352929052205460ff165b8061090e5750836001600160a01b0316610903846103a7565b6001600160a01b0316145b949350505050565b826001600160a01b03166109298261057d565b6001600160a01b03161461094f5760405162461bcd60e51b815260040161044390611497565b6001600160a01b0382166109b15760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b6064820152608401610443565b826001600160a01b03166109c48261057d565b6001600160a01b0316146109ea5760405162461bcd60e51b815260040161044390611497565b600081815260046020908152604080832080546001600160a01b03199081169091556001600160a01b0387811680865260038552838620805460001901905590871680865283862080546001019055868652600290945282852080549092168417909155905184937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a46104e4565b6006546001600160a01b031633146106755760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610443565b61052d828260405180602001604052806000815250610d83565b6000610af98261057d565b9050610b048261057d565b600083815260046020908152604080832080546001600160a01b03199081169091556001600160a01b0385168085526003845282852080546000190190558785526002909352818420805490911690555192935084927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908390a461052d565b600680546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b816001600160a01b0316836001600160a01b03161415610c385760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c6572000000000000006044820152606401610443565b6001600160a01b03838116600081815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b610cb0848484610916565b610cbc84848484610db6565b6106c35760405162461bcd60e51b815260040161044390611445565b60606000610ce583610ec3565b600101905060008167ffffffffffffffff811115610d1357634e487b7160e01b600052604160045260246000fd5b6040519080825280601f01601f191660200182016040528015610d3d576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a8504945084610d7657610d7b565b610d47565b509392505050565b610d8d8383610f9b565b610d9a6000848484610db6565b6104e45760405162461bcd60e51b815260040161044390611445565b60006001600160a01b0384163b15610eb857604051630a85bd0160e11b81526001600160a01b0385169063150b7a0290610dfa9033908990889088906004016113a8565b602060405180830381600087803b158015610e1457600080fd5b505af1925050508015610e44575060408051601f3d908101601f19168201909252610e4191810190611319565b60015b610e9e573d808015610e72576040519150601f19603f3d011682016040523d82523d6000602084013e610e77565b606091505b508051610e965760405162461bcd60e51b815260040161044390611445565b805181602001fd5b6001600160e01b031916630a85bd0160e11b14905061090e565b506001949350505050565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b8310610f025772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310610f2e576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310610f4c57662386f26fc10000830492506010015b6305f5e1008310610f64576305f5e100830492506008015b6127108310610f7857612710830492506004015b60648310610f8a576064830492506002015b600a831061030d5760010192915050565b6001600160a01b038216610ff15760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f20616464726573736044820152606401610443565b6000818152600260205260409020546001600160a01b0316156110565760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e746564000000006044820152606401610443565b6000818152600260205260409020546001600160a01b0316156110bb5760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e746564000000006044820152606401610443565b6001600160a01b038216600081815260036020908152604080832080546001019055848352600290915280822080546001600160a01b0319168417905551839291907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a461052d565b80356001600160a01b038116811461031057600080fd5b60006020828403121561114f578081fd5b61074d82611127565b6000806040838503121561116a578081fd5b61117383611127565b915061118160208401611127565b90509250929050565b60008060006060848603121561119e578081fd5b6111a784611127565b92506111b560208501611127565b9150604084013590509250925092565b600080600080608085870312156111da578081fd5b6111e385611127565b93506111f160208601611127565b925060408501359150606085013567ffffffffffffffff80821115611214578283fd5b818701915087601f830112611227578283fd5b81358181111561123957611239611543565b604051601f8201601f19908116603f0116810190838211818310171561126157611261611543565b816040528281528a6020848701011115611279578586fd5b82602086016020830137918201602001949094529598949750929550505050565b600080604083850312156112ac578182fd5b6112b583611127565b9150602083013580151581146112c9578182fd5b809150509250929050565b600080604083850312156112e6578182fd5b6112ef83611127565b946020939093013593505050565b60006020828403121561130e578081fd5b813561074d81611559565b60006020828403121561132a578081fd5b815161074d81611559565b600060208284031215611346578081fd5b5035919050565b600081518084526113658160208601602086016114dc565b601f01601f19169290920160200192915050565b6000835161138b8184602088016114dc565b83519083019061139f8183602088016114dc565b01949350505050565b6001600160a01b03858116825284166020820152604081018390526080606082018190526000906113db9083018461134d565b9695505050505050565b60006020825261074d602083018461134d565b6020808252602d908201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560408201526c1c881bdc88185c1c1c9bdd9959609a1b606082015260800190565b60208082526032908201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560408201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b606082015260800190565b60208082526025908201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060408201526437bbb732b960d91b606082015260800190565b60005b838110156114f75781810151838201526020016114df565b838111156106c35750506000910152565b60028104600182168061151c57607f821691505b6020821081141561153d57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b6001600160e01b03198116811461057a57600080fdfea264697066735822122004b4845a3825244943d3ce0f7d20aaac8e24eabce796a7b908426fa0a1b5f83a64736f6c63430008020033",
}

// ERC721TokenTestABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC721TokenTestMetaData.ABI instead.
var ERC721TokenTestABI = ERC721TokenTestMetaData.ABI

// ERC721TokenTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC721TokenTestMetaData.Bin instead.
var ERC721TokenTestBin = ERC721TokenTestMetaData.Bin

// DeployERC721TokenTest deploys a new Ethereum contract, binding an instance of ERC721TokenTest to it.
func DeployERC721TokenTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC721TokenTest, error) {
	parsed, err := ERC721TokenTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC721TokenTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721TokenTest{ERC721TokenTestCaller: ERC721TokenTestCaller{contract: contract}, ERC721TokenTestTransactor: ERC721TokenTestTransactor{contract: contract}, ERC721TokenTestFilterer: ERC721TokenTestFilterer{contract: contract}}, nil
}

// ERC721TokenTest is an auto generated Go binding around an Ethereum contract.
type ERC721TokenTest struct {
	ERC721TokenTestCaller     // Read-only binding to the contract
	ERC721TokenTestTransactor // Write-only binding to the contract
	ERC721TokenTestFilterer   // Log filterer for contract events
}

// ERC721TokenTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC721TokenTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721TokenTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC721TokenTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721TokenTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC721TokenTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721TokenTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC721TokenTestSession struct {
	Contract     *ERC721TokenTest  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721TokenTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC721TokenTestCallerSession struct {
	Contract *ERC721TokenTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ERC721TokenTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC721TokenTestTransactorSession struct {
	Contract     *ERC721TokenTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ERC721TokenTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC721TokenTestRaw struct {
	Contract *ERC721TokenTest // Generic contract binding to access the raw methods on
}

// ERC721TokenTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC721TokenTestCallerRaw struct {
	Contract *ERC721TokenTestCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721TokenTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC721TokenTestTransactorRaw struct {
	Contract *ERC721TokenTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721TokenTest creates a new instance of ERC721TokenTest, bound to a specific deployed contract.
func NewERC721TokenTest(address common.Address, backend bind.ContractBackend) (*ERC721TokenTest, error) {
	contract, err := bindERC721TokenTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTest{ERC721TokenTestCaller: ERC721TokenTestCaller{contract: contract}, ERC721TokenTestTransactor: ERC721TokenTestTransactor{contract: contract}, ERC721TokenTestFilterer: ERC721TokenTestFilterer{contract: contract}}, nil
}

// NewERC721TokenTestCaller creates a new read-only instance of ERC721TokenTest, bound to a specific deployed contract.
func NewERC721TokenTestCaller(address common.Address, caller bind.ContractCaller) (*ERC721TokenTestCaller, error) {
	contract, err := bindERC721TokenTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestCaller{contract: contract}, nil
}

// NewERC721TokenTestTransactor creates a new write-only instance of ERC721TokenTest, bound to a specific deployed contract.
func NewERC721TokenTestTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721TokenTestTransactor, error) {
	contract, err := bindERC721TokenTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestTransactor{contract: contract}, nil
}

// NewERC721TokenTestFilterer creates a new log filterer instance of ERC721TokenTest, bound to a specific deployed contract.
func NewERC721TokenTestFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721TokenTestFilterer, error) {
	contract, err := bindERC721TokenTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestFilterer{contract: contract}, nil
}

// bindERC721TokenTest binds a generic wrapper to an already deployed contract.
func bindERC721TokenTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC721TokenTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721TokenTest *ERC721TokenTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721TokenTest.Contract.ERC721TokenTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721TokenTest *ERC721TokenTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.ERC721TokenTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721TokenTest *ERC721TokenTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.ERC721TokenTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721TokenTest *ERC721TokenTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC721TokenTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721TokenTest *ERC721TokenTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721TokenTest *ERC721TokenTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ERC721TokenTest *ERC721TokenTestCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ERC721TokenTest *ERC721TokenTestSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721TokenTest.Contract.BalanceOf(&_ERC721TokenTest.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ERC721TokenTest.Contract.BalanceOf(&_ERC721TokenTest.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721TokenTest *ERC721TokenTestCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721TokenTest *ERC721TokenTestSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721TokenTest.Contract.GetApproved(&_ERC721TokenTest.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ERC721TokenTest.Contract.GetApproved(&_ERC721TokenTest.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721TokenTest *ERC721TokenTestCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721TokenTest *ERC721TokenTestSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721TokenTest.Contract.IsApprovedForAll(&_ERC721TokenTest.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC721TokenTest.Contract.IsApprovedForAll(&_ERC721TokenTest.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721TokenTest *ERC721TokenTestCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721TokenTest *ERC721TokenTestSession) Name() (string, error) {
	return _ERC721TokenTest.Contract.Name(&_ERC721TokenTest.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) Name() (string, error) {
	return _ERC721TokenTest.Contract.Name(&_ERC721TokenTest.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC721TokenTest *ERC721TokenTestCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC721TokenTest *ERC721TokenTestSession) Owner() (common.Address, error) {
	return _ERC721TokenTest.Contract.Owner(&_ERC721TokenTest.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) Owner() (common.Address, error) {
	return _ERC721TokenTest.Contract.Owner(&_ERC721TokenTest.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721TokenTest *ERC721TokenTestCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721TokenTest *ERC721TokenTestSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721TokenTest.Contract.OwnerOf(&_ERC721TokenTest.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ERC721TokenTest.Contract.OwnerOf(&_ERC721TokenTest.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721TokenTest *ERC721TokenTestCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721TokenTest *ERC721TokenTestSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721TokenTest.Contract.SupportsInterface(&_ERC721TokenTest.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ERC721TokenTest.Contract.SupportsInterface(&_ERC721TokenTest.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721TokenTest *ERC721TokenTestCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721TokenTest *ERC721TokenTestSession) Symbol() (string, error) {
	return _ERC721TokenTest.Contract.Symbol(&_ERC721TokenTest.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) Symbol() (string, error) {
	return _ERC721TokenTest.Contract.Symbol(&_ERC721TokenTest.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721TokenTest *ERC721TokenTestCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721TokenTest *ERC721TokenTestSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721TokenTest.Contract.TokenURI(&_ERC721TokenTest.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ERC721TokenTest.Contract.TokenURI(&_ERC721TokenTest.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Approve(&_ERC721TokenTest.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Approve(&_ERC721TokenTest.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Burn(&_ERC721TokenTest.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Burn(&_ERC721TokenTest.TransactOpts, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _id) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) Mint(opts *bind.TransactOpts, _to common.Address, _id *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "mint", _to, _id)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _id) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) Mint(_to common.Address, _id *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Mint(&_ERC721TokenTest.TransactOpts, _to, _id)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _to, uint256 _id) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) Mint(_to common.Address, _id *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Mint(&_ERC721TokenTest.TransactOpts, _to, _id)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC721TokenTest *ERC721TokenTestSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.RenounceOwnership(&_ERC721TokenTest.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.RenounceOwnership(&_ERC721TokenTest.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SafeTransferFrom(&_ERC721TokenTest.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SafeTransferFrom(&_ERC721TokenTest.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SafeTransferFrom0(&_ERC721TokenTest.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SafeTransferFrom0(&_ERC721TokenTest.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SetApprovalForAll(&_ERC721TokenTest.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SetApprovalForAll(&_ERC721TokenTest.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.TransferFrom(&_ERC721TokenTest.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.TransferFrom(&_ERC721TokenTest.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.TransferOwnership(&_ERC721TokenTest.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.TransferOwnership(&_ERC721TokenTest.TransactOpts, newOwner)
}

// ERC721TokenTestApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC721TokenTest contract.
type ERC721TokenTestApprovalIterator struct {
	Event *ERC721TokenTestApproval // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestApproval)
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
		it.Event = new(ERC721TokenTestApproval)
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
func (it *ERC721TokenTestApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestApproval represents a Approval event raised by the ERC721TokenTest contract.
type ERC721TokenTestApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ERC721TokenTestApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestApprovalIterator{contract: _ERC721TokenTest.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestApproval)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseApproval(log types.Log) (*ERC721TokenTestApproval, error) {
	event := new(ERC721TokenTestApproval)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TokenTestApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC721TokenTest contract.
type ERC721TokenTestApprovalForAllIterator struct {
	Event *ERC721TokenTestApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestApprovalForAll)
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
		it.Event = new(ERC721TokenTestApprovalForAll)
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
func (it *ERC721TokenTestApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestApprovalForAll represents a ApprovalForAll event raised by the ERC721TokenTest contract.
type ERC721TokenTestApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC721TokenTestApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestApprovalForAllIterator{contract: _ERC721TokenTest.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestApprovalForAll)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseApprovalForAll(log types.Log) (*ERC721TokenTestApprovalForAll, error) {
	event := new(ERC721TokenTestApprovalForAll)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TokenTestOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ERC721TokenTest contract.
type ERC721TokenTestOwnershipTransferredIterator struct {
	Event *ERC721TokenTestOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestOwnershipTransferred)
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
		it.Event = new(ERC721TokenTestOwnershipTransferred)
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
func (it *ERC721TokenTestOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestOwnershipTransferred represents a OwnershipTransferred event raised by the ERC721TokenTest contract.
type ERC721TokenTestOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ERC721TokenTestOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestOwnershipTransferredIterator{contract: _ERC721TokenTest.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestOwnershipTransferred)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseOwnershipTransferred(log types.Log) (*ERC721TokenTestOwnershipTransferred, error) {
	event := new(ERC721TokenTestOwnershipTransferred)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TokenTestTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC721TokenTest contract.
type ERC721TokenTestTransferIterator struct {
	Event *ERC721TokenTestTransfer // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestTransfer)
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
		it.Event = new(ERC721TokenTestTransfer)
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
func (it *ERC721TokenTestTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestTransfer represents a Transfer event raised by the ERC721TokenTest contract.
type ERC721TokenTestTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ERC721TokenTestTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestTransferIterator{contract: _ERC721TokenTest.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestTransfer)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseTransfer(log types.Log) (*ERC721TokenTestTransfer, error) {
	event := new(ERC721TokenTestTransfer)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
