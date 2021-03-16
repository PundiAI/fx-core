package simulation_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/capability/keeper"
	"github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/functionx/fx-core/app"
	ibctransferkeeper "github.com/functionx/fx-core/x/ibc/applications/transfer/keeper"
	simulation2 "github.com/functionx/fx-core/x/ibc/applications/transfer/simulation"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"
)

func TestDecodeStore(t *testing.T) {
	config := app.MakeEncodingConfig()
	transferKeeper := ibctransferkeeper.NewKeeper(config.Marshaler, nil, types.Subspace{}, nil, nil, nil, nil, keeper.ScopedKeeper{})
	//app := simapp.Setup(false)
	dec := simulation2.NewDecodeStore(transferKeeper)

	trace := ibctransfertypes.DenomTrace{
		BaseDenom: "uatom",
		Path:      "transfer/channelToA",
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   ibctransfertypes.PortKey,
				Value: []byte(ibctransfertypes.PortID),
			},
			{
				Key:   ibctransfertypes.DenomTraceKey,
				Value: transferKeeper.MustMarshalDenomTrace(trace),
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
			},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"PortID", fmt.Sprintf("Port A: %s\nPort B: %s", ibctransfertypes.PortID, ibctransfertypes.PortID)},
		{"DenomTrace", fmt.Sprintf("DenomTrace A: %s\nDenomTrace B: %s", trace.IBCDenom(), trace.IBCDenom())},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
