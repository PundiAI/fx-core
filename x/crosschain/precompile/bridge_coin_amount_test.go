package precompile_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
)

func TestBridgeCoinAmountMethod_ABI(t *testing.T) {
	bridgeCoinAmount := precompile.NewBridgeCoinAmountMethod(nil).Method
	assert.Equal(t, 2, len(bridgeCoinAmount.Inputs))
	assert.Equal(t, 1, len(bridgeCoinAmount.Outputs))
}
