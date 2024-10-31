package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestCrosschainABI(t *testing.T) {
	crosschainABI := precompile.NewCrosschainABI()

	require.Len(t, crosschainABI.Method.Inputs, 6)
	require.Len(t, crosschainABI.Method.Outputs, 1)

	require.Len(t, crosschainABI.Event.Inputs, 8)
}
