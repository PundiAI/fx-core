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
