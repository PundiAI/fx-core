package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func TestGetMetadata_Validate(t *testing.T) {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	for _, _ = range []string{NetworkDevnet(), NetworkTestnet(), NetworkMainnet()} {
		for _, metadata := range GetMetadata() {
			err := metadata.Validate()
			assert.NoError(t, err)
		}
	}
}
