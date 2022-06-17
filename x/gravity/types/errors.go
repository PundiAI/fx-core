package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid                 = sdkerrors.Register(ModuleName, 2, "invalid")
	ErrEmpty                   = sdkerrors.Register(ModuleName, 3, "empty")
	ErrUnknown                 = sdkerrors.Register(ModuleName, 4, "unknown")
	ErrDuplicate               = sdkerrors.Register(ModuleName, 5, "duplicate")
	ErrNonContiguousEventNonce = sdkerrors.Register(ModuleName, 6, "non contiguous event nonce")
)
