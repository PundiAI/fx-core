package types

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ethereumAddressRegular = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

// ethereumContractAddressLen is the length of contract address strings
const ethereumContractAddressLen = 42

// ethereumAddrLessThan migrates the Ethereum address less than function
func ethereumAddrLessThan(e, o string) bool {
	return bytes.Compare([]byte(e)[:], []byte(o)[:]) == -1
}

// ValidateEthereumAddress validates the ethereum address strings
func ValidateEthereumAddress(addr string) error {
	if addr == "" {
		return ErrEmpty
	}
	if len(addr) != ethereumContractAddressLen {
		return sdkerrors.Wrap(ErrExternalAddress, fmt.Sprintf("address(%s) of the wrong length exp(%d) actual(%d)", addr, len(addr), ethereumContractAddressLen))
	}
	if !ethereumAddressRegular.MatchString(addr) {
		return sdkerrors.Wrap(ErrExternalAddress, fmt.Sprintf("address(%s) doesn't pass regex", addr))
	}
	// add ethereum address checksum check 2021-09-02.
	if !common.IsHexAddress(addr) {
		return sdkerrors.Wrap(ErrExternalAddress, fmt.Sprintf("invalid address: %s", addr))
	}
	expectAddress := common.HexToAddress(addr).Hex()
	if expectAddress != addr {
		return sdkerrors.Wrap(ErrExternalAddress, fmt.Sprintf("invalid address got:%s, expected:%s", addr, expectAddress))
	}
	return nil
}

/////////////////////////
//      ERC20Token     //
/////////////////////////

func NewERC20Token(amount sdk.Int, contract string) ERC20Token {
	return ERC20Token{Amount: amount, Contract: contract}
}

// ValidateBasic permforms stateless validation
func (m ERC20Token) ValidateBasic() error {
	if err := ValidateEthereumAddress(m.Contract); err != nil {
		return sdkerrors.Wrap(err, "invalid contract address")
	}
	if !m.Amount.IsPositive() {
		return errors.New("invalid amount")
	}
	return nil
}
