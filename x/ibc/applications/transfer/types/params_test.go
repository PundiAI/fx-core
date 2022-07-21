package types_test

import (
	"testing"

	"github.com/functionx/fx-core/v2/x/ibc/applications/transfer/types"

	"github.com/stretchr/testify/require"
)

func TestValidateParams(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate())
	require.NoError(t, types.NewParams(true, false).Validate())
}
