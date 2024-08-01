package precompile

import (
	"errors"
	"fmt"
	"math/big"
	"time"

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

	fxcontract "github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
)

type RedelegationMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewRedelegationMethod(keeper *Keeper) *RedelegationMethod {
	return &RedelegationMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["redelegate"],
		Event:  fxstakingtypes.GetABI().Events["Redelegate"],
	}
}

func (m *RedelegationMethod) IsReadonly() bool {
	return false
}

func (m *RedelegationMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *RedelegationMethod) RequiredGas() uint64 {
	return 60_000
}

func (m *RedelegationMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error { // withdraw src reward
		sender := sdk.AccAddress(contract.Caller().Bytes())
		shares := sdkmath.LegacyNewDecFromBigInt(args.Shares)

		valSrcAddr := args.GetValidatorSrc()
		// check src validator
		validatorSrc, found := m.stakingKeeper.GetValidator(ctx, valSrcAddr)
		if !found {
			return fmt.Errorf("validator src not found: %s", valSrcAddr.String())
		}

		// check delegation
		delegation, found := m.stakingKeeper.GetDelegation(ctx, sender, valSrcAddr)
		if !found {
			return fmt.Errorf("delegation not found")
		}
		if delegation.Shares.LT(shares) {
			return fmt.Errorf("insufficient shares to redelegate")
		}

		// check dst validator
		valDstAddr := args.GetValidatorDst()
		if _, found = m.stakingKeeper.GetValidator(ctx, valDstAddr); !found {
			return fmt.Errorf("validator dst not found: %s", valDstAddr.String())
		}

		withdrawAddr := m.distrKeeper.GetDelegatorWithdrawAddr(ctx, sender)
		beforeDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)

		// redelegate
		completionTime, err := m.stakingKeeper.BeginRedelegation(ctx, sender, valSrcAddr, valDstAddr, shares)
		if err != nil {
			return err
		}

		redelAmount := validatorSrc.TokensFromShares(shares).TruncateInt()
		afterDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)
		rewardCoin := afterDelBalance.Sub(beforeDelBalance)

		// add redelegate log
		data, topic, err := m.NewRedelegationEvent(contract.Caller(), args.ValidatorSrc, args.ValidatorDst, args.Shares, redelAmount.BigInt(), completionTime.Unix())
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

		// add redelegate event
		RedelegateEmitEvents(ctx, sender, valSrcAddr, valDstAddr, redelAmount, completionTime)

		result, err = m.PackOutput(redelAmount.BigInt(), rewardCoin.Amount.BigInt(), completionTime.Unix())
		return err
	})

	return result, err
}

func (m *RedelegationMethod) NewRedelegationEvent(sender common.Address, validatorSrc, validatorDst string, shares, amount *big.Int, completionTime int64) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validatorSrc, validatorDst, shares, amount, big.NewInt(completionTime))
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *RedelegationMethod) PackInput(args fxstakingtypes.RedelegateArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.ValidatorSrc, args.ValidatorDst, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *RedelegationMethod) UnpackInput(data []byte) (*fxstakingtypes.RedelegateArgs, error) {
	args := new(fxstakingtypes.RedelegateArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *RedelegationMethod) PackOutput(redelAmount, reward *big.Int, completionTime int64) ([]byte, error) {
	return m.Method.Outputs.Pack(redelAmount, reward, big.NewInt(completionTime))
}

func (m *RedelegationMethod) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *RedelegationMethod) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingRedelegate, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseRedelegate(*log)
}

type RedelegateMethodV2 struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewRedelegateV2Method(keeper *Keeper) *RedelegateMethodV2 {
	return &RedelegateMethodV2{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["redelegateV2"],
		Event:  fxstakingtypes.GetABI().Events["RedelegateV2"],
	}
}

func (m *RedelegateMethodV2) IsReadonly() bool {
	return false
}

func (m *RedelegateMethodV2) GetMethodId() []byte {
	return m.Method.ID
}

func (m *RedelegateMethodV2) RequiredGas() uint64 {
	return 60_000
}

func (m *RedelegateMethodV2) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		resp, err := m.stakingMsgServer.BeginRedelegate(sdk.WrapSDKContext(ctx), &stakingtypes.MsgBeginRedelegate{
			DelegatorAddress:    sdk.AccAddress(contract.Caller().Bytes()).String(),
			ValidatorSrcAddress: args.ValidatorSrc,
			ValidatorDstAddress: args.ValidatorDst,
			Amount:              m.NewStakingCoin(args.Amount),
		})
		if err != nil {
			return err
		}

		// add redelegate log
		data, topic, err := m.NewRedelegationEvent(contract.Caller(), args.ValidatorSrc, args.ValidatorDst, args.Amount, resp.CompletionTime.Unix())
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

func (m *RedelegateMethodV2) NewRedelegationEvent(sender common.Address, validatorSrc, validatorDst string, amount *big.Int, completionTime int64) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validatorSrc, validatorDst, amount, big.NewInt(completionTime))
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *RedelegateMethodV2) PackInput(args fxstakingtypes.RedelegateV2Args) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.ValidatorSrc, args.ValidatorDst, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *RedelegateMethodV2) UnpackInput(data []byte) (*fxstakingtypes.RedelegateV2Args, error) {
	args := new(fxstakingtypes.RedelegateV2Args)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *RedelegateMethodV2) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m *RedelegateMethodV2) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *RedelegateMethodV2) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingRedelegateV2, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseRedelegateV2(*log)
}

func RedelegateEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, validatorSrc, validatorDst sdk.ValAddress, amount sdkmath.Int, completionTime time.Time) {
	if amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "redelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", fxtypes.DefaultDenom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeRedelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, validatorSrc.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, validatorDst.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, delegator.String()),
		),
	})
}
