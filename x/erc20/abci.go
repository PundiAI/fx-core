package erc20

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/x/erc20/keeper"
)

// BeginBlock update system contract
func BeginBlock(ctx sdk.Context, k keeper.Keeper) {
	if !k.HasInit(ctx) {
		return
	}
	cacheCtx, commit := ctx.CacheContext()
	if err := k.UpgradeSystemContract(cacheCtx); err != nil {
		panic(fmt.Sprintf("upgrade system contract error %v", err))
	}
	ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
	commit()
}
