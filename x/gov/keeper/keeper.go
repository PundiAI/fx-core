package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
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

	minDeposit := keeper.NeedMinDeposit(ctx, proposal)
	if proposal.Status == govv1.StatusDepositPeriod && sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(minDeposit) {
		keeper.ActivateVotingPeriod(ctx, proposal)
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

func (keeper Keeper) ActivateVotingPeriod(ctx sdk.Context, proposal govv1.Proposal) {
	startTime := ctx.BlockHeader().Time
	proposal.VotingStartTime = &startTime
	votingPeriod := keeper.GetVotingPeriod(ctx, types.ExtractMsgTypeURL(proposal.Messages))
	endTime := proposal.VotingStartTime.Add(*votingPeriod)
	proposal.VotingEndTime = &endTime
	proposal.Status = govv1.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.Id, *proposal.DepositEndTime)
	keeper.InsertActiveProposalQueue(ctx, proposal.Id, *proposal.VotingEndTime)
}

func (keeper Keeper) NeedMinDeposit(ctx sdk.Context, proposal govv1.Proposal) sdk.Coins {
	var minDeposit sdk.Coins
	msgTypeURL := types.ExtractMsgTypeURL(proposal.Messages)
	isEGF, needDeposit := types.CheckEGFProposalMsg(proposal.Messages)
	if isEGF {
		minDeposit = keeper.EGFProposalMinDeposit(ctx, msgTypeURL, needDeposit)
	} else {
		minDeposit = keeper.GetMinDeposit(ctx, msgTypeURL)
	}
	return minDeposit
}

func (keeper Keeper) EGFProposalMinDeposit(ctx sdk.Context, msgType string, claimCoin sdk.Coins) sdk.Coins {
	egfParams := keeper.GetEGFParams(ctx)
	egfDepositThreshold := egfParams.EgfDepositThreshold
	claimRatio := egfParams.ClaimRatio
	claimAmount := claimCoin.AmountOf(fxtypes.DefaultDenom)
	if claimAmount.LTE(egfDepositThreshold.Amount) {
		return sdk.NewCoins(keeper.GetMinInitialDeposit(ctx, msgType))
	}
	ratio := sdk.MustNewDecFromStr(claimRatio)
	initialDeposit := sdk.NewDecFromInt(claimAmount).Mul(ratio).TruncateInt()
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initialDeposit))
}

func (keeper Keeper) InitFxGovParams(ctx sdk.Context) error {
	params := keeper.GetParams(ctx, "")
	erc20Params := types.Erc20ProposalParams(params.MinDeposit, params.MinInitialDeposit, params.VotingPeriod, types.DefaultErc20Quorum.String(), params.MaxDepositPeriod, params.Threshold, params.VetoThreshold)
	if err := keeper.SetAllParams(ctx, erc20Params); err != nil {
		return err
	}
	evmParams := types.EVMProposalParams(params.MinDeposit, params.MinInitialDeposit, &types.DefaultEvmVotingPeriod, types.DefaultEvmQuorum.String(), params.MaxDepositPeriod, params.Threshold, params.VetoThreshold)
	if err := keeper.SetAllParams(ctx, evmParams); err != nil {
		return err
	}
	egfParams := types.EGFProposalParams(params.MinDeposit, params.MinInitialDeposit, &types.DefaultEgfVotingPeriod, params.Quorum, params.MaxDepositPeriod, params.Threshold, params.VetoThreshold)
	if err := keeper.SetAllParams(ctx, egfParams); err != nil {
		return err
	}
	if err := keeper.SetEGFParams(ctx, types.DefaultEGFParams()); err != nil {
		return err
	}
	return nil
}
