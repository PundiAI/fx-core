package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, &data.Params)

	// ensure erc20 module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the erc20 module account has not been set")
	}

	for _, pair := range data.TokenPairs {
		k.AddTokenPair(ctx, pair)
	}

	if _, found := k.GetTokenPair(ctx, fxtypes.DefaultDenom); !found {
		pair, err := k.RegisterCoin(ctx, fxtypes.GetFXMetaData(fxtypes.DefaultDenom))
		if err != nil {
			panic(fmt.Sprintf("register default denom error %s", err.Error()))
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRegisterCoin,
			sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
		))
	}
}

// ExportGenesis export module status
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:     k.GetParams(ctx),
		TokenPairs: k.GetAllTokenPairs(ctx),
	}
}
