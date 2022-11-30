package app

import (
	"fmt"

	"github.com/spf13/cobra"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v2 "github.com/functionx/fx-core/v2/app/upgrades/v020"
)

// AddPreUpgradeCommand add pre-upgrade command
func AddPreUpgradeCommand(rootCmd *cobra.Command) {
	// v020 pre-upgrade command
	rootCmd.AddCommand(v2.PreUpgradeCmd())
}

func (app *App) setUpgradeHandler() {
	// v2 upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(
		v2.UpgradeName,
		v2.CreateUpgradeHandler(
			app.keys,
			app.mm,
			app.configurator,
			app.BankKeeper,
			app.ParamsKeeper,
			app.IBCKeeper,
			app.TransferKeeper,
			app.Erc20Keeper,
		),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}
	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	switch upgradeInfo.Name {
	case v2.UpgradeName:
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, v2.GetStoreUpgrades()))
	}
}
