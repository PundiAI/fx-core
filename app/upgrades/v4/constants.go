package v4

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v3/app/upgrades"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "fxv4",
	CreateUpgradeHandler: createUpgradeHandler,
	PreUpgradeCmd:        preUpgradeCmd(),
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Deleted: []string{
				gravitytypes.ModuleName,
			},
		}
	},
}
