package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
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
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.BondedOracle(ctx, msg)
	}
}

func (k msgServer) AddDelegate(ctx context.Context, msg *types.MsgAddDelegate) (*types.MsgAddDelegateResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.AddDelegate(ctx, msg)
	}
}

func (k msgServer) ReDelegate(ctx context.Context, msg *types.MsgReDelegate) (*types.MsgReDelegateResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.ReDelegate(ctx, msg)
	}
}

func (k msgServer) EditBridger(ctx context.Context, msg *types.MsgEditBridger) (*types.MsgEditBridgerResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.EditBridger(ctx, msg)
	}
}

func (k msgServer) WithdrawReward(ctx context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.WithdrawReward(ctx, msg)
	}
}

func (k msgServer) UnbondedOracle(ctx context.Context, msg *types.MsgUnbondedOracle) (*types.MsgUnbondedOracleResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.UnbondedOracle(ctx, msg)
	}
}

func (k msgServer) OracleSetConfirm(ctx context.Context, msg *types.MsgOracleSetConfirm) (*types.MsgOracleSetConfirmResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.OracleSetConfirm(ctx, msg)
	}
}

func (k msgServer) SendToExternal(ctx context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.SendToExternal(ctx, msg)
	}
}

func (k msgServer) ConfirmBatch(ctx context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.ConfirmBatch(ctx, msg)
	}
}

func (k msgServer) BridgeCallConfirm(ctx context.Context, msg *types.MsgBridgeCallConfirm) (*types.MsgBridgeCallConfirmResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.BridgeCallConfirm(ctx, msg)
	}
}

func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.UpdateParams(ctx, msg)
	}
}

func (k msgServer) UpdateChainOracles(ctx context.Context, msg *types.MsgUpdateChainOracles) (*types.MsgUpdateChainOraclesResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.UpdateChainOracles(ctx, msg)
	}
}

func (k msgServer) getMsgServerByChainName(chainName string) (types.MsgServer, error) {
	msgServerRouter := k.routerKeeper.Router()
	if !msgServerRouter.HasRoute(chainName) {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("Unrecognized crosschain type:%s", chainName)
	}
	return msgServerRouter.GetRoute(chainName).MsgServer, nil
}

func (k msgServer) Claim(ctx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.Claim(ctx, msg)
	}
}

func (k msgServer) Confirm(ctx context.Context, msg *types.MsgConfirm) (*types.MsgConfirmResponse, error) {
	if server, err := k.getMsgServerByChainName(msg.GetChainName()); err != nil {
		return nil, err
	} else {
		return server.Confirm(ctx, msg)
	}
}
