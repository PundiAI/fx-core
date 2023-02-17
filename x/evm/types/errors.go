package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/evmos/ethermint/x/evm/types"
)

const (
	codeErrABIPack   = 10001
	codeErrABIUnpack = 10002
)

var (
	ErrABIPack   = sdkerrors.Register(types.ModuleName, codeErrABIPack, "failed abi pack args")
	ErrABIUnpack = sdkerrors.Register(types.ModuleName, codeErrABIUnpack, "failed abi unpack data")
)
