package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func PruneBatchConfirmKey(cdc codec.BinaryCodec, store sdk.KVStore) {
	lastOutgoingBatchNonce := sdk.BigEndianToUint64(store.Get(types.KeyLastOutgoingBatchID))

	iter := sdk.KVStorePrefixIterator(store, types.BatchConfirmKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgConfirmBatch)
		cdc.MustUnmarshal(iter.Value(), confirm)
		if lastOutgoingBatchNonce > confirm.Nonce {
			store.Delete(iter.Key())
		}
	}
}

func PruneOracleSetConfirmKey(cdc codec.BinaryCodec, store sdk.KVStore) {
	latestOracleSetNonce := sdk.BigEndianToUint64(store.Get(types.LatestOracleSetNonce))

	iter := sdk.KVStorePrefixIterator(store, types.OracleSetConfirmKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgOracleSetConfirm)
		cdc.MustUnmarshal(iter.Value(), confirm)
		if latestOracleSetNonce > confirm.Nonce {
			store.Delete(iter.Key())
		}
	}
}

func PruneEvidence(store sdk.KVStore) {
	iter := sdk.KVStorePrefixIterator(store, types.PastExternalSignatureCheckpointKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
