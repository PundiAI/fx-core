package precompile

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type ValidatorListMethod struct {
	*Keeper
	ValidatorListABI
}

func NewValidatorListMethod(keeper *Keeper) *ValidatorListMethod {
	return &ValidatorListMethod{
		Keeper:           keeper,
		ValidatorListABI: NewValidatorListABI(),
	}
}

func (m *ValidatorListMethod) IsReadonly() bool {
	return true
}

func (m *ValidatorListMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *ValidatorListMethod) RequiredGas() uint64 {
	return 1_000
}

func (m *ValidatorListMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	cacheCtx := stateDB.Context()

	bondedVals, err := m.stakingKeeper.GetLastValidators(cacheCtx)
	if err != nil {
		return nil, err
	}

	valAddrs := make([]string, 0, len(bondedVals))
	switch args.GetSortBy() {
	case fxstakingtypes.ValidatorSortByPower:
		valAddrs = validatorListPower(bondedVals)
	case fxstakingtypes.ValidatorSortByMissed:
		valAddrs, err = m.ValidatorListMissedBlock(cacheCtx, bondedVals)
		if err != nil {
			return nil, err
		}
	}

	return m.PackOutput(valAddrs)
}

type ValidatorListABI struct {
	abi.Method
}

func NewValidatorListABI() ValidatorListABI {
	return ValidatorListABI{
		Method: stakingABI.Methods["validatorList"],
	}
}

func (m ValidatorListABI) PackInput(args fxstakingtypes.ValidatorListArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.SortBy)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m ValidatorListABI) UnpackInput(data []byte) (*fxstakingtypes.ValidatorListArgs, error) {
	args := new(fxstakingtypes.ValidatorListArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m ValidatorListABI) PackOutput(valList []string) ([]byte, error) {
	return m.Method.Outputs.Pack(valList)
}

func (m ValidatorListABI) UnpackOutput(data []byte) ([]string, error) {
	unpack, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return unpack[0].([]string), nil
}
