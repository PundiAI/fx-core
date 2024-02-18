package ibcrouter_test

import (
	"fmt"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	_ "github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/x/ibc/ibcrouter"
	"github.com/functionx/fx-core/v7/x/ibc/ibcrouter/parser"
)

func TestParseIncomingTransferField(t *testing.T) {
	testCases := []struct {
		name                string
		input               string
		expForward          bool
		expThisChainAddress string
		expFinalDestination string
		expPort             string
		expChannel          string
		expPass             bool
	}{
		{
			name:       "error - no-forward error thisChainAddress",
			input:      "fx1av497q6kjky9j5v3z95668d57q9ha80pwe45qy",
			expPass:    true,
			expForward: false,
		},
		{
			name:    "error - no-forward empty thisChainAddress",
			input:   "",
			expPass: false,
		},
		{
			name:    "error - forward empty thisChainAddress",
			input:   "|transfer/channel-0:cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPass: false,
		},
		{
			name:    "error - forward empty destinationAddress",
			input:   "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4|transfer/channel-0:",
			expPass: false,
		},
		{
			name:    "error - forward thisChain address is 0x",
			input:   common.BytesToAddress(tmrand.Bytes(20)).Hex() + "|transfer/channel-0:",
			expPass: false,
		},
		{
			name:                "ok - no-forward",
			input:               "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPass:             true,
			expThisChainAddress: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
		},
		{
			name:                "ok - forward empty portID",
			input:               "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4|/channel-0:cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPass:             true,
			expThisChainAddress: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPort:             "",
			expChannel:          "channel-0",
			expFinalDestination: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expForward:          true,
		},
		{
			name:                "ok - forward empty channelID",
			input:               "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4|transfer/:cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPass:             true,
			expThisChainAddress: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPort:             "transfer",
			expChannel:          "",
			expFinalDestination: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expForward:          true,
		},
		{
			name:                "ok - forward",
			input:               "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4|transfer/channel-0:cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPass:             true,
			expThisChainAddress: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expPort:             "transfer",
			expChannel:          "channel-0",
			expFinalDestination: "cosmos1av497q6kjky9j5v3z95668d57q9ha80p5fypd4",
			expForward:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pt, err := parser.ParseReceiverData(tc.input)
			if tc.expPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				return
			}

			if !tc.expForward {
				require.False(t, pt.ShouldForward)
				return
			}
			require.EqualValues(t, tc.expThisChainAddress, pt.HostAccAddr.String())
			require.EqualValues(t, tc.expFinalDestination, pt.Destination)
			require.EqualValues(t, tc.expPort, pt.Port)
			require.EqualValues(t, tc.expChannel, pt.Channel)
		})
	}
}

func TestGetDenomByIBCPacket(t *testing.T) {
	testCases := []struct {
		name          string
		sourcePort    string
		sourceChannel string
		destPort      string
		destChannel   string
		packetDenom   string
		expDenom      string
	}{
		{
			name:          "source token FX",
			sourcePort:    "transfer",
			sourceChannel: "channel-0",
			destPort:      "transfer",
			destChannel:   "channel-1",
			packetDenom:   "transfer/channel-0/FX",
			expDenom:      "FX",
		},
		{
			name:          "source token - eth0x61CAf09780f6F227B242EA64997a36c94a40Aa3a",
			sourcePort:    "transfer",
			sourceChannel: "channel-0",
			destPort:      "transfer",
			destChannel:   "channel-1",
			packetDenom:   "transfer/channel-0/eth0x61CAf09780f6F227B242EA64997a36c94a40Aa3a",
			expDenom:      "eth0x61CAf09780f6F227B242EA64997a36c94a40Aa3a",
		},
		{
			name:          "dest token - atom",
			sourcePort:    "transfer",
			sourceChannel: "channel-0",
			destPort:      "transfer",
			destChannel:   "channel-1",
			packetDenom:   "atom",
			expDenom:      transfertypes.ParseDenomTrace(fmt.Sprintf("%s/%s/%s", "transfer", "channel-1", "atom")).IBCDenom(),
		},
		{
			name:          "dest token - ibc denom a->b  b->c",
			sourcePort:    "transfer",
			sourceChannel: "channel-0",
			destPort:      "transfer",
			destChannel:   "channel-1",
			packetDenom:   "transfer/channel-2/atom",
			expDenom:      transfertypes.ParseDenomTrace(fmt.Sprintf("%s/%s/%s", "transfer", "channel-1", "transfer/channel-2/atom")).IBCDenom(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualValue := ibcrouter.GetDenomByIBCPacket(tc.sourcePort, tc.sourceChannel, tc.destPort, tc.destChannel, tc.packetDenom)
			require.EqualValues(t, tc.expDenom, actualValue)
		})
	}
}
