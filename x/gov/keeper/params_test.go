package keeper_test

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
	"github.com/functionx/fx-core/v7/x/gov/types"
)

func (suite *KeeperTestSuite) TestParams() {
	timeDuration := time.Second * 60 * 60 * 24 * 7
	timeDurationErr := -time.Second * 60 * 60 * 24 * 7

	testCases := []struct {
		name   string
		params types.Params
		expErr bool
		errStr string
	}{
		{
			name: "Valid Params",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: false,
		},
		{
			name: "Invalid MinDeposit",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        nil,
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: true,
			errStr: "invalid minimum deposit: ",
		},
		{
			name: "Invalid MinInitialDeposit",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.Coin{},
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: true,
			errStr: "invalid minimum initial deposit: <nil>",
		},
		{
			name: "Nil MaxDepositPeriod",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  nil,
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: true,
			errStr: "maximum deposit period must not be nil: 0",
		},
		{
			name: "Negative Quorum",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "-0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: true,
			errStr: "quorom cannot be negative: -0.200000000000000000",
		},
		{
			name: "Negative Threshold",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "0.2",
				Threshold:         "-0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: true,
			errStr: "vote threshold must be positive: -0.500000000000000000",
		},
		{
			name: "Negative VetoThreshold",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "-0.334",
				VotingPeriod:      &timeDuration, // 1 week
			},
			expErr: true,
			errStr: "veto threshold must be positive: -0.334000000000000000",
		},
		{
			name: "Nil VotingPeriod",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDuration, // 1 week
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      nil,
			},
			expErr: true,
			errStr: "voting period must not be nil: 0",
		},
		{
			name: "Nil MaxDepositPeriod",
			params: types.Params{
				MsgType:           sdk.MsgTypeURL(&erc20types.MsgRegisterCoin{}),
				MinDeposit:        sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000))),
				MinInitialDeposit: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000)),
				MaxDepositPeriod:  &timeDurationErr, // 1 week
				Quorum:            "0.2",
				Threshold:         "0.5",
				VetoThreshold:     "0.334",
				VotingPeriod:      &timeDuration,
			},
			expErr: true,
			errStr: fmt.Sprintf("maximum deposit period must be positive: %v", timeDurationErr.String()),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			err := tc.params.ValidateBasic()
			if !tc.expErr {
				suite.Require().NoError(err, "expected no error, got %v", err)

				_, err := suite.MsgServer.UpdateParams(suite.ctx, types.NewMsgUpdateParams(suite.govAcct, tc.params))
				suite.Require().NoError(err, "expected no error, got %v", err)

				response, err := suite.queryClient.Params(suite.ctx, &types.QueryParamsRequest{MsgType: tc.params.MsgType})
				suite.Require().NoError(err, "expected no error, got %v", err)
				params := response.Params
				suite.Require().EqualValues(params.String(), tc.params.String())
			} else {
				suite.Require().EqualError(err, tc.errStr)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestEGFParams_ValidateBasic() {
	testCases := []struct {
		name          string
		p             *types.EGFParams
		expectedError error
	}{
		{
			name: "Valid EGFParams",
			p: &types.EGFParams{
				EgfDepositThreshold: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).MulRaw(1e18)), // 2%
				ClaimRatio:          "0.3",
			},
			expectedError: nil,
		},
		{
			name: "Invalid EGF Deposit Threshold",
			p: &types.EGFParams{
				EgfDepositThreshold: sdk.Coin{},
				ClaimRatio:          "0.5",
			},
			expectedError: fmt.Errorf("invalid Egf Deposit Threshold: <nil>"),
		},
		{
			name: "Invalid EGF Claim Ratio - Negative",
			p: &types.EGFParams{
				EgfDepositThreshold: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).MulRaw(1e18)),
				ClaimRatio:          "-0.5",
			},
			expectedError: fmt.Errorf("egf claim ratio cannot be negative: -0.500000000000000000"),
		},
		{
			name: "Invalid EGF Claim Ratio - Too Large",
			p: &types.EGFParams{
				EgfDepositThreshold: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).MulRaw(1e18)),
				ClaimRatio:          "2.0",
			},
			expectedError: fmt.Errorf("egf claim ratio too large: 2.000000000000000000"),
		},
		{
			name: "Invalid EGF Claim Ratio - Invalid String",
			p: &types.EGFParams{
				EgfDepositThreshold: sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).MulRaw(1e18)),
				ClaimRatio:          "invalid_ratio",
			},
			expectedError: fmt.Errorf("invalid egf claim ratio string: failed to set decimal string with base 10: invalid_ratio000000000000000000"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.p.ValidateBasic()
			if tc.expectedError == nil {
				suite.Require().NoError(err, "expected no error, got %v", err)

				_, err := suite.MsgServer.UpdateEGFParams(suite.ctx, types.NewMsgUpdateEGFParams(suite.govAcct, *tc.p))
				suite.Require().NoError(err, "expected no error, got %v", err)

				response, err := suite.queryClient.EGFParams(suite.ctx, &types.QueryEGFParamsRequest{})
				suite.Require().NoError(err, "expected no error, got %v", err)
				params := response.Params
				suite.Require().EqualValues(params.String(), tc.p.String())
			} else {
				suite.Require().EqualError(err, tc.expectedError.Error(), "expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
