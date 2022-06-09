package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ MsgValidateBasic = &EthereumMsgValidate{}

// EthereumMsgValidate
type EthereumMsgValidate struct{}

func (b EthereumMsgValidate) MsgCreateOracleBridgerValidate(m MsgCreateOracleBridger) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrOracleAddress, m.OracleAddress)
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrExternalAddress, m.ExternalAddress)
	}
	if !m.DelegateAmount.IsValid() || !m.DelegateAmount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidCoin, m.DelegateAmount.String())
	}
	return nil
}

func (b EthereumMsgValidate) MsgAddOracleDelegateValidate(m MsgAddOracleDelegate) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrOracleAddress, m.OracleAddress)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidCoin, m.Amount.String())
	}
	return nil
}

func (b EthereumMsgValidate) MsgEditOracleValidate(m MsgEditOracle) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b EthereumMsgValidate) MsgWithdrawRewardValidate(m MsgWithdrawReward) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b EthereumMsgValidate) MsgOracleSetConfirmValidate(m MsgOracleSetConfirm) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrExternalAddress, m.ExternalAddress)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}

func (b EthereumMsgValidate) MsgOracleSetUpdatedClaimValidate(m MsgOracleSetUpdatedClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "members len == 0")
	}
	for _, member := range m.Members {
		if err = ValidateEthereumAddress(member.ExternalAddress); err != nil {
			return sdkerrors.Wrap(ErrExternalAddress, member.ExternalAddress)
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(ErrInvalid, "member power == 0")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	return nil
}

func (b EthereumMsgValidate) MsgBridgeTokenClaimValidate(m MsgBridgeTokenClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrTokenContractAddress, m.TokenContract)
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not decode hex channelIbc string: %s", m.ChannelIbc)
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "token name is empty")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "token symbol is empty")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	return nil
}

func (b EthereumMsgValidate) MsgSendToExternalClaimValidate(m MsgSendToExternalClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrTokenContractAddress, m.TokenContract)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "batch nonce == 0")
	}
	return nil
}

func (b EthereumMsgValidate) MsgSendToFxClaimValidate(m MsgSendToFxClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateEthereumAddress(m.Sender); err != nil {
		return sdkerrors.Wrap(ErrExternalAddress, m.Sender)
	}
	if err = ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrTokenContractAddress, m.TokenContract)
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Receiver)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalid, "amount cannot be negative")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not decode hex targetIbc string: %s", m.TargetIbc)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "block height == 0")
	}
	return nil
}

func (b EthereumMsgValidate) MsgSendToExternalValidate(m MsgSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if err = ValidateEthereumAddress(m.Dest); err != nil {
		return sdkerrors.Wrap(ErrExternalAddress, m.Dest)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidCoin, m.Amount.String())
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(ErrInvalid, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidCoin, m.BridgeFee.String())
	}
	return nil
}

func (b EthereumMsgValidate) MsgCancelSendToExternalValidate(m MsgCancelSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if m.TransactionId == 0 {
		return sdkerrors.Wrap(ErrInvalid, "transaction id == 0")
	}
	return nil
}

func (b EthereumMsgValidate) MsgRequestBatchValidate(m MsgRequestBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(ErrInvalid, fmt.Sprintf("denom is empty:%s", m.Denom))
	}
	if !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(ErrInvalid, fmt.Sprintf("minimum fee is positive:%s", m.MinimumFee.String()))
	}
	if err = ValidateEthereumAddress(m.FeeReceive); err != nil {
		return sdkerrors.Wrap(ErrExternalAddress, m.FeeReceive)
	}
	if m.BaseFee == nil || m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return ErrInvalidRequestBatchBaseFee
	}
	return nil
}

func (b EthereumMsgValidate) MsgConfirmBatchValidate(m MsgConfirmBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrBridgerAddress, m.BridgerAddress)
	}
	if err = ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrExternalAddress, m.ExternalAddress)
	}
	if err = ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrTokenContractAddress, m.TokenContract)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}
