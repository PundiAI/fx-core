package precompiles

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type VMCall struct {
	evm    *vm.EVM
	maxGas uint64
}

func NewVMCall(evm *vm.EVM) *VMCall {
	return &VMCall{evm: evm, maxGas: contract.DefaultGasCap / 10}
}

func (v *VMCall) QueryContract(_ context.Context, from, contract common.Address, abi abi.ABI, method string, res interface{}, args ...interface{}) error {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return types.ErrABIPack.Wrap(err.Error())
	}
	ret, _, err := v.evm.StaticCall(vm.AccountRef(from), contract, data, v.maxGas)
	if err != nil {
		return err
	}
	if err = abi.UnpackIntoInterface(res, method, ret); err != nil {
		return types.ErrABIUnpack.Wrap(err.Error())
	}
	return nil
}

func (v *VMCall) ApplyContract(_ context.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, constructorData...)
	if err != nil {
		return nil, types.ErrABIPack.Wrap(err.Error())
	}
	ret, leftoverGas, vmErr := v.evm.Call(vm.AccountRef(from), contract, data, v.maxGas, value)
	var vmError string
	if vmErr != nil {
		vmError = vmErr.Error()
	}
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: v.maxGas - leftoverGas,
		VmError: vmError,
		Ret:     ret,
	}, nil
}
