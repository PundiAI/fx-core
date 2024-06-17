package types

import (
	"encoding/hex"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgUpdateEGFParams{}
	_ sdk.Msg = &MsgUpdateStore{}
)

const (
	TypeMsgUpdateParams    = "fx_update_params"
	TypeMsgUpdateEGFParams = "fx_update_egf_params"
	TypeMsgUpdateStore     = "fx_update_store"
)

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{Authority: authority, Params: params}
}

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
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
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

func NewMsgUpdateEGFParams(authority string, params EGFParams) *MsgUpdateEGFParams {
	return &MsgUpdateEGFParams{Authority: authority, Params: params}
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
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
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

func NewMsgUpdateStore(authority string, stores []UpdateStore) *MsgUpdateStore {
	return &MsgUpdateStore{Authority: authority, Stores: stores}
}

func (m *MsgUpdateStore) Route() string { return types.ModuleName }

func (m *MsgUpdateStore) Type() string { return TypeMsgUpdateStore }

func (m *MsgUpdateStore) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpdateStore) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m *MsgUpdateStore) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "authority")
	}
	if len(m.Stores) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("submitted store changes are empty")
	}
	for _, s := range m.Stores {
		if len(s.Space) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("store space is empty")
		}
		if len(s.Key) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("store key is empty")
		}
		if _, err := hex.DecodeString(s.Key); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid store key")
		}
		if len(s.OldValue) > 0 {
			if _, err := hex.DecodeString(s.OldValue); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrap("invalid old store value")
			}
		}
		if len(s.Value) > 0 {
			if _, err := hex.DecodeString(s.Value); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrap("invalid store value")
			}
		}
	}
	return nil
}

func (us UpdateStore) String() string {
	out, _ := json.Marshal(us)
	return string(out)
}

func (us UpdateStore) KeyToBytes() []byte {
	b, err := hex.DecodeString(us.Key)
	if err != nil {
		panic(err)
	}
	return b
}

func (us UpdateStore) OldValueToBytes() []byte {
	if len(us.OldValue) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(us.OldValue)
	if err != nil {
		panic(err)
	}
	return b
}

func (us UpdateStore) ValueToBytes() []byte {
	if len(us.Value) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(us.Value)
	if err != nil {
		panic(err)
	}
	return b
}
