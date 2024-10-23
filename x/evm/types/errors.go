package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/evmos/ethermint/x/evm/types"
)

var (
	ErrABIPack   = errorsmod.Register(types.ModuleName, 10001, "failed abi pack args")
	ErrABIUnpack = errorsmod.Register(types.ModuleName, 10002, "failed abi unpack data")
)
