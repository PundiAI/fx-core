package keeper

import (
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func (keeper Keeper) IteratorInactiveProposal(ctx sdk.Context, t time.Time, fn func(proposal v1.Proposal) (bool, error)) error {
	rng := collections.NewPrefixUntilPairRange[time.Time, uint64](t)
	return keeper.InactiveProposalsQueue.Walk(ctx, rng, func(key1 collections.Pair[time.Time, uint64], _ uint64) (bool, error) {
		proposal, err := keeper.Proposals.Get(ctx, key1.K2())
		if err != nil {
			return false, err
		}
		return fn(proposal)
	})
}

func (keeper Keeper) IteratorActiveProposal(ctx sdk.Context, t time.Time, fn func(proposal v1.Proposal) (bool, error)) error {
	rngT := collections.NewPrefixUntilPairRange[time.Time, uint64](t)
	return keeper.ActiveProposalsQueue.Walk(ctx, rngT, func(key collections.Pair[time.Time, uint64], _ uint64) (bool, error) {
		proposal, err := keeper.Proposals.Get(ctx, key.K2())
		if err != nil {
			return false, err
		}
		return fn(proposal)
	})
}

func (keeper Keeper) HasDeposit(ctx sdk.Context, proposalId uint64, depositor sdk.AccAddress) (bool, error) {
	key := collections.Join(proposalId, depositor)
	return keeper.Deposits.Has(ctx, key)
}

func (keeper Keeper) HasVote(ctx sdk.Context, proposalId uint64, voter sdk.AccAddress) (bool, error) {
	key := collections.Join(proposalId, voter)
	return keeper.Votes.Has(ctx, key)
}
