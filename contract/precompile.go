package contract

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type PrecompileMethod interface {
	GetMethodId() []byte
	RequiredGas() uint64
	IsReadonly() bool
	Run(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract) ([]byte, error)
}
