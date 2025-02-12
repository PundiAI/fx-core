package types

const (
	EventTypeContractEvent = "observation"
	AttributeKeyClaimType  = "claim_type"
	AttributeKeyClaimHash  = "claim_hash"

	EventTypeOracleSetUpdate   = "oracle_set_update"
	AttributeKeyOracleSetNonce = "oracle_set_nonce"
	AttributeKeyOracleSetLen   = "oracle_set_len"

	EventTypeSendToExternal          = "send_to_external"
	EventTypeOutgoingBatch           = "outgoing_batch"
	EventTypeOutgoingBatchCanceled   = "outgoing_batch_canceled"
	AttributeKeyOutgoingTxID         = "outgoing_tx_id"
	AttributeKeyPendingOutgoingTxID  = "pending_outgoing_tx_id"
	AttributeKeyOutgoingTxIds        = "outgoing_tx_ids"
	AttributeKeyOutgoingBatchNonce   = "batch_nonce"
	AttributeKeyOutgoingBatchTimeout = "outgoing_batch_timeout"

	EventTypeEvmTransfer      = "evm_transfer"
	EventTypeBridgeCallEvent  = "bridge_call_event"
	EventTypeBridgeCall       = "bridge_call"
	EventTypeBridgeCallFailed = "bridge_call_failed"
	EventTypeBridgeCallResend = "bridge_call_resend"
	EventTypeBridgeCallRefund = "bridge_call_refund"
	EventTypeBridgeCallResult = "bridge_call_result"

	AttributeKeyBridgeCallResendOldNonce = "old_bridge_call_nonce"
	AttributeKeyBridgeCallResendNewNonce = "new_bridge_call_nonce"
	AttributeKeyIBCSequence              = "ibc_sequence"
	AttributeKeyBridgeCallNonce          = "bridge_call_nonce"
	AttributeKeyBridgeCallResultNonce    = "bridge_call_result_nonce"
	AttributeKeyRefundAddr               = "refund_addr"
	AttributeKeyBlockHeight              = "block_height"
	AttributeKeyStateSuccess             = "state_success"
	AttributeKeyErrCause                 = "err_cause"
	AttributeKeyEventNonce               = "event_nonce"
)
