package v3

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v3/app/upgrades"
)

const (
	UpgradeName = "v3"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *store.StoreUpgrades {
		return &store.StoreUpgrades{
			Added: []string{},
			Deleted: []string{
				"other",
			},
		}
	},
}
