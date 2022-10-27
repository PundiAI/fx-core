package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
)

func TestGetMetadata_Validate(t *testing.T) {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	testnetMetadata := []banktypes.Metadata{wfxMetadata, testnetPUNDIXMetadata, testnetEthUSDTMetadata, testnetUSDCMetadata, testnetDAIMetadata,
		testnetPURSEMetadata, testnetUSDJMetadata, testnetUSDFMetadata, testnetTronUSDTMetadata, testnetLINKMetadata}
	mainnetMetadata := []banktypes.Metadata{wfxMetadata, mainnetPUNDIXMetadata, mainnetPURSEMetadata}

	for _, m := range append(testnetMetadata, mainnetMetadata...) {
		err := m.Validate()
		assert.NoError(t, err)
	}
}
