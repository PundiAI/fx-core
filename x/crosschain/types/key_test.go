package types_test

import (
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func TestGetBatchConfirmKey(t *testing.T) {
	tokenContract := helpers.GenHexAddress().Hex()
	batchNonce := tmrand.Uint64()
	keys1 := append(types.BatchConfirmKey, append([]byte(tokenContract), sdk.Uint64ToBigEndian(batchNonce)...)...)
	keys2 := types.GetBatchConfirmKey(tokenContract, batchNonce, []byte{})
	assert.Equal(t, keys1, keys2)
}
