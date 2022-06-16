package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrDuplicate                  = sdkerrors.Register(ModuleName, 2, "duplicate")
	ErrInvalid                    = sdkerrors.Register(ModuleName, 3, "invalid")
	ErrUnknown                    = sdkerrors.Register(ModuleName, 5, "unknown")
	ErrEmpty                      = sdkerrors.Register(ModuleName, 6, "empty")
	ErrNonContiguousEventNonce    = sdkerrors.Register(ModuleName, 9, "non contiguous event nonce")
	ErrInvalidRequestBatchBaseFee = sdkerrors.Register(ModuleName, 10, "invalid request batch base fee")
)
