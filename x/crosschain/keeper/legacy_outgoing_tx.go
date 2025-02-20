package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

// --- OUTGOING TX POOL --- //

// Deprecated: OutgoingTxPoolKey indexes the last nonce for the outgoing tx pool
var OutgoingTxPoolKey = []byte{0x18}

// IterateOutgoingTxPool iterates through all outgoing txs in the pool
func (k Keeper) IterateOutgoingTxPool(ctx sdk.Context, cb func(types.OutgoingTransferTx) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, OutgoingTxPoolKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		tx := new(types.OutgoingTransferTx)
		k.cdc.MustUnmarshal(iter.Value(), tx)
		if cb(*tx) {
			break
		}
	}
}

// DeleteOutgoingTxPool deletes an outgoing tx from the pool
func (k Keeper) DeleteOutgoingTxPool(ctx sdk.Context, fee types.ERC20Token, txID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetOutgoingTxPoolKey(fee, txID))
}

// GetOutgoingTxPoolKey returns the following key format
func GetOutgoingTxPoolKey(fee types.ERC20Token, id uint64) []byte {
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	return append(OutgoingTxPoolKey, append([]byte(fee.Contract), append(amount, sdk.Uint64ToBigEndian(id)...)...)...)
}
