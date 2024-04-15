package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrERC20Disabled          = errorsmod.Register(ModuleName, 2, "erc20 module is disabled")
	ErrInternalTokenPair      = errorsmod.Register(ModuleName, 3, "internal ethereum token mapping error")
	ErrTokenPairNotFound      = errorsmod.Register(ModuleName, 4, "token pair not found")
	ErrTokenPairAlreadyExists = errorsmod.Register(ModuleName, 5, "token pair already exists")
	ErrUndefinedOwner         = errorsmod.Register(ModuleName, 6, "undefined owner of contract pair")
	ErrUnexpectedEvent        = errorsmod.Register(ModuleName, 7, "unexpected event")
	ErrABIUnpack              = errorsmod.Register(ModuleName, 9, "contract ABI unpack failed")
	ErrInvalidMetadata        = errorsmod.Register(ModuleName, 10, "invalid metadata")
	ErrERC20TokenPairDisabled = errorsmod.Register(ModuleName, 11, "erc20 token pair is disabled")
	ErrInvalidDenom           = errorsmod.Register(ModuleName, 12, "invalid denom")
	ErrInvalidAlias           = errorsmod.Register(ModuleName, 15, "invalid alias")
	ErrInsufficientLiquidity  = errorsmod.Register(ModuleName, 16, "insufficient liquidity")

	// Deprecated
	// ErrABIPack                = errorsmod.Register(ModuleName, 8, "contract ABI pack failed")
	// ErrInvalidTarget          = errorsmod.Register(ModuleName, 13, "invalid target")
	// ErrInternalRouter         = errorsmod.Register(ModuleName, 14, "internal router error")
)

func IsInsufficientLiquidityErr(err error) bool {
	return errorsmod.IsOf(err, ErrInsufficientLiquidity)
}
