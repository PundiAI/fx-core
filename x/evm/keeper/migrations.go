package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/evmos/ethermint/x/evm/keeper"
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
