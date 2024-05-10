package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
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

// GetTokenPair - get registered token pair from the token or denom
func (k Keeper) GetTokenPair(ctx sdk.Context, tokenOrDenom string) (types.TokenPair, bool) {
	var id []byte
	if common.IsHexAddress(tokenOrDenom) {
		id = k.getERC20Map(ctx, common.HexToAddress(tokenOrDenom))
	} else {
		id = k.getDenomMap(ctx, tokenOrDenom)
	}
	if id == nil {
		return types.TokenPair{}, false
	}
	return k.getTokenPair(ctx, id)
}

func (k Keeper) GetTokenPairByAddress(ctx sdk.Context, address common.Address) (types.TokenPair, bool) {
	pairID := k.getERC20Map(ctx, address)
	if len(pairID) == 0 {
		return types.TokenPair{}, false
	}
	return k.getTokenPair(ctx, pairID)
}

func (k Keeper) getTokenPair(ctx sdk.Context, id []byte) (types.TokenPair, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(types.KeyPrefixTokenPair, id...))
	if len(bz) == 0 {
		return types.TokenPair{}, false
	}
	var tokenPair types.TokenPair
	k.cdc.MustUnmarshal(bz, &tokenPair)
	return tokenPair, true
}

func (k Keeper) AddTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := ctx.KVStore(k.storeKey)
	id := tokenPair.GetID()
	// set pair
	store.Set(append(types.KeyPrefixTokenPair, id...), k.cdc.MustMarshal(&tokenPair))

	// set denom map
	store.Set(append(types.KeyPrefixTokenPairByDenom, []byte(tokenPair.Denom)...), id)

	// set erc20 map
	store.Set(append(types.KeyPrefixTokenPairByERC20, tokenPair.GetERC20Contract().Bytes()...), id)
}

func (k Keeper) SetTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.KeyPrefixTokenPair, tokenPair.GetID()...), k.cdc.MustMarshal(&tokenPair))
}

// RemoveTokenPair removes a token pair.
func (k Keeper) RemoveTokenPair(ctx sdk.Context, tokenPair types.TokenPair) {
	store := ctx.KVStore(k.storeKey)
	id := tokenPair.GetID()
	// delete token pair
	store.Delete(append(types.KeyPrefixTokenPair, id...))

	// delete denom map
	store.Delete(append(types.KeyPrefixTokenPairByDenom, []byte(tokenPair.Denom)...))

	// delete erc20 map
	store.Delete(append(types.KeyPrefixTokenPairByERC20, tokenPair.GetERC20Contract().Bytes()...))

	// delete denom alias
	if md, found := k.HasDenomAlias(ctx, tokenPair.Denom); found {
		k.DeleteAliasesDenom(ctx, md.DenomUnits[0].Aliases...)
	}
}

// getERC20Map returns the token pair id for the given address
func (k Keeper) getERC20Map(ctx sdk.Context, erc20 common.Address) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(append(types.KeyPrefixTokenPairByERC20, erc20.Bytes()...))
}

// IsERC20Registered check if registered ERC20 token is registered
func (k Keeper) IsERC20Registered(ctx sdk.Context, erc20 common.Address) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(types.KeyPrefixTokenPairByERC20, erc20.Bytes()...))
}

// getDenomMap returns the token pair id for the given denomination
func (k Keeper) getDenomMap(ctx sdk.Context, denom string) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(append(types.KeyPrefixTokenPairByDenom, []byte(denom)...))
}

// IsDenomRegistered check if registered coin denom is registered
func (k Keeper) IsDenomRegistered(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(types.KeyPrefixTokenPairByDenom, []byte(denom)...))
}

// SetAliasesDenom sets the aliases for the denomination
func (k Keeper) SetAliasesDenom(ctx sdk.Context, denom string, aliases ...string) {
	store := ctx.KVStore(k.storeKey)
	for _, alias := range aliases {
		store.Set(append(types.KeyPrefixAliasDenom, []byte(alias)...), []byte(denom))
	}
}

// GetAliasDenom returns the denom for the given alias
func (k Keeper) GetAliasDenom(ctx sdk.Context, alias string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(append(types.KeyPrefixAliasDenom, []byte(alias)...))
	if value == nil {
		return "", false
	}
	return string(value), true
}

// DeleteAliasesDenom deletes the denom-alias for the given alias
func (k Keeper) DeleteAliasesDenom(ctx sdk.Context, aliases ...string) {
	store := ctx.KVStore(k.storeKey)
	for _, alias := range aliases {
		store.Delete(append(types.KeyPrefixAliasDenom, []byte(alias)...))
	}
}

// IsAliasDenomRegistered check if registered coin alias is registered
func (k Keeper) IsAliasDenomRegistered(ctx sdk.Context, alias string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(append(types.KeyPrefixAliasDenom, []byte(alias)...))
}

// IsOriginDenom check denom origin
// denom must be eth0x...|bsc0x...|tronTL...|ibc/ABC...
// origin: FX, fxUSD(fxevm-ERC20)
// cross: PUNDIX, PURSE, Other(usdt,usdc)
func (k Keeper) IsOriginDenom(ctx sdk.Context, denom string) bool {
	// exclude ethereum FX
	if strings.EqualFold(denom, fxtypes.DefaultDenom) {
		return true
	}

	// exclude PUNDIX|PURSE
	if k.IsDenomRegistered(ctx, denom) {
		return false
	}

	// exclude denom not register (may be add bridge token, but not register)
	baseDenom, found := k.GetAliasDenom(ctx, denom)
	if !found {
		// false mean to mint it, because we don't know what type it is.
		// if type not origin, we need mint it, false is correct
		// if type is origin(FX/fxUSD token in ethereum), we need unlock it, false is incorrect,
		// so there is a problem here, we need remove token from bridge contract and fix it when upgrade chain
		ctx.Logger().Info("denom not register, but add to bridge contract", "denom", denom)
		return false
	}

	// exclude other chain FX token
	if strings.EqualFold(baseDenom, fxtypes.DefaultDenom) {
		return false
	}

	tokenPair, found := k.GetTokenPair(ctx, baseDenom)
	if !found {
		ctx.Logger().Info("alias register, but denom token pair not found", "alias", denom, "denom", baseDenom)
		return false
	}

	return tokenPair.IsNativeERC20()
}

func (k Keeper) IsConvertedAlias(ctx sdk.Context, denom string) bool {
	// exclude ethereum FX
	if strings.EqualFold(denom, fxtypes.DefaultDenom) {
		return false
	}

	// exclude PUNDIX|PURSE
	if k.IsDenomRegistered(ctx, denom) {
		return false
	}

	// exclude denom not register (may be add bridge token, but not register)
	baseDenom, found := k.GetAliasDenom(ctx, denom)
	if !found {
		return false
	}
	return k.checkConvertedDenom(baseDenom)
}

func (k Keeper) IsOriginOrConvertedDenom(ctx sdk.Context, denom string) bool {
	if k.IsOriginDenom(ctx, denom) {
		return true
	}
	return k.IsConvertedAlias(ctx, denom)
}

func (k Keeper) checkConvertedDenom(baseDenom string) bool {
	if baseDenom == fxtypes.DefaultDenom {
		return true
	}
	if strings.HasPrefix(baseDenom, ibctransfertypes.DenomPrefix+"/") {
		return true
	}

	for _, chainName := range k.chainsName {
		if len(baseDenom) > len(chainName) && strings.HasPrefix(baseDenom, chainName) &&
			crosschaintypes.ValidateExternalAddr(chainName, strings.TrimPrefix(baseDenom, chainName)) == nil {
			return true
		}
	}
	return false
}
