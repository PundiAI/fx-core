package v8

import (
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	legacytypes "github.com/pundiai/fx-core/v8/types/legacy"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (m Migrator) migrateKeys(ctx sdk.Context) error {
	store := ctx.KVStore(m.storeKey)
	if err := m.migrateParams(ctx, store); err != nil {
		return err
	}
	return m.migrateTokenPair(ctx, store)
}

func (m Migrator) migrateParams(ctx sdk.Context, store storetypes.KVStore) error {
	bz := store.Get(ParamsKey)
	if len(bz) == 0 {
		return nil
	}
	var legacyParams legacytypes.LegacyERC20Params
	m.cdc.MustUnmarshal(bz, &legacyParams)

	store.Delete(ParamsKey)
	return m.keeper.Params.Set(ctx, types.Params{EnableErc20: legacyParams.EnableErc20})
}

func (m Migrator) migrateTokenPair(ctx sdk.Context, store storetypes.KVStore) error {
	iterator := storetypes.KVStorePrefixIterator(store, KeyPrefixTokenPair)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var tokenPair types.ERC20Token
		m.cdc.MustUnmarshal(iterator.Value(), &tokenPair)
		md, found := m.bankKeeper.GetDenomMetaData(ctx, tokenPair.GetDenom())
		if !found {
			return sdkerrors.ErrKeyNotFound.Wrapf("metadata not found: %s", tokenPair.GetDenom())
		}
		if md.Base == fxtypes.LegacyFXDenom || md.Base == strings.ToLower(md.Symbol) {
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
