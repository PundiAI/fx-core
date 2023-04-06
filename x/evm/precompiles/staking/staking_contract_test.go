// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package staking_test

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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Approve\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Undelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"validatorShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611e60806100206000396000f3fe6080604052600436106100915760003560e01c80639ddb511a116100595780639ddb511a146101ca578063aa2fddc7146101fb578063bd5b93d214610239578063bf98d77214610276578063d5c498eb146102b357610091565b80631fc38a8b1461009657806331fb67c2146100d4578063495aff711461011157806351af513a1461014e5780638dfc88971461018b575b600080fd5b3480156100a257600080fd5b506100bd60048036038101906100b89190611586565b6102f1565b6040516100cb929190611618565b60405180910390f35b3480156100e057600080fd5b506100fb60048036038101906100f69190611641565b610318565b604051610108919061168a565b60405180910390f35b34801561011d57600080fd5b50610138600480360381019061013391906116a5565b61032f565b604051610145919061168a565b60405180910390f35b34801561015a57600080fd5b5061017560048036038101906101709190611714565b610345565b604051610182919061168a565b60405180910390f35b34801561019757600080fd5b506101b260048036038101906101ad9190611770565b610359565b6040516101c1939291906117cc565b60405180910390f35b6101e460048036038101906101df9190611641565b6103ba565b6040516101f2929190611618565b60405180910390f35b34801561020757600080fd5b50610222600480360381019061021d9190611803565b61040f565b604051610230929190611618565b60405180910390f35b34801561024557600080fd5b50610260600480360381019061025b9190611803565b610434565b60405161026d919061188d565b60405180910390f35b34801561028257600080fd5b5061029d60048036038101906102989190611641565b61044f565b6040516102aa919061168a565b60405180910390f35b3480156102bf57600080fd5b506102da60048036038101906102d59190611714565b61047d565b6040516102e8929190611618565b60405180910390f35b60008060008061030388888888610495565b91509150818193509350505094509492505050565b6000806103248361056c565b905080915050919050565b600061033c848484610639565b90509392505050565b60006103518383610708565b905092915050565b60008060008060008061036c88886107d5565b925092509250866000896040516103839190611919565b908152602001604051809103902060008282546103a0919061195f565b925050819055508282829550955095505050509250925092565b6000806000806103c9856108ac565b91509150816000866040516103de9190611919565b908152602001604051809103902060008282546103fb9190611993565b925050819055508181935093505050915091565b60008060008061042087878761097d565b915091508181935093505050935093915050565b600080610442858585610a52565b9050809150509392505050565b6000818051602081018201805184825260208301602085012081835280955050505050506000915090505481565b60008061048a8484610b23565b915091509250929050565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff166104c089898989610bf4565b6040516104cd9190611a0e565b6000604051808303816000865af19150503d806000811461050a576040519150601f19603f3d011682016040523d82523d6000602084013e61050f565b606091505b509150915061055482826040518060400160405280601381526020017f7472616e7366657246726f6d206661696c656400000000000000000000000000815250610c94565b61055d81610d5b565b93509350505094509492505050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff1661059385610d86565b6040516105a09190611a0e565b6000604051808303816000865af19150503d80600081146105dd576040519150601f19603f3d011682016040523d82523d6000602084013e6105e2565b606091505b509150915061062782826040518060400160405280600f81526020017f7769746864726177206661696c65640000000000000000000000000000000000815250610c94565b61063081610e1d565b92505050919050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff16610662878787610e3f565b60405161066f9190611a0e565b600060405180830381855afa9150503d80600081146106aa576040519150601f19603f3d011682016040523d82523d6000602084013e6106af565b606091505b50915091506106f482826040518060400160405280601081526020017f616c6c6f77616e6365206661696c656400000000000000000000000000000000815250610c94565b6106fd81610edc565b925050509392505050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff166107308686610efe565b60405161073d9190611a0e565b600060405180830381855afa9150503d8060008114610778576040519150601f19603f3d011682016040523d82523d6000602084013e61077d565b606091505b50915091506107c282826040518060400160405280601881526020017f64656c65676174696f6e52657761726473206661696c65640000000000000000815250610c94565b6107cb81610f98565b9250505092915050565b600080600080600061100373ffffffffffffffffffffffffffffffffffffffff166108008888610fba565b60405161080d9190611a0e565b6000604051808303816000865af19150503d806000811461084a576040519150601f19603f3d011682016040523d82523d6000602084013e61084f565b606091505b509150915061089482826040518060400160405280601181526020017f756e64656c6567617465206661696c6564000000000000000000000000000000815250610c94565b61089d81611054565b94509450945050509250925092565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff16346108d58761108a565b6040516108e29190611a0e565b60006040518083038185875af1925050503d806000811461091f576040519150601f19603f3d011682016040523d82523d6000602084013e610924565b606091505b509150915061096982826040518060400160405280600f81526020017f64656c6567617465206661696c65640000000000000000000000000000000000815250610c94565b61097281611121565b935093505050915091565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff166109a788888861114c565b6040516109b49190611a0e565b6000604051808303816000865af19150503d80600081146109f1576040519150601f19603f3d011682016040523d82523d6000602084013e6109f6565b606091505b5091509150610a3b82826040518060400160405280600f81526020017f7472616e73666572206661696c65640000000000000000000000000000000000815250610c94565b610a44816111e9565b935093505050935093915050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff16610a7b878787611214565b604051610a889190611a0e565b6000604051808303816000865af19150503d8060008114610ac5576040519150601f19603f3d011682016040523d82523d6000602084013e610aca565b606091505b5091509150610b0f82826040518060400160405280600e81526020017f617070726f7665206661696c6564000000000000000000000000000000000000815250610c94565b610b18816112b1565b925050509392505050565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff16610b4c87876112d3565b604051610b599190611a0e565b600060405180830381855afa9150503d8060008114610b94576040519150601f19603f3d011682016040523d82523d6000602084013e610b99565b606091505b5091509150610bde82826040518060400160405280601181526020017f64656c65676174696f6e206661696c6564000000000000000000000000000000815250610c94565b610be78161136d565b9350935050509250929050565b606084848484604051602401610c0d9493929190611a7e565b6040516020818303038152906040527f1fc38a8b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050949350505050565b82610d5657600082806020019051810190610caf9190611b3a565b9050600182511015610cf857806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cef9190611b83565b60405180910390fd5b8181604051602001610d0b929190611bf1565b6040516020818303038152906040526040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d4d9190611b83565b60405180910390fd5b505050565b60008060008084806020019051810190610d759190611c35565b915091508181935093505050915091565b606081604051602401610d999190611b83565b6040516020818303038152906040527f31fb67c2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b60008082806020019051810190610e349190611c75565b905080915050919050565b6060838383604051602401610e5693929190611ca2565b6040516020818303038152906040527f495aff71000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509392505050565b60008082806020019051810190610ef39190611c75565b905080915050919050565b60608282604051602401610f13929190611ce0565b6040516020818303038152906040527f51af513a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b60008082806020019051810190610faf9190611c75565b905080915050919050565b60608282604051602401610fcf929190611d10565b6040516020818303038152906040527f8dfc8897000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b600080600080600080868060200190518101906110719190611d40565b9250925092508282829550955095505050509193909250565b60608160405160240161109d9190611b83565b6040516020818303038152906040527f9ddb511a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b6000806000808480602001905181019061113b9190611c35565b915091508181935093505050915091565b606083838360405160240161116393929190611d93565b6040516020818303038152906040527faa2fddc7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509392505050565b600080600080848060200190518101906112039190611c35565b915091508181935093505050915091565b606083838360405160240161122b93929190611d93565b6040516020818303038152906040527fbd5b93d2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509392505050565b600080828060200190518101906112c89190611dfd565b905080915050919050565b606082826040516024016112e8929190611ce0565b6040516020818303038152906040527fd5c498eb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b600080600080848060200190518101906113879190611c35565b915091508181935093505050915091565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6113ff826113b6565b810181811067ffffffffffffffff8211171561141e5761141d6113c7565b5b80604052505050565b6000611431611398565b905061143d82826113f6565b919050565b600067ffffffffffffffff82111561145d5761145c6113c7565b5b611466826113b6565b9050602081019050919050565b82818337600083830152505050565b600061149561149084611442565b611427565b9050828152602081018484840111156114b1576114b06113b1565b5b6114bc848285611473565b509392505050565b600082601f8301126114d9576114d86113ac565b5b81356114e9848260208601611482565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061151d826114f2565b9050919050565b61152d81611512565b811461153857600080fd5b50565b60008135905061154a81611524565b92915050565b6000819050919050565b61156381611550565b811461156e57600080fd5b50565b6000813590506115808161155a565b92915050565b600080600080608085870312156115a05761159f6113a2565b5b600085013567ffffffffffffffff8111156115be576115bd6113a7565b5b6115ca878288016114c4565b94505060206115db8782880161153b565b93505060406115ec8782880161153b565b92505060606115fd87828801611571565b91505092959194509250565b61161281611550565b82525050565b600060408201905061162d6000830185611609565b61163a6020830184611609565b9392505050565b600060208284031215611657576116566113a2565b5b600082013567ffffffffffffffff811115611675576116746113a7565b5b611681848285016114c4565b91505092915050565b600060208201905061169f6000830184611609565b92915050565b6000806000606084860312156116be576116bd6113a2565b5b600084013567ffffffffffffffff8111156116dc576116db6113a7565b5b6116e8868287016114c4565b93505060206116f98682870161153b565b925050604061170a8682870161153b565b9150509250925092565b6000806040838503121561172b5761172a6113a2565b5b600083013567ffffffffffffffff811115611749576117486113a7565b5b611755858286016114c4565b92505060206117668582860161153b565b9150509250929050565b60008060408385031215611787576117866113a2565b5b600083013567ffffffffffffffff8111156117a5576117a46113a7565b5b6117b1858286016114c4565b92505060206117c285828601611571565b9150509250929050565b60006060820190506117e16000830186611609565b6117ee6020830185611609565b6117fb6040830184611609565b949350505050565b60008060006060848603121561181c5761181b6113a2565b5b600084013567ffffffffffffffff81111561183a576118396113a7565b5b611846868287016114c4565b93505060206118578682870161153b565b925050604061186886828701611571565b9150509250925092565b60008115159050919050565b61188781611872565b82525050565b60006020820190506118a2600083018461187e565b92915050565b600081519050919050565b600081905092915050565b60005b838110156118dc5780820151818401526020810190506118c1565b60008484015250505050565b60006118f3826118a8565b6118fd81856118b3565b935061190d8185602086016118be565b80840191505092915050565b600061192582846118e8565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061196a82611550565b915061197583611550565b925082820390508181111561198d5761198c611930565b5b92915050565b600061199e82611550565b91506119a983611550565b92508282019050808211156119c1576119c0611930565b5b92915050565b600081519050919050565b600081905092915050565b60006119e8826119c7565b6119f281856119d2565b9350611a028185602086016118be565b80840191505092915050565b6000611a1a82846119dd565b915081905092915050565b600082825260208201905092915050565b6000611a41826118a8565b611a4b8185611a25565b9350611a5b8185602086016118be565b611a64816113b6565b840191505092915050565b611a7881611512565b82525050565b60006080820190508181036000830152611a988187611a36565b9050611aa76020830186611a6f565b611ab46040830185611a6f565b611ac16060830184611609565b95945050505050565b6000611add611ad884611442565b611427565b905082815260208101848484011115611af957611af86113b1565b5b611b048482856118be565b509392505050565b600082601f830112611b2157611b206113ac565b5b8151611b31848260208601611aca565b91505092915050565b600060208284031215611b5057611b4f6113a2565b5b600082015167ffffffffffffffff811115611b6e57611b6d6113a7565b5b611b7a84828501611b0c565b91505092915050565b60006020820190508181036000830152611b9d8184611a36565b905092915050565b7f3a20000000000000000000000000000000000000000000000000000000000000600082015250565b6000611bdb6002836118b3565b9150611be682611ba5565b600282019050919050565b6000611bfd82856118e8565b9150611c0882611bce565b9150611c1482846118e8565b91508190509392505050565b600081519050611c2f8161155a565b92915050565b60008060408385031215611c4c57611c4b6113a2565b5b6000611c5a85828601611c20565b9250506020611c6b85828601611c20565b9150509250929050565b600060208284031215611c8b57611c8a6113a2565b5b6000611c9984828501611c20565b91505092915050565b60006060820190508181036000830152611cbc8186611a36565b9050611ccb6020830185611a6f565b611cd86040830184611a6f565b949350505050565b60006040820190508181036000830152611cfa8185611a36565b9050611d096020830184611a6f565b9392505050565b60006040820190508181036000830152611d2a8185611a36565b9050611d396020830184611609565b9392505050565b600080600060608486031215611d5957611d586113a2565b5b6000611d6786828701611c20565b9350506020611d7886828701611c20565b9250506040611d8986828701611c20565b9150509250925092565b60006060820190508181036000830152611dad8186611a36565b9050611dbc6020830185611a6f565b611dc96040830184611609565b949350505050565b611dda81611872565b8114611de557600080fd5b50565b600081519050611df781611dd1565b92915050565b600060208284031215611e1357611e126113a2565b5b6000611e2184828501611de8565b9150509291505056fea2646970667358221220229c65e1e45b127ddbc0022b111c11bfdfbda084ac8220e9dc07f93e3ae7273764736f6c63430008130033",
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

