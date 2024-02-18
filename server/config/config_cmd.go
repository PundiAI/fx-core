package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/functionx/fx-core/v7/client/cli"
)

func CmdHandler(cmd *cobra.Command, args []string) error {
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
	if err := a.v.Unmarshal(a.config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ZeroFields = true
	}); err != nil {
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
	if err := c.v.Unmarshal(c.config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ZeroFields = true
	}); err != nil {
		return err
	}
	tmcfg.WriteConfigFile(c.configName, c.config)
	return nil
}

func newConfig(v *viper.Viper, configName string) (cmdConfig, error) {
	if strings.HasSuffix(configName, "app.toml") {
		configData := Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &appTomlConfig{config: &configData, v: v, configName: configName}, nil
	} else if strings.HasSuffix(configName, "config.toml") {
		configData := tmcfg.Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &configTomlConfig{config: &configData, v: v, configName: configName}, nil
	} else {
		return nil, fmt.Errorf("invalid config file: %s, (support: %v)", configName, strings.Join([]string{"app.toml", "config.toml"}, "/"))
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
