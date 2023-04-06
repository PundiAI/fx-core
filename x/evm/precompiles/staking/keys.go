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
	TransferFromGas      = 500000

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
)

var (
	precompileAddress = common.HexToAddress(fxtypes.StakingAddress)
	JsonABI           = StakingMetaData.ABI
)

func GetPrecompileAddress() common.Address {
	return precompileAddress
}
