package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type DelegationRewardsMethod struct {
	*Keeper
	abi.Method
}

func NewDelegationRewardsMethod(keeper *Keeper) *DelegationRewardsMethod {
	return &DelegationRewardsMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["delegationRewards"],
	}
}

func (m *DelegationRewardsMethod) IsReadonly() bool {
	return true
}

func (m *DelegationRewardsMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *DelegationRewardsMethod) RequiredGas() uint64 {
	return 30_000
}

func (m *DelegationRewardsMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}
	stateDB := evm.StateDB.(types.ExtStateDB)
	cacheCtx := stateDB.CacheContext()

	valAddr := args.GetValidator()
	validator, found := m.stakingKeeper.GetValidator(cacheCtx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}
	delegation, found := m.stakingKeeper.GetDelegation(cacheCtx, args.Delegator.Bytes(), valAddr)
	if !found {
		return m.PackOutput(big.NewInt(0))
	}

	endingPeriod := m.distrKeeper.IncrementValidatorPeriod(cacheCtx, validator)
	rewards := m.distrKeeper.CalculateDelegationRewards(cacheCtx, validator, delegation, endingPeriod)

	return m.PackOutput(rewards.AmountOf(m.stakingDenom).TruncateInt().BigInt())
}

func (m *DelegationRewardsMethod) PackInput(args fxstakingtypes.DelegationRewardsArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Delegator)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *DelegationRewardsMethod) UnpackInput(data []byte) (*fxstakingtypes.DelegationRewardsArgs, error) {
	args := new(fxstakingtypes.DelegationRewardsArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *DelegationRewardsMethod) PackOutput(amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount)
}

func (m *DelegationRewardsMethod) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}
