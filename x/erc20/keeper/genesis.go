package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

// InitGenesis import module genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	// ensure erc20 module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		return sdkerrors.ErrNotFound.Wrapf("module account %s", types.ModuleName)
	}

	metadata := fxtypes.NewDefaultMetadata()
	_, err := k.RegisterNativeCoin(ctx, metadata.Name, metadata.Symbol, fxtypes.DenomUnit)
	if err != nil {
		return sdkerrors.ErrLogic.Wrapf("failed to register native coin: %s", err.Error())
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
