package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	crosschainkeeper "github.com/functionx/fx-core/v2/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v2/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v2/x/tron/types"
)

var _ crosschaintypes.MsgServer = TronMsgServer{}

type TronMsgServer struct {
	crosschainkeeper.EthereumMsgServer
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) crosschaintypes.MsgServer {
	return &TronMsgServer{crosschainkeeper.EthereumMsgServer{Keeper: keeper.Keeper}}
}

// ConfirmBatch handles MsgConfirmBatch
func (s TronMsgServer) ConfirmBatch(c context.Context, msg *crosschaintypes.MsgConfirmBatch) (*crosschaintypes.MsgConfirmBatchResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := s.GetOutgoingTxBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "couldn't find batch")
	}

	checkpoint, err := trontypes.GetCheckpointConfirmBatch(batch, s.GetGravityID(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "checkpoint generation")
	}

	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.Nonce, msg.TokenContract, oracleAddr) != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrDuplicate, "signature")
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
func (s TronMsgServer) OracleSetConfirm(c context.Context, msg *crosschaintypes.MsgOracleSetConfirm) (*crosschaintypes.MsgOracleSetConfirmResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleSet := s.GetOracleSet(ctx, msg.Nonce)
	if oracleSet == nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "couldn't find oracleSet")
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
		return nil, sdkerrors.Wrap(crosschaintypes.ErrDuplicate, "signature")
	}
	s.SetOracleSetConfirm(ctx, oracleAddr, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return &crosschaintypes.MsgOracleSetConfirmResponse{}, nil
}

func (s TronMsgServer) confirmHandlerCommon(ctx sdk.Context, bridgerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "signature decoding")
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
		return nil, sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "got %s, expected %s", signatureAddr, oracle.ExternalAddress)
	}
	if oracle.BridgerAddress != bridgerAddr.String() {
		return nil, sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "got %s, expected %s", bridgerAddr, oracle.BridgerAddress)
	}
	if err = trontypes.ValidateTronSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return oracle.GetOracle(), nil
}
