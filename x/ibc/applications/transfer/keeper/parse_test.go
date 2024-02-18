package keeper

import (
	"bytes"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

func TestParseReceiveAndAmountByPacket(t *testing.T) {
	type expect struct {
		address []byte
		amount  sdkmath.Int
		fee     sdkmath.Int
	}
	testCases := []struct {
		name      string
		packet    types.FungibleTokenPacketData
		isEvmAddr bool
		expPass   bool
		err       error
		expect    expect
	}{
		{
			"no router - expect address is receive",
			types.FungibleTokenPacketData{Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "0"},
			false, true, nil,
			expect{address: sdk.AccAddress("receive1"), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(0)},
		},
		{
			"no router - expect fee is 0, input 1",
			types.FungibleTokenPacketData{Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "0"},
			false, true, nil,
			expect{address: sdk.AccAddress("receive1"), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(0)},
		},
		{
			"no router - receive is 0x address",
			types.FungibleTokenPacketData{Receiver: common.BytesToAddress([]byte{0x1}).String(), Amount: "1", Fee: "0"},
			true, true, nil,
			expect{address: sdk.AccAddress(common.BytesToAddress([]byte{0x1}).Bytes()), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(0)},
		},
		{
			"router - expect address is sender",
			types.FungibleTokenPacketData{Sender: sdk.AccAddress("sender1").String(), Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "0", Router: "erc20"},
			false, true, nil,
			expect{address: sdk.AccAddress("sender1"), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(0)},
		},
		{
			"router - expect fee is 1, input 1",
			types.FungibleTokenPacketData{Sender: sdk.AccAddress("sender1").String(), Receiver: sdk.AccAddress("receive1").String(), Amount: "1", Fee: "1", Router: "erc20"},
			false, true, nil,
			expect{address: sdk.AccAddress("sender1"), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(1)},
		},
		{
			"router - expect address is sender, input eip address",
			types.FungibleTokenPacketData{Sender: "0x50194ffc34DB0fb3De90A4eE75dB66e868AD7820", Receiver: "0xa5d890DA1b82B69383DbB5768B42138e0Ee435c8", Amount: "1", Fee: "1", Router: "erc20"},
			false, true, nil,
			expect{address: common.HexToAddress("0x50194ffc34DB0fb3De90A4eE75dB66e868AD7820").Bytes(), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(1)},
		},
		{
			"router - expect address is sender, input eip address",
			types.FungibleTokenPacketData{Sender: "0x50194ffc34DB0fb3De90A4eE75dB66e868AD7820", Receiver: "0xa5d890DA1b82B69383DbB5768B42138e0Ee435c8", Amount: "1", Fee: "1", Router: "erc20"},
			false, true, nil,
			expect{address: common.HexToAddress("0x50194ffc34DB0fb3De90A4eE75dB66e868AD7820").Bytes(), amount: sdkmath.NewIntFromUint64(1), fee: sdkmath.NewIntFromUint64(1)},
		},
		{
			"error router - expect error, sender eip address is lowercase",
			types.FungibleTokenPacketData{Sender: "0x50194ffc34db0fb3de90a4ee75db66e868ad7820", Receiver: "0x50194ffc34DB0fb3De90A4eE75dB66e868AD7820", Amount: "1", Fee: "1", Router: "erc20"},
			false, false,
			fmt.Errorf("decoding bech32 failed: invalid character not part of charset: 98\nmismatch expected: 0x50194ffc34DB0fb3De90A4eE75dB66e868AD7820, got: 0x50194ffc34db0fb3de90a4ee75db66e868ad7820"),
			expect{address: []byte{}, amount: sdkmath.Int{}, fee: sdkmath.Int{}},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualAddress, isEvmAddr, actualAmount, actualFee, err := parseReceiveAndAmountByPacket(tc.packet)
			if tc.expPass {
				require.NoError(t, err, "valid test case %d failed: %v", i, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, err.Error(), tc.err.Error())
			}
			require.Truef(t, bytes.Equal(tc.expect.address, actualAddress.Bytes()), "expected %s, actual %s", sdk.AccAddress(tc.expect.address).String(), actualAddress.String())
			require.EqualValues(t, tc.expect.amount.String(), actualAmount.String())
			require.EqualValues(t, tc.expect.fee.String(), actualFee.String())
			require.EqualValues(t, tc.isEvmAddr, isEvmAddr)
		})
	}
}

func TestParseAmountAndFeeByPacket(t *testing.T) {
	type expect struct {
		amount sdkmath.Int
		fee    sdkmath.Int
	}
	testCases := []struct {
		name    string
		packet  types.FungibleTokenPacketData
		expPass bool
		errStr  string
		expect  expect
	}{
		{
			"pass - no router only amount ",
			types.FungibleTokenPacketData{Amount: "1"},
			true, "",
			expect{amount: sdkmath.NewInt(1), fee: sdkmath.ZeroInt()},
		},
		{
			"error - amount is empty",
			types.FungibleTokenPacketData{Amount: ""},
			false,
			"unable to parse transfer amount () into sdkmath.Int: invalid token amount",
			expect{amount: sdkmath.Int{}, fee: sdkmath.Int{}},
		},
		{
			"error - fee is empty",
			types.FungibleTokenPacketData{Amount: "1", Fee: "", Router: "aaa"},
			false,
			"fee amount is invalid:: invalid token amount",
			expect{amount: sdkmath.Int{}, fee: sdkmath.Int{}},
		},
		{
			"error - fee is negative",
			types.FungibleTokenPacketData{Amount: "1", Fee: "-1", Router: "aaa"},
			false,
			"fee amount is invalid:-1: invalid token amount",
			expect{amount: sdkmath.Int{}, fee: sdkmath.Int{}},
		},
		{
			"pass - fee is zero",
			types.FungibleTokenPacketData{Amount: "1", Fee: "0", Router: "aaa"},
			true,
			"",
			expect{amount: sdkmath.NewInt(1), fee: sdkmath.ZeroInt()},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualAmount, actualFee, err := parseAmountAndFeeByPacket(tc.packet)
			if tc.expPass {
				require.NoError(t, err, "valid test case %d failed: %v", i, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, tc.errStr, err.Error())
			}
			require.EqualValues(t, tc.expect.amount.String(), actualAmount.String())
			require.EqualValues(t, tc.expect.fee.String(), actualFee.String())
		})
	}
}

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
