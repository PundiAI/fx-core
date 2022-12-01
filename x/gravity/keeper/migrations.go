package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	v2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	cdc             codec.BinaryCodec
	sk              v2.StakingKeeper
	ak              v2.AccountKeeper
	bk              v2.BankKeeper
	gravityStoreKey sdk.StoreKey
	ethStoreKey     sdk.StoreKey
	legacyAmino     *codec.LegacyAmino
	paramsStoreKey  sdk.StoreKey
}

// NewMigrator returns a new Migrator.
func NewMigrator(cdc codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino, paramsStoreKey sdk.StoreKey,
	gravityStoreKey sdk.StoreKey, ethStoreKey sdk.StoreKey,
	sk v2.StakingKeeper, ak v2.AccountKeeper, bk v2.BankKeeper) Migrator {
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
	if err := v2.MigrateBank(ctx, m.ak, m.bk, ethtypes.ModuleName); err != nil {
		return err
	}
	multiStore := ctx.MultiStore()
	gravityStore := multiStore.GetKVStore(m.gravityStoreKey)
	ethStore := multiStore.GetKVStore(m.ethStoreKey)
	paramsStore := multiStore.GetKVStore(m.paramsStoreKey)
	v2.MigrateValidatorToOracle(ctx, m.cdc, gravityStore, ethStore, m.sk, m.bk)
	if err := v2.MigrateParams(m.legacyAmino, paramsStore, ethtypes.ModuleName); err != nil {
		return err
	}
	v2.MigrateStore(m.cdc, gravityStore, ethStore)
	return v2.CleanKVStore(gravityStore)
}
