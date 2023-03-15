package v4

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/functionx/fx-core/v3/app/upgrades"
	arbitrumtypes "github.com/functionx/fx-core/v3/x/arbitrum/types"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	optimismtypes "github.com/functionx/fx-core/v3/x/optimism/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "fxv4",
	CreateUpgradeHandler: createUpgradeHandler,
	PreUpgradeCmd:        preUpgradeCmd(),
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Added: []string{
				arbitrumtypes.ModuleName,
				optimismtypes.ModuleName,
			},
			Deleted: []string{
				gravitytypes.ModuleName,
			},
		}
	},
}

type DenomAlias struct {
	Denom string
	Alias string
}
