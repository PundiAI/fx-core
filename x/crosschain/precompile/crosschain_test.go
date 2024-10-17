package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestCrosschainABI(t *testing.T) {
	crosschain := precompile.NewCrosschainMethod(nil)

	require.Equal(t, 6, len(crosschain.Method.Inputs))
	require.Equal(t, 1, len(crosschain.Method.Outputs))

	require.Equal(t, 8, len(crosschain.Event.Inputs))
}
