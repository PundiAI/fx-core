package staking

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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/evm/types"
)

func (c *Contract) Undelegate(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("undelegate method not readonly")
	}
	// parse args
	var args UndelegateArgs
	err := types.ParseMethodArgs(UndelegateMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	_, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}
	sender := sdk.AccAddress(contract.Caller().Bytes())
	evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom

	// withdraw rewards if delegation exist, add reward to evm state balance
	reward := big.NewInt(0)
	if _, found = c.stakingKeeper.GetDelegation(ctx, sender, valAddr); found {
		if reward, err = c.withdraw(ctx, evm, contract.Caller(), valAddr, evmDenom); err != nil {
			return nil, err
		}
	}

	unDelAmount, completionTime, err := Undelegate(ctx, c.stakingKeeper, c.bankKeeper, sender, valAddr, sdk.NewDecFromBigInt(args.Shares), evmDenom)
	if err != nil {
		return nil, fmt.Errorf("undelegate failed: %s", err.Error())
	}

	// add undelegate log
	if err := c.AddLog(evm, UndelegateEvent, []common.Hash{contract.Caller().Hash()},
		args.Validator, args.Shares, unDelAmount.BigInt(), big.NewInt(completionTime.Unix())); err != nil {
		return nil, err
	}
	// add undelegate event
	UndelegateEmitEvents(ctx, sender, valAddr, unDelAmount, completionTime)

	return UndelegateMethod.Outputs.Pack(unDelAmount.BigInt(), reward, big.NewInt(completionTime.Unix()))
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
