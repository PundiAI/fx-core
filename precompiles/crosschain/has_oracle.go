package crosschain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/contract"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/evm/types"
)

type HasOracleMethod struct {
	*Keeper
	HasOracleABI
}

func NewHasOracleMethod(keeper *Keeper) *HasOracleMethod {
	return &HasOracleMethod{
		Keeper:       keeper,
		HasOracleABI: NewHasOracleABI(),
	}
}

func (m *HasOracleMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *HasOracleMethod) RequiredGas() uint64 {
	return 1_000
}

func (m *HasOracleMethod) IsReadonly() bool {
	return true
}

func (m *HasOracleMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}
	stateDB := evm.StateDB.(types.ExtStateDB)

	router, has := m.Keeper.router.GetRoute(args.Chain)
	if !has {
		return nil, fmt.Errorf("chain not support: %s", args.Chain)
	}
	hasOracle := router.HasOracleAddrByExternalAddr(stateDB.Context(), crosschaintypes.ExternalAddrToStr(args.Chain, args.ExternalAddress.Bytes()))
	return m.PackOutput(hasOracle)
}

type HasOracleABI struct {
	abi.Method
}

func NewHasOracleABI() HasOracleABI {
	return HasOracleABI{
		Method: crosschainABI.Methods["hasOracle"],
	}
}

func (m HasOracleABI) PackInput(args contract.HasOracleArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Chain, args.ExternalAddress)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m HasOracleABI) UnpackInput(data []byte) (*contract.HasOracleArgs, error) {
	args := new(contract.HasOracleArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m HasOracleABI) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m HasOracleABI) UnpackOutput(data []byte) (bool, error) {
	result, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return result[0].(bool), nil
}
