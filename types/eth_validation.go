package types

import (
	"bytes"
	"fmt"
	"math"
	"regexp"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	EthereumAddressRegular = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

// EthereumContractAddressLen is the length of contract address strings
const EthereumContractAddressLen = 42

// IsEmptyHash returns true if the hash corresponds to an empty ethereum hex hash.
func IsEmptyHash(hash string) bool {
	return bytes.Equal(common.HexToHash(hash).Bytes(), common.Hash{}.Bytes())
}

// IsZeroEthereumAddress returns true if the address corresponds to an empty ethereum hex address.
func IsZeroEthereumAddress(address string) bool {
	return bytes.Equal(common.HexToAddress(address).Bytes(), common.Address{}.Bytes())
}

// ValidateEthereumAddress validates the ethereum address strings
func ValidateEthereumAddress(address string) error {
	if address == "" {
		return fmt.Errorf("empty")
	}
	if len(address) != EthereumContractAddressLen {
		return fmt.Errorf("invalid address (%s) of the wrong length exp (%d) actual (%d)", address, len(address), EthereumContractAddressLen)
	}
	if !EthereumAddressRegular.MatchString(address) {
		return fmt.Errorf("invalid address (%s) doesn't pass regex", address)
	}
	// add ethereum address checksum check 2021-09-02.
	if !common.IsHexAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}
	ethAddress := common.HexToAddress(address).Hex()
	if ethAddress != address {
		return fmt.Errorf("invalid address got: %s, expected: %s", address, ethAddress)
	}
	return nil
}

// SafeInt64 checks for overflows while casting a uint64 to int64 value.
func SafeInt64(value uint64) (int64, error) {
	if value > uint64(math.MaxInt64) {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrInvalidHeight, "uint64 value %v cannot exceed %v", value, int64(math.MaxInt64))
	}

	return int64(value), nil
}
