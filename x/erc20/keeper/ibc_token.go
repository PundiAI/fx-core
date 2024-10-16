package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) GetIBCToken(ctx context.Context, baseDenom, channel string) (types.IBCToken, error) {
	return k.IBCToken.Get(ctx, collections.Join(baseDenom, channel))
}

func (k Keeper) AddIBCToken(ctx context.Context, baseDenom, channel, ibcDenom string) error {
	key := collections.Join(baseDenom, channel)
	has, err := k.IBCToken.Has(ctx, key)
	if err != nil {
		return err
	}
	if has {
		return types.ErrExists.Wrapf("channel: %s base denom: %s", channel, baseDenom)
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
