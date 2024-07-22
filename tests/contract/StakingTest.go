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

// StakingTestMetaData contains all meta data concerning the StakingTest contract.
var StakingTestMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"ApproveShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valSrc\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valDst\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Redelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valSrc\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"valDst\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"RedelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"TransferShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Undelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"UndelegateV2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowanceShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approveShares\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"delegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_valSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_valDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"redelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_valSrc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_valDst\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"redelegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFromShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"undelegateV2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"validatorShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611755806100206000396000f3fe6080604052600436106100dd5760003560e01c80638dfc88971161007f578063d5c498eb11610059578063d5c498eb14610280578063dc6ffc7d146102a0578063de2b3451146102c0578063ee226c66146102e057600080fd5b80638dfc8897146102155780639ddb511a14610235578063bf98d7721461024857600080fd5b806351af513a116100bb57806351af513a1461017a5780636d7880351461019a5780637b625c0f146101ba5780637dd0209d146101da57600080fd5b8063161298c1146100e257806331fb67c21461011c57806349da433e1461014a575b600080fd5b3480156100ee57600080fd5b506101026100fd3660046111ae565b610300565b604080519283526020830191909152015b60405180910390f35b34801561012857600080fd5b5061013c610137366004611205565b61031f565b604051908152602001610113565b34801561015657600080fd5b5061016a6101653660046111ae565b610332565b6040519015158152602001610113565b34801561018657600080fd5b5061013c61019536600461123a565b610349565b3480156101a657600080fd5b5061016a6101b5366004611288565b610355565b3480156101c657600080fd5b5061013c6101d53660046112cd565b61040d565b3480156101e657600080fd5b506101fa6101f536600461132b565b610422565b60408051938452602084019290925290820152606001610113565b34801561022157600080fd5b506101fa610230366004611288565b6104b2565b610102610243366004611205565b61050b565b34801561025457600080fd5b5061013c610263366004611205565b805160208183018101805160008252928201919093012091525481565b34801561028c57600080fd5b5061010261029b36600461123a565b61055c565b3480156102ac57600080fd5b506101026102bb366004611398565b610574565b3480156102cc57600080fd5b5061016a6102db366004611288565b610595565b3480156102ec57600080fd5b5061016a6102fb36600461132b565b6105bf565b60008060008061031187878761062e565b909890975095505050505050565b60008061032b836106e5565b9392505050565b600080610340858585610784565b95945050505050565b600061032b8383610835565b6000814710156103a35760405162461bcd60e51b8152602060048201526014602482015273696e73756666696369656e742062616c616e636560601b60448201526064015b60405180910390fd5b604051636d78803560e01b815261100390636d788035906103ca908690869060040161145a565b6020604051808303816000875af11580156103e9573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032b919061147c565b600061041a8484846108e1565b949350505050565b60008060008060008061043689898961098e565b9250925092508660008a60405161044d919061149e565b9081526020016040518091039020600082825461046a91906114d0565b9250508190555086600089604051610482919061149e565b9081526020016040518091039020600082825461049f91906114e7565b9091555092999198509650945050505050565b6000806000806000806104c58888610a44565b925092509250866000896040516104dc919061149e565b908152602001604051809103902060008282546104f991906114d0565b90915550929891975095509350505050565b60008060008061051b8534610af7565b9150915081600086604051610530919061149e565b9081526020016040518091039020600082825461054d91906114e7565b90915550919590945092505050565b6000806105698484610ba5565b915091509250929050565b60008060008061058688888888610c3d565b90999098509650505050505050565b60405163de2b345160e01b81526000906110039063de2b3451906103ca908690869060040161145a565b604051637711363360e11b81526000906110039063ee226c66906105eb908790879087906004016114ff565b6020604051808303816000875af115801561060a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041a919061147c565b6000808080611003610641888888610cfd565b60405161064e919061149e565b6000604051808303816000865af19150503d806000811461068b576040519150601f19603f3d011682016040523d82523d6000602084013e610690565b606091505b50915091506106ce8282604051806040016040528060168152602001751d1c985b9cd9995c881cda185c995cc819985a5b195960521b815250610d47565b6106d781610dc6565b935093505050935093915050565b600080806110036106f585610dec565b604051610702919061149e565b6000604051808303816000865af19150503d806000811461073f576040519150601f19603f3d011682016040523d82523d6000602084013e610744565b606091505b509150915061077b82826040518060400160405280600f81526020016e1dda5d1a191c985dc819985a5b1959608a1b815250610d47565b61041a81610e2f565b60008080611003610796878787610e46565b6040516107a3919061149e565b6000604051808303816000865af19150503d80600081146107e0576040519150601f19603f3d011682016040523d82523d6000602084013e6107e5565b606091505b5091509150610822828260405180604001604052806015815260200174185c1c1c9bdd99481cda185c995cc819985a5b1959605a1b815250610d47565b61082b81610e90565b9695505050505050565b600080806110036108468686610ea7565b604051610853919061149e565b600060405180830381855afa9150503d806000811461088e576040519150601f19603f3d011682016040523d82523d6000602084013e610893565b606091505b50915091506108d882826040518060400160405280601881526020017f64656c65676174696f6e52657761726473206661696c65640000000000000000815250610d47565b61034081610e2f565b600080806110036108f3878787610eee565b604051610900919061149e565b600060405180830381855afa9150503d806000811461093b576040519150601f19603f3d011682016040523d82523d6000602084013e610940565b606091505b509150915061098582826040518060400160405280601781526020017f616c6c6f77616e636520736861726573206661696c6564000000000000000000815250610d47565b61082b81610e2f565b6000808080806110036109a2898989610f38565b6040516109af919061149e565b6000604051808303816000865af19150503d80600081146109ec576040519150601f19603f3d011682016040523d82523d6000602084013e6109f1565b606091505b5091509150610a2a8282604051806040016040528060118152602001701c9959195b1959d85d194819985a5b1959607a1b815250610d47565b610a3381610f82565b945094509450505093509350939050565b600080808080611003610a578888610faf565b604051610a64919061149e565b6000604051808303816000865af19150503d8060008114610aa1576040519150601f19603f3d011682016040523d82523d6000602084013e610aa6565b606091505b5091509150610adf8282604051806040016040528060118152602001701d5b99195b1959d85d194819985a5b1959607a1b815250610d47565b610ae881610f82565b94509450945050509250925092565b600080808061100385610b0988610ff6565b604051610b16919061149e565b60006040518083038185875af1925050503d8060008114610b53576040519150601f19603f3d011682016040523d82523d6000602084013e610b58565b606091505b5091509150610b8f82826040518060400160405280600f81526020016e19195b1959d85d194819985a5b1959608a1b815250610d47565b610b9881610dc6565b9350935050509250929050565b6000808080611003610bb78787611039565b604051610bc4919061149e565b600060405180830381855afa9150503d8060008114610bff576040519150601f19603f3d011682016040523d82523d6000602084013e610c04565b606091505b5091509150610b8f82826040518060400160405280601181526020017019195b1959d85d1a5bdb8819985a5b1959607a1b815250610d47565b6000808080611003610c5189898989611080565b604051610c5e919061149e565b6000604051808303816000865af19150503d8060008114610c9b576040519150601f19603f3d011682016040523d82523d6000602084013e610ca0565b606091505b5091509150610ce582826040518060400160405280601a81526020017f7472616e7366657246726f6d20736861726573206661696c6564000000000000815250610d47565b610cee81610dc6565b93509350505094509492505050565b6060838383604051602401610d1493929190611535565b60408051601f198184030181529190526020810180516001600160e01b031663161298c160e01b17905290509392505050565b82610dc157600082806020019051810190610d629190611563565b9050600182511015610d88578060405162461bcd60e51b815260040161039a91906115d1565b8181604051602001610d9b9291906115e4565b60408051601f198184030181529082905262461bcd60e51b825261039a916004016115d1565b505050565b60008060008084806020019051810190610de09190611621565b90969095509350505050565b606081604051602401610dff91906115d1565b60408051601f198184030181529190526020810180516001600160e01b03166318fdb3e160e11b17905292915050565b6000808280602001905181019061032b9190611645565b6060838383604051602401610e5d93929190611535565b60408051601f198184030181529190526020810180516001600160e01b03166324ed219f60e11b17905290509392505050565b6000808280602001905181019061032b919061147c565b60608282604051602401610ebc92919061165e565b60408051601f198184030181529190526020810180516001600160e01b03166328d7a89d60e11b179052905092915050565b6060838383604051602401610f0593929190611688565b60408051601f198184030181529190526020810180516001600160e01b0316637b625c0f60e01b17905290509392505050565b6060838383604051602401610f4f939291906114ff565b60408051601f198184030181529190526020810180516001600160e01b0316637dd0209d60e01b17905290509392505050565b60008060008060008086806020019051810190610f9f91906116bb565b9199909850909650945050505050565b60608282604051602401610fc492919061145a565b60408051601f198184030181529190526020810180516001600160e01b0316638dfc889760e01b179052905092915050565b60608160405160240161100991906115d1565b60408051601f198184030181529190526020810180516001600160e01b0316634eeda88d60e11b17905292915050565b6060828260405160240161104e92919061165e565b60408051601f198184030181529190526020810180516001600160e01b031663d5c498eb60e01b179052905092915050565b60608484848460405160240161109994939291906116e9565b60408051601f198184030181529190526020810180516001600160e01b031663dc6ffc7d60e01b1790529050949350505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561110c5761110c6110cd565b604052919050565b600067ffffffffffffffff82111561112e5761112e6110cd565b50601f01601f191660200190565b600082601f83011261114d57600080fd5b813561116061115b82611114565b6110e3565b81815284602083860101111561117557600080fd5b816020850160208301376000918101602001919091529392505050565b80356001600160a01b03811681146111a957600080fd5b919050565b6000806000606084860312156111c357600080fd5b833567ffffffffffffffff8111156111da57600080fd5b6111e68682870161113c565b9350506111f560208501611192565b9150604084013590509250925092565b60006020828403121561121757600080fd5b813567ffffffffffffffff81111561122e57600080fd5b61041a8482850161113c565b6000806040838503121561124d57600080fd5b823567ffffffffffffffff81111561126457600080fd5b6112708582860161113c565b92505061127f60208401611192565b90509250929050565b6000806040838503121561129b57600080fd5b823567ffffffffffffffff8111156112b257600080fd5b6112be8582860161113c565b95602094909401359450505050565b6000806000606084860312156112e257600080fd5b833567ffffffffffffffff8111156112f957600080fd5b6113058682870161113c565b93505061131460208501611192565b915061132260408501611192565b90509250925092565b60008060006060848603121561134057600080fd5b833567ffffffffffffffff8082111561135857600080fd5b6113648783880161113c565b9450602086013591508082111561137a57600080fd5b506113878682870161113c565b925050604084013590509250925092565b600080600080608085870312156113ae57600080fd5b843567ffffffffffffffff8111156113c557600080fd5b6113d18782880161113c565b9450506113e060208601611192565b92506113ee60408601611192565b9396929550929360600135925050565b60005b83811015611419578181015183820152602001611401565b83811115611428576000848401525b50505050565b600081518084526114468160208601602086016113fe565b601f01601f19169290920160200192915050565b60408152600061146d604083018561142e565b90508260208301529392505050565b60006020828403121561148e57600080fd5b8151801515811461032b57600080fd5b600082516114b08184602087016113fe565b9190910192915050565b634e487b7160e01b600052601160045260246000fd5b6000828210156114e2576114e26114ba565b500390565b600082198211156114fa576114fa6114ba565b500190565b606081526000611512606083018661142e565b8281036020840152611524818661142e565b915050826040830152949350505050565b606081526000611548606083018661142e565b6001600160a01b039490941660208301525060400152919050565b60006020828403121561157557600080fd5b815167ffffffffffffffff81111561158c57600080fd5b8201601f8101841361159d57600080fd5b80516115ab61115b82611114565b8181528560208385010111156115c057600080fd5b6103408260208301602086016113fe565b60208152600061032b602083018461142e565b600083516115f68184602088016113fe565b6101d160f51b90830190815283516116158160028401602088016113fe565b01600201949350505050565b6000806040838503121561163457600080fd5b505080516020909101519092909150565b60006020828403121561165757600080fd5b5051919050565b604081526000611671604083018561142e565b905060018060a01b03831660208301529392505050565b60608152600061169b606083018661142e565b6001600160a01b0394851660208401529290931660409091015292915050565b6000806000606084860312156116d057600080fd5b8351925060208401519150604084015190509250925092565b6080815260006116fc608083018761142e565b6001600160a01b039586166020840152939094166040820152606001529291505056fea2646970667358221220650c5a9c36b440e635e7d03c16fb3e5d0d9643687ad0f6ce31981e819c7ea04064736f6c634300080a0033",
}

// StakingTestABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingTestMetaData.ABI instead.
var StakingTestABI = StakingTestMetaData.ABI

// StakingTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingTestMetaData.Bin instead.
var StakingTestBin = StakingTestMetaData.Bin

// DeployStakingTest deploys a new Ethereum contract, binding an instance of StakingTest to it.
func DeployStakingTest(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakingTest, error) {
	parsed, err := StakingTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingTestBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingTest{StakingTestCaller: StakingTestCaller{contract: contract}, StakingTestTransactor: StakingTestTransactor{contract: contract}, StakingTestFilterer: StakingTestFilterer{contract: contract}}, nil
}

// StakingTest is an auto generated Go binding around an Ethereum contract.
type StakingTest struct {
	StakingTestCaller     // Read-only binding to the contract
	StakingTestTransactor // Write-only binding to the contract
	StakingTestFilterer   // Log filterer for contract events
}

// StakingTestCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingTestCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTestTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingTestFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTestSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingTestSession struct {
	Contract     *StakingTest      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingTestCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingTestCallerSession struct {
	Contract *StakingTestCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakingTestTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTestTransactorSession struct {
	Contract     *StakingTestTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakingTestRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingTestRaw struct {
	Contract *StakingTest // Generic contract binding to access the raw methods on
}

// StakingTestCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingTestCallerRaw struct {
	Contract *StakingTestCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTestTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTestTransactorRaw struct {
	Contract *StakingTestTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingTest creates a new instance of StakingTest, bound to a specific deployed contract.
func NewStakingTest(address common.Address, backend bind.ContractBackend) (*StakingTest, error) {
	contract, err := bindStakingTest(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingTest{StakingTestCaller: StakingTestCaller{contract: contract}, StakingTestTransactor: StakingTestTransactor{contract: contract}, StakingTestFilterer: StakingTestFilterer{contract: contract}}, nil
}

// NewStakingTestCaller creates a new read-only instance of StakingTest, bound to a specific deployed contract.
func NewStakingTestCaller(address common.Address, caller bind.ContractCaller) (*StakingTestCaller, error) {
	contract, err := bindStakingTest(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTestCaller{contract: contract}, nil
}

// NewStakingTestTransactor creates a new write-only instance of StakingTest, bound to a specific deployed contract.
func NewStakingTestTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTestTransactor, error) {
	contract, err := bindStakingTest(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTestTransactor{contract: contract}, nil
}

// NewStakingTestFilterer creates a new log filterer instance of StakingTest, bound to a specific deployed contract.
func NewStakingTestFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingTestFilterer, error) {
	contract, err := bindStakingTest(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingTestFilterer{contract: contract}, nil
}

// bindStakingTest binds a generic wrapper to an already deployed contract.
func bindStakingTest(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingTestMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingTest *StakingTestRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingTest.Contract.StakingTestCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingTest *StakingTestRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingTest.Contract.StakingTestTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingTest *StakingTestRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingTest.Contract.StakingTestTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingTest *StakingTestCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingTest.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingTest *StakingTestTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingTest.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingTest *StakingTestTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingTest.Contract.contract.Transact(opts, method, params...)
}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestCaller) AllowanceShares(opts *bind.CallOpts, _val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "allowanceShares", _val, _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestSession) AllowanceShares(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StakingTest.Contract.AllowanceShares(&_StakingTest.CallOpts, _val, _owner, _spender)
}

// AllowanceShares is a free data retrieval call binding the contract method 0x7b625c0f.
//
// Solidity: function allowanceShares(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) AllowanceShares(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StakingTest.Contract.AllowanceShares(&_StakingTest.CallOpts, _val, _owner, _spender)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestCaller) Delegation(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "delegation", _val, _del)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Delegation(_val string, _del common.Address) (*big.Int, *big.Int, error) {
	return _StakingTest.Contract.Delegation(&_StakingTest.CallOpts, _val, _del)
}

// Delegation is a free data retrieval call binding the contract method 0xd5c498eb.
//
// Solidity: function delegation(string _val, address _del) view returns(uint256, uint256)
func (_StakingTest *StakingTestCallerSession) Delegation(_val string, _del common.Address) (*big.Int, *big.Int, error) {
	return _StakingTest.Contract.Delegation(&_StakingTest.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestCaller) DelegationRewards(opts *bind.CallOpts, _val string, _del common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "delegationRewards", _val, _del)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _StakingTest.Contract.DelegationRewards(&_StakingTest.CallOpts, _val, _del)
}

// DelegationRewards is a free data retrieval call binding the contract method 0x51af513a.
//
// Solidity: function delegationRewards(string _val, address _del) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) DelegationRewards(_val string, _del common.Address) (*big.Int, error) {
	return _StakingTest.Contract.DelegationRewards(&_StakingTest.CallOpts, _val, _del)
}

// ValidatorShares is a free data retrieval call binding the contract method 0xbf98d772.
//
// Solidity: function validatorShares(string ) view returns(uint256)
func (_StakingTest *StakingTestCaller) ValidatorShares(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "validatorShares", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorShares is a free data retrieval call binding the contract method 0xbf98d772.
//
// Solidity: function validatorShares(string ) view returns(uint256)
func (_StakingTest *StakingTestSession) ValidatorShares(arg0 string) (*big.Int, error) {
	return _StakingTest.Contract.ValidatorShares(&_StakingTest.CallOpts, arg0)
}

// ValidatorShares is a free data retrieval call binding the contract method 0xbf98d772.
//
// Solidity: function validatorShares(string ) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) ValidatorShares(arg0 string) (*big.Int, error) {
	return _StakingTest.Contract.ValidatorShares(&_StakingTest.CallOpts, arg0)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestTransactor) ApproveShares(opts *bind.TransactOpts, _val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "approveShares", _val, _spender, _shares)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestSession) ApproveShares(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.ApproveShares(&_StakingTest.TransactOpts, _val, _spender, _shares)
}

// ApproveShares is a paid mutator transaction binding the contract method 0x49da433e.
//
// Solidity: function approveShares(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestTransactorSession) ApproveShares(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.ApproveShares(&_StakingTest.TransactOpts, _val, _spender, _shares)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) Delegate(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "delegate", _val)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Delegate(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Delegate(&_StakingTest.TransactOpts, _val)
}

