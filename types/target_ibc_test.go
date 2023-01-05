package types

import (
	"testing"

	"github.com/stretchr/testify/require"
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
		targetIBC string
		expect    expect
	}{
		{name: "normal ibc data hex fx/transfer/channel-0 to targetIBC ",
			targetIBC: "fx/transfer/channel-0",
			expect: expect{
				prefix:  "fx",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "normal ibc data hex 0x/transfer/channel-0 to targetIBC ",
			targetIBC: "0x/transfer/channel-0",
			expect: expect{
				prefix:  "0x",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "normal ibc data hex upper prefix 0X/transfer/channel-0 to targetIBC ",
			targetIBC: "0X/transfer/channel-0",
			expect: expect{
				prefix:  "0X",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "no prefix ibc data /transfer/channel-0",
			targetIBC: "/transfer/channel-0",
			expect: expect{
				prefix:  "",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "no prefix and no port ibc data /channel-0",
			targetIBC: "/channel-0",
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      false,
			},
		},
		{name: "empty ibc data ''",
			targetIBC: "''",
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      false,
			},
		},
		{name: "two slash ibc data //",
			targetIBC: "//",
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      true,
			},
		},
	}

	for _, tc := range testCases {
		targetIBC, isOk := ParseTargetIBC(tc.targetIBC)
		require.EqualValues(t, tc.expect.ok, isOk, tc.name)
		if !tc.expect.ok {
			return
		}
		require.EqualValues(t, tc.expect.prefix, targetIBC.Prefix, tc.name)
		require.EqualValues(t, tc.expect.port, targetIBC.SourcePort, tc.name)
		require.EqualValues(t, tc.expect.channel, targetIBC.SourceChannel, tc.name)
	}
}
