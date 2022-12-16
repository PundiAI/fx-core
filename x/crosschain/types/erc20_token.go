package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

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