// Delegate is a paid mutator transaction binding the contract method 0x9ddb511a.
//
// Solidity: function delegate(string _val) payable returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Delegate(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Delegate(&_StakingTest.TransactOpts, _val)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactor) DelegateV2(opts *bind.TransactOpts, _val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "delegateV2", _val, _amount)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestSession) DelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.DelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// DelegateV2 is a paid mutator transaction binding the contract method 0x6d788035.
//
// Solidity: function delegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactorSession) DelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.DelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// Redelegate is a paid mutator transaction binding the contract method 0x7dd0209d.
//
// Solidity: function redelegate(string _valSrc, string _valDst, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestTransactor) Redelegate(opts *bind.TransactOpts, _valSrc string, _valDst string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "redelegate", _valSrc, _valDst, _shares)
}

// Redelegate is a paid mutator transaction binding the contract method 0x7dd0209d.
//
// Solidity: function redelegate(string _valSrc, string _valDst, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestSession) Redelegate(_valSrc string, _valDst string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Redelegate(&_StakingTest.TransactOpts, _valSrc, _valDst, _shares)
}

// Redelegate is a paid mutator transaction binding the contract method 0x7dd0209d.
//
// Solidity: function redelegate(string _valSrc, string _valDst, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Redelegate(_valSrc string, _valDst string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Redelegate(&_StakingTest.TransactOpts, _valSrc, _valDst, _shares)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactor) RedelegateV2(opts *bind.TransactOpts, _valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "redelegateV2", _valSrc, _valDst, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestSession) RedelegateV2(_valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.RedelegateV2(&_StakingTest.TransactOpts, _valSrc, _valDst, _amount)
}

