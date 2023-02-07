package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func PruneStore(cdc codec.BinaryCodec, kvStore sdk.KVStore) {
	pruneEvidence(kvStore)
	pruneBatchConfirmKey(cdc, kvStore)
	pruneOracleSetConfirmKey(cdc, kvStore)
	kvStore.Delete(types.LastProposalBlockHeight) // nolint:staticcheck
}

func pruneBatchConfirmKey(cdc codec.BinaryCodec, store sdk.KVStore) {
	iterator := sdk.KVStorePrefixIterator(store, types.OutgoingTxBatchKey)
	defer iterator.Close()

	var batchs []types.OutgoingTxBatch
	for ; iterator.Valid(); iterator.Next() {
		var batch types.OutgoingTxBatch
		cdc.MustUnmarshal(iterator.Value(), &batch)
		batchs = append(batchs, batch)
	}

	iter := sdk.KVStorePrefixIterator(store, types.BatchConfirmKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgConfirmBatch)
		cdc.MustUnmarshal(iter.Value(), confirm)
		found := false
		for _, batch := range batchs {
			if batch.BatchNonce == confirm.Nonce && batch.TokenContract == confirm.TokenContract {
				found = true
				break
			}
		}
		if !found {
			store.Delete(iter.Key())
		}
	}
}

func pruneOracleSetConfirmKey(cdc codec.BinaryCodec, store sdk.KVStore) {
	lastSlashedOracleSetNonce := sdk.BigEndianToUint64(store.Get(types.LastSlashedOracleSetNonce))

	iter := sdk.KVStorePrefixIterator(store, types.OracleSetConfirmKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		confirm := new(types.MsgOracleSetConfirm)
		cdc.MustUnmarshal(iter.Value(), confirm)
		if lastSlashedOracleSetNonce >= confirm.Nonce {
			store.Delete(iter.Key())
		}
	}
}

func pruneEvidence(store sdk.KVStore) {
	iter := sdk.KVStorePrefixIterator(store, types.PastExternalSignatureCheckpointKey) // nolint:staticcheck
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
