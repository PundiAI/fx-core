package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

// AttestationHandler Handle is the entry point for Attestation processing.
func (k Keeper) AttestationHandler(ctx sdk.Context, externalClaim types.ExternalClaim) error {
	switch claim := externalClaim.(type) {
	case *types.MsgSendToFxClaim, *types.MsgBridgeCallClaim, *types.MsgBridgeCallResultClaim:
		k.SavePendingExecuteClaim(ctx, externalClaim)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	case *types.MsgBridgeTokenClaim:
		return k.AddBridgeTokenExecuted(ctx, claim)

	case *types.MsgOracleSetUpdatedClaim:
		return k.UpdateOracleSetExecuted(ctx, claim)

	default:
		return errorsmod.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}

func (k Keeper) ExecuteClaim(ctx sdk.Context, eventNonce uint64) error {
	externalClaim, found := k.GetPendingExecuteClaim(ctx, eventNonce)
	if !found {
		return errortypes.ErrInvalidRequest.Wrap("claim not found")
	}
	k.DeletePendingExecuteClaim(ctx, eventNonce)
	switch claim := externalClaim.(type) {
	case *types.MsgSendToFxClaim:
		return k.SendToFxExecuted(ctx, claim)
	case *types.MsgBridgeCallClaim:
		return k.BridgeCallHandler(ctx, claim)
	case *types.MsgBridgeCallResultClaim:
		k.BridgeCallResultHandler(ctx, claim)
	default:
		return errortypes.ErrInvalidRequest.Wrapf("invalid claim type: %s", claim.GetType())
	}
	return nil
}
