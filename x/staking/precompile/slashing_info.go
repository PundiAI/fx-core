package precompile

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type SlashingInfoMethod struct {
	*Keeper
	abi.Method
}

func NewSlashingInfoMethod(keeper *Keeper) *SlashingInfoMethod {
	return &SlashingInfoMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["slashingInfo"],
	}
}

func (m *SlashingInfoMethod) IsReadonly() bool {
	return true
}

func (m *SlashingInfoMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *SlashingInfoMethod) RequiredGas() uint64 {
	return 1_000
}

func (m *SlashingInfoMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	cacheCtx := stateDB.Context()

	validator, err := m.Keeper.stakingKeeper.GetValidator(cacheCtx, args.GetValidator())
	if err != nil {
		return nil, err
	}

	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return nil, err
	}

	signingInfo, err := m.Keeper.slashingKeeper.GetValidatorSigningInfo(cacheCtx, consAddr)
	if err != nil {
		return nil, err
	}
	return m.PackOutput(validator.Jailed, signingInfo.MissedBlocksCounter)
}

func (m *SlashingInfoMethod) PackInput(args fxstakingtypes.SlashingInfoArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *SlashingInfoMethod) UnpackInput(data []byte) (*fxstakingtypes.SlashingInfoArgs, error) {
	args := new(fxstakingtypes.SlashingInfoArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *SlashingInfoMethod) PackOutput(jailed bool, missed int64) ([]byte, error) {
	return m.Method.Outputs.Pack(jailed, big.NewInt(missed))
}

func (m *SlashingInfoMethod) UnpackOutput(data []byte) (bool, *big.Int, error) {
	unpack, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, nil, err
	}
	return unpack[0].(bool), unpack[1].(*big.Int), nil
}
