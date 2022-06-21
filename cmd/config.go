package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/functionx/fx-core/app/cli"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/viper"

	"github.com/mitchellh/mapstructure"

	"github.com/cosmos/cosmos-sdk/client"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	ethermintconfig "github.com/evmos/ethermint/server/config"
)

const (
	configFileName = "config.toml"
	appFileName    = "app.toml"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update app.toml and config.toml files to the latest version, default only missing parts are added",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			rootDir := serverCtx.Config.RootDir
			fileName := filepath.Join(rootDir, "config", configFileName)
			tmcfg.WriteConfigFile(fileName, serverCtx.Config)
			serverCtx.Logger.Info("Update config.toml is successful", "fileName", fileName)

			config.SetConfigTemplate(DefaultConfigTemplate())
			appConfig := DefaultConfig()
			if err := serverCtx.Viper.Unmarshal(appConfig); err != nil {
				return err
			}
			fileName = filepath.Join(rootDir, "config", appFileName)
			config.WriteConfigFile(fileName, appConfig)
			serverCtx.Logger.Info("Update app.toml is successful", "fileName", fileName)
			return nil
		},
	}
	return cmd
}

func appTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app.toml [key] [value]",
		Short: "Create or query an `~/.fxcore/config/apptoml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SetConfigTemplate(DefaultConfigTemplate())
			return runConfigCmd(cmd, append([]string{appFileName}, args...))
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

func configTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config.toml [key] [value]",
		Short: "Create or query an `~/.fxcore/config/config.toml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCmd(cmd, append([]string{configFileName}, args...))
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

func runConfigCmd(cmd *cobra.Command, args []string) error {
	serverCtx := server.GetServerContextFromCmd(cmd)
	clientCtx := client.GetClientContextFromCmd(cmd)

	configName := filepath.Join(serverCtx.Config.RootDir, "config", args[0])
	cfg, err := newConfig(serverCtx.Viper, configName)
	if err != nil {
		return err
	}

	// is len(args) == 1, get config file content
	if len(args) == 1 {
		return cfg.output(clientCtx)
	}

	// 2. is len(args) == 2, get config key value
	if len(args) == 2 {
		return output(clientCtx, serverCtx.Viper.Get(args[1]))
	}

	serverCtx.Viper.Set(args[1], args[2])
	return cfg.save()
}

type cmdConfig interface {
	save() error
	output(ctx client.Context) error
}

var (
	_ cmdConfig = &appTomlConfig{}
	_ cmdConfig = &configTomlConfig{}
)

type appTomlConfig struct {
	v          *viper.Viper
	config     *Config
	configName string
}

func (a *appTomlConfig) output(ctx client.Context) error {
	return output(ctx, a.config)
}

func (a *appTomlConfig) save() error {
	if err := a.v.Unmarshal(a.config); err != nil {
		return err
	}
	config.WriteConfigFile(a.configName, a.config)
	return nil
}

type configTomlConfig struct {
	v          *viper.Viper
	config     *tmcfg.Config
	configName string
}

func (c *configTomlConfig) output(ctx client.Context) error {
	type outputConfig struct {
		tmcfg.BaseConfig `mapstructure:",squash"`
		RPC              tmcfg.RPCConfig             `mapstructure:"rpc"`
		P2P              tmcfg.P2PConfig             `mapstructure:"p2p"`
		Mempool          tmcfg.MempoolConfig         `mapstructure:"mempool"`
		StateSync        tmcfg.StateSyncConfig       `mapstructure:"statesync"`
		FastSync         tmcfg.FastSyncConfig        `mapstructure:"fastsync"`
		Consensus        tmcfg.ConsensusConfig       `mapstructure:"consensus"`
		TxIndex          tmcfg.TxIndexConfig         `mapstructure:"tx_index"`
		Instrumentation  tmcfg.InstrumentationConfig `mapstructure:"instrumentation"`
	}
	return output(ctx, outputConfig{
		BaseConfig:      c.config.BaseConfig,
		RPC:             *c.config.RPC,
		P2P:             *c.config.P2P,
		Mempool:         *c.config.Mempool,
		StateSync:       *c.config.StateSync,
		FastSync:        *c.config.FastSync,
		Consensus:       *c.config.Consensus,
		TxIndex:         *c.config.TxIndex,
		Instrumentation: *c.config.Instrumentation,
	})
}

func (c *configTomlConfig) save() error {
	if err := c.v.Unmarshal(c.config); err != nil {
		return err
	}
	tmcfg.WriteConfigFile(c.configName, c.config)
	return nil
}

func newConfig(v *viper.Viper, configName string) (cmdConfig, error) {
	if strings.HasSuffix(configName, appFileName) {
		var configData = Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &appTomlConfig{config: &configData, v: v, configName: configName}, nil
	} else if strings.HasSuffix(configName, configFileName) {
		var configData = tmcfg.Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &configTomlConfig{config: &configData, v: v, configName: configName}, nil
	} else {
		return nil, fmt.Errorf("invalid config file: %s, (support: %v)", configName, strings.Join([]string{appFileName, configFileName}, "/"))
	}
}

func output(ctx client.Context, content interface{}) error {
	var mapData map[string]interface{}
	if err := mapstructure.Decode(content, &mapData); err != nil {
		var data interface{}
		if err := mapstructure.Decode(content, &data); err != nil {
			return err
		}
		return cli.PrintOutput(ctx, data)
	}
	return cli.PrintOutput(ctx, mapData)
}

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

	EVM     *ethermintconfig.EVMConfig     `mapstructure:"evm"`
	JSONRPC *ethermintconfig.JSONRPCConfig `mapstructure:"json-rpc"`
	TLS     *ethermintconfig.TLSConfig     `mapstructure:"tls"`
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config:       *config.DefaultConfig(),
		BypassMinFee: DefaultBypassMinFee(),
		EVM:          ethermintconfig.DefaultEVMConfig(),
		JSONRPC:      ethermintconfig.DefaultJSONRPCConfig(),
		TLS:          ethermintconfig.DefaultTLSConfig(),
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
