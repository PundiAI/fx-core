package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func (k Keeper) RefundAfter(ctx sdk.Context, channel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error {
	// check exist
	if !k.DeleteIBCTransferRelation(ctx, channel, sequence) {
		return nil
	}
	cacheCtx, commit := ctx.CacheContext()
	if err := k.TransferAfter(cacheCtx, sender, common.BytesToAddress(sender.Bytes()).String(),
		amount, sdk.NewCoin(amount.Denom, sdkmath.ZeroInt()), false); err != nil {
		return err
	}
	commit()
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

func (k Keeper) HookOutgoingRefund(ctx sdk.Context, moduleName string, txID uint64, sender sdk.AccAddress, totalCoin sdk.Coin) error {
	if _, err := k.ConvertCoin(sdk.WrapSDKContext(ctx), &types.MsgConvertCoin{
		Coin:     totalCoin,
		Receiver: common.BytesToAddress(sender.Bytes()).String(),
		Sender:   sender.String(),
	}); err != nil {
		return err
	}

	k.DeleteOutgoingTransferRelation(ctx, moduleName, txID)
	return nil
}

func (k Keeper) SetOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingTransferKey(moduleName, txID), []byte{})
}

func (k Keeper) DeleteOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingTransferKey(moduleName, txID))
}

func (k Keeper) HasOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOutgoingTransferKey(moduleName, txID))
}
