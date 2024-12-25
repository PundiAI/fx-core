package staking

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type SlashingInfoMethod struct {
	*Keeper
	SlashingInfoABI
}

func NewSlashingInfoMethod(keeper *Keeper) *SlashingInfoMethod {
	return &SlashingInfoMethod{
		Keeper:          keeper,
		SlashingInfoABI: NewSlashingInfoABI(),
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

type SlashingInfoABI struct {
	abi.Method
}

func NewSlashingInfoABI() SlashingInfoABI {
	return SlashingInfoABI{
		Method: stakingABI.Methods["slashingInfo"],
	}
}

func (m SlashingInfoABI) PackInput(args fxcontract.SlashingInfoArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m SlashingInfoABI) UnpackInput(data []byte) (*fxcontract.SlashingInfoArgs, error) {
	args := new(fxcontract.SlashingInfoArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m SlashingInfoABI) PackOutput(jailed bool, missed int64) ([]byte, error) {
	return m.Method.Outputs.Pack(jailed, big.NewInt(missed))
}

func (m SlashingInfoABI) UnpackOutput(data []byte) (bool, *big.Int, error) {
	unpack, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, nil, err
	}
	return unpack[0].(bool), unpack[1].(*big.Int), nil
}
