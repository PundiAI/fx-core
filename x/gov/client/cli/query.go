package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/spf13/cobra"

	fxgovtypes "github.com/functionx/fx-core/v7/x/gov/types"
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
		cli.GetCmdQueryProposal(),
		cli.GetCmdQueryProposals(),
		cli.GetCmdQueryVote(),
		cli.GetCmdQueryVotes(),
		GetCmdQueryParams(),
		GetCmdQueryEGFParams(),
		cli.GetCmdQueryProposer(),
		cli.GetCmdQueryDeposit(),
		cli.GetCmdQueryDeposits(),
		cli.GetCmdQueryTally(),
	)

	return govQueryCmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %s query gov params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := fxgovtypes.NewQueryClient(clientCtx)
			ctx := cmd.Context()
			msgType, err := cmd.Flags().GetString("msg-type")
			if err != nil {
				return err
			}
			res, err := queryClient.Params(
				ctx,
				&fxgovtypes.QueryParamsRequest{MsgType: msgType},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String("msg-type", "", "proto name to the type (pointer to struct) of the protocol buffer")
	return cmd
}

// GetCmdQueryEGFParams implements the query params command.
func GetCmdQueryEGFParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "egf-params",
		Short: "Query the EGF parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the EGF parameters for the governance process.

Example:
$ %s query gov egf-params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := fxgovtypes.NewQueryClient(clientCtx)
			ctx := cmd.Context()
			res, err := queryClient.EGFParams(
				ctx,
				&fxgovtypes.QueryEGFParamsRequest{},
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
