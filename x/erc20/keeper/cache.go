package keeper

import (
	"context"

	"cosmossdk.io/collections"
)

func (k Keeper) HasCache(ctx context.Context, key string) (bool, error) {
	return k.Cache.Has(ctx, key)
}

func (k Keeper) SetCache(ctx context.Context, key string) error {
	return k.Cache.Set(ctx, key, collections.NoValue{})
}

func (k Keeper) DeleteCache(ctx context.Context, key string) error {
	return k.Cache.Remove(ctx, key)
}
