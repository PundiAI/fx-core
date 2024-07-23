package contract

import (
	"github.com/ethereum/go-ethereum/core/vm"
)

type PrecompileMethod interface {
	GetMethodId() []byte
	RequiredGas() uint64
	IsReadonly() bool
	Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error)
}
