package cli

import (
	"encoding/hex"

	"github.com/tendermint/tendermint/libs/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func QueryStoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store <store name> <hex key>",
		Short: "Query for a chain store",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			hexKey, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			data, height, err := clientCtx.QueryStore(hexKey, args[0])
			if err != nil {
				return err
			}
			bts, err := json.Marshal(map[string]interface{}{
				"BlockHeight": height,
				"Data":        hex.EncodeToString(data),
			})
			if err != nil {
				return err
			}
			return PrintOutput(clientCtx, bts)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
