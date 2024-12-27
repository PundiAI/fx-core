package types

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type VMCall struct {
	evm      *vm.EVM
	contract *vm.Contract
}

func NewVMCall(evm *vm.EVM, contract *vm.Contract) *VMCall {
	return &VMCall{evm: evm, contract: contract}
}

func (v *VMCall) QueryContract(_ context.Context, from, contract common.Address, abi abi.ABI, method string, res interface{}, args ...interface{}) error {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return types.ErrABIPack.Wrap(err.Error())
	}
	ret, _, err := v.evm.StaticCall(vm.AccountRef(from), contract, data, v.contract.Gas)
	if err != nil {
		return err
	}
	if err = abi.UnpackIntoInterface(res, method, ret); err != nil {
		return types.ErrABIUnpack.Wrap(err.Error())
	}
	return nil
}

func (v *VMCall) ApplyContract(_ context.Context, _, _ common.Address, _ *big.Int, _ abi.ABI, _ string, _ ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	return nil, errors.New("not implemented")
}
