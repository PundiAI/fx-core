// nolint:staticcheck
package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func CmdUpdateChainOraclesProposal(chainName string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-crosschain-oracles [oracles]",
		Short: "Submit a update cross chain oracles proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			oracles := strings.Split(args[0], ",")
			for i, oracle := range oracles {
				oracleAddr, err := sdk.AccAddressFromBech32(oracle)
				if err != nil {
					return err
				}
				oracles[i] = oracleAddr.String()
			}
			proposal := &types.UpdateChainOraclesProposal{
				Title:       title,
				Description: description,
				Oracles:     oracles,
				ChainName:   chainName,
			}
			fromAddress := cliCtx.GetFromAddress()
			msg, err := govv1betal.NewMsgSubmitProposal(proposal, deposit, fromAddress)
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)
	_ = cmd.MarkFlagRequired(cli.FlagDeposit)
	return cmd
}
