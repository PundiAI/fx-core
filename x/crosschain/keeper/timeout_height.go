package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

// CalExternalTimeoutHeight This gets the timeout height in External blocks.
func (k Keeper) CalExternalTimeoutHeight(ctx sdk.Context, getTimeoutCallback func(params types.Params) uint64) uint64 {
	currentFxHeight := ctx.BlockHeight()
	// we store the last observed Cosmos and Ethereum heights, we do not concern ourselves if these values
	// are zero because no batch can be produced if the last Ethereum block height is not first populated by a deposit event.
	heights := k.GetLastObservedBlockHeight(ctx)
	if heights.ExternalBlockHeight == 0 {
		return 0
	}
	params := k.GetParams(ctx)
	// we project how long it has been in milliseconds since the last Ethereum block height was observed
	projectedMillis := (uint64(currentFxHeight) - heights.BlockHeight) * params.AverageBlockTime
	// we convert that projection into the current Ethereum height using the average Ethereum block time in millis
	projectedCurrentEthereumHeight := (projectedMillis / params.AverageExternalBlockTime) + heights.ExternalBlockHeight
	// we convert our target time for block timeouts (lets say 12 hours) into a number of blocks to
	// place on top of our projection of the current Ethereum block height.
	timeout := getTimeoutCallback(params)
	blocksToAdd := timeout / params.AverageExternalBlockTime
	return projectedCurrentEthereumHeight + blocksToAdd
}

func GetBridgeCallTimeout(params types.Params) uint64 {
	return params.BridgeCallTimeout
}

func GetExternalBatchTimeout(params types.Params) uint64 {
	return params.ExternalBatchTimeout
}
