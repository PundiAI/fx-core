package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) HasToken(ctx context.Context, index string) (bool, error) {
	return k.DenomIndex.Has(ctx, index)
}

func (k Keeper) GetBaseDenom(ctx context.Context, index string) (string, error) {
	return k.DenomIndex.Get(ctx, index)
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
		return sdkerrors.ErrInvalidRequest.Wrapf("bridge token %s already exists", contract)
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
	return k.DenomIndex.Set(ctx, types.NewBridgeDenom(chainName, contract), baseDenom)
}
