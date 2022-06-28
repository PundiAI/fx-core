package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetChainId_mainnet(t *testing.T) {
	require.Equal(t, ChainId(), mainnetChainId)

	SetChainId(mainnetChainId)
	require.Equal(t, ChainId(), mainnetChainId)

	SetChainId(testnetChainId)
	require.Equal(t, ChainId(), mainnetChainId)
}

func TestSetChainId_testnet(t *testing.T) {
	require.Equal(t, ChainId(), mainnetChainId)

	SetChainId(testnetChainId)
	require.Equal(t, ChainId(), testnetChainId)

	SetChainId(mainnetChainId)
	require.Equal(t, ChainId(), testnetChainId)
}

func TestSetChainId_invalid(t *testing.T) {
	require.Equal(t, ChainId(), mainnetChainId)
	require.Panics(t, func() {
		SetChainId("hello")
	})
}
