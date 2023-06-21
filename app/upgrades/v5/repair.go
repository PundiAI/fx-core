package v5

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v5/types"
	fxstakingkeeper "github.com/functionx/fx-core/v5/x/staking/keeper"
)

func RepairSlashPeriod(ctx sdk.Context, sk fxstakingkeeper.Keeper, dk distrkeeper.Keeper) {
	slashPeriod := delegatorNotSlash(ctx, sk, dk)
	vals, newSlashPeriods := addSlashPeriodTestnetFXV4(ctx, slashPeriod)
	for _, val := range vals {
		valAddr, err := sdk.ValAddressFromBech32(val)
		if err != nil {
			panic(err)
		}
		fixSlashPeriodTestnetFXV4(ctx, dk, valAddr, newSlashPeriods[val])
	}
}

func fixSlashPeriodTestnetFXV4(ctx sdk.Context, dk distrkeeper.Keeper, val sdk.ValAddress, periods []SlashPeriod) {
	logger := ctx.Logger()

	currentRewards := dk.GetValidatorCurrentRewards(ctx, val)
	currentHistoricalRewards := dk.GetValidatorHistoricalRewards(ctx, val, currentRewards.Period)

	periodHistoricalRewards := make(map[uint64]distrtypes.ValidatorHistoricalRewards, len(periods))
	for _, p := range periods {
		if p.Delegator == nil {
			continue
		}
		// get period starting info
		startingInfo := dk.GetDelegatorStartingInfo(ctx, val, p.Delegator)
		historicalRewards := dk.GetValidatorHistoricalRewards(ctx, val, startingInfo.PreviousPeriod)
		periodHistoricalRewards[p.Period] = historicalRewards
		dk.DeleteValidatorHistoricalReward(ctx, val, startingInfo.PreviousPeriod)
	}

	slashBefore := false
	for idx, p := range periods {
		if p.Delegator != nil {
			// set new period historical rewards
			dk.SetValidatorHistoricalRewards(ctx, val, p.Period, periodHistoricalRewards[p.Period])
		}

		if p.Height < fxtypes.TestnetBlockHeightV4 ||
			p.Delegator != nil && !slashBefore {
			continue
		}

		if p.Delegator == nil {
			referenceCount := uint32(1)
			if idx == len(periods)-1 {
				referenceCount += 1
			}
			logger.Info("add slash period", "validator", val, "height", p.Height, "period", p.Period, "referenceCount", referenceCount)
			lastHistoricalRewards := periodHistoricalRewards[periods[idx-1].Period] // todo
			historicalRewards := distrtypes.NewValidatorHistoricalRewards(lastHistoricalRewards.CumulativeRewardRatio, referenceCount)
			dk.SetValidatorHistoricalRewards(ctx, val, p.Period, historicalRewards)
			periodHistoricalRewards[p.Period] = historicalRewards

			// add slash period
			fraction, _ := sdk.NewDecFromStr("0.001") // todo
			slashEvent := distrtypes.NewValidatorSlashEvent(p.Period, fraction)
			dk.SetValidatorSlashEvent(ctx, val, p.Height, p.Period, slashEvent)
			slashBefore = true
		} else {
			logger.Info("migrate slash period", "validator", val, "height", p.Height, "period", p.Period, "delegator", p.Delegator)
			// set new starting info
			startingInfo := dk.GetDelegatorStartingInfo(ctx, val, p.Delegator)
			startingInfo.PreviousPeriod = p.Period
			dk.SetDelegatorStartingInfo(ctx, val, p.Delegator, startingInfo)
		}

		if idx == len(periods)-1 {
			currentRewards.Period = p.Period + 1
			dk.SetValidatorCurrentRewards(ctx, val, currentRewards)
			dk.SetValidatorHistoricalRewards(ctx, val, currentRewards.Period, currentHistoricalRewards)
		}
	}
}

