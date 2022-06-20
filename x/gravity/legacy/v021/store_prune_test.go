package v021_test

import (
	"testing"

	v021 "github.com/functionx/fx-core/x/gravity/legacy/v021"
	"github.com/functionx/fx-core/x/gravity/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMigratePruneKey_IbcSequenceMigration(t *testing.T) {
	gravityKey := sdk.NewKVStoreKey("gravity")
	ctx := testutil.DefaultContext(gravityKey, sdk.NewTransientStoreKey("transient_test"))
	store := ctx.KVStore(gravityKey)

	sourcePort, sourceChannel := "a", "b"
	for i := 1; i < 100; i++ {
		store.Set(types.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i)), sdk.Uint64ToBigEndian(uint64(i)))
	}

	// migrate before check key exists
	for i := 1; i < 100; i++ {
		require.True(t, store.Has(types.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i))))
	}
	v021.MigratePruneKey(ctx.KVStore(gravityKey), types.IbcSequenceHeightKey)

	// migrate after check key not exists
	for i := 1; i < 100; i++ {
		require.False(t, store.Has(types.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i))))
	}
}
