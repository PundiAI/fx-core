package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/functionx/fx-core/x/intrarelayer/client/cli"
	"github.com/functionx/fx-core/x/intrarelayer/client/rest"
)

var (
	RegisterCoinProposalHandler     = govclient.NewProposalHandler(cli.NewRegisterCoinProposalCmd, rest.RegisterCoinProposalRESTHandler)
	RegisterFIP20ProposalHandler    = govclient.NewProposalHandler(cli.NewRegisterFIP20ProposalCmd, rest.RegisterFIP20ProposalRESTHandler)
	ToggleTokenRelayProposalHandler = govclient.NewProposalHandler(cli.NewToggleTokenRelayProposalCmd, rest.ToggleTokenRelayRESTHandler)
)
