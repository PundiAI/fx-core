package types

import (
	"fmt"
	"strings"
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
			denom:     fmt.Sprintf("%s/%s/%s", types.PortID, PundixChannel, MainnetPundixUnWrapDenom),
			expected:  true,
			wrapDenom: PundixWrapDenom,
		},
		{
			name:      "testnet eth0xpundix to pundix",
			chainId:   TestnetChainId,
			channel:   PundixChannel,
			denom:     fmt.Sprintf("%s/%s/%s", types.PortID, PundixChannel, TestnetPundixUnWrapDenom),
			expected:  true,
			wrapDenom: PundixWrapDenom,
		},
		{
			name:      "mainnet FX to DefaultDenom",
			chainId:   MainnetChainId,
			channel:   "channel-2716",
			denom:     "transfer/channel-2716/FX",
			expected:  true,
			wrapDenom: DefaultDenom,
		},
		{
			name:      "testnet FX to DefaultDenom",
			chainId:   TestnetChainId,
			channel:   "channel-2716",
			denom:     "transfer/channel-2716/FX",
			expected:  true,
			wrapDenom: DefaultDenom,
		},
		{
			name:     "no need wrap",
			chainId:  MainnetChainId,
			channel:  PundixChannel,
			denom:    strings.ToLower(tmrand.Str(6)),
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			needWrap, wrapDenom, packetDenom := OnRecvDenomNeedWrap(tc.chainId, types.PortID, tc.channel, tc.denom)
			require.Equal(t, tc.expected, needWrap)
			if !needWrap {
				return
			}
			require.Equal(t, tc.wrapDenom, wrapDenom)
			require.Equal(t, fmt.Sprintf("%s/%s/%s", types.PortID, tc.channel, wrapDenom), packetDenom)
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
