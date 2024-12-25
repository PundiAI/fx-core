package crosschain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
)

func TestCrosschainABI(t *testing.T) {
	crosschainABI := crosschain.NewCrosschainABI()

	require.Len(t, crosschainABI.Method.Inputs, 6)
	require.Len(t, crosschainABI.Method.Outputs, 1)

	require.Len(t, crosschainABI.Event.Inputs, 8)
}
