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
	EventTypeSendToExternalCanceled = "send_to_external_canceled"
	EventTypeIncreaseBridgeFee      = "increase_bridge_fee"
	AttributeKeyOutgoingTxID        = "outgoing_tx_id"
	AttributeKeyPendingOutgoingTxID = "pending_outgoing_tx_id"
	AttributeKeyIncreaseFee         = "increase_fee"
	EventTypeAddPendingRewards      = "add_pending_rewards"
	AttributeKeyPendingID           = "pending_id"
	AttributeKeyPendingRewards      = "pending_rewards"
	AttributeKeyPendingType         = "pending_type"

	EventTypeOutgoingBatch           = "outgoing_batch"
	EventTypeOutgoingBatchCanceled   = "outgoing_batch_canceled"
	AttributeKeyOutgoingTxIds        = "outgoing_tx_ids"
	AttributeKeyOutgoingBatchNonce   = "batch_nonce"
	AttributeKeyOutgoingBatchTimeout = "outgoing_batch_timeout"

	EventTypeIbcTransfer         = "ibc_transfer"
	AttributeKeyIbcSendSequence  = "ibc_send_sequence"
	AttributeKeyIbcSourcePort    = "ibc_source_port"
	AttributeKeyIbcSourceChannel = "ibc_source_channel"

	EventTypeEvmTransfer = "evm_transfer"

	EventTypeBridgeCallEvent = "bridge_call_event"

	EventTypeBridgeCall         = "bridge_call"
	EventTypePendingBridgeCall  = "pending_bridge_call"
	AttributeKeyBridgeCallNonce = "bridge_call_nonce"

	EventTypeProvideLiquidity         = "provide_liquidity"
	AttributeKeyProvideLiquidityTxIds = "provide_liquidity_tx_ids"

	EventTypeBridgeCallRefundOut = "bridge_call_refund_out"

	EventTypeBridgeCallRefund = "bridge_call_refund"
	AttributeKeyRefund        = "refund"

	EventTypeBridgeCallResult = "bridge_call_result"
)

// precompile contract events
var (
	CancelSendToExternalEvent = GetABI().Events[CancelSendToExternalEventName]
	CrossChainEvent           = GetABI().Events[CrossChainEventName]
	IncreaseBridgeFeeEvent    = GetABI().Events[IncreaseBridgeFeeEventName]
	BridgeCallEvent           = GetABI().Events[BridgeCallEventName]
)
