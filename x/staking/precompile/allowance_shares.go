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
	AllowanceSharesABI
}

func NewAllowanceSharesMethod(keeper *Keeper) *AllowanceSharesMethod {
	return &AllowanceSharesMethod{
		Keeper:             keeper,
		AllowanceSharesABI: NewAllowanceSharesABI(),
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

type AllowanceSharesABI struct {
	abi.Method
}

func NewAllowanceSharesABI() AllowanceSharesABI {
	return AllowanceSharesABI{
		Method: stakingABI.Methods["allowanceShares"],
	}
}

func (m AllowanceSharesABI) PackInput(args fxstakingtypes.AllowanceSharesArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Owner, args.Spender)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m AllowanceSharesABI) UnpackInput(data []byte) (*fxstakingtypes.AllowanceSharesArgs, error) {
	args := new(fxstakingtypes.AllowanceSharesArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m AllowanceSharesABI) PackOutput(amount *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(amount)
}

func (m AllowanceSharesABI) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}
