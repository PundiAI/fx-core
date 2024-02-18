package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v7/x/evm/types"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
)

func (c *Contract) TransferShares(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("transfer method not readonly")
	}
	// parse args
	var args TransferSharesArgs
	if err := types.ParseMethodArgs(TransferSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	token, reward, err := c.handlerTransferShares(ctx, evm, valAddr, contract.Caller(), args.To, args.Shares)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(contract.Caller().Bytes()).String()),
		),
	})
	return TransferSharesMethod.Outputs.Pack(token, reward)
}

func (c *Contract) TransferFromShares(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("transferFrom method not readonly")
	}
	// parse args
	var args TransferFromSharesArgs
	if err := types.ParseMethodArgs(TransferFromSharesMethod, &args, contract.Input[4:]); err != nil {
		return nil, err
	}

	valAddr := args.GetValidator()
	spender := contract.Caller()
	if err := c.decrementAllowance(ctx, valAddr, args.From.Bytes(), spender.Bytes(), args.Shares); err != nil {
		return nil, err
	}
	token, reward, err := c.handlerTransferShares(ctx, evm, valAddr, args.From, args.To, args.Shares)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evmtypes.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(contract.Caller().Bytes()).String()),
		),
	})
	return TransferFromSharesMethod.Outputs.Pack(token, reward)
}

func (c *Contract) decrementAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, decrease *big.Int) error {
	allowance := c.stakingKeeper.GetAllowance(ctx, valAddr, owner, spender)
	if allowance.Cmp(decrease) < 0 {
		return fmt.Errorf("transfer shares exceeds allowance(%s < %s)", allowance.String(), decrease.String())
	}
	newAllowance := big.NewInt(0).Sub(allowance, decrease)
	c.stakingKeeper.SetAllowance(ctx, valAddr, owner, spender, newAllowance)
	return nil
}

func (c *Contract) handlerTransferShares(
	ctx sdk.Context,
	evm *vm.EVM,
	valAddr sdk.ValAddress,
	from, to common.Address,
	sharesInt *big.Int,
) (*big.Int, *big.Int, error) {
	validator, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}
	fromDel, found := c.stakingKeeper.GetDelegation(ctx, from.Bytes(), valAddr)
	if !found {
		return nil, nil, errors.New("from delegation not found")
	}
	// if from has receiving redelegation, can't transfer shares
	if c.stakingKeeper.HasReceivingRedelegation(ctx, from.Bytes(), valAddr) {
		return nil, nil, errors.New("from has receiving redelegation")
	}

	shares := sdk.NewDecFromBigInt(sharesInt) // TODO share with sdk.Precision
	if fromDel.GetShares().LT(shares) {
		return nil, nil, fmt.Errorf("insufficient shares(%s < %s)", fromDel.GetShares().TruncateInt().String(), shares.TruncateInt().String())
	}

	// withdraw reward
	evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom
	_, err := c.withdraw(ctx, evm, from, valAddr, evmDenom)
	if err != nil {
		return nil, nil, fmt.Errorf("withdraw failed: %s", err.Error())
	}

	// get to delegation
	toReward := big.NewInt(0)
	toDel, toDelFound := c.stakingKeeper.GetDelegation(ctx, to.Bytes(), valAddr)
	if !toDelFound {
		toDel = stakingtypes.NewDelegation(to.Bytes(), valAddr, sdk.ZeroDec())
		// if address to not delegate before, increase validator period
		_ = c.distrKeeper.IncrementValidatorPeriod(ctx, validator)
	} else {
		// withdraw to address reward
		toReward, err = c.withdraw(ctx, evm, to, valAddr, evmDenom)
		if err != nil {
			return nil, nil, fmt.Errorf("withdraw failed: %s", err.Error())
		}
	}

	// update from delegate, delete it if shares zero
	fromDelStartingInfo := c.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, from.Bytes())
	fromDel.Shares = fromDel.Shares.Sub(shares)
	if fromDel.GetShares().IsZero() {
		// if shares zero, remove delegation and delete starting info and reference count
		if err := c.stakingKeeper.RemoveDelegation(ctx, fromDel); err != nil {
			return nil, nil, err
		}
		// decrement previous period reference count
		decrementReferenceCount(c.distrKeeper, ctx, valAddr, fromDelStartingInfo.PreviousPeriod)
		c.distrKeeper.DeleteDelegatorStartingInfo(ctx, valAddr, from.Bytes())
	} else {
		c.stakingKeeper.SetDelegation(ctx, fromDel)
		// update from starting info
		fromDelStartingInfo.Stake = validator.TokensFromSharesTruncated(fromDel.GetShares())
		c.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, from.Bytes(), fromDelStartingInfo)
	}

	// update to delegate, set starting info if to not delegate before
	toDel.Shares = toDel.Shares.Add(shares)
	c.stakingKeeper.SetDelegation(ctx, toDel)
	if !toDelFound {
		// if to not delegate before, last period reference count - 1 and set starting info
		previousPeriod := c.distrKeeper.GetValidatorCurrentRewards(ctx, valAddr).Period - 1
		incrementReferenceCount(c.distrKeeper, ctx, valAddr, previousPeriod)

		stakeToken := validator.TokensFromSharesTruncated(shares)
		toDelStartingInfo := distrtypes.NewDelegatorStartingInfo(previousPeriod, stakeToken, uint64(ctx.BlockHeight()))
		c.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo)
	} else {
		// update to starting info
		toDelStartingInfo := c.distrKeeper.GetDelegatorStartingInfo(ctx, valAddr, to.Bytes())
		toDelStartingInfo.Stake = validator.TokensFromSharesTruncated(toDel.GetShares())
		c.distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, to.Bytes(), toDelStartingInfo)
	}

	// calculate token from shares
	token := validator.TokensFromShares(shares).TruncateInt()

	// add log
	if err := c.AddLog(evm, TransferSharesEvent, []common.Hash{from.Hash(), to.Hash()},
		valAddr.String(), shares.TruncateInt().BigInt(), token.BigInt()); err != nil {
		return nil, nil, err
	}

	// add emit event
	TransferSharesEmitEvents(ctx, valAddr, from.Bytes(), to.Bytes(), sdkmath.NewIntFromBigInt(sharesInt), token)

	return token.BigInt(), toReward, nil
}

// increment the reference count for a historical rewards value
func incrementReferenceCount(k DistrKeeper, ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func decrementReferenceCount(k DistrKeeper, ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}

func TransferSharesEmitEvents(ctx sdk.Context, validator sdk.ValAddress, from, recipient sdk.AccAddress, shares, token sdkmath.Int) {
	if shares.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, evmtypes.ModuleName, "transfer_shares")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", evmtypes.TypeMsgEthereumTx},
				float32(shares.Int64()),
				[]metrics.Label{telemetry.NewLabel("validator", validator.String())},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			fxstakingtypes.EventTypeTransferShares,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, validator.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyFrom, from.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyRecipient, recipient.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyShares, shares.String()),
			sdk.NewAttribute(fxstakingtypes.AttributeKeyAmount, token.String()),
		),
	})
}
