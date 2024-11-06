package staking_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/precompiles/staking"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func TestStakingWithdrawABI(t *testing.T) {
	withdrawABI := staking.NewWithdrawABI()

	require.Len(t, withdrawABI.Method.Inputs, 1)
	require.Len(t, withdrawABI.Method.Outputs, 1)

	require.Len(t, withdrawABI.Event.Inputs, 3)
}

func (suite *StakingPrecompileTestSuite) TestWithdraw() {
	testCases := []struct {
		name     string
		malleate func(val sdk.ValAddress, shares sdkmath.LegacyDec) (contract.WithdrawArgs, error)
		error    func(errArgs []string) string
		result   bool
	}{
		{
			name: "ok",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (contract.WithdrawArgs, error) {
				return contract.WithdrawArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (contract.WithdrawArgs, error) {
				newVal := val.String() + "1"
				return contract.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "failed validator not found",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (contract.WithdrawArgs, error) {
				newVal := sdk.ValAddress(suite.signer.Address().Bytes()).String()

				return contract.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("validator does not exist")
			},
			result: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			operator := suite.GetFirstValAddr()
			delAddr := suite.GetDelAddr()
			delAmt := suite.GetStakingBalance(delAddr.Bytes())

			res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
				Validator: operator.String(),
				Amount:    delAmt.BigInt(),
			})
			suite.Require().False(res.Failed(), res.VmError)

			suite.Commit()

			chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, delAddr.Bytes())
			suite.Require().True(chainBalances.IsZero(), chainBalances.String())
			totalBefore, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
			suite.Require().NoError(err)

			delegation := suite.GetDelegation(delAddr.Bytes(), operator)

			args, expectErr := tc.malleate(operator, delegation.Shares)

			res, _ = suite.WithError(expectErr).Withdraw(suite.Ctx, args)
			if tc.result {
				suite.Require().False(res.Failed(), res.VmError)

				totalAfter, err := suite.App.BankKeeper.TotalSupply(suite.Ctx, &banktypes.QueryTotalSupplyRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(totalAfter, totalBefore)

				withdrawABI := staking.NewWithdrawABI()
				reward, err := withdrawABI.UnpackOutput(res.Ret)
				suite.Require().NoError(err)

				chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, delAddr.Bytes())
				suite.Require().True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == withdrawABI.Event.ID.String() {
						suite.Require().Equal(contract.StakingAddress, log.Address)

						event, err := withdrawABI.UnpackEvent(log.ToEthereum())
						suite.Require().NoError(err)
						suite.Require().Equal(event.Sender, delAddr)
						suite.Require().Equal(event.Validator, operator.String())
						suite.Require().Equal(event.Reward.String(), chainBalances.AmountOf(fxtypes.DefaultDenom).BigInt().String())
						existLog = true
					}
				}
				suite.Require().True(existLog)

				existEvent := false
				for _, event := range suite.Ctx.EventManager().Events() {
					if event.Type == distritypes.EventTypeWithdrawRewards {
						for _, attr := range event.Attributes {
							if attr.Key == distritypes.AttributeKeyValidator {
								suite.Require().Equal(attr.Value, operator.String())
								existEvent = true
							}
							if attr.Key == sdk.AttributeKeyAmount {
								suite.Require().Equal(attr.Value, sdk.NewCoin(fxtypes.DefaultDenom, chainBalances.AmountOf(fxtypes.DefaultDenom)).String())
								existEvent = true
							}
							if attr.Key == distritypes.AttributeKeyDelegator {
								suite.Require().Equal(attr.Value, sdk.AccAddress(delAddr.Bytes()).String())
								existEvent = true
							}
						}
					}
				}
				suite.Require().True(existEvent)
			}
		})
	}
}
