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
	ExternalAddressRegular = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

// ExternalContractAddressLen is the length of contract address strings
const ExternalContractAddressLen = 42

// IsEmptyHash returns true if the hash corresponds to an empty ethereum hex hash.
func IsEmptyHash(hash string) bool {
	return bytes.Equal(common.HexToHash(hash).Bytes(), common.Hash{}.Bytes())
}

// IsZeroAddress returns true if the address corresponds to an empty ethereum hex address.
func IsZeroAddress(address string) bool {
	return bytes.Equal(common.HexToAddress(address).Bytes(), common.Address{}.Bytes())
}

// ValidateAddress validates the ethereum address strings
func ValidateAddress(address string) error {
	if address == "" {
		return fmt.Errorf("empty")
	}
	if len(address) != ExternalContractAddressLen {
		return fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", address, len(address), ExternalContractAddressLen)
	}
	if !ExternalAddressRegular.MatchString(address) {
		return fmt.Errorf("address(%s) doesn't pass regex", address)
	}
	// add ethereum address checksum check 2021-09-02.
	if !common.IsHexAddress(address) {
		return fmt.Errorf("invalid address: %s", address)
	}
	expectAddress := common.HexToAddress(address).Hex()
	if expectAddress != address {
		return fmt.Errorf("invalid address got:%s, expected:%s", address, expectAddress)
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
