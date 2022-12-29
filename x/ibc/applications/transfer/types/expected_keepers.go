package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TransactionHook interface {
	TransferAfter(ctx sdk.Context, sender, receive string, coins, fee sdk.Coin) error
}

type RefundHook interface {
	RefundAfter(ctx sdk.Context, sourcePort, sourceChannel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error
}

type AckHook interface {
	AckAfter(ctx sdk.Context, sourcePort, sourceChannel string, sequence uint64) error
}
