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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"safeMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100e1565b600054610100900460ff161561008e5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116146100df576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b608051612398620001196000396000818161061d0152818161065d015281816107180152818161075801526107eb01526123986000f3fe60806040526004361061012a5760003560e01c806370a08231116100ab578063a22cb4651161006f578063a22cb46514610313578063b88d4fde14610333578063c87b56dd14610353578063d204c45e14610373578063e985e9c514610393578063f2fde38b146103dc57600080fd5b806370a0823114610296578063715018a6146102b65780638129fc1c146102cb5780638da5cb5b146102e057806395d89b41146102fe57600080fd5b80633659cfe6116100f25780633659cfe61461020057806342842e0e146102205780634f1ef2861461024057806352d1902d146102535780636352211e1461027657600080fd5b806301ffc9a71461012f57806306fdde0314610164578063081812fc14610186578063095ea7b3146101be57806323b872dd146101e0575b600080fd5b34801561013b57600080fd5b5061014f61014a366004611cfe565b6103fc565b60405190151581526020015b60405180910390f35b34801561017057600080fd5b5061017961040d565b60405161015b9190611d73565b34801561019257600080fd5b506101a66101a1366004611d86565b61049f565b6040516001600160a01b03909116815260200161015b565b3480156101ca57600080fd5b506101de6101d9366004611dbb565b6104c6565b005b3480156101ec57600080fd5b506101de6101fb366004611de5565b6105e1565b34801561020c57600080fd5b506101de61021b366004611e21565b610612565b34801561022c57600080fd5b506101de61023b366004611de5565b6106f2565b6101de61024e366004611ee8565b61070d565b34801561025f57600080fd5b506102686107de565b60405190815260200161015b565b34801561028257600080fd5b506101a6610291366004611d86565b610891565b3480156102a257600080fd5b506102686102b1366004611e21565b6108f1565b3480156102c257600080fd5b506101de610977565b3480156102d757600080fd5b506101de61098b565b3480156102ec57600080fd5b5060c9546001600160a01b03166101a6565b34801561030a57600080fd5b50610179610af7565b34801561031f57600080fd5b506101de61032e366004611f36565b610b06565b34801561033f57600080fd5b506101de61034e366004611f72565b610b11565b34801561035f57600080fd5b5061017961036e366004611d86565b610b49565b34801561037f57600080fd5b506101de61038e366004611fda565b610b54565b34801561039f57600080fd5b5061014f6103ae366004612032565b6001600160a01b039182166000908152606a6020908152604080832093909416825291909152205460ff1690565b3480156103e857600080fd5b506101de6103f7366004611e21565b610b8d565b600061040782610c03565b92915050565b60606065805461041c90612065565b80601f016020809104026020016040519081016040528092919081815260200182805461044890612065565b80156104955780601f1061046a57610100808354040283529160200191610495565b820191906000526020600020905b81548152906001019060200180831161047857829003601f168201915b5050505050905090565b60006104aa82610c28565b506000908152606960205260409020546001600160a01b031690565b60006104d182610891565b9050806001600160a01b0316836001600160a01b031614156105445760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b60648201526084015b60405180910390fd5b336001600160a01b0382161480610560575061056081336103ae565b6105d25760405162461bcd60e51b815260206004820152603d60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206f7220617070726f76656420666f7220616c6c000000606482015260840161053b565b6105dc8383610c87565b505050565b6105eb3382610cf5565b6106075760405162461bcd60e51b815260040161053b906120a0565b6105dc838383610d74565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016141561065b5760405162461bcd60e51b815260040161053b906120ed565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106a460008051602061231c833981519152546001600160a01b031690565b6001600160a01b0316146106ca5760405162461bcd60e51b815260040161053b90612139565b6106d381610ed8565b604080516000808252602082019092526106ef91839190610ee0565b50565b6105dc83838360405180602001604052806000815250610b11565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614156107565760405162461bcd60e51b815260040161053b906120ed565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661079f60008051602061231c833981519152546001600160a01b031690565b6001600160a01b0316146107c55760405162461bcd60e51b815260040161053b90612139565b6107ce82610ed8565b6107da82826001610ee0565b5050565b6000306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461087e5760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c0000000000000000606482015260840161053b565b5060008051602061231c83398151915290565b6000818152606760205260408120546001600160a01b0316806104075760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b604482015260640161053b565b60006001600160a01b03821661095b5760405162461bcd60e51b815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f7420612076616044820152683634b21037bbb732b960b91b606482015260840161053b565b506001600160a01b031660009081526068602052604090205490565b61097f61104b565b61098960006110a5565b565b600054610100900460ff16158080156109ab5750600054600160ff909116105b806109c55750303b1580156109c5575060005460ff166001145b610a285760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b606482015260840161053b565b6000805460ff191660011790558015610a4b576000805461ff0019166101001790555b610a976040518060400160405280600f81526020016e115490cdcc8c551bdad95b95195cdd608a1b8152506040518060400160405280600381526020016215151560ea1b8152506110f7565b610a9f611128565b610aa761114f565b610aaf611128565b80156106ef576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60606066805461041c90612065565b6107da33838361117e565b610b1b3383610cf5565b610b375760405162461bcd60e51b815260040161053b906120a0565b610b438484848461124d565b50505050565b606061040782611280565b610b5c61104b565b6000610b6861015f5490565b9050610b7961015f80546001019055565b610b8383826113a0565b6105dc81836113ba565b610b9561104b565b6001600160a01b038116610bfa5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161053b565b6106ef816110a5565b60006001600160e01b03198216632483248360e11b148061040757506104078261148c565b6000818152606760205260409020546001600160a01b03166106ef5760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b604482015260640161053b565b600081815260696020526040902080546001600160a01b0319166001600160a01b0384169081179091558190610cbc82610891565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b600080610d0183610891565b9050806001600160a01b0316846001600160a01b03161480610d4857506001600160a01b038082166000908152606a602090815260408083209388168352929052205460ff165b80610d6c5750836001600160a01b0316610d618461049f565b6001600160a01b0316145b949350505050565b826001600160a01b0316610d8782610891565b6001600160a01b031614610dad5760405162461bcd60e51b815260040161053b90612185565b6001600160a01b038216610e0f5760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b606482015260840161053b565b826001600160a01b0316610e2282610891565b6001600160a01b031614610e485760405162461bcd60e51b815260040161053b90612185565b600081815260696020908152604080832080546001600160a01b03199081169091556001600160a01b0387811680865260688552838620805460001901905590871680865283862080546001019055868652606790945282852080549092168417909155905184937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4505050565b6106ef61104b565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610f13576105dc836114dc565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015610f6d575060408051601f3d908101601f19168201909252610f6a918101906121ca565b60015b610fd05760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b606482015260840161053b565b60008051602061231c833981519152811461103f5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b606482015260840161053b565b506105dc838383611578565b60c9546001600160a01b031633146109895760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161053b565b60c980546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff1661111e5760405162461bcd60e51b815260040161053b906121e3565b6107da828261159d565b600054610100900460ff166109895760405162461bcd60e51b815260040161053b906121e3565b600054610100900460ff166111765760405162461bcd60e51b815260040161053b906121e3565b6109896115eb565b816001600160a01b0316836001600160a01b031614156111e05760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c657200000000000000604482015260640161053b565b6001600160a01b038381166000818152606a6020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b611258848484610d74565b6112648484848461161b565b610b435760405162461bcd60e51b815260040161053b9061222e565b606061128b82610c28565b600082815260976020526040812080546112a490612065565b80601f01602080910402602001604051908101604052809291908181526020018280546112d090612065565b801561131d5780601f106112f25761010080835404028352916020019161131d565b820191906000526020600020905b81548152906001019060200180831161130057829003601f168201915b50505050509050600061135260408051808201909152600f81526e1a5c199cce8bcbdd195cdd0b5d5c9b608a1b602082015290565b9050805160001415611365575092915050565b81511561139757808260405160200161137f929190612280565b60405160208183030381529060405292505050919050565b610d6c84611719565b6107da8282604051806020016040528060008152506117a4565b6000828152606760205260409020546001600160a01b03166114355760405162461bcd60e51b815260206004820152602e60248201527f45524337323155524953746f726167653a2055524920736574206f66206e6f6e60448201526d32bc34b9ba32b73a103a37b5b2b760911b606482015260840161053b565b6000828152609760209081526040909120825161145492840190611c4f565b506040518281527ff8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce79060200160405180910390a15050565b60006001600160e01b031982166380ac58cd60e01b14806114bd57506001600160e01b03198216635b5e139f60e01b145b8061040757506301ffc9a760e01b6001600160e01b0319831614610407565b6001600160a01b0381163b6115495760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840161053b565b60008051602061231c83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611581836117d7565b60008251118061158e5750805b156105dc57610b438383611817565b600054610100900460ff166115c45760405162461bcd60e51b815260040161053b906121e3565b81516115d7906065906020850190611c4f565b5080516105dc906066906020840190611c4f565b600054610100900460ff166116125760405162461bcd60e51b815260040161053b906121e3565b610989336110a5565b60006001600160a01b0384163b1561170e57604051630a85bd0160e11b81526001600160a01b0385169063150b7a029061165f9033908990889088906004016122af565b6020604051808303816000875af192505050801561169a575060408051601f3d908101601f19168201909252611697918101906122e2565b60015b6116f4573d8080156116c8576040519150601f19603f3d011682016040523d82523d6000602084013e6116cd565b606091505b5080516116ec5760405162461bcd60e51b815260040161053b9061222e565b805181602001fd5b6001600160e01b031916630a85bd0160e11b149050610d6c565b506001949350505050565b606061172482610c28565b600061175260408051808201909152600f81526e1a5c199cce8bcbdd195cdd0b5d5c9b608a1b602082015290565b90506000815111611772576040518060200160405280600081525061179d565b8061177c8461183c565b60405160200161178d929190612280565b6040516020818303038152906040525b9392505050565b6117ae83836118d9565b6117bb600084848461161b565b6105dc5760405162461bcd60e51b815260040161053b9061222e565b6117e0816114dc565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061179d838360405180606001604052806027815260200161233c60279139611a64565b6060600061184983611adc565b600101905060008167ffffffffffffffff81111561186957611869611e3c565b6040519080825280601f01601f191660200182016040528015611893576020820181803683370190505b5090508181016020015b600019016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846118cc576118d1565b61189d565b509392505050565b6001600160a01b03821661192f5760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f2061646472657373604482015260640161053b565b6000818152606760205260409020546001600160a01b0316156119945760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e74656400000000604482015260640161053b565b6000818152606760205260409020546001600160a01b0316156119f95760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e74656400000000604482015260640161053b565b6001600160a01b038216600081815260686020908152604080832080546001019055848352606790915280822080546001600160a01b0319168417905551839291907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b6060600080856001600160a01b031685604051611a8191906122ff565b600060405180830381855af49150503d8060008114611abc576040519150601f19603f3d011682016040523d82523d6000602084013e611ac1565b606091505b5091509150611ad286838387611bb4565b9695505050505050565b60008072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b8310611b1b5772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef81000000008310611b47576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc100008310611b6557662386f26fc10000830492506010015b6305f5e1008310611b7d576305f5e100830492506008015b6127108310611b9157612710830492506004015b60648310611ba3576064830492506002015b600a83106104075760010192915050565b60608315611c20578251611c19576001600160a01b0385163b611c195760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161053b565b5081610d6c565b610d6c8383815115611c355781518083602001fd5b8060405162461bcd60e51b815260040161053b9190611d73565b828054611c5b90612065565b90600052602060002090601f016020900481019282611c7d5760008555611cc3565b82601f10611c9657805160ff1916838001178555611cc3565b82800160010185558215611cc3579182015b82811115611cc3578251825591602001919060010190611ca8565b50611ccf929150611cd3565b5090565b5b80821115611ccf5760008155600101611cd4565b6001600160e01b0319811681146106ef57600080fd5b600060208284031215611d1057600080fd5b813561179d81611ce8565b60005b83811015611d36578181015183820152602001611d1e565b83811115610b435750506000910152565b60008151808452611d5f816020860160208601611d1b565b601f01601f19169290920160200192915050565b60208152600061179d6020830184611d47565b600060208284031215611d9857600080fd5b5035919050565b80356001600160a01b0381168114611db657600080fd5b919050565b60008060408385031215611dce57600080fd5b611dd783611d9f565b946020939093013593505050565b600080600060608486031215611dfa57600080fd5b611e0384611d9f565b9250611e1160208501611d9f565b9150604084013590509250925092565b600060208284031215611e3357600080fd5b61179d82611d9f565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff80841115611e6d57611e6d611e3c565b604051601f8501601f19908116603f01168101908282118183101715611e9557611e95611e3c565b81604052809350858152868686011115611eae57600080fd5b858560208301376000602087830101525050509392505050565b600082601f830112611ed957600080fd5b61179d83833560208501611e52565b60008060408385031215611efb57600080fd5b611f0483611d9f565b9150602083013567ffffffffffffffff811115611f2057600080fd5b611f2c85828601611ec8565b9150509250929050565b60008060408385031215611f4957600080fd5b611f5283611d9f565b915060208301358015158114611f6757600080fd5b809150509250929050565b60008060008060808587031215611f8857600080fd5b611f9185611d9f565b9350611f9f60208601611d9f565b925060408501359150606085013567ffffffffffffffff811115611fc257600080fd5b611fce87828801611ec8565b91505092959194509250565b60008060408385031215611fed57600080fd5b611ff683611d9f565b9150602083013567ffffffffffffffff81111561201257600080fd5b8301601f8101851361202357600080fd5b611f2c85823560208401611e52565b6000806040838503121561204557600080fd5b61204e83611d9f565b915061205c60208401611d9f565b90509250929050565b600181811c9082168061207957607f821691505b6020821081141561209a57634e487b7160e01b600052602260045260246000fd5b50919050565b6020808252602d908201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560408201526c1c881bdc88185c1c1c9bdd9959609a1b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b60208082526025908201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060408201526437bbb732b960d91b606082015260800190565b6000602082840312156121dc57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60208082526032908201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560408201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b606082015260800190565b60008351612292818460208801611d1b565b8351908301906122a6818360208801611d1b565b01949350505050565b6001600160a01b0385811682528416602082015260408101839052608060608201819052600090611ad290830184611d47565b6000602082840312156122f457600080fd5b815161179d81611ce8565b60008251612311818460208701611d1b565b919091019291505056fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220669c6943431236e79addb2ca9dbff0fbcd604df8d4f6116c23f4ad669d66a05664736f6c634300080a0033",
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

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ERC721TokenTest *ERC721TokenTestCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ERC721TokenTest.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ERC721TokenTest *ERC721TokenTestSession) ProxiableUUID() ([32]byte, error) {
	return _ERC721TokenTest.Contract.ProxiableUUID(&_ERC721TokenTest.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ERC721TokenTest *ERC721TokenTestCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ERC721TokenTest.Contract.ProxiableUUID(&_ERC721TokenTest.CallOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ERC721TokenTest *ERC721TokenTestSession) Initialize() (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Initialize(&_ERC721TokenTest.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) Initialize() (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.Initialize(&_ERC721TokenTest.TransactOpts)
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

// SafeMint is a paid mutator transaction binding the contract method 0xd204c45e.
//
// Solidity: function safeMint(address to, string uri) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) SafeMint(opts *bind.TransactOpts, to common.Address, uri string) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "safeMint", to, uri)
}

// SafeMint is a paid mutator transaction binding the contract method 0xd204c45e.
//
// Solidity: function safeMint(address to, string uri) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) SafeMint(to common.Address, uri string) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SafeMint(&_ERC721TokenTest.TransactOpts, to, uri)
}

// SafeMint is a paid mutator transaction binding the contract method 0xd204c45e.
//
// Solidity: function safeMint(address to, string uri) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) SafeMint(to common.Address, uri string) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.SafeMint(&_ERC721TokenTest.TransactOpts, to, uri)
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

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ERC721TokenTest *ERC721TokenTestSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.UpgradeTo(&_ERC721TokenTest.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.UpgradeTo(&_ERC721TokenTest.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ERC721TokenTest *ERC721TokenTestTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ERC721TokenTest.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ERC721TokenTest *ERC721TokenTestSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.UpgradeToAndCall(&_ERC721TokenTest.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ERC721TokenTest *ERC721TokenTestTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ERC721TokenTest.Contract.UpgradeToAndCall(&_ERC721TokenTest.TransactOpts, newImplementation, data)
}

// ERC721TokenTestAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC721TokenTest contract.
type ERC721TokenTestAdminChangedIterator struct {
	Event *ERC721TokenTestAdminChanged // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestAdminChanged)
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
		it.Event = new(ERC721TokenTestAdminChanged)
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
func (it *ERC721TokenTestAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestAdminChanged represents a AdminChanged event raised by the ERC721TokenTest contract.
type ERC721TokenTestAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC721TokenTestAdminChangedIterator, error) {

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestAdminChangedIterator{contract: _ERC721TokenTest.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestAdminChanged)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseAdminChanged(log types.Log) (*ERC721TokenTestAdminChanged, error) {
	event := new(ERC721TokenTestAdminChanged)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

// ERC721TokenTestBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the ERC721TokenTest contract.
type ERC721TokenTestBatchMetadataUpdateIterator struct {
	Event *ERC721TokenTestBatchMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestBatchMetadataUpdate)
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
		it.Event = new(ERC721TokenTestBatchMetadataUpdate)
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
func (it *ERC721TokenTestBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the ERC721TokenTest contract.
type ERC721TokenTestBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*ERC721TokenTestBatchMetadataUpdateIterator, error) {

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestBatchMetadataUpdateIterator{contract: _ERC721TokenTest.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestBatchMetadataUpdate)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
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

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseBatchMetadataUpdate(log types.Log) (*ERC721TokenTestBatchMetadataUpdate, error) {
	event := new(ERC721TokenTestBatchMetadataUpdate)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TokenTestBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC721TokenTest contract.
type ERC721TokenTestBeaconUpgradedIterator struct {
	Event *ERC721TokenTestBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestBeaconUpgraded)
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
		it.Event = new(ERC721TokenTestBeaconUpgraded)
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
func (it *ERC721TokenTestBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestBeaconUpgraded represents a BeaconUpgraded event raised by the ERC721TokenTest contract.
type ERC721TokenTestBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC721TokenTestBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestBeaconUpgradedIterator{contract: _ERC721TokenTest.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestBeaconUpgraded)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseBeaconUpgraded(log types.Log) (*ERC721TokenTestBeaconUpgraded, error) {
	event := new(ERC721TokenTestBeaconUpgraded)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TokenTestInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ERC721TokenTest contract.
type ERC721TokenTestInitializedIterator struct {
	Event *ERC721TokenTestInitialized // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestInitialized)
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
		it.Event = new(ERC721TokenTestInitialized)
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
func (it *ERC721TokenTestInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestInitialized represents a Initialized event raised by the ERC721TokenTest contract.
type ERC721TokenTestInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterInitialized(opts *bind.FilterOpts) (*ERC721TokenTestInitializedIterator, error) {

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestInitializedIterator{contract: _ERC721TokenTest.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestInitialized) (event.Subscription, error) {

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestInitialized)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseInitialized(log types.Log) (*ERC721TokenTestInitialized, error) {
	event := new(ERC721TokenTestInitialized)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC721TokenTestMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the ERC721TokenTest contract.
type ERC721TokenTestMetadataUpdateIterator struct {
	Event *ERC721TokenTestMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestMetadataUpdate)
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
		it.Event = new(ERC721TokenTestMetadataUpdate)
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
func (it *ERC721TokenTestMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestMetadataUpdate represents a MetadataUpdate event raised by the ERC721TokenTest contract.
type ERC721TokenTestMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*ERC721TokenTestMetadataUpdateIterator, error) {

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestMetadataUpdateIterator{contract: _ERC721TokenTest.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestMetadataUpdate)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseMetadataUpdate(log types.Log) (*ERC721TokenTestMetadataUpdate, error) {
	event := new(ERC721TokenTestMetadataUpdate)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ERC721TokenTestUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC721TokenTest contract.
type ERC721TokenTestUpgradedIterator struct {
	Event *ERC721TokenTestUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC721TokenTestUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC721TokenTestUpgraded)
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
		it.Event = new(ERC721TokenTestUpgraded)
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
func (it *ERC721TokenTestUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC721TokenTestUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC721TokenTestUpgraded represents a Upgraded event raised by the ERC721TokenTest contract.
type ERC721TokenTestUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC721TokenTest *ERC721TokenTestFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC721TokenTestUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC721TokenTestUpgradedIterator{contract: _ERC721TokenTest.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC721TokenTest *ERC721TokenTestFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC721TokenTestUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC721TokenTest.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC721TokenTestUpgraded)
				if err := _ERC721TokenTest.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ERC721TokenTest *ERC721TokenTestFilterer) ParseUpgraded(log types.Log) (*ERC721TokenTestUpgraded, error) {
	event := new(ERC721TokenTestUpgraded)
	if err := _ERC721TokenTest.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
