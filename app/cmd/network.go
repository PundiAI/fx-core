package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/functionx/fx-core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func Network() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "network",
		Args:    cobra.NoArgs,
		Short:   "Show fxcored network and upgrade info",
		Example: "fxcored network",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			outputBytes, err := json.Marshal(map[string]interface{}{
				"Network":                                types.Network(),
				"GravityPruneValsetsAndAttestationBlock": fmt.Sprintf("%d", types.GravityPruneValsetsAndAttestationBlock()),
				"GravityValsetSlashBlock":                fmt.Sprintf("%d", types.GravityValsetSlashBlock()),
				"CrossChainSupportBscBlock":              fmt.Sprintf("%d", types.CrossChainSupportBscBlock()),
				"CrossChainSupportTronBlock":             fmt.Sprintf("%d", types.CrossChainSupportTronBlock()),
				"EIP155ChainID":                          fmt.Sprintf("%d", types.EIP155ChainID()),
				"EvmSupportBlock":                        fmt.Sprintf("%d", types.EvmSupportBlock()),
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
