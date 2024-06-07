package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) AddBridgeTokenExecuted(ctx sdk.Context, claim *types.MsgBridgeTokenClaim) error {
	// Check if it already exists
	if has := k.HasBridgeToken(ctx, claim.TokenContract); has {
		return types.ErrInvalid.Wrap("bridge token is exist")
	}

	k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token",
		claim.TokenContract, "channelIbc", claim.ChannelIbc)
	if claim.Symbol == fxtypes.DefaultDenom {
		// Check if denom exists
		if !k.bankKeeper.HasDenomMetaData(ctx, claim.Symbol) {
			return types.ErrUnknown.Wrapf("denom not found %s", claim.Symbol)
		}

		if uint64(fxtypes.DenomUnit) != claim.Decimals {
			return types.ErrInvalid.Wrapf("%s denom decimals not match %d, expect %d", fxtypes.DefaultDenom,
				claim.Decimals, fxtypes.DenomUnit)
		}

		k.AddBridgeToken(ctx, claim.TokenContract, fxtypes.DefaultDenom)
		return nil
	}

	denom, err := k.SetIbcDenomTrace(ctx, claim.TokenContract, claim.ChannelIbc)
	if err != nil {
		return err
	}
	k.AddBridgeToken(ctx, claim.TokenContract, denom)
	return nil
}

func (k Keeper) GetBridgeTokenDenom(ctx sdk.Context, tokenContract string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetDenomToTokenKey(tokenContract))
	if len(data) == 0 {
		return nil
	}
	return &types.BridgeToken{
		Denom: string(data),
		Token: tokenContract,
	}
}

func (k Keeper) GetDenomBridgeToken(ctx sdk.Context, denom string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetTokenToDenomKey(denom))
	if len(data) == 0 {
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
	denom := types.NewBridgeDenom(k.moduleName, token)
	denomTrace, err := fxtypes.GetIbcDenomTrace(denom, channelIBC)
	if err != nil {
		return denom, err
	}
	if denomTrace.Path != "" {
		k.ibcTransferKeeper.SetDenomTrace(ctx, denomTrace)
		return denomTrace.IBCDenom(), nil
	}
	return denomTrace.BaseDenom, nil
}

func (k Keeper) TransferBridgeCoinToExternal(ctx sdk.Context, sender sdk.AccAddress, targetCoin sdk.Coin) error {
	// lock coins in module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, k.moduleName, sdk.NewCoins(targetCoin)); err != nil {
		return err
	}
	isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, targetCoin.Denom)
	if isOriginOrConverted {
		return nil
	}
	// If it is an external blockchain asset, burn vouchers to send them back to external blockchain
	return k.bankKeeper.BurnCoins(ctx, k.moduleName, sdk.NewCoins(targetCoin))
}
