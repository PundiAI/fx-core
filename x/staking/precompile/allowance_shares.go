package precompile

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type AllowanceSharesMethod struct {
	*Keeper
	abi.Method
}

func NewAllowanceSharesMethod(keeper *Keeper) *AllowanceSharesMethod {
	return &AllowanceSharesMethod{
		Keeper: keeper,
		Method: stakingABI.Methods["allowanceShares"],
	}
}

func (m *AllowanceSharesMethod) IsReadonly() bool {
	return true
}

func (m *AllowanceSharesMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *AllowanceSharesMethod) RequiredGas() uint64 {
	return 5_000
}

func (m *AllowanceSharesMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	allowance := m.Keeper.stakingKeeper.GetAllowance(stateDB.Context(), args.GetValidator(), args.Owner.Bytes(), args.Spender.Bytes())
	return m.PackOutput(allowance)
}

func (m *AllowanceSharesMethod) PackInput(args fxstakingtypes.AllowanceSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Owner, args.Spender)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *AllowanceSharesMethod) UnpackInput(data []byte) (*fxstakingtypes.AllowanceSharesArgs, error) {
	args := new(fxstakingtypes.AllowanceSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *AllowanceSharesMethod) PackOutput(amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount)
}

func (m *AllowanceSharesMethod) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}
