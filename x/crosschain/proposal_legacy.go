// nolint:staticcheck
package crosschain

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// NewCrosschainProposalHandler
// Deprecated: instead of defining gov router proposal handlers, the v0.46 gov execution models is based on sdk.Msgs
func NewCrosschainProposalHandler(k keeper.RouterKeeper) govv1betal.Handler {
	return func(ctx sdk.Context, content govv1betal.Content) error {
		switch c := content.(type) {
		case *types.UpdateChainOraclesProposal:
			router := k.Router()
			if !router.HasRoute(c.ChainName) {
				return errorsmod.Wrap(errortypes.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", c.ChainName))
			}
			return router.GetRoute(c.ChainName).ProposalServer.UpdateChainOraclesProposal(ctx, c)
		default:
			return errorsmod.Wrapf(errortypes.ErrUnknownRequest, "Unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}
