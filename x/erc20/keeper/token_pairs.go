package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// GetAllTokenPairs - get all registered token tokenPairs
func (k Keeper) GetAllTokenPairs(ctx sdk.Context) []types.TokenPair {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixTokenPair)
	defer iterator.Close()

	var tokenPairs []types.TokenPair
	for ; iterator.Valid(); iterator.Next() {
		var tokenPair types.TokenPair
		k.cdc.MustUnmarshal(iterator.Value(), &tokenPair)
		tokenPairs = append(tokenPairs, tokenPair)
	}
	return tokenPairs
}

// GetTokenPairID returns the pair id from either of the registered tokens.
func (k Keeper) GetTokenPairID(ctx sdk.Context, token string) []byte {
	if common.IsHexAddress(token) {
		return k.GetERC20Map(ctx, common.HexToAddress(token))
	}
	return k.GetDenomMap(ctx, token)
}

func (k Keeper) GetTokenPairByAddress(ctx sdk.Context, address common.Address) (types.TokenPair, bool) {
	pairID := k.GetERC20Map(ctx, address)
	if len(pairID) == 0 {
		return types.TokenPair{}, false
	}
	return k.GetTokenPair(ctx, pairID)
}

// GetTokenPair - get registered token pair from the identifier
func (k Keeper) GetTokenPair(ctx sdk.Context, id []byte) (types.TokenPair, bool) {
	if id == nil {
		return types.TokenPair{}, false
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	bz := store.Get(id)
	if len(bz) == 0 {
		return types.TokenPair{}, false
	}

	var tokenPair types.TokenPair
	k.cdc.MustUnmarshal(bz, &tokenPair)
	return tokenPair, true
}

// SetTokenPair stores a token pair
func (k Keeper) SetTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	key := tokenPair.GetID()
	bz := k.cdc.MustMarshal(&tokenPair)
	store.Set(key, bz)
}

func (k Keeper) AddTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	id := tokenPair.GetID()
	bz := k.cdc.MustMarshal(&tokenPair)
	store.Set(id, bz)

	k.SetDenomMap(ctx, tokenPair.Denom, id)
	k.SetERC20Map(ctx, tokenPair.GetERC20Contract(), id)
}

// RemoveTokenPair removes a token pair.
func (k Keeper) RemoveTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	id := tokenPair.GetID()
	k.DeleteTokenPair(ctx, id)
	k.DeleteERC20Map(ctx, tokenPair.GetERC20Contract())
	k.DeleteDenomMap(ctx, tokenPair.Denom)

	//delete denom alias
	if md, found := k.HasDenomAlias(ctx, tokenPair.Denom); found {
		k.DeleteAliasesDenom(ctx, md.DenomUnits[0].Aliases...)
	}
}

// DeleteTokenPair deletes the token pair for the given id
func (k Keeper) DeleteTokenPair(ctx sdk.Context, id []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	store.Delete(id)
}

// GetERC20Map returns the token pair id for the given address
func (k Keeper) GetERC20Map(ctx sdk.Context, erc20 common.Address) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20)
	return store.Get(erc20.Bytes())
}

// GetDenomMap returns the token pair id for the given denomination
func (k Keeper) GetDenomMap(ctx sdk.Context, denom string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	return store.Get([]byte(denom))
}

// SetERC20Map sets the token pair id for the given address
func (k Keeper) SetERC20Map(ctx sdk.Context, erc20 common.Address, id []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20)
	store.Set(erc20.Bytes(), id)
}

// DeleteERC20Map deletes the token pair id for the given address
func (k Keeper) DeleteERC20Map(ctx sdk.Context, erc20 common.Address) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20)
	store.Delete(erc20.Bytes())
}

// SetDenomMap sets the token pair id for the denomination
func (k Keeper) SetDenomMap(ctx sdk.Context, denom string, id []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	store.Set([]byte(denom), id)
}

// DeleteDenomMap deletes the token pair id for the given denom
func (k Keeper) DeleteDenomMap(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	store.Delete([]byte(denom))
}

// IsTokenPairRegistered - check if registered token tokenPair is registered
func (k Keeper) IsTokenPairRegistered(ctx sdk.Context, id []byte) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPair)
	return store.Has(id)
}

// IsERC20Registered check if registered ERC20 token is registered
func (k Keeper) IsERC20Registered(ctx sdk.Context, erc20 common.Address) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByERC20)
	return store.Has(erc20.Bytes())
}

// IsDenomRegistered check if registered coin denom is registered
func (k Keeper) IsDenomRegistered(ctx sdk.Context, denom string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenPairByDenom)
	return store.Has([]byte(denom))
}

// SetAliasesDenom sets the aliases for the denomination
func (k Keeper) SetAliasesDenom(ctx sdk.Context, denom string, aliases ...string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAliasDenom)
	for _, alias := range aliases {
		store.Set([]byte(alias), []byte(denom))
	}
}

// GetAliasDenom returns the denom for the given alias
func (k Keeper) GetAliasDenom(ctx sdk.Context, alias string) (string, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAliasDenom)
	value := store.Get([]byte(alias))
	if value == nil {
		return "", false
	}
	return string(value), true
}

// DeleteAliasesDenom deletes the denom-alias for the given alias
func (k Keeper) DeleteAliasesDenom(ctx sdk.Context, aliases ...string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAliasDenom)
	for _, alias := range aliases {
		store.Delete([]byte(alias))
	}
}

// IsAliasDenomRegistered check if registered coin alias is registered
func (k Keeper) IsAliasDenomRegistered(ctx sdk.Context, alias string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAliasDenom)
	return store.Has([]byte(alias))
}
