package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v7/x/migrate/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the migrate module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdGetMigrateRecord(), CmdGetMigrateAccount())
	return cmd
}

func CmdGetMigrateAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account [bech32/hex address]",
		Short:   "Query address migrate info",
		Example: fmt.Sprintf("$ %s q migrate account fx1plg.../0xdf9...", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			input := args[0]
			var addr sdk.AccAddress
			if common.IsHexAddress(input) {
				hexAddr := common.HexToAddress(input)
				addr = hexAddr.Bytes()
			} else {
				var err error
				addr, err = sdk.AccAddressFromBech32(input)
				if err != nil {
					return fmt.Errorf("invalid address %s, must be hex or bech32", input)
				}
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			bankClient := banktypes.NewQueryClient(clientCtx)
			balances, err := bankClient.AllBalances(cmd.Context(), &banktypes.QueryAllBalancesRequest{
				Address:    addr.String(),
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			stakingClient := stakingtypes.NewQueryClient(clientCtx)
			delegations, err := stakingClient.DelegatorDelegations(cmd.Context(), &stakingtypes.QueryDelegatorDelegationsRequest{
				DelegatorAddr: addr.String(),
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}
			unbondingDelegations, err := stakingClient.DelegatorUnbondingDelegations(cmd.Context(), &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
				DelegatorAddr: addr.String(),
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}
			redelegations, err := stakingClient.Redelegations(cmd.Context(), &stakingtypes.QueryRedelegationsRequest{
				DelegatorAddr: addr.String(),
				Pagination:    pageReq,
			})
			if err != nil {
				return err
			}
			govClient := govv1betal.NewQueryClient(clientCtx)
			depositProposals, err := govClient.Proposals(cmd.Context(), &govv1betal.QueryProposalsRequest{
				ProposalStatus: govv1betal.StatusDepositPeriod,
				Depositor:      addr.String(),
				Pagination:     pageReq,
			})
			if err != nil {
				return err
			}
			voteProposals, err := govClient.Proposals(cmd.Context(), &govv1betal.QueryProposalsRequest{
				ProposalStatus: govv1betal.StatusVotingPeriod,
				Voter:          addr.String(),
				Pagination:     pageReq,
			})
			if err != nil {
				return err
			}

			info := map[string]interface{}{
				"bank": map[string]interface{}{
					"balance": balances.Balances,
				},
				"staking": map[string]interface{}{
					"delegate":     delegations.DelegationResponses,
					"unbonding":    unbondingDelegations.UnbondingResponses,
					"redelegation": redelegations.RedelegationResponses,
				},
				"gov": map[string]interface{}{
					"deposit": depositProposals.Proposals,
					"vote":    voteProposals.Proposals,
				},
			}
			bz, err := json.Marshal(info)
			if err != nil {
				return err
			}
			return clientCtx.PrintBytes(bz)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "migrate info")
	return cmd
}

func CmdGetMigrateRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "record [bech32/hex address]",
		Short:   "Query the migrate record of address",
		Example: fmt.Sprintf("$ %s q migrate record fx1plg.../0xdf9...", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var addr string
			if common.IsHexAddress(args[0]) {
				addr = common.HexToAddress(args[0]).String()
			} else {
				if acc, err := sdk.AccAddressFromBech32(args[0]); err == nil {
					addr = acc.String()
				}
			}
			if len(addr) == 0 {
				return fmt.Errorf("must be bech32 or hex address: %s", args[0])
			}

			res, err := queryClient.MigrateRecord(cmd.Context(), &types.QueryMigrateRecordRequest{Address: addr})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
