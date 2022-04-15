package erc20

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/erc20/keeper"
	"github.com/functionx/fx-core/x/erc20/types"
)

// NewHandler defines the erc20 module handler instance
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		if ctx.BlockHeight() < fxtype.EvmSupportBlock() || !k.HasInit(ctx) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "erc20 module not enable")
		}
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgConvertCoin:
			res, err := k.ConvertCoin(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgConvertERC20:
			res, err := k.ConvertERC20(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			err := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, err
		}
	}
}
