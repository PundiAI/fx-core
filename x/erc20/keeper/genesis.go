package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(ctx sdk.Context, accountKeeper authkeeper.AccountKeeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)

	// ensure erc20 module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the erc20 module account has not been set")
	}

	for _, pair := range data.TokenPairs {
		k.AddTokenPair(ctx, pair)
	}
}

// ExportGenesis export module status
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:     k.GetParams(ctx),
		TokenPairs: k.GetAllTokenPairs(ctx),
	}
}
