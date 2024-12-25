package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/gov/types"
)

func TestDefaultInitGenesisCustomParams(t *testing.T) {
	defaultInitGenesisCustomParams := types.DefaultInitGenesisCustomParams()
	require.Len(t, defaultInitGenesisCustomParams, 6)

	helpers.AssertJsonFile(t, "./init_genesis_custom_params.json", defaultInitGenesisCustomParams)
}
