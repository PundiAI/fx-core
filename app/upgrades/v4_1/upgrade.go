package v4_1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v4/app/keepers"
	v4 "github.com/functionx/fx-core/v4/app/upgrades/v4"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

func createUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// testnet upgrade
		if ctx.ChainID() == fxtypes.TestnetChainId {
			cacheCtx, commit := ctx.CacheContext()

			// update logic code
			v4.UpdateLogicCode(cacheCtx, app.EvmKeeper)

			commit()
			ctx.Logger().Info("Upgrade complete")
			return fromVM, nil
		}

		// mainnet upgrade
		return v4.CreateUpgradeHandler(mm, configurator, app)(ctx, plan, fromVM)
	}
}
