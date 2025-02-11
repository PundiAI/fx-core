package v8

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pundiai/fx-core/v8/app/keepers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20keeper "github.com/pundiai/fx-core/v8/x/erc20/keeper"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func updateMainnetPundiAI(ctx sdk.Context, app *keepers.AppKeepers) error {
	if err := migrateErc20FXToPundiAI(ctx, app.Erc20Keeper); err != nil {
		return err
	}
	if err := updateFXBridgeDenom(ctx, app.Erc20Keeper); err != nil {
		return err
	}
	return addMainnetPundiAIBridgeToken(ctx, app.Erc20Keeper)
}

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

func updateFXBridgeDenom(ctx sdk.Context, keeper erc20keeper.Keeper) error {
	pundiaiERC20Token, err := keeper.GetERC20Token(ctx, fxtypes.DefaultDenom)
	if err != nil {
		return err
	}
	if err = keeper.DenomIndex.Set(ctx, pundiaiERC20Token.Erc20Address, fxtypes.DefaultDenom); err != nil {
		return err
	}
	bridgeToken, err := keeper.GetBridgeToken(ctx, ethtypes.ModuleName, fxtypes.LegacyFXDenom)
	if err != nil {
		return err
	}
	if err = keeper.BridgeToken.Remove(ctx, collections.Join(ethtypes.ModuleName, fxtypes.LegacyFXDenom)); err != nil {
		return err
	}
	if err = keeper.DenomIndex.Remove(ctx, fxtypes.LegacyFXDenom); err != nil {
		return err
	}
	bridgeToken.Denom = fxtypes.FXDenom
	if err = keeper.BridgeToken.Set(ctx, collections.Join(ethtypes.ModuleName, fxtypes.FXDenom), bridgeToken); err != nil {
		return err
	}
	bridgeDenom := erc20types.NewBridgeDenom(ethtypes.ModuleName, bridgeToken.Contract)
	return keeper.DenomIndex.Set(ctx, bridgeDenom, fxtypes.FXDenom)
}

func addMainnetPundiAIBridgeToken(ctx sdk.Context, keeper erc20keeper.Keeper) error {
	pundiaiToken := GetMainnetBridgeToken(ctx)
	return keeper.AddBridgeToken(ctx, fxtypes.DefaultDenom, ethtypes.ModuleName, pundiaiToken.String(), false)
}
