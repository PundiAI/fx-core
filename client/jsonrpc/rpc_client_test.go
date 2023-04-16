package jsonrpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/functionx/fx-core/v4/client/jsonrpc"
)

func TestNewWsClient(t *testing.T) {
	t.SkipNow()
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
