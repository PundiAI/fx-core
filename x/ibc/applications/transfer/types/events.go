package types

// IBC transfer events
const (
	EventTypeReceive      = "ibc_receive"
	EventTypeReceiveRoute = "ibc_receive_route"
	AttributeKeyRecvError = "error"

	AttributeKeyRouteSuccess = "success"
	AttributeKeyRoute        = "route"
	AttributeKeyRouteError   = "error"

	EventTypeIBCCall            = "ibc_call"
	AttributeKeyIBCCallType     = "ibc_call_type"
	AttributeKeyIBCCallErrCause = "ibc_call_err_cause"
	AttributeKeyIBCCallSuccess  = "ibc_call_success"
)
