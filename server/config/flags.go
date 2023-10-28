package config

const (
	// BypassMinFeeMsgTypesKey defines the configuration key for the
	// BypassMinFeeMsgTypes value.
	BypassMinFeeMsgTypesKey = "bypass-min-fee.msg-types" //nolint:gosec #nosec G101

	BypassMinFeeMsgMaxGasUsageKey = "bypass-min-fee.msg-max-gas-usage" //nolint:gosec #nosec G101
)

const (
	DefaultGasCap uint64 = 30000000
)
