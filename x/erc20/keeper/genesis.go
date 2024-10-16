package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	// ensure erc20 module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		return fmt.Errorf("the erc20 module account has not been set")
	}

	_, err := k.RegisterNativeCoin(ctx, fxtypes.DefaultDenom, fxtypes.DefaultDenom, fxtypes.DenomUnit)
	if err != nil {
		return fmt.Errorf("register default denom error %s", err.Error())
	}
	return nil
}

// ExportGenesis export module status
func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &types.GenesisState{
		Params: params,
	}, nil
}
