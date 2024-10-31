package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestBridgeCoinAmountMethod_ABI(t *testing.T) {
	bridgeCoinAmountABI := precompile.NewBridgeCoinAmountABI()
	assert.Len(t, bridgeCoinAmountABI.Inputs, 2)
	assert.Len(t, bridgeCoinAmountABI.Outputs, 1)
}
