package v7

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v7/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v7.1.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{}
	},
}
