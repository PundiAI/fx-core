package precompile

import (
	"errors"
	"math/big"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
)

type DelegationMethod struct {
	*Keeper
	DelegationABI
}

func NewDelegationMethod(keeper *Keeper) *DelegationMethod {
	return &DelegationMethod{
		Keeper:        keeper,
		DelegationABI: NewDelegationABI(),
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

type DelegationABI struct {
	abi.Method
}

func NewDelegationABI() DelegationABI {
	return DelegationABI{
		Method: stakingABI.Methods["delegation"],
	}
}

func (m DelegationABI) PackInput(args fxcontract.DelegationArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Delegator)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m DelegationABI) UnpackInput(data []byte) (*fxcontract.DelegationArgs, error) {
	args := new(fxcontract.DelegationArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m DelegationABI) PackOutput(shares, amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(shares, amount)
}

func (m DelegationABI) UnpackOutput(data []byte) (*big.Int, *big.Int, error) {
	unpack, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, nil, err
	}
	shares := unpack[0].(*big.Int)
	amount := unpack[1].(*big.Int)
	return shares, amount, nil
}
