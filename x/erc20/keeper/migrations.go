package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v3 "github.com/functionx/fx-core/v3/x/erc20/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper        Keeper
	channelKeeper v3.Channelkeeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, channelKeeper v3.Channelkeeper) Migrator {
	return Migrator{
		keeper:        keeper,
		channelKeeper: channelKeeper,
	}
}

func (k Migrator) Migrate2to3(ctx sdk.Context) error {
	kvStore := ctx.KVStore(k.keeper.storeKey)
	v3.PruneExpirationIBCTransferRelation(ctx, kvStore, k.channelKeeper)
	v3.MigrateIBCTransferRelation(ctx, kvStore, k.keeper)
	return nil
}
