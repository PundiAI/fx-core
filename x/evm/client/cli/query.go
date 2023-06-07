package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/evmos/ethermint/x/evm/client/cli"
	"github.com/evmos/ethermint/x/evm/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the evm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cli.GetStorageCmd(),
		cli.GetCodeCmd(),
		cli.GetParamsCmd(),
	)
	return cmd
}
