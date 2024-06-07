package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// AttestationHandler Handle is the entry point for Attestation processing.
//
//gocyclo:ignore
func (k Keeper) AttestationHandler(ctx sdk.Context, externalClaim types.ExternalClaim) error {
	switch claim := externalClaim.(type) {
	case *types.MsgSendToFxClaim:
		return k.SendToFxExecuted(ctx, claim)

	case *types.MsgSendToExternalClaim:
		k.OutgoingTxBatchExecuted(ctx, claim.TokenContract, claim.BatchNonce)

	case *types.MsgBridgeTokenClaim:
		return k.AddBridgeTokenExecuted(ctx, claim)

	case *types.MsgOracleSetUpdatedClaim:
		return k.UpdateOracleSetExecuted(ctx, claim)

	case *types.MsgBridgeCallClaim:
		return k.BridgeCallHandler(ctx, claim)

	case *types.MsgBridgeCallResultClaim:
		k.BridgeCallResultHandler(ctx, claim)

	default:
		return errorsmod.Wrapf(types.ErrInvalid, "event type: %s", claim.GetType())
	}
	return nil
}
