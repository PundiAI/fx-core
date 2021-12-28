package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/functionx/fx-core/x/intrarelayer/client/cli"
	"github.com/functionx/fx-core/x/intrarelayer/client/rest"
)

var (
	InitIntrarelayerParamsProposalHandler = govclient.NewProposalHandler(cli.NewInitIntrarelayerParamsProposalCmd, rest.InitIntrarelayerParamsProposalRESTHandler)
	RegisterCoinProposalHandler           = govclient.NewProposalHandler(cli.NewRegisterCoinProposalCmd, rest.RegisterCoinProposalRESTHandler)
	RegisterERC20ProposalHandler          = govclient.NewProposalHandler(cli.NewRegisterERC20ProposalCmd, rest.RegisterERC20ProposalRESTHandler)
	ToggleTokenRelayProposalHandler       = govclient.NewProposalHandler(cli.NewToggleTokenRelayProposalCmd, rest.ToggleTokenRelayRESTHandler)
)
