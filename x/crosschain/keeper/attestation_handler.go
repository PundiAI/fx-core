package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// AttestationHandler Handle is the entry point for Attestation processing.
//
//gocyclo:ignore
func (k Keeper) AttestationHandler(ctx sdk.Context, externalClaim types.ExternalClaim) error {
	switch claim := externalClaim.(type) {
	case *types.MsgSendToFxClaim:
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
			if err := k.bankKeeper.MintCoins(ctx, k.moduleName, sdk.NewCoins(coin)); err != nil {
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

		k.HandlePendingOutgoingTx(ctx, receiveAddr, externalClaim.GetEventNonce(), bridgeToken)

	case *types.MsgBridgeCallClaim:
		k.CreateBridgeAccount(ctx, claim.TxOrigin)
		return k.BridgeCallHandler(ctx, claim)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	case *types.MsgBridgeTokenClaim:
		// Check if it already exists
		if has := k.HasBridgeToken(ctx, claim.TokenContract); has {
			return types.ErrInvalid.Wrap("bridge token is exist")
		}

		k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token",
			claim.TokenContract, "channelIbc", claim.ChannelIbc)
		if claim.Symbol == types.NativeDenom {
			// Check if denom exists
			if !k.bankKeeper.HasDenomMetaData(ctx, claim.Symbol) {
				return types.ErrUnknown.Wrapf("denom not found %s", claim.Symbol)
			}

			if fxtypes.DenomUnit != uint32(claim.Decimals) {
				return types.ErrInvalid.Wrapf("%s denom decimals not match %d, expect %d", types.NativeDenom,
					claim.Decimals, fxtypes.DenomUnit)
			}

			k.AddBridgeToken(ctx, claim.TokenContract, types.NativeDenom)
			return nil
		}

		denom, err := k.SetIbcDenomTrace(ctx, claim.TokenContract, claim.ChannelIbc)
		if err != nil {
			return err
		}
		k.AddBridgeToken(ctx, claim.TokenContract, denom)

	case *types.MsgOracleSetUpdatedClaim:
		observedOracleSet := &types.OracleSet{
			Nonce:   claim.OracleSetNonce,
			Members: claim.Members,
		}
		// check the contents of the validator set against the store
		if claim.OracleSetNonce != 0 {
			trustedOracleSet := k.GetOracleSet(ctx, claim.OracleSetNonce)
			if trustedOracleSet == nil {
				ctx.Logger().Error("Received attestation for a oracle set which does not exist in store", "oracleSetNonce", claim.OracleSetNonce, "claim", claim)
				return errorsmod.Wrapf(types.ErrInvalid, "attested oracleSet (%v) does not exist in store", claim.OracleSetNonce)
			}
			// overwrite the height, since it's not part of the claim
			observedOracleSet.Height = trustedOracleSet.Height

			if _, err := trustedOracleSet.Equal(observedOracleSet); err != nil {
				panic(fmt.Sprintf("Potential bridge highjacking: observed oracleSet (%+v) does not match stored oracleSet (%+v)! %s", observedOracleSet, trustedOracleSet, err.Error()))
			}
		}
		k.SetLastObservedOracleSet(ctx, observedOracleSet)

	case *types.MsgBridgeCallResultClaim:
		k.CreateBridgeAccount(ctx, claim.TxOrigin)
		outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(ctx, claim.Nonce)
		if !found {
			panic(fmt.Errorf("bridge call not found for nonce %d", claim.Nonce))
		}
		if !claim.GetResult() {
			k.HandleOutgoingBridgeCallRefund(ctx, outgoingBridgeCall)
		}
		k.DeleteOutgoingBridgeCallRecord(ctx, claim.Nonce)
		return nil
	default:
		return errorsmod.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
