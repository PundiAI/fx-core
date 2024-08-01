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

type UndelegateMethod struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewUndelegateMethod(keeper *Keeper) *UndelegateMethod {
	return &UndelegateMethod{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["undelegate"],
		Event:  fxstakingtypes.GetABI().Events["Undelegate"],
	}
}

func (m *UndelegateMethod) IsReadonly() bool {
	return false
}

func (m *UndelegateMethod) GetMethodId() []byte {
	return m.Method.ID
}

func (m *UndelegateMethod) RequiredGas() uint64 {
	return 45_000
}

func (m *UndelegateMethod) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		_, found := m.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return fmt.Errorf("validator not found: %s", valAddr.String())
		}
		sender := sdk.AccAddress(contract.Caller().Bytes())
		withdrawAddr := m.distrKeeper.GetDelegatorWithdrawAddr(ctx, sender)
		beforeDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)

		unDelAmount, completionTime, err := Undelegate(ctx, m.stakingKeeper, m.bankKeeper, sender, valAddr, sdkmath.LegacyNewDecFromBigInt(args.Shares), m.stakingDenom)
		if err != nil {
			return fmt.Errorf("undelegate failed: %s", err.Error())
		}

		afterDelBalance := m.bankKeeper.GetBalance(ctx, withdrawAddr, m.stakingDenom)
		rewardCoin := afterDelBalance.Sub(beforeDelBalance)

		// add undelegate log
		data, topic, err := m.NewUndelegateEvent(contract.Caller(), args.Validator, args.Shares, unDelAmount.BigInt(), completionTime.Unix())
		if err != nil {
			return err
		}
		EmitEvent(evm, data, topic)

		// add undelegate event
		UndelegateEmitEvents(ctx, sender, valAddr, unDelAmount, completionTime)

		result, err = m.PackOutput(unDelAmount.BigInt(), rewardCoin.Amount.BigInt(), completionTime.Unix())
		return err
	})
	return result, err
}

func (m *UndelegateMethod) NewUndelegateEvent(sender common.Address, validator string, shares, amount *big.Int, completionTime int64) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, shares, amount, big.NewInt(completionTime))
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *UndelegateMethod) PackInput(args fxstakingtypes.UndelegateArgs) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Shares)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *UndelegateMethod) UnpackInput(data []byte) (*fxstakingtypes.UndelegateArgs, error) {
	args := new(fxstakingtypes.UndelegateArgs)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *UndelegateMethod) PackOutput(undelAmount, reward *big.Int, completionTime int64) ([]byte, error) {
	return m.Method.Outputs.Pack(undelAmount, reward, big.NewInt(completionTime))
}

func (m *UndelegateMethod) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *UndelegateMethod) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingUndelegate, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseUndelegate(*log)
}

type UndelegateV2Method struct {
	*Keeper
	abi.Method
	abi.Event
}

func NewUndelegateV2Method(keeper *Keeper) *UndelegateV2Method {
	return &UndelegateV2Method{
		Keeper: keeper,
		Method: fxstakingtypes.GetABI().Methods["undelegateV2"],
		Event:  fxstakingtypes.GetABI().Events["UndelegateV2"],
	}
}

func (m *UndelegateV2Method) IsReadonly() bool {
	return false
}

func (m *UndelegateV2Method) GetMethodId() []byte {
	return m.Method.ID
}

func (m *UndelegateV2Method) RequiredGas() uint64 {
	return 45_000
}

func (m *UndelegateV2Method) Run(evm *vm.EVM, contract *vm.Contract) ([]byte, error) {
	args, err := m.UnpackInput(contract.Input)
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	if err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		resp, err := m.stakingMsgServer.Undelegate(sdk.WrapSDKContext(ctx), &stakingtypes.MsgUndelegate{
			DelegatorAddress: sdk.AccAddress(contract.Caller().Bytes()).String(),
			ValidatorAddress: args.Validator,
			Amount:           m.NewStakingCoin(args.Amount),
		})
		if err != nil {
			return err
		}

		// add undelegate log
		data, topic, err := m.NewUndelegateEvent(contract.Caller(), args.Validator, args.Amount, resp.CompletionTime.Unix())
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

func (m *UndelegateV2Method) NewUndelegateEvent(sender common.Address, validator string, amount *big.Int, completionTime int64) (data []byte, topic []common.Hash, err error) {
	data, topic, err = types.PackTopicData(m.Event, []common.Hash{sender.Hash()}, validator, amount, big.NewInt(completionTime))
	if err != nil {
		return nil, nil, err
	}
	return data, topic, nil
}

func (m *UndelegateV2Method) PackInput(args fxstakingtypes.UndelegateV2Args) ([]byte, error) {
	arguments, err := m.Method.Inputs.Pack(args.Validator, args.Amount)
	if err != nil {
		return nil, err
	}
	return append(m.GetMethodId(), arguments...), nil
}

func (m *UndelegateV2Method) UnpackInput(data []byte) (*fxstakingtypes.UndelegateV2Args, error) {
	args := new(fxstakingtypes.UndelegateV2Args)
	err := types.ParseMethodArgs(m.Method, args, data[4:])
	return args, err
}

func (m *UndelegateV2Method) PackOutput(result bool) ([]byte, error) {
	return m.Method.Outputs.Pack(result)
}

func (m *UndelegateV2Method) UnpackOutput(data []byte) (bool, error) {
	amount, err := m.Method.Outputs.Unpack(data)
	if err != nil {
		return false, err
	}
	return amount[0].(bool), nil
}

func (m *UndelegateV2Method) UnpackEvent(log *ethtypes.Log) (*fxcontract.IStakingUndelegateV2, error) {
	if log == nil {
		return nil, errors.New("empty log")
	}
	filterer, err := fxcontract.NewIStakingFilterer(common.Address{}, nil)
	if err != nil {
		return nil, err
	}
	return filterer.ParseUndelegateV2(*log)
}

func Undelegate(ctx sdk.Context, sk StakingKeeper, bk BankKeeper, delAddr sdk.AccAddress,
	valAddr sdk.ValAddress, shares sdk.Dec, bondDenom string,
) (sdkmath.Int, time.Time, error) {
	validator, found := sk.GetValidator(ctx, valAddr)
	if !found {
		return sdkmath.Int{}, time.Time{}, stakingtypes.ErrNoDelegatorForAddress
	}

	if sk.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return sdkmath.Int{}, time.Time{}, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := sk.Unbond(ctx, delAddr, valAddr, shares)
	if err != nil {
		return sdkmath.Int{}, time.Time{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		coins := sdk.NewCoins(sdk.NewCoin(bondDenom, returnAmount))
		if err := bk.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, coins); err != nil {
			return sdkmath.Int{}, time.Time{}, err
		}
	}

	completionTime := ctx.BlockHeader().Time.Add(sk.UnbondingTime(ctx))
	ubd := sk.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	sk.InsertUBDQueue(ctx, ubd, completionTime)

	return returnAmount, completionTime, nil
}

func UndelegateEmitEvents(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, amount sdkmath.Int, completionTime time.Time) {
	if amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "undelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", fxtypes.DefaultDenom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeUnbond,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
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
