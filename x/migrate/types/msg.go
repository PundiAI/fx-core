package types

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

const TypeMsgMigrateAccount = "migrate_account"

var (
	_ sdk.Msg = &MsgMigrateAccount{}
)

// NewMsgMigrateAccount returns a new MsgMigrateAccount
func NewMsgMigrateAccount(from sdk.AccAddress, to common.Address, signature string) *MsgMigrateAccount {
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, err.Error())
	}
	//check to address
	if err := fxtypes.ValidateEthereumAddress(m.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, err.Error())
	}
	toAddress := common.HexToAddress(m.To)

	//check same account
	if bytes.Equal(fromAddress.Bytes(), toAddress.Bytes()) {
		return sdkerrors.Wrap(ErrSameAccount, m.From)
	}

	//check signature
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrInvalidSignature, "signature is empty")
	}
	sig, err := hex.DecodeString(m.Signature)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSignature, "could not hex decode signature: %s", m.Signature)
	}
	pubKey, err := crypto.SigToPub(MigrateAccountSignatureHash(fromAddress, toAddress.Bytes()), sig)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidSignature, "sig to pub key error: %s", err)
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(address.Bytes(), toAddress.Bytes()) {
		return sdkerrors.Wrapf(ErrInvalidSignature, "signature key not equal to address, expected %s, got %s", m.To, address.String())
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

func MigrateAccountSignatureHash(from, to []byte) []byte {
	return crypto.Keccak256([]byte(MigrateAccountSignaturePrefix), from, to)
}
