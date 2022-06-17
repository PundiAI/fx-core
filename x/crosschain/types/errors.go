package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid   = sdkerrors.Register(ModuleName, 2, "invalid")
	ErrEmpty     = sdkerrors.Register(ModuleName, 3, "empty")
	ErrUnknown   = sdkerrors.Register(ModuleName, 4, "unknown")
	ErrDuplicate = sdkerrors.Register(ModuleName, 5, "duplicate")

	ErrNonContiguousEventNonce = sdkerrors.Register(ModuleName, 6, "non contiguous event nonce")
	ErrNoFoundOracle           = sdkerrors.Register(ModuleName, 7, "no found oracle")
	ErrOracleNotOnLine         = sdkerrors.Register(ModuleName, 8, "oracle not on line")

	ErrDelegateAmountBelowMinimum = sdkerrors.Register(ModuleName, 9, "delegate amount must be greater than oracle stake threshold")
	ErrDelegateAmountBelowMaximum = sdkerrors.Register(ModuleName, 10, "delegate amount must be less than double oracle stake threshold")
)
