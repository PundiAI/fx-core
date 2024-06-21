package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func TestBridgeCallABI(t *testing.T) {
	crosschainABI := types.GetABI()

	method := crosschainABI.Methods[types.BridgeCallMethodName]
	require.Equal(t, method, types.BridgeCallMethod)
	require.Equal(t, 8, len(types.BridgeCallMethod.Inputs))
	require.Equal(t, 1, len(types.BridgeCallMethod.Outputs))
}
