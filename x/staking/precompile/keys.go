package precompile

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

var (
	stakingAddress = common.HexToAddress(contract.StakingAddress)
	stakingABI     = contract.MustABIJson(contract.IStakingMetaData.ABI)
)

func GetAddress() common.Address {
	return stakingAddress
}

func GetABI() abi.ABI {
	return stakingABI
}
