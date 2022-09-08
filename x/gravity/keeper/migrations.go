package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	v3 "github.com/functionx/fx-core/v3/x/gravity/legacy/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	cdc             codec.BinaryCodec
	sk              v3.StakingKeeper
	ak              v3.AccountKeeper
	bk              v3.BankKeeper
	gravityStoreKey sdk.StoreKey
	ethStoreKey     sdk.StoreKey
	legacyAmino     *codec.LegacyAmino
	paramsStoreKey  sdk.StoreKey
}

// NewMigrator returns a new Migrator.
func NewMigrator(cdc codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino, paramsStoreKey sdk.StoreKey,
	gravityStoreKey sdk.StoreKey, ethStoreKey sdk.StoreKey,
	sk v3.StakingKeeper, ak v3.AccountKeeper, bk v3.BankKeeper) Migrator {
	return Migrator{
		cdc:             cdc,
		sk:              sk,
		ak:              ak,
		bk:              bk,
		gravityStoreKey: gravityStoreKey,
		ethStoreKey:     ethStoreKey,
		paramsStoreKey:  paramsStoreKey,
		legacyAmino:     legacyAmino,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	if err := v3.MigrateBank(ctx, m.ak, m.bk, ethtypes.ModuleName); err != nil {
		return err
	}
	gravityStore := ctx.KVStore(m.gravityStoreKey)
	ethStore := ctx.KVStore(m.ethStoreKey)
	v3.MigrateValidatorToOracle(ctx, m.cdc, gravityStore, ethStore, m.sk)
	if err := v3.MigrateParams(ctx, m.legacyAmino, m.paramsStoreKey, ethtypes.ModuleName); err != nil {
		return err
	}
	v3.MigrateStore(m.cdc, gravityStore, ethStore)
	return nil
}
