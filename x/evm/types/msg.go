package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

var _ sdk.Msg = &MsgCallContract{}

// GetSignBytes implements the LegacyMsg interface.
func (m *MsgCallContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(evmtypes.ModuleCdc.MustMarshalJSON(m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgCallContract) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgCallContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}
	if err := fxtypes.ValidateEthereumAddress(m.ContractAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid contract address: %s", err.Error())
	}
	if len(m.Data) == 0 {
		return errortypes.ErrInvalidAddress.Wrap("invalid data ")
	}
	return nil
}
