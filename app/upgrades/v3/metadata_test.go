package v3

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

func TestGetMetadata_Validate(t *testing.T) {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	for _, m := range append(GetMetadata(fxtypes.MainnetChainId), GetMetadata(fxtypes.TestnetChainId)...) {
		err := m.Validate()
		assert.NoError(t, err)
	}
}
