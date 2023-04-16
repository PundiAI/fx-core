package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v4/x/crosschain/types"
)

// nolint:staticcheck
// Deprecated
type proposalServer interface {
	UpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error
}

var _ proposalServer = Keeper{}

// UpdateChainOraclesProposal
// nolint:staticcheck
// Deprecated: v0.46 gov execution models is based on sdk.Msgs
func (k Keeper) UpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error {
	k.Logger(ctx).Info("handle update chain oracles proposal", "proposal", proposal.String())
	return k.UpdateChainOracles(ctx, proposal.Oracles)
}
