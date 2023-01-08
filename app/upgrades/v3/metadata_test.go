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

	for _, metadata := range append(getMetadata(fxtypes.MainnetChainId), getMetadata(fxtypes.TestnetChainId)...) {
		assert.NoError(t, metadata.Validate())
		assert.NoError(t, fxtypes.ValidateMetadata(metadata))
	}
}
