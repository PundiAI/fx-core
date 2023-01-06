package keeper

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (k Keeper) RefundAfter(ctx sdk.Context, port, channel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error {
	// check exist
	if !k.DeleteIBCTransferHash(ctx, port, channel, sequence) {
		return nil
	}
	return k.TransferAfter(ctx, sender.String(), common.BytesToAddress(sender.Bytes()).String(), amount, sdk.NewCoin(amount.Denom, sdk.ZeroInt()))
}

func (k Keeper) AckAfter(ctx sdk.Context, port, channel string, sequence uint64) error {
	k.DeleteIBCTransferHash(ctx, port, channel, sequence)
	return nil
}

func (k Keeper) SetIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64, hash common.Hash) {
	// todo delete unused port and hash
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetIBCTransferKey(port, channel, sequence), hash.Bytes())
}

func (k Keeper) DeleteIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) bool {
	if !k.hasIBCTransferHash(ctx, port, channel, sequence) {
		return false
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetIBCTransferKey(port, channel, sequence))
	return true
}

func (k Keeper) hasIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetIBCTransferKey(port, channel, sequence))
}

func (k Keeper) IterateIBCTransferHash(ctx sdk.Context, cb func(port, channel string, sequence uint64) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixIBCTransfer)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := bytes.TrimPrefix(iter.Key(), types.KeyPrefixIBCTransfer)
		split := strings.Split(string(key), "/")
		if len(split) != 3 {
			panic(fmt.Sprintf("invalid key: %s", string(key)))
		}
		port := split[0]
		channel := split[1]
		sequence, err := strconv.ParseUint(split[2], 10, 64)
		if err != nil {
			panic(fmt.Sprintf("parse sequence %s error %s", split[2], err.Error()))
		}
		if cb(port, channel, sequence) {
			return
		}
	}
}
