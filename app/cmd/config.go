package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/functionx/fx-core/app/fxcore"
	fxserver "github.com/functionx/fx-core/server/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
)

const (
	configFileName = "config.toml"
	appFileName    = "app.toml"
)

var (
	supportConfigs = []string{configFileName, appFileName}
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("config <%s> [key] [value]", strings.Join(supportConfigs, "/")),
		Short: "Update or query an application configuration file",
		Long: `
fxcored config app.toml   // 1. show app.toml content
fxcored config app.toml minimum-gas-prices  // 2. show app.toml minimul-gas-prices value
fxcored config app.toml minimum-gas-prices 4000000000000FX  // 3. update app.toml minimul-gas-prices value to 4000000000000FX
`,
		Args: cobra.RangeArgs(1, 3),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			customAppTemplate, _ := fxserver.AppConfig(fxcore.MintDenom)
			config.SetConfigTemplate(customAppTemplate)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			serverCtx := server.GetServerContextFromCmd(cmd)
			clientCtx := client.GetClientContextFromCmd(cmd)

			operatorConfig, err := newConfig(args[0], serverCtx)
			if err != nil {
				return err
			}

			// is len(args) == 1, get config file content
			if len(args) == 1 {
				return operatorConfig.output(clientCtx)
			}

			// 2. is len(args) == 2, get config key value
			if len(args) == 2 {
				return output(clientCtx, clientCtx.Viper.Get(args[1]))
			}

			serverCtx.Viper.Set(args[1], args[2])
			configPath := filepath.Join(serverCtx.Viper.GetString(flags.FlagHome), "config")
			if err = operatorConfig.save(serverCtx, configPath); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

type cmdConfig interface {
	output(clientCtx client.Context) error
	save(clientCtx *server.Context, configPath string) error
}

var (
	_ cmdConfig = appTomlConfig{}
	_ cmdConfig = configTomlConfig{}
)

type appTomlConfig struct {
	config *fxserver.Config
}

func (a appTomlConfig) output(clientCtx client.Context) error {
	return output(clientCtx, a.config)
}

func (a appTomlConfig) save(clientCtx *server.Context, configPath string) error {
	if err := clientCtx.Viper.Unmarshal(a.config); err != nil {
		return err
	}
	configPath = filepath.Join(configPath, appFileName)
	config.WriteConfigFile(configPath, a.config)
	return nil
}

type configTomlConfig struct {
	config *tmcfg.Config
}

func (c configTomlConfig) output(clientCtx client.Context) error {
	return output(clientCtx, c.config)
}

func (c configTomlConfig) save(clientCtx *server.Context, configPath string) error {
	if err := clientCtx.Viper.Unmarshal(c.config); err != nil {
		return err
	}
	configPath = filepath.Join(configPath, configFileName)
	tmcfg.WriteConfigFile(configPath, c.config)
	return nil
}

func newConfig(configName string, clientCtx *server.Context) (cmdConfig, error) {
	switch configName {
	case appFileName:
		var configData = fxserver.Config{}
		if err := clientCtx.Viper.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &appTomlConfig{config: &configData}, nil
	case configFileName:
		var configData = tmcfg.Config{}
		if err := clientCtx.Viper.Unmarshal(&configData); err != nil {
			return nil, err
		}
		return &configTomlConfig{config: &configData}, nil
	default:
		return nil, errors.New(fmt.Sprintf("invalid config file:%s, (support: %v)", configName, strings.Join(supportConfigs, "/")))
	}
}

func output(clientCtx client.Context, content interface{}) error {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	return PrintOutput(clientCtx, data)
}
