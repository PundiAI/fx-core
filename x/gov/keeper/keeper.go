package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/x/gov/types"
)

type Keeper struct {
	govkeeper.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	bankKeeper govtypes.BankKeeper
	sk         govtypes.StakingKeeper
}

func NewKeeper(bk govtypes.BankKeeper, sk govtypes.StakingKeeper, key sdk.StoreKey, gk govkeeper.Keeper) Keeper {
	return Keeper{
		storeKey:   key,
		bankKeeper: bk,
		sk:         sk,
		Keeper:     gk,
	}
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}

	// Check if proposal is still depositable
	if (proposal.Status != govtypes.StatusDepositPeriod) && (proposal.Status != govtypes.StatusVotingPeriod) {
		return false, sdkerrors.Wrapf(govtypes.ErrInactiveProposal, "%d", proposalID)
	}

	// update the governance module's account coins pool
	err := keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, govtypes.ModuleName, depositAmount)
	if err != nil {
		return false, err
	}

	// first deposit
	first := proposal.TotalDeposit.IsZero()
	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount...)
	keeper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false

	var minDeposit sdk.Coins
	isEGF := types.CommunityPoolSpendByRouter == proposal.ProposalRoute() && types.ProposalTypeCommunityPoolSpend == proposal.ProposalType()
	if isEGF {
		cpsp, ok := proposal.GetContent().(*distrtypes.CommunityPoolSpendProposal)
		if !ok {
			return false, sdkerrors.Wrapf(govtypes.ErrInvalidProposalType, "%d", proposalID)
		}
		minDeposit = SupportEGFProposalTotalDeposit(first, cpsp.Amount)
	} else {
		minDeposit = keeper.GetDepositParams(ctx).MinDeposit
	}
	if proposal.Status == govtypes.StatusDepositPeriod && proposal.TotalDeposit.IsAllGTE(minDeposit) {
		if ctx.BlockHeight() >= fxtypes.UpgradeTrigonometric2Block() && isEGF {
			keeper.EGFActivateVotingPeriod(ctx, proposal)
		} else {
			keeper.ActivateVotingPeriod(ctx, proposal)
		}

		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)

	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount...)
	} else {
		deposit = govtypes.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	// called when deposit has been added to a proposal, however the proposal may not be active
	keeper.AfterProposalDeposit(ctx, proposalID, depositorAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	keeper.SetDeposit(ctx, deposit)

	return activatedVotingPeriod, nil
}

func (keeper Keeper) EGFActivateVotingPeriod(ctx sdk.Context, proposal govtypes.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	twoWeek := time.Hour * 24 * 14
	if votingPeriod < twoWeek {
		votingPeriod = twoWeek
	}
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	proposal.Status = govtypes.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.ProposalId, proposal.DepositEndTime)
	keeper.InsertActiveProposalQueue(ctx, proposal.ProposalId, proposal.VotingEndTime)
}

func SupportEGFProposalTotalDeposit(first bool, claimCoin sdk.Coins) sdk.Coins {
	// minimum collateral amount for initializing EGF proposals
	if claimCoin.IsAllLTE(types.DepositProposalThreshold) && first {
		return types.InitialDeposit
	}
	initialDeposit := types.InitialDeposit
	for _, coin := range claimCoin {
		amount := coin.Amount.ToDec().Mul(types.ClaimRatio).TruncateInt()
		initialDeposit = initialDeposit.Add(sdk.NewCoin(coin.Denom, amount))
	}
	return initialDeposit
}
