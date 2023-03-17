package staking

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

type Contract struct {
	ctx           sdk.Context
	evm           *vm.EVM
	bankKeeper    BankKeeper
	distrKeeper   DistrKeeper
	stakingKeeper StakingKeeper
	evmKeeper     EvmKeeper
}

func NewPrecompiledContract(
	ctx sdk.Context,
	evm *vm.EVM,
	bankKeeper BankKeeper,
	stakingKeeper StakingKeeper,
	distrKeeper DistrKeeper,
	evmKeeper EvmKeeper,
) *Contract {
	return &Contract{
		ctx:           ctx,
		evm:           evm,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		evmKeeper:     evmKeeper,
	}
}

func (c *Contract) Address() common.Address {
	return StakingAddress
}

func (c *Contract) IsStateful() bool {
	return true
}

func (c *Contract) RequiredGas(input []byte) uint64 {
	if len(input) <= 4 {
		return 0
	}
	switch string(input[:4]) {
	case string(DelegateMethod.ID):
		return DelegateGas
	case string(UndelegateMethod.ID):
		return UndelegateGas
	case string(WithdrawMethod.ID):
		return WithdrawGas
	case string(DelegationMethod.ID):
		return DelegationGas
	default:
		return 0
	}
}

func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(contract.Input) <= 4 {
		return types.PackRetError("invalid input")
	}
	// parse input
	switch string(contract.Input[:4]) {
	case string(DelegateMethod.ID):
		ret, err = c.Delegate(evm, contract, readonly)
	case string(UndelegateMethod.ID):
		ret, err = c.Undelegate(evm, contract, readonly)
	case string(WithdrawMethod.ID):
		ret, err = c.Withdraw(evm, contract, readonly)
	case string(DelegationMethod.ID):
		ret, err = c.Delegation(evm, contract, readonly)
	default:
		err = errors.New("unknown method")
	}

	if err != nil {
		return types.PackRetError(err.Error())
	}
	return ret, nil
}
