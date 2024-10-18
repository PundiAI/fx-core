package v8

import (
	"context"
	"encoding/hex"
	"strings"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/v8/app/keepers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
	erc20v8 "github.com/functionx/fx-core/v8/x/erc20/migrations/v8"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
	"github.com/functionx/fx-core/v8/x/gov/keeper"
	fxgovv8 "github.com/functionx/fx-core/v8/x/gov/migrations/v8"
	fxstakingv8 "github.com/functionx/fx-core/v8/x/staking/migrations/v8"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()

		cacheCtx.Logger().Info("start to run migrations...", "module", "upgrade", "plan", plan.Name)
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		removeStoreKeys(cacheCtx, app.GetKey(stakingtypes.StoreKey), fxstakingv8.GetRemovedStoreKeys())

		if err = migrationGovCustomParam(cacheCtx, app.GovKeeper, app.GetKey(govtypes.StoreKey)); err != nil {
			return fromVM, err
		}

		if err = updateArbitrumParams(cacheCtx, app.ArbitrumKeeper); err != nil {
			return fromVM, err
		}

		updateMetadata(cacheCtx, app.BankKeeper)

		removeStoreKeys(cacheCtx, app.GetKey(erc20types.StoreKey), erc20v8.GetRemovedStoreKeys())

		commit()
		cacheCtx.Logger().Info("upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func updateArbitrumParams(ctx sdk.Context, keeper crosschainkeeper.Keeper) error {
	params := keeper.GetParams(ctx)
	params.AverageExternalBlockTime = 250
	return keeper.SetParams(ctx, &params)
}

func migrationGovCustomParam(ctx sdk.Context, keeper *keeper.Keeper, storeKey *storetypes.KVStoreKey) error {
	// 1. delete fxParams key
	removeStoreKeys(ctx, storeKey, fxgovv8.GetRemovedStoreKeys())

	// 2. init custom params
	return keeper.InitCustomParams(ctx)
}

func removeStoreKeys(ctx sdk.Context, storeKey *storetypes.KVStoreKey, keys [][]byte) {
	store := ctx.KVStore(storeKey)

	deleteFn := func(key []byte) {
		iterator := storetypes.KVStorePrefixIterator(store, key)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			store.Delete(iterator.Key())
			ctx.Logger().Info("remove store key", "kvStore", storeKey.Name(),
				"prefix", hex.EncodeToString(key), "key", string(iterator.Key()))
		}
	}

	for _, key := range keys {
		deleteFn(key)
	}
}

func updateMetadata(ctx sdk.Context, bankKeeper bankkeeper.Keeper) {
	mds := bankKeeper.GetAllDenomMetaData(ctx)
	for _, md := range mds {
		if md.Base == fxtypes.DefaultDenom || len(md.DenomUnits) == 0 || len(md.DenomUnits[0].Aliases) == 0 {
			continue
		}
		// remove alias
		md.DenomUnits[0].Aliases = []string{}

		// update pundix/purse base denom
		newBase := strings.ToLower(md.Symbol)
		if md.Base != newBase {
			md.Base = newBase
			md.DenomUnits[0].Denom = newBase
		}

		bankKeeper.SetDenomMetaData(ctx, md)
	}
}
