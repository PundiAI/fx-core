package app

import (
	"fmt"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v4 "github.com/functionx/fx-core/v4/app/upgrades/v4"
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
	if upgradeInfo.Name == v4.Upgrade.UpgradeName {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, v4.Upgrade.StoreUpgrades()))
	}
}

func (app *App) setupUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v4.Upgrade.UpgradeName,
		v4.Upgrade.CreateUpgradeHandler(
			app.mm,
			app.configurator,
			app.AppKeepers,
		),
	)
}
