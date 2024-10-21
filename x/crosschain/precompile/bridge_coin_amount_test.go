package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestBridgeCoinAmountMethod_ABI(t *testing.T) {
	bridgeCoinAmount := precompile.NewBridgeCoinAmountMethod(nil).Method
	assert.Len(t, bridgeCoinAmount.Inputs, 2)
	assert.Len(t, bridgeCoinAmount.Outputs, 1)
}
