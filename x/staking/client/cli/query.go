package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"

	fxstakingtypes "github.com/functionx/fx-core/v3/x/staking/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingQueryCmd.AddCommand(
		cli.GetCmdQueryDelegation(),
		cli.GetCmdQueryDelegations(),
		cli.GetCmdQueryUnbondingDelegation(),
		cli.GetCmdQueryUnbondingDelegations(),
		cli.GetCmdQueryRedelegation(),
		cli.GetCmdQueryRedelegations(),
		cli.GetCmdQueryValidator(),
		cli.GetCmdQueryValidators(),
		cli.GetCmdQueryValidatorDelegations(),
		cli.GetCmdQueryValidatorUnbondingDelegations(),
		cli.GetCmdQueryValidatorRedelegations(),
		cli.GetCmdQueryHistoricalInfo(),
		cli.GetCmdQueryParams(),
		cli.GetCmdQueryPool(),
		GetCmdQueryValidatorLPToken(),
	)

	return stakingQueryCmd
}

// GetCmdQueryValidatorLPToken implements the ValidatorLPToken query command.
func GetCmdQueryValidatorLPToken() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()
	cmd := &cobra.Command{
		Use:   "validator-lp-token [validator-addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the validator lp token address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query lp token address.

Example:
$ %s query staking validator-lp-token %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryClient := fxstakingtypes.NewQueryClient(clientCtx)
			res, err := queryClient.ValidatorLPToken(cmd.Context(), &fxstakingtypes.QueryValidatorLPTokenRequest{
				ValidatorAddr: valAddr.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.LpToken)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
