package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"

	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

// TronContractAddressLen is the length of contract address strings
const TronContractAddressLen = 34

// ValidateTronAddress validates the ethereum address strings
func ValidateTronAddress(addr string) error {
	if addr == "" {
		return crosschaintypes.ErrEmpty
	}
	if len(addr) != TronContractAddressLen {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("address (%s) of the wrong length exp (%d) actual (%d)", addr, len(addr), TronContractAddressLen))
	}

	tronAddress, err := common.DecodeCheck(addr)
	if err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("address: %s", addr))
	}
	expectAddress := common.EncodeCheck(tronAddress[:])
	if expectAddress != addr {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("address got: %s, expected: %s", addr, expectAddress))
	}
	return nil
}
