package types_test

import (
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	autytypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
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

func TestBridgeCallSender(t *testing.T) {
	assert.Equal(t, "0x8D5C3128408b212F7F0Dc206a981fC16c079DE19", common.BytesToAddress(autytypes.NewModuleAddress(types.BridgeCallSender)).String())
}
