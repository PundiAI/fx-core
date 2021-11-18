package cmd

import (
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <config file name> <key> [value]",
		Short: "Update or query an application configuration file",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			if args[0] != "app.toml" && args[0] != "config.toml" {
				return errors.New("invalid config file(support: app.toml,config.toml)")
			}
			clientCtx := client.GetClientContextFromCmd(cmd)
			if len(args) == 2 {
				data, err := json.MarshalIndent(serverCtx.Viper.Get(args[1]), "", "\t")
				if err != nil {
					return err
				}
				return PrintOutput(clientCtx, data)
			}
			if len(args) == 3 {
				serverCtx.Viper.Set(args[1], args[2])
			}
			rootDir := serverCtx.Viper.GetString(flags.FlagHome)
			configPath := filepath.Join(rootDir, "config")
			switch args[0] {
			case "app.toml":
				var appConfig = config.Config{}
				if err := serverCtx.Viper.Unmarshal(&appConfig); err != nil {
					return err
				}
				if len(args) == 1 {
					data, err := json.MarshalIndent(appConfig, "", "  ")
					if err != nil {
						return err
					}
					return PrintOutput(clientCtx, data)
				}
				configPath = filepath.Join(configPath, "app.toml")
				config.WriteConfigFile(configPath, &appConfig)
			case "config.toml":
				var tmConfig = tmcfg.Config{}
				if err := serverCtx.Viper.Unmarshal(&tmConfig); err != nil {
					return err
				}
				if len(args) == 1 {
					data, err := json.MarshalIndent(tmConfig, "", "  ")
					if err != nil {
						return err
					}
					return PrintOutput(clientCtx, data)
				}
				configPath := filepath.Join(configPath, "config.toml")
				tmcfg.WriteConfigFile(configPath, &tmConfig)
			}
			return nil
		},
	}
	return cmd
}
