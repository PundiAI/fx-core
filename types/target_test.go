package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/types"
)

func TestParseTargetIBC(t *testing.T) {
	type expect struct {
		target  string
		prefix  string
		port    string
		channel string
		isIBC   bool
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
				isIBC:   true,
			},
		},
		{
			name:      "normal ibc data hex 0x/transfer/channel-0 to targetStr ",
			targetStr: "0x/transfer/channel-0",
			expect: expect{
				prefix:  "0x",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "normal ibc data hex upper prefix 0X/transfer/channel-0 to targetStr ",
			targetStr: "0X/transfer/channel-0",
			expect: expect{
				prefix:  "0X",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "no prefix ibc data /transfer/channel-0",
			targetStr: "/transfer/channel-0",
			expect: expect{
				prefix:  "",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "no prefix and no port ibc data /channel-0",
			targetStr: "/channel-0",
			expect: expect{
				target: "/channel-0",
				isIBC:  false,
			},
		},
		{
			name:      "empty ibc data ''",
			targetStr: "''",
			expect: expect{
				target: "''",
				isIBC:  false,
			},
		},
		{
			name:      "two slash ibc data //",
			targetStr: "//",
			expect: expect{
				isIBC: true,
			},
		},
		{
			name:      "chain prefix",
			targetStr: "chain/gravity",
			expect: expect{
				target: "gravity",
				isIBC:  false,
			},
		},
		{
			name:      "chain prefix, empty module",
			targetStr: "chain/",
			expect: expect{
				target: "",
				isIBC:  false,
			},
		},
		{
			name:      "module",
			targetStr: "gravity",
			expect: expect{
				target: "gravity",
				isIBC:  false,
			},
		},
		{
			name:      "empty",
			targetStr: "",
			expect: expect{
				target: "",
				isIBC:  false,
			},
		},
		{
			name:      "ibc prefix with channel/prefix",
			targetStr: "ibc/0/px",
			expect: expect{
				prefix:  "px",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with channel/prefix, but empty address prefix",
			targetStr: "ibc/0/",
			expect: expect{
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with channel/prefix, but empty channel sequence",
			targetStr: "ibc//px",
			expect: expect{
				prefix:  "px",
				port:    "transfer",
				channel: "channel-",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with channel/prefix, but empty channel sequence and address prefix",
			targetStr: "ibc//",
			expect: expect{
				port:    "transfer",
				channel: "channel-",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel",
			targetStr: "ibc/px/transfer/channel-0",
			expect: expect{
				prefix:  "px",
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty address prefix",
			targetStr: "ibc//transfer/channel-0",
			expect: expect{
				port:    "transfer",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port",
			targetStr: "ibc/px//channel-0",
			expect: expect{
				prefix:  "px",
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty channel",
			targetStr: "ibc/px/transfer/",
			expect: expect{
				prefix: "px",
				port:   "transfer",
				isIBC:  true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port and address prefix",
			targetStr: "ibc///channel-0",
			expect: expect{
				channel: "channel-0",
				isIBC:   true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty port and channel",
			targetStr: "ibc/px//",
			expect: expect{
				prefix: "px",
				isIBC:  true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty prefix and channel",
			targetStr: "ibc//transfer/",
			expect: expect{
				port:  "transfer",
				isIBC: true,
			},
		},
		{
			name:      "ibc prefix with prefix/port/channel, but empty all",
			targetStr: "ibc///",
			expect: expect{
				isIBC: true,
			},
		},
		{
			name:      "ibc prefix with '/'",
			targetStr: "ibc/",
			expect: expect{
				target: "ibc/",
				isIBC:  false,
			},
		},
	}

	for _, tc := range testCases {
		target := types.ParseFxTarget(tc.targetStr)
		require.EqualValues(t, tc.expect.isIBC, target.IsIBC(), tc.name)
		require.EqualValues(t, tc.expect.target, target.GetTarget(), tc.name)
		if tc.expect.isIBC {
			require.EqualValues(t, tc.expect.prefix, target.Prefix, tc.name)
			require.EqualValues(t, tc.expect.port, target.SourcePort, tc.name)
			require.EqualValues(t, tc.expect.channel, target.SourceChannel, tc.name)
		}
	}
}
