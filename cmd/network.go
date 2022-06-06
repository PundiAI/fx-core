package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/app/cli"

	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/functionx/fx-core/types"
)

func networkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Args:    cobra.NoArgs,
		Short:   "Show fxcored network and upgrade info",
		Example: "fxcored network",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			output := map[string]interface{}{
				"ChainId":       types.ChainID,
				"Network":       types.Network(),
				"EIP155ChainID": fmt.Sprintf("%d", types.EIP155ChainID()),
			}
			return cli.PrintOutput(clientCtx, output)
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}
