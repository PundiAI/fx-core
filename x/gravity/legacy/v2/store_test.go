package v2_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	v2 "github.com/functionx/fx-core/x/gravity/legacy/v2"
	"github.com/functionx/fx-core/x/gravity/types"
)

func TestIbcSequenceMigration(t *testing.T) {
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
	err := v2.MigrateStore(ctx, gravityKey)
	require.NoError(t, err)

	// migrate after check key not exists
	for i := 1; i < 100; i++ {
		require.False(t, store.Has(types.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i))))
	}
}
