package store_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func Benchmark_GetStoreValue1(b *testing.B) {
	storeKey := sdk.NewKVStoreKey("test")
	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)
	assert.NoError(b, ms.LoadLatestVersion())
	store := ms.GetKVStore(storeKey)

	count := 10000
	for i := 0; i < count; i++ {
		key := append([]byte{0x1}, sdk.Uint64ToBigEndian(uint64(i))...)
		store.Set(key, []byte{1, 2, 3})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var data [][]byte
		iter := sdk.KVStorePrefixIterator(store, []byte{0x1})
		for ; iter.Valid(); iter.Next() {
			iter.Value()
			data = append(data, iter.Value())
		}
		assert.Equal(b, count, len(data))
	}
}

func Benchmark_GetStoreValue2(b *testing.B) {
	storeKey := sdk.NewKVStoreKey("test")
	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)
	assert.NoError(b, ms.LoadLatestVersion())
	store := ms.GetKVStore(storeKey)

	count := 10000
	for i := 0; i < count; i++ {
		key := append([]byte{0x1}, sdk.Uint64ToBigEndian(uint64(i))...)
		store.Set(key, []byte{1, 2, 3})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var data [][]byte
		for j := 0; j < count; j++ {
			key := append([]byte{0x1}, sdk.Uint64ToBigEndian(uint64(i))...)
			store.Get(key)
			data = append(data, store.Get(key))
		}
		assert.Equal(b, count, len(data))
	}
}
