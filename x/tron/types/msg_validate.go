package types

import (
	"encoding/hex"
	"fmt"

	crosschaintype "github.com/functionx/fx-core/x/crosschain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ crosschaintype.MsgValidateBasic = &MsgValidateBasic{}

// MsgValidateBasic
type MsgValidateBasic struct {
	crosschaintype.EthereumMsgValidateBasic
}

func (b MsgValidateBasic) MsgSetOrchestratorAddressValidate(m crosschaintype.MsgSetOrchestratorAddress) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Oracle); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOracleAddress, m.Oracle)
	}
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.Orchestrator)
	}
	if err = ValidateExternalAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, m.ExternalAddress)
	}
	if !m.Deposit.IsValid() || !m.Deposit.IsPositive() {
		return sdkerrors.Wrap(crosschaintype.ErrInvalidCoin, m.Deposit.String())
	}
	return nil
}

func (b MsgValidateBasic) MsgOracleSetConfirmValidate(m crosschaintype.MsgOracleSetConfirm) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.OrchestratorAddress)
	}
	if err = ValidateExternalAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, m.ExternalAddress)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(crosschaintype.ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}

func (b MsgValidateBasic) MsgOracleSetUpdatedClaimValidate(m crosschaintype.MsgOracleSetUpdatedClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.Orchestrator)
	}
	if len(m.Members) == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "members len == 0")
	}
	for _, member := range m.Members {
		if err = ValidateExternalAddress(member.ExternalAddress); err != nil {
			return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, member.ExternalAddress)
		}
		if member.Power == 0 {
			return sdkerrors.Wrap(crosschaintype.ErrInvalid, "member power == 0")
		}
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "block height == 0")
	}
	return nil
}

func (b MsgValidateBasic) MsgBridgeTokenClaimValidate(m crosschaintype.MsgBridgeTokenClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.Orchestrator)
	}
	if err = ValidateExternalAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrTokenContractAddress, m.TokenContract)
	}
	if _, err = hex.DecodeString(m.ChannelIbc); len(m.ChannelIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(crosschaintype.ErrInvalid, "could not decode hex channelIbc string: %s", m.ChannelIbc)
	}
	if len(m.Name) == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "token name is empty")
	}
	if len(m.Symbol) == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "token symbol is empty")
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "block height == 0")
	}
	return nil
}

func (b MsgValidateBasic) MsgSendToExternalClaimValidate(m crosschaintype.MsgSendToExternalClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.Orchestrator)
	}
	if err = ValidateExternalAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrTokenContractAddress, m.TokenContract)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "block height == 0")
	}
	if m.BatchNonce == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "batch nonce == 0")
	}
	return nil
}

func (b MsgValidateBasic) MsgSendToFxClaimValidate(m crosschaintype.MsgSendToFxClaim) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.Orchestrator)
	}
	if err = ValidateExternalAddress(m.Sender); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, m.Sender)
	}
	if err = ValidateExternalAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrTokenContractAddress, m.TokenContract)
	}
	if _, err = sdk.AccAddressFromBech32(m.Receiver); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Receiver)
	}
	if m.Amount.IsNil() || m.Amount.IsNegative() {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "amount cannot be negative")
	}
	if _, err = hex.DecodeString(m.TargetIbc); len(m.TargetIbc) > 0 && err != nil {
		return sdkerrors.Wrapf(crosschaintype.ErrInvalid, "could not decode hex targetIbc string: %s", m.TargetIbc)
	}
	if m.EventNonce == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "event nonce == 0")
	}
	if m.BlockHeight == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "block height == 0")
	}
	return nil
}

func (b MsgValidateBasic) MsgSendToExternalValidate(m crosschaintype.MsgSendToExternal) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if err = ValidateExternalAddress(m.Dest); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, m.Dest)
	}
	if !m.Amount.IsValid() || !m.Amount.IsPositive() {
		return sdkerrors.Wrap(crosschaintype.ErrInvalidCoin, m.Amount.String())
	}
	if m.Amount.Denom != m.BridgeFee.Denom {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, fmt.Sprintf("fee and amount must be the same type %s != %s", m.Amount.Denom, m.BridgeFee.Denom))
	}
	if !m.BridgeFee.IsValid() || !m.BridgeFee.IsPositive() {
		return sdkerrors.Wrap(crosschaintype.ErrInvalidCoin, m.BridgeFee.String())
	}
	return nil
}

func (b MsgValidateBasic) MsgRequestBatchValidate(m crosschaintype.MsgRequestBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Sender)
	}
	if len(m.Denom) <= 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, fmt.Sprintf("denom is empty:%s", m.Denom))
	}
	if !m.MinimumFee.IsPositive() {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, fmt.Sprintf("minimum fee is positive:%s", m.MinimumFee.String()))
	}
	if err = ValidateExternalAddress(m.FeeReceive); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, m.FeeReceive)
	}
	return nil
}

func (b MsgValidateBasic) MsgConfirmBatchValidate(m crosschaintype.MsgConfirmBatch) (err error) {
	if _, err = sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrOrchestratorAddress, m.OrchestratorAddress)
	}
	if err = ValidateExternalAddress(m.ExternalAddress); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrExternalAddress, m.ExternalAddress)
	}
	if err = ValidateExternalAddress(m.TokenContract); err != nil {
		return sdkerrors.Wrap(crosschaintype.ErrTokenContractAddress, m.TokenContract)
	}
	if len(m.Signature) == 0 {
		return sdkerrors.Wrap(crosschaintype.ErrInvalid, "signature is empty")
	}
	if _, err = hex.DecodeString(m.Signature); err != nil {
		return sdkerrors.Wrapf(crosschaintype.ErrInvalid, "could not hex decode signature: %s", m.Signature)
	}
	return nil
}
