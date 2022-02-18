package types

import (
	"bytes"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

var (
	_ sdk.Msg = &MsgConvertCoin{}
	_ sdk.Msg = &MsgConvertFIP20{}
)

const (
	TypeMsgConvertCoin  = "convert_coin"
	TypeMsgConvertFIP20 = "convert_FIP20"
)

// NewMsgConvertCoin creates a new instance of MsgConvertCoin
func NewMsgConvertCoin(coin sdk.Coin, receiver common.Address, sender sdk.AccAddress) *MsgConvertCoin { // nolint: interfacer
	return &MsgConvertCoin{
		Coin:     coin,
		Receiver: receiver.Hex(),
		Sender:   sender.String(),
	}
}

// Route should return the name of the module
func (msg MsgConvertCoin) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConvertCoin) Type() string { return TypeMsgConvertCoin }

// ValidateBasic runs stateless checks on the message
func (msg MsgConvertCoin) ValidateBasic() error {
	if err := ValidateIntrarelayerDenom(msg.Coin.Denom); err != nil {
		if err := ibctransfertypes.ValidateIBCDenom(msg.Coin.Denom); err != nil {
			return err
		}
	}

	if !msg.Coin.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "cannot mint a non-positive amount")
	}
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	if !common.IsHexAddress(msg.Receiver) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver hex address %s", msg.Receiver)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgConvertCoin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConvertCoin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}

	return []sdk.AccAddress{addr}
}

// NewMsgConvertFIP20 creates a new instance of MsgConvertFIP20
func NewMsgConvertFIP20(amount sdk.Int, sender, receiver sdk.AccAddress, contract common.Address, bech32PubKey string) *MsgConvertFIP20 { // nolint: interfacer
	return &MsgConvertFIP20{
		ContractAddress: contract.String(),
		Amount:          amount,
		Receiver:        receiver.String(),
		Sender:          sender.String(),
		PubKey:          bech32PubKey,
	}
}

// Route should return the name of the module
func (msg MsgConvertFIP20) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConvertFIP20) Type() string { return TypeMsgConvertFIP20 }

// ValidateBasic runs stateless checks on the message
func (msg MsgConvertFIP20) ValidateBasic() error {
	if !common.IsHexAddress(msg.ContractAddress) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract hex address '%s'", msg.ContractAddress)
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "cannot mint a non-positive amount")
	}
	_, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid receiver address")
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid public key")
	}
	decompressPubkey, err := crypto.DecompressPubkey(pubKey.Bytes())
	if err != nil {
		return sdkerrors.Wrap(err, "invalid public key, can not decompress")
	}
	ethSender := crypto.PubkeyToAddress(*decompressPubkey)

	if !bytes.Equal(sender, pubKey.Address().Bytes()) &&
		!bytes.Equal(sender, ethSender.Bytes()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "public key does not match sender address")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgConvertFIP20) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConvertFIP20) GetSigners() []sdk.AccAddress {
	acc, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{acc}
}

func (msg MsgConvertFIP20) HexAddress() common.Address {
	pubKey, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	decompressPubkey, _ := crypto.DecompressPubkey(pubKey.Bytes())
	return crypto.PubkeyToAddress(*decompressPubkey)
}

func PubKeyToEIP55Address(pubKey cryptotypes.PubKey) (common.Address, error) {
	uncompressedPubKey, err := crypto.DecompressPubkey(pubKey.Bytes())
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*uncompressedPubKey), nil
}
