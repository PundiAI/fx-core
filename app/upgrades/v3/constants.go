package v3

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v3/app/upgrades"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
)

const (
	UpgradeName = "v3"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *store.StoreUpgrades {
		return &store.StoreUpgrades{
			Added: []string{
				avalanchetypes.ModuleName,
				ethtypes.ModuleName,
			},
			Deleted: []string{
				"other",
			},
		}
	},
}
