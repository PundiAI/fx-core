package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

func (k Keeper) GetBridgeTokenDenom(ctx sdk.Context, tokenContract string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)

	data := store.Get(types.GetDenomToTokenKey(tokenContract))
	if len(data) <= 0 {
		return nil
	}
	var bridgeToken types.BridgeToken
	k.cdc.MustUnmarshal(data, &bridgeToken)
	return &bridgeToken
}

func (k Keeper) GetDenomByBridgeToken(ctx sdk.Context, denom string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)

	data := store.Get(types.GetTokenToDenomKey(denom))
	if len(data) <= 0 {
		return nil
	}
	var bridgeToken types.BridgeToken
	k.cdc.MustUnmarshal(data, &bridgeToken)
	return &bridgeToken
}

func (k Keeper) hasBridgeToken(ctx sdk.Context, tokenContract string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDenomToTokenKey(tokenContract))
}

func (k Keeper) addBridgeToken(ctx sdk.Context, token, _, channelIBC string) (string, error) {
	store := ctx.KVStore(k.storeKey)
	decodeChannelIBC, err := hex.DecodeString(channelIBC)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "decode channel ibc err")
	}

	decodeChannelIBCStr := string(decodeChannelIBC)
	denom := fmt.Sprintf("%s%s", k.moduleName, token)
	if len(decodeChannelIBCStr) > 0 {
		denomTrace := ibctransfertypes.DenomTrace{
			Path:      decodeChannelIBCStr,
			BaseDenom: denom,
		}
		k.ibcTransferKeeper.SetDenomTrace(ctx, denomTrace)
		denom = denomTrace.IBCDenom()
	}
	store.Set(types.GetTokenToDenomKey(denom), k.cdc.MustMarshal(&types.BridgeToken{
		Token:      token,
		ChannelIbc: decodeChannelIBCStr,
	}))
	store.Set(types.GetDenomToTokenKey(token), k.cdc.MustMarshal(&types.BridgeToken{
		Denom:      denom,
		ChannelIbc: decodeChannelIBCStr,
	}))
	return denom, nil
}

// IterateBridgeTokenToDenom iterates over token to denom relations
func (k Keeper) IterateBridgeTokenToDenom(ctx sdk.Context, cb func([]byte, *types.BridgeToken) bool) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.TokenToDenomKey)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bridgeToken types.BridgeToken
		k.cdc.MustUnmarshal(iter.Value(), &bridgeToken)
		bridgeToken.Denom = string(iter.Value())
		// cb returns true to stop early
		if cb(iter.Key(), &bridgeToken) {
			break
		}
	}
}
