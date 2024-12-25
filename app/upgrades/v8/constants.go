package v8

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/pundiai/fx-core/v8/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v8.0.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Deleted: []string{
				"fxtransfer",
			},
		}
	},
}
