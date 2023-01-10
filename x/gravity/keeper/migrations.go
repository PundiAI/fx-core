package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
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
func NewMigrator(cdc codec.BinaryCodec, legacyAmino *codec.LegacyAmino,
	paramsStoreKey sdk.StoreKey, gravityStoreKey sdk.StoreKey, ethStoreKey sdk.StoreKey,
	sk v2.StakingKeeper, ak v2.AccountKeeper, bk v2.BankKeeper,
) Migrator {
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
	ctx.Logger().Info("migrating module gravity to module eth", "module", "gravity")
	if err := v2.MigrateBank(ctx, m.ak, m.bk, ethtypes.ModuleName); err != nil {
		return err
	}
	gravityStore := ctx.MultiStore().GetKVStore(m.gravityStoreKey)
	ethStore := ctx.MultiStore().GetKVStore(m.ethStoreKey)
	paramsStore := ctx.MultiStore().GetKVStore(m.paramsStoreKey)
	if err := v2.MigrateParams(m.legacyAmino, paramsStore, ethtypes.ModuleName); err != nil {
		return err
	}
	ctx.Logger().Info("params has been migrated successfully", "module", "gravity")

	var metadatas []banktypes.Metadata
	m.bk.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		metadatas = append(metadatas, metadata)
		return false
	})
	v2.MigrateBridgeTokenFromMetadatas(metadatas, ethStore)
	ctx.Logger().Info("bridge token has been migrated successfully", "module", "gravity")

	v2.MigrateValidatorToOracle(ctx, m.cdc, gravityStore, ethStore, m.sk, m.bk)

	v2.MigrateStore(m.cdc, gravityStore, ethStore)
	ctx.Logger().Info("store key has been migrated successfully", "module", "gravity",
		"LatestOracleSetNonce", sdk.BigEndianToUint64(ethStore.Get(crosschaintypes.LatestOracleSetNonce)))
	return nil
}
