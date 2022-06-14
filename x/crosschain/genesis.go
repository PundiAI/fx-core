package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/crosschain/types"
)

// InitGenesis import module genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	if state.Params != nil {
		k.SetParams(ctx, state.Params)
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
