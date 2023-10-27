package v6

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v6/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v6.0.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{}
	},
}
