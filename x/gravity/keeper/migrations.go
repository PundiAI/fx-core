package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	v045 "github.com/functionx/fx-core/x/gravity/legacy/v045"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper          Keeper
	cdc             codec.BinaryCodec
	sk              v045.StakingKeeper
	ak              v045.AccountKeeper
	bk              v045.BankKeeper
	gravityStoreKey sdk.StoreKey
	ethStoreKey     sdk.StoreKey

	legacyAmino *codec.LegacyAmino
	paramsKey   sdk.StoreKey
}

// NewMigrator returns a new Migrator.
func NewMigrator(cdc codec.BinaryCodec, keeper Keeper, sk v045.StakingKeeper, ak v045.AccountKeeper, bk v045.BankKeeper,
	gravityStoreKey sdk.StoreKey, ethStoreKey sdk.StoreKey,
	legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) Migrator {
	return Migrator{
		keeper:          keeper,
		cdc:             cdc,
		sk:              sk,
		ak:              ak,
		bk:              bk,
		gravityStoreKey: gravityStoreKey,
		ethStoreKey:     ethStoreKey,

		legacyAmino: legacyAmino,
		paramsKey:   paramsKey,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	if err := v045.MigrateBank(ctx, m.ak, m.bk); err != nil {
		return err
	}
	gravityStore := ctx.KVStore(m.gravityStoreKey)
	ethStore := ctx.KVStore(m.ethStoreKey)
	v045.MigrateValidatorToOracle(ctx, m.cdc, gravityStore, ethStore, m.sk)
	if err := v045.MigrateParams(ctx, m.legacyAmino, m.paramsKey); err != nil {
		return err
	}
	v045.MigrateStore(ctx, m.cdc, gravityStore, ethStore)
	return nil
}
