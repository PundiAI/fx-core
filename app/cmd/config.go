package cmd

import (
	"encoding/json"
	"fmt"
	fxconfig "github.com/functionx/fx-core/server/config"
	"path/filepath"
	"strings"

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

func AppTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app.toml [key] [value]",
		Short: "Create or query an `.fxcore/config/apptoml` file",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.SetConfigTemplate(fxconfig.DefaultConfigTemplate())
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCmd(cmd, append([]string{appFileName}, args...))
		},
		Args: cobra.RangeArgs(0, 2),
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

func ConfigTomlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config.toml [key] [value]",
		Short: "Create or query an `.fxcore/config/config.toml` file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCmd(cmd, append([]string{configFileName}, args...))
		},
		Args: cobra.RangeArgs(0, 2),
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
		return operatorConfig.output(clientCtx)
	}

	// 2. is len(args) == 2, get config key value
	if len(args) == 2 {
		return output(clientCtx, serverCtx.Viper.Get(args[1]))
	}

	serverCtx.Viper.Set(args[1], args[2])
	configPath := filepath.Join(serverCtx.Viper.GetString(flags.FlagHome), "config")
	if err = operatorConfig.save(serverCtx, configPath); err != nil {
		return err
	}
	return nil
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
	config *fxconfig.Config
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
		var configData = fxconfig.Config{}
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
		return nil, fmt.Errorf("invalid config file:%s, (support: %v)", configName, strings.Join([]string{appFileName, configFileName}, "/"))
	}
}

func output(clientCtx client.Context, content interface{}) error {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	return clientCtx.PrintOutput(data)
}
