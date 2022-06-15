package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/types"
)

const (
	// ProposalTypeCommunityPoolSpend defines the type for a CommunityPoolSpendProposal
	ProposalTypeCommunityPoolSpend = "CommunityPoolSpend"
	CommunityPoolSpendByRouter     = "distribution"
)

var (
	InitialDeposit           = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(1000).Mul(sdk.NewInt(1e18))))
	ClaimRatio               = sdk.NewDecWithPrec(1, 1)
	DepositProposalThreshold = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100000).Mul(sdk.NewInt(1e18))))
)
