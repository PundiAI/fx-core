package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func TestGetBatchConfirmKey(t *testing.T) {
	tokenContract := helpers.GenerateAddress().Hex()
	batchNonce := tmrand.Uint64()
	keys1 := append(types.BatchConfirmKey, append([]byte(tokenContract), sdk.Uint64ToBigEndian(batchNonce)...)...)
	keys2 := types.GetBatchConfirmKey(tokenContract, batchNonce, []byte{})
	assert.Equal(t, keys1, keys2)
}
