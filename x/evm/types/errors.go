package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func PackRetError(str string) ([]byte, error) {
	pack, _ := abi.Arguments{{Type: TypeString}}.Pack(str)
	return pack, errors.New(str)
}
