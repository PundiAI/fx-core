package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestCrosschainABI(t *testing.T) {
	crosschain := precompile.NewCrosschainMethod(nil)

	require.Len(t, crosschain.Method.Inputs, 6)
	require.Len(t, crosschain.Method.Outputs, 1)

	require.Len(t, crosschain.Event.Inputs, 8)
}
