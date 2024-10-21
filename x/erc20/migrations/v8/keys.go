package v8

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func (m Migrator) migrateKeys(ctx sdk.Context) error {
	store := ctx.KVStore(m.storeKey)
	if err := m.migrateParams(ctx, store); err != nil {
		return err
	}

	return m.migrateRelationToCache(ctx, store)
}

func (m Migrator) migrateRelationToCache(ctx sdk.Context, store storetypes.KVStore) error {
	// 1. migrate ibc transfer relation
	iterator := storetypes.KVStorePrefixIterator(store, KeyPrefixIBCTransfer)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		newKey := key[len(KeyPrefixIBCTransfer):]
		if err := m.keeper.Cache.Set(ctx, string(newKey), sdkmath.ZeroInt()); err != nil {
			return err
		}
		store.Delete(key)
	}

	// 2. migrate outgoing transfer relation
	iterator = storetypes.KVStorePrefixIterator(store, KeyPrefixOutgoingTransfer)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		originTokenKey := OutgoingTransferKeyToOriginTokenKey(key)
		if err := m.keeper.Cache.Set(ctx, originTokenKey, sdkmath.ZeroInt()); err != nil {
			return err
		}
		store.Delete(key)
	}
	return nil
}

func (m Migrator) migrateParams(ctx sdk.Context, store storetypes.KVStore) error {
	bz := store.Get(ParamsKey)
	if len(bz) == 0 {
		return nil
	}
	var legacyParams types.LegacyParams
	m.cdc.MustUnmarshal(bz, &legacyParams)

	store.Delete(ParamsKey)
	return m.keeper.Params.Set(ctx, types.Params{EnableErc20: legacyParams.EnableErc20})
}

func OutgoingTransferKeyToOriginTokenKey(key []byte) string {
	// key = prefix + moduleName(string-) + sdk.Uint64ToBigEndian(id)(len = 8)
	// 1. remove prefix
	key = key[len(KeyPrefixOutgoingTransfer):]
	// 2. get moduleName
	moduleName := string(key[:len(key)-8])
	// 3. get id
	id := sdk.BigEndianToUint64(key[len(key)-8:])
	// 4. new originTokenKey
	return crosschaintypes.NewOriginTokenKey(moduleName, id)
}
