package v8

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func migrateStakingModule(ctx sdk.Context, keeper *stakingkeeper.Keeper) error {
	params, err := keeper.GetParams(ctx)
	if err != nil {
		return err
	}
	params.BondDenom = fxtypes.DefaultDenom
	if err = keeper.SetParams(ctx, params); err != nil {
		return err
	}
	validators, err := keeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	for _, validator := range validators {
		validator.Tokens = fxtypes.SwapAmount(validator.Tokens)
		validator.MinSelfDelegation = fxtypes.SwapAmount(validator.MinSelfDelegation)
		if err = keeper.SetValidator(ctx, validator); err != nil {
			return err
		}
		err = keeper.SetValidatorByPowerIndex(ctx, validator)
		if err != nil {
			return err
		}
	}

	iterErr := keeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := 0; i < len(red.Entries); i++ {
			red.Entries[i].InitialBalance = fxtypes.SwapAmount(red.Entries[i].InitialBalance)
		}
		if err = keeper.SetRedelegation(ctx, red); err != nil {
			return true
		}
		return false
	})
	if iterErr != nil {
		return iterErr
	}
	if err != nil {
		return err
	}

	iterErr = keeper.IterateUnbondingDelegations(ctx, func(index int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := 0; i < len(ubd.Entries); i++ {
			ubd.Entries[i].Balance = fxtypes.SwapAmount(ubd.Entries[i].Balance)
			ubd.Entries[i].InitialBalance = fxtypes.SwapAmount(ubd.Entries[i].InitialBalance)
		}
		if err = keeper.SetUnbondingDelegation(ctx, ubd); err != nil {
			return true
		}
		return false
	})
	if iterErr != nil {
		return iterErr
	}
	return err
}

func CheckStakingModule(t *testing.T, ctx sdk.Context, keeper *stakingkeeper.Keeper) {
	t.Helper()

	params, err := keeper.GetParams(ctx)
	require.NoError(t, err)
	assert.Equal(t, fxtypes.DefaultDenom, params.BondDenom)

	delegations, err := keeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	delegationByValidator := make(map[string][]stakingtypes.Delegation)
	for i := 0; i < len(delegations); i++ {
		delegationByValidator[delegations[i].ValidatorAddress] = append(delegationByValidator[delegations[i].ValidatorAddress], delegations[i])
	}
	err = keeper.IterateValidators(ctx, func(_ int64, validator stakingtypes.ValidatorI) (stop bool) {
		delegations = delegationByValidator[validator.GetOperator()]
		totalShare := sdkmath.LegacyZeroDec()
		for i := 0; i < len(delegations); i++ {
			totalShare = totalShare.Add(delegations[i].Shares)
		}
		assert.Equal(t, totalShare, validator.GetDelegatorShares())
		return false
	})
	require.NoError(t, err)
}
