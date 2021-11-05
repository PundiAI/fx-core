package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
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
				return errors.New("invalid config file")
			}
			if len(args) == 2 {
				fmt.Println(serverCtx.Viper.Get(args[1]))
				return nil
			}
			if len(args) <= 1 {
				data, err := json.MarshalIndent(serverCtx.Viper.AllSettings(), "", "\t")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}
			serverCtx.Viper.Set(args[1], args[2])

			rootDir := serverCtx.Viper.GetString(flags.FlagHome)
			configPath := filepath.Join(rootDir, "config")
			tmCfgFile := filepath.Join(configPath, "config.toml")
			appCfgFile := filepath.Join(configPath, "app.toml")

			switch args[0] {
			case "app.toml":
				var appConfig = config.Config{}
				if err := serverCtx.Viper.Unmarshal(&appConfig); err != nil {
					return err
				}
				config.WriteConfigFile(appCfgFile, &appConfig)
			case "config.toml":
				var tmConfig = tmcfg.Config{}
				if err := serverCtx.Viper.Unmarshal(&tmConfig); err != nil {
					return err
				}
				tmcfg.WriteConfigFile(tmCfgFile, &tmConfig)
			}

			fmt.Printf("update configuration file: %s\n", configPath)
			return nil
		},
	}
	return cmd
}
