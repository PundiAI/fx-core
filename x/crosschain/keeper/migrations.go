package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crosschainv3 "github.com/functionx/fx-core/v3/x/crosschain/migrations/v3"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace types.Subspace
}

func NewMigrator(k Keeper, ss types.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	// todo  update params
	//	m.legacySubspace.Set(ctx, types.ParamsStoreKeySignedWindow, uint64(30_000))
	//	m.legacySubspace.Set(ctx, types.ParamsStoreSlashFraction, sdk.NewDecWithPrec(8, 1)) // 80%

	crosschainv3.PruneStore(m.keeper.cdc, ctx.KVStore(m.keeper.storeKey))
	ctx.Logger().Info("prune store has been successfully", "module", m.keeper.moduleName)

	kvStore := ctx.KVStore(m.keeper.storeKey)
	crosschainv3.MigrateBridgeToken(m.keeper.cdc, kvStore)
	ctx.Logger().Info("bridge token has been migrated successfully", "module", m.keeper.moduleName)

	// fix oracle delegate
	validatorsByPower := m.keeper.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validatorsByPower) <= 0 {
		panic("no found bonded validator")
	}
	validator := validatorsByPower[0].GetOperator()
	oracles := m.keeper.GetAllOracles(ctx, false)
	return crosschainv3.MigrateDepositToStaking(ctx, m.keeper.moduleName, m.keeper.stakingKeeper, m.keeper.stakingMsgServer, m.keeper.bankKeeper, oracles, validator)
}

// Migrate3to4 migrates the x/crosschain module state from the consensus version 3 to
// version 4. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/crosschain
// module state.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	var currParams types.Params
	m.legacySubspace.GetParamSet(ctx, &currParams)
	if err := currParams.ValidateBasic(); err != nil {
		return err
	}
	bz := m.keeper.cdc.MustMarshal(&currParams)
	ctx.KVStore(m.keeper.storeKey).Set(types.ParamsKey, bz)
	return nil
}
