package keeper

import (
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v8/x/tron/types"
)

func (k Keeper) ConfirmHandler(ctx sdk.Context, confirm types.Confirm) error {
	switch c := confirm.(type) {
	case *types.MsgConfirmBatch:
		if err := k.BatchConfirmHandler(ctx, c); err != nil {
			return err
		}

	case *types.MsgOracleSetConfirm:
		if err := k.OracleSetConfirmHandler(ctx, c); err != nil {
			return err
		}

	case *types.MsgBridgeCallConfirm:
		if err := k.BridgeCallConfirmHandler(ctx, c); err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, confirm.GetChainName()),
		sdk.NewAttribute(sdk.AttributeKeySender, confirm.GetBridgerAddress()),
	))
	return nil
}

func (k Keeper) BatchConfirmHandler(ctx sdk.Context, msg *types.MsgConfirmBatch) error {
	// fetch the outgoing batch given the nonce
	batch := k.GetOutgoingTxBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return errorsmod.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	var checkpoint []byte
	var err error
	if k.moduleName == trontypes.ModuleName {
		checkpoint, err = trontypes.GetCheckpointConfirmBatch(batch, k.GetGravityID(ctx))
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, "checkpoint generation")
		}
	} else {
		checkpoint, err = batch.GetCheckpoint(k.GetGravityID(ctx))
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, err.Error())
		}
	}

	oracleAddr, err := k.ValidateConfirmSign(ctx, msg.BridgerAddress, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return err
	}
	// check if we already have this confirm
	if k.GetBatchConfirm(ctx, msg.TokenContract, msg.Nonce, oracleAddr) != nil {
		return errorsmod.Wrap(types.ErrDuplicate, "signature")
	}
	k.SetBatchConfirm(ctx, oracleAddr, msg)
	return nil
}

func (k Keeper) OracleSetConfirmHandler(ctx sdk.Context, msg *types.MsgOracleSetConfirm) error {
	oracleSet := k.GetOracleSet(ctx, msg.Nonce)
	if oracleSet == nil {
		return errorsmod.Wrap(types.ErrInvalid, "couldn't find oracleSet")
	}

	var checkpoint []byte
	var err error
	if k.moduleName == trontypes.ModuleName {
		checkpoint, err = trontypes.GetCheckpointOracleSet(oracleSet, k.GetGravityID(ctx))
		if err != nil {
			return err
		}
	} else {
		checkpoint, err = oracleSet.GetCheckpoint(k.GetGravityID(ctx))
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, err.Error())
		}
	}

	oracleAddr, err := k.ValidateConfirmSign(ctx, msg.BridgerAddress, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return err
	}
	// check if we already have this confirm
	if k.GetOracleSetConfirm(ctx, msg.Nonce, oracleAddr) != nil {
		return errorsmod.Wrap(types.ErrDuplicate, "signature")
	}
	k.SetOracleSetConfirm(ctx, oracleAddr, msg)
	return nil
}

func (k Keeper) BridgeCallConfirmHandler(ctx sdk.Context, msg *types.MsgBridgeCallConfirm) error {
	outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(ctx, msg.Nonce)
	if !found {
		return errorsmod.Wrap(types.ErrInvalid, "couldn't find outgoing bridge call")
	}

	var checkpoint []byte
	var err error
	if k.moduleName == trontypes.ModuleName {
		checkpoint, err = trontypes.GetCheckpointBridgeCall(outgoingBridgeCall, k.GetGravityID(ctx))
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, err.Error())
		}
	} else {
		checkpoint, err = outgoingBridgeCall.GetCheckpoint(k.GetGravityID(ctx))
		if err != nil {
			return errorsmod.Wrap(types.ErrInvalid, err.Error())
		}
	}

	oracleAddr, err := k.ValidateConfirmSign(ctx, msg.BridgerAddress, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return err
	}

	if k.HasBridgeCallConfirm(ctx, msg.Nonce, oracleAddr) {
		return errorsmod.Wrap(types.ErrDuplicate, "signature")
	}
	k.SetBridgeCallConfirm(ctx, oracleAddr, msg)
	return nil
}

func (k Keeper) ValidateConfirmSign(ctx sdk.Context, bridgerAddr, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "signature decoding")
	}

	oracleAddr, found := k.GetOracleAddrByExternalAddr(ctx, signatureAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}

	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}

	if oracle.ExternalAddress != signatureAddr {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "got %s, expected %s", signatureAddr, oracle.ExternalAddress)
	}
	if oracle.BridgerAddress != bridgerAddr {
		return nil, errorsmod.Wrapf(types.ErrInvalid, "got %s, expected %s", bridgerAddr, oracle.BridgerAddress)
	}
	if k.moduleName == trontypes.ModuleName {
		if err = trontypes.ValidateTronSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
		}
	} else {
		if err = types.ValidateEthereumSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
		}
	}
	return oracleAddr, nil
}
