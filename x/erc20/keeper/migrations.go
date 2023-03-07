package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v3 "github.com/functionx/fx-core/v3/x/erc20/migrations/v3"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	channelKeeper  v3.Channelkeeper
	legacySubspace types.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, channelKeeper v3.Channelkeeper, ss types.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		channelKeeper:  channelKeeper,
		legacySubspace: ss,
	}
}

// Migrate3to4 migrates the x/erc20 module state from the consensus version 3 to
// version 4. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/erc20
// module state.
func (k Migrator) Migrate3to4(ctx sdk.Context) error {
	var currParams types.Params
	k.legacySubspace.GetParamSet(ctx, &currParams)
	if err := currParams.Validate(); err != nil {
		return err
	}
	bz := k.keeper.cdc.MustMarshal(&currParams)
	ctx.KVStore(k.keeper.storeKey).Set(types.ParamsKey, bz)
	return nil
}
