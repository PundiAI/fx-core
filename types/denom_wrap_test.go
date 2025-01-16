package types

import (
	"strings"
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/require"
)

func Test_OnRecvDenomNeedWrap(t *testing.T) {
	testCases := []struct {
		name      string
		chainId   string
		channel   string
		denom     string
		expected  bool
		wrapDenom string
	}{
		{
			name:      "mainnet ethPundix to pundix",
			chainId:   MainnetChainId,
			channel:   PundixChannel,
			denom:     MainnetPundixUnWrapDenom,
			expected:  true,
			wrapDenom: PundixWrapDenom,
		},
		{
			name:      "mainnet osmosis FX to DefaultDenom",
			chainId:   MainnetChainId,
			channel:   MainnetOsmosisChannel,
			denom:     IBCFXDenom,
			expected:  true,
			wrapDenom: DefaultDenom,
		},
		{
			name:      "mainnet pundix FX to DefaultDenom",
			chainId:   MainnetChainId,
			channel:   PundixChannel,
			denom:     IBCFXDenom,
			expected:  true,
			wrapDenom: DefaultDenom,
		},
		{
			name:      "testnet eth0xpundix to pundix",
			chainId:   TestnetChainId,
			channel:   PundixChannel,
			denom:     TestnetPundixUnWrapDenom,
			expected:  true,
			wrapDenom: PundixWrapDenom,
		},
		{
			name:      "testnet osmosis FX to DefaultDenom",
			chainId:   TestnetChainId,
			channel:   TestnetOsmosisChannel,
			denom:     IBCFXDenom,
			expected:  true,
			wrapDenom: DefaultDenom,
		},
		{
			name:      "testnet pundix FX to DefaultDenom",
			chainId:   TestnetChainId,
			channel:   PundixChannel,
			denom:     IBCFXDenom,
			expected:  true,
			wrapDenom: DefaultDenom,
		},
		{
			name:      "no need wrap",
			chainId:   MainnetChainId,
			channel:   PundixChannel,
			denom:     strings.ToLower(tmrand.Str(6)),
			expected:  false,
			wrapDenom: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			needWrap, wrapDenom := OnRecvDenomNeedWrap(tc.chainId, tc.channel, tc.denom)
			require.Equal(t, tc.expected, needWrap)
			require.Equal(t, tc.wrapDenom, wrapDenom)
		})
	}
}

func Test_SendPacketDenomNeedWrap(t *testing.T) {
	testCases := []struct {
		name      string
		chainId   string
		channel   string
		denom     string
		expected  bool
		wrapDenom string
	}{
		{
			name:      "mainnet pundix to ethPundix",
			chainId:   MainnetChainId,
			channel:   PundixChannel,
			denom:     PundixWrapDenom,
			expected:  true,
			wrapDenom: MainnetPundixUnWrapDenom,
		},
		{
			name:      "testnet pundix to ethPundix",
			chainId:   TestnetChainId,
			channel:   PundixChannel,
			denom:     PundixWrapDenom,
			expected:  true,
			wrapDenom: TestnetPundixUnWrapDenom,
		},
		{
			name:      "no need wrap",
			chainId:   MainnetChainId,
			channel:   PundixChannel,
			denom:     strings.ToLower(tmrand.Str(6)),
			expected:  false,
			wrapDenom: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			needWrap, wrapDenom := SendPacketDenomNeedWrap(tc.chainId, tc.channel, tc.denom)
			require.Equal(t, tc.expected, needWrap)
			require.Equal(t, tc.wrapDenom, wrapDenom)
		})
	}
}
