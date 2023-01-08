package v2

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Channelkeeper interface {
	HasPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) bool
}

type erc20Keeper interface {
	IterateIBCTransferHash(ctx sdk.Context, cb func(port, channel string, sequence uint64) bool)
	DeleteIBCTransferHash(ctx sdk.Context, port, channel string, sequence uint64) bool
}

func PruneExpirationIBCTransferHash(ctx sdk.Context, erc20Keeper erc20Keeper, channelKeeper Channelkeeper) {
	counts := make(map[string]uint64, 10)
	erc20Keeper.IterateIBCTransferHash(ctx, func(port, channel string, sequence uint64) bool {
		found := channelKeeper.HasPacketCommitment(ctx, port, channel, sequence)
		if found {
			return false
		}
		erc20Keeper.DeleteIBCTransferHash(ctx, port, channel, sequence)
		// statistics count
		counts[fmt.Sprintf("%s/%s", port, channel)] += 1
		return false
	})
	for portChannel, count := range counts {
		ctx.Logger().Info("delete expiration ibc transfer hash", "module", "erc20",
			portChannel, strconv.FormatUint(count, 10))
	}
}
