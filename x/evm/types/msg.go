package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/contract"
)

var _ sdk.Msg = &MsgCallContract{}

// GetSignBytes implements the LegacyMsg interface.
func (m *MsgCallContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgCallContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgCallContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := contract.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return errorsmod.Wrap(err, "contract address")
	}
	if len(m.Data) == 0 {
		return errors.New("empty data")
	}
	return nil
}
