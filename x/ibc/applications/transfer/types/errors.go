package types

import (
	errorsmod "cosmossdk.io/errors"
)

var ErrMemoNotSupport = errorsmod.Register(CompatibleModuleName, 104, "memo not support")
