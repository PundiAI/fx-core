package bank

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type TransferFromModuleToAccountMethod struct {
	*Keeper
	TransferFromModuleToAccountABI
	accessControlAddr common.Address
}

func NewTransferFromModuleToAccountMethod(keeper *Keeper) *TransferFromModuleToAccountMethod {
	return &TransferFromModuleToAccountMethod{
		Keeper:                         keeper,
		TransferFromModuleToAccountABI: NewTransferFromModuleToAccountABI(),
		accessControlAddr:              common.Address{}, // TODO: set access control address
	}
}

func (m *TransferFromModuleToAccountMethod) IsReadonly() bool {
	return false
}

func (m *TransferFromModuleToAccountMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *TransferFromModuleToAccountMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *TransferFromModuleToAccountMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if !fxcontract.IsZeroEthAddress(m.accessControlAddr) && m.accessControlAddr != contract.Caller() {
		return nil, errors.New("access denied")
	}
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		return m.TransferFromModuleToAccount(ctx, args)
	}); err != nil {
		return nil, err
	}
	return m.PackOutput(true)
}

type TransferFromModuleToAccountABI struct {
	abi.Method
}

func NewTransferFromModuleToAccountABI() TransferFromModuleToAccountABI {
	return TransferFromModuleToAccountABI{
		Method: bankAbi.Methods["transferFromModuleToAccount"],
	}
}

func (m TransferFromModuleToAccountABI) PackInput(args fxcontract.TransferFromModuleToAccountArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Module, args.Account, args.Token, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m TransferFromModuleToAccountABI) UnpackInput(data []byte) (*fxcontract.TransferFromModuleToAccountArgs, error) {
	args := new(fxcontract.TransferFromModuleToAccountArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m TransferFromModuleToAccountABI) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m TransferFromModuleToAccountABI) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

type TransferFromAccountToModuleMethod struct {
	*Keeper
	TransferFromAccountToModuleABI
	accessControlAddr common.Address
}

func NewTransferFromAccountToModuleMethod(keeper *Keeper) *TransferFromAccountToModuleMethod {
	return &TransferFromAccountToModuleMethod{
		Keeper:                         keeper,
		TransferFromAccountToModuleABI: NewTransferFromAccountToModuleABI(),
		accessControlAddr:              common.Address{}, // TODO: set access control address
	}
}

func (m *TransferFromAccountToModuleMethod) IsReadonly() bool {
	return false
}

func (m *TransferFromAccountToModuleMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *TransferFromAccountToModuleMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *TransferFromAccountToModuleMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	if !fxcontract.IsZeroEthAddress(m.accessControlAddr) && m.accessControlAddr != contract.Caller() {
		return nil, errors.New("access denied")
	}
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		return m.TransferFromAccountToModule(ctx, args)
	}); err != nil {
		return nil, err
	}
	return m.PackOutput(true)
}

type TransferFromAccountToModuleABI struct {
	abi.Method
}

func NewTransferFromAccountToModuleABI() TransferFromAccountToModuleABI {
	return TransferFromAccountToModuleABI{
		Method: bankAbi.Methods["transferFromAccountToModule"],
	}
}

func (m TransferFromAccountToModuleABI) PackInput(args fxcontract.TransferFromAccountToModuleArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Module, args.Account, args.Token, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m TransferFromAccountToModuleABI) UnpackInput(data []byte) (*fxcontract.TransferFromAccountToModuleArgs, error) {
	args := new(fxcontract.TransferFromAccountToModuleArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m TransferFromAccountToModuleABI) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m TransferFromAccountToModuleABI) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}
