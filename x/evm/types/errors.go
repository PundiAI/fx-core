package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/evmos/ethermint/x/evm/types"
)

const (
	codeErrABIPack   = 10001
	codeErrABIUnpack = 10002
)

var (
	ErrABIPack   = errorsmod.Register(types.ModuleName, codeErrABIPack, "failed abi pack args")
	ErrABIUnpack = errorsmod.Register(types.ModuleName, codeErrABIUnpack, "failed abi unpack data")
)
