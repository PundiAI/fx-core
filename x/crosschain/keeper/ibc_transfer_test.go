package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCovertIbcData(t *testing.T) {

	type expect struct {
		prefix  string
		port    string
		channel string
		ok      bool
	}
	testCases := []struct {
		name    string
		ibcData string
		expect  expect
	}{
		{name: "normal ibc data",
			ibcData: "66782f7472616e736665722f6368616e6e656c2d30",
			expect: expect{
				prefix:  "fx",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "normal ibc data hex fx/transfer/channel-0 to ibcData ",
			ibcData: hex.EncodeToString([]byte("fx/transfer/channel-0")),
			expect: expect{
				prefix:  "fx",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "no prefix ibc data /transfer/channel-0",
			ibcData: hex.EncodeToString([]byte("/transfer/channel-0")),
			expect: expect{
				prefix:  "",
				port:    "transfer",
				channel: "channel-0",
				ok:      true,
			},
		},
		{name: "no prefix and no port ibc data /channel-0",
			ibcData: hex.EncodeToString([]byte("/channel-0")),
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      false,
			},
		},
		{name: "empty ibc data ''",
			ibcData: hex.EncodeToString([]byte("''")),
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      false,
			},
		},
		{name: "two slash ibc data //",
			ibcData: hex.EncodeToString([]byte("//")),
			expect: expect{
				prefix:  "",
				port:    "",
				channel: "",
				ok:      true,
			},
		},
	}

	for _, tc := range testCases {
		addressPrefix, sourcePort, sourceChannel, isOk := covertIbcData(tc.ibcData)
		require.EqualValues(t, tc.expect.ok, isOk, tc.name)
		if !tc.expect.ok {
			return
		}
		require.EqualValues(t, tc.expect.prefix, addressPrefix, tc.name)
		require.EqualValues(t, tc.expect.port, sourcePort, tc.name)
		require.EqualValues(t, tc.expect.channel, sourceChannel, tc.name)
	}
}
