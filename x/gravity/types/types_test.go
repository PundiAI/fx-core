package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	_ "github.com/functionx/fx-core/app"
	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/gravity/types"
)

func TestCovertIbcPacketReceiveAddressByPrefix(t *testing.T) {
	fxtypes.ChangeNetworkForTest(fxtypes.NetworkDevnet())
	testcases := []struct {
		name    string
		height  int64
		prefix  string
		address sdk.AccAddress
		expect  string
		err     error
	}{
		{
			name:    "normal: not support 0x prefix, expect px prefix",
			height:  fxtypes.EvmV1SupportBlock() - 1,
			prefix:  "px",
			address: sdk.AccAddress("____________________"),
			expect:  "px1ta047h6lta047h6lta047h6lta047h6l0zfr4s",
			err:     nil,
		},
		{
			name:    "normal: support 0x prefix, expect px prefix",
			height:  fxtypes.EvmV1SupportBlock(),
			prefix:  "px",
			address: sdk.AccAddress("____________________"),
			expect:  "px1ta047h6lta047h6lta047h6lta047h6l0zfr4s",
			err:     nil,
		},
		{
			name:    "normal: support 0x prefix after, expect px prefix",
			height:  fxtypes.EvmV1SupportBlock() + 1,
			prefix:  "px",
			address: sdk.AccAddress("____________________"),
			expect:  "px1ta047h6lta047h6lta047h6lta047h6l0zfr4s",
			err:     nil,
		},
		{
			name:    "normal: not support 0x prefix, expect 0x prefix",
			height:  fxtypes.EvmV1SupportBlock() - 1,
			prefix:  "0x",
			address: sdk.AccAddress("____________________"),
			expect:  "0x1ta047h6lta047h6lta047h6lta047h6lvx3jz6",
			err:     nil,
		},
		{
			name:    "normal: support 0x prefix, expect eth address",
			height:  fxtypes.EvmV1SupportBlock(),
			prefix:  "0x",
			address: sdk.AccAddress("____________________"),
			expect:  "0x5f5f5f5f5f5F5f5F5F5F5F5f5F5f5f5F5F5F5F5f",
			err:     nil,
		},
		{
			name:    "normal: support 0x prefix after, expect eth address",
			height:  fxtypes.EvmV1SupportBlock() + 1,
			prefix:  "0x",
			address: sdk.AccAddress("____________________"),
			expect:  "0x5f5f5f5f5f5F5f5F5F5F5F5f5F5f5f5F5F5F5F5f",
			err:     nil,
		},
		{
			name:    "normal: support 0x prefix after, lower prefix, expect eth address",
			height:  fxtypes.EvmV1SupportBlock() + 1,
			prefix:  "0x",
			address: sdk.AccAddress("____________________"),
			expect:  "0x5f5f5f5f5f5F5f5F5F5F5F5f5F5f5f5F5F5F5F5f",
			err:     nil,
		},
		{
			name:    "normal: support 0x prefix after, upper prefix, expect eth address",
			height:  fxtypes.EvmV1SupportBlock() + 1,
			prefix:  "0X",
			address: sdk.AccAddress("____________________"),
			expect:  "0x5f5f5f5f5f5F5f5F5F5F5F5f5F5f5f5F5F5F5F5f",
			err:     nil,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			actual, err := types.CovertIbcPacketReceiveAddressByPrefix(testcase.height, testcase.prefix, testcase.address)
			if testcase.err != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, testcase.expect, actual)
		})
	}

}
