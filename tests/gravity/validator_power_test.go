package test_gravity

import (
	"context"
	"testing"

	"github.com/functionx/fx-core/x/gravity/types"
)

func TestQueryValidatorPowerChanger(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	grpcClient, err := grpcNewClient("localhost:9090")
	if err != nil {
		t.Fatal(err)
	}
	gravityQueryClient := types.NewQueryClient(grpcClient)

	currentValsetResp, err := gravityQueryClient.CurrentValset(context.Background(), &types.QueryCurrentValsetRequest{})
	if err != nil {
		t.Fatal(err)
	}
	currentValset := currentValsetResp.GetValset()
	t.Logf("CurrentValset nonce:[%v], height:[%v], valSetSize:[%v]", currentValset.Nonce, currentValset.Height, len(currentValset.Members))
	for _, member := range currentValset.Members {
		t.Logf("ethAddress:[%v], power:[%v]", member.EthAddress, member.Power)
	}
	lastValsetRequestsResp, err := gravityQueryClient.LastValsetRequests(context.Background(), &types.QueryLastValsetRequestsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	latestValset := lastValsetRequestsResp.GetValsets()
	t.Logf("\n\n\nLatestValset nonce:[%v], height:[%v], valSetSize:[%v]", currentValset.Nonce, currentValset.Height, len(currentValset.Members))
	for index, valset := range latestValset {
		t.Logf("valset index:[%v], nonce:[%v], height:[%v], valSetSize:[%v]", index, valset.Nonce, valset.Height, len(valset.Members))
		for _, member := range valset.Members {
			t.Logf("ethAddress:[%v], power:[%v]", member.EthAddress, member.Power)
		}
	}
	powerDiff := types.BridgeValidators(currentValset.Members).PowerDiff(latestValset[0].Members)
	t.Logf("powerDiff:[%.8f]", powerDiff)
}
