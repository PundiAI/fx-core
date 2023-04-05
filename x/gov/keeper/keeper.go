package keeper

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/gov/types"
)

type Keeper struct {
	govkeeper.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey storetypes.StoreKey

	bankKeeper govtypes.BankKeeper
	sk         govtypes.StakingKeeper

	config types.Config

	cdc codec.BinaryCodec

	authority string
}

func NewKeeper(bk govtypes.BankKeeper, sk govtypes.StakingKeeper, key storetypes.StoreKey, gk govkeeper.Keeper, config types.Config, cdc codec.BinaryCodec, authority string) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	// If not set by app developer, set to default value.
	if config.MaxTitleLen == 0 {
		config.MaxTitleLen = types.DefaultConfig().MaxTitleLen
	}
	if config.MaxSummaryLen == 0 {
		config.MaxSummaryLen = types.DefaultConfig().MaxSummaryLen
	}
	if config.MaxMetadataLen == 0 {
		config.MaxMetadataLen = types.DefaultConfig().MaxMetadataLen
	}
	return Keeper{
		storeKey:   key,
		bankKeeper: bk,
		sk:         sk,
		Keeper:     gk,
		config:     config,
		cdc:        cdc,
		authority:  authority,
	}
}

func (keeper Keeper) Config() types.Config {
	return keeper.config
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return false, errorsmod.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}

	// Check if proposal is still depositable
	if (proposal.Status != govv1.StatusDepositPeriod) && (proposal.Status != govv1.StatusVotingPeriod) {
		return false, errorsmod.Wrapf(govtypes.ErrInactiveProposal, "%d", proposalID)
	}

	// update the governance module's account coins pool
	err := keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, govtypes.ModuleName, depositAmount)
	if err != nil {
		return false, err
	}

	// Update proposal
	proposal.TotalDeposit = sdk.NewCoins(proposal.TotalDeposit...).Add(depositAmount...)
	keeper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false

	var minDeposit sdk.Coins
	var isEVM bool
	isEGF, needDeposit := types.CheckEGFProposalMsg(proposal.Messages)
	if isEGF {
		minDeposit = keeper.EGFProposalMinDeposit(ctx, needDeposit)
	} else {
		isEVM = types.CheckEVMProposalMsg(proposal.Messages)
		minDeposit = keeper.GetDepositParams(ctx).MinDeposit
	}
	if proposal.Status == govv1.StatusDepositPeriod && sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(minDeposit) {
		fxParams := keeper.GetParams(ctx)
		if isEGF {
			keeper.afreshActivateVotingPeriod(ctx, fxParams.EgfVotingPeriod, proposal)
		} else if isEVM {
			keeper.afreshActivateVotingPeriod(ctx, fxParams.EvmVotingPeriod, proposal)
		} else {
			keeper.ActivateVotingPeriod(ctx, proposal)
		}

		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)

	if found {
		deposit.Amount = sdk.NewCoins(deposit.Amount...).Add(depositAmount...)
	} else {
		deposit = govv1.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	// called when deposit has been added to a proposal, however the proposal may not be active
	keeper.AfterProposalDeposit(ctx, proposalID, depositorAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		govtypes.EventTypeProposalDeposit,
		sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
		sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
	))

	keeper.SetDeposit(ctx, deposit)

	return activatedVotingPeriod, nil
}

func (keeper Keeper) EGFProposalMinDeposit(ctx sdk.Context, claimCoin sdk.Coins) sdk.Coins {
	fxParams := keeper.GetParams(ctx)
	egfDepositThreshold := fxParams.GetEgfDepositThreshold()
	claimRatio := fxParams.GetClaimRatio()
	claimAmount := claimCoin.AmountOf(fxtypes.DefaultDenom)
	if claimAmount.LTE(egfDepositThreshold.Amount) {
		return keeper.GetInitialDeposit(ctx)
	}
	ratio, _ := sdk.NewDecFromStr(claimRatio)
	initialDeposit := sdk.NewDecFromInt(claimAmount).Mul(ratio).TruncateInt()
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initialDeposit))
}

func (keeper Keeper) GetInitialDeposit(ctx sdk.Context) sdk.Coins {
	return sdk.NewCoins(keeper.GetParams(ctx).MinInitialDeposit)
}

func (keeper Keeper) afreshActivateVotingPeriod(ctx sdk.Context, newVotingPeriod *time.Duration, proposal govv1.Proposal) {
	blockTime := ctx.BlockHeader().Time
	proposal.VotingStartTime = &blockTime
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	if *votingPeriod < *newVotingPeriod {
		votingPeriod = newVotingPeriod
	}
	votingStartTime := *proposal.VotingStartTime
	add := votingStartTime.Add(*votingPeriod)
	proposal.VotingEndTime = &add
	proposal.Status = govv1.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.Id, *proposal.DepositEndTime)
	keeper.InsertActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)
}
