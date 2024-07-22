package precompile

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/functionx/fx-core/v7/x/evm/types"
)

type Contract struct {
	bankKeeper    BankKeeper
	distrKeeper   DistrKeeper
	stakingKeeper StakingKeeper
	evmKeeper     EvmKeeper
}

func NewPrecompiledContract(bankKeeper BankKeeper, stakingKeeper StakingKeeper, distrKeeper DistrKeeper, evmKeeper EvmKeeper) *Contract {
	return &Contract{
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		evmKeeper:     evmKeeper,
	}
}

func (c *Contract) Address() common.Address {
	return stakingAddress
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
	case string(RedelegateMethod.ID):
		return RedelegateGas
	case string(WithdrawMethod.ID):
		return WithdrawGas
	case string(DelegationMethod.ID):
		return DelegationGas
	case string(DelegationRewardsMethod.ID):
		return DelegationRewardsGas
	case string(TransferSharesMethod.ID):
		return TransferSharesGas
	case string(ApproveSharesMethod.ID):
		return ApproveSharesGas
	case string(AllowanceSharesMethod.ID):
		return AllowanceSharesGas
	case string(TransferFromSharesMethod.ID):
		return TransferFromSharesGas
	default:
		return 0
	}
}

//gocyclo:ignore
func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) (ret []byte, err error) {
	if len(contract.Input) <= 4 {
		return types.PackRetError(errors.New("invalid input"))
	}

	switch string(contract.Input[:4]) {
	case string(DelegateMethod.ID):
		ret, err = c.Delegate(evm, contract, readonly)
	case string(UndelegateMethod.ID):
		ret, err = c.Undelegate(evm, contract, readonly)
	case string(RedelegateMethod.ID):
		ret, err = c.Redelegation(evm, contract, readonly)
	case string(WithdrawMethod.ID):
		ret, err = c.Withdraw(evm, contract, readonly)
	case string(DelegationMethod.ID):
		ret, err = c.Delegation(evm, contract, readonly)
	case string(DelegationRewardsMethod.ID):
		ret, err = c.DelegationRewards(evm, contract, readonly)
	case string(TransferSharesMethod.ID):
		ret, err = c.TransferShares(evm, contract, readonly)
	case string(ApproveSharesMethod.ID):
		ret, err = c.ApproveShares(evm, contract, readonly)
	case string(AllowanceSharesMethod.ID):
		ret, err = c.AllowanceShares(evm, contract, readonly)
	case string(TransferFromSharesMethod.ID):
		ret, err = c.TransferFromShares(evm, contract, readonly)
	case string(DelegateV2Method.ID):
		ret, err = c.DelegateV2(evm, contract, readonly)
	case string(UndelegateV2Method.ID):
		ret, err = c.UndelegateV2(evm, contract, readonly)
	case string(RedelegateV2Method.ID):
		ret, err = c.RedelegationV2(evm, contract, readonly)

	default:
		err = errors.New("unknown method")
	}

	if err != nil {
		return types.PackRetError(err)
	}

	return ret, nil
}

func (c *Contract) AddLog(evm *vm.EVM, event abi.Event, topics []common.Hash, args ...interface{}) error {
	data, newTopic, err := types.PackTopicData(event, topics, args...)
	if err != nil {
		return err
	}
	evm.StateDB.AddLog(&ethtypes.Log{
		Address:     c.Address(),
		Topics:      newTopic,
		Data:        data,
		BlockNumber: evm.Context.BlockNumber.Uint64(),
	})
	return nil
}
