package types

import (
	"fmt"

	"github.com/fbsobreira/gotron-sdk/pkg/common"
)

// ExternalContractAddressLen is the length of contract address strings
const ExternalContractAddressLen = 34

// ValidateExternalAddress validates the ethereum address strings
func ValidateExternalAddress(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty")
	}
	if len(addr) != ExternalContractAddressLen {
		return fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", addr, len(addr), ExternalContractAddressLen)
	}

	tronAddress, err := common.DecodeCheck(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %s", addr)
	}
	expectAddress := common.EncodeCheck(tronAddress[:])
	if expectAddress != addr {
		return fmt.Errorf("invalid address got:%s, expected:%s", addr, expectAddress)
	}
	return nil
}
