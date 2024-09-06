package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/contract"
)

const (
	codeErrABIPack   = 10001
	codeErrABIUnpack = 10002
)

var (
	ErrABIPack   = errorsmod.Register(types.ModuleName, codeErrABIPack, "failed abi pack args")
	ErrABIUnpack = errorsmod.Register(types.ModuleName, codeErrABIUnpack, "failed abi unpack data")
)

func PackRetError(err error) ([]byte, error) {
	pack, _ := abi.Arguments{{Type: contract.TypeString}}.Pack(err.Error())
	return pack, err
}

func PackRetErrV2(err error) ([]byte, error) {
	pack, _ := contract.GetErrorABI().Pack("Error", err.Error())
	return pack, err
}
