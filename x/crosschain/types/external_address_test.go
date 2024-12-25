package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/pundiai/fx-core/v8/app"
	arbitrumtypes "github.com/pundiai/fx-core/v8/x/arbitrum/types"
	avalanchetypes "github.com/pundiai/fx-core/v8/x/avalanche/types"
	bsctypes "github.com/pundiai/fx-core/v8/x/bsc/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	optimismtypes "github.com/pundiai/fx-core/v8/x/optimism/types"
	polygontypes "github.com/pundiai/fx-core/v8/x/polygon/types"
	trontypes "github.com/pundiai/fx-core/v8/x/tron/types"
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
