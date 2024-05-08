package types

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
)

// TronContractAddressLen is the length of contract address strings
const TronContractAddressLen = 34

// ValidateTronAddress validates the ethereum address strings
func ValidateTronAddress(address string) error {
	if address == "" {
		return errors.New("empty")
	}
	if len(address) != TronContractAddressLen {
		return errors.New("wrong length")
	}

	tronAddr, err := common.DecodeCheck(address)
	if err != nil {
		return errors.New("doesn't pass format validation")
	}
	expectAddress := common.EncodeCheck(tronAddr)
	if expectAddress != address {
		return fmt.Errorf("mismatch expected: %s, got: %s", expectAddress, address)
	}
	return nil
}

var _ crosschaintypes.ExternalAddress = &TronAddress{}

type TronAddress struct{}

func (b TronAddress) ValidateExternalAddress(addr string) error {
	return ValidateTronAddress(addr)
}

func (b TronAddress) ExternalAddressToAccAddress(addr string) (sdk.AccAddress, error) {
	tronAddr, err := tronaddress.Base58ToAddress(addr)
	if err != nil {
		return nil, err
	}
	return tronAddr.Bytes()[1:], nil
}
