package v7

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/functionx/fx-core/v7/app/upgrades"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v7.3.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		return &storetypes.StoreUpgrades{
			Added: []string{
				crisistypes.ModuleName,
				consensustypes.ModuleName,
			},
		}
	},
}
