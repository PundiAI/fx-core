package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/gov/types"
)

func TestDefaultInitGenesisCustomParams(t *testing.T) {
	defaultInitGenesisCustomParams := types.DefaultInitGenesisCustomParams()
	require.EqualValues(t, 6, len(defaultInitGenesisCustomParams))

	helpers.AssertJsonFile(t, "./init_genesis_custom_params.json", defaultInitGenesisCustomParams)
}
