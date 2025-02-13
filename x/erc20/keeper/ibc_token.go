package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) GetIBCToken(ctx context.Context, baseDenom, channel string) (types.IBCToken, error) {
	return k.IBCToken.Get(ctx, collections.Join(baseDenom, channel))
}

func (k Keeper) GetBaseIBCTokens(ctx context.Context, baseDenom string) ([]types.IBCToken, error) {
	rng := collections.NewPrefixedPairRange[string, string](baseDenom)
	iter, err := k.IBCToken.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}
	kvs, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	tokens := make([]types.IBCToken, 0, len(kvs))
	for _, kv := range kvs {
		tokens = append(tokens, kv.Value)
	}
	return tokens, nil
}

func (k Keeper) AddIBCToken(ctx context.Context, baseDenom, channel, ibcDenom string) error {
	key := collections.Join(baseDenom, channel)
	has, err := k.IBCToken.Has(ctx, key)
	if err != nil {
		return err
	}
	if has {
		return sdkerrors.ErrInvalidRequest.Wrapf("ibc token %s already exists", ibcDenom)
	}
	ibcToken := types.IBCToken{
		Channel:  channel,
		IbcDenom: ibcDenom,
	}
	if err = k.IBCToken.Set(ctx, key, ibcToken); err != nil {
		return err
	}
	return k.DenomIndex.Set(ctx, ibcToken.IbcDenom, baseDenom)
}
