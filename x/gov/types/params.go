package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

const (
	// ProposalTypeCommunityPoolSpend defines the type for a CommunityPoolSpendProposal
	ProposalTypeCommunityPoolSpend = "CommunityPoolSpend"
	CommunityPoolSpendByRouter     = "distribution"
)

func EGFProposalMinDeposit(claimCoin sdk.Coins) sdk.Coins {
	var (
		ClaimRatio          = sdk.NewDecWithPrec(1, 1) // 10%
		EGFDepositThreshold = sdkmath.NewInt(10_000).Mul(sdkmath.NewInt(1e18))
	)

	claimAmount := claimCoin.AmountOf(fxtypes.DefaultDenom)
	if claimAmount.LTE(EGFDepositThreshold) {
		return GetInitialDeposit()
	}
	initialDeposit := claimAmount.ToDec().Mul(ClaimRatio).TruncateInt()
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initialDeposit))
}

func GetInitialDeposit() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))))
}
