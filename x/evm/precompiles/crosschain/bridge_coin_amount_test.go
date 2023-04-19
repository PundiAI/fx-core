package crosschain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v4/x/evm/precompiles/crosschain"
)

func TestBridgeCoinAmountABI(t *testing.T) {
	crosschainABI := crosschain.GetABI()

	method := crosschainABI.Methods[crosschain.BridgeCoinAmountMethodName]
	require.Equal(t, method, crosschain.BridgeCoinAmountMethod)
	require.Equal(t, 2, len(crosschain.BridgeCoinAmountMethod.Inputs))
	require.Equal(t, 1, len(crosschain.BridgeCoinAmountMethod.Outputs))
}
