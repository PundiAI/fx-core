package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/functionx/fx-core/v7/x/erc20/client/cli"
)

var (
	LegacyRegisterCoinProposalHandler          = govclient.NewProposalHandler(cli.NewLegacyRegisterCoinProposalCmd)
	LegacyRegisterERC20ProposalHandler         = govclient.NewProposalHandler(cli.NewLegacyRegisterERC20ProposalCmd)
	LegacyToggleTokenConversionProposalHandler = govclient.NewProposalHandler(cli.NewLegacyToggleTokenConversionProposalCmd)
	LegacyUpdateDenomAliasProposalHandler      = govclient.NewProposalHandler(cli.NewLegacyUpdateDenomAliasProposalCmd)
)