// Allowance is a free data retrieval call binding the contract method 0x495aff71.
//
// Solidity: function allowance(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestCaller) Allowance(opts *bind.CallOpts, _val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingTest.contract.Call(opts, &out, "allowance", _val, _owner, _spender)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0x495aff71.
//
// Solidity: function allowance(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestSession) Allowance(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StakingTest.Contract.Allowance(&_StakingTest.CallOpts, _val, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0x495aff71.
//
// Solidity: function allowance(string _val, address _owner, address _spender) view returns(uint256)
func (_StakingTest *StakingTestCallerSession) Allowance(_val string, _owner common.Address, _spender common.Address) (*big.Int, error) {
	return _StakingTest.Contract.Allowance(&_StakingTest.CallOpts, _val, _owner, _spender)
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

// Approve is a paid mutator transaction binding the contract method 0xbd5b93d2.
//
// Solidity: function approve(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestTransactor) Approve(opts *bind.TransactOpts, _val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "approve", _val, _spender, _shares)
}

// Approve is a paid mutator transaction binding the contract method 0xbd5b93d2.
//
// Solidity: function approve(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestSession) Approve(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Approve(&_StakingTest.TransactOpts, _val, _spender, _shares)
}

// Approve is a paid mutator transaction binding the contract method 0xbd5b93d2.
//
// Solidity: function approve(string _val, address _spender, uint256 _shares) returns(bool)
func (_StakingTest *StakingTestTransactorSession) Approve(_val string, _spender common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Approve(&_StakingTest.TransactOpts, _val, _spender, _shares)
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

// Transfer is a paid mutator transaction binding the contract method 0xaa2fddc7.
//
// Solidity: function transfer(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) Transfer(opts *bind.TransactOpts, _val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "transfer", _val, _to, _shares)
}

