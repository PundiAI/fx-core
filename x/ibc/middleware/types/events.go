package types

// IBC transfer events
const (
	EventTypeReceive      = "ibc_receive"
	AttributeKeyRecvError = "error"

	EventTypeIBCCall            = "ibc_call"
	AttributeKeyIBCCallType     = "ibc_call_type"
	AttributeKeyIBCCallErrCause = "ibc_call_err_cause"
	AttributeKeyIBCCallSuccess  = "ibc_call_success"
)