func addSlashPeriodTestnetFXV4(ctx sdk.Context, slashPeriod map[string][]SlashPeriod) ([]string, map[string][]SlashPeriod) {
	vals := make([]string, 0, len(slashPeriod))
	newSlashPeriod := make(map[string][]SlashPeriod, len(slashPeriod))
	for val, periods := range slashPeriod {
		if len(periods) == 0 {
			ctx.Logger().Info("skip delegation empty", "address", val)
			continue
		}
		heights, ok := ValidatorSlashHeightTestnetFXV4[val]
		if !ok {
			ctx.Logger().Info("skip validator not found", "address", val)
			continue
		}

		for _, height := range heights {
			periods = append(periods, SlashPeriod{
				Delegator: nil,
				Height:    uint64(height),
				Period:    0,
			})
		}
		sort.SliceStable(periods, func(i, j int) bool {
			return periods[i].Height < periods[j].Height
		})

		// fix period with slash height
		lastPeriod := uint64(0)
		newPeriods := make([]SlashPeriod, 0, len(periods))
		for _, del := range periods {
			if del.Delegator == nil && del.Period == 0 || // slash height
				del.Period <= lastPeriod { // after slash height
				del.Period = lastPeriod + 1
			}
			lastPeriod = del.Period
			newPeriods = append(newPeriods, del)
		}

		newSlashPeriod[val] = newPeriods
		vals = append(vals, val)
	}

	// sort by validator address
	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})

	return vals, newSlashPeriod
}

func delegatorNotSlash(ctx sdk.Context, sk fxstakingkeeper.Keeper, dk distrkeeper.Keeper) map[string][]SlashPeriod {
	valDels := validatorNotSlash(ctx, sk, dk)
	sk.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) (stop bool) {
		if _, ok := valDels[del.ValidatorAddress]; !ok {
			return false
		}
		startingInfo := dk.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
		valDels[del.ValidatorAddress] = append(valDels[del.ValidatorAddress], SlashPeriod{
			Delegator: del.GetDelegatorAddr(),
			Height:    startingInfo.Height,
			Period:    startingInfo.PreviousPeriod,
		})
		return false
	})
	for val, dels := range valDels {
		sort.SliceStable(dels, func(i, j int) bool {
			return dels[i].Period < dels[j].Period
		})
		valDels[val] = dels
	}
	return valDels
}

func validatorNotSlash(ctx sdk.Context, sk fxstakingkeeper.Keeper, dk distrkeeper.Keeper) map[string][]SlashPeriod {
	vals := make(map[string][]SlashPeriod, 50)
	sk.IterateAllDelegations(ctx, func(del stakingtypes.Delegation) (stop bool) {
		if _, ok := vals[del.GetValidatorAddr().String()]; ok {
			return false
		}
		val, found := sk.GetValidator(ctx, del.GetValidatorAddr())
		if !found {
			ctx.Logger().Error("validator not found", "validator", del.GetValidatorAddr().String())
			// if validator not found, skip
			return false
		}
		currentStake := val.TokensFromShares(del.GetShares())
		stake := calculateSlashStake(ctx, dk, del)

		marginOfErr := sdk.SmallestDec().MulInt64(3)
		if !stake.LTE(currentStake.Add(marginOfErr)) {
			// if delegate after slash, no error
			vals[del.GetValidatorAddr().String()] = make([]SlashPeriod, 0, 30)
		}
		return false
	})
	return vals
}

func calculateSlashStake(ctx sdk.Context, dk distrkeeper.Keeper, del stakingtypes.Delegation) sdk.Dec {
	startingInfo := dk.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake

	startingHeight := startingInfo.Height
	endingHeight := uint64(ctx.BlockHeight())

	if endingHeight > startingHeight {
		dk.IterateValidatorSlashEventsBetween(ctx, del.GetValidatorAddr(), startingHeight, endingHeight,
			func(height uint64, event distrtypes.ValidatorSlashEvent) (stop bool) {
				endingPeriod := event.ValidatorPeriod
				if endingPeriod > startingPeriod {
					stake = stake.MulTruncate(sdk.OneDec().Sub(event.Fraction))
					startingPeriod = endingPeriod
				}
				return false
			},
		)
	}
	return stake
}
