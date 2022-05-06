package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"

	fxconfig "github.com/functionx/fx-core/server/config"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
)

const (
	configFileName = "config.toml"
	appFileName    = "app.toml"
)

func appTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app.toml [key] [value]",
		Short: "Create or query an `~/.fxcore/config/apptoml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SetConfigTemplate(fxconfig.DefaultConfigTemplate())
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

	operatorConfig, err := newConfig(args[0], serverCtx)
	if err != nil {
		return err
	}

	// is len(args) == 1, get config file content
	if len(args) == 1 {
		return operatorConfig.output(clientCtx.PrintOutput)
	}

	// 2. is len(args) == 2, get config key value
	if len(args) == 2 {
		return output(clientCtx.PrintOutput, serverCtx.Viper.Get(args[1]))
	}

	serverCtx.Viper.Set(args[1], args[2])
	configPath := filepath.Join(serverCtx.Viper.GetString(flags.FlagHome), "config")
	if err = operatorConfig.save(configPath); err != nil {
		return err
	}
	return nil
}

type cmdConfig interface {
	save(configPath string) error
	output(printOutput func(out []byte) error) error
}

var (
	_ cmdConfig = &appTomlConfig{}
	_ cmdConfig = &configTomlConfig{}
)

type appTomlConfig struct {
	ctx    *server.Context
	config *fxconfig.Config
}

func (a *appTomlConfig) output(printOutput func(out []byte) error) error {
	return output(printOutput, a.config)
}

func (a *appTomlConfig) save(configPath string) error {
	if err := a.ctx.Viper.Unmarshal(a.config); err != nil {
		return err
	}
	configPath = filepath.Join(configPath, appFileName)
	config.WriteConfigFile(configPath, a.config)
	return nil
}

type configTomlConfig struct {
	ctx    *server.Context
	config *tmcfg.Config
}

func (c *configTomlConfig) output(printOutput func(out []byte) error) error {
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
	return output(printOutput, outputConfig{
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

func (c *configTomlConfig) save(configPath string) error {
	if err := c.ctx.Viper.Unmarshal(c.config); err != nil {
		return err
	}
	configPath = filepath.Join(configPath, configFileName)
	tmcfg.WriteConfigFile(configPath, c.config)
	return nil
}

func newConfig(configName string, ctx *server.Context) (cmdConfig, error) {
	switch configName {
	case appFileName:
		var configData = fxconfig.Config{}
		if err := ctx.Viper.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &appTomlConfig{config: &configData, ctx: ctx}, nil
	case configFileName:
		var configData = tmcfg.Config{}
		if err := ctx.Viper.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &configTomlConfig{config: &configData, ctx: ctx}, nil
	default:
		return nil, fmt.Errorf("invalid config file: %s, (support: %v)", configName, strings.Join([]string{appFileName, configFileName}, "/"))
	}
}

func output(printOutput func(out []byte) error, content interface{}) error {
	var mapData map[string]interface{}
	if err := mapstructure.Decode(content, &mapData); err != nil {
		return err
	}
	data, err := json.MarshalIndent(mapData, "", "  ")
	if err != nil {
		return err
	}
	return printOutput(data)
}
