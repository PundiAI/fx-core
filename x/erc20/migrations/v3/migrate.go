package v3

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type erc20Keeper interface {
	SetIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64)
}

type migrateRelation struct {
	channel  string
	sequence uint64
}

func MigrateIBCTransferRelation(ctx sdk.Context, store sdk.KVStore, erc20Keeper erc20Keeper) {
	counts := make(map[string]uint64)
	migrates := make([]migrateRelation, 0)

	iterateIBCTransferRelationLegacy(store, func(port, channel string, sequence uint64) bool {
		migrates = append(migrates, migrateRelation{channel: channel, sequence: sequence})
		deleteIBCTransferRelationLegacy(store, port, channel, sequence)
		// statistics count
		counts[fmt.Sprintf("%s/%s", port, channel)] += 1
		return false
	})

	for _, mr := range migrates {
		erc20Keeper.SetIBCTransferRelation(ctx, mr.channel, mr.sequence)
	}

	for portChannel, count := range counts {
		ctx.Logger().Info("migrate legacy ibc transfer relation", "module", "erc20",
			portChannel, strconv.FormatUint(count, 10))
	}
}
