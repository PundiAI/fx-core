package v8

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/app/upgrades/store"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

var (
	// Deprecated: DenomToTokenKey prefixes the index of asset denom to external token
	DenomToTokenKey = []byte{0x26}

	// Deprecated: TokenToDenomKey prefixes the index of assets external token to denom
	TokenToDenomKey = []byte{0x27}

	// Deprecated: remove data in upgrade
	BridgeCallFromMsgKey = []byte{0x51}
)

// Deprecated: GetTokenToDenomKey returns the following key format
func GetTokenToDenomKey(denom string) []byte {
	return append(TokenToDenomKey, []byte(denom)...)
}

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{
		DenomToTokenKey,
		TokenToDenomKey,
		BridgeCallFromMsgKey,
	}
}

func Migrate(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store.RemoveStoreKeys(ctx, storeKey, GetRemovedStoreKeys())
	migrateOutgoingTxBatchBlockKey(ctx, storeKey, cdc)
	return nil
}

func migrateOutgoingTxBatchBlockKey(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	kvStore := ctx.KVStore(storeKey)
	iter := storetypes.KVStorePrefixIterator(kvStore, types.OutgoingTxBatchBlockKey)
	defer iter.Close()

	batches := make([]*types.OutgoingTxBatch, 0, 100)
	for ; iter.Valid(); iter.Next() {
		batch := new(types.OutgoingTxBatch)
		cdc.MustUnmarshal(iter.Value(), batch)
		batches = append(batches, batch)
	}

	for _, batch := range batches {
		blockKey := types.GetOutgoingTxBatchBlockKey(batch.Block, batch.BatchNonce)
		kvStore.Set(blockKey, cdc.MustMarshal(batch))
	}
}
