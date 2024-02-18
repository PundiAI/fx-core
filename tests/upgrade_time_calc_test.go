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

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func TestCalculateUpgradeHeight(t *testing.T) {
	helpers.SkipTest(t)
	blockInterval := 20000

	// example: UPGRADE_TIME=2023-08-10T08:00:00Z
	upgradeTime := os.Getenv("UPGRADE_TIME")
	if len(upgradeTime) == 0 {
		upgradeTime = time.Now().AddDate(0, 0, 14).Format(time.RFC3339)
	}
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

	beforeHeight := latestHeight - int64(blockInterval)
	beforeBlock, err := rpcClient.Block(ctx, &beforeHeight)
	require.NoError(t, err)
	blockTime := float64(latestTime.Unix()-beforeBlock.Block.Time.Unix()) / float64(blockInterval)
	t.Logf("Avg blcok time:%.4f", blockTime)

	blockCount := int64(float64(expectTime.Unix()-latestTime.Unix()) / blockTime)
	t.Logf("Remaining Blocks:%d", blockCount)
	t.Logf("Expected Blcok:%d", latestHeight+blockCount)
}
