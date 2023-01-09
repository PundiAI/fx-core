package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/functionx/fx-core/v3/x/erc20/legacy/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper        Keeper
	channelKeeper v2.Channelkeeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, channelKeeper v2.Channelkeeper) Migrator {
	return Migrator{
		keeper:        keeper,
		channelKeeper: channelKeeper,
	}
}

func (k Migrator) Migrate2to3(ctx sdk.Context) error {
	kvStore := ctx.KVStore(k.keeper.storeKey)
	v2.PruneExpirationIBCTransferRelation(ctx, kvStore, k.channelKeeper)
	v2.MigrateIBCTransferRelation(ctx, kvStore, k.keeper)
	return nil
}
