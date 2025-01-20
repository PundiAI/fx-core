package v8

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	crosschaintypes "github.com/pundiai/fx-core/v8/x/crosschain/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (m Migrator) migrateKeys(ctx sdk.Context) error {
	store := ctx.KVStore(m.storeKey)
	if err := m.migrateParams(ctx, store); err != nil {
		return err
	}
	if err := m.migrateTokenPair(ctx, store); err != nil {
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

func (m Migrator) migrateTokenPair(ctx sdk.Context, store storetypes.KVStore) error {
	fxDenom := fxtypes.OriginalFXDenom()
	iterator := storetypes.KVStorePrefixIterator(store, KeyPrefixTokenPair)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tokenPair types.ERC20Token
		m.cdc.MustUnmarshal(iterator.Value(), &tokenPair)
		md, found := m.bankKeeper.GetDenomMetaData(ctx, tokenPair.GetDenom())
		if !found {
			return sdkerrors.ErrKeyNotFound.Wrapf("metadata not found: %s", tokenPair.GetDenom())
		}
		if md.Base == fxDenom || md.Base == strings.ToLower(md.Symbol) {
			if err := m.keeper.ERC20Token.Set(ctx, md.Base, tokenPair); err != nil {
				return err
			}
			if err := m.keeper.DenomIndex.Set(ctx, tokenPair.Erc20Address, md.Base); err != nil {
				return err
			}
			continue
		}
		tokenPair.Denom = md.Base
		if !strings.Contains(md.Base, strings.ToLower(md.Symbol)) {
			// reset pundix and purse
			tokenPair.Denom = strings.ToLower(md.Symbol)
		}
		if err := m.keeper.ERC20Token.Set(ctx, tokenPair.Denom, tokenPair); err != nil {
			return err
		}
		if err := m.keeper.DenomIndex.Set(ctx, tokenPair.Erc20Address, tokenPair.Denom); err != nil {
			return err
		}
	}
	return nil
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
