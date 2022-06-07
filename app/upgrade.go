package app

import (
	"fmt"

	store "github.com/cosmos/cosmos-sdk/store/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	upgradev2 "github.com/functionx/fx-core/app/upgrades/v2"
)

func (app *App) setUpgradeHandler() {
	// set upgrade handler v2
	app.UpgradeKeeper.SetUpgradeHandler(
		upgradev2.UpgradeName, upgradev2.CreateUpgradeHandler(app.mm, app.configurator,
			app.GetKey(banktypes.StoreKey), app.BankKeeper, app.IBCKeeper, app.Erc20Keeper,
		),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}
	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
	var storeUpgrades *store.StoreUpgrades
	switch upgradeInfo.Name {
	case upgradev2.UpgradeName:
		storeUpgrades = upgradev2.GetStoreUpgrades()
	}
	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
