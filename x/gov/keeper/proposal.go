package keeper

import (
	"context"
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

// ActivateVotingPeriod activates the voting period of a proposal
func (keeper Keeper) ActivateVotingPeriod(ctx context.Context, proposal v1.Proposal) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	startTime := sdkCtx.BlockHeader().Time
	proposal.VotingStartTime = &startTime
	var votingPeriod *time.Duration
	params, err := keeper.Params.Get(ctx)
	if err != nil {
		return err
	}

	if proposal.Expedited {
		votingPeriod = params.ExpeditedVotingPeriod
	} else {
		votingPeriod = params.VotingPeriod
	}

	votingPeriod = keeper.GetCustomMsgVotingPeriod(ctx, votingPeriod, proposal)

	endTime := proposal.VotingStartTime.Add(*votingPeriod)
	proposal.VotingEndTime = &endTime
	proposal.Status = v1.StatusVotingPeriod
	err = keeper.SetProposal(ctx, proposal)
	if err != nil {
		return err
	}

	err = keeper.InactiveProposalsQueue.Remove(ctx, collections.Join(*proposal.DepositEndTime, proposal.Id))
	if err != nil {
		return err
	}

	return keeper.ActiveProposalsQueue.Set(ctx, collections.Join(*proposal.VotingEndTime, proposal.Id), proposal.Id)
}

func (keeper Keeper) GetCustomMsgVotingPeriod(ctx context.Context, defaultVotingPeriod *time.Duration, proposal v1.Proposal) *time.Duration {
	msgType := getProposalMsgType(proposal)
	if customParams, found := keeper.GetCustomParams(ctx, msgType); found {
		return customParams.VotingPeriod
	}
	return defaultVotingPeriod
}

func getProposalMsgType(proposal v1.Proposal) string {
	message := proposal.GetMessages()
	for _, msg := range message {
		return sdk.MsgTypeURL(msg)
	}
	return ""
}
