package crosschain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/precompiles/crosschain"
)

func TestBridgeCoinAmountMethod_ABI(t *testing.T) {
	bridgeCoinAmountABI := crosschain.NewBridgeCoinAmountABI()
	assert.Len(t, bridgeCoinAmountABI.Inputs, 2)
	assert.Len(t, bridgeCoinAmountABI.Outputs, 1)
}
