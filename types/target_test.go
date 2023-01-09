package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/types"
)

func TestParseTargetIBC(t *testing.T) {
	type expect struct {
		prefix  string
		port    string
		channel string
		ok      bool
	}
	testCases := []struct {
		name      string
		targetStr string
		expect    expect
	}{
		{
			name:      "normal ibc data hex fx/transfer/channel-0 to targetStr ",
			targetStr: "fx/transfer/channel-0",
			expect: expect{
				prefix:  "fx",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{
			name:      "normal ibc data hex 0x/transfer/channel-0 to targetStr ",
			targetStr: "0x/transfer/channel-0",
			expect: expect{
				prefix:  "0x",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{
			name:      "normal ibc data hex upper prefix 0X/transfer/channel-0 to targetStr ",
			targetStr: "0X/transfer/channel-0",
			expect: expect{
				prefix:  "0X",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{
			name:      "no prefix ibc data /transfer/channel-0",
			targetStr: "/transfer/channel-0",
			expect: expect{
				prefix:  "",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{
			name:      "no prefix and no port ibc data /channel-0",
			targetStr: "/channel-0",
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      false,
			},
		},
		{
			name:      "empty ibc data ''",
			targetStr: "''",
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      false,
			},
		},
		{
			name:      "two slash ibc data //",
			targetStr: "//",
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      true,
			},
		},
	}

	for _, tc := range testCases {
		target := types.ParseFxTarget(tc.targetStr)
		require.EqualValues(t, tc.expect.ok, target.IsIBC(), tc.name)
		if !tc.expect.ok {
			return
		}
		require.EqualValues(t, tc.expect.prefix, target.Prefix, tc.name)
		require.EqualValues(t, tc.expect.port, target.SourcePort, tc.name)
		require.EqualValues(t, tc.expect.channel, target.SourceChannel, tc.name)
	}
}
