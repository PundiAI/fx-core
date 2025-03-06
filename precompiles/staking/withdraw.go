package staking

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	fxcontract "github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/x/evm/types"
)

type WithdrawMethod struct {
	*Keeper
	WithdrawABI
}

func NewWithdrawMethod(keeper *Keeper) *WithdrawMethod {
	return &WithdrawMethod{
		Keeper:      keeper,
		WithdrawABI: NewWithdrawABI(),
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
	if contract.Value().Sign() != 0 {
		return nil, errors.New("msg.value must be zero")
	}

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

type WithdrawABI struct {
	abi.Method
	abi.Event
}

func NewWithdrawABI() WithdrawABI {
	return WithdrawABI{
		Method: stakingABI.Methods["withdraw"],
		Event:  stakingABI.Events["Withdraw"],
	}
}

func (m WithdrawABI) NewWithdrawEvent(sender common.Address, validator string, reward *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, reward)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m WithdrawABI) PackInput(args fxcontract.WithdrawArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator)
	if err != nil {
		return nil, err
	}
	return append(m.Method.ID, arguments...), nil
}

func (m WithdrawABI) UnpackInput(data []byte) (*fxcontract.WithdrawArgs, error) {
	args := new(fxcontract.WithdrawArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m WithdrawABI) PackOutput(reward *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(reward)
}

func (m WithdrawABI) UnpackOutput(data []byte) (*big.Int, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return nil, err
	}
	return amount[0].(*big.Int), nil
}

func (m WithdrawABI) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingWithdraw, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseWithdraw(*log)
}
