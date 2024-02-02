package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// SetPastExternalSignatureCheckpoint puts the checkpoint of a oracle set, batch into a set
// in order to prove later that it existed at one point.
func (k Keeper) SetPastExternalSignatureCheckpoint(ctx sdk.Context, checkpoint []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPastExternalSignatureCheckpointKey(uint64(ctx.BlockHeight()), checkpoint), []byte{0x1})
}

// IteratePastExternalSignatureCheckpoint iterates through all PastExternalSignatureCheckpoint in the half-open interval [start,end)
func (k Keeper) IteratePastExternalSignatureCheckpoint(ctx sdk.Context, start uint64, end uint64, cb func([]byte) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.PastExternalSignatureCheckpointKey, sdk.Uint64ToBigEndian(start)...) // nolint:staticcheck
	endKey := append(types.PastExternalSignatureCheckpointKey, sdk.Uint64ToBigEndian(end)...)     // nolint:staticcheck
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// cb returns true to stop early
		if cb(iter.Key()[9:]) {
			break
		}
	}
}
