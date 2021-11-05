package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

// Handle is the entry point for Attestation processing.
func (k *Keeper) AttestationHandler(ctx sdk.Context, _ types.Attestation, externalClaim types.ExternalClaim) error {
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
			return sdkerrors.Wrap(err, "invalid receiver address")
		}
		if err := k.bankKeeper.MintCoins(ctx, k.moduleName, coins); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins")
		}
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, coins); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}

		event := sdk.NewEvent(
			types.EventTypeSendToFx,
			sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
			sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", claim.EventNonce)),
		)

		sourcePort, sourceChannel, nextChannelSendSequence, isOk := k.handleIbcTransfer(ctx, claim, receiveAddr, coin)
		if isOk {
			event = event.
				AppendAttributes(sdk.NewAttribute(types.AttributeKeyAttestationHandlerIbcChannelSendSequence, fmt.Sprintf("%d", nextChannelSendSequence))).
				AppendAttributes(sdk.NewAttribute(types.AttributeKeyAttestationHandlerIbcChannelSourcePort, sourcePort)).
				AppendAttributes(sdk.NewAttribute(types.AttributeKeyAttestationHandlerIbcChannelSourceChannel, sourceChannel))
			k.SetIbcSequenceHeight(ctx, sourcePort, sourceChannel, nextChannelSendSequence, uint64(ctx.BlockHeight()))
		}
		// broadcast event
		ctx.EventManager().EmitEvent(event)

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
		var coinDenom string
		var err error
		if coinDenom, err = k.addBridgeToken(ctx, claim.TokenContract, claim.Symbol, claim.ChannelIbc); err != nil {
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
