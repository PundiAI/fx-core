package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	clienthttp "github.com/tendermint/tendermint/rpc/client/http"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

func TestHeight(t *testing.T) {
	t.SkipNow()
	upgradeTime := "2023-05-30T08:00:00Z"
	expectTime, err := time.Parse(time.RFC3339, upgradeTime)
	require.NoError(t, err)
	t.Logf("Expected upgrade time:%s", expectTime)

	jsonRpcUrl := os.Getenv("JSON_RPC_URL")
	if len(jsonRpcUrl) == 0 {
		jsonRpcUrl = "https://testnet-fx-json.functionx.io:26657"
	}
	httpClient, err := jsonrpcclient.DefaultHTTPClient(jsonRpcUrl)
	require.NoError(t, err)
	httpClient.Transport = http.DefaultTransport
	rpcClient, err := clienthttp.NewWithClient(jsonRpcUrl, fmt.Sprintf("%s/websocket", jsonRpcUrl), httpClient)
	require.NoError(t, err)

	ctx := context.Background()
	status, err := rpcClient.Status(ctx)
	require.NoError(t, err)

	latestHeight := status.SyncInfo.LatestBlockHeight
	latestTime := status.SyncInfo.LatestBlockTime

	require.Truef(t, expectTime.After(latestTime), "The upgrade time has expired\nExpect:%s\nCurrent:%s", expectTime, latestTime)

	beforeHeight := latestHeight - 20000
	beforeBlock, err := rpcClient.Block(ctx, &beforeHeight)
	require.NoError(t, err)
	blockInterval := float64(latestTime.Unix()-beforeBlock.Block.Time.Unix()) / float64(20000)
	t.Logf("Avg blcok time:%.4f", blockInterval)

	blockCount := int64(float64(expectTime.Unix()-latestTime.Unix()) / blockInterval)
	t.Logf("Remaining Blocks:%d", blockCount)
	t.Logf("Expected Blcok:%d", latestHeight+blockCount)
}
