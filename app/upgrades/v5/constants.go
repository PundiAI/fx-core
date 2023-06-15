package v5

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v5/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v5.0.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{}
	},
}
