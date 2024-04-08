package tests_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestParams() {
	testCases := []struct {
		name      string
		input     *types.Params
		expectErr bool
		expErrMsg string
	}{
		{
			name: "set GravityId cannpt be empty",
			input: &types.Params{
				GravityId:                         "",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "gravityId cannpt be empty",
		},
		{
			name: "set Invalid average block time",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  10,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "invalid average block time, too short for latency limitations",
		},
		{
			name: "set invalid average external block time",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          10,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "invalid average external block time, too short for latency limitations",
		},
		{
			name: "set Invalid signed window too short",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      1,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "invalid signed window too short",
		},
		{
			name: "set oracle delegate denom must FX",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 sdk.NewCoin("PX", sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "oracle delegate denom must FX",
		},
		{
			name: "set Invalid ibc transfer timeout too short",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          1,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "invalid ibc transfer timeout too short",
		},
		{
			name: "powet change percent too large",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDec(2),            // 200%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "powet change percent too large",
		},
		{
			name: "set slash factor too large",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDec(2),            // 200%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 types.DefaultBridgeCallTimeout,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "slash factor too large",
		},
		{
			name: "invalid bridge call timeout",
			input: &types.Params{
				GravityId:                         "fx-gravity-id",
				AverageBlockTime:                  7_000,
				AverageExternalBlockTime:          5_000,
				ExternalBatchTimeout:              12 * 3600 * 1000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
				OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 types.NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
				BridgeCallTimeout:                 1,
				Oracles:                           nil,
			},
			expectErr: true,
			expErrMsg: "invalid bridge call timeout",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			expected := suite.Keeper().GetParams(suite.ctx)
			err := suite.Keeper().SetParams(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				expected = *tc.input
				suite.Require().NoError(err)
			}
			params := suite.Keeper().GetParams(suite.ctx)
			suite.Require().Equal(expected, params)
		})
	}
}
