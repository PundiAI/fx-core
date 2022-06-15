package v2

import (
	"fmt"
	"strings"

	ethtypes "github.com/functionx/fx-core/x/eth/types"

	migratetypes "github.com/functionx/fx-core/x/migrate/types"

	erc20types "github.com/functionx/fx-core/x/erc20/types"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	fxtypes "github.com/functionx/fx-core/types"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(
	mm *module.Manager, configurator module.Configurator,
	bankStoreKey *sdk.KVStoreKey, bankKeeper bankKeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper, erc20Keeper erc20keeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		// update FX metadata
		if err := UpdateFXMetadata(cacheCtx, bankKeeper, bankStoreKey); err != nil {
			return nil, err
		}

		// set max expected block time parameter. Replace the default with your expected value
		// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
		ibcKeeper.ConnectionKeeper.SetParams(cacheCtx, ibcconnectiontypes.DefaultParams())

		// cosmos-sdk 0.42.x from version must be empty
		if len(fromVM) != 0 {
			panic("invalid from version map")
		}

		for n, m := range mm.Modules {
			//NOTE: fromVM empty
			if initGenesis[n] {
				continue
			}
			if v, ok := runMigrates[n]; ok {
				fromVM[n] = v
				continue
			}
			fromVM[n] = m.ConsensusVersion()
		}

		if mm.OrderMigrations == nil {
			mm.OrderMigrations = migrationsOrder(mm.ModuleNames())
		}
		cacheCtx.Logger().Info("start to run module v2 migrations...")
		toVersion, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return nil, fmt.Errorf("run migrations error %s", err.Error())
		}

		// register coin
		for _, metadata := range fxtypes.GetMetadata() {
			cacheCtx.Logger().Info("add metadata", "coin", metadata.String())
			pair, err := erc20Keeper.RegisterCoin(cacheCtx, metadata)
			if err != nil {
				return nil, fmt.Errorf("register %s error %s", metadata.Base, err.Error())
			}
			cacheCtx.EventManager().EmitEvent(sdk.NewEvent(
				erc20types.EventTypeRegisterCoin,
				sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
			))
		}

		//commit upgrade
		commit()

		return toVersion, nil
	}
}

func UpdateFXMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper, key *types.KVStoreKey) error {
	md := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if err := md.Validate(); err != nil {
		return fmt.Errorf("invalid %s metadata", fxtypes.DefaultDenom)
	}
	ctx.Logger().Info("update metadata", "metadata", md.String())
	//delete fx
	deleteMetadata(ctx, key, strings.ToLower(fxtypes.DefaultDenom))
	//set FX
	bankKeeper.SetDenomMetaData(ctx, md)
	return nil
}

func deleteMetadata(ctx sdk.Context, key *types.KVStoreKey, base ...string) {
	store := ctx.KVStore(key)
	for _, b := range base {
		store.Delete(banktypes.DenomMetadataKey(b))
	}
}

func migrationsOrder(modules []string) []string {
	modules = module.DefaultMigrationsOrder(modules)
	for i, name := range modules {
		if name == migratetypes.ModuleName {
			modules = append(append(modules[:i], modules[i+1:]...), name)
			return modules
		}
		// eth module
		if name == ethtypes.ModuleName {
			modules = append([]string{name}, append(modules[:i], modules[i+1:]...)...)
			return modules
		}
	}
	return modules
}
