package types

import (
	"encoding/hex"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateFXParams{}
	_ sdk.Msg = &MsgUpdateEGFParams{}
	_ sdk.Msg = &MsgUpdateStore{}
	_ sdk.Msg = &MsgUpdateSwitchParams{}
)

func NewMsgUpdateFXParams(authority string, params Params) *MsgUpdateFXParams {
	return &MsgUpdateFXParams{Authority: authority, Params: params}
}

func (m *MsgUpdateFXParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("params err: %s", err.Error())
	}
	return nil
}

func NewMsgUpdateEGFParams(authority string, params EGFParams) *MsgUpdateEGFParams {
	return &MsgUpdateEGFParams{Authority: authority, Params: params}
}

func (m *MsgUpdateEGFParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("params err: %s", err.Error())
	}
	return nil
}

func NewMsgUpdateStore(authority string, updateStores []UpdateStore) *MsgUpdateStore {
	return &MsgUpdateStore{Authority: authority, UpdateStores: updateStores}
}

func (m *MsgUpdateStore) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if len(m.UpdateStores) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("stores are empty")
	}
	for _, updateStore := range m.UpdateStores {
		if len(updateStore.Space) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("store space is empty")
		}
		if len(updateStore.Key) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("store key is empty")
		}
		if _, err := hex.DecodeString(updateStore.Key); err != nil {
			return sdkerrors.ErrInvalidRequest.Wrap("invalid store key")
		}
		if len(updateStore.OldValue) > 0 {
			if _, err := hex.DecodeString(updateStore.OldValue); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrap("invalid old store value")
			}
		}
		if len(updateStore.Value) > 0 {
			if _, err := hex.DecodeString(updateStore.Value); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrap("invalid store value")
			}
		}
	}
	return nil
}

func (us *UpdateStore) String() string {
	out, _ := json.Marshal(us)
	return string(out)
}

func (us *UpdateStore) KeyToBytes() []byte {
	b, err := hex.DecodeString(us.Key)
	if err != nil {
		panic(err)
	}
	return b
}

func (us *UpdateStore) OldValueToBytes() []byte {
	if len(us.OldValue) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(us.OldValue)
	if err != nil {
		panic(err)
	}
	return b
}

func (us *UpdateStore) ValueToBytes() []byte {
	if len(us.Value) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(us.Value)
	if err != nil {
		panic(err)
	}
	return b
}

func NewMsgUpdateSwitchParams(authority string, params SwitchParams) *MsgUpdateSwitchParams {
	return &MsgUpdateSwitchParams{Authority: authority, Params: params}
}

func (m *MsgUpdateSwitchParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap("authority")
	}
	if err := m.Params.ValidateBasic(); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("params err: %s", err.Error())
	}
	return nil
}
