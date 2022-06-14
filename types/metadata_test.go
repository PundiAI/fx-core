package types

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func TestGetMetadata_Validate(t *testing.T) {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	devnetMetadata := []banktypes.Metadata{wfxMetadata, devnetPUNDIXMetadata, devnetPURSEMetadata, devnetUSDTMetadata}
	testnetMetadata := []banktypes.Metadata{wfxMetadata, testnetPUNDIXMetadata, testnetEthUSDTMetadata, testnetUSDCMetadata, testnetDAIMetadata,
		testnetPURSEMetadata, testnetUSDJMetadata, testnetUSDFMetadata, testnetTronUSDTMetadata, testnetLINKMetadata}
	mainnetMetadata := []banktypes.Metadata{wfxMetadata, mainnetPUNDIXMetadata, mainnetPURSEMetadata, mainnetTronUSDTMetadata, mainnetPolygonUSDTMetadata}

	for _, m := range append(append(devnetMetadata, testnetMetadata...), mainnetMetadata...) {
		err := m.Validate()
		assert.NoError(t, err)
	}
}
