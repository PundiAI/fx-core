package v4_1

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v5/app/upgrades"
	v4 "github.com/functionx/fx-core/v5/app/upgrades/v4"
	fxtypes "github.com/functionx/fx-core/v5/types"
)

// Deprecated: Please use v4.2.x
func Upgrade() upgrades.Upgrade {
	upgrade := upgrades.Upgrade{
		UpgradeName:          "v4.1.x",
		CreateUpgradeHandler: createUpgradeHandler,
		StoreUpgrades:        v4.Upgrade.StoreUpgrades,
	}

	// if testnet, store has been upgraded in v4
	if fxtypes.ChainId() == fxtypes.TestnetChainId {
		upgrade.StoreUpgrades = func() *storetypes.StoreUpgrades {
			return &storetypes.StoreUpgrades{}
		}
	}

	return upgrade
}
