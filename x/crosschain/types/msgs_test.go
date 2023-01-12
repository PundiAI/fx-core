package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/functionx/fx-core/v3/app"
	avalanchetypes "github.com/functionx/fx-core/v3/x/avalanche/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestValidateModuleName(t *testing.T) {
	for _, name := range []string{
		gravitytypes.ModuleName,
		ethtypes.ModuleName,
		bsctypes.ModuleName,
		polygontypes.ModuleName,
		trontypes.ModuleName,
		avalanchetypes.ModuleName,
	} {
		assert.NoError(t, types.ValidateModuleName(name))
	}
}
