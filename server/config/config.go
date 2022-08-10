package config

import (
	"github.com/cosmos/cosmos-sdk/server/config"
	ethermintconfig "github.com/evmos/ethermint/server/config"
)

// BypassMinFee defines custom that will bypass minimum fee checks during CheckTx.
type BypassMinFee struct {
	// MsgTypes defines custom message types the operator may set that
	// will bypass minimum fee checks during CheckTx.
	MsgTypes []string `mapstructure:"msg-types"`
}

// DefaultBypassMinFee returns the default BypassMinFee configuration
func DefaultBypassMinFee() BypassMinFee {
	return BypassMinFee{
		MsgTypes: []string{},
	}
}

type Config struct {
	config.Config `mapstructure:",squash"`

	// BypassMinFeeMsgTypes defines custom that will bypass minimum fee checks during CheckTx.
	BypassMinFee BypassMinFee `mapstructure:"bypass-min-fee"`

	EVM     ethermintconfig.EVMConfig     `mapstructure:"evm"`
	JSONRPC ethermintconfig.JSONRPCConfig `mapstructure:"json-rpc"`
	TLS     ethermintconfig.TLSConfig     `mapstructure:"tls"`
}

// AppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func AppConfig(denom string) (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := config.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In ethermint, we set the min gas prices to 0.
	if denom != "" {
		srvCfg.MinGasPrices = "0" + denom
	}

	customAppConfig := Config{
		Config:  *srvCfg,
		EVM:     *ethermintconfig.DefaultEVMConfig(),
		JSONRPC: *ethermintconfig.DefaultJSONRPCConfig(),
		TLS:     *ethermintconfig.DefaultTLSConfig(),
	}

	customAppTemplate := DefaultConfigTemplate()

	return customAppTemplate, customAppConfig
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config:       *config.DefaultConfig(),
		BypassMinFee: DefaultBypassMinFee(),
		EVM:          *ethermintconfig.DefaultEVMConfig(),
		JSONRPC:      *ethermintconfig.DefaultJSONRPCConfig(),
		TLS:          *ethermintconfig.DefaultTLSConfig(),
	}
}

func DefaultConfigTemplate() string {
	return config.DefaultConfigTemplate + CustomConfigTemplate + ethermintconfig.DefaultConfigTemplate
}

const CustomConfigTemplate = `
###############################################################################
###                        Custom Fx Configuration                        ###
###############################################################################
[bypass-min-fee]
# MsgTypes defines custom message types the operator may set that will bypass minimum fee checks during CheckTx.
# Example:
# ["/ibc.core.channel.v1.MsgRecvPacket", "/ibc.core.channel.v1.MsgAcknowledgement", ...]
msg-types = [{{ range .BypassMinFee.MsgTypes }}{{ printf "%q, " . }}{{end}}]
`
