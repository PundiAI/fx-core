package parser_test

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v3/x/ibc/ibcrouter/parser"
)

func TestParseReceiverDataTransfer(t *testing.T) {
	data := "cosmos1vzxkv3lxccnttr9rs0002s93sgw72h7ghukuhs|transfer/channel-0:cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k"
	pt, err := parser.ParseReceiverData(data)

	require.NoError(t, err)
	require.True(t, pt.ShouldForward)
	require.Equal(t, pt.HostAccAddr.String(), "cosmos1vzxkv3lxccnttr9rs0002s93sgw72h7ghukuhs")
	require.Equal(t, pt.Destination, "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k")
	require.Equal(t, pt.Port, "transfer")
	require.Equal(t, pt.Channel, "channel-0")
}

func TestParseReceiverDataNoTransfer(t *testing.T) {
	data := "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k"
	pt, err := parser.ParseReceiverData(data)

	require.NoError(t, err)
	require.False(t, pt.ShouldForward)
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
