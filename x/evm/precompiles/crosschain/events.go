package crosschain

var (
	CancelSendToExternalEvent = GetABI().Events[CancelSendToExternalEventName]
	CrossChainEvent           = GetABI().Events[CrossChainEventName]
	IncreaseBridgeFeeEvent    = GetABI().Events[IncreaseBridgeFeeEventName]
	BridgeCallEvent           = GetABI().Events[BridgeCallEventName]
)