// RedelegateV2 is a paid mutator transaction binding the contract method 0xee226c66.
//
// Solidity: function redelegateV2(string _valSrc, string _valDst, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactorSession) RedelegateV2(_valSrc string, _valDst string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.RedelegateV2(&_StakingTest.TransactOpts, _valSrc, _valDst, _amount)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) TransferFromShares(opts *bind.TransactOpts, _val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "transferFromShares", _val, _from, _to, _shares)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) TransferFromShares(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferFromShares(&_StakingTest.TransactOpts, _val, _from, _to, _shares)
}

// TransferFromShares is a paid mutator transaction binding the contract method 0xdc6ffc7d.
//
// Solidity: function transferFromShares(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) TransferFromShares(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferFromShares(&_StakingTest.TransactOpts, _val, _from, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) TransferShares(opts *bind.TransactOpts, _val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "transferShares", _val, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) TransferShares(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferShares(&_StakingTest.TransactOpts, _val, _to, _shares)
}

// TransferShares is a paid mutator transaction binding the contract method 0x161298c1.
//
// Solidity: function transferShares(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) TransferShares(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferShares(&_StakingTest.TransactOpts, _val, _to, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestTransactor) Undelegate(opts *bind.TransactOpts, _val string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "undelegate", _val, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestSession) Undelegate(_val string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Undelegate(&_StakingTest.TransactOpts, _val, _shares)
}

