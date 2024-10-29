package precompile_test

import (
	"fmt"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distritypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/staking/precompile"
)

func TestStakingWithdrawABI(t *testing.T) {
	withdrawMethod := precompile.NewWithdrawMethod(nil)

	require.Len(t, withdrawMethod.Method.Inputs, 1)
	require.Len(t, withdrawMethod.Method.Outputs, 1)

	require.Len(t, withdrawMethod.Event.Inputs, 3)
}

func (suite *PrecompileTestSuite) TestWithdraw() {
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
		{
			name: "contract - ok",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (contract.WithdrawArgs, error) {
				return contract.WithdrawArgs{
					Validator: val.String(),
				}, nil
			},
			result: true,
		},
		{
			name: "contract - failed invalid validator address",
			malleate: func(val sdk.ValAddress, shares sdkmath.LegacyDec) (contract.WithdrawArgs, error) {
				newVal := val.String() + "1"
				return contract.WithdrawArgs{
					Validator: newVal,
				}, fmt.Errorf("invalid validator address: %s", newVal)
			},
			result: false,
		},
		{
			name: "contract - failed validator not found",
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
			delAmt := helpers.NewRandAmount()

			signer := suite.NewSigner()
			suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, delAmt))

			suite.WithContract(suite.stakingAddr)
			delAddr := signer.Address()
			if strings.HasPrefix(tc.name, "contract") {
				suite.WithContract(suite.stakingTestAddr)
				delAddr = suite.stakingTestAddr
				suite.MintToken(delAddr.Bytes(), sdk.NewCoin(fxtypes.DefaultDenom, delAmt))
			}

			suite.WithSigner(signer)
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

				abi := precompile.NewWithdrawABI()
				reward, err := abi.UnpackOutput(res.Ret)
				suite.Require().NoError(err)

				chainBalances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, delAddr.Bytes())
				suite.Require().True(chainBalances.AmountOf(fxtypes.DefaultDenom).Equal(sdkmath.NewIntFromBigInt(reward)), chainBalances.String())

				existLog := false
				for _, log := range res.Logs {
					if log.Topics[0] == abi.Event.ID.String() {
						suite.Require().Equal(log.Address, suite.stakingAddr.String())

						event, err := abi.UnpackEvent(log.ToEthereum())
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
