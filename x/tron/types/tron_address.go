package types

import (
	"errors"
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/common"
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
