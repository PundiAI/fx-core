package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGetBatchConfirmKey(t *testing.T) {
	tokenContract := helpers.GenerateAddress().Hex()
	batchNonce := rand.Uint64()
	keys1 := append(types.BatchConfirmKey, append([]byte(tokenContract), sdk.Uint64ToBigEndian(batchNonce)...)...)
	keys2 := types.GetBatchConfirmKey(tokenContract, batchNonce, []byte{})
	assert.Equal(t, keys1, keys2)
}
