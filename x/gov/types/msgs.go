package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgUpdateEGFParams{}
)

const (
	TypeMsgUpdateParams    = "fx_update_params"
	TypeMsgUpdateEGFParams = "fx_update_egf_params"
)

// Route returns the MsgUpdateParams message route.
func (m *MsgUpdateParams) Route() string { return types.ModuleName }

// Type returns the MsgUpdateParams message type.
func (m *MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	return nil
}

// Route returns the MsgUpdateParams message route.
func (m *MsgUpdateEGFParams) Route() string { return types.ModuleName }

// Type returns the MsgUpdateParams message type.
func (m *MsgUpdateEGFParams) Type() string { return TypeMsgUpdateEGFParams }

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m *MsgUpdateEGFParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateEGFParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateEGFParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "params")
	}
	return nil
}
