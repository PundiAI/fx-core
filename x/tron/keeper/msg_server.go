package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v7/x/tron/types"
)

var _ crosschaintypes.MsgServer = msgServer{}

type msgServer struct {
	crosschainkeeper.MsgServer
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) crosschaintypes.MsgServer {
	return &msgServer{crosschainkeeper.MsgServer{Keeper: keeper.Keeper}}
}

// ConfirmBatch handles MsgConfirmBatch
func (s msgServer) ConfirmBatch(c context.Context, msg *crosschaintypes.MsgConfirmBatch) (*crosschaintypes.MsgConfirmBatchResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := s.GetOutgoingTxBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "couldn't find batch")
	}

	checkpoint, err := trontypes.GetCheckpointConfirmBatch(batch, s.GetGravityID(ctx))
	if err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "checkpoint generation")
	}

	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.TokenContract, msg.Nonce, oracleAddr) != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrDuplicate, "signature")
	}
	s.SetBatchConfirm(ctx, oracleAddr, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return &crosschaintypes.MsgConfirmBatchResponse{}, nil
}

// OracleSetConfirm handles MsgOracleSetConfirm
func (s msgServer) OracleSetConfirm(c context.Context, msg *crosschaintypes.MsgOracleSetConfirm) (*crosschaintypes.MsgOracleSetConfirmResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleSet := s.GetOracleSet(ctx, msg.Nonce)
	if oracleSet == nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "couldn't find oracleSet")
	}

	checkpoint, err := trontypes.GetCheckpointOracleSet(oracleSet, s.GetGravityID(ctx))
	if err != nil {
		return nil, err
	}
	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetOracleSetConfirm(ctx, msg.Nonce, oracleAddr) != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrDuplicate, "signature")
	}
	s.SetOracleSetConfirm(ctx, oracleAddr, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return &crosschaintypes.MsgOracleSetConfirmResponse{}, nil
}

func (s msgServer) BridgeCallConfirm(c context.Context, msg *crosschaintypes.MsgBridgeCallConfirm) (*crosschaintypes.MsgBridgeCallConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// TODO bridge call record
	record, found := s.GetRefundRecord(ctx, msg.Nonce)
	if !found {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "couldn't find bridge call record")
	}

	snapshotOracle, found := s.GetSnapshotOracle(ctx, record.OracleSetNonce)
	if !found {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "couldn't find snapshot oracle")
	}
	if !snapshotOracle.HasExternalAddress(msg.ExternalAddress) {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "external address not in snapshot oracle")
	}

	checkpoint, err := trontypes.GetCheckpointBridgeCall(record, s.GetGravityID(ctx))
	if err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, err.Error())
	}
	sigBytes, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "signature decoding")
	}

	if err = trontypes.ValidateTronSignature(checkpoint, sigBytes, msg.ExternalAddress); err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", msg.ExternalAddress, hex.EncodeToString(checkpoint), sigBytes))
	}

	externalAddr := crosschaintypes.ExternalAddressToAccAddress(s.ModuleName(), msg.ExternalAddress)
	if _, found = s.GetBridgeCallConfirm(ctx, msg.Nonce, externalAddr); found {
		return nil, errorsmod.Wrap(crosschaintypes.ErrDuplicate, "signature")
	}
	s.SetBridgeCallConfirm(ctx, externalAddr, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return &crosschaintypes.MsgBridgeCallConfirmResponse{}, nil
}

func (s msgServer) confirmHandlerCommon(ctx sdk.Context, bridgerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, "signature decoding")
	}

	oracleAddr, found := s.GetOracleByExternalAddress(ctx, signatureAddr)
	if !found {
		return nil, crosschaintypes.ErrNoFoundOracle
	}

	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, crosschaintypes.ErrNoFoundOracle
	}

	if oracle.ExternalAddress != signatureAddr {
		return nil, errorsmod.Wrapf(crosschaintypes.ErrInvalid, "got %s, expected %s", signatureAddr, oracle.ExternalAddress)
	}
	if oracle.BridgerAddress != bridgerAddr.String() {
		return nil, errorsmod.Wrapf(crosschaintypes.ErrInvalid, "got %s, expected %s", bridgerAddr, oracle.BridgerAddress)
	}
	if err = trontypes.ValidateTronSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
		return nil, errorsmod.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return oracle.GetOracle(), nil
}
