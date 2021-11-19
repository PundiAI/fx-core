package cmd

import (
	sdkclient "github.com/cosmos/cosmos-sdk/client/config"
	"github.com/spf13/cobra"
)

func ClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "Fxcored client",
	}
	cmd.AddCommand(sdkclient.Cmd())
	return cmd
}
