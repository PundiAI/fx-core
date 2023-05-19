package v4_2

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v4/app/upgrades"
	v4 "github.com/functionx/fx-core/v4/app/upgrades/v4"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

func Upgrade() upgrades.Upgrade {
	upgrade := upgrades.Upgrade{
		UpgradeName:          "v4.2.x",
		CreateUpgradeHandler: createUpgradeHandler,
		StoreUpgrades:        v4.Upgrade.StoreUpgrades,
	}

	// if testnet, store has been upgraded in fxv4
	if fxtypes.ChainId() == fxtypes.TestnetChainId {
		upgrade.StoreUpgrades = func() *storetypes.StoreUpgrades {
			return &storetypes.StoreUpgrades{}
		}
	}

	return upgrade
}
