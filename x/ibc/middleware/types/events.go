package types

// IBC transfer events
const (
	EventTypeReceive      = "ibc_receive"
	EventTypeReceiveRoute = "ibc_receive_route"
	EventTypeIBCCall      = "ibc_call"

	AttributeKeyRoute    = "route"
	AttributeKeyError    = "error"
	AttributeKeySuccess  = "success"
	AttributeKeyType     = "type"
	AttributeKeyErrCause = "err_cause"
)
