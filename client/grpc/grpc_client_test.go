package grpc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v3/client/grpc"
)

func TestClient_QueryBalances(t *testing.T) {
	t.SkipNow()
	client, err := grpc.DailClient("http://127.0.0.1:9090")
	assert.NoError(t, err)

	balances, err := client.WithBlockHeight(1).QueryBalances("fx1ausfqqwyqn83e8x4l46qc2ydrqn0e3wnep02fs")
	assert.NoError(t, err)
	assert.False(t, balances.IsAllPositive(), balances.String())
}
