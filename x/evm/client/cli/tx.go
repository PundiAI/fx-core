package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/evmos/ethermint/x/evm/client/cli"
	"github.com/spf13/cobra"

	"github.com/evmos/ethermint/x/evm/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "EVM transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(cli.NewRawTxCmd())
	return cmd
}
