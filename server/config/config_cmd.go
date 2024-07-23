package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CmdHandler(cmd *cobra.Command, args []string) error {
	serverCtx := server.GetServerContextFromCmd(cmd)
	clientCtx := client.GetClientContextFromCmd(cmd)

	configName := filepath.Join(serverCtx.Config.RootDir, "config", args[0])
	cfg, err := NewConfig(serverCtx.Viper, configName)
	if err != nil {
		return err
	}

	// is len(args) == 1, get Config file content
	if len(args) == 1 {
		return cfg.Output(clientCtx)
	}

	// 2. is len(args) == 2, get Config key value
	if len(args) == 2 {
		return Output(clientCtx, serverCtx.Viper.Get(args[1]))
	}

	serverCtx.Viper.Set(args[1], args[2])
	return cfg.Save()
}

type CmdConfig interface {
	Save() error
	Output(ctx client.Context) error
}

var (
	_ CmdConfig = &AppToml{}
	_ CmdConfig = &TmConfigToml{}
)

type AppToml struct {
	v          *viper.Viper
	Config     *Config
	configName string
}

func (a *AppToml) Output(ctx client.Context) error {
	return Output(ctx, a.Config)
}

func (a *AppToml) Save() error {
	if err := a.v.Unmarshal(a.Config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ZeroFields = true
	}); err != nil {
		return err
	}
	srvconfig.WriteConfigFile(a.configName, a.Config)
	return nil
}

type TmConfigToml struct {
	v          *viper.Viper
	Config     *tmcfg.Config
	configName string
}

func (c *TmConfigToml) Output(ctx client.Context) error {
	type outputConfig struct {
		tmcfg.BaseConfig `mapstructure:",squash"`
		RPC              tmcfg.RPCConfig             `mapstructure:"rpc"`
		P2P              tmcfg.P2PConfig             `mapstructure:"p2p"`
		Mempool          tmcfg.MempoolConfig         `mapstructure:"mempool"`
		StateSync        tmcfg.StateSyncConfig       `mapstructure:"statesync"`
		BlockSync        tmcfg.BlockSyncConfig       `mapstructure:"blocksync"`
		Consensus        tmcfg.ConsensusConfig       `mapstructure:"consensus"`
		Storage          tmcfg.StorageConfig         `mapstructure:"storage"`
		TxIndex          tmcfg.TxIndexConfig         `mapstructure:"tx_index"`
		Instrumentation  tmcfg.InstrumentationConfig `mapstructure:"instrumentation"`
	}
	return Output(ctx, outputConfig{
		BaseConfig:      c.Config.BaseConfig,
		RPC:             *c.Config.RPC,
		P2P:             *c.Config.P2P,
		Mempool:         *c.Config.Mempool,
		StateSync:       *c.Config.StateSync,
		BlockSync:       *c.Config.BlockSync,
		Consensus:       *c.Config.Consensus,
		Storage:         *c.Config.Storage,
		TxIndex:         *c.Config.TxIndex,
		Instrumentation: *c.Config.Instrumentation,
	})
}

func (c *TmConfigToml) Save() error {
	if err := c.v.Unmarshal(c.Config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.ZeroFields = true
	}); err != nil {
		return err
	}
	tmcfg.WriteConfigFile(c.configName, c.Config)
	return nil
}

func NewConfig(v *viper.Viper, configName string) (CmdConfig, error) {
	if strings.HasSuffix(configName, "app.toml") {
		configData := Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &AppToml{Config: &configData, v: v, configName: configName}, nil
	} else if strings.HasSuffix(configName, "config.toml") {
		configData := tmcfg.Config{}
		if err := v.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &TmConfigToml{Config: &configData, v: v, configName: configName}, nil
	} else {
		return nil, fmt.Errorf("invalid Config file: %s, (support: %v)", configName, strings.Join([]string{"app.toml", "config.toml"}, "/"))
	}
}

func Output(ctx client.Context, content interface{}) error {
	var mapData map[string]interface{}
	if err := mapstructure.Decode(content, &mapData); err != nil {
		var data interface{}
		if err = mapstructure.Decode(content, &data); err != nil {
			return err
		}
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return ctx.PrintRaw(raw)
	}
	raw, err := json.Marshal(mapData)
	if err != nil {
		return err
	}
	return ctx.PrintRaw(raw)
}
