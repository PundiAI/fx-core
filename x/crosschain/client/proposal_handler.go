package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v7/x/crosschain/client/cli"
)

var LegacyUpdateChainOraclesProposalHandler = govclient.NewProposalHandler(NewLegacyUpdateChainOraclesProposalCmd)

func NewLegacyUpdateChainOraclesProposalCmd() *cobra.Command {
	var chainName string
	cmd := cli.CmdUpdateChainOraclesProposal(chainName)
	cmd.Flags().StringVarP(&chainName, "chain-name", "", "", "cross chain name")
	_ = cmd.MarkFlagRequired("chain-name")
	return cmd
}
