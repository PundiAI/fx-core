package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrFeeDenomNotMatchTokenDenom = sdkerrors.Register(ModuleName, 100, "invalid fee denom, must match token denom")
)
