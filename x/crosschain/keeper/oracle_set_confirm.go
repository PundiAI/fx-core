package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

// --- ORACLE SET CONFIRMS --- //

// GetOracleSetConfirm returns a oracleSet confirmation by a nonce and external address
func (k Keeper) GetOracleSetConfirm(ctx sdk.Context, nonce uint64, oracleAddr sdk.AccAddress) *types.MsgOracleSetConfirm {
	store := ctx.KVStore(k.storeKey)
	entity := store.Get(types.GetOracleSetConfirmKey(nonce, oracleAddr))
	if entity == nil {
		return nil
	}
	confirm := types.MsgOracleSetConfirm{}
	k.cdc.MustUnmarshal(entity, &confirm)
	return &confirm
}

// SetOracleSetConfirm sets a oracleSet confirmation
func (k Keeper) SetOracleSetConfirm(ctx sdk.Context, oracleAddr sdk.AccAddress, oracleSetConfirm *types.MsgOracleSetConfirm) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOracleSetConfirmKey(oracleSetConfirm.Nonce, oracleAddr)
	store.Set(key, k.cdc.MustMarshal(oracleSetConfirm))
}

// IterateOracleSetConfirmByNonce iterates through all oracleSet confirms by nonce
func (k Keeper) IterateOracleSetConfirmByNonce(ctx sdk.Context, nonce uint64, cb func(*types.MsgOracleSetConfirm) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.GetOracleSetConfirmKey(nonce, sdk.AccAddress{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgOracleSetConfirm)
		k.cdc.MustUnmarshal(iter.Value(), confirm)
		// cb returns true to stop early
		if cb(confirm) {
			break
		}
	}
}

func (k Keeper) DeleteOracleSetConfirm(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.GetOracleSetConfirmKey(nonce, sdk.AccAddress{}))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// --- LAST SLASHED ORACLE SET NONCE --- //

// SetLastSlashedOracleSetNonce sets the latest slashed oracleSet nonce
func (k Keeper) SetLastSlashedOracleSetNonce(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastSlashedOracleSetNonce, sdk.Uint64ToBigEndian(nonce))
}

// GetLastSlashedOracleSetNonce returns the latest slashed oracleSet nonce
func (k Keeper) GetLastSlashedOracleSetNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(types.LastSlashedOracleSetNonce))
}
