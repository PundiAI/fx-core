package transfer_test

import (
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"

	_ "github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

func TestUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name   string
		data   []byte
		pass   bool
		expErr error
		exp    types.FungibleTokenPacketData
	}{
		{
			name:   "fx transfer packet - no router",
			data:   types.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "", "0").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "",
				Fee:      "",
			},
		},
		{
			name:   "fx transfer packet - router with 0fee",
			data:   types.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "router", "0").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "router",
				Fee:      "0",
			},
		},
		{
			name:   "fx transfer packet - router with empty fee",
			data:   types.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "router", "").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "router",
				Fee:      "",
			},
		},
		{
			name:   "ibc transfer packet",
			data:   transfertypes.NewFungibleTokenPacketData("FX", "100", "Add1", "Add2", "").GetBytes(),
			pass:   true,
			expErr: nil,
			exp: types.FungibleTokenPacketData{
				Denom:    "FX",
				Amount:   "100",
				Sender:   "Add1",
				Receiver: "Add2",
				Router:   "",
				Fee:      "",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var packet types.FungibleTokenPacketData
			err := types.ModuleCdc.UnmarshalJSON(testCase.data, &packet)
			if testCase.pass {
				require.NoError(t, err)
				require.EqualValues(t, testCase.exp, packet)
				require.EqualValues(t, testCase.exp.GetDenom(), packet.GetDenom())
				require.EqualValues(t, testCase.exp.GetAmount(), packet.GetAmount())
				require.EqualValues(t, testCase.exp.GetSender(), packet.GetSender())
				require.EqualValues(t, testCase.exp.GetReceiver(), packet.GetReceiver())
				require.EqualValues(t, testCase.exp.GetRouter(), packet.GetRouter())
				require.EqualValues(t, testCase.exp.GetFee(), packet.GetFee())
			} else {
				require.ErrorIs(t, err, testCase.expErr)
			}
		})
	}
}
