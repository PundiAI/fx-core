package staking

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

const (
	// DelegateGas default delegate use 0.76FX
	DelegateGas = 400000 // if gas price 500Gwei, fee is 0.2FX
	// UndelegateGas default undelegate use 0.82FX
	UndelegateGas = 600000 // if gas price 500Gwei, fee is 0.3FX
	// WithdrawGas default withdraw use 0.56FX
	WithdrawGas           = 300000 // if gas price 500Gwei, fee is 0.15FX
	DelegationGas         = 200000 // if gas price 500Gwei, fee is 0.1FX
	DelegationRewardsGas  = 200000
	TransferSharesGas     = 400000
	ApproveSharesGas      = 200000
	AllowanceSharesGas    = 100000
	TransferFromSharesGas = 500000

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
	stakingABI     = fxtypes.MustABIJson(StakingMetaData.ABI)
)

func GetAddress() common.Address {
	return stakingAddress
}

func GetABI() abi.ABI {
	return stakingABI
}
