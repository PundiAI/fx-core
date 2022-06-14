package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v045 "github.com/functionx/fx-core/x/crosschain/legacy/v045"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper    Keeper
	sk        v045.StakingKeeper
	paramsKey sdk.StoreKey
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, sk v045.StakingKeeper, paramsKey sdk.StoreKey) Migrator {
	return Migrator{
		keeper:    keeper,
		sk:        sk,
		paramsKey: paramsKey,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	if err := v045.MigrateParams(ctx, &m.keeper.paramSpace, m.keeper.moduleName, m.paramsKey); err != nil {
		return err
	}
	oracles, err := v045.MigrateOracle(ctx, m.keeper.cdc, m.keeper.storeKey, m.sk)
	if err != nil {
		return err
	}
	if err := v045.MigrateDepositToStaking(ctx, oracles, m.sk); err != nil {
		return err
	}
	v045.MigratePruneIbcSequenceKey(ctx, m.keeper.storeKey)
	return nil
}
