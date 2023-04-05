package keeper_test

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/gov/types"
)

func (suite *KeeperTestSuite) TestFXParams() {
	egfVotingPeriod := time.Hour * 24 * 7
	evmVotingPeriod := time.Hour * 24 * 7

	ErrEgfVotingPeriod := time.Duration(0)
	ErrEvmVotingPeriod := time.Duration(0)
	testCases := []struct {
		name   string
		params types.Params
		expErr bool
		errStr string
	}{
		{
			name: "Valid Params",
			params: types.Params{
				MinInitialDeposit:   sdk.NewInt64Coin(fxtypes.DefaultDenom, 100),
				EgfDepositThreshold: sdk.NewInt64Coin(fxtypes.DefaultDenom, 50),
				ClaimRatio:          "0.5",
				Erc20Quorum:         "0.1",
				EvmQuorum:           "0.1",
				EgfVotingPeriod:     &egfVotingPeriod,
				EvmVotingPeriod:     &evmVotingPeriod,
			},
			expErr: false,
		},
		{
			name: "Invalid ClaimRatio",
			params: types.Params{
				MinInitialDeposit:   sdk.NewInt64Coin(fxtypes.DefaultDenom, 100),
				EgfDepositThreshold: sdk.NewInt64Coin(fxtypes.DefaultDenom, 50),
				ClaimRatio:          "-0.5",
				Erc20Quorum:         "0.1",
				EvmQuorum:           "0.1",
				EgfVotingPeriod:     &egfVotingPeriod,
				EvmVotingPeriod:     &evmVotingPeriod,
			},
			expErr: true,
			errStr: "claimRatio cannot be negative: -0.500000000000000000",
		},
		{
			name: "Invalid Erc20Quorum",
			params: types.Params{
				MinInitialDeposit:   sdk.NewInt64Coin(fxtypes.DefaultDenom, 100),
				EgfDepositThreshold: sdk.NewInt64Coin(fxtypes.DefaultDenom, 50),
				ClaimRatio:          "0.5",
				Erc20Quorum:         "1.2",
				EvmQuorum:           "0.1",
				EgfVotingPeriod:     &egfVotingPeriod,
				EvmVotingPeriod:     &evmVotingPeriod,
			},
			expErr: true,
			errStr: "erc20Quorum too large: 1.2",
		},
		{
			name: "Invalid EvmQuorum",
			params: types.Params{
				MinInitialDeposit:   sdk.NewInt64Coin(fxtypes.DefaultDenom, 100),
				EgfDepositThreshold: sdk.NewInt64Coin(fxtypes.DefaultDenom, 50),
				ClaimRatio:          "0.5",
				Erc20Quorum:         "0.1",
				EvmQuorum:           "-0.1",
				EgfVotingPeriod:     &egfVotingPeriod,
				EvmVotingPeriod:     &evmVotingPeriod,
			},
			expErr: true,
			errStr: "evmQuorum cannot be negative: -0.100000000000000000",
		},
		{
			name: "Invalid EgfVotingPeriod",
			params: types.Params{
				MinInitialDeposit:   sdk.NewInt64Coin(fxtypes.DefaultDenom, 100),
				EgfDepositThreshold: sdk.NewInt64Coin(fxtypes.DefaultDenom, 50),
				ClaimRatio:          "0.5",
				Erc20Quorum:         "0.1",
				EvmQuorum:           "0.1",
				EgfVotingPeriod:     &ErrEgfVotingPeriod,
				EvmVotingPeriod:     &evmVotingPeriod,
			},
			expErr: true,
			errStr: "egf voting period must be positive: 0s",
		},
		{
			name: "Invalid EvmVotingPeriod",
			params: types.Params{
				MinInitialDeposit:   sdk.NewInt64Coin(fxtypes.DefaultDenom, 100),
				EgfDepositThreshold: sdk.NewInt64Coin(fxtypes.DefaultDenom, 50),
				ClaimRatio:          "0.5",
				Erc20Quorum:         "0.1",
				EvmQuorum:           "0.1",
				EgfVotingPeriod:     &egfVotingPeriod,
				EvmVotingPeriod:     &ErrEvmVotingPeriod,
			},
			expErr: true,
			errStr: "evm voting period must be positive: 0s",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			err := tc.params.ValidateBasic()
			if tc.expErr {
				suite.Require().EqualValues(err.Error(), tc.errStr)
			} else {
				suite.Require().NoError(err)
				_, err = suite.MsgServer.UpdateParams(suite.ctx, &types.MsgUpdateParams{Authority: suite.govAcct, Params: tc.params})
				suite.Require().NoError(err)
				params := suite.app.GovKeeper.GetParams(suite.ctx)
				suite.Require().EqualValues(params.MinInitialDeposit, tc.params.MinInitialDeposit)
				suite.Require().EqualValues(params.EgfDepositThreshold, tc.params.EgfDepositThreshold)
				suite.Require().EqualValues(params.ClaimRatio, tc.params.ClaimRatio)
				suite.Require().EqualValues(params.Erc20Quorum, tc.params.Erc20Quorum)
				suite.Require().EqualValues(params.EvmQuorum, tc.params.EvmQuorum)
				suite.Require().EqualValues(params.EgfVotingPeriod, tc.params.EgfVotingPeriod)
				suite.Require().EqualValues(params.EgfVotingPeriod, tc.params.EgfVotingPeriod)
			}
		})
	}
}
