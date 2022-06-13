package crosschain

import (
	"fmt"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/gogo/protobuf/proto"

	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/crosschain/types"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k keeper.RouterKeeper) sdk.Handler {
	moduleHandlerRouter := k.Router()
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {

		// check is cross chain msg.
		var msgServer types.MsgServer
		var chainName string
		switch msg := msg.(type) {
		case types.CrossChainMsg:
			chainName = msg.GetChainName()
			if !moduleHandlerRouter.HasRoute(chainName) {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", chainName))
			}
			msgServer = k.Router().GetRoute(chainName).MsgServer
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "not cross chain msg")
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var res proto.Message
		var err error
		sdkCtx := sdk.WrapSDKContext(ctx)
		switch msg := msg.(type) {
		case *types.MsgCreateOracleBridger:
			res, err = msgServer.CreateOracleBridger(sdkCtx, msg)
		case *types.MsgAddOracleDelegate:
			res, err = msgServer.AddOracleDelegate(sdkCtx, msg)
		case *types.MsgOracleSetConfirm:
			res, err = msgServer.OracleSetConfirm(sdkCtx, msg)
		case *types.MsgOracleSetUpdatedClaim:
			res, err = msgServer.OracleSetUpdateClaim(sdkCtx, msg)
		case *types.MsgBridgeTokenClaim:
			res, err = msgServer.BridgeTokenClaim(sdkCtx, msg)

		case *types.MsgSendToFxClaim:
			res, err = msgServer.SendToFxClaim(sdkCtx, msg)

		case *types.MsgSendToExternalClaim:
			res, err = msgServer.SendToExternalClaim(sdkCtx, msg)
		case *types.MsgSendToExternal:
			res, err = msgServer.SendToExternal(sdkCtx, msg)
		case *types.MsgCancelSendToExternal:
			res, err = msgServer.CancelSendToExternal(sdkCtx, msg)

		case *types.MsgRequestBatch:
			res, err = msgServer.RequestBatch(sdkCtx, msg)
		case *types.MsgConfirmBatch:
			res, err = msgServer.ConfirmBatch(sdkCtx, msg)
		default:
			err = sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized %s Msg type: %T", chainName, msg))
		}
		return sdk.WrapServiceResult(ctx, res, err)
	}
}

func NewCrossChainProposalHandler(k keeper.RouterKeeper) govtypes.Handler {
	moduleHandlerRouter := k.Router()
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UpdateCrossChainOraclesProposal:
			if !moduleHandlerRouter.HasRoute(c.ChainName) {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", c.ChainName))
			}
			return keeper.HandleUpdateCrossChainOraclesProposal(ctx, k.Router().GetRoute(c.ChainName).MsgServer, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}
