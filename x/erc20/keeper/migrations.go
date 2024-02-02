package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v4 "github.com/functionx/fx-core/v7/x/erc20/migrations/v4"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace types.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, ss types.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: ss,
	}
}

// Migrate3to4 migrates the x/erc20 module state from the consensus version 3 to
// version 4. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/erc20
// module state.
func (k Migrator) Migrate3to4(ctx sdk.Context) error {
	if err := v4.MigratorParam(ctx, k.legacySubspace, k.keeper.storeKey, k.keeper.cdc); err != nil {
		return err
	}
	return nil
}
