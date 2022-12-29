package types

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

var _ MsgValidateBasic = &MsgValidate{}

type MsgValidate struct{}

func (b MsgValidate) MsgBondedOracleValidate(m *MsgBondedOracle) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "oracle address")
	}
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "external address")
	}
	if !m.DelegateAmount.IsValid() || m.DelegateAmount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalid, "delegate amount")
	}
	if m.OracleAddress == m.BridgerAddress {
		return sdkerrors.Wrap(ErrInvalid, "same address")
	}
	return nil
}

func (b MsgValidate) MsgAddDelegateValidate(m *MsgAddDelegate) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "oracle address")
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalid, "amount")
	}
	return nil
}

func (b MsgValidate) MsgReDelegateValidate(m *MsgReDelegate) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "oracle address")
	}
	if _, err = sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "validator address")
	}
	return nil
}

func (b MsgValidate) MsgEditBridgerValidate(m *MsgEditBridger) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "oracle address")
	}
	if _, err = sdk.ValAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if m.OracleAddress == m.BridgerAddress {
		return sdkerrors.Wrap(ErrInvalid, "same address")
	}
	return nil
}

func (b MsgValidate) MsgWithdrawRewardValidate(m *MsgWithdrawReward) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "oracle address")
	}
	return nil
}

func (b MsgValidate) MsgUnbondedOracleValidate(m *MsgUnbondedOracle) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OracleAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "oracle address")
	}
	return nil
}

func (b MsgValidate) MsgOracleSetConfirmValidate(m *MsgOracleSetConfirm) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "external address")
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}

func (b MsgValidate) MsgOracleSetUpdatedClaimValidate(m *MsgOracleSetUpdatedClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "members")
	}
	for _, member := range m.Members {
		if err = fxtypes.ValidateEthereumAddress(member.ExternalAddress); err != nil {
			return sdkerrors.Wrap(ErrInvalid, "external address")
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(ErrEmpty, "member power")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrUnknown, "block height")
	}
	return nil
}

func (b MsgValidate) MsgBridgeTokenClaimValidate(m *MsgBridgeTokenClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "token contract")
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not decode hex channelIbc string: %s", m.ChannelIbc)
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "token name")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "token symbol")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrUnknown, "block height")
	}
	return nil
}

func (b MsgValidate) MsgSendToExternalClaimValidate(m *MsgSendToExternalClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "token contract")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrUnknown, "block height")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.Wrap(ErrUnknown, "batch nonce")
	}
	return nil
}

func (b MsgValidate) MsgSendToFxClaimValidate(m *MsgSendToFxClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.Sender); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "sender address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "token contract")
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "receiver address")
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalid, "amount cannot be negative")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not decode hex targetIbc string: %s", m.TargetIbc)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(ErrUnknown, "event nonce")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(ErrUnknown, "block height")
	}
	return nil
}

func (b MsgValidate) MsgSendToExternalValidate(m *MsgSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "sender address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.Dest); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "dest")
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalid, "amount")
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(ErrInvalid, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return sdkerrors.Wrap(ErrInvalid, "bridge fee")
	}
	return nil
}

func (b MsgValidate) MsgCancelSendToExternalValidate(m *MsgCancelSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "sender address")
	}
	if m.TransactionId == 0 {
		return sdkerrors.Wrap(ErrUnknown, "transaction id")
	}
	return nil
}

func (b MsgValidate) MsgRequestBatchValidate(m *MsgRequestBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "sender address")
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(ErrUnknown, "denom")
	}
	if m.MinimumFee.IsNil() || !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(ErrInvalid, "minimum fee")
	}
	if err = fxtypes.ValidateEthereumAddress(m.FeeReceive); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "fee receive address")
	}
	if m.BaseFee.IsNil() || m.BaseFee.IsNegative() {
		return sdkerrors.Wrap(ErrInvalid, "base fee")
	}
	return nil
}

func (b MsgValidate) MsgConfirmBatchValidate(m *MsgConfirmBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.BridgerAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "bridger address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "external address")
	}
	if err = fxtypes.ValidateEthereumAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(ErrInvalid, "token contract")
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "signature")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}
