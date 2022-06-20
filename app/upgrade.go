package app

import (
	"fmt"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	upgrade "github.com/functionx/fx-core/app/upgrades/v020"
)

func (app *App) setUpgradeHandler() {
	// set upgrade handler v2
	app.UpgradeKeeper.SetUpgradeHandler(
		upgrade.UpgradeName, upgrade.CreateUpgradeHandler(app.keys, app.mm, app.configurator,
			app.BankKeeper, app.AccountKeeper, app.ParamsKeeper, app.IBCKeeper, app.TransferKeeper, app.Erc20Keeper),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}
	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
	switch upgradeInfo.Name {
	case upgrade.UpgradeName:
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgrade.GetStoreUpgrades()))
	}
}
