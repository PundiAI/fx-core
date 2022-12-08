package crosschain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/functionx/fx-core/v3/x/crosschain/keeper"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func NewChainProposalHandler(k keeper.RouterKeeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UpdateChainOraclesProposal:
			router := k.Router()
			if !router.HasRoute(c.ChainName) {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", c.ChainName))
			}
			return router.GetRoute(c.ChainName).ProposalServer.UpdateChainOraclesProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}
