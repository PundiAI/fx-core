package precompile_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/crosschain/precompile"
	ethtypes "github.com/functionx/fx-core/v8/x/eth/types"
)

func TestExecuteClaimMethod_ABI(t *testing.T) {
	executeClaim := precompile.NewExecuteClaimMethod(nil)

	methodStr := `function executeClaim(string _chain, uint256 _eventNonce) returns(bool _result)`
	assert.Equal(t, methodStr, executeClaim.Method.String())

	eventStr := `event ExecuteClaimEvent(address indexed _sender, uint256 _eventNonce, string _chain)`
	assert.Equal(t, eventStr, executeClaim.Event.String())
}

func TestExecuteClaimMethod_PackInput(t *testing.T) {
	executeClaim := precompile.NewExecuteClaimMethod(nil)
	input, err := executeClaim.PackInput(contract.ExecuteClaimArgs{
		Chain:      ethtypes.ModuleName,
		EventNonce: big.NewInt(1),
	})
	require.NoError(t, err)
	expected := "4ac3bdc30000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000036574680000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expected, hex.EncodeToString(input))
}
