package precompiles

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (v *VMCall) ApplyContract(_ context.Context, from, contract common.Address, value *big.Int, cABI abi.ABI, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := cABI.Pack(method, args...)
	if err != nil {
		return nil, types.ErrABIPack.Wrap(err.Error())
	}
	res, err := v.ExecuteEVM(sdk.Context{}, from, &contract, value, v.maxGas, data)
	if err != nil {
		return res, err
	}
	if res.Failed() {
		errStr := res.VmError
		if res.VmError == vm.ErrExecutionReverted.Error() {
			if cause, err := abi.UnpackRevert(common.CopyBytes(res.Ret)); err == nil {
				errStr = cause
			}
		}
		return res, evmtypes.ErrVMExecution.Wrap(errStr)
	}
	return res, nil
}

func (v *VMCall) ExecuteEVM(_ sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte) (*evmtypes.MsgEthereumTxResponse, error) {
	if value == nil {
		value = big.NewInt(0)
	}
	if contract == nil {
		contract = &common.Address{}
	}
	ret, leftoverGas, vmErr := v.evm.Call(vm.AccountRef(from), *contract, data, gasLimit, value)
	var vmError string
	if vmErr != nil {
		vmError = vmErr.Error()
	}
	// the error handling is done in the caller
	return &evmtypes.MsgEthereumTxResponse{
		GasUsed: gasLimit - leftoverGas,
		VmError: vmError,
		Ret:     ret,
	}, nil
}
