package crosschain

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/crosschain/types"
	"github.com/gogo/protobuf/proto"
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
			if !moduleHandlerRouter.HasRoute(msg.GetChainName()) {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", msg.GetChainName()))
			}
			chainName = msg.GetChainName()
			msgServer = k.Router().GetRoute(msg.GetChainName()).MsgServer
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("not cross chain msg"))
		}
		ignoreCommitKeyNameMapByHeight := rootmulti.GetIgnoreCommitKeyNameMapByHeight(ctx.BlockHeight())
		if _, ok := ignoreCommitKeyNameMapByHeight[chainName]; ok {
			panic(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("module not enable")))
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("module not enable"))
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var res proto.Message
		var err error
		sdkCtx := sdk.WrapSDKContext(ctx)
		switch msg := msg.(type) {
		case *types.MsgSetOrchestratorAddress:
			res, err = msgServer.SetOrchestratorAddress(sdkCtx, msg)
		case *types.MsgAddOracleDeposit:
			res, err = msgServer.AddOracleDeposit(sdkCtx, msg)
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
			err = sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized bsc Msg type: %v", msg.Type()))
		}
		return sdk.WrapServiceResult(ctx, res, err)
	}
}

func NewCrossChainProposalHandler(k keeper.RouterKeeper) govtypes.Handler {
	moduleHandlerRouter := k.Router()
	return func(ctx sdk.Context, content govtypes.Content) error {
		ignoreCommitKeyNameMapByHeight := rootmulti.GetIgnoreCommitKeyNameMapByHeight(ctx.BlockHeight())
		switch c := content.(type) {
		case *types.InitCrossChainParamsProposal:
			if !moduleHandlerRouter.HasRoute(c.ChainName) {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", c.ChainName))
			}
			if _, ok := ignoreCommitKeyNameMapByHeight[c.ChainName]; ok {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("module not enable"))
			}
			return k.Router().GetRoute(c.ChainName).MsgServer.HandleInitCrossChainParamsProposal(ctx, c)
		case *types.UpdateChainOraclesProposal:
			if !moduleHandlerRouter.HasRoute(c.ChainName) {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type:%s", c.ChainName))
			}
			if _, ok := ignoreCommitKeyNameMapByHeight[c.ChainName]; ok {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("module not enable"))
			}
			return k.Router().GetRoute(c.ChainName).MsgServer.HandleUpdateChainOraclesProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}
