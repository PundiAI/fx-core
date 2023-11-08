package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (s MsgServer) SendToExternal(goCtx context.Context, msg *crosschaintypes.MsgSendToExternal) (*crosschaintypes.MsgSendToExternalResponse, error) {
	if sdk.UnwrapSDKContext(goCtx).BlockHeight() > 1e7 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unsupported method")
	}
	return s.MsgServer.SendToExternal(goCtx, msg)
}
