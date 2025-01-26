package types

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	evmtypes "github.com/pundiai/fx-core/v8/x/evm/types"
)

var (
	DefaultCustomParamVotingPeriod    = time.Hour * 24 * 7  // Default period for deposits & voting  7 days
	DefaultEGFCustomParamVotingPeriod = time.Hour * 24 * 14 // Default egf period for deposits & voting  14 days

	EGFCustomParamDepositRatio     = sdkmath.LegacyNewDecWithPrec(1, 1) // 10%
	DefaultCustomParamDepositRatio = sdkmath.LegacyZeroDec()

	DefaultCustomParamQuorum40 = sdkmath.LegacyNewDecWithPrec(4, 1)
	DefaultCustomParamQuorum25 = sdkmath.LegacyNewDecWithPrec(25, 2) // 25%
)

var (
	FxSwitchParamsKey = []byte{0x92}
	CustomParamsKey   = []byte{0x93}
)

type InitGenesisCustomParams struct {
	MsgType string
	Params  CustomParams
}

func NewInitGenesisCustomParams(msgType string, params CustomParams) InitGenesisCustomParams {
	return InitGenesisCustomParams{
		MsgType: msgType,
		Params:  params,
	}
}

func NewCustomParams(depositRatio string, votingPeriod time.Duration, quorum string) CustomParams {
	return CustomParams{
		DepositRatio: depositRatio,
		VotingPeriod: &votingPeriod,
		Quorum:       quorum,
	}
}

func DefaultInitGenesisCustomParams() []InitGenesisCustomParams {
	var customParams []InitGenesisCustomParams
	customParams = append(customParams, newEGFCustomParams())
	customParams = append(customParams, newOtherCustomParams()...)
	return customParams
}

func newEGFCustomParams() InitGenesisCustomParams {
	return NewInitGenesisCustomParams(
		sdk.MsgTypeURL(&distributiontypes.MsgCommunityPoolSpend{}),
		NewCustomParams(EGFCustomParamDepositRatio.String(), DefaultEGFCustomParamVotingPeriod, DefaultCustomParamQuorum40.String()),
	)
}

func newOtherCustomParams() []InitGenesisCustomParams {
	defaultParams := NewCustomParams(DefaultCustomParamDepositRatio.String(), DefaultCustomParamVotingPeriod, DefaultCustomParamQuorum25.String())
	customMsgTypes := []string{
		// erc20 proposal
		sdk.MsgTypeURL(&erc20types.MsgRegisterNativeCoin{}),
		sdk.MsgTypeURL(&erc20types.MsgRegisterNativeERC20{}),
		sdk.MsgTypeURL(&erc20types.MsgRegisterBridgeToken{}),
		sdk.MsgTypeURL(&erc20types.MsgToggleTokenConversion{}),

		// evm proposal
		sdk.MsgTypeURL(&evmtypes.MsgCallContract{}),
	}
	params := make([]InitGenesisCustomParams, 0, len(customMsgTypes))
	for _, msgType := range customMsgTypes {
		params = append(params, NewInitGenesisCustomParams(msgType, defaultParams))
	}

	return params
}

func (p *SwitchParams) ValidateBasic() error {
	duplicate := make(map[string]bool)
	for _, precompile := range p.DisablePrecompiles {
		if duplicate[precompile] {
			return fmt.Errorf("duplicate precompile: %s", precompile)
		}
		duplicate[precompile] = true
	}

	duplicate = make(map[string]bool)
	for _, msgType := range p.DisableMsgTypes {
		if duplicate[msgType] {
			return fmt.Errorf("duplicate msg type: %s", msgType)
		}
		duplicate[msgType] = true
	}
	return nil
}
