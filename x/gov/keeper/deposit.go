package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal.
// Activates voting period when appropriate and returns true in that case, else returns false.
//
//nolint:gocyclo // copy from cosmos-sdk
func (keeper Keeper) AddDeposit(ctx context.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Coins) (bool, error) {
	// Checks to see if proposal exists
	proposal, err := keeper.Proposals.Get(ctx, proposalID)
	if err != nil {
		return false, err
	}

	// Check if proposal is still depositable
	if (proposal.Status != v1.StatusDepositPeriod) && (proposal.Status != v1.StatusVotingPeriod) {
		return false, govtypes.ErrInactiveProposal.Wrapf("%d", proposalID)
	}

	// Check coins to be deposited match the proposal's deposit params
	params, err := keeper.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	minDepositAmount := proposal.GetMinDepositFromParams(params)
	minDepositRatio, err := sdkmath.LegacyNewDecFromStr(params.GetMinDepositRatio())
	if err != nil {
		return false, err
	}

	// the deposit must only contain valid denoms (listed in the min deposit param)
	if err = keeper.validateDepositDenom(params, depositAmount); err != nil {
		return false, err
	}

	// If minDepositRatio is set, the deposit must be equal or greater than minDepositAmount*minDepositRatio
	// for at least one denom. If minDepositRatio is zero we skip this check.
	if !minDepositRatio.IsZero() {
		var (
			depositThresholdMet bool
			thresholds          []string
		)
		for _, minDep := range minDepositAmount {
			// calculate the threshold for this denom, and hold a list to later return a useful error message
			threshold := sdk.NewCoin(minDep.GetDenom(), minDep.Amount.ToLegacyDec().Mul(minDepositRatio).TruncateInt())
			thresholds = append(thresholds, threshold.String())

			found, deposit := depositAmount.Find(minDep.Denom)
			if !found { // if not found, continue, as we know the deposit contains at least 1 valid denom
				continue
			}

			// Once we know at least one threshold has been met, we can break. The deposit
			// might contain other denoms but we don't care.
			if deposit.IsGTE(threshold) {
				depositThresholdMet = true
				break
			}
		}

		// the threshold must be met with at least one denom, if not, return the list of minimum deposits
		if !depositThresholdMet {
			return false, errors.Wrapf(govtypes.ErrMinDepositTooSmall, "received %s but need at least one of the following: %s", depositAmount, strings.Join(thresholds, ","))
		}
	}

	// update the governance module's account coins pool
	err = keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, govtypes.ModuleName, depositAmount)
	if err != nil {
		return false, err
	}

	// Update proposal
	proposal.TotalDeposit = sdk.NewCoins(proposal.TotalDeposit...).Add(depositAmount...)
	err = keeper.SetProposal(ctx, proposal)
	if err != nil {
		return false, err
	}

	// get proposal minDepositAmount
	minDepositAmount, err = keeper.GetMinDepositAmountFromProposalMsgs(ctx, minDepositAmount, proposal)
	if err != nil {
		return false, err
	}

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == v1.StatusDepositPeriod && sdk.NewCoins(proposal.TotalDeposit...).IsAllGTE(minDepositAmount) {
		err = keeper.ActivateVotingPeriod(ctx, proposal)
		if err != nil {
			return false, err
		}

		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, err := keeper.Deposits.Get(ctx, collections.Join(proposalID, depositorAddr))
	switch {
	case err == nil:
		// deposit exists
		deposit.Amount = sdk.NewCoins(deposit.Amount...).Add(depositAmount...)
	case errors.IsOf(err, collections.ErrNotFound):
		// deposit doesn't exist
		deposit = v1.NewDeposit(proposalID, depositorAddr, depositAmount)
	default:
		// failed to get deposit
		return false, err
	}

	// called when deposit has been added to a proposal, however the proposal may not be active
	err = keeper.Hooks().AfterProposalDeposit(ctx, proposalID, depositorAddr)
	if err != nil {
		return false, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeProposalDeposit,
			sdk.NewAttribute(govtypes.AttributeKeyDepositor, depositorAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(govtypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	err = keeper.SetDeposit(ctx, deposit)
	if err != nil {
		return false, err
	}

	return activatedVotingPeriod, nil
}

func (keeper Keeper) GetMinDepositAmountFromProposalMsgs(ctx context.Context, defaultMinDeposit sdk.Coins, proposal v1.Proposal) (sdk.Coins, error) {
	message := proposal.GetMessages()
	totalCommunityPoolSpendAmount := sdk.NewCoins()
	egfMsgTypeURL := sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{})
	for _, msg := range message {
		if !strings.EqualFold(sdk.MsgTypeURL(msg), egfMsgTypeURL) {
			return defaultMinDeposit, nil
		}

		communityPoolSpendProposal := msg.GetCachedValue().(*distributiontypes.MsgCommunityPoolSpend)
		totalCommunityPoolSpendAmount = totalCommunityPoolSpendAmount.Add(communityPoolSpendProposal.Amount...)
	}
	// check egf params is set
	egfParams, err := keeper.CustomerParams.Get(ctx, egfMsgTypeURL)
	switch {
	case err == nil:
	case errors.IsOf(err, collections.ErrNotFound):
		return defaultMinDeposit, nil
	default:
		return nil, err
	}

	minDepositRatio, err := sdkmath.LegacyNewDecFromStr(egfParams.DepositRatio)
	if err != nil {
		return nil, err
	}
	if minDepositRatio.IsZero() {
		return defaultMinDeposit, nil
	}

	minDepositCoins := totalCommunityPoolSpendAmount
	for i := range minDepositCoins {
		minDepositCoins[i].Amount = sdkmath.LegacyNewDecFromInt(minDepositCoins[i].Amount).Mul(minDepositRatio).RoundInt()
	}

	// If the egf deposit amount is less than the default amount, the default amount is used
	if minDepositCoins.IsAllLT(defaultMinDeposit) {
		return defaultMinDeposit, nil
	}

	return minDepositCoins, nil
}

// validateInitialDeposit validates if initial deposit is greater than or equal to the minimum
// required at the time of proposal submission. This threshold amount is determined by
// the deposit parameters. Returns nil on success, error otherwise.
func (keeper Keeper) validateInitialDeposit(params v1.Params, initialDeposit sdk.Coins, expedited bool) error {
	if !initialDeposit.IsValid() || initialDeposit.IsAnyNegative() {
		return sdkerrors.ErrInvalidCoins.Wrap(initialDeposit.String())
	}

	minInitialDepositRatio, err := sdkmath.LegacyNewDecFromStr(params.MinInitialDepositRatio)
	if err != nil {
		return err
	}
	if minInitialDepositRatio.IsZero() {
		return nil
	}

	var minDepositCoins sdk.Coins
	if expedited {
		minDepositCoins = params.ExpeditedMinDeposit
	} else {
		minDepositCoins = params.MinDeposit
	}

	for i := range minDepositCoins {
		minDepositCoins[i].Amount = sdkmath.LegacyNewDecFromInt(minDepositCoins[i].Amount).Mul(minInitialDepositRatio).RoundInt()
	}
	if !initialDeposit.IsAllGTE(minDepositCoins) {
		return govtypes.ErrMinDepositTooSmall.Wrapf("was (%s), need (%s)", initialDeposit, minDepositCoins)
	}
	return nil
}

// validateDepositDenom validates if the deposit denom is accepted by the governance module.
func (keeper Keeper) validateDepositDenom(params v1.Params, depositAmount sdk.Coins) error {
	acceptedDenoms := make(map[string]bool, len(params.MinDeposit))
	denoms := make([]string, 0, len(params.MinDeposit))
	for _, coin := range params.MinDeposit {
		acceptedDenoms[coin.Denom] = true
		denoms = append(denoms, coin.Denom)
	}

	for _, coin := range depositAmount {
		if _, ok := acceptedDenoms[coin.Denom]; !ok {
			return govtypes.ErrInvalidDepositDenom.Wrapf("deposited %s, but gov accepts only the following denom(s): %v", depositAmount, denoms)
		}
	}

	return nil
}
