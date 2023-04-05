package types

import (
	"fmt"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	fxtypes "github.com/functionx/fx-core/v3/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	evmtypes "github.com/functionx/fx-core/v3/x/evm/types"
)

var (
	DefaultMinInitialDeposit   = sdkmath.NewInt(1000).Mul(sdkmath.NewInt(1e18))
	DefaultEgfDepositThreshold = sdkmath.NewInt(10_000).Mul(sdkmath.NewInt(1e18))
	DefaultClaimRatio          = sdk.NewDecWithPrec(1, 1)  // 10%
	DefaultErc20Quorum         = sdk.NewDecWithPrec(25, 2) // 25%
	DefaultEvmQuorum           = sdk.NewDecWithPrec(25, 2) // 25%
	DefaultEgfVotingPeriod     = time.Hour * 24 * 14       // Default egf period for deposits & voting  14 days
	DefaultEvmVotingPeriod     = time.Hour * 24 * 14       // Default evm period for deposits & voting  2 days
)

// NewParams creates a new Params instance with given values.
func NewParams(
	minInitialDeposit, egfDepositThreshold sdk.Coin, claimRatio, erc20Quorum, evmQuorum string, egfVotingPeriod, evmVotingPeriod *time.Duration,
) *Params {
	return &Params{
		MinInitialDeposit:   minInitialDeposit,
		EgfDepositThreshold: egfDepositThreshold,
		ClaimRatio:          claimRatio,
		Erc20Quorum:         erc20Quorum,
		EvmQuorum:           evmQuorum,
		EgfVotingPeriod:     egfVotingPeriod,
		EvmVotingPeriod:     evmVotingPeriod,
	}
}

// DefaultParams returns the default governance params
func DefaultParams() *Params {
	return NewParams(
		sdk.NewCoin(fxtypes.DefaultDenom, DefaultMinInitialDeposit),
		sdk.NewCoin(fxtypes.DefaultDenom, DefaultEgfDepositThreshold),
		DefaultClaimRatio.String(),
		DefaultErc20Quorum.String(),
		DefaultEvmQuorum.String(),
		&DefaultEgfVotingPeriod,
		&DefaultEvmVotingPeriod)
}

// ValidateBasic performs basic validation on governance parameters.
//
//gocyclo:ignore
func (p *Params) ValidateBasic() error {
	if !p.MinInitialDeposit.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", p.MinInitialDeposit.String())
	}
	if !p.EgfDepositThreshold.IsValid() {
		return fmt.Errorf("invalid minimum deposit: %s", p.EgfDepositThreshold.String())
	}
	claimRatio, err := sdk.NewDecFromStr(p.ClaimRatio)
	if err != nil {
		return fmt.Errorf("invalid claimRatio string: %w", err)
	}
	if claimRatio.IsNegative() {
		return fmt.Errorf("claimRatio cannot be negative: %s", claimRatio)
	}
	if claimRatio.GT(sdk.OneDec()) {
		return fmt.Errorf("claimRatio too large: %s", p.ClaimRatio)
	}
	erc20Quorum, err := sdk.NewDecFromStr(p.Erc20Quorum)
	if err != nil {
		return fmt.Errorf("invalid erc20Quorum string: %w", err)
	}
	if erc20Quorum.IsNegative() {
		return fmt.Errorf("erc20Quorum cannot be negative: %s", erc20Quorum)
	}
	if erc20Quorum.GT(sdk.OneDec()) {
		return fmt.Errorf("erc20Quorum too large: %s", p.Erc20Quorum)
	}
	evmQuorum, err := sdk.NewDecFromStr(p.EvmQuorum)
	if err != nil {
		return fmt.Errorf("invalid evmQuorum string: %w", err)
	}
	if evmQuorum.IsNegative() {
		return fmt.Errorf("evmQuorum cannot be negative: %s", evmQuorum)
	}
	if evmQuorum.GT(sdk.OneDec()) {
		return fmt.Errorf("evmQuorum too large: %s", p.EvmQuorum)
	}
	if p.EgfVotingPeriod == nil {
		return fmt.Errorf("egf Voting Period period must not be nil: %d", p.EgfVotingPeriod)
	}
	if p.EgfVotingPeriod.Seconds() <= 0 {
		return fmt.Errorf("egf voting period must be positive: %s", p.EgfVotingPeriod)
	}
	if p.EvmVotingPeriod == nil {
		return fmt.Errorf("evm Voting Period period must not be nil: %d", p.EvmVotingPeriod)
	}
	if p.EvmVotingPeriod.Seconds() <= 0 {
		return fmt.Errorf("evm voting period must be positive: %s", p.EvmVotingPeriod)
	}
	return nil
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

func CheckEVMProposalMsg(msgs []*codectypes.Any) bool {
	for _, msg := range msgs {
		if strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&govv1.MsgExecLegacyContent{})) {
			return false
		}
		if strings.EqualFold(msg.TypeUrl, sdk.MsgTypeURL(&evmtypes.MsgCallContract{})) {
			return true
		}
	}
	return false
}
