package gravity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	"github.com/functionx/fx-core/v3/x/gravity/keeper"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k crosschainkeeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgValsetConfirm:
			// nolint
			res, err := msgServer.ValsetConfirm(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSendToEth:
			// nolint
			res, err := msgServer.SendToEth(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRequestBatch:
			// nolint
			res, err := msgServer.RequestBatch(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgConfirmBatch:
			// nolint
			res, err := msgServer.ConfirmBatch(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDepositClaim:
			// nolint
			res, err := msgServer.DepositClaim(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgWithdrawClaim:
			// nolint
			res, err := msgServer.WithdrawClaim(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCancelSendToEth:
			// nolint
			res, err := msgServer.CancelSendToEth(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgValsetUpdatedClaim:
			// nolint
			res, err := msgServer.ValsetUpdateClaim(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Gravity Msg type: %T", msg))
		}
	}
}
