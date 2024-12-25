package keeper

import (
	"bytes"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func TestParsePacketAddress(t *testing.T) {
	testCases := []struct {
		name      string
		address   string
		expPass   bool
		isEvmAddr bool
		err       error
		expect    []byte
	}{
		{"normal fx address", sdk.AccAddress("abc").String(), true, false, nil, sdk.AccAddress("abc")},
		{"normal eip address", "0x2652554541Eff910C154fB643d2b167D743434EA", true, true, nil, common.HexToAddress("0x2652554541Eff910C154fB643d2b167D743434EA").Bytes()},

		{"err bech32 address - kc74", "fx1yef9232palu3ps25ldjr62ck046rgd8292kc74", false, false, fmt.Errorf("decoding bech32 failed: invalid checksum (expected 92kc73 got 92kc74)\nwrong length"), []byte{}},
		{"err lowercase eip address", "0x2652554541eff910c154fb643d2b167d743434ea", false, false, fmt.Errorf("decoding bech32 failed: invalid checksum (expected j389ls got 3434ea)\nmismatch expected: 0x2652554541Eff910C154fB643d2b167D743434EA, got: 0x2652554541eff910c154fb643d2b167d743434ea"), []byte{}},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualAddress, isEvmAddr, err := fxtypes.ParseAddress(tc.address)
			if tc.expPass {
				require.EqualValues(t, tc.isEvmAddr, isEvmAddr)
				require.NoError(t, err, "valid test case %d failed: %v", i, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, err.Error(), tc.err.Error())
			}
			require.Truef(t, bytes.Equal(tc.expect, actualAddress.Bytes()), "expected %s, actual %s", sdk.AccAddress(tc.expect).String(), actualAddress.String())
		})
	}
}
