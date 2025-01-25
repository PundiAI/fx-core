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

func upgradeTestnet(ctx sdk.Context, app *keepers.AppKeepers) error {
	return updateTestnetFXBridgeToken(ctx, app.Erc20Keeper)
}

func updateTestnetFXBridgeToken(ctx sdk.Context, erc20Keeper erc20keeper.Keeper) error {
	bridgeToken, err := erc20Keeper.GetBridgeToken(ctx, ethtypes.ModuleName, fxtypes.LegacyFXDenom)
	if err != nil {
		return err
	}
	if err = erc20Keeper.BridgeToken.Remove(ctx, collections.Join(ethtypes.ModuleName, fxtypes.LegacyFXDenom)); err != nil {
		return err
	}
	if err = erc20Keeper.DenomIndex.Remove(ctx, fxtypes.LegacyFXDenom); err != nil {
		return err
	}
	bridgeToken.Denom = fxtypes.FXDenom
	if err = erc20Keeper.BridgeToken.Set(ctx, collections.Join(ethtypes.ModuleName, fxtypes.FXDenom), bridgeToken); err != nil {
		return err
	}
	bridgeDenom := erc20types.NewBridgeDenom(ethtypes.ModuleName, bridgeToken.Contract)
	return erc20Keeper.DenomIndex.Set(ctx, bridgeDenom, fxtypes.FXDenom)
}
