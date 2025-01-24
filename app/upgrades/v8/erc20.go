package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
)

func migrateErc20FXToPundiAI(ctx sdk.Context, keeper erc20keeper.Keeper) error {
	erc20Token, err := keeper.GetERC20Token(ctx, fxtypes.LegacyFXDenom)
	if err != nil {
		return err
	}
	erc20Token.Denom = fxtypes.DefaultDenom
	if err = keeper.ERC20Token.Set(ctx, erc20Token.Denom, erc20Token); err != nil {
		return err
	}
	return keeper.ERC20Token.Remove(ctx, fxtypes.LegacyFXDenom)
}
