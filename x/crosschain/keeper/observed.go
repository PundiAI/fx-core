package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

// --- OBSERVED EVENT NONCE --- //

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(types.LastObservedEventNonceKey))
}

// SetLastObservedEventNonce sets the latest observed event nonce
func (k Keeper) SetLastObservedEventNonce(ctx sdk.Context, eventNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedEventNonceKey, sdk.Uint64ToBigEndian(eventNonce))
}

// --- OBSERVED BLOCK HEIGHT --- //

// GetLastObservedBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedBlockHeight(ctx sdk.Context) types.LastObservedBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedBlockHeightKey)
	if len(bytes) == 0 {
		return types.LastObservedBlockHeight{
			ExternalBlockHeight: 0,
			BlockHeight:         0,
		}
	}
	height := types.LastObservedBlockHeight{}
	k.cdc.MustUnmarshal(bytes, &height)
	return height
}

// SetLastObservedBlockHeight sets the block height in the store.
func (k Keeper) SetLastObservedBlockHeight(ctx sdk.Context, externalBlockHeight, blockHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	height := types.LastObservedBlockHeight{
		ExternalBlockHeight: externalBlockHeight,
		BlockHeight:         blockHeight,
	}
	store.Set(types.LastObservedBlockHeightKey, k.cdc.MustMarshal(&height))
}

// --- LAST OBSERVED ORACLE SET --- //

// GetLastObservedOracleSet retrieves the last observed oracle set from the store
// WARNING: This value is not an up to date oracle set on Ethereum, it is a oracle set
// that AT ONE POINT was the one in the bridge on Ethereum. If you assume that it's up
// to date you may break the bridge
func (k Keeper) GetLastObservedOracleSet(ctx sdk.Context) *types.OracleSet {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.LastObservedOracleSetKey)

	if len(bytes) == 0 {
		return nil
	}
	oracleSet := types.OracleSet{}
	k.cdc.MustUnmarshal(bytes, &oracleSet)
	return &oracleSet
}

// SetLastObservedOracleSet updates the last observed oracle set in the store
func (k Keeper) SetLastObservedOracleSet(ctx sdk.Context, oracleSet *types.OracleSet) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastObservedOracleSetKey, k.cdc.MustMarshal(oracleSet))
}
