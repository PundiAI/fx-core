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

type DelegationMethod struct {
	*Keeper
	abi.Method
}

func NewDelegationMethod(keeper *Keeper) *DelegationMethod {
	return &DelegationMethod{
		Keeper: keeper,
		Method: stakingABI.Methods["delegation"],
	}
}

func (m *DelegationMethod) IsReadonly() bool {
	return true
}

func (m *DelegationMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *DelegationMethod) RequiredGas() uint64 {
	return 30_000
}

func (m *DelegationMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	ctx := stateDB.Context()

	valAddr := args.GetValidator()
	validator, err := m.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	delegation, err := m.stakingKeeper.GetDelegation(ctx, args.Delegator.Bytes(), valAddr)
	if err != nil {
		if !errors.Is(err, stakingtypes.ErrNoDelegation) {
			return nil, err
		}
		return m.PackOutput(big.NewInt(0), big.NewInt(0))
	}

	delegationAmt := delegation.GetShares().MulInt(validator.GetTokens()).Quo(validator.GetDelegatorShares())
	return m.PackOutput(delegation.GetShares().TruncateInt().BigInt(), delegationAmt.TruncateInt().BigInt())
}

func (m *DelegationMethod) PackInput(args fxstakingtypes.DelegationArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Delegator)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *DelegationMethod) UnpackInput(data []byte) (*fxstakingtypes.DelegationArgs, error) {
	args := new(fxstakingtypes.DelegationArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *DelegationMethod) PackOutput(shares, amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(shares, amount)
}

func (m *DelegationMethod) UnpackOutput(data []byte) (*big.Int, *big.Int, error) {
	unpack, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, nil, err
	}
	shares := unpack[0].(*big.Int)
	amount := unpack[1].(*big.Int)
	return shares, amount, nil
}
