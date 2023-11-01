package v6

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v6/app/keepers"
	crosschainkeeper "github.com/functionx/fx-core/v6/x/crosschain/keeper"
	layer2types "github.com/functionx/fx-core/v6/x/layer2/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		ctx.Logger().Info("start to run v6 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		commit()
		ctx.Logger().Info("Upgrade complete")
		return toVM, nil
	}
}

func MigrateMetadata(ctx sdk.Context, bankKeeper bankkeeper.Keeper) {
	bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		address, ok := Layer2GenesisTokenAddress[metadata.Symbol]
		if !ok {
			return false
		}
		if len(metadata.DenomUnits) > 0 {
			metadata.DenomUnits[0].Aliases = append(metadata.DenomUnits[0].Aliases,
				fmt.Sprintf("%s%s", layer2types.ModuleName, address))
		}
		return false
	})
}

func MigrateLayer2Module(ctx sdk.Context, bankKeeper bankkeeper.Keeper, layer2CrossChainKeeper crosschainkeeper.Keeper) {
	for _, address := range Layer2GenesisTokenAddress {
		fxTokenDenom := fmt.Sprintf("%s%s", layer2types.ModuleName, address)
		layer2CrossChainKeeper.AddBridgeToken(ctx, address, fxTokenDenom)
	}
}
