package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// Deprecated
type proposalServer interface {
	UpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error // nolint:staticcheck
}

var _ proposalServer = Keeper{}

// UpdateChainOraclesProposal
// Deprecated: v0.46 gov execution models is based on sdk.Msgs
func (k Keeper) UpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error { // nolint:staticcheck
	k.Logger(ctx).Info("handle update chain oracles proposal", "proposal", proposal.String())
	return k.UpdateProposalOracles(ctx, proposal.Oracles)
}
