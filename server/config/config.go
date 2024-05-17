package config

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ethermintconfig "github.com/evmos/ethermint/server/config"
	"github.com/spf13/viper"
)

// BypassMinFee defines custom that will bypass minimum fee checks during CheckTx.
type BypassMinFee struct {
	// MsgTypes defines custom message types the operator may set that
	// will bypass minimum fee checks during CheckTx.
	MsgTypes       []string `mapstructure:"msg-types"`
	MsgMaxGasUsage uint64   `mapstructure:"msg-max-gas-usage"`
}

func (f BypassMinFee) Validate() error {
	for _, msgType := range f.MsgTypes {
		if strings.TrimSpace(msgType) != msgType || !strings.HasPrefix(msgType, "/") {
			return fmt.Errorf("invalid message type: %s", msgType)
		}
	}
	return nil
}

// DefaultBypassMinFee returns the default BypassMinFee configuration
func DefaultBypassMinFee() BypassMinFee {
	return BypassMinFee{
		MsgTypes:       []string{},
		MsgMaxGasUsage: uint64(300_000),
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

func GetConfig(v *viper.Viper) (*Config, error) {
	cfg, err := ethermintconfig.GetConfig(v)
	if err != nil {
		return nil, err
	}
	return &Config{
		Config:       cfg.Config,
		BypassMinFee: BypassMinFee{},
		EVM:          cfg.EVM,
		JSONRPC:      cfg.JSONRPC,
		TLS:          cfg.TLS,
	}, nil
}

// ValidateBasic returns an error any of the application configuration fields are invalid
func (c *Config) ValidateBasic() error {
	if err := c.BypassMinFee.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid bypass-min-fee config value: %s", err.Error())
	}

	if err := c.EVM.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid evm config value: %s", err.Error())
	}

	if err := c.JSONRPC.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid json-rpc config value: %s", err.Error())
	}

	if err := c.TLS.Validate(); err != nil {
		return errorsmod.Wrapf(errortypes.ErrAppConfig, "invalid tls config value: %s", err.Error())
	}

	return c.Config.ValidateBasic()
}

func (c *Config) ToEthermintConfig() *ethermintconfig.Config {
	return &ethermintconfig.Config{
		Config:  c.Config,
		EVM:     c.EVM,
		JSONRPC: c.JSONRPC,
		TLS:     c.TLS,
	}
}

// AppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func AppConfig(mintGasPrice sdk.Coin) (string, interface{}) {
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
	srvCfg.MinGasPrices = mintGasPrice.String()
	srvCfg.Rosetta.DenomToSuggest = mintGasPrice.Denom

	customAppConfig := Config{
		Config:       *srvCfg,
		BypassMinFee: DefaultBypassMinFee(),
		EVM:          *ethermintconfig.DefaultEVMConfig(),
		JSONRPC:      *ethermintconfig.DefaultJSONRPCConfig(),
		TLS:          *ethermintconfig.DefaultTLSConfig(),
	}

	customAppConfig.JSONRPC.GasCap = DefaultGasCap

	customAppTemplate := DefaultConfigTemplate()

	return customAppTemplate, customAppConfig
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	cfg := &Config{
		Config:       *config.DefaultConfig(),
		BypassMinFee: DefaultBypassMinFee(),
		EVM:          *ethermintconfig.DefaultEVMConfig(),
		JSONRPC:      *ethermintconfig.DefaultJSONRPCConfig(),
		TLS:          *ethermintconfig.DefaultTLSConfig(),
	}
	cfg.JSONRPC.GasCap = DefaultGasCap
	return cfg
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

# MsgMaxGasUsage defines gas consumption threshold .Default: 300000
msg-max-gas-usage = {{ .BypassMinFee.MsgMaxGasUsage }}
`
