package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/functionx/fx-core/x/migrate/types"
)

type GovMigrate struct {
	govKey    sdk.StoreKey
	govKeeper types.GovKeeper
}

func NewGovMigrate(govKey sdk.StoreKey, govKeeper types.GovKeeper) MigrateI {
	return &GovMigrate{
		govKey:    govKey,
		govKeeper: govKeeper,
	}
}

func (m *GovMigrate) Validate(ctx sdk.Context, _ Keeper, from, to sdk.AccAddress) error {
	votingParams := m.govKeeper.GetVotingParams(ctx)
	activeIter := m.govKeeper.ActiveProposalQueueIterator(ctx, ctx.BlockTime().Add(votingParams.VotingPeriod))
	defer activeIter.Close()
	for ; activeIter.Valid(); activeIter.Next() {
		//check vote
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(activeIter.Key())
		_, fromVoteFound := m.govKeeper.GetVote(ctx, proposalID, from)
		_, toVoteFound := m.govKeeper.GetVote(ctx, proposalID, to)
		if fromVoteFound && toVoteFound {
			return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, %s has voted proposal %d", to.String(), proposalID)
		}
	}
	return nil
}

func (m *GovMigrate) Execute(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error {
	govStore := ctx.KVStore(m.govKey)

	depositParams := m.govKeeper.GetDepositParams(ctx)
	inactiveIter := m.govKeeper.InactiveProposalQueueIterator(ctx, ctx.BlockTime().Add(depositParams.MaxDepositPeriod))
	defer inactiveIter.Close()
	for ; inactiveIter.Valid(); inactiveIter.Next() {
		proposalID, _ := govtypes.SplitInactiveProposalQueueKey(inactiveIter.Key())
		//migrate deposit
		if fromDeposit, fromFound := m.govKeeper.GetDeposit(ctx, proposalID, from); fromFound {
			amount := fromDeposit.Amount
			toDeposit, toFound := m.govKeeper.GetDeposit(ctx, proposalID, to)
			if toFound {
				amount = amount.Add(toDeposit.Amount...)
			}
			fromDeposit.Depositor = to.String()
			fromDeposit.Amount = amount
			govStore.Delete(govtypes.DepositKey(fromDeposit.ProposalId, from))
			govStore.Set(govtypes.DepositKey(fromDeposit.ProposalId, to), k.cdc.MustMarshal(&fromDeposit))
		}
	}

	votingParams := m.govKeeper.GetVotingParams(ctx)
	activeIter := m.govKeeper.ActiveProposalQueueIterator(ctx, ctx.BlockTime().Add(votingParams.VotingPeriod))
	defer activeIter.Close()
	for ; activeIter.Valid(); activeIter.Next() {
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(activeIter.Key())
		//migrate deposit
		if fromDeposit, depositFound := m.govKeeper.GetDeposit(ctx, proposalID, from); depositFound {
			amount := fromDeposit.Amount
			toDeposit, toFound := m.govKeeper.GetDeposit(ctx, proposalID, to)
			if toFound {
				amount = amount.Add(toDeposit.Amount...)
			}
			fromDeposit.Depositor = to.String()
			fromDeposit.Amount = amount
			govStore.Delete(govtypes.DepositKey(proposalID, from))
			govStore.Set(govtypes.DepositKey(proposalID, to), k.cdc.MustMarshal(&fromDeposit))
		}
		//migrate vote
		if fromVote, voteFound := m.govKeeper.GetVote(ctx, proposalID, from); voteFound {
			_, toFound := m.govKeeper.GetVote(ctx, proposalID, to)
			if toFound {
				return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, %s has voted proposal %d", to.String(), proposalID)
			}
			fromVote.Voter = to.String()
			govStore.Delete(govtypes.VoteKey(proposalID, from))
			govStore.Set(govtypes.VoteKey(proposalID, to), k.cdc.MustMarshal(&fromVote))
		}
	}
	return nil
}
