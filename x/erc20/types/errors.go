package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrDisabled       = errorsmod.Register(ModuleName, 2, "erc20 module is disabled")
	ErrExists         = errorsmod.Register(ModuleName, 3, "token already exists")
	ErrUndefinedOwner = errorsmod.Register(ModuleName, 4, "undefined owner of erc20 contract")
)
