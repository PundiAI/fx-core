package app

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	upgradev2 "github.com/functionx/fx-core/app/upgrades/v2"
)

func (app *App) setUpgradeHandler() {
	// set upgrade handler v2
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradev2.UpgradeName, upgradev2.CreateUpgradeHandler(app.mm, app.configurator,
			app.GetKey(banktypes.StoreKey), app.BankKeeper, app.AccountKeeper,
			app.ParamsKeeper, app.IBCKeeper, app.Erc20Keeper),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}
	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
	switch upgradeInfo.Name {
	case upgradev2.UpgradeName:
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, upgradev2.GetStoreUpgrades()))
	}
}
