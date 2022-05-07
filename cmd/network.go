package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

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
			outputBytes, err := json.Marshal(map[string]interface{}{
				"ChainId":                               types.ChainID,
				"Network":                               types.Network(),
				"GravityPruneValsetAndAttestationBlock": fmt.Sprintf("%d", types.GravityPruneValsetAndAttestationBlock()),
				"GravityValsetSlashBlock":               fmt.Sprintf("%d", types.GravityValsetSlashBlock()),
				"CrossChainSupportBscBlock":             fmt.Sprintf("%d", types.CrossChainSupportBscBlock()),
				"CrossChainSupportPolygonAndTronBlock":  fmt.Sprintf("%d", types.CrossChainSupportPolygonAndTronBlock()),
				"EIP155ChainID":                         fmt.Sprintf("%d", types.EIP155ChainID()),
				"EvmV0SupportBlock":                     fmt.Sprintf("%d", types.EvmV0SupportBlock()),
				"EvmV1SupportBlock":                     fmt.Sprintf("%d", types.EvmV1SupportBlock()),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(outputBytes)
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "text", "Output format (text|json)")
	return cmd
}
