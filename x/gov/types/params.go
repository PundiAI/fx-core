package types

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
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
	initialDeposit := sdk.NewDecFromInt(claimAmount).Mul(ClaimRatio).TruncateInt()
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, initialDeposit))
}

func GetInitialDeposit() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))))
}

func CheckEGFProposalMsg(msgs []*codectypes.Any) (bool, sdk.Coins) {
	totalCommunityPoolSpendAmount := sdk.NewCoins()
	for _, msg := range msgs {
		// v1beta1 legacy MsgServer interface.from a legacy Content
		if strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&govv1.MsgExecLegacyContent{})) {
			legacyContent := msg.GetCachedValue().(*govv1.MsgExecLegacyContent)
			content := legacyContent.GetContent()
			if !strings.EqualFold(content.TypeUrl, "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal") {
				return false, nil
			}
			communityPoolSpendProposal := content.GetCachedValue().(*distributiontypes.CommunityPoolSpendProposal)
			totalCommunityPoolSpendAmount = totalCommunityPoolSpendAmount.Add(communityPoolSpendProposal.Amount...)
		} else {
			// TODO v1 MsgServer MsgCommunityPoolSpend pending
			// CommunityPoolSpendProposal is no msg type yet
			return false, nil
		}
	}
	return true, totalCommunityPoolSpendAmount
}

func CheckErc20ProposalMsg(msgs []*codectypes.Any) bool {
	for _, msg := range msgs {
		// v1beta1 legacy MsgServer interface.from a legacy Content
		if strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&govv1.MsgExecLegacyContent{})) {
			legacyContent := msg.GetCachedValue().(*govv1.MsgExecLegacyContent)
			content := legacyContent.GetContent()
			switch content.TypeUrl {
			case "/fx.erc20.v1.RegisterCoinProposal",
				"/fx.erc20.v1.RegisterERC20Proposal",
				"/fx.erc20.v1.ToggleTokenConversionProposal",
				"/fx.erc20.v1.UpdateDenomAliasProposal":
				return true
			default:
				return false
			}
		}
		if strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{})) ||
			strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&erc20types.MsgRegisterERC20{})) ||
			strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&erc20types.MsgToggleTokenConversion{})) ||
			strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&erc20types.MsgUpdateDenomAlias{})) {
			return true
		}
	}
	return false
}

// todo  add evm update contract
