package cmd

import (
	"encoding/hex"
	"fmt"
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
			fmt.Printf("block height: %d, data: %s\n", height, hex.EncodeToString(data))
			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
