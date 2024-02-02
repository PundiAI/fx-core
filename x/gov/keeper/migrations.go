package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	fxv046 "github.com/functionx/fx-core/v7/x/gov/migrations/v046"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	cdc    codec.BinaryCodec
	keeper Keeper
	govkeeper.Migrator
}

// NewMigrator returns a new Migrator.
func NewMigrator(cdc codec.Codec, keeper Keeper) Migrator {
	return Migrator{
		cdc:      cdc,
		keeper:   keeper,
		Migrator: govkeeper.NewMigrator(keeper.Keeper),
	}
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return fxv046.MigrateStore(ctx, m.keeper.storeKey, m.cdc)
}
