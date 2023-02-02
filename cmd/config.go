package cmd

import (
	"path/filepath"

	sdkcfg "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	fxcfg "github.com/functionx/fx-core/v3/server/config"
)

const (
	configFileName = "config.toml"
	appFileName    = "app.toml"
)

// configCmd returns a CLI command to interactively create an application CLI
// config file.
func configCmd() *cobra.Command {
	cmd := sdkcfg.Cmd()
	cmd.AddCommand(
		updateCfgCmd(),
		appTomlCfgCmd(),
		configTomlCfgCmd(),
	)
	return cmd
}

func updateCfgCmd() *cobra.Command {
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

			config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
			appConfig := fxcfg.DefaultConfig()
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

func appTomlCfgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app.toml [key] [value]",
		Short: "Create or query an `config/app.toml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SetConfigTemplate(fxcfg.DefaultConfigTemplate())
			return fxcfg.NewConfigCmd(cmd, append([]string{appFileName}, args...))
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}

func configTomlCfgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config.toml [key] [value]",
		Short: "Create or query an `config/config.toml` file",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fxcfg.NewConfigCmd(cmd, append([]string{configFileName}, args...))
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}
