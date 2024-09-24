package keeper

import (
	"context"
	"fmt"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) AddBridgeTokenExecuted(ctx sdk.Context, claim *types.MsgBridgeTokenClaim) error {
	// Check if it already exists
	if has := k.HasBridgeToken(ctx, claim.TokenContract); has {
		return types.ErrInvalid.Wrap("bridge token is exist")
	}

	k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token",
		claim.TokenContract, "channelIbc", claim.ChannelIbc)
	if claim.Symbol == fxtypes.DefaultDenom {
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
	iter := storetypes.KVStorePrefixIterator(store, types.TokenToDenomKey)
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

func (k Keeper) HasDenom(ctx context.Context, denom string) (bool, error) {
	ok := k.bankKeeper.HasDenomMetaData(ctx, denom)
	return ok, nil
}

func (k Keeper) GetAliases(ctx context.Context, denom string) ([]string, error) {
	metadata, ok := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !ok {
		return nil, fmt.Errorf("denom %s not found", denom)
	}
	if len(metadata.DenomUnits) == 0 {
		return nil, fmt.Errorf("denom %s denom units is empty", denom)
	}
	aliases := metadata.DenomUnits[0].Aliases
	if len(aliases) == 0 {
		return nil, fmt.Errorf("denom %s aliases is empty", denom)
	}
	return aliases, nil
}

func (k Keeper) GetAllBridgeTokens(ctx context.Context) ([]types.BridgeToken, error) {
	panic("not implemented") // TODO implement me
}

func (k Keeper) UpdateAliases(ctx context.Context, denom string, aliases ...string) error {
	metadata, ok := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !ok {
		return fmt.Errorf("denom %s not found", denom)
	}
	if len(metadata.DenomUnits) == 0 {
		return fmt.Errorf("denom %s denom units is empty", denom)
	}
	metadata.DenomUnits[0].Aliases = aliases
	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	return nil
}

func (k Keeper) SetBridgeToken(ctx context.Context, name, symbol string, decimals uint32, aliases ...string) error {
	if ok := k.bankKeeper.HasDenomMetaData(ctx, strings.ToLower(name)); ok {
		return fmt.Errorf("denom %s already exist", name)
	}
	metadata := fxtypes.GetCrossChainMetadataManyToOne(name, symbol, decimals, aliases...)
	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	return nil
}
