package migrate

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/migrate/types"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(server types.MsgServer) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		//check module enable
		if ctx.BlockHeight() < fxtypes.EvmV1SupportBlock() {
			return nil, sdkerrors.Wrap(types.InvalidRequest, "migrate module not enable")
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgMigrateAccount:
			res, err := server.MigrateAccount(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized migrate msg type: %v", msg.Type()))
		}
	}
}
