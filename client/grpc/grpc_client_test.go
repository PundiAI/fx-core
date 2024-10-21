package grpc_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/client/grpc"
	"github.com/functionx/fx-core/v8/testutil/helpers"
)

func TestClient_QueryBalances(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	client, err := grpc.DailClient("http://127.0.0.1:9090")
	require.NoError(t, err)

	balances, err := client.WithBlockHeight(1).QueryBalances("fx1ausfqqwyqn83e8x4l46qc2ydrqn0e3wnep02fs")
	require.NoError(t, err)
	assert.False(t, balances.IsAllPositive(), balances.String())
}

func TestClient_GetChainId(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	client, err := grpc.DailClient("http://127.0.0.1:9090")
	require.NoError(t, err)

	chainId, err := client.GetChainId()
	require.NoError(t, err)
	assert.Equal(t, "fxcore", chainId)

	account, err := client.QueryAccount("fx17w0adeg64ky0daxwd2ugyuneellmjgnxed28x3")
	require.NoError(t, err)

	assert.NotNil(t, account.GetPubKey())
	_, err = types.NewAnyWithValue(&ethsecp256k1.PubKey{Key: account.GetPubKey().Bytes()})
	require.NoError(t, err)
}
