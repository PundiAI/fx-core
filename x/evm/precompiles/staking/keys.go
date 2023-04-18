package staking

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v4/contract"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

const (
	DelegateGas           = 40000 // 98000 - 160000 // 165000
	UndelegateGas         = 45000 // 94000 - 163000 // 172000
	WithdrawGas           = 30000 // 94000 // 120000
	DelegationGas         = 30000 // 98000
	DelegationRewardsGas  = 30000 // 94000
	TransferSharesGas     = 50000 // 134000 - 190000
	ApproveSharesGas      = 10000 // 4400
	AllowanceSharesGas    = 5000  // 1240
	TransferFromSharesGas = 60000 // 134000 - 200000

	DelegateMethodName           = "delegate"
	UndelegateMethodName         = "undelegate"
	WithdrawMethodName           = "withdraw"
	DelegationMethodName         = "delegation"
	DelegationRewardsMethodName  = "delegationRewards"
	ApproveSharesMethodName      = "approveShares"
	AllowanceSharesMethodName    = "allowanceShares"
	TransferSharesMethodName     = "transferShares"
	TransferFromSharesMethodName = "transferFromShares"

	DelegateEventName       = "Delegate"
	UndelegateEventName     = "Undelegate"
	WithdrawEventName       = "Withdraw"
	TransferSharesEventName = "TransferShares"
	ApproveSharesEventName  = "ApproveShares"
)

var (
	stakingAddress = common.HexToAddress(fxtypes.StakingAddress)
	stakingABI     = fxtypes.MustABIJson(contract.IStakingMetaData.ABI)
)

func GetAddress() common.Address {
	return stakingAddress
}

func GetABI() abi.ABI {
	return stakingABI
}
