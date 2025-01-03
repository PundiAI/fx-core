package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) HasToken(ctx context.Context, denom string) (bool, error) {
	return k.DenomIndex.Has(ctx, denom)
}

func (k Keeper) GetBaseDenom(ctx context.Context, token string) (string, error) {
	return k.DenomIndex.Get(ctx, token)
}

func (k Keeper) GetBridgeToken(ctx context.Context, chainName, baseDenom string) (types.BridgeToken, error) {
	return k.BridgeToken.Get(ctx, collections.Join(chainName, baseDenom))
}

func (k Keeper) GetBridgeTokens(ctx context.Context, chainName string) ([]types.BridgeToken, error) {
	rng := collections.NewPrefixedPairRange[string, string](chainName)
	iter, err := k.BridgeToken.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}

	kvs, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	tokens := make([]types.BridgeToken, 0, len(kvs))
	for _, kv := range kvs {
		tokens = append(tokens, kv.Value)
	}

	return tokens, nil
}

func (k Keeper) AddBridgeToken(ctx context.Context, baseDenom, chainName, contract string, isNative bool) error {
	key := collections.Join(chainName, baseDenom)
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