// Transfer is a paid mutator transaction binding the contract method 0xaa2fddc7.
//
// Solidity: function transfer(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) Transfer(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Transfer(&_StakingTest.TransactOpts, _val, _to, _shares)
}

// Transfer is a paid mutator transaction binding the contract method 0xaa2fddc7.
//
// Solidity: function transfer(string _val, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) Transfer(_val string, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.Transfer(&_StakingTest.TransactOpts, _val, _to, _shares)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x1fc38a8b.
//
// Solidity: function transferFrom(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactor) TransferFrom(opts *bind.TransactOpts, _val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.contract.Transact(opts, "transferFrom", _val, _from, _to, _shares)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x1fc38a8b.
//
// Solidity: function transferFrom(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestSession) TransferFrom(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferFrom(&_StakingTest.TransactOpts, _val, _from, _to, _shares)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x1fc38a8b.
//
// Solidity: function transferFrom(string _val, address _from, address _to, uint256 _shares) returns(uint256, uint256)
func (_StakingTest *StakingTestTransactorSession) TransferFrom(_val string, _from common.Address, _to common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _StakingTest.Contract.TransferFrom(&_StakingTest.TransactOpts, _val, _from, _to, _shares)
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

// StakingTestApproveIterator is returned from FilterApprove and is used to iterate over the raw logs and unpacked data for Approve events raised by the StakingTest contract.
type StakingTestApproveIterator struct {
	Event *StakingTestApprove // Event containing the contract specifics and raw log

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
func (it *StakingTestApproveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestApprove)
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
		it.Event = new(StakingTestApprove)
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
func (it *StakingTestApproveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestApproveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestApprove represents a Approve event raised by the StakingTest contract.
type StakingTestApprove struct {
	Owner     common.Address
	Spender   common.Address
	Validator string
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterApprove is a free log retrieval operation binding the contract event 0x757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020.
//
// Solidity: event Approve(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) FilterApprove(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StakingTestApproveIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Approve", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestApproveIterator{contract: _StakingTest.contract, event: "Approve", logs: logs, sub: sub}, nil
}

// WatchApprove is a free log subscription operation binding the contract event 0x757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020.
//
// Solidity: event Approve(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) WatchApprove(opts *bind.WatchOpts, sink chan<- *StakingTestApprove, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Approve", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestApprove)
				if err := _StakingTest.contract.UnpackLog(event, "Approve", log); err != nil {
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

// ParseApprove is a log parse operation binding the contract event 0x757aaea72f5be0810fcdd3e3cf20ae8e5dc931f9290e900b0ce1a792eb4f3020.
//
// Solidity: event Approve(address indexed owner, address indexed spender, string validator, uint256 shares)
func (_StakingTest *StakingTestFilterer) ParseApprove(log types.Log) (*StakingTestApprove, error) {
	event := new(StakingTestApprove)
	if err := _StakingTest.contract.UnpackLog(event, "Approve", log); err != nil {
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

// StakingTestTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StakingTest contract.
type StakingTestTransferIterator struct {
	Event *StakingTestTransfer // Event containing the contract specifics and raw log

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
func (it *StakingTestTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTestTransfer)
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
		it.Event = new(StakingTestTransfer)
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
func (it *StakingTestTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTestTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTestTransfer represents a Transfer event raised by the StakingTest contract.
type StakingTestTransfer struct {
	From      common.Address
	To        common.Address
	Validator string
	Shares    *big.Int
	Token     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0x39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813.
//
// Solidity: event Transfer(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StakingTestTransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingTest.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StakingTestTransferIterator{contract: _StakingTest.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0x39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813.
//
// Solidity: event Transfer(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StakingTestTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StakingTest.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTestTransfer)
				if err := _StakingTest.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0x39ce25cb9eb8e7d68c58325d6493b7f9347812d615d19e33dc399364e792a813.
//
// Solidity: event Transfer(address indexed from, address indexed to, string validator, uint256 shares, uint256 token)
func (_StakingTest *StakingTestFilterer) ParseTransfer(log types.Log) (*StakingTestTransfer, error) {
	event := new(StakingTestTransfer)
	if err := _StakingTest.contract.UnpackLog(event, "Transfer", log); err != nil {
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
