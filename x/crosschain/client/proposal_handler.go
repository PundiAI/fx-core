package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v3/x/crosschain/client/cli"
)

var UpdateChainOraclesProposalHandler = govclient.NewProposalHandler(NewLegacyUpdateChainOraclesProposalCmd, func(context client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update_chain_oracles",
		Handler:  nil,
	}
})

func NewLegacyUpdateChainOraclesProposalCmd() *cobra.Command {
	var chainName string
	cmd := cli.CmdUpdateChainOraclesProposal(chainName)
	cmd.Flags().StringVarP(&chainName, "chain-name", "", "", "cross chain name")
	_ = cmd.MarkFlagRequired("chain-name")
	return cmd
}