// Undelegate is a paid mutator transaction binding the contract method 0x8dfc8897.
//
// Solidity: function undelegate(string _val, uint256 _shares) returns(uint256, uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Undelegate(_val string, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Undelegate(&_StakingTest.TransactOpts, _val, _shares)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactor) UndelegateV2(opts *bind.TransactOpts, _val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "undelegateV2", _val, _amount)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestSession) UndelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.UndelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// UndelegateV2 is a paid mutator transaction binding the contract method 0xde2b3451.
//
// Solidity: function undelegateV2(string _val, uint256 _amount) returns(bool _result)
func (_StakingTest *StakingTestTransactorSession) UndelegateV2(_val string, _amount *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.UndelegateV2(&_StakingTest.TransactOpts, _val, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256)
func (_StakingTest *StakingTestTransactor) Withdraw(opts *bind.TransactOpts, _val string) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "withdraw", _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256)
func (_StakingTest *StakingTestSession) Withdraw(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Withdraw(&_StakingTest.TransactOpts, _val)
}

// Withdraw is a paid mutator transaction binding the contract method 0x31fb67c2.
//
// Solidity: function withdraw(string _val) returns(uint256)
func (_StakingTest *StakingTestTransactorSession) Withdraw(_val string) (*types.Transaction, error) {
	return _StakingTest.Contract.Withdraw(&_StakingTest.TransactOpts, _val)
}

