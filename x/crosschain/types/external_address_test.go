package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/functionx/fx-core/v8/app"
	arbitrumtypes "github.com/functionx/fx-core/v8/x/arbitrum/types"
	avalanchetypes "github.com/functionx/fx-core/v8/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v8/x/bsc/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
	optimismtypes "github.com/functionx/fx-core/v8/x/optimism/types"
	polygontypes "github.com/functionx/fx-core/v8/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

func TestValidateModuleName(t *testing.T) {
	for _, name := range []string{
		ethtypes.ModuleName,
		bsctypes.ModuleName,
		polygontypes.ModuleName,
		trontypes.ModuleName,
		avalanchetypes.ModuleName,
		arbitrumtypes.ModuleName,
		optimismtypes.ModuleName,
	} {
		require.NoError(t, types.ValidateModuleName(name))
	}
}
