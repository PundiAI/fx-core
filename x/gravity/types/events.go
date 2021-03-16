package types

const (
	EventTypeObservation              = "observation"
	EventTypeOutgoingBatch            = "outgoing_batch"
	EventTypeMultisigUpdateRequest    = "multisig_update_request"
	EventTypeOutgoingBatchCanceled    = "outgoing_batch_canceled"
	EventTypeBridgeWithdrawalReceived = "withdrawal_received"
	EventTypeBridgeDepositReceived    = "deposit_received"
	EventTypeBridgeWithdrawCanceled   = "withdraw_canceled"

	EventTypeAttestationHandlerDeposit = "attestation_handler_deposit_claim"

	AttributeKeyAttestationID    = "attestation_id"
	AttributeKeyBatchConfirmKey  = "batch_confirm_key"
	AttributeKeyValsetConfirmKey = "valset_confirm_key"
	AttributeKeyMultisigID       = "multisig_id"
	AttributeKeyOutgoingBatchID  = "batch_id"
	AttributeKeyOutgoingTXID     = "outgoing_tx_id"
	AttributeKeyAttestationType  = "attestation_type"
	AttributeKeyContract         = "bridge_contract"
	AttributeKeyBatchNonceTxIds  = "batch_nonce_tx_ids"

	AttributeKeyWithdrawalTokenContract = "withdrawal_token_contract"
	AttributeKeyWithdrawalSender        = "withdrawal_sender"
	AttributeKeyWithdrawalReceiver      = "withdrawal_receiver"
	AttributeKeyWithdrawalAmount        = "withdrawal_amount"
	AttributeKeyWithdrawalFee           = "withdrawal_fee"
	AttributeKeyNonce                   = "nonce"
	AttributeKeyValsetNonce             = "valset_nonce"
	AttributeKeyBatchNonce              = "batch_nonce"
	AttributeKeyBridgeChainID           = "bridge_chain_id"
	AttributeKeySetOperatorAddr         = "set_operator_address"

	AttributeKeyAttestationHandlerNonce                   = "nonce"
	AttributeKeyAttestationHandlerTokenContract           = "token_contract"
	AttributeKeyAttestationHandlerAmount                  = "amount"
	AttributeKeyAttestationHandlerEthereumSender          = "ethereum_sender"
	AttributeKeyAttestationHandlerFxReceiver              = "fx_receiver"
	AttributeKeyAttestationHandlerTargetIbc               = "target_ibc"
	AttributeKeyAttestationHandlerIbcChannelSendSequence  = "ibc_channel_send_sequence"
	AttributeKeyAttestationHandlerIbcChannelSourcePort    = "ibc_channel_source_port"
	AttributeKeyAttestationHandlerIbcChannelSourceChannel = "ibc_channel_source_channel"

	AttributeValueCategory = ModuleName
)
