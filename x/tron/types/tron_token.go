package types

import (
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

// TronContractAddressLen is the length of contract address strings
const TronContractAddressLen = 34

// ValidateTronAddress validates the ethereum address strings
func ValidateTronAddress(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty")
	}
	if len(addr) != TronContractAddressLen {
		return fmt.Errorf("invalid address (%s) of the wrong length exp (%d) actual (%d)", addr, len(addr), TronContractAddressLen)
	}

	tronAddress, err := common.DecodeCheck(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %s", addr)
	}
	expectAddress := common.EncodeCheck(tronAddress[:])
	if expectAddress != addr {
		return fmt.Errorf("invalid address got: %s, expected: %s", addr, expectAddress)
	}
	return nil
}
