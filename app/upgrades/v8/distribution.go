package v8

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	fxstakingkeeper "github.com/pundiai/fx-core/v8/x/staking/keeper"
)

func migrateDistribution(ctx sdk.Context, stakingKeeper *fxstakingkeeper.Keeper, distrKeeper distributionkeeper.Keeper) error {
	if err := migrateValidatorAccumulatedCommission(ctx, distrKeeper); err != nil {
		return err
	}
	if err := migrateValidatorOutstandingRewards(ctx, distrKeeper); err != nil {
		return err
	}
	if err := migrateValidatorCurrentRewards(ctx, distrKeeper); err != nil {
		return err
	}
	if err := migrateDelegatorStartingInfos(ctx, stakingKeeper, distrKeeper); err != nil {
		return err
	}
	if err := migrateValidatorHistoricalRewards(ctx, distrKeeper); err != nil {
		return err
	}
	return migrateFeePool(ctx, distrKeeper)
}

func migrateValidatorHistoricalRewards(ctx sdk.Context, keeper distributionkeeper.Keeper) error {
	var err error
	keeper.IterateValidatorHistoricalRewards(ctx, func(valAddr sdk.ValAddress, period uint64, rewards types.ValidatorHistoricalRewards) bool {
		newRewards := make(sdk.DecCoins, 0, len(rewards.CumulativeRewardRatio))
		for _, coin := range rewards.CumulativeRewardRatio {
			newDenom := coin.Denom
			if newDenom == fxtypes.LegacyFXDenom {
				newDenom = fxtypes.DefaultDenom
			}
			newRewards = append(newRewards, sdk.NewDecCoinFromDec(newDenom, coin.Amount))
		}
		rewards.CumulativeRewardRatio = newRewards
		err = keeper.SetValidatorHistoricalRewards(ctx, valAddr, period, rewards)
		return err != nil
	})
	return err
}

func migrateDelegatorStartingInfos(ctx sdk.Context, stakingKeeper *fxstakingkeeper.Keeper, distrKeeper distributionkeeper.Keeper) error {
	var err error
	validatorMap := make(map[string]stakingtypes.Validator)
	distrKeeper.IterateDelegatorStartingInfos(ctx, func(valAddr sdk.ValAddress, addr sdk.AccAddress, info types.DelegatorStartingInfo) bool {
		var delegation stakingtypes.Delegation
		delegation, err = stakingKeeper.GetDelegation(ctx, addr, valAddr)
		if err != nil {
			return true
		}
		validator, ok := validatorMap[delegation.ValidatorAddress]
		if !ok {
			validator, err = stakingKeeper.GetValidator(ctx, valAddr)
			if err != nil {
				return true
			}
			validatorMap[delegation.ValidatorAddress] = validator
		}
		info.Stake = validator.TokensFromSharesTruncated(delegation.Shares)
		err = distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, addr, info)
		return err != nil
	})
	return err
}

func migrateFeePool(ctx sdk.Context, distrKeeper distributionkeeper.Keeper) error {
	feePool, err := distrKeeper.FeePool.Get(ctx)
	if err != nil {
		return err
	}

	feePool.CommunityPool = fxtypes.SwapDecCoins(feePool.CommunityPool)
	return distrKeeper.FeePool.Set(ctx, feePool)
}

func migrateValidatorAccumulatedCommission(ctx sdk.Context, distrKeeper distributionkeeper.Keeper) error {
	var err error
	distrKeeper.IterateValidatorAccumulatedCommissions(ctx, func(addr sdk.ValAddress, commission types.ValidatorAccumulatedCommission) bool {
		newCommission := fxtypes.SwapDecCoins(commission.Commission)
		commission.Commission = newCommission
		err = distrKeeper.SetValidatorAccumulatedCommission(ctx, addr, commission)
		return err != nil
	})

	return err
}

func migrateValidatorCurrentRewards(ctx sdk.Context, distrKeeper distributionkeeper.Keeper) error {
	var err error
	distrKeeper.IterateValidatorCurrentRewards(ctx, func(addr sdk.ValAddress, rewards types.ValidatorCurrentRewards) bool {
		newRewards := fxtypes.SwapDecCoins(rewards.Rewards)
		rewards.Rewards = newRewards
		err = distrKeeper.SetValidatorCurrentRewards(ctx, addr, rewards)
		return err != nil
	})
	return err
}

func migrateValidatorOutstandingRewards(ctx sdk.Context, distrKeeper distributionkeeper.Keeper) error {
	var err error
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(addr sdk.ValAddress, rewards types.ValidatorOutstandingRewards) bool {
		newRewards := fxtypes.SwapDecCoins(rewards.Rewards)
		rewards.Rewards = newRewards
		err = distrKeeper.SetValidatorOutstandingRewards(ctx, addr, rewards)
		return err != nil
	})

	return err
}

func CheckDistributionModule(t *testing.T, ctx sdk.Context, distrKeeper distributionkeeper.Keeper) {
	t.Helper()

	denomCheckFn := func(decCoins sdk.DecCoins, msg ...string) {
		if decCoins.IsZero() {
			return
		}
		require.Truef(t, decCoins.AmountOf(fxtypes.LegacyFXDenom).IsZero(), decCoins.String(), msg)
		require.Falsef(t, decCoins.AmountOf(fxtypes.DefaultDenom).IsZero(), decCoins.String(), msg)
	}

	// check fee pool
	feePool, err := distrKeeper.FeePool.Get(ctx)
	require.NoError(t, err)
	denomCheckFn(feePool.GetCommunityPool())

	// check validator commission
	distrKeeper.IterateValidatorAccumulatedCommissions(ctx, func(val sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool) {
		denomCheckFn(commission.GetCommission(), fmt.Sprintf("val:%s,commission:%s", val.String(), commission.String()))
		return false
	})

	// check validator outstanding rewards
	distrKeeper.IterateValidatorOutstandingRewards(ctx, func(val sdk.ValAddress, rewards types.ValidatorOutstandingRewards) (stop bool) {
		denomCheckFn(rewards.GetRewards(), fmt.Sprintf("val:%s,outstanding rewards:%s", val.String(), rewards.String()))
		return false
	})

	// check validator current rewards
	distrKeeper.IterateValidatorCurrentRewards(ctx, func(val sdk.ValAddress, rewards types.ValidatorCurrentRewards) (stop bool) {
		denomCheckFn(rewards.GetRewards(), fmt.Sprintf("val:%s,current rewards:%s", val.String(), rewards.String()))
		return false
	})

	// check validator historical rewards
	distrKeeper.IterateValidatorHistoricalRewards(ctx, func(val sdk.ValAddress, period uint64, rewards types.ValidatorHistoricalRewards) (stop bool) {
		denomCheckFn(rewards.GetCumulativeRewardRatio(), fmt.Sprintf("val:%s,historical rewards:%s", val.String(), rewards.String()))
		return false
	})
}
