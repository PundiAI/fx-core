package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

// SetPastExternalSignatureCheckpoint puts the checkpoint of a oracle set, batch, or logic call into a set
// in order to prove later that it existed at one point.
func (k Keeper) SetPastExternalSignatureCheckpoint(ctx sdk.Context, checkpoint []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPastExternalSignatureCheckpointKey(checkpoint), []byte{0x1})
}

// GetPastExternalSignatureCheckpoint tells you whether a given checkpoint has ever existed
func (k Keeper) GetPastExternalSignatureCheckpoint(ctx sdk.Context, checkpoint []byte) (found bool) {
	store := ctx.KVStore(k.storeKey)
	if bytes.Equal(store.Get(types.GetPastExternalSignatureCheckpointKey(checkpoint)), []byte{0x1}) {
		return true
	} else {
		return false
	}
}
