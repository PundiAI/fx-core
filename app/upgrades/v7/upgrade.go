package v7

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	autytypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v7/app/keepers"
	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	fxevmkeeper "github.com/functionx/fx-core/v7/x/evm/keeper"
)

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, app *keepers.AppKeepers) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Testnet skip
		if fxtypes.ChainId() == fxtypes.TestnetChainId {
			return fromVM, nil
		}
		// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module.
		baseAppLegacySS, found := app.ParamsKeeper.GetSubspace(baseapp.Paramspace)
		if !found {
			panic("baseapp subspace not found")
		}
		baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

		cacheCtx, commit := ctx.CacheContext()

		ctx.Logger().Info("start to run migrations...", "module", "upgrade", "plan", plan.Name)
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		UpdateWFXLogicCode(cacheCtx, app.EvmKeeper)
		UpdateFIP20LogicCode(cacheCtx, app.EvmKeeper)
		crosschainBridgeCallFrom := autytypes.NewModuleAddress(crosschaintypes.ModuleName)
		if account := app.AccountKeeper.GetAccount(ctx, crosschainBridgeCallFrom); account == nil {
			app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, crosschainBridgeCallFrom))
		}

		commit()
		ctx.Logger().Info("upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func UpdateWFXLogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) {
	wfx := contract.GetWFX()
	if err := keeper.UpdateContractCode(ctx, wfx.Address, wfx.Code); err != nil {
		ctx.Logger().Error("update WFX contract", "module", "upgrade", "err", err.Error())
	} else {
		ctx.Logger().Info("update WFX contract", "module", "upgrade", "codeHash", wfx.CodeHash())
	}
}

func UpdateFIP20LogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) {
	fip20 := contract.GetFIP20()
	if err := keeper.UpdateContractCode(ctx, fip20.Address, fip20.Code); err != nil {
		ctx.Logger().Error("update FIP20 contract", "module", "upgrade", "err", err.Error())
	} else {
		ctx.Logger().Info("update FIP20 contract", "module", "upgrade", "codeHash", fip20.CodeHash())
	}
}
