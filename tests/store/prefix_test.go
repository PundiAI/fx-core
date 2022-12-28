package store_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func TestPrefixStore(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")

	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)
	assert.NoError(t, ms.LoadLatestVersion())

	store := ms.GetKVStore(storeKey)
	store.Set([]byte{1, 1}, []byte{1, 1})
	store.Set([]byte{1, 2}, []byte{2, 2})
	store.Set([]byte{1}, []byte{3, 3})

	newStore := prefix.NewStore(store, []byte{1})
	newStore.Set([]byte{4}, []byte{4, 4})
	iter := newStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		t.Log(iter.Key(), iter.Value())
		//newStore.Delete(iter.Key())
	}

	iterator := sdk.KVStorePrefixIterator(store, []byte{1})
	for ; iterator.Valid(); iterator.Next() {
		t.Log(iterator.Key(), iterator.Value())
	}
}
