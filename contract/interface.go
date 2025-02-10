package contract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type EvmKeeper interface {
	Caller
	DeployContract(ctx sdk.Context, from common.Address, abi abi.ABI, bin []byte, args ...interface{}) (common.Address, error)
	CreateContractWithCode(ctx sdk.Context, address common.Address, code []byte) error
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}
