package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlock update system contract
func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	if !k.HasInit(ctx) {
		return
	}
	cacheCtx, commit := ctx.CacheContext()
	if err := contracts.UpgradeSystemContract(cacheCtx, k.evmKeeper); err != nil {
		ctx.Logger().Error("begin block upgrade system contract", "error", err, "height", ctx.BlockHeight())
		//TODO need record, if failed, current height abi can not load
	} else {
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		commit()
	}
}
