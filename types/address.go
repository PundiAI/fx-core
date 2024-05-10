package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
)

func ParseAddress(addr string) (accAddr sdk.AccAddress, isEvmAddr bool, err error) {
	_, bytes, decodeErr := bech32.DecodeAndConvert(addr)
	if decodeErr == nil {
		return bytes, false, nil
	}
	ethAddrError := contract.ValidateEthereumAddress(addr)
	if ethAddrError == nil {
		return common.HexToAddress(addr).Bytes(), true, nil
	}
	return nil, false, errors.Join(decodeErr, ethAddrError)
}
