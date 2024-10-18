package v8_test

import (
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
	v8 "github.com/functionx/fx-core/v8/x/erc20/migrations/v8"
)

func TestOutgoingTransferKeyToOriginTokenKey(t *testing.T) {
	moduleName := tmrand.Str(5)
	txID := uint64(tmrand.Int63n(100000))
	oldKey := append(append(v8.KeyPrefixOutgoingTransfer, []byte(moduleName)...), sdk.Uint64ToBigEndian(txID)...)
	expectKey := types.NewOriginTokenKey(moduleName, txID)
	actual := v8.OutgoingTransferKeyToOriginTokenKey(oldKey)
	require.EqualValues(t, expectKey, actual)
}
