package types

import (
	"encoding/hex"
	"fmt"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ crosschaintypes.MsgValidateBasic = &TronMsgValidate{}

// TronMsgValidate
type TronMsgValidate struct {
	crosschaintypes.EthereumMsgValidate
}

func (b TronMsgValidate) MsgCreateOracleBridgerValidate(m crosschaintypes.MsgCreateOracleBridger) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrOracleAddress, m.OracleAddress)
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateTronAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, m.ExternalAddress)
	}
	if !m.DelegateAmount.IsValid() || !m.DelegateAmount.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalidCoin, m.DelegateAmount.String())
	}
	return nil
}

func (b TronMsgValidate) MsgOracleSetConfirmValidate(m crosschaintypes.MsgOracleSetConfirm) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateTronAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, m.ExternalAddress)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}

func (b TronMsgValidate) MsgOracleSetUpdatedClaimValidate(m crosschaintypes.MsgOracleSetUpdatedClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "members len == 0")
	}
	for _, member := range m.Members {
		if err = ValidateTronAddress(member.ExternalAddress); err != nil {
			return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, member.ExternalAddress)
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "member power == 0")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "block height == 0")
	}
	return nil
}

func (b TronMsgValidate) MsgBridgeTokenClaimValidate(m crosschaintypes.MsgBridgeTokenClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrTokenContractAddress, m.TokenContract)
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not decode hex channelIbc string: %s", m.ChannelIbc)
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "token name is empty")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "token symbol is empty")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "block height == 0")
	}
	return nil
}

func (b TronMsgValidate) MsgSendToExternalClaimValidate(m crosschaintypes.MsgSendToExternalClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrTokenContractAddress, m.TokenContract)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "block height == 0")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "batch nonce == 0")
	}
	return nil
}

func (b TronMsgValidate) MsgSendToFxClaimValidate(m crosschaintypes.MsgSendToFxClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateTronAddress(m.Sender); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, m.Sender)
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrTokenContractAddress, m.TokenContract)
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Receiver)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "amount cannot be negative")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not decode hex targetIbc string: %s", m.TargetIbc)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "block height == 0")
	}
	return nil
}

func (b TronMsgValidate) MsgSendToExternalValidate(m crosschaintypes.MsgSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if err = ValidateTronAddress(m.Dest); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, m.Dest)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalidCoin, m.Amount.String())
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalidCoin, m.BridgeFee.String())
	}
	return nil
}

func (b TronMsgValidate) MsgRequestBatchValidate(m crosschaintypes.MsgRequestBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("denom is empty:%s", m.Denom))
	}
	if !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, fmt.Sprintf("minimum fee is positive:%s", m.MinimumFee.String()))
	}
	if err = ValidateTronAddress(m.FeeReceive); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, m.FeeReceive)
	}
	return nil
}

func (b TronMsgValidate) MsgConfirmBatchValidate(m crosschaintypes.MsgConfirmBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateTronAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrExternalAddress, m.ExternalAddress)
	}
	if err = ValidateTronAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintypes.ErrTokenContractAddress, m.TokenContract)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(crosschaintypes.ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}
