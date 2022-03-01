package v1

import (
	"bytes"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
	typescommon "github.com/functionx/fx-core/x/migrate/types/common"
)

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
func (m *MsgMigrateAccount) Route() string { return typescommon.RouterKey }

// Type should return the action
func (m *MsgMigrateAccount) Type() string { return "migrate_account" }

// ValidateBasic performs stateless checks
func (m *MsgMigrateAccount) ValidateBasic() (err error) {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return err
	}
	to, err := sdk.AccAddressFromBech32(m.To)
	if err != nil {
		return err
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(typescommon.ErrInvalidSignature, "signature is empty")
	}
	sig, err := hex.DecodeString(m.Signature)
	if err != nil {
		return sdkerrors.Wrapf(typescommon.ErrInvalidSignature, "could not hex decode signature: %s", m.Signature)
	}
	pubKey, err := crypto.SigToPub(MigrateAccountSignatureHash(from, to), sig)
	if err != nil {
		return sdkerrors.Wrapf(typescommon.ErrInvalidSignature, "unmarshal pub key error: %v", err)
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(address.Bytes(), to.Bytes()) { //fx address byte not equal recover address byte
		return sdkerrors.Wrap(typescommon.ErrInvalidAddress, "signature key not equal to address")
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
	return crypto.Keccak256([]byte(typescommon.MigrateAccountSignaturePrefix), from.Bytes(), to.Bytes())
}
