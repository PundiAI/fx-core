package v3

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Channelkeeper todo module depend interface
type Channelkeeper interface {
	HasPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) bool
}

func PruneExpirationIBCTransferRelation(ctx sdk.Context, store sdk.KVStore, channelKeeper Channelkeeper) {
	counts := make(map[string]uint64)
	iterateIBCTransferRelationLegacy(store, func(port, channel string, sequence uint64) bool {
		found := channelKeeper.HasPacketCommitment(ctx, port, channel, sequence)
		if found {
			return false
		}
		deleteIBCTransferRelationLegacy(store, port, channel, sequence)
		// statistics count
		counts[fmt.Sprintf("%s/%s", port, channel)] += 1
		return false
	})
	for portChannel, count := range counts {
		ctx.Logger().Info("delete expiration ibc transfer hash", "module", "erc20",
			portChannel, strconv.FormatUint(count, 10))
	}
}
