package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RefundHook interface {
	RefundAfter(ctx sdk.Context, sourceChannel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error
	AckAfter(ctx sdk.Context, sourceChannel string, sequence uint64) error
}
