package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func TestParseTargetIBC(t *testing.T) {
	t.Parallel()
	type expect struct {
		target  string
		prefix  string
		channel string
		isIBC   bool
		pass    bool
	}
	testCases := []struct {
		name      string
		targetStr string
		expect    expect
	}{
		{
			name:      "normal ibc data hex fx/transfer/channel-0 to targetStr (legacy: no support)",
			targetStr: "fx/transfer/channel-0",
		},
		{
			name:      "normal ibc data hex 0x/transfer/channel-0 to targetStr (legacy: no support)",
			targetStr: "0x/transfer/channel-0",
		},
		{
			name:      "normal ibc data hex upper prefix 0X/transfer/channel-0 to targetStr (legacy: no support)",
			targetStr: "0X/transfer/channel-0",
		},
		{
			name:      "no prefix ibc data /transfer/channel-0",
			targetStr: "/transfer/channel-0",
		},
		{
			name:      "no prefix and no port ibc data /channel-0",
			targetStr: "/channel-0",
		},
		{
			name:      "empty ibc data ''",
			targetStr: "''",
			expect: expect{
				pass:   true,
				target: "''",
			},
		},
		{
			name:      "two slash ibc data //",
			targetStr: "//",
		},
		{
			name:      "chain prefix",
			targetStr: "chain/gravity",
			expect: expect{
				pass:   true,
				target: "eth",
			},
		},
		{
			name:      "chain prefix, empty module",
			targetStr: "chain/",
			expect: expect{
				pass:   true,
				target: "",
			},
		},
		{
			name:      "module",
			targetStr: "gravity",
			expect: expect{
				pass:   true,
				target: "eth",
			},
		},
		{
			name:      "empty",
			targetStr: "",
			expect: expect{
				pass:   true,
				target: "",
			},
		},
		{
			name:      "ibc prefix with channel/prefix",
			targetStr: "ibc/0/px",
			expect: expect{
				pass:    true,
				prefix:  "px",
				channel: "channel-0",
				isIBC:   true,
				target:  "eth", // temp
			},
		},
		{
			name:      "empty address prefix",
			targetStr: "ibc/0/",
		},
		{
			name:      "empty channel sequence",
			targetStr: "ibc//px",
		},
		{
			name:      "empty channel sequence and address prefix",
			targetStr: "ibc//",
		},
		{
			name:      "ibc prefix with prefix/port/channel",
			targetStr: "ibc/px/transfer/channel-0",
			expect: expect{
				pass:    true,
				prefix:  "px",
				channel: "channel-0",
				isIBC:   true,
				target:  "eth", // temp
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty address prefix",
			targetStr: "ibc//transfer/channel-0",
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port",
			targetStr: "ibc/px//channel-0",
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty channel",
			targetStr: "ibc/px/transfer/",
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port and address prefix",
			targetStr: "ibc///channel-0",
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port and channel",
			targetStr: "ibc/px//",
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty prefix and channel",
			targetStr: "ibc//transfer/",
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty all",
			targetStr: "ibc///",
		},
		{
			name:      "ibc prefix with '/'",
			targetStr: "ibc/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			target, err := types.ParseFxTarget(tc.targetStr)
			if tc.expect.pass {
				require.EqualValues(t, tc.expect.isIBC, target.IsIBC())
				require.EqualValues(t, tc.expect.prefix, target.Bech32Prefix)
				require.EqualValues(t, tc.expect.channel, target.IBCChannel)
				require.EqualValues(t, tc.expect.target, target.GetModuleName())
			} else {
				require.Error(t, err)
			}
		})
	}
}
