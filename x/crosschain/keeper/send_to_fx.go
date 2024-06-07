package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) SendToFxExecuted(ctx sdk.Context, claim *types.MsgSendToFxClaim) error {
	bridgeToken := k.GetBridgeTokenDenom(ctx, claim.TokenContract)
	if bridgeToken == nil {
		return errorsmod.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	coin := sdk.NewCoin(bridgeToken.Denom, claim.Amount)
	receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
	if err != nil {
		return errorsmod.Wrap(types.ErrInvalid, "receiver address")
	}
	isOriginOrConverted := k.erc20Keeper.IsOriginOrConvertedDenom(ctx, bridgeToken.Denom)
	if !isOriginOrConverted {
		// If it is not fxcore originated, mint the coins (aka vouchers)
		if err = k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
			return errorsmod.Wrapf(err, "mint vouchers coins")
		}
	}
	if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, sdk.NewCoins(coin)); err != nil {
		return errorsmod.Wrap(err, "transfer vouchers")
	}

	// convert to base denom
	cacheCtx, commit := ctx.CacheContext()
	targetCoin, err := k.erc20Keeper.ConvertDenomToTarget(cacheCtx, receiveAddr, coin, fxtypes.ParseFxTarget(fxtypes.ERC20Target))
	if err != nil {
		k.Logger(ctx).Info("failed to convert base denom", "error", err)
		return nil
	}
	commit()

	// relay transfer
	if err = k.RelayTransferHandler(ctx, claim.EventNonce, claim.TargetIbc, receiveAddr, targetCoin); err != nil {
		k.Logger(ctx).Info("failed to relay transfer", "error", err)
		return nil
	}

	k.HandlePendingOutgoingTx(ctx, receiveAddr, claim.GetEventNonce(), bridgeToken)
	return nil
}
