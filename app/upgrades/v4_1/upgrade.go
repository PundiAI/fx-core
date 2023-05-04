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
			testnetHandler(ctx, app)
			return fromVM, nil
		}

		// mainnet upgrade
		return mainnetHandler(ctx, mm, configurator, app, fromVM)
	}
}

func testnetHandler(ctx sdk.Context, app *keepers.AppKeepers) {
	cacheCtx, commit := ctx.CacheContext()

	// update logic code
	v4.UpdateLogicCode(cacheCtx, app.EvmKeeper)

	commit()
	ctx.Logger().Info("Upgrade complete")
}

func mainnetHandler(
	ctx sdk.Context,
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
	fromVM module.VersionMap,
) (module.VersionMap, error) {
	cacheCtx, commit := ctx.CacheContext()
	// 1. initialize the evm module account
	v4.CreateEvmModuleAccount(cacheCtx, app.AccountKeeper)

	// 2. init go fx params
	v4.InitGovFXParams(cacheCtx, app.GovKeeper)

	// 3. update logic code
	v4.UpdateLogicCode(cacheCtx, app.EvmKeeper)

	cacheCtx.Logger().Info("start to run v4 migrations...", "module", "upgrade")
	toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
	if err != nil {
		return fromVM, err
	}

	// 4. update arbitrum and optimism denom alias, after bank module migration, because bank module migrates to fixing the bank denom bug
	// discovered in https://github.com/cosmos/cosmos-sdk/pull/13821
	v4.UpdateDenomAliases(cacheCtx, app.Erc20Keeper)

	// 5. reset cross chain module oracle delegate, bind oracle delegate starting info
	err = v4.ResetCrossChainModuleOracleDelegate(cacheCtx, app.CrossChainKeepers, app.StakingKeeper, app.DistrKeeper)
	if err != nil {
		return fromVM, err
	}

	// 6. remove bsc oracles
	v4.RemoveBscOracle(cacheCtx, app.BscKeeper)

	commit()
	ctx.Logger().Info("Upgrade complete")
	return toVM, nil
}
