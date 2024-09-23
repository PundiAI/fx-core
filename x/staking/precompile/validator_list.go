package precompile

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type ValidatorListMethod struct {
	*Keeper
	abi.Method
}

type ValidatorList struct {
	ValAddr      string
	MissedBlocks int64
}

func NewValidatorListMethod(keeper *Keeper) *ValidatorListMethod {
	return &ValidatorListMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["validatorList"],
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
		valAddrs = m.ValidatorListPower(bondedVals)
	case fxstakingtypes.ValidatorSortByMissed:
		valAddrs, err = m.ValidatorListMissedBlock(cacheCtx, bondedVals)
		if err != nil {
			return nil, err
		}
	}

	return m.PackOutput(valAddrs)
}

func (m *ValidatorListMethod) ValidatorListPower(bondedVals []stakingtypes.Validator) []string {
	valAddrs := make([]string, 0, len(bondedVals))
	for _, val := range bondedVals {
		valAddrs = append(valAddrs, val.OperatorAddress)
	}
	return valAddrs
}

func (m *ValidatorListMethod) ValidatorListMissedBlock(ctx sdk.Context, bondedVals []stakingtypes.Validator) ([]string, error) {
	valList := make([]ValidatorList, 0, len(bondedVals))
	for _, val := range bondedVals {
		consAddr, err := val.GetConsAddr()
		if err != nil {
			return nil, err
		}
		info, err := m.slashingKeeper.GetValidatorSigningInfo(ctx, consAddr)
		if err != nil {
			return nil, err
		}
		valList = append(valList, ValidatorList{
			ValAddr:      val.OperatorAddress,
			MissedBlocks: info.MissedBlocksCounter,
		})
	}
	sort.Slice(valList, func(i, j int) bool {
		return valList[i].MissedBlocks > valList[j].MissedBlocks
	})
	valAddrs := make([]string, 0, len(valList))
	for _, l := range valList {
		valAddrs = append(valAddrs, l.ValAddr)
	}
	return valAddrs, nil
}

func (m *ValidatorListMethod) PackInput(args fxstakingtypes.ValidatorListArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.SortBy)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *ValidatorListMethod) UnpackInput(data []byte) (*fxstakingtypes.ValidatorListArgs, error) {
	args := new(fxstakingtypes.ValidatorListArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *ValidatorListMethod) PackOutput(valList []string) ([]byte, error) {
	return m.Method.Outputs.Pack(valList)
}

func (m *ValidatorListMethod) UnpackOutput(data []byte) ([]string, error) {
	unpack, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return unpack[0].([]string), nil
}
