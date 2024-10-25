package precompile

import (
	"errors"
	"math/big"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type DelegationRewardsMethod struct {
	*Keeper
	DelegationRewardsABI
}

func NewDelegationRewardsMethod(keeper *Keeper) *DelegationRewardsMethod {
	return &DelegationRewardsMethod{
		Keeper:               keeper,
		DelegationRewardsABI: NewDelegationRewardsABI(),
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
	cacheCtx := stateDB.Context()

	valAddr := args.GetValidator()
	validator, err := m.stakingKeeper.GetValidator(cacheCtx, valAddr)
	if err != nil {
		return nil, err
	}
	delegation, err := m.stakingKeeper.GetDelegation(cacheCtx, args.Delegator.Bytes(), valAddr)
	if err != nil {
		if !errors.Is(err, stakingtypes.ErrNoDelegation) {
			return nil, err
		}
		return m.PackOutput(big.NewInt(0))
	}

	endingPeriod, err := m.distrKeeper.IncrementValidatorPeriod(cacheCtx, validator)
	if err != nil {
		return nil, err
	}
	rewards, err := m.distrKeeper.CalculateDelegationRewards(cacheCtx, validator, delegation, endingPeriod)
	if err != nil {
		return nil, err
	}

	return m.PackOutput(rewards.AmountOf(m.stakingDenom).TruncateInt().BigInt())
}

type DelegationRewardsABI struct {
	abi.Method
}

func NewDelegationRewardsABI() DelegationRewardsABI {
	return DelegationRewardsABI{
		Method: stakingABI.Methods["delegationRewards"],
	}
}

func (m DelegationRewardsABI) PackInput(args fxstakingtypes.DelegationRewardsArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Delegator)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m DelegationRewardsABI) UnpackInput(data []byte) (*fxstakingtypes.DelegationRewardsArgs, error) {
	args := new(fxstakingtypes.DelegationRewardsArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m DelegationRewardsABI) PackOutput(amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount)
}

func (m DelegationRewardsABI) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}
