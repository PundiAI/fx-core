package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v4 "github.com/functionx/fx-core/v7/x/crosschain/migrations/v4"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
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

// Migrate3to4 migrates the x/crosschain module state from the consensus version 3 to
// version 4. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/crosschain
// module state.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	if err := v4.MigratorParam(ctx, m.legacySubspace, m.keeper.storeKey, m.keeper.cdc); err != nil {
		return err
	}
	return nil
}
