package cli

import (
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

// PreUpgradeCmd called by cosmovisor
func PreUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "Pre-upgrade called by cosmovisor, before migrations upgrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			serverCtx.Logger.Info("pre-upgrade", "action", "update app.toml and config.toml")
			return updateCmd().RunE(cmd, args)
		},
	}
	return cmd
}
