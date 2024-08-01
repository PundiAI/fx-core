package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
)

type Keeper struct {
	*govkeeper.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey storetypes.StoreKey

	bankKeeper govtypes.BankKeeper
	sk         govtypes.StakingKeeper

	cdc codec.BinaryCodec

	authority string

	storeKeys map[string]*storetypes.KVStoreKey
}

func NewKeeper(bk govtypes.BankKeeper, sk govtypes.StakingKeeper, keys map[string]*storetypes.KVStoreKey, gk *govkeeper.Keeper, cdc codec.BinaryCodec, authority string) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	return &Keeper{
		storeKey:   keys[govtypes.StoreKey],
		bankKeeper: bk,
		sk:         sk,
		Keeper:     gk,
		cdc:        cdc,
		authority:  authority,
		storeKeys:  keys,
	}
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
	keeper.Hooks().AfterProposalDeposit(ctx, proposalID, depositorAddr)

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
	ratio := sdkmath.LegacyMustNewDecFromStr(claimRatio)
	initialDeposit := sdkmath.LegacyNewDecFromInt(claimAmount).Mul(ratio).TruncateInt()
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initialDeposit))
}

func (keeper Keeper) InitFxGovParams(ctx sdk.Context) error {
	params := keeper.GetFXParams(ctx, "")
	erc20Params := types.Erc20ProposalParams(params.MinDeposit, params.MinInitialDeposit, params.VotingPeriod,
		types.DefaultErc20Quorum.String(), params.MaxDepositPeriod, params.Threshold, params.VetoThreshold,
		params.MinInitialDepositRatio, params.BurnVoteQuorum, params.BurnProposalDepositPrevote, params.BurnVoteVeto)
	if err := keeper.SetAllParams(ctx, erc20Params); err != nil {
		return err
	}
	evmParams := types.EVMProposalParams(params.MinDeposit, params.MinInitialDeposit, &types.DefaultEvmVotingPeriod,
		types.DefaultEvmQuorum.String(), params.MaxDepositPeriod, params.Threshold, params.VetoThreshold,
		params.MinInitialDepositRatio, params.BurnVoteQuorum, params.BurnProposalDepositPrevote, params.BurnVoteVeto)
	if err := keeper.SetAllParams(ctx, evmParams); err != nil {
		return err
	}
	egfParams := types.EGFProposalParams(params.MinDeposit, params.MinInitialDeposit, &types.DefaultEgfVotingPeriod,
		params.Quorum, params.MaxDepositPeriod, params.Threshold, params.VetoThreshold,
		params.MinInitialDepositRatio, params.BurnVoteQuorum, params.BurnProposalDepositPrevote, params.BurnVoteVeto)
	if err := keeper.SetAllParams(ctx, egfParams); err != nil {
		return err
	}
	if err := keeper.SetEGFParams(ctx, types.DefaultEGFParams()); err != nil {
		return err
	}
	return nil
}

func (keeper Keeper) CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error {
	switchParams := keeper.GetSwitchParams(ctx)
	return CheckContractAddressIsDisabled(switchParams.DisablePrecompiles, contractAddress, methodId)
}

func CheckContractAddressIsDisabled(disabledPrecompiles []string, addr common.Address, methodId []byte) error {
	if len(disabledPrecompiles) == 0 {
		return nil
	}

	addrStr := strings.ToLower(addr.String())
	methodIdStr := hex.EncodeToString(methodId)
	addrMethodId := fmt.Sprintf("%s/%s", addrStr, methodIdStr)
	for _, disabledPrecompile := range disabledPrecompiles {
		disabledPrecompile = strings.ToLower(disabledPrecompile)
		if disabledPrecompile == addrStr {
			return errors.New("precompile address is disabled")
		}

		if disabledPrecompile == addrMethodId {
			return fmt.Errorf("precompile method %s is disabled", methodIdStr)
		}
	}

	return nil
}
