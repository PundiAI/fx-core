package v7

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/functionx/fx-core/v7/app/upgrades"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          "v7.4.x",
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: func() *storetypes.StoreUpgrades {
		if fxtypes.ChainId() == fxtypes.TestnetChainId {
			return &storetypes.StoreUpgrades{}
		}
		return &storetypes.StoreUpgrades{
			Added: []string{
				crisistypes.ModuleName,
				consensustypes.ModuleName,
			},
		}
	},
}
