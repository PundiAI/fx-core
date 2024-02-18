package jsonrpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/v7/client/jsonrpc"
	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestNewWsClient(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	client, err := jsonrpc.NewWsClient("ws://localhost:26657/websocket", context.Background())
	assert.NoError(t, err)
	client.Logger = log.NewTMLogger(os.Stdout)
	responses := make(chan jsonrpc.RPCResponse, 1024)
	id, err := client.Subscribe("tm.event='NewBlockHeader'", responses)
	assert.NoError(t, err)
	defer client.Unsubscribe(id)
	for resp := range responses {
		t.Log(resp.Error, resp.JSONRPC, resp.ID, string(resp.Result))
	}
}

func TestQueryAccount(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	rpc := jsonrpc.NewNodeRPC(jsonrpc.NewClient("http://localhost:26657"))
	account, err := rpc.QueryAccount("fx17w0adeg64ky0daxwd2ugyuneellmjgnxed28x3")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account.GetSequence())
}
