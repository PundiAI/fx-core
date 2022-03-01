package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	typescommon "github.com/functionx/fx-core/x/migrate/types/common"
	typesv1 "github.com/functionx/fx-core/x/migrate/types/v1"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        typescommon.ModuleName,
		Short:                      "Querying commands for the eth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdGetMigrateRecord())

	for _, command := range cmd.Commands() {
		flags.AddQueryFlagsToCmd(command)
	}
	return cmd
}

func CmdGetMigrateRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate-record [address]",
		Short:   "Query the migrate record of address",
		Example: fmt.Sprintf("%s q migrate migrate-record fx1plglgtkj4kj7z2q0jqgyw8exfnahwu8rlu6kzm", version.AppName),
		Aliases: []string{"mr"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := typesv1.NewQueryClient(clientCtx)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := queryClient.QueryMigrateRecord(cmd.Context(), &typesv1.QueryMigrateRecordRequest{Address: addr.String()})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
