package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/functionx/fx-core/v3/x/erc20/client/cli"
)

var (
	RegisterCoinProposalHandler          = govclient.NewProposalHandler(cli.NewRegisterCoinProposalCmd)
	RegisterERC20ProposalHandler         = govclient.NewProposalHandler(cli.NewRegisterERC20ProposalCmd)
	ToggleTokenConversionProposalHandler = govclient.NewProposalHandler(cli.NewToggleTokenConversionProposalCmd)
	UpdateDenomAliasProposalHandler      = govclient.NewProposalHandler(cli.NewUpdateDenomAliasProposalCmd)
)
