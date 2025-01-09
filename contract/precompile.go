package contract

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type PrecompileMethod interface {
	GetMethodId() []byte
	RequiredGas() uint64
	IsReadonly() bool
	Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error)
}

func EmitEvent(evm *vm.EVM, address common.Address, data []byte, topics []common.Hash) {
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     address,
		Topics:      topics,
		Data:        data,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
}
