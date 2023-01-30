package tests

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app"
	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ibctransfertypes "github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"
)

func TestAminoEncode(t *testing.T) {
	testcases := []struct {
		name     string
		expected string
		msg      interface{}
	}{
		{
			name:     "upgrade-SoftwareUpgradeProposal",
			expected: `{"type":"cosmos-sdk/SoftwareUpgradeProposal","value":{"description":"foo","plan":{"height":"123","info":"foo","name":"foo","time":"0001-01-01T00:00:00Z"},"title":"v2"}}`,
			msg: upgradetypes.SoftwareUpgradeProposal{
				Title:       "v2",
				Description: "foo",
				Plan: upgradetypes.Plan{
					Name:   "foo",
					Time:   time.Time{},
					Height: 123,
					Info:   "foo",
				},
			},
		},
		{
			name:     "ibc-MsgTransfer",
			expected: `{"type":"fxtransfer/MsgTransfer","value":{"fee":{"amount":"0","denom":"FX"},"receiver":"0x001","sender":"fx1001","source_channel":"channel-0","source_port":"transfer","timeout_height":{},"timeout_timestamp":"1675063442000000000","token":{"amount":"1","denom":"FX"}}}`,
			msg: ibctransfertypes.MsgTransfer{
				SourcePort:       "transfer",
				SourceChannel:    "channel-0",
				Token:            sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1)),
				Sender:           "fx1001",
				Receiver:         "0x001",
				TimeoutHeight:    clienttypes.Height{},
				TimeoutTimestamp: 1675063442000000000,
				Router:           "",
				Fee:              sdk.NewCoin(fxtypes.DefaultDenom, sdk.ZeroInt()),
				Memo:             "",
			},
		},
		{
			name:     "erc20-RegisterCoinProposal",
			expected: `{"description":"foo","metadata":{"base":"test","denom_units":[{"aliases":["ethtest"],"denom":"test"},{"denom":"TEST","exponent":18}],"description":"test","display":"test","name":"test name","symbol":"TEST"},"title":"v2"}`,
			msg: erc20types.RegisterCoinProposal{
				Title:       "v2",
				Description: "foo",
				Metadata: types.Metadata{
					Description: "test",
					DenomUnits: []*types.DenomUnit{
						{
							Denom:    "test",
							Exponent: 0,
							Aliases: []string{
								"ethtest",
							},
						},
						{
							Denom:    "TEST",
							Exponent: 18,
							Aliases:  []string{},
						},
					},
					Base:    "test",
					Display: "test",
					Name:    "test name",
					Symbol:  "TEST",
				},
			},
		},
	}

	encode := app.MakeEncodingConfig()
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			aminoJson, err := encode.Amino.MarshalJSON(testcase.msg)
			require.NoError(t, err)
			require.Equal(t, testcase.expected, string(sdk.MustSortJSON(aminoJson)))
		})
	}
}
