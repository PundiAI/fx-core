package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
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
		receiveAddr, err := sdk.AccAddressFromBech32(claim.Receiver)
		if err != nil {
			return sdkerrors.Wrap(types.ErrInvalid, "receiver address")
		}
		if bridgeToken.Denom != fxtypes.DefaultDenom {
			// If it is not fxcore originated, mint the coins (aka vouchers)
			if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
				return sdkerrors.Wrapf(err, "mint vouchers coins")
			}
		}
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.moduleName, receiveAddr, sdk.NewCoins(coin)); err != nil {
			return sdkerrors.Wrap(err, "transfer vouchers")
		}

		k.HandlerRelayTransfer(ctx, claim.EventNonce, claim.TargetIbc, receiveAddr, coin)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	case *types.MsgBridgeTokenClaim:
		// Check if it already exists
		isExist := k.HasBridgeToken(ctx, claim.TokenContract)
		if isExist {
			return sdkerrors.Wrap(types.ErrInvalid, "bridge token is exist")
		}
		k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token", claim.TokenContract, "channelIbc", claim.ChannelIbc)
		if claim.Symbol == fxtypes.DefaultDenom {
			// Check if denom exists
			metadata, found := k.bankKeeper.GetDenomMetaData(ctx, claim.Symbol)
			if !found {
				return sdkerrors.Wrap(
					types.ErrUnknown,
					fmt.Sprintf("denom not found %s", claim.Symbol))
			}

			// Check if attributes of ERC20 match fx denom
			if claim.Name != metadata.Name {
				return sdkerrors.Wrap(
					types.ErrInvalid,
					fmt.Sprintf("ERC20 name %s does not match denom display %s", claim.Name, metadata.Description))
			}

			if claim.Symbol != metadata.Symbol {
				return sdkerrors.Wrap(
					types.ErrInvalid,
					fmt.Sprintf("ERC20 symbol %s does not match denom display %s", claim.Symbol, metadata.Display))
			}

			if fxtypes.DenomUnit != uint32(claim.Decimals) {
				return sdkerrors.Wrap(
					types.ErrInvalid,
					fmt.Sprintf("ERC20 decimals %d does not match denom decimals %d", claim.Decimals, fxtypes.DenomUnit))
			}
			k.AddBridgeToken(ctx, claim.TokenContract, claim.Symbol)
			return nil
		}

		denom, err := k.SetIbcDenomTrace(ctx, claim.TokenContract, claim.ChannelIbc)
		if err != nil {
			return err
		}
		k.AddBridgeToken(ctx, claim.TokenContract, denom)
		k.Logger(ctx).Info("add bridge token success", "symbol", claim.Symbol, "token", claim.TokenContract, "denom", denom, "channelIbc", claim.ChannelIbc)

	case *types.MsgOracleSetUpdatedClaim:
		observedOracleSet := types.OracleSet{
			Nonce:   claim.OracleSetNonce,
			Members: claim.Members,
		}
		// check the contents of the validator set against the store
		if claim.OracleSetNonce != 0 {
			trustedOracleSet := k.GetOracleSet(ctx, claim.OracleSetNonce)
			if trustedOracleSet == nil {
				ctx.Logger().Error("Received attestation for a oracle set which does not exist in store", "oracleSetNonce", claim.OracleSetNonce, "claim", claim)
				return sdkerrors.Wrapf(types.ErrInvalid, "attested oracleSet (%v) does not exist in store", claim.OracleSetNonce)
			}
			// overwrite the height, since it's not part of the claim
			observedOracleSet.Height = trustedOracleSet.Height

			if _, err := trustedOracleSet.Equal(observedOracleSet); err != nil {
				panic(fmt.Sprintf("Potential bridge highjacking: observed oracleSet (%+v) does not match stored oracleSet (%+v)! %s", observedOracleSet, trustedOracleSet, err.Error()))
			}
		}
		k.SetLastObservedOracleSet(ctx, observedOracleSet)
	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
