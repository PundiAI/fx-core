package precompile

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

type DelegateMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewDelegateMethod(keeper *Keeper) *DelegateMethod {
	return &DelegateMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["delegate"],
		Event:  fxstakingtypes.GetABI().Events["Delegate"],
	}
}

func (m *DelegateMethod) IsReadonly() bool {
	return false
}

func (m *DelegateMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *DelegateMethod) RequiredGas() uint64 {
	return 40_000
}

func (m *DelegateMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	amount := contract.Value()
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid delegate amount: %s", amount.String())
	}
	valAddr := args.GetValidator()

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		val, found := m.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return fmt.Errorf("validator not found: %s", valAddr.String())
		}

		sender := sdk.AccAddress(contract.Caller().Bytes())
		delCoin := m.NewStakingCoin(amount)
		if err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, contract.Address().Bytes(), evmtypes.ModuleName, sdk.NewCoins(delCoin)); err != nil {
			return err
		}
		if err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, evmtypes.ModuleName, sender, sdk.NewCoins(delCoin)); err != nil {
			return err
		}

		withdrawAddr := m.distrKeeper.GetDelegatorWithdrawAddr(ctx, sender)
		beforeDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)
		if withdrawAddr.Equals(sender) {
			beforeDelBalance = beforeDelBalance.Sub(delCoin)
		}

		// delegate amount
		shares, err := m.stakingKeeper.Delegate(ctx, sender, sdkmath.NewIntFromBigInt(amount), stakingtypes.Unbonded, val, true)
		if err != nil {
			return err
		}

		afterDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)
		rewardCoin := afterDelBalance.Sub(beforeDelBalance)

		// add delegate event
		DelegateEmitEvents(ctx, sender, valAddr, amount, shares)

		data, topic, err := m.NewDelegateEvent(contract.Caller(), args.Validator, amount, shares.TruncateInt().BigInt())
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

		result, err = m.PackOutput(shares.TruncateInt().BigInt(), rewardCoin.Amount.BigInt())
		return err
	})

	return result, err
}

func (m *DelegateMethod) NewDelegateEvent(sender common.Address, validator string, amount, shares *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, amount, shares)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *DelegateMethod) PackInput(args fxstakingtypes.DelegateArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *DelegateMethod) UnpackInput(data []byte) (*fxstakingtypes.DelegateArgs, error) {
	args := new(fxstakingtypes.DelegateArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *DelegateMethod) PackOutput(share, reward *big.Int) ([]byte, error) {
	return m.Method.Outputs.Pack(share, reward)
}

func (m *DelegateMethod) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *DelegateMethod) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingDelegate, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseDelegate(*log)
}

type DelegateV2Method struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewDelegateV2Method(keeper *Keeper) *DelegateV2Method {
	return &DelegateV2Method{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["delegateV2"],
		Event:  fxstakingtypes.GetABI().Events["DelegateV2"],
	}
}

func (m *DelegateV2Method) IsReadonly() bool {
	return false
}

func (m *DelegateV2Method) GetMethodId() []byte {
	return m.Method.ID
}

func (m *DelegateV2Method) RequiredGas() uint64 {
	return 40_000
}

func (m *DelegateV2Method) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)
	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		if _, err = m.stakingMsgServer.Delegate(sdk.WrapSDKContext(ctx), &stakingtypes.MsgDelegate{
			DelegatorAddress: sdk.AccAddress(contract.Caller().Bytes()).String(),
			ValidatorAddress: args.Validator,
			Amount:           m.NewStakingCoin(args.Amount),
		}); err != nil {
			return err
		}

		// add delegate log
		data, topic, err := m.NewDelegateEvent(contract.Caller(), args.Validator, args.Amount)
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)
		return nil
	}); err != nil {
		return nil, err
	}
	return m.PackOutput(true)
}

func (m *DelegateV2Method) NewDelegateEvent(sender common.Address, validator string, amount *big.Int) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, amount)
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *DelegateV2Method) PackInput(args fxstakingtypes.DelegateV2Args) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *DelegateV2Method) UnpackInput(data []byte) (*fxstakingtypes.DelegateV2Args, error) {
	args := new(fxstakingtypes.DelegateV2Args)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *DelegateV2Method) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m *DelegateV2Method) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *DelegateV2Method) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingDelegateV2, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseDelegateV2(*log)
}

func DelegateEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, amount *big.Int, newShares sdk.Dec) {
	if amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", fxtypes.DefaultDenom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	})
}
