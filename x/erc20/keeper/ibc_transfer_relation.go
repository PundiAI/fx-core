package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (k Keeper) RefundAfter(ctx sdk.Context, channel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error {
	// check exist
	if !k.DeleteIBCTransferRelation(ctx, channel, sequence) {
		return nil
	}
	cacheCtx, commit := ctx.CacheContext()
	if err := k.TransferAfter(cacheCtx, sender.String(), common.BytesToAddress(sender.Bytes()).String(),
		amount, sdk.NewCoin(amount.Denom, sdk.ZeroInt())); err != nil {
		return err
	}
	commit()
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
	return nil
}

func (k Keeper) AckAfter(ctx sdk.Context, channel string, sequence uint64) error {
	k.DeleteIBCTransferRelation(ctx, channel, sequence)
	return nil
}

func (k Keeper) SetIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIBCTransferKey(channel, sequence), []byte{})
}

func (k Keeper) DeleteIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64) bool {
	if !k.hasIBCTransferRelation(ctx, channel, sequence) {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIBCTransferKey(channel, sequence))
	return true
}

func (k Keeper) hasIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetIBCTransferKey(channel, sequence))
}
