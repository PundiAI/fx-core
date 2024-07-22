package cli

import (
	"context"
	"encoding/json"
	"strconv"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// BlockCommand returns the verified block data for a given heights
func BlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [height]",
		Short: "Get verified data for a the block at given height",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			var height *int64

			// optional height
			if len(args) > 0 {
				h, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				if h > 0 {
					tmp := int64(h)
					height = &tmp
				}
			}

			// get the node
			node, err := clientCtx.GetNode()
			if err != nil {
				return err
			}

			// header -> BlockchainInfo
			// header, tx -> Block
			// results -> BlockResults
			res, err := node.Block(context.Background(), height)
			if err != nil {
				return err
			}

			raw, err := json.Marshal(res)
			if err != nil {
				return err
			}
			return clientCtx.PrintRaw(raw)
		},
	}
	cmd.Flags().StringP(tmcli.OutputFlag, "o", "json", "Output format (text|json)")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	return cmd
}
