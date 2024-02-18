package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

type RefundHook interface {
	RefundAfter(ctx sdk.Context, sourceChannel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error
	AckAfter(ctx sdk.Context, sourceChannel string, sequence uint64) error
}

type Erc20Keeper interface {
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
	ConvertCoin(goCtx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
}
