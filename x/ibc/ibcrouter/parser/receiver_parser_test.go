package parser_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/x/ibc/ibcrouter/parser"
)

func TestParseReceiverDataTransfer(t *testing.T) {
	ethAddr := common.BytesToAddress(rand.Bytes(20)).Hex()
	cosmosAddr := sdk.AccAddress(rand.Bytes(20)).String()
	testCases := []struct {
		name       string
		data       string
		expectPass bool
		expectData *parser.ParsedReceiver
	}{
		{
			name:       "pass - normal transfer",
			data:       cosmosAddr + "|transfer/channel-0:" + cosmosAddr,
			expectPass: true,
			expectData: &parser.ParsedReceiver{
				ShouldForward: true,
				HostAccAddr:   sdk.MustAccAddressFromBech32(cosmosAddr),
				Destination:   cosmosAddr,
				Port:          "transfer",
				Channel:       "channel-0",
			},
		},
		{
			name:       "pass - normal transfer - receive is 0x address",
			data:       cosmosAddr + "|transfer/channel-0:" + ethAddr,
			expectPass: true,
			expectData: &parser.ParsedReceiver{
				ShouldForward: true,
				HostAccAddr:   sdk.MustAccAddressFromBech32(cosmosAddr),
				Destination:   ethAddr,
				Port:          "transfer",
				Channel:       "channel-0",
			},
		},
		{
			name:       "pass - no forward",
			data:       cosmosAddr,
			expectPass: true,
			expectData: &parser.ParsedReceiver{
				ShouldForward: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			pt, err := parser.ParseReceiverData(testCase.data)

			if !testCase.expectPass {
				require.Error(t, err)
				return

			}
			require.NoError(t, err)
			if !testCase.expectData.ShouldForward {
				require.False(t, pt.ShouldForward)
				return
			}
			checkParsedReceiverData(t, pt, testCase.expectData)
		})
	}
}

func checkParsedReceiverData(t *testing.T, expect *parser.ParsedReceiver, actual *parser.ParsedReceiver) {
	require.Equal(t, expect.ShouldForward, actual.ShouldForward)
	require.Equal(t, expect.HostAccAddr, actual.HostAccAddr)
	require.Equal(t, expect.Destination, actual.Destination)
	require.Equal(t, expect.Port, actual.Port)
	require.Equal(t, expect.Channel, actual.Channel)
}

func TestParseReceiverDataErrors(t *testing.T) {
	testCases := []struct {
		name          string
		data          string
		errStartsWith string
	}{
		{
			"unparsable transfer field",
			"",
			"unparsable receiver",
		},
		{
			"unparsable transfer field",
			"abc:def:",
			"unparsable receiver",
		},
		{
			"missing pipe",
			"transfer/channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
			"formatting incorrect",
		},
		{
			"invalid this chain address",
			"somm16plylpsgxechajltx9yeseqexzdzut9g8vla4k|transfer/channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
			"decoding bech32 failed",
		},
		{
			"invalid this chain address",
			common.BytesToAddress(rand.Bytes(20)).Hex() + "|transfer/channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
			"decoding bech32 failed",
		},
		{
			"invalid this chain address - ethereum address",
			fmt.Sprintf("%s|transfer/channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k", common.BytesToAddress(rand.Bytes(20))),
			"decoding bech32 failed",
		},
		{
			"missing slash",
			"cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k|transfer\\channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
			"formatting incorrect",
		},
		{
			"missing slash",
			"cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k|transfer\\channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k",
			"formatting incorrect",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.ParseReceiverData(tc.data)
			require.Error(t, err)
			require.Equal(t, err.Error()[:len(tc.errStartsWith)], tc.errStartsWith)
		})
	}
}
