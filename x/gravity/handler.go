package gravity

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/gravity/keeper"
	"github.com/functionx/fx-core/x/gravity/types"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewLegacyMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgSendToEth:
			res, err := msgServer.SendToEth(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized %s Msg type: %T", types.ModuleName, msg))
		}
	}
}
