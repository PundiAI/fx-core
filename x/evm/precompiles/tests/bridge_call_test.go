package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/x/evm/precompiles/crosschain"
)

func TestBridgeCallABI(t *testing.T) {
	crosschainABI := crosschain.GetABI()

	method := crosschainABI.Methods[crosschain.BridgeCallMethodName]
	require.Equal(t, method, crosschain.BridgeCallMethod)
	require.Equal(t, 8, len(crosschain.BridgeCallMethod.Inputs))
	require.Equal(t, 1, len(crosschain.BridgeCallMethod.Outputs))
}
