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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"ApproveShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"token\",\"type\":\"uint256\"}],\"name\":\"TransferShares\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"completionTime\",\"type\":\"uint256\"}],\"name\":\"Undelegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"validator\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowanceShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"approveShares\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"delegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_del\",\"type\":\"address\"}],\"name\":\"delegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferFromShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"transferShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"undelegate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"validatorShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_val\",\"type\":\"string\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611e60806100206000396000f3fe6080604052600436106100915760003560e01c80638dfc8897116100595780638dfc8897146101c85780639ddb511a14610207578063bf98d77214610238578063d5c498eb14610275578063dc6ffc7d146102b357610091565b8063161298c11461009657806331fb67c2146100d457806349da433e1461011157806351af513a1461014e5780637b625c0f1461018b575b600080fd5b3480156100a257600080fd5b506100bd60048036038101906100b89190611586565b6102f1565b6040516100cb929190611604565b60405180910390f35b3480156100e057600080fd5b506100fb60048036038101906100f6919061162d565b610316565b6040516101089190611676565b60405180910390f35b34801561011d57600080fd5b5061013860048036038101906101339190611586565b61032d565b60405161014591906116ac565b60405180910390f35b34801561015a57600080fd5b50610175600480360381019061017091906116c7565b610348565b6040516101829190611676565b60405180910390f35b34801561019757600080fd5b506101b260048036038101906101ad9190611723565b61035c565b6040516101bf9190611676565b60405180910390f35b3480156101d457600080fd5b506101ef60048036038101906101ea9190611792565b610372565b6040516101fe939291906117ee565b60405180910390f35b610221600480360381019061021c919061162d565b6103d3565b60405161022f929190611604565b60405180910390f35b34801561024457600080fd5b5061025f600480360381019061025a919061162d565b610428565b60405161026c9190611676565b60405180910390f35b34801561028157600080fd5b5061029c600480360381019061029791906116c7565b610456565b6040516102aa929190611604565b60405180910390f35b3480156102bf57600080fd5b506102da60048036038101906102d59190611825565b61046e565b6040516102e8929190611604565b60405180910390f35b600080600080610302878787610495565b915091508181935093505050935093915050565b6000806103228361056a565b905080915050919050565b60008061033b858585610637565b9050809150509392505050565b60006103548383610708565b905092915050565b60006103698484846107d5565b90509392505050565b60008060008060008061038588886108a4565b9250925092508660008960405161039c9190611919565b908152602001604051809103902060008282546103b9919061195f565b925050819055508282829550955095505050509250925092565b6000806000806103e28561097b565b91509150816000866040516103f79190611919565b908152602001604051809103902060008282546104149190611993565b925050819055508181935093505050915091565b6000818051602081018201805184825260208301602085012081835280955050505050506000915090505481565b6000806104638484610a4c565b915091509250929050565b60008060008061048088888888610b1d565b91509150818193509350505094509492505050565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff166104bf888888610bf4565b6040516104cc9190611a0e565b6000604051808303816000865af19150503d8060008114610509576040519150601f19603f3d011682016040523d82523d6000602084013e61050e565b606091505b509150915061055382826040518060400160405280601681526020017f7472616e7366657220736861726573206661696c656400000000000000000000815250610c91565b61055c81610d58565b935093505050935093915050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff1661059185610d83565b60405161059e9190611a0e565b6000604051808303816000865af19150503d80600081146105db576040519150601f19603f3d011682016040523d82523d6000602084013e6105e0565b606091505b509150915061062582826040518060400160405280600f81526020017f7769746864726177206661696c65640000000000000000000000000000000000815250610c91565b61062e81610e1a565b92505050919050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff16610660878787610e3c565b60405161066d9190611a0e565b6000604051808303816000865af19150503d80600081146106aa576040519150601f19603f3d011682016040523d82523d6000602084013e6106af565b606091505b50915091506106f482826040518060400160405280601581526020017f617070726f766520736861726573206661696c65640000000000000000000000815250610c91565b6106fd81610ed9565b925050509392505050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff166107308686610efb565b60405161073d9190611a0e565b600060405180830381855afa9150503d8060008114610778576040519150601f19603f3d011682016040523d82523d6000602084013e61077d565b606091505b50915091506107c282826040518060400160405280601881526020017f64656c65676174696f6e52657761726473206661696c65640000000000000000815250610c91565b6107cb81610f95565b9250505092915050565b600080600061100373ffffffffffffffffffffffffffffffffffffffff166107fe878787610fb7565b60405161080b9190611a0e565b600060405180830381855afa9150503d8060008114610846576040519150601f19603f3d011682016040523d82523d6000602084013e61084b565b606091505b509150915061089082826040518060400160405280601781526020017f616c6c6f77616e636520736861726573206661696c6564000000000000000000815250610c91565b61089981611054565b925050509392505050565b600080600080600061100373ffffffffffffffffffffffffffffffffffffffff166108cf8888611076565b6040516108dc9190611a0e565b6000604051808303816000865af19150503d8060008114610919576040519150601f19603f3d011682016040523d82523d6000602084013e61091e565b606091505b509150915061096382826040518060400160405280601181526020017f756e64656c6567617465206661696c6564000000000000000000000000000000815250610c91565b61096c81611110565b94509450945050509250925092565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff16346109a487611146565b6040516109b19190611a0e565b60006040518083038185875af1925050503d80600081146109ee576040519150601f19603f3d011682016040523d82523d6000602084013e6109f3565b606091505b5091509150610a3882826040518060400160405280600f81526020017f64656c6567617465206661696c65640000000000000000000000000000000000815250610c91565b610a41816111dd565b935093505050915091565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff16610a758787611208565b604051610a829190611a0e565b600060405180830381855afa9150503d8060008114610abd576040519150601f19603f3d011682016040523d82523d6000602084013e610ac2565b606091505b5091509150610b0782826040518060400160405280601181526020017f64656c65676174696f6e206661696c6564000000000000000000000000000000815250610c91565b610b10816112a2565b9350935050509250929050565b60008060008061100373ffffffffffffffffffffffffffffffffffffffff16610b48898989896112cd565b604051610b559190611a0e565b6000604051808303816000865af19150503d8060008114610b92576040519150601f19603f3d011682016040523d82523d6000602084013e610b97565b606091505b5091509150610bdc82826040518060400160405280601a81526020017f7472616e7366657246726f6d20736861726573206661696c6564000000000000815250610c91565b610be58161136d565b93509350505094509492505050565b6060838383604051602401610c0b93929190611a7e565b6040516020818303038152906040527f161298c1000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509392505050565b82610d5357600082806020019051810190610cac9190611b2c565b9050600182511015610cf557806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cec9190611b75565b60405180910390fd5b8181604051602001610d08929190611be3565b6040516020818303038152906040526040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d4a9190611b75565b60405180910390fd5b505050565b60008060008084806020019051810190610d729190611c27565b915091508181935093505050915091565b606081604051602401610d969190611b75565b6040516020818303038152906040527f31fb67c2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b60008082806020019051810190610e319190611c67565b905080915050919050565b6060838383604051602401610e5393929190611a7e565b6040516020818303038152906040527f49da433e000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509392505050565b60008082806020019051810190610ef09190611cc0565b905080915050919050565b60608282604051602401610f10929190611ced565b6040516020818303038152906040527f51af513a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b60008082806020019051810190610fac9190611c67565b905080915050919050565b6060838383604051602401610fce93929190611d1d565b6040516020818303038152906040527f7b625c0f000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509392505050565b6000808280602001905181019061106b9190611c67565b905080915050919050565b6060828260405160240161108b929190611d5b565b6040516020818303038152906040527f8dfc8897000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b6000806000806000808680602001905181019061112d9190611d8b565b9250925092508282829550955095505050509193909250565b6060816040516024016111599190611b75565b6040516020818303038152906040527f9ddb511a000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050919050565b600080600080848060200190518101906111f79190611c27565b915091508181935093505050915091565b6060828260405160240161121d929190611ced565b6040516020818303038152906040527fd5c498eb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050905092915050565b600080600080848060200190518101906112bc9190611c27565b915091508181935093505050915091565b6060848484846040516024016112e69493929190611dde565b6040516020818303038152906040527fdc6ffc7d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050509050949350505050565b600080600080848060200190518101906113879190611c27565b915091508181935093505050915091565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6113ff826113b6565b810181811067ffffffffffffffff8211171561141e5761141d6113c7565b5b80604052505050565b6000611431611398565b905061143d82826113f6565b919050565b600067ffffffffffffffff82111561145d5761145c6113c7565b5b611466826113b6565b9050602081019050919050565b82818337600083830152505050565b600061149561149084611442565b611427565b9050828152602081018484840111156114b1576114b06113b1565b5b6114bc848285611473565b509392505050565b600082601f8301126114d9576114d86113ac565b5b81356114e9848260208601611482565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061151d826114f2565b9050919050565b61152d81611512565b811461153857600080fd5b50565b60008135905061154a81611524565b92915050565b6000819050919050565b61156381611550565b811461156e57600080fd5b50565b6000813590506115808161155a565b92915050565b60008060006060848603121561159f5761159e6113a2565b5b600084013567ffffffffffffffff8111156115bd576115bc6113a7565b5b6115c9868287016114c4565b93505060206115da8682870161153b565b92505060406115eb86828701611571565b9150509250925092565b6115fe81611550565b82525050565b600060408201905061161960008301856115f5565b61162660208301846115f5565b9392505050565b600060208284031215611643576116426113a2565b5b600082013567ffffffffffffffff811115611661576116606113a7565b5b61166d848285016114c4565b91505092915050565b600060208201905061168b60008301846115f5565b92915050565b60008115159050919050565b6116a681611691565b82525050565b60006020820190506116c1600083018461169d565b92915050565b600080604083850312156116de576116dd6113a2565b5b600083013567ffffffffffffffff8111156116fc576116fb6113a7565b5b611708858286016114c4565b92505060206117198582860161153b565b9150509250929050565b60008060006060848603121561173c5761173b6113a2565b5b600084013567ffffffffffffffff81111561175a576117596113a7565b5b611766868287016114c4565b93505060206117778682870161153b565b92505060406117888682870161153b565b9150509250925092565b600080604083850312156117a9576117a86113a2565b5b600083013567ffffffffffffffff8111156117c7576117c66113a7565b5b6117d3858286016114c4565b92505060206117e485828601611571565b9150509250929050565b600060608201905061180360008301866115f5565b61181060208301856115f5565b61181d60408301846115f5565b949350505050565b6000806000806080858703121561183f5761183e6113a2565b5b600085013567ffffffffffffffff81111561185d5761185c6113a7565b5b611869878288016114c4565b945050602061187a8782880161153b565b935050604061188b8782880161153b565b925050606061189c87828801611571565b91505092959194509250565b600081519050919050565b600081905092915050565b60005b838110156118dc5780820151818401526020810190506118c1565b60008484015250505050565b60006118f3826118a8565b6118fd81856118b3565b935061190d8185602086016118be565b80840191505092915050565b600061192582846118e8565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061196a82611550565b915061197583611550565b925082820390508181111561198d5761198c611930565b5b92915050565b600061199e82611550565b91506119a983611550565b92508282019050808211156119c1576119c0611930565b5b92915050565b600081519050919050565b600081905092915050565b60006119e8826119c7565b6119f281856119d2565b9350611a028185602086016118be565b80840191505092915050565b6000611a1a82846119dd565b915081905092915050565b600082825260208201905092915050565b6000611a41826118a8565b611a4b8185611a25565b9350611a5b8185602086016118be565b611a64816113b6565b840191505092915050565b611a7881611512565b82525050565b60006060820190508181036000830152611a988186611a36565b9050611aa76020830185611a6f565b611ab460408301846115f5565b949350505050565b6000611acf611aca84611442565b611427565b905082815260208101848484011115611aeb57611aea6113b1565b5b611af68482856118be565b509392505050565b600082601f830112611b1357611b126113ac565b5b8151611b23848260208601611abc565b91505092915050565b600060208284031215611b4257611b416113a2565b5b600082015167ffffffffffffffff811115611b6057611b5f6113a7565b5b611b6c84828501611afe565b91505092915050565b60006020820190508181036000830152611b8f8184611a36565b905092915050565b7f3a20000000000000000000000000000000000000000000000000000000000000600082015250565b6000611bcd6002836118b3565b9150611bd882611b97565b600282019050919050565b6000611bef82856118e8565b9150611bfa82611bc0565b9150611c0682846118e8565b91508190509392505050565b600081519050611c218161155a565b92915050565b60008060408385031215611c3e57611c3d6113a2565b5b6000611c4c85828601611c12565b9250506020611c5d85828601611c12565b9150509250929050565b600060208284031215611c7d57611c7c6113a2565b5b6000611c8b84828501611c12565b91505092915050565b611c9d81611691565b8114611ca857600080fd5b50565b600081519050611cba81611c94565b92915050565b600060208284031215611cd657611cd56113a2565b5b6000611ce484828501611cab565b91505092915050565b60006040820190508181036000830152611d078185611a36565b9050611d166020830184611a6f565b9392505050565b60006060820190508181036000830152611d378186611a36565b9050611d466020830185611a6f565b611d536040830184611a6f565b949350505050565b60006040820190508181036000830152611d758185611a36565b9050611d8460208301846115f5565b9392505050565b600080600060608486031215611da457611da36113a2565b5b6000611db286828701611c12565b9350506020611dc386828701611c12565b9250506040611dd486828701611c12565b9150509250925092565b60006080820190508181036000830152611df88187611a36565b9050611e076020830186611a6f565b611e146040830185611a6f565b611e2160608301846115f5565b9594505050505056fea26469706673582212202ea2ba76fb940f3e36d6674ae6ed987ce66d86b4b8175c6cc5d554416363edc364736f6c63430008130033",
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
