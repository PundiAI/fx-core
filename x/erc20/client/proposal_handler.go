package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/functionx/fx-core/x/erc20/client/cli"
	"github.com/functionx/fx-core/x/erc20/client/rest"
)

var (
	RegisterCoinProposalHandler     = govclient.NewProposalHandler(cli.NewRegisterCoinProposalCmd, rest.RegisterCoinProposalRESTHandler)
	ToggleTokenRelayProposalHandler = govclient.NewProposalHandler(cli.NewToggleTokenRelayProposalCmd, rest.ToggleTokenRelayRESTHandler)
	RegisterERC20ProposalHandler    = govclient.NewProposalHandler(cli.NewRegisterERC20ProposalCmd, rest.RegisterERC20ProposalRESTHandler)
)
