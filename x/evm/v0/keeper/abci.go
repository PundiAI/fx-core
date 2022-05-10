package keeper

import (
	fxtypes "github.com/functionx/fx-core/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// BeginBlock sets the sdk Context and EIP155 chain id to the Keeper.
func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	if ctx.BlockHeight() < fxtypes.EvmV0SupportBlock() ||
		ctx.BlockHeight() >= fxtypes.EvmV1SupportBlock() {
		return
	}
	k.WithContext(ctx)
	k.WithChainID(ctx)
}

// EndBlock also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func (k *Keeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) {
	if ctx.BlockHeight() < fxtypes.EvmV0SupportBlock() ||
		ctx.BlockHeight() >= fxtypes.EvmV1SupportBlock() {
		return
	}
	if !k.HasInit(ctx) {
		return
	}
	// Gas costs are handled within msg handler so costs should be ignored
	infCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	k.WithContext(infCtx)

	bloom := ethtypes.BytesToBloom(k.GetBlockBloomTransient().Bytes())
	k.EmitBlockBloomEvent(infCtx, bloom)

	k.WithContext(ctx)
}
