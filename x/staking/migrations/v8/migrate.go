package v8

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// Deprecated: do not use, remove in v8
	ValidatorOperatorKey = []byte{0x91}
	// Deprecated: do not use, remove in v8
	ConsensusPubKey = []byte{0x92}
	// Deprecated: do not use, remove in v8
	ConsensusProcessKey = []byte{0x93}
)

func GetRemovedValidatorStoreKeys() [][]byte {
	return [][]byte{ValidatorOperatorKey, ConsensusPubKey, ConsensusProcessKey}
}

func DeleteMigrationValidatorStore(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
) {
	store := ctx.KVStore(storeKey)
	removeKeys := GetRemovedValidatorStoreKeys()
	for _, key := range removeKeys {
		iterator := storetypes.KVStorePrefixIterator(store, key)
		for ; iterator.Valid(); iterator.Next() {
			store.Delete(iterator.Key())
		}
	}
}
