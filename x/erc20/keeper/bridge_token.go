package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) HasToken(ctx context.Context, denom string) (bool, error) {
	return k.DenomIndex.Has(ctx, denom)
}

func (k Keeper) GetBridgeToken(ctx context.Context, baseDenom, chainName string) (types.BridgeToken, error) {
	return k.BridgeToken.Get(ctx, collections.Join(baseDenom, chainName))
}

func (k Keeper) GetBaseDenom(ctx context.Context, token string) (string, error) {
	return k.DenomIndex.Get(ctx, token)
}

func (k Keeper) AddBridgeToken(ctx context.Context, baseDenom string, chainName string, contract string, isNative bool) error {
	key := collections.Join(baseDenom, chainName)
	has, err := k.BridgeToken.Has(ctx, key)
	if err != nil {
		return err
	}
	if has {
		return types.ErrExists.Wrapf("%s base denom: %s", chainName, baseDenom)
	}
	bridgeToken := types.BridgeToken{
		IsNative:  isNative,
		Denom:     baseDenom,
		Contract:  contract,
		ChainName: chainName,
	}
	if err = k.BridgeToken.Set(ctx, key, bridgeToken); err != nil {
		return err
	}
	return k.DenomIndex.Set(ctx, bridgeToken.BridgeDenom(), baseDenom)
}
