package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrFeeDenomNotMatchTokenDenom = errorsmod.Register(ModuleName, 100, "invalid fee denom, must match token denom")
	ErrRouterNotFound             = errorsmod.Register(ModuleName, 103, "router not found")
	ErrMemoNotSupport             = errorsmod.Register(ModuleName, 104, "memo not support")
)
