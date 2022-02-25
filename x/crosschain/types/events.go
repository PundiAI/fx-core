package types

const (
	EventTypeObservation            = "observation"
	EventTypeOutgoingBatch          = "outgoing_batch"
	EventTypeOracleSetRequest       = "multisig_update_request"
	EventTypeOutgoingBatchCanceled  = "outgoing_batch_canceled"
	EventTypeSendToExternalReceived = "send_to_external_received"
	EventTypeSendToExternalCanceled = "send_to_external_canceled"
	EventTypeSendToFx               = "send_to_fx"

	AttributeKeyAttestationID       = "attestation_id"
	AttributeKeyBatchConfirmKey     = "batch_confirm_key"
	AttributeKeyOracleSetConfirmKey = "valset_confirm_key"
	AttributeKeySetOperatorAddr     = "set_operator_address"

	AttributeKeyAttestationType = "attestation_type"
	AttributeKeyEventNonce      = "event_nonce"
	AttributeKeyOracleSetNonce  = "oracle_set_nonce"
	AttributeKeyOutgoingTXID    = "outgoing_tx_id"
	AttributeKeyBatchNonceTxIds = "batch_nonce_tx_ids"
	AttributeKeyBatchNonce      = "batch_nonce"

	AttributeKeyAttestationHandlerIbcChannelSendSequence  = "ibc_channel_send_sequence"
	AttributeKeyAttestationHandlerIbcChannelSourcePort    = "ibc_channel_source_port"
	AttributeKeyAttestationHandlerIbcChannelSourceChannel = "ibc_channel_source_channel"

	AttributeKeyAttestationHandlerEvmTransfer = "evm_transfer"
)
