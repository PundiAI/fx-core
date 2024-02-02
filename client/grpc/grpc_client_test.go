package grpc_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/client/grpc"
	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestClient_QueryBalances(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	client, err := grpc.DailClient("http://127.0.0.1:9090")
	assert.NoError(t, err)

	balances, err := client.WithBlockHeight(1).QueryBalances("fx1ausfqqwyqn83e8x4l46qc2ydrqn0e3wnep02fs")
	assert.NoError(t, err)
	assert.False(t, balances.IsAllPositive(), balances.String())
}

func TestClient_GetChainId(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	client, err := grpc.DailClient("http://127.0.0.1:9090")
	assert.NoError(t, err)

	chainId, err := client.GetChainId()
	assert.NoError(t, err)
	assert.Equal(t, "fxcore", chainId)

	account, err := client.QueryAccount("fx17w0adeg64ky0daxwd2ugyuneellmjgnxed28x3")
	assert.NoError(t, err)

	assert.NotNil(t, account.GetPubKey())
	pubKey, err := types.NewAnyWithValue(&ethsecp256k1.PubKey{Key: account.GetPubKey().Bytes()})
	assert.NoError(t, err)
	t.Log(pubKey)
	t.Log(account)
}
