package types

const (
	EventTypeContractEvent = "observation"
	AttributeKeyClaimType  = "claim_type"
	AttributeKeyEventNonce = "event_nonce"
	AttributeKeyClaimHash  = "claim_hash"

	AttributeKeyBlockHeight  = "block_height"
	AttributeKeyStateSuccess = "state_success"
	AttributeKeyErrCause     = "err_cause"

	EventTypeOracleSetUpdate   = "oracle_set_update"
	AttributeKeyOracleSetNonce = "oracle_set_nonce"
	AttributeKeyOracleSetLen   = "oracle_set_len"

	EventTypeSendToExternal         = "send_to_external"
	AttributeKeyOutgoingTxID        = "outgoing_tx_id"
	AttributeKeyPendingOutgoingTxID = "pending_outgoing_tx_id"

	EventTypeOutgoingBatch           = "outgoing_batch"
	EventTypeOutgoingBatchCanceled   = "outgoing_batch_canceled"
	AttributeKeyOutgoingTxIds        = "outgoing_tx_ids"
	AttributeKeyOutgoingBatchNonce   = "batch_nonce"
	AttributeKeyOutgoingBatchTimeout = "outgoing_batch_timeout"

	EventTypeEvmTransfer = "evm_transfer"

	EventTypeBridgeCallEvent = "bridge_call_event"

	EventTypeBridgeCall         = "bridge_call"
	AttributeKeyBridgeCallNonce = "bridge_call_nonce"

	EventTypeBridgeCallFailed              = "bridge_call_failed"
	AttributeKeyBridgeCallFailedRefundAddr = "refund_addr"

	EventTypeBridgeCallResend            = "bridge_call_resend"
	AttributeKeyBridgeCallResendOldNonce = "old_bridge_call_nonce"
	AttributeKeyBridgeCallResendNewNonce = "new_bridge_call_nonce"

	EventTypeBridgeCallRefund = "bridge_call_refund"
	AttributeKeyRefund        = "refund"

	EventTypeBridgeCallResult = "bridge_call_result"
)
