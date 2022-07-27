package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/functionx/fx-core/app/cli"

	fxtypes "github.com/functionx/fx-core/types"

	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
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
				"ChainId":                                fxtypes.ChainID,
				"Network":                                fxtypes.Network(),
				"GravityPruneValsetsAndAttestationBlock": fmt.Sprintf("%d", fxtypes.GravityPruneValsetsAndAttestationBlock()),
				"GravityValsetSlashBlock":                fmt.Sprintf("%d", fxtypes.GravityValsetSlashBlock()),
				"CrossChainSupportBscBlock":              fmt.Sprintf("%d", fxtypes.CrossChainSupportBscBlock()),
				"CrossChainSupportTronBlock":             fmt.Sprintf("%d", fxtypes.CrossChainSupportTronBlock()),
				"CrossChainSupportPolygonBlock":          fmt.Sprintf("%d", fxtypes.CrossChainSupportPolygonBlock()),
			})
			if err != nil {
				return err
			}
			return cli.PrintOutput(clientCtx, outputBytes)
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}
