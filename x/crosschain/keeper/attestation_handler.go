package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/v2/x/crosschain/types"
)

// AttestationHandler Handle is the entry point for Attestation processing.
func (k Keeper) AttestationHandler(ctx sdk.Context, externalClaim types.ExternalClaim) error {
	switch claim := externalClaim.(type) {
	case *types.MsgSendToFxClaim:
		bridgeToken := k.GetBridgeTokenDenom(ctx, claim.TokenContract)
		if bridgeToken == nil {
			return sdkerrors.Wrap(types.ErrInvalid, "bridge token is not exist")
		}

		coin := sdk.NewCoin(bridgeToken.Denom, claim.Amount)
		coins := sdk.NewCoins(coin)

		receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalid, "receiver address")
		}
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, coins); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins")
		}
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, coins); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}

		k.handlerRelayTransfer(ctx, claim, receiveAddr, coin)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)
		return nil

	case *types.MsgBridgeTokenClaim:
		// Check if it already exists
		isExist := k.hasBridgeToken(ctx, claim.TokenContract)
		if isExist {
			return sdkerrors.Wrap(types.ErrInvalid, "bridge token is exist")
		}
		k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token", claim.TokenContract, "channelIbc", claim.ChannelIbc)

		coinDenom, err := k.AddBridgeToken(ctx, claim.TokenContract, claim.ChannelIbc)
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("add bridge token success", "symbol", claim.Symbol, "token", claim.TokenContract, "channelIbc", claim.ChannelIbc, "coinDenom", coinDenom)

	case *types.MsgOracleSetUpdatedClaim:
		k.SetLastObservedOracleSet(ctx, types.OracleSet{
			Nonce:   claim.OracleSetNonce,
			Members: claim.Members,
		})

	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
