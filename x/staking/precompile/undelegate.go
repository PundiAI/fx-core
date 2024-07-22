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
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
	fxstakingkeeper "github.com/functionx/fx-core/v7/x/staking/keeper"
)

func (c *Contract) Undelegate(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("undelegate method not readonly")
	}
	var args UndelegateArgs
	err := types.ParseMethodArgs(UndelegateMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		_, found := c.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return fmt.Errorf("validator not found: %s", valAddr.String())
		}
		sender := sdk.AccAddress(contract.Caller().Bytes())
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		withdrawAddr := c.distrKeeper.GetDelegatorWithdrawAddr(ctx, sender)
		beforeDelBalance := c.bankKeeper.GetBalance(ctx, withdrawAddr, evmDenom)

		unDelAmount, completionTime, err := Undelegate(ctx, c.stakingKeeper, c.bankKeeper, sender, valAddr, sdkmath.LegacyNewDecFromBigInt(args.Shares), evmDenom)
		if err != nil {
			return fmt.Errorf("undelegate failed: %s", err.Error())
		}

		afterDelBalance := c.bankKeeper.GetBalance(ctx, withdrawAddr, evmDenom)
		rewardCoin := afterDelBalance.Sub(beforeDelBalance)

		// add undelegate log
		if err = c.AddLog(evm, UndelegateEvent, []common.Hash{contract.Caller().Hash()},
			args.Validator, args.Shares, unDelAmount.BigInt(), big.NewInt(completionTime.Unix())); err != nil {
			return err
		}
		// add undelegate event
		UndelegateEmitEvents(ctx, sender, valAddr, unDelAmount, completionTime)

		result, err = UndelegateMethod.Outputs.Pack(unDelAmount.BigInt(), rewardCoin.Amount.BigInt(), big.NewInt(completionTime.Unix()))
		return err
	})

	return result, err
}

func (c *Contract) UndelegateV2(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("undelegate method not readonly")
	}
	var args UndelegateV2Args
	err := types.ParseMethodArgs(UndelegateV2Method, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(types.ExtStateDB)

	var result []byte
	err = stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
		sender := sdk.AccAddress(contract.Caller().Bytes())
		evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
		undelCoin := sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(args.Amount))

		fxStakingKeeper, ok := c.stakingKeeper.(fxstakingkeeper.Keeper)
		if !ok {
			return errortypes.ErrNotSupported
		}
		impl := stakingkeeper.NewMsgServerImpl(fxStakingKeeper.Keeper)
		undelResp, err := impl.Undelegate(sdk.WrapSDKContext(ctx), &stakingtypes.MsgUndelegate{
			DelegatorAddress: sender.String(),
			ValidatorAddress: args.Validator,
			Amount:           undelCoin,
		})
		if err != nil {
			return err
		}

		// add undelegate log
		if err = c.AddLog(evm, UndelegateV2Event, []common.Hash{contract.Caller().Hash()},
			args.Validator, args.Amount, big.NewInt(undelResp.CompletionTime.Unix())); err != nil {
			return err
		}
		result, err = UndelegateMethod.Outputs.Pack(true)
		return err
	})

	return result, err
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
