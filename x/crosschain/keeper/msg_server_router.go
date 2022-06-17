package keeper

import (
	"context"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

type msgServer struct {
	routerKeeper RouterKeeper
}

// NewMsgServerRouterImpl returns an implementation of the crosschain router MsgServer interface
// for the provided Keeper.
func NewMsgServerRouterImpl(routerKeeper RouterKeeper) types.MsgServer {
	return &msgServer{routerKeeper: routerKeeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) BondedOracle(ctx context.Context, msg *types.MsgBondedOracle) (*types.MsgBondedOracleResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.BondedOracle(ctx, msg)
	}
}

func (k msgServer) AddDelegate(ctx context.Context, msg *types.MsgAddDelegate) (*types.MsgAddDelegateResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.AddDelegate(ctx, msg)
	}
}

func (k msgServer) EditOracle(ctx context.Context, msg *types.MsgEditOracle) (*types.MsgEditOracleResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.EditOracle(ctx, msg)
	}
}

func (k msgServer) WithdrawReward(ctx context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.WithdrawReward(ctx, msg)
	}
}

func (k msgServer) UnbondedOracle(ctx context.Context, msg *types.MsgUnbondedOracle) (*types.MsgUnbondedOracleResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.UnbondedOracle(ctx, msg)
	}
}

func (k msgServer) OracleSetConfirm(ctx context.Context, msg *types.MsgOracleSetConfirm) (*types.MsgOracleSetConfirmResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.OracleSetConfirm(ctx, msg)
	}
}

func (k msgServer) OracleSetUpdateClaim(ctx context.Context, msg *types.MsgOracleSetUpdatedClaim) (*types.MsgOracleSetUpdatedClaimResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.OracleSetUpdateClaim(ctx, msg)
	}
}

func (k msgServer) BridgeTokenClaim(ctx context.Context, msg *types.MsgBridgeTokenClaim) (*types.MsgBridgeTokenClaimResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeTokenClaim(ctx, msg)
	}
}

func (k msgServer) SendToFxClaim(ctx context.Context, msg *types.MsgSendToFxClaim) (*types.MsgSendToFxClaimResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.SendToFxClaim(ctx, msg)
	}
}

func (k msgServer) SendToExternal(ctx context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.SendToExternal(ctx, msg)
	}
}

func (k msgServer) CancelSendToExternal(ctx context.Context, msg *types.MsgCancelSendToExternal) (*types.MsgCancelSendToExternalResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.CancelSendToExternal(ctx, msg)
	}
}

func (k msgServer) SendToExternalClaim(ctx context.Context, msg *types.MsgSendToExternalClaim) (*types.MsgSendToExternalClaimResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.SendToExternalClaim(ctx, msg)
	}
}

func (k msgServer) RequestBatch(ctx context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.RequestBatch(ctx, msg)
	}
}

func (k msgServer) ConfirmBatch(ctx context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	if queryServer, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return queryServer.ConfirmBatch(ctx, msg)
	}
}

func (k msgServer) getMsgServerByChainName(chainName string) (types.MsgServer, error) {
	msgServerRouter := k.routerKeeper.Router()
	if !msgServerRouter.HasRoute(chainName) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", chainName))
	}
	return msgServerRouter.GetRoute(chainName).MsgServer, nil
}
