package app

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/functionx/fx-core/v8/app/upgrades"
	nextversion "github.com/functionx/fx-core/v8/app/upgrades/v8"
)

func (app *App) GetUpgrade() upgrades.Upgrade {
	return nextversion.Upgrade
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	plan := app.GetUpgrade()
	if upgradeInfo.Name == plan.UpgradeName {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, plan.StoreUpgrades()))
	}
}

func (app *App) setupUpgradeHandlers() {
	plan := app.GetUpgrade()
	app.UpgradeKeeper.SetUpgradeHandler(
		plan.UpgradeName,
		plan.CreateUpgradeHandler(
			app.mm,
			app.configurator,
			&app.AppKeepers,
		),
	)
}
