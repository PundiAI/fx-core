package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestCrossChainABI(t *testing.T) {
	crossChain := precompile.NewCrossChainMethod(nil)

	require.Equal(t, 6, len(crossChain.Method.Inputs))
	require.Equal(t, 1, len(crossChain.Method.Outputs))

	require.Equal(t, 8, len(crossChain.Event.Inputs))
}
