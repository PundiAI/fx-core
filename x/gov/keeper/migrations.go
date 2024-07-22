package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/exported"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	v4 "github.com/functionx/fx-core/v7/x/gov/migrations/v4"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper *Keeper
	govkeeper.Migrator
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper *Keeper, legacySubspace exported.ParamSubspace) Migrator {
	return Migrator{
		keeper:   keeper,
		Migrator: govkeeper.NewMigrator(keeper.Keeper, legacySubspace),
	}
}

// Migrate3to4 migrates from version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	if err := m.Migrator.Migrate3to4(ctx); err != nil {
		return err
	}
	params := m.keeper.Keeper.GetParams(ctx)
	return v4.MigrateFXParams(ctx, m.keeper.storeKey, m.keeper.cdc, params)
}
