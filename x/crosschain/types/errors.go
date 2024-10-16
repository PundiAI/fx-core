package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalid                 = errorsmod.Register(ModuleName, 2, "invalid")
	ErrNonContiguousEventNonce = errorsmod.Register(ModuleName, 6, "non contiguous event nonce")

	ErrNoFoundOracle   = errorsmod.Register(ModuleName, 7, "no found oracle")
	ErrOracleNotOnLine = errorsmod.Register(ModuleName, 8, "oracle not on line")

	ErrDelegateAmountBelowMinimum = errorsmod.Register(ModuleName, 9, "delegate amount must be greater than oracle stake threshold")
	ErrDelegateAmountAboveMaximum = errorsmod.Register(ModuleName, 10, "delegate amount must be less than double oracle stake threshold")
)
