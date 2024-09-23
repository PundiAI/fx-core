package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/gov/types"
)

type Keeper struct {
	*govkeeper.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey storetypes.StoreKey

	authKeeper govtypes.AccountKeeper
	bankKeeper govtypes.BankKeeper
	sk         govtypes.StakingKeeper

	cdc codec.BinaryCodec

	authority string

	storeKeys map[string]*storetypes.KVStoreKey
}

func NewKeeper(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, sk govtypes.StakingKeeper, keys map[string]*storetypes.KVStoreKey, gk *govkeeper.Keeper, cdc codec.BinaryCodec, authority string) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	return &Keeper{
		storeKey:   keys[govtypes.StoreKey],
		authKeeper: ak,
		bankKeeper: bk,
		sk:         sk,
		Keeper:     gk,
		cdc:        cdc,
		authority:  authority,
		storeKeys:  keys,
	}
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
