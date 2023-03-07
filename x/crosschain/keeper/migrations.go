package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
