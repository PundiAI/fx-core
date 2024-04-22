package types

import (
	"bytes"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

const (
	TypeMsgGrantPrivilege      = "grant_privilege"
	TypeMsgEditConsensusPubKey = "edit_consensus_pubkey"
)

var (
	_ sdk.Msg                            = &MsgGrantPrivilege{}
	_ codectypes.UnpackInterfacesMessage = (*MsgGrantPrivilege)(nil)
	_ sdk.Msg                            = &MsgEditConsensusPubKey{}
	_ codectypes.UnpackInterfacesMessage = (*MsgEditConsensusPubKey)(nil)
)

func NewMsgGrantPrivilege(val sdk.ValAddress, from sdk.AccAddress, pubKey cryptotypes.PubKey) (*MsgGrantPrivilege, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return nil, err
	}
	return &MsgGrantPrivilege{
		ValidatorAddress: val.String(),
		FromAddress:      from.String(),
		ToPubkey:         pkAny,
	}, nil
}

func (m *MsgGrantPrivilege) Route() string { return stakingtypes.RouterKey }

func (m *MsgGrantPrivilege) Type() string { return TypeMsgGrantPrivilege }

func (m *MsgGrantPrivilege) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	fromAddress, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	pk, err := ProtoAnyToAccountPubKey(m.ToPubkey)
	if err != nil {
		return err
	}
	toAddress := sdk.AccAddress(pk.Address())
	// check same account
	if bytes.Equal(fromAddress.Bytes(), toAddress.Bytes()) {
		return errortypes.ErrInvalidRequest.Wrap("same account")
	}
	return nil
}

func (m *MsgGrantPrivilege) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgGrantPrivilege) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.FromAddress)}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *MsgGrantPrivilege) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(m.ToPubkey, &pubKey)
}

func ProtoAnyToAccountPubKey(any *codectypes.Any) (cryptotypes.PubKey, error) {
	if any == nil {
		return nil, errortypes.ErrInvalidPubKey.Wrap("empty pubkey")
	}
	pk, ok := any.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errortypes.ErrInvalidPubKey.Wrapf("expecting cryptotypes.PubKey, got %T", any.GetCachedValue())
	}
	_, ok1 := pk.(*secp256k1.PubKey)
	_, ok2 := pk.(*ethsecp256k1.PubKey)
	if !ok1 && !ok2 {
		return nil, errortypes.ErrInvalidPubKey.Wrapf("expecting *secp256k1.PubKey or *ethsecp256k1.PubKey, got %T", pk)
	}
	return pk, nil
}

func GrantPrivilegeSignatureData(val, from, to []byte) []byte {
	prefixLen := len(GrantPrivilegeSignaturePrefix)
	valLen := len(val)
	fromLen := len(from)
	toLen := len(to)
	data := make([]byte, prefixLen+valLen+fromLen+toLen)
	copy(data[:prefixLen], GrantPrivilegeSignaturePrefix)
	copy(data[prefixLen:prefixLen+valLen], val)
	copy(data[prefixLen+valLen:prefixLen+valLen+fromLen], from)
	copy(data[prefixLen+valLen+fromLen:], to)
	return data
}

func NewMsgEditConsensusPubKey(val sdk.ValAddress, from sdk.AccAddress, pubKey cryptotypes.PubKey) (*MsgEditConsensusPubKey, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgEditConsensusPubKey{
		ValidatorAddress: val.String(),
		From:             from.String(),
		Pubkey:           pkAny,
	}, nil
}

func (m *MsgEditConsensusPubKey) Route() string { return stakingtypes.RouterKey }

func (m *MsgEditConsensusPubKey) Type() string { return TypeMsgEditConsensusPubKey }

func (m *MsgEditConsensusPubKey) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	if m.Pubkey == nil {
		return stakingtypes.ErrEmptyValidatorPubKey
	}
	return nil
}

func (m *MsgEditConsensusPubKey) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgEditConsensusPubKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *MsgEditConsensusPubKey) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(m.Pubkey, &pubKey)
}
