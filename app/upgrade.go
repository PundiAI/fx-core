package app

import (
	"fmt"

	store "github.com/cosmos/cosmos-sdk/store/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	upgradev2 "github.com/functionx/fx-core/app/upgrades/v2"
)

func (myApp *App) setUpgradeHandler() {
	// set upgrade handler v2
	myApp.UpgradeKeeper.SetUpgradeHandler(
		upgradev2.UpgradeName, upgradev2.CreateUpgradeHandler(myApp.mm, myApp.configurator,
			myApp.GetKey(banktypes.StoreKey), myApp.BankKeeper, myApp.IBCKeeper, myApp.Erc20Keeper,
		),
	)

	upgradeInfo, err := myApp.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}
	if myApp.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
	var storeUpgrades *store.StoreUpgrades
	switch upgradeInfo.Name {
	case upgradev2.UpgradeName:
		storeUpgrades = upgradev2.GetStoreUpgrades()
	}
	if storeUpgrades != nil {
		myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
