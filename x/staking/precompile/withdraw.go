package precompile

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type WithdrawMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewWithdrawMethod(keeper *Keeper) *WithdrawMethod {
	return &WithdrawMethod{
		Keeper: keeper,
		Method: stakingABI.Methods["withdraw"],
		Event:  stakingABI.Events["Withdraw"],
	}
}

func (m *WithdrawMethod) IsReadonly() bool {
	return false
}

func (m *WithdrawMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *WithdrawMethod) RequiredGas() uint64 {
	return 30_000
}

func (m *WithdrawMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		sender := sdk.AccAddress(contract.Caller().Bytes())

		resp, err := m.distrMsgServer.WithdrawDelegatorReward(ctx, &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender.String(),
			ValidatorAddress: args.GetValidator().String(),
		})
		if err != nil {
			return err
		}

		// add withdraw log
		bigInt := resp.Amount.AmountOf(m.stakingDenom).BigInt()
		data, topic, err := m.NewWithdrawEvent(contract.Caller(), args.GetValidator().String(), bigInt)
		if err != nil {
			return err
		}
		fxcontract.EmitEvent(evm, stakingAddress, data, topic)

		result, err = m.PackOutput(bigInt)
		return err
	})

	return result, err
}

func (m *WithdrawMethod) NewWithdrawEvent(sender common.Address, validator string, reward *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, reward)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *WithdrawMethod) PackInput(args fxstakingtypes.WithdrawArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *WithdrawMethod) UnpackInput(data []byte) (*fxstakingtypes.WithdrawArgs, error) {
	args := new(fxstakingtypes.WithdrawArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *WithdrawMethod) PackOutput(reward *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(reward)
}

func (m *WithdrawMethod) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}

func (m *WithdrawMethod) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingWithdraw, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseWithdraw(*log)
}
