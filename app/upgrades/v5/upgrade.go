package v5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v5/app/keepers"
	fxtypes "github.com/functionx/fx-core/v5/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		if ctx.ChainID() == fxtypes.TestnetChainId {
			RepairSlashPeriod(ctx, app.StakingKeeper, app.DistrKeeper)
		}

		ctx.Logger().Info("start to run v5 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		commit()
		ctx.Logger().Info("Upgrade complete")
		return toVM, nil
	}
}
