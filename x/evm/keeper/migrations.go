package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/x/evm/keeper"

	v3 "github.com/functionx/fx-core/v3/x/evm/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper.Migrator
	legacyAmino    *codec.LegacyAmino
	paramsStoreKey storetypes.StoreKey
}

// NewMigrator returns a new Migrator.
func NewMigrator(k *Keeper, legacyAmino *codec.LegacyAmino, paramsStoreKey storetypes.StoreKey) Migrator {
	return Migrator{
		Migrator:       keeper.NewMigrator(*k.Keeper),
		legacyAmino:    legacyAmino,
		paramsStoreKey: paramsStoreKey,
	}
}

// Migrate2to3 migrates the store from consensus version v2 to v3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	v3.MigrateParams(ctx, m.legacyAmino, m.paramsStoreKey)
	return m.Migrator.Migrate2to3(ctx)
}
