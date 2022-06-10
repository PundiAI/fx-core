package keeper

import (
	"context"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/x/eth/types"
	"github.com/functionx/fx-core/x/gravity/types"
)

type legacyMsgServer interface {
	SendToExternal(c context.Context, msg *crosschaintypes.MsgSendToExternal) (*crosschaintypes.MsgSendToExternalResponse, error)
}

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewLegacyMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SendToEth(ctx context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	_, err := k.legacyMsgServer.SendToExternal(ctx, &crosschaintypes.MsgSendToExternal{
		Sender:    msg.Sender,
		Dest:      msg.EthDest,
		Amount:    msg.Amount,
		BridgeFee: msg.BridgeFee,
		ChainName: ethtypes.ModuleName,
	})
	if err != nil {
		return nil, err
	}
	return &types.MsgSendToEthResponse{}, nil
}
