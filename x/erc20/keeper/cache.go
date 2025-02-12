package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) HasCache(ctx context.Context, key string) (bool, error) {
	return k.Cache.Has(ctx, key)
}

func (k Keeper) SetCache(ctx context.Context, key string, amount sdkmath.Int) error {
	found, err := k.HasCache(ctx, key)
	if err != nil {
		return err
	}
	if found {
		return sdkerrors.ErrInvalidRequest.Wrapf("cache %s already exists", key)
	}
	return k.Cache.Set(ctx, key, amount)
}

func (k Keeper) DeleteCache(ctx context.Context, key string) error {
	return k.Cache.Remove(ctx, key)
}

func (k Keeper) GetCache(ctx context.Context, key string) (sdkmath.Int, error) {
	return k.Cache.Get(ctx, key)
}

func (k Keeper) ReSetCache(ctx context.Context, oldKey, newKey string) error {
	amount, err := k.Cache.Get(ctx, oldKey)
	if err == nil {
		if err = k.Cache.Set(ctx, newKey, amount); err != nil {
			return err
		}
		return k.Cache.Remove(ctx, oldKey)
	}

	if !errors.IsOf(err, collections.ErrNotFound) {
		return err
	}

	return nil
}
