package staking

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v6/contract"
	fxtypes "github.com/functionx/fx-core/v6/types"
)

const (
	DelegateGas           = 40_000 // 98000 - 160000 // 165000
	UndelegateGas         = 45_000 // 94000 - 163000 // 172000
	RedelegateGas         = 60_000 // undelegate_gas+delegate_gas+withdraw_gas*2
	WithdrawGas           = 30_000 // 94000 // 120000
	DelegationGas         = 30_000 // 98000
	DelegationRewardsGas  = 30_000 // 94000
	TransferSharesGas     = 50_000 // 134000 - 190000
	ApproveSharesGas      = 10_000 // 4400
	AllowanceSharesGas    = 5_000  // 1240
	TransferFromSharesGas = 60_000 // 134000 - 200000

	DelegateMethodName           = "delegate"
	UndelegateMethodName         = "undelegate"
	RedelegateMethodName         = "redelegate"
	WithdrawMethodName           = "withdraw"
	DelegationMethodName         = "delegation"
	DelegationRewardsMethodName  = "delegationRewards"
	ApproveSharesMethodName      = "approveShares"
	AllowanceSharesMethodName    = "allowanceShares"
	TransferSharesMethodName     = "transferShares"
	TransferFromSharesMethodName = "transferFromShares"

	DelegateEventName       = "Delegate"
	UndelegateEventName     = "Undelegate"
	RedelegateEventName     = "Redelegate"
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
