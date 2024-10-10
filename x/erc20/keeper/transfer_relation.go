package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) IbcRefund(ctx sdk.Context, channel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error {
	// check exist
	if !k.DeleteIBCTransferRelation(ctx, channel, sequence) {
		return nil
	}
	_, err := k.ConvertCoin(ctx, &types.MsgConvertCoin{
		Coin:     amount,
		Receiver: common.BytesToAddress(sender.Bytes()).String(),
		Sender:   sender.String(),
	})
	return err
}

func (k Keeper) SetIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIBCTransferKey(channel, sequence), []byte{})
}

func (k Keeper) DeleteIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetIBCTransferKey(channel, sequence)
	if !store.Has(key) {
		return false
	}
	store.Delete(key)
	return true
}

func (k Keeper) HookOutgoingRefund(ctx sdk.Context, moduleName string, txID uint64, sender sdk.AccAddress, totalCoin sdk.Coin) error {
	if _, err := k.ConvertCoin(ctx, &types.MsgConvertCoin{
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
