package app

import (
	"fmt"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v4/app/upgrades/v4_1"
)

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	plan := v4_1.Upgrade()
	if upgradeInfo.Name == plan.UpgradeName {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, plan.StoreUpgrades()))
	}
}

func (app *App) setupUpgradeHandlers() {
	plan := v4_1.Upgrade()
	app.UpgradeKeeper.SetUpgradeHandler(
		plan.UpgradeName,
		plan.CreateUpgradeHandler(
			app.mm,
			app.configurator,
			app.AppKeepers,
		),
	)
}
