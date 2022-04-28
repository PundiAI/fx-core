package types

import (
	"bytes"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

const TypeMsgMigrateAccount = "migrate_account"

var (
	_ sdk.Msg = &MsgMigrateAccount{}
)

// NewMsgMigrateAccount returns a new MsgMigrateAccount
func NewMsgMigrateAccount(from, to sdk.AccAddress, signature string) *MsgMigrateAccount {
	return &MsgMigrateAccount{
		From:      from.String(),
		To:        to.String(),
		Signature: signature,
	}
}

// Route should return the name of the module
func (m *MsgMigrateAccount) Route() string { return RouterKey }

// Type should return the action
func (m *MsgMigrateAccount) Type() string { return TypeMsgMigrateAccount }

// ValidateBasic performs stateless checks
func (m *MsgMigrateAccount) ValidateBasic() error {
	fromAddress, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	toAddress, err := sdk.AccAddressFromBech32(m.To)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address (%s)", err)
	}

	if fromAddress.Equals(toAddress) {
		return sdkerrors.Wrap(ErrSameAccount, m.From)
	}

	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrInvalidSignature, "signature is empty")
	}
	sig, err := hex.DecodeString(m.Signature)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSignature, "could not hex decode signature: %s", m.Signature)
	}
	pubKey, err := crypto.SigToPub(MigrateAccountSignatureHash(fromAddress, toAddress), sig)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSignature, "sig to pub key error: %s", err)
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(address.Bytes(), toAddress.Bytes()) {
		return sdkerrors.Wrapf(ErrInvalidSignature, "signature key not equal to address, expected %s, got %s", m.To, sdk.AccAddress(address.Bytes()).String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgMigrateAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgMigrateAccount) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func MigrateAccountSignatureHash(from, to sdk.AccAddress) []byte {
	return crypto.Keccak256([]byte(MigrateAccountSignaturePrefix), from.Bytes(), to.Bytes())
}
