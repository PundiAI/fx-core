package types

import (
	"fmt"

	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
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

	tronAddr, err := common.DecodeCheck(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %s", addr)
	}
	expectAddress := common.EncodeCheck(tronAddr[:])
	if expectAddress != addr {
		return fmt.Errorf("invalid address got: %s, expected: %s", addr, expectAddress)
	}
	return nil
}

func AddressFromHex(str string) string {
	bytes, _ := common.FromHex(str)
	return tronAddress.Address(append([]byte{tronAddress.TronBytePrefix}, bytes...)).String()
}
