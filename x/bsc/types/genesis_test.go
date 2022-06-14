package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	fxtypes "github.com/functionx/fx-core/types"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"reflect"
	"testing"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want *crosschaintypes.GenesisState
	}{
		{
			name: "bsc default genesis",
			want: &crosschaintypes.GenesisState{
				Params: &crosschaintypes.Params{
					GravityId:                         "fx-bridge-bsc",
					AverageBlockTime:                  5000,
					ExternalBatchTimeout:              24 * 3600 * 1e3,
					AverageExternalBlockTime:          5000,
					SignedWindow:                      20000,
					SlashFraction:                     sdk.NewDec(1).Quo(sdk.NewDec(1000)),
					OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					IbcTransferTimeoutHeight:          20000,
					DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10000).MulRaw(1e18)),
					DelegateMultiple:                  10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultGenesisState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}
