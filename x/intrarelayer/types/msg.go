package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

var (
	_ sdk.Msg = &MsgConvertCoin{}
	_ sdk.Msg = &MsgConvertERC20{}
)

const (
	TypeMsgConvertCoin  = "convert_coin"
	TypeMsgConvertERC20 = "convert_ERC20"
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

// NewMsgConvertERC20 creates a new instance of MsgConvertERC20
func NewMsgConvertERC20(amount sdk.Int, receiver sdk.AccAddress, contract common.Address, bech32PubKey string) *MsgConvertERC20 { // nolint: interfacer

	return &MsgConvertERC20{
		ContractAddress: contract.String(),
		Amount:          amount,
		Receiver:        receiver.String(),
		PubKey:          bech32PubKey,
	}
}

// Route should return the name of the module
func (msg MsgConvertERC20) Route() string { return RouterKey }

// Type should return the action
func (msg MsgConvertERC20) Type() string { return TypeMsgConvertERC20 }

// ValidateBasic runs stateless checks on the message
func (msg MsgConvertERC20) ValidateBasic() error {
	if !common.IsHexAddress(msg.ContractAddress) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract hex address '%s'", msg.ContractAddress)
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "cannot mint a non-positive amount")
	}
	_, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid reciver address")
	}
	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid bench32 string '%s'", msg.PubKey)
	}
	if _, err := crypto.DecompressPubkey(pubKey.Bytes()); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "can not decompress pub key %s", msg.PubKey)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgConvertERC20) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgConvertERC20) GetSigners() []sdk.AccAddress {
	pubKey, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	return []sdk.AccAddress{pubKey.Address().Bytes()}
}

func (msg MsgConvertERC20) AccAddress() sdk.AccAddress {
	pubKey, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	return sdk.AccAddress(pubKey.Address().Bytes())
}
func (msg MsgConvertERC20) HexAddress() common.Address {
	pubKey, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, msg.PubKey)
	uncompressedPubKey, _ := crypto.DecompressPubkey(pubKey.Bytes())
	return crypto.PubkeyToAddress(*uncompressedPubKey)
}

func PubKeyToEIP55Address(pubKey cryptotypes.PubKey) (common.Address, error) {
	uncompressedPubKey, err := crypto.DecompressPubkey(pubKey.Bytes())
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*uncompressedPubKey), nil
}
