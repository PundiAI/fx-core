package v7

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/functionx/fx-core/v7/app/keepers"
	"github.com/functionx/fx-core/v7/contract"
	fxevmkeeper "github.com/functionx/fx-core/v7/x/evm/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		ctx.Logger().Info("start to run migrations...", "module", "upgrade", "plan", plan.Name)
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		UpdateWFXLogicCode(cacheCtx, app.EvmKeeper)
		UpdateFIP20LogicCode(cacheCtx, app.EvmKeeper)

		commit()
		ctx.Logger().Info("Upgrade complete", "module", "upgrade")
		return toVM, nil
	}
}

func UpdateWFXLogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) {
	wfx := contract.GetWFX()
	if err := keeper.UpdateContractCode(ctx, wfx.Address, wfx.Code); err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update WFX contract", "module", "upgrade", "codeHash", wfx.CodeHash())
}

func UpdateFIP20LogicCode(ctx sdk.Context, keeper *fxevmkeeper.Keeper) {
	fip20 := contract.GetFIP20()
	if err := keeper.UpdateContractCode(ctx, fip20.Address, fip20.Code); err != nil {
		panic(fmt.Sprintf("update wfx logic code error: %s", err.Error()))
	}
	ctx.Logger().Info("update WFX contract", "module", "upgrade", "codeHash", fip20.CodeHash())
}
