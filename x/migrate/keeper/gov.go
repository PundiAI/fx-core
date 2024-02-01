package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v6/x/migrate/types"
)

type GovMigrate struct {
	govKey    storetypes.StoreKey
	govKeeper types.GovKeeper
}

func NewGovMigrate(govKey storetypes.StoreKey, govKeeper types.GovKeeper) MigrateI {
	return &GovMigrate{
		govKey:    govKey,
		govKeeper: govKeeper,
	}
}

func (m *GovMigrate) Validate(ctx sdk.Context, _ codec.BinaryCodec, from sdk.AccAddress, to common.Address) error {
	votingParams := m.govKeeper.GetVotingParams(ctx)
	activeIter := m.govKeeper.ActiveProposalQueueIterator(ctx, ctx.BlockTime().Add(*votingParams.VotingPeriod))
	defer activeIter.Close()
	for ; activeIter.Valid(); activeIter.Next() {
		// check vote
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(activeIter.Key())
		_, fromVoteFound := m.govKeeper.GetVote(ctx, proposalID, from)
		_, toVoteFound := m.govKeeper.GetVote(ctx, proposalID, to.Bytes())
		if fromVoteFound && toVoteFound {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "can not migrate, both from and to have voting proposal %d", proposalID)
		}
	}
	return nil
}

func (m *GovMigrate) Execute(ctx sdk.Context, cdc codec.BinaryCodec, from sdk.AccAddress, to common.Address) error {
	govStore := ctx.KVStore(m.govKey)
	events := make([]sdk.Event, 0, 10)

	depositParams := m.govKeeper.GetDepositParams(ctx)
	inactiveIter := m.govKeeper.InactiveProposalQueueIterator(ctx, ctx.BlockTime().Add(*depositParams.MaxDepositPeriod))
	defer inactiveIter.Close()
	for ; inactiveIter.Valid(); inactiveIter.Next() {
		proposalID, _ := govtypes.SplitInactiveProposalQueueKey(inactiveIter.Key())
		// migrate deposit
		if fromDeposit, fromFound := m.govKeeper.GetDeposit(ctx, proposalID, from); fromFound {
			amount := fromDeposit.Amount
			toDeposit, toFound := m.govKeeper.GetDeposit(ctx, proposalID, to.Bytes())
			if toFound {
				amount = sdk.NewCoins(amount...).Add(toDeposit.Amount...)
			}

			events = append(events,
				sdk.NewEvent(
					types.EventTypeMigrateGovDeposit,
					sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
					sdk.NewAttribute(sdk.AttributeKeyAmount, sdk.NewCoins(fromDeposit.Amount...).String()),
				),
			)

			fromDeposit.Depositor = sdk.AccAddress(to.Bytes()).String()
			fromDeposit.Amount = amount
			govStore.Delete(govtypes.DepositKey(fromDeposit.ProposalId, from))
			govStore.Set(govtypes.DepositKey(fromDeposit.ProposalId, to.Bytes()), cdc.MustMarshal(&fromDeposit))
		}
	}

	votingParams := m.govKeeper.GetVotingParams(ctx)
	activeIter := m.govKeeper.ActiveProposalQueueIterator(ctx, ctx.BlockTime().Add(*votingParams.VotingPeriod))
	defer activeIter.Close()
	for ; activeIter.Valid(); activeIter.Next() {
		proposalID, _ := govtypes.SplitActiveProposalQueueKey(activeIter.Key())
		// migrate deposit
		if fromDeposit, depositFound := m.govKeeper.GetDeposit(ctx, proposalID, from); depositFound {
			amount := fromDeposit.Amount
			toDeposit, toFound := m.govKeeper.GetDeposit(ctx, proposalID, to.Bytes())
			if toFound {
				amount = sdk.NewCoins(amount...).Add(toDeposit.Amount...)
			}

			events = append(events,
				sdk.NewEvent(
					types.EventTypeMigrateGovDeposit,
					sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
					sdk.NewAttribute(sdk.AttributeKeyAmount, sdk.NewCoins(fromDeposit.Amount...).String()),
				),
			)

			fromDeposit.Depositor = sdk.AccAddress(to.Bytes()).String()
			fromDeposit.Amount = amount
			govStore.Delete(govtypes.DepositKey(proposalID, from))
			govStore.Set(govtypes.DepositKey(proposalID, to.Bytes()), cdc.MustMarshal(&fromDeposit))
		}
		// migrate vote
		if fromVote, voteFound := m.govKeeper.GetVote(ctx, proposalID, from); voteFound {
			_, toFound := m.govKeeper.GetVote(ctx, proposalID, to.Bytes())
			if toFound {
				return errorsmod.Wrapf(types.ErrInvalidAddress, "can not migrate, both from and to have voting proposal %d", proposalID)
			}
			fromVote.Voter = sdk.AccAddress(to.Bytes()).String()
			govStore.Delete(govtypes.VoteKey(proposalID, from))
			govStore.Set(govtypes.VoteKey(proposalID, to.Bytes()), cdc.MustMarshal(&fromVote))

			var options govv1.WeightedVoteOptions = fromVote.Options
			// add events
			events = append(events,
				sdk.NewEvent(
					types.EventTypeMigrateGovVote,
					sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
					sdk.NewAttribute(types.AttributeKeyVoteOption, options.String()),
				),
			)
		}
	}

	if len(events) > 0 {
		ctx.EventManager().EmitEvents(events)
	}
	return nil
}
