package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	crosschaintypes "github.com/functionx/fx-core/v2/x/crosschain/types"
)

var _ crosschaintypes.MsgValidateBasic = &TronMsgValidate{}

type TronMsgValidate struct {
	crosschaintypes.EthereumMsgValidate
}

func (b TronMsgValidate) MsgBondedOracleValidate(m crosschaintypes.MsgBondedOracle) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "oracle address")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if err = ValidateTronAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "external address")
	}
	if !m.DelegateAmount.IsValid() || !m.DelegateAmount.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "delegate amount")
	}
	return nil
}

func (b TronMsgValidate) MsgOracleSetConfirmValidate(m crosschaintypes.MsgOracleSetConfirm) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if err = ValidateTronAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "external address")
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrEmpty, "signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}

func (b TronMsgValidate) MsgOracleSetUpdatedClaimValidate(m crosschaintypes.MsgOracleSetUpdatedClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrEmpty, "members")
	}
	for _, member := range m.Members {
		if err = ValidateTronAddress(member.ExternalAddress); err != nil {
			return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "external address")
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(crosschaintypes.ErrEmpty, "member power")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "block height")
	}
	return nil
}

func (b TronMsgValidate) MsgBridgeTokenClaimValidate(m crosschaintypes.MsgBridgeTokenClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "token contract")
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not decode hex channelIbc string: %s", m.ChannelIbc)
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrEmpty, "token name")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrEmpty, "token symbol")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "block height")
	}
	return nil
}

func (b TronMsgValidate) MsgSendToExternalClaimValidate(m crosschaintypes.MsgSendToExternalClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "token contract")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "block height")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "batch nonce")
	}
	return nil
}

func (b TronMsgValidate) MsgSendToFxClaimValidate(m crosschaintypes.MsgSendToFxClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if err = ValidateTronAddress(m.Sender); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "sender address")
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "token contract")
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "receiver address")
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "amount cannot be negative")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not decode hex targetIbc string: %s", m.TargetIbc)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "block height")
	}
	return nil
}

func (b TronMsgValidate) MsgSendToExternalValidate(m crosschaintypes.MsgSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "sender address")
	}
	if err = ValidateTronAddress(m.Dest); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "dest")
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "amount")
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridge fee")
	}
	return nil
}

func (b TronMsgValidate) MsgRequestBatchValidate(m crosschaintypes.MsgRequestBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "sender address")
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrUnknown, "denom")
	}
	if m.MinimumFee.IsNil() || !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "minimum fee")
	}
	if err = ValidateTronAddress(m.FeeReceive); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "fee receive address")
	}
	if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "base fee")
	}
	return nil
}

func (b TronMsgValidate) MsgConfirmBatchValidate(m crosschaintypes.MsgConfirmBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	if err = ValidateTronAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "external address")
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "token contract")
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrEmpty, "signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}
