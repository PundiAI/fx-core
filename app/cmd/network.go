package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/functionx/fx-core/app/fxcore"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/app"
)

func Network() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Args:    cobra.NoArgs,
		Short:   "Show fxcored network and upgrade info",
		Example: "fxcored network",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			outputBytes, err := json.Marshal(map[string]interface{}{
				"ChainId":                                fxcore.ChainID,
				"Network":                                app.Network(),
				"GravityPruneValsetsAndAttestationBlock": fmt.Sprintf("%d", app.GravityPruneValsetsAndAttestationBlock()),
				"GravityValsetSlashBlock":                fmt.Sprintf("%d", app.GravityValsetSlashBlock()),
				"CrossChainSupportBscBlock":              fmt.Sprintf("%d", app.CrossChainSupportBscBlock()),
				"CrossChainSupportTronBlock":             fmt.Sprintf("%d", app.CrossChainSupportTronBlock()),
				"CrossChainSupportPolygonBlock":          fmt.Sprintf("%d", app.CrossChainSupportPolygonBlock()),
			})
			if err != nil {
				return err
			}
			return PrintOutput(clientCtx, outputBytes)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