// StakingTestApproveSharesIterator is returned from FilterApproveShares and is used to iterate over the raw logs and unpacked data for ApproveShares events raised by the StakingTest contract.
type StakingTestApproveSharesIterator struct {
	Event *StakingTestApproveShares // Event containing the contract specifics and raw log

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
func (it *StakingTestApproveSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestApproveShares)
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
		it.Event = new(StakingTestApproveShares)
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
func (it *StakingTestApproveSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestApproveSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestApproveShares represents a ApproveShares event raised by the StakingTest contract.
type StakingTestApproveShares struct {
	Owner     common.Address
	Spender   common.Address
	Validator string
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterApproveShares is a free log retrieval operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) FilterApproveShares(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StakingTestApproveSharesIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "ApproveShares", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestApproveSharesIterator{contract: _StakingTest.contract, event: "ApproveShares", logs: logs, sub: sub}, nil
}

// WatchApproveShares is a free log subscription operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) WatchApproveShares(opts *bind.WatchOpts, sink chan<- *StakingTestApproveShares, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "ApproveShares", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestApproveShares)
				if err := _StakingTest.contract.UnpackLog(event, "ApproveShares", log); err != nil {
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

// ParseApproveShares is a log parse operation binding the contract event 0xbd99ef1c86c593a90a79f794ca07759c5a04cf54bf800cfb77bb0b9fdb9bc04a.
//
// Solidity: event ApproveShares(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) ParseApproveShares(log types.Log) (*StakingTestApproveShares, error) {
	event := new(StakingTestApproveShares)
	if err := _StakingTest.contract.UnpackLog(event, "ApproveShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the StakingTest contract.
type StakingTestDelegateIterator struct {
	Event *StakingTestDelegate // Event containing the contract specifics and raw log

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
func (it *StakingTestDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestDelegate)
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
		it.Event = new(StakingTestDelegate)
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
func (it *StakingTestDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestDelegate represents a Delegate event raised by the StakingTest contract.
type StakingTestDelegate struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_StakingTest *StakingTestFilterer) FilterDelegate(opts *bind.FilterOpts, delegator []common.Address) (*StakingTestDelegateIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Delegate", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestDelegateIterator{contract: _StakingTest.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_StakingTest *StakingTestFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *StakingTestDelegate, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Delegate", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestDelegate)
				if err := _StakingTest.contract.UnpackLog(event, "Delegate", log); err != nil {
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

// ParseDelegate is a log parse operation binding the contract event 0x5a5adf903ba232ef17ed8be4ef872e1f60d17c5ba26a1ecbf44e388a672b118a.
//
// Solidity: event Delegate(address indexed delegator, string validator, uint256 amount, uint256 shares)
func (_StakingTest *StakingTestFilterer) ParseDelegate(log types.Log) (*StakingTestDelegate, error) {
	event := new(StakingTestDelegate)
	if err := _StakingTest.contract.UnpackLog(event, "Delegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestDelegateV2Iterator is returned from FilterDelegateV2 and is used to iterate over the raw logs and unpacked data for DelegateV2 events raised by the StakingTest contract.
type StakingTestDelegateV2Iterator struct {
	Event *StakingTestDelegateV2 // Event containing the contract specifics and raw log

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
func (it *StakingTestDelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestDelegateV2)
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
		it.Event = new(StakingTestDelegateV2)
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
func (it *StakingTestDelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestDelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestDelegateV2 represents a DelegateV2 event raised by the StakingTest contract.
type StakingTestDelegateV2 struct {
	Delegator common.Address
	Validator string
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegateV2 is a free log retrieval operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_StakingTest *StakingTestFilterer) FilterDelegateV2(opts *bind.FilterOpts, delegator []common.Address) (*StakingTestDelegateV2Iterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "DelegateV2", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestDelegateV2Iterator{contract: _StakingTest.contract, event: "DelegateV2", logs: logs, sub: sub}, nil
}

// WatchDelegateV2 is a free log subscription operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_StakingTest *StakingTestFilterer) WatchDelegateV2(opts *bind.WatchOpts, sink chan<- *StakingTestDelegateV2, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "DelegateV2", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestDelegateV2)
				if err := _StakingTest.contract.UnpackLog(event, "DelegateV2", log); err != nil {
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

// ParseDelegateV2 is a log parse operation binding the contract event 0x330852c9460e583c049d932477c038fca307363fa8c1083a332905a68b821f10.
//
// Solidity: event DelegateV2(address indexed delegator, string validator, uint256 amount)
func (_StakingTest *StakingTestFilterer) ParseDelegateV2(log types.Log) (*StakingTestDelegateV2, error) {
	event := new(StakingTestDelegateV2)
	if err := _StakingTest.contract.UnpackLog(event, "DelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestRedelegateIterator is returned from FilterRedelegate and is used to iterate over the raw logs and unpacked data for Redelegate events raised by the StakingTest contract.
type StakingTestRedelegateIterator struct {
	Event *StakingTestRedelegate // Event containing the contract specifics and raw log

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
func (it *StakingTestRedelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestRedelegate)
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
		it.Event = new(StakingTestRedelegate)
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
func (it *StakingTestRedelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestRedelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestRedelegate represents a Redelegate event raised by the StakingTest contract.
type StakingTestRedelegate struct {
	Sender         common.Address
	ValSrc         string
	ValDst         string
	Shares         *big.Int
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRedelegate is a free log retrieval operation binding the contract event 0x14e0e9558f524ca41364e4e284ebe7aabee65559c8ea32a6fca4d812e0a1d9e6.
//
// Solidity: event Redelegate(address indexed sender, string valSrc, string valDst, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) FilterRedelegate(opts *bind.FilterOpts, sender []common.Address) (*StakingTestRedelegateIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Redelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestRedelegateIterator{contract: _StakingTest.contract, event: "Redelegate", logs: logs, sub: sub}, nil
}

// WatchRedelegate is a free log subscription operation binding the contract event 0x14e0e9558f524ca41364e4e284ebe7aabee65559c8ea32a6fca4d812e0a1d9e6.
//
// Solidity: event Redelegate(address indexed sender, string valSrc, string valDst, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchRedelegate(opts *bind.WatchOpts, sink chan<- *StakingTestRedelegate, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Redelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestRedelegate)
				if err := _StakingTest.contract.UnpackLog(event, "Redelegate", log); err != nil {
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

// ParseRedelegate is a log parse operation binding the contract event 0x14e0e9558f524ca41364e4e284ebe7aabee65559c8ea32a6fca4d812e0a1d9e6.
//
// Solidity: event Redelegate(address indexed sender, string valSrc, string valDst, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) ParseRedelegate(log types.Log) (*StakingTestRedelegate, error) {
	event := new(StakingTestRedelegate)
	if err := _StakingTest.contract.UnpackLog(event, "Redelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestRedelegateV2Iterator is returned from FilterRedelegateV2 and is used to iterate over the raw logs and unpacked data for RedelegateV2 events raised by the StakingTest contract.
type StakingTestRedelegateV2Iterator struct {
	Event *StakingTestRedelegateV2 // Event containing the contract specifics and raw log

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
func (it *StakingTestRedelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestRedelegateV2)
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
		it.Event = new(StakingTestRedelegateV2)
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
func (it *StakingTestRedelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestRedelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestRedelegateV2 represents a RedelegateV2 event raised by the StakingTest contract.
type StakingTestRedelegateV2 struct {
	Sender         common.Address
	ValSrc         string
	ValDst         string
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRedelegateV2 is a free log retrieval operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) FilterRedelegateV2(opts *bind.FilterOpts, sender []common.Address) (*StakingTestRedelegateV2Iterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "RedelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestRedelegateV2Iterator{contract: _StakingTest.contract, event: "RedelegateV2", logs: logs, sub: sub}, nil
}

// WatchRedelegateV2 is a free log subscription operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchRedelegateV2(opts *bind.WatchOpts, sink chan<- *StakingTestRedelegateV2, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "RedelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestRedelegateV2)
				if err := _StakingTest.contract.UnpackLog(event, "RedelegateV2", log); err != nil {
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

// ParseRedelegateV2 is a log parse operation binding the contract event 0xdcf3a72a725100ce405b1ea62706114bec51d16536bf2cf868772ca440fe0da9.
//
// Solidity: event RedelegateV2(address indexed sender, string valSrc, string valDst, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) ParseRedelegateV2(log types.Log) (*StakingTestRedelegateV2, error) {
	event := new(StakingTestRedelegateV2)
	if err := _StakingTest.contract.UnpackLog(event, "RedelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestTransferSharesIterator is returned from FilterTransferShares and is used to iterate over the raw logs and unpacked data for TransferShares events raised by the StakingTest contract.
type StakingTestTransferSharesIterator struct {
	Event *StakingTestTransferShares // Event containing the contract specifics and raw log

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
func (it *StakingTestTransferSharesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestTransferShares)
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
		it.Event = new(StakingTestTransferShares)
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
func (it *StakingTestTransferSharesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestTransferSharesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestTransferShares represents a TransferShares event raised by the StakingTest contract.
type StakingTestTransferShares struct {
	From      common.Address
	To        common.Address
	Validator string
	Shares    *big.Int
	Token     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferShares is a free log retrieval operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) FilterTransferShares(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingTestTransferSharesIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "TransferShares", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestTransferSharesIterator{contract: _StakingTest.contract, event: "TransferShares", logs: logs, sub: sub}, nil
}

// WatchTransferShares is a free log subscription operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) WatchTransferShares(opts *bind.WatchOpts, sink chan<- *StakingTestTransferShares, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "TransferShares", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestTransferShares)
				if err := _StakingTest.contract.UnpackLog(event, "TransferShares", log); err != nil {
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

// ParseTransferShares is a log parse operation binding the contract event 0x77a2ac7846d05ab9937faf9bf901529bef4b499a2939e632f99b3fab92448344.
//
// Solidity: event TransferShares(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) ParseTransferShares(log types.Log) (*StakingTestTransferShares, error) {
	event := new(StakingTestTransferShares)
	if err := _StakingTest.contract.UnpackLog(event, "TransferShares", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestUndelegateIterator is returned from FilterUndelegate and is used to iterate over the raw logs and unpacked data for Undelegate events raised by the StakingTest contract.
type StakingTestUndelegateIterator struct {
	Event *StakingTestUndelegate // Event containing the contract specifics and raw log

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
func (it *StakingTestUndelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestUndelegate)
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
		it.Event = new(StakingTestUndelegate)
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
func (it *StakingTestUndelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestUndelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestUndelegate represents a Undelegate event raised by the StakingTest contract.
type StakingTestUndelegate struct {
	Sender         common.Address
	Validator      string
	Shares         *big.Int
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUndelegate is a free log retrieval operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) FilterUndelegate(opts *bind.FilterOpts, sender []common.Address) (*StakingTestUndelegateIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Undelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestUndelegateIterator{contract: _StakingTest.contract, event: "Undelegate", logs: logs, sub: sub}, nil
}

// WatchUndelegate is a free log subscription operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchUndelegate(opts *bind.WatchOpts, sink chan<- *StakingTestUndelegate, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Undelegate", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestUndelegate)
				if err := _StakingTest.contract.UnpackLog(event, "Undelegate", log); err != nil {
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

// ParseUndelegate is a log parse operation binding the contract event 0xadff14cd34035a6bbb90fbe80979f36398f244f1885f7612e6e33a05a0b90d0f.
//
// Solidity: event Undelegate(address indexed sender, string validator, uint256 shares, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) ParseUndelegate(log types.Log) (*StakingTestUndelegate, error) {
	event := new(StakingTestUndelegate)
	if err := _StakingTest.contract.UnpackLog(event, "Undelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestUndelegateV2Iterator is returned from FilterUndelegateV2 and is used to iterate over the raw logs and unpacked data for UndelegateV2 events raised by the StakingTest contract.
type StakingTestUndelegateV2Iterator struct {
	Event *StakingTestUndelegateV2 // Event containing the contract specifics and raw log

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
func (it *StakingTestUndelegateV2Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestUndelegateV2)
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
		it.Event = new(StakingTestUndelegateV2)
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
func (it *StakingTestUndelegateV2Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestUndelegateV2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestUndelegateV2 represents a UndelegateV2 event raised by the StakingTest contract.
type StakingTestUndelegateV2 struct {
	Sender         common.Address
	Validator      string
	Amount         *big.Int
	CompletionTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUndelegateV2 is a free log retrieval operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) FilterUndelegateV2(opts *bind.FilterOpts, sender []common.Address) (*StakingTestUndelegateV2Iterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "UndelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestUndelegateV2Iterator{contract: _StakingTest.contract, event: "UndelegateV2", logs: logs, sub: sub}, nil
}

// WatchUndelegateV2 is a free log subscription operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) WatchUndelegateV2(opts *bind.WatchOpts, sink chan<- *StakingTestUndelegateV2, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "UndelegateV2", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestUndelegateV2)
				if err := _StakingTest.contract.UnpackLog(event, "UndelegateV2", log); err != nil {
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

// ParseUndelegateV2 is a log parse operation binding the contract event 0x4d3e71c3e3ff90f64b7095a17eb6b6cdd1ca0f0563102ef30415f73cb64b866f.
//
// Solidity: event UndelegateV2(address indexed sender, string validator, uint256 amount, uint256 completionTime)
func (_StakingTest *StakingTestFilterer) ParseUndelegateV2(log types.Log) (*StakingTestUndelegateV2, error) {
	event := new(StakingTestUndelegateV2)
	if err := _StakingTest.contract.UnpackLog(event, "UndelegateV2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingTestWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the StakingTest contract.
type StakingTestWithdrawIterator struct {
	Event *StakingTestWithdraw // Event containing the contract specifics and raw log

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
func (it *StakingTestWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestWithdraw)
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
		it.Event = new(StakingTestWithdraw)
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
func (it *StakingTestWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestWithdraw represents a Withdraw event raised by the StakingTest contract.
type StakingTestWithdraw struct {
	Sender    common.Address
	Validator string
	Reward    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_StakingTest *StakingTestFilterer) FilterWithdraw(opts *bind.FilterOpts, sender []common.Address) (*StakingTestWithdrawIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestWithdrawIterator{contract: _StakingTest.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_StakingTest *StakingTestFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StakingTestWithdraw, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Withdraw", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestWithdraw)
				if err := _StakingTest.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x901c03da5d88eb3d62ab4617e7b7d17d86db16356823a7971127d5181a842fef.
//
// Solidity: event Withdraw(address indexed sender, string validator, uint256 reward)
func (_StakingTest *StakingTestFilterer) ParseWithdraw(log types.Log) (*StakingTestWithdraw, error) {
	event := new(StakingTestWithdraw)
	if err := _StakingTest.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
