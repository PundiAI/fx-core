package staking

import (
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

const (
	// DelegateGas default delegate use 0.76FX
	DelegateGas = 400000 // if gas price 500Gwei, fee is 0.2FX
	// UndelegateGas default undelegate use 0.82FX
	UndelegateGas = 600000 // if gas price 500Gwei, fee is 0.3FX
	// WithdrawGas default withdraw use 0.56FX
	WithdrawGas          = 300000 // if gas price 500Gwei, fee is 0.15FX
	DelegationGas        = 200000 // if gas price 500Gwei, fee is 0.1FX
	DelegationRewardsGas = 200000
	TransferGas          = 400000
	ApproveGas           = 200000
	AllowanceGas         = 100000

	DelegateMethodName          = "delegate"
	UndelegateMethodName        = "undelegate"
	WithdrawMethodName          = "withdraw"
	DelegationMethodName        = "delegation"
	DelegationRewardsMethodName = "delegationRewards"
	ApproveMethodName           = "approve"
	AllowanceMethodName         = "allowance"
	TransferMethodName          = "transfer"
	TransferFromMethodName      = "transferFrom"

	DelegateEventName   = "Delegate"
	UndelegateEventName = "Undelegate"
	WithdrawEventName   = "Withdraw"
	TransferEventName   = "Transfer"
	ApproveEventName    = "Approve"

	JsonABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"delegator","type":"address"},{"indexed":false,"internalType":"string","name":"validator","type":"string"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"shares","type":"uint256"}],"name":"Delegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"validator","type":"string"},{"indexed":false,"internalType":"uint256","name":"shares","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"completionTime","type":"uint256"}],"name":"Undelegate","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"string","name":"validator","type":"string"},{"indexed":false,"internalType":"uint256","name":"reward","type":"uint256"}],"name":"Withdraw","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"string","name":"validator","type":"string"},{"indexed":false,"internalType":"uint256","name":"shares","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"token","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"string","name":"validator","type":"string"},{"indexed":false,"internalType":"uint256","name":"shares","type":"uint256"}],"name":"Approve","type":"event"},{"type":"function","name":"delegate","inputs":[{"name":"validator","type":"string"}],"outputs":[{"name":"shares","type":"uint256"},{"name":"reward","type":"uint256"}],"payable":true,"stateMutability":"payable"},{"type":"function","name":"undelegate","inputs":[{"name":"validator","type":"string"},{"name":"shares","type":"uint256"}],"outputs":[{"name":"amount","type":"uint256"},{"name":"reward","type":"uint256"},{"name":"completionTime","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"withdraw","inputs":[{"name":"validator","type":"string"}],"outputs":[{"name":"reward","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"transfer","inputs":[{"name":"validator","type":"string"},{"name":"to","type":"address"},{"name":"shares","type":"uint256"}],"outputs":[{"name":"token","type":"uint256"},{"name":"reward","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"approve","inputs":[{"name":"validator","type":"string"},{"name":"spender","type":"address"},{"name":"shares","type":"uint256"}],"outputs":[{"name":"result","type":"bool"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"transferFrom","inputs":[{"name":"validator","type":"string"},{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"shares","type":"uint256"}],"outputs":[{"name":"token","type":"uint256"},{"name":"reward","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"delegation","inputs":[{"name":"validator","type":"string"},{"name":"delegator","type":"address"}],"outputs":[{"name":"shares","type":"uint256"},{"name":"delegate","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"delegationRewards","inputs":[{"name":"validator","type":"string"},{"name":"delegator","type":"address"}],"outputs":[{"name":"rewards","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"},{"type":"function","name":"allowance","inputs":[{"name":"validator","type":"string"},{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"outputs":[{"name":"shares","type":"uint256"}],"payable":false,"stateMutability":"nonpayable"}]`
)

var precompileAddress = common.HexToAddress(fxtypes.StakingAddress)

func GetPrecompileAddress() common.Address {
	return precompileAddress
}
