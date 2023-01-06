package keeper

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func (k Keeper) GetBridgeTokenDenom(ctx sdk.Context, tokenContract string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetDenomToTokenKey(tokenContract))
	if len(data) <= 0 {
		return nil
	}
	return &types.BridgeToken{
		Denom: string(data),
		Token: tokenContract,
	}
}

func (k Keeper) GetDenomByBridgeToken(ctx sdk.Context, denom string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetTokenToDenomKey(denom))
	if len(data) <= 0 {
		return nil
	}
	return &types.BridgeToken{
		Denom: denom,
		Token: string(data),
	}
}

func (k Keeper) HasBridgeToken(ctx sdk.Context, tokenContract string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDenomToTokenKey(tokenContract))
}

func (k Keeper) AddBridgeToken(ctx sdk.Context, token, denom string) {
	store := ctx.KVStore(k.storeKey)
	// todo need remove after test completion
	if denom != fxtypes.DefaultDenom && !strings.HasPrefix(denom, k.moduleName) {
		panic("invalid denom: " + denom)
	}
	store.Set(types.GetTokenToDenomKey(denom), []byte(token))
	store.Set(types.GetDenomToTokenKey(token), []byte(denom))
}

func (k Keeper) IterateBridgeTokenToDenom(ctx sdk.Context, cb func(*types.BridgeToken) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenToDenomKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		bridgeToken := &types.BridgeToken{
			Denom: string(iter.Key()[len(types.TokenToDenomKey):]),
			Token: string(iter.Value()),
		}
		// cb returns true to stop early
		if cb(bridgeToken) {
			break
		}
	}
}

func (k Keeper) SetIbcDenomTrace(ctx sdk.Context, token, channelIBC string) (string, error) {
	channelPath, err := hex.DecodeString(channelIBC)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "decode channel ibc err")
	}

	// todo need check path
	path := string(channelPath)
	denom := fmt.Sprintf("%s%s", k.moduleName, token)
	if len(path) > 0 {
		denomTrace := ibctransfertypes.DenomTrace{
			Path:      path,
			BaseDenom: denom,
		}
		k.ibcTransferKeeper.SetDenomTrace(ctx, denomTrace)
	}
	return denom, nil
}
