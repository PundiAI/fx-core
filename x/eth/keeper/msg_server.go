package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	crosschainkeeper "github.com/functionx/fx-core/v6/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v6/x/crosschain/types"
)

var _ crosschaintypes.MsgServer = MsgServer{}

type MsgServer struct {
	crosschainkeeper.MsgServer
}

func NewMsgServerImpl(keeper crosschainkeeper.Keeper) crosschaintypes.MsgServer {
	return MsgServer{
		MsgServer: crosschainkeeper.MsgServer{Keeper: keeper},
	}
}

func (s MsgServer) SendToExternal(c context.Context, msg *crosschaintypes.MsgSendToExternal) (*crosschaintypes.MsgSendToExternalResponse, error) {
	return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "not supported")
}
