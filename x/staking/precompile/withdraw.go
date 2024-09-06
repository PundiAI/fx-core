package precompile

import (
	"errors"
	"math/big"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

func WithdrawEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, amount sdk.Coins) {
	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	)
}

type WithdrawMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewWithdrawMethod(keeper *Keeper) *WithdrawMethod {
	return &WithdrawMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["withdraw"],
		Event:  fxstakingtypes.GetABI().Events["Withdraw"],
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

		resp, err := m.distrMsgServer.WithdrawDelegatorReward(sdk.WrapSDKContext(ctx), &distrtypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender.String(),
			ValidatorAddress: args.GetValidator().String(),
		})
		if err != nil {
			return err
		}
		// add withdraw event
		WithdrawEmitEvents(ctx, sender, resp.Amount)

		// add withdraw log
		bigInt := resp.Amount.AmountOf(m.stakingDenom).BigInt()
		data, topic, err := m.NewWithdrawEvent(contract.Caller(), args.GetValidator().String(), bigInt)
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

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

func (m *WithdrawMethod) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
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
