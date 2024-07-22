package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

// GenesisCoreCommand adds core sdk's sub-commands into genesis command:
// -> gentx, migrate, collect-gentxs, validate-genesis, add-genesis-account
func GenesisCoreCommand(moduleBasics module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "genesis",
		Short:                      "Application's genesis-related subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		AddGenesisAccountCmd(defaultNodeHome),
		GenTxCmd(moduleBasics, banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		cli.ValidateGenesisCmd(moduleBasics),
	)

	return cmd
}
