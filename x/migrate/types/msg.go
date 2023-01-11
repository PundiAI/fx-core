package types

import (
	"bytes"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

const TypeMsgMigrateAccount = "migrate_account"

var _ sdk.Msg = &MsgMigrateAccount{}

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
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	// check to address
	if err := fxtypes.ValidateEthereumAddress(m.To); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}
	toAddress := common.HexToAddress(m.To)

	// check same account
	if bytes.Equal(fromAddress.Bytes(), toAddress.Bytes()) {
		return sdkerrors.ErrInvalidRequest.Wrap("same account")
	}

	// check signature
	if len(m.Signature) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty signature")
	}
	sig, err := hex.DecodeString(m.Signature)
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("could not hex decode signature")
	}
	pubKey, err := crypto.SigToPub(MigrateAccountSignatureHash(fromAddress, toAddress.Bytes()), sig)
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("sig to pub key error")
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(address.Bytes(), toAddress.Bytes()) {
		return sdkerrors.ErrInvalidRequest.Wrap("signature key not equal to address")
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
