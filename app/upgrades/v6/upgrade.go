package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v6/app/keepers"
	migratekeeper "github.com/functionx/fx-core/v6/x/migrate/keeper"
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

func AutoUndelegate(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper) []stakingtypes.Delegation {
	var delegations []stakingtypes.Delegation
	stakingKeeper.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) (stop bool) {
		delegations = append(delegations, delegation)
		delegator := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		valAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		if delegator.Equals(valAddress) {
			return false
		}
		if _, err := stakingKeeper.Undelegate(ctx, delegator, valAddress, delegation.Shares); err != nil {
			panic(err)
		}
		return false
	})
	return delegations
}

func ExportDelegate(ctx sdk.Context, delegations []stakingtypes.Delegation, migrateKeeper migratekeeper.Keeper) []stakingtypes.Delegation {
	for i := 0; i < len(delegations); i++ {
		delegation := delegations[i]
		delegator := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		if !migrateKeeper.HasMigratedDirectionTo(ctx, common.BytesToAddress(delegator.Bytes())) {
			delegations = append(delegations[:i], delegations[i+1:]...)
			i--
			continue
		}
	}
	return delegations
}
