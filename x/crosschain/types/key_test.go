package types_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func TestGetBatchConfirmKey(t *testing.T) {
	tokenContract := helpers.GenerateAddress().Hex()
	batchNonce := rand.Uint64()
	keys1 := append(types.BatchConfirmKey, append([]byte(tokenContract), sdk.Uint64ToBigEndian(batchNonce)...)...)
	keys2 := types.GetBatchConfirmKey(tokenContract, batchNonce, []byte{})
	assert.Equal(t, keys1, keys2)
}
