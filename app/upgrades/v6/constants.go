package v6

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v6/app/upgrades"
	layer2types "github.com/functionx/fx-core/v6/x/layer2/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v6.0.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Added: []string{layer2types.ModuleName},
		}
	},
}
