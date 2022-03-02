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

func (m *GovMigrate) Validate(_ sdk.Context, _ Keeper, _, _ sdk.AccAddress) error {
	return nil
}

func (m *GovMigrate) Execute(ctx sdk.Context, k Keeper, from, to sdk.AccAddress) error {
	govStore := ctx.KVStore(m.govKey)

	depositParams := m.govKeeper.GetDepositParams(ctx)
	inactiveIter := m.govKeeper.InactiveProposalQueueIterator(ctx, ctx.BlockTime().Add(depositParams.MaxDepositPeriod))
	defer inactiveIter.Close()
	for ; inactiveIter.Valid(); inactiveIter.Next() {
		proposalID, _ := govtypes.SplitInactiveProposalQueueKey(inactiveIter.Key())
		//migrate vote
		fromDeposit, fromFound := m.govKeeper.GetDeposit(ctx, proposalID, from)
		if fromFound {
			amount := fromDeposit.Amount
			toDeposit, toFound := m.govKeeper.GetDeposit(ctx, proposalID, to)
			if toFound {
				amount = amount.Add(toDeposit.Amount...)
			}
			fromDeposit.Depositor = to.String()
			fromDeposit.Amount = amount
			govStore.Delete(govtypes.DepositKey(fromDeposit.ProposalId, from))
			govStore.Set(govtypes.DepositKey(fromDeposit.ProposalId, to), k.cdc.MustMarshalBinaryBare(&fromDeposit))
		}
	}

	votingParams := m.govKeeper.GetVotingParams(ctx)
	activeIter := m.govKeeper.ActiveProposalQueueIterator(ctx, ctx.BlockTime().Add(votingParams.VotingPeriod))
	defer activeIter.Close()
	for ; activeIter.Valid(); activeIter.Next() {
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(activeIter.Key())
		//migrate deposit
		fromDeposit, depositFound := m.govKeeper.GetDeposit(ctx, proposalID, from)
		if depositFound {
			amount := fromDeposit.Amount
			toDeposit, toFound := m.govKeeper.GetDeposit(ctx, proposalID, to)
			if toFound {
				amount = amount.Add(toDeposit.Amount...)
			}
			fromDeposit.Depositor = to.String()
			fromDeposit.Amount = amount
			govStore.Delete(govtypes.DepositKey(proposalID, from))
			govStore.Set(govtypes.DepositKey(proposalID, to), k.cdc.MustMarshalBinaryBare(&fromDeposit))
		}
		//migrate vote
		fromVote, voteFound := m.govKeeper.GetVote(ctx, proposalID, from)
		if voteFound {
			_, toFound := m.govKeeper.GetVote(ctx, proposalID, to)
			//TODO error or remain unchanged or cover ???
			if toFound {
				return sdkerrors.Wrapf(types.ErrInvalidAddress, "can not migrate, %s has voted proposal %d", to.String(), proposalID)
			}
			fromVote.Voter = to.String()
			govStore.Delete(govtypes.VoteKey(proposalID, from))
			govStore.Set(govtypes.VoteKey(proposalID, to), k.cdc.MustMarshalBinaryBare(&fromVote))
		}
	}
	return nil
}
