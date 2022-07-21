package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/v2/types"
)

// ethereumAddrLessThan migrates the Ethereum address less than function
func ethereumAddrLessThan(e, o string) bool {
	return bytes.Compare([]byte(e)[:], []byte(o)[:]) == -1
}

/////////////////////////
//      ERC20Token     //
/////////////////////////

func NewERC20Token(amount sdk.Int, contract string) ERC20Token {
	return ERC20Token{Amount: amount, Contract: contract}
}

// ValidateBasic permforms stateless validation
func (m ERC20Token) ValidateBasic() error {
	if err := fxtypes.ValidateEthereumAddress(m.Contract); err != nil {
		return sdkerrors.Wrap(err, "invalid contract address")
	}
	if !m.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalid, "amount")
	}
	return nil
}
