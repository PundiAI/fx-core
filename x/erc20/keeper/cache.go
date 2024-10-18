package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
)

func (k Keeper) HasCache(ctx context.Context, key string) (bool, error) {
	return k.Cache.Has(ctx, key)
}

func (k Keeper) SetCache(ctx context.Context, key string, amount sdkmath.Int) error {
	return k.Cache.Set(ctx, key, amount)
}

func (k Keeper) DeleteCache(ctx context.Context, key string) error {
	return k.Cache.Remove(ctx, key)
}