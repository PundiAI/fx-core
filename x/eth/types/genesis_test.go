package types

import (
	"reflect"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want *crosschaintypes.GenesisState
	}{
		{
			name: "eth default genesis",
			want: &crosschaintypes.GenesisState{
				Params: crosschaintypes.Params{
					GravityId:                         "fx-bridge-eth",
					AverageBlockTime:                  7_000,
					AverageExternalBlockTime:          15_000,
					ExternalBatchTimeout:              12 * 3600 * 1000,
					SignedWindow:                      30_000,
					SlashFraction:                     sdk.NewDec(8).Quo(sdk.NewDec(10)),
					OracleSetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
					IbcTransferTimeoutHeight:          20_000,
					DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18)),
					DelegateMultiple:                  10,
					BridgeCallTimeout:                 crosschaintypes.DefaultBridgeCallTimeout,
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
