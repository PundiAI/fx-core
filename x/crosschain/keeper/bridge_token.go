package keeper

import (
	"context"
	"fmt"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (k Keeper) AddBridgeTokenExecuted(ctx sdk.Context, claim *types.MsgBridgeTokenClaim) error {
	k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token",
		claim.TokenContract, "channelIbc", claim.ChannelIbc)
	bridgeDenom := types.NewBridgeDenom(k.moduleName, claim.TokenContract)
	denomToken := bridgeDenom

	// Check if it already exists
	if has := k.HasBridgeToken(ctx, bridgeDenom); has {
		return types.ErrInvalid.Wrapf("bridge token is exist %s", bridgeDenom)
	}

	if claim.Symbol == fxtypes.DefaultDenom {
		if uint64(fxtypes.DenomUnit) != claim.Decimals {
			return types.ErrInvalid.Wrapf("%s denom decimals not match %d, expect %d", fxtypes.DefaultDenom,
				claim.Decimals, fxtypes.DenomUnit)
		}
		k.AddBridgeToken(ctx, fxtypes.DefaultDenom, bridgeDenom)
		denomToken = fxtypes.DefaultDenom
	}

	k.AddBridgeToken(ctx, bridgeDenom, denomToken)
	return nil
}

func (k Keeper) GetBridgeDenomByContract(ctx sdk.Context, tokenContract string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bridgeDenom := types.NewBridgeDenom(k.moduleName, tokenContract)
	data := store.Get(types.GetBridgeDenomKey(bridgeDenom))
	if len(data) == 0 {
		return "", false
	}
	result := string(data)
	// result = (value == key ？key : value)
	if bridgeDenom == result {
		result = bridgeDenom
	}
	return result, true
}

func (k Keeper) GetContractByBridgeDenom(ctx sdk.Context, bridgeDenom string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetBridgeDenomKey(bridgeDenom))
	if len(data) == 0 {
		return "", false
	}
	result := string(data)
	if bridgeDenom == result || result == fxtypes.DefaultDenom {
		return types.BridgeDenomToContract(k.moduleName, bridgeDenom), true
	}
	// bridgeDenom should not be eth0xfx
	return types.BridgeDenomToContract(k.moduleName, result), true
}

func (k Keeper) HasBridgeToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetTokenToDenomKey(denom))
}

func (k Keeper) AddBridgeToken(ctx sdk.Context, bridgeDenom, baseDenom string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeDenomKey(bridgeDenom), []byte(baseDenom))
}

func (k Keeper) IteratorBridgeDenomWithContract(ctx sdk.Context, cb func(token *types.BridgeToken) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.BridgeDenomKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		if string(iter.Value()) == fxtypes.DefaultDenom {
			continue
		}

		bridgeToken := &types.BridgeToken{
			Denom: string(iter.Key()[len(types.BridgeDenomKey):]),
			Token: types.BridgeDenomToContract(k.moduleName, string(iter.Value())),
		}
		// cb returns true to stop early
		if cb(bridgeToken) {
			break
		}
	}
}

// Deprecated: do not use
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

func (k Keeper) HasToken(ctx context.Context, denom string) (bool, error) {
	ok := k.bankKeeper.HasDenomMetaData(ctx, denom)
	return ok, nil
}

func (k Keeper) GetBridgeDenom(ctx context.Context, denom string) ([]string, error) {
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

func (k Keeper) GetBaseDenom(ctx context.Context, alias string) (string, error) {
	var baseDenom string
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		if len(metadata.DenomUnits) == 0 {
			return false
		}
		for _, a := range metadata.DenomUnits[0].Aliases {
			if a == alias {
				baseDenom = metadata.Base
				return true
			}
		}
		return false
	})
	if baseDenom == "" {
		return baseDenom, fmt.Errorf("alias %s not found", alias)
	}
	return baseDenom, nil
}

func (k Keeper) GetAllTokens(ctx context.Context) ([]string, error) {
	panic("not implemented") // TODO implement me
}

func (k Keeper) UpdateBridgeDenom(ctx context.Context, denom string, bridgeDenoms ...string) error {
	metadata, ok := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !ok {
		return fmt.Errorf("denom %s not found", denom)
	}
	if len(metadata.DenomUnits) == 0 {
		return fmt.Errorf("denom %s denom units is empty", denom)
	}
	metadata.DenomUnits[0].Aliases = bridgeDenoms
	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	return nil
}

func (k Keeper) SetToken(ctx context.Context, name, symbol string, decimals uint32, bridgeDenoms ...string) error {
	if ok := k.bankKeeper.HasDenomMetaData(ctx, strings.ToLower(name)); ok {
		return fmt.Errorf("denom %s already exist", name)
	}
	metadata := fxtypes.GetCrossChainMetadataManyToOne(name, symbol, decimals, bridgeDenoms...)
	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	return nil
}
