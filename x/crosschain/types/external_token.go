package types

import (
	"bytes"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"regexp"
)

var (
	ExternalAddressRegular = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

// ExternalContractAddressLen is the length of contract address strings
const ExternalContractAddressLen = 42

// ExternalAddrLessThan migrates the Ethereum address less than function
func ExternalAddrLessThan(e, o string) bool {
	return bytes.Compare([]byte(e)[:], []byte(o)[:]) == -1
}

// ValidateExternalAddress validates the ethereum address strings
func ValidateExternalAddress(addr string) error {
	if addr == "" {
		return fmt.Errorf("empty")
	}
	if len(addr) != ExternalContractAddressLen {
		return fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", addr, len(addr), ExternalContractAddressLen)
	}
	if !ExternalAddressRegular.MatchString(addr) {
		return fmt.Errorf("address(%s) doesn't pass regex", addr)
	}
	// add ethereum address checksum check 2021-09-02.
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("invalid address: %s", addr)
	}
	expectAddress := common.HexToAddress(addr).Hex()
	if expectAddress != addr {
		return fmt.Errorf("invalid address got:%s, expected:%s", addr, expectAddress)
	}
	return nil
}

/////////////////////////
//     ExternalToken   //
/////////////////////////

func NewExternalToken(amount uint64, contract string) *ExternalToken {
	return &ExternalToken{Amount: sdk.NewIntFromUint64(amount), Contract: contract}
}

func NewExternalTokenBySdkInt(amount sdk.Int, contract string) *ExternalToken {
	return &ExternalToken{Amount: amount, Contract: contract}
}

// ValidateBasic permforms stateless validation
func (m *ExternalToken) ValidateBasic() error {
	if err := ValidateExternalAddress(m.Contract); err != nil {
		return sdkerrors.Wrap(err, "ethereum address")
	}
	if !m.Amount.IsPositive() {
		return errors.New("invalid amount")
	}
	return nil
}
