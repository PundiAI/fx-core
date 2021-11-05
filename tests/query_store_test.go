package tests

import (
	"context"
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/jsonrpc/client"

	"github.com/functionx/fx-core/x/gravity/types"
)

func TestStore(t *testing.T) {
	if !testing.Short() {
		t.SkipNow()
	}
	httpClient, err := client.DefaultHTTPClient("tcp://localhost:26657")
	require.NoError(t, err)
	cli, err := http.NewWithClient("http://localhost:26657", "/websocket", httpClient)
	require.NoError(t, err)
	abciQueryRes, err := cli.ABCIQuery(context.Background(), "store/gravity/key", types.GetDenomToERC20Key(""))
	require.NoError(t, err)

	if len(abciQueryRes.Response.Value) <= 0 {
		t.Log(abciQueryRes.Response)
		t.Fatal("not found key data")
	}
	// abciQueryRes.Response.Value
	t.Log(string(abciQueryRes.Response.Value))
}

func TestGetLastEventNonceByValidatorKey(t *testing.T) {
	bech32, err := sdk.ValAddressFromBech32("fxvaloper1nklmn350dvyykphgnr3308rn9kpnj32gaphy4f")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(types.GetLastEventNonceByValidatorKey(bech32)))

	accAddr, err := sdk.AccAddressFromBech32("fx1qllms2p25gec8fn4xvyak83g856xdltp4wc335")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(types.GetOrchestratorAddressKey(accAddr)))

	t.Log(hex.EncodeToString(types.LastObservedEventNonceKey))

	hexBytes, err := hex.DecodeString("000000000000049b")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(types.UInt64FromBytes(hexBytes))

	t.Log(hex.EncodeToString(types.GetValsetKey(255)))
}
