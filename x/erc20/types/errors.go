package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrERC20Disabled          = errorsmod.Register(ModuleName, 2, "erc20 module is disabled")
	ErrUndefinedOwner         = errorsmod.Register(ModuleName, 6, "undefined owner of erc20 contract")
	ErrERC20TokenPairDisabled = errorsmod.Register(ModuleName, 11, "erc20 token is disabled")
)
