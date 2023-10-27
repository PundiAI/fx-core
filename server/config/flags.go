package config

const (
	// BypassMinFeeMsgTypesKey defines the configuration key for the
	// BypassMinFeeMsgTypes value.
	BypassMinFeeMsgTypesKey = "bypass-min-fee.msg-types" //nolint:gosec

	BypassMinFeeMsgMaxGasUsageKey = "bypass-min-fee.msg-max-gas-usage" //nolint:gosec
)

const (
	DefaultGasCap uint64 = 30000000
)
