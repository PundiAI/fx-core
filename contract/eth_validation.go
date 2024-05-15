package contract

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// EthereumContractAddressLen is the length of contract address strings
	EthereumContractAddressLen = 42

	// EthereumAddressPrefix is the address prefix address
	EthereumAddressPrefix = "0x"
)

var ethereumAddressRegular = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

// IsEmptyHash returns true if the hash corresponds to an empty ethereum hex hash.
func IsEmptyHash(hash common.Hash) bool {
	return bytes.Equal(hash.Bytes(), common.Hash{}.Bytes())
}

// IsZeroEthAddress returns true if the address corresponds to an empty ethereum address.
func IsZeroEthAddress(address common.Address) bool {
	return bytes.Equal(address.Bytes(), common.Address{}.Bytes())
}

// ValidateEthereumAddress validates the ethereum address strings
func ValidateEthereumAddress(address string) error {
	if address == "" {
		return errors.New("empty")
	}
	if len(address) != EthereumContractAddressLen {
		return errors.New("wrong length")
	}
	if !ethereumAddressRegular.MatchString(address) {
		return errors.New("invalid format")
	}
	// add ethereum address checksum check 2021-09-02.
	if !common.IsHexAddress(address) {
		return errors.New("doesn't pass format validation")
	}
	expectAddress := common.HexToAddress(address).Hex()
	if expectAddress != address {
		return fmt.Errorf("mismatch expected: %s, got: %s", expectAddress, address)
	}
	return nil
}
