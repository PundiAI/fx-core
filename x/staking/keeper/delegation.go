package keeper

import (
	"bytes"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Unbond unbonds a particular delegation and perform associated store operations.
func (k Keeper) Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (amount sdk.Int, err error) {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return amount, types.ErrNoDelegatorForAddress
	}

	// call the before-delegation-modified hook
	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	// ensure that we have enough shares to remove
	if delegation.Shares.LT(shares) {
		return amount, errorsmod.Wrap(types.ErrNotEnoughDelegationShares, delegation.Shares.String())
	}

	// get validator
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return amount, types.ErrNoValidatorFound
	}

	// subtract shares from delegation
	delegation.Shares = delegation.Shares.Sub(shares)

	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return amount, err
	}

	isValidatorOperator := delegatorAddress.Equals(validator.GetOperator())

	// If the delegation is the operator of the validator and undelegating will decrease the validator's
	// self-delegation below their minimum, we jail the validator.
	if isValidatorOperator && !validator.Jailed &&
		validator.TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
		k.jailValidator(ctx, validator)
		validator = k.mustGetValidator(ctx, validator.GetOperator())
	}

	// remove the delegation
	if delegation.Shares.IsZero() {
		k.RemoveDelegation(ctx, delegation)
	} else {
		k.SetDelegation(ctx, delegation)
		// call the after delegation modification hook
		k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr())
	}

	// NOTE: this logic by fx-core: begin
	lpTokenContract, found := k.GetValidatorLPToken(ctx, validator.GetOperator())
	if !found {
		return amount, errortypes.ErrInvalidRequest.Wrapf("lpToken contract not found for validator")
	}

	if err = k.BurnLPToken(ctx, lpTokenContract, delAddr, shares); err != nil {
		return amount, err
	}
	// NOTE: this logic by fx-core: end

	// remove the shares and coins from the validator
	// NOTE that the amount is later (in keeper.Delegation) moved between staking module pools
	validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)

	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		k.RemoveValidator(ctx, validator.GetOperator())
	}

	return amount, nil
}

// getBeginInfo returns the completion time and height of a redelegation, along
// with a boolean signaling if the redelegation is complete based on the source
// validator.
func (k Keeper) getBeginInfo(
	ctx sdk.Context, valSrcAddr sdk.ValAddress,
) (completionTime time.Time, height int64, completeNow bool) {
	validator, found := k.GetValidator(ctx, valSrcAddr)

	// TODO: When would the validator not be found?
	switch {
	case !found || validator.IsBonded():
		// the longest wait - just unbonding period from now
		completionTime = ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
		height = ctx.BlockHeight()

		return completionTime, height, false

	case validator.IsUnbonded():
		return completionTime, height, true

	case validator.IsUnbonding():
		return validator.UnbondingTime, validator.UnbondingHeight, false

	default:
		panic(fmt.Sprintf("unknown validator status: %s", validator.Status))
	}
}

// Undelegate unbonds an amount of delegator shares from a given validator. It
// will verify that the unbonding entries between the delegator and validator
// are not exceeded and unbond the staked tokens (based on shares) by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the staking EndBlocker.
func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (time.Time, error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, types.ErrNoDelegatorForAddress
	}

	if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, types.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := k.Unbond(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		k.bondedTokensToNotBonded(ctx, returnAmount)
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// BeginRedelegation begins unbonding / redelegation and creates a redelegation
// record.
func (k Keeper) BeginRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (completionTime time.Time, err error) {
	if bytes.Equal(valSrcAddr, valDstAddr) {
		return time.Time{}, types.ErrSelfRedelegation
	}

	dstValidator, found := k.GetValidator(ctx, valDstAddr)
	if !found {
		return time.Time{}, types.ErrBadRedelegationDst
	}

	srcValidator, found := k.GetValidator(ctx, valSrcAddr)
	if !found {
		return time.Time{}, types.ErrBadRedelegationDst
	}

	// check if this is a transitive redelegation
	if k.HasReceivingRedelegation(ctx, delAddr, valSrcAddr) {
		return time.Time{}, types.ErrTransitiveRedelegation
	}

	if k.HasMaxRedelegationEntries(ctx, delAddr, valSrcAddr, valDstAddr) {
		return time.Time{}, types.ErrMaxRedelegationEntries
	}

	returnAmount, err := k.Unbond(ctx, delAddr, valSrcAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	if returnAmount.IsZero() {
		return time.Time{}, types.ErrTinyRedelegationAmount
	}

	sharesCreated, err := k.Delegate(ctx, delAddr, returnAmount, srcValidator.GetStatus(), dstValidator, false)
	if err != nil {
		return time.Time{}, err
	}

	// create the unbonding delegation
	completionTime, height, completeNow := k.getBeginInfo(ctx, valSrcAddr)

	if completeNow { // no need to create the redelegation object
		return completionTime, nil
	}

	red := k.SetRedelegationEntry(
		ctx, delAddr, valSrcAddr, valDstAddr,
		height, completionTime, returnAmount, sharesAmount, sharesCreated,
	)
	k.InsertRedelegationQueue(ctx, red, completionTime)

	return completionTime, nil
}
