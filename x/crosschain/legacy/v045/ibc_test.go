package v045_test

import (
	"testing"

	v042 "github.com/functionx/fx-core/x/crosschain/legacy/v042"
	v045 "github.com/functionx/fx-core/x/crosschain/legacy/v045"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestIbcSequenceMigration(t *testing.T) {
	gravityKey := sdk.NewKVStoreKey("gravity")
	ctx := testutil.DefaultContext(gravityKey, sdk.NewTransientStoreKey("transient_test"))
	store := ctx.KVStore(gravityKey)

	sourcePort, sourceChannel := "a", "b"
	for i := 1; i < 100; i++ {
		store.Set(v042.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i)), sdk.Uint64ToBigEndian(uint64(i)))
	}

	// migrate before check key exists
	for i := 1; i < 100; i++ {
		require.True(t, store.Has(v042.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i))))
	}
	require.NoError(t, v045.MigratePruneIbcSequenceKey(ctx, gravityKey))

	// migrate after check key not exists
	for i := 1; i < 100; i++ {
		require.False(t, store.Has(v042.GetIbcSequenceHeightKey(sourcePort, sourceChannel, uint64(i))))
	}
}
