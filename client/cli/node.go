package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/spf13/cobra"
)

func QueryGasPricesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gas-prices",
		Short: "Query node gas prices",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := node.NewQueryServer(clientCtx)
			res, err := queryClient.Config(context.Background(), &node.ConfigRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
