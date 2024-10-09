package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	fxgovtypes "github.com/functionx/fx-core/v8/x/gov/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(
		GetCmdQuerySwitchParams(),
		GetCmdQueryCustomParams(),
	)

	return govQueryCmd
}

func GetCmdQuerySwitchParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch-params",
		Short: "Query the Switch parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the Switch parameters for the governance process.
Example:
$ %s query gov switch-params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := fxgovtypes.NewQueryClient(clientCtx)
			ctx := cmd.Context()
			res, err := queryClient.SwitchParams(ctx, &fxgovtypes.QuerySwitchParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryCustomParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-params",
		Short: "Query the custom params by msg url of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the msgUrl custom param for the governance process.

Example:
$ %s query gov custom-params /cosmos.distribution.v1beta1.MsgCommunityPoolSpend
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := fxgovtypes.NewQueryClient(clientCtx)
			ctx := cmd.Context()

			res, err := queryClient.CustomParams(
				ctx,
				&fxgovtypes.QueryCustomParamsRequest{MsgUrl: args[0]},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
