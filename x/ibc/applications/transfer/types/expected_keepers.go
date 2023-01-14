package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

type RefundHook interface {
	RefundAfter(ctx sdk.Context, sourceChannel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error
	AckAfter(ctx sdk.Context, sourceChannel string, sequence uint64) error
}

type Erc20Keeper interface {
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
}
