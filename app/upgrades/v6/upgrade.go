package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v6/app/keepers"
	govtypes "github.com/functionx/fx-core/v6/x/gov/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		if err := UpdateParams(cacheCtx, app); err != nil {
			return nil, err
		}

		ctx.Logger().Info("start to run v6 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		commit()
		ctx.Logger().Info("Upgrade complete")
		return toVM, nil
	}
}

func UpdateParams(cacheCtx sdk.Context, app *keepers.AppKeepers) error {
	mintParams := app.MintKeeper.GetParams(cacheCtx)
	mintParams.InflationMax = sdk.ZeroDec()
	mintParams.InflationMin = sdk.ZeroDec()
	if err := mintParams.Validate(); err != nil {
		return err
	}
	app.MintKeeper.SetParams(cacheCtx, mintParams)

	distrParams := app.DistrKeeper.GetParams(cacheCtx)
	distrParams.CommunityTax = sdk.ZeroDec()
	distrParams.BaseProposerReward = sdk.ZeroDec()
	distrParams.BonusProposerReward = sdk.ZeroDec()
	if err := distrParams.ValidateBasic(); err != nil {
		return err
	}
	app.DistrKeeper.SetParams(cacheCtx, distrParams)

	stakingParams := app.StakingKeeper.GetParams(cacheCtx)
	stakingParams.UnbondingTime = 1
	if err := stakingParams.Validate(); err != nil {
		return err
	}
	app.StakingKeeper.SetParams(cacheCtx, stakingParams)

	govTallyParams := app.GovKeeper.GetTallyParams(cacheCtx)
	govTallyParams.Quorum = sdk.OneDec().String()        // 100%
	govTallyParams.Threshold = sdk.OneDec().String()     // 100%
	govTallyParams.VetoThreshold = sdk.OneDec().String() // 100%
	app.GovKeeper.SetTallyParams(cacheCtx, govTallyParams)

	app.GovKeeper.IterateParams(cacheCtx, func(param *govtypes.Params) (stop bool) {
		param.Quorum = sdk.OneDec().String()        // 100%
		param.Threshold = sdk.OneDec().String()     // 100%
		param.VetoThreshold = sdk.OneDec().String() // 100%
		if err := param.ValidateBasic(); err != nil {
			panic(err)
		}
		if err := app.GovKeeper.SetParams(cacheCtx, param); err != nil {
			panic(err)
		}
		return false
	})
	return nil
}
